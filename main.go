package main

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sensu/sensu-cloudwatch-check/common"
	"github.com/sensu/sensu-cloudwatch-check/presets"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	sensuAWS "github.com/sensu/sensu-plugin-sdk/aws"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	//Base Sensu plugin configs
	sensu.PluginConfig
	//AWS specific Sensu plugin configs
	sensuAWS.AWSPluginConfig
	//Additional configs for this check command
	Namespace              string
	MetricName             string
	DimensionFilterStrings []string
	DimensionFilters       []types.DimensionFilter
	Verbose                bool
	DryRun                 bool
	RecentlyActive         bool
	MaxPages               int
	PeriodMinutes          int
	StatsList              []string
	PresetName             string
	Preset                 presets.ServicePreset
	// TODO: replace dryrun HELP with something useful
	ServiceExplorer bool
	// TODO: add support for json config
	ConfigString string
}

type MetricQueryMap struct {
	Id         string
	Label      string
	Namespace  string
	MetricName string
	Dimensions []types.Dimension
	Metric     *types.Metric
}

type MetricConfig struct {
	Measurement      string   `json:"measurement"`
	Namespace        string   `json:"namespace"`
	MetricName       string   `json:"metric-name"`
	Stat             string   `json:"stat"`
	DimensionFilters []string `json:"dimension-filters"`
}

var (
	//initialize Sensu plugin Config object
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-cloudwatch-check",
			Short:    "Sensu Cloudwatch Check",
			Keyspace: "sensu.io/plugins/sensu-cloudwatch-check/config",
		},
	}
	//initialize options list with custom options
	options = []*sensu.PluginConfigOption{
		/* TODO: Add support for json config
		&sensu.PluginConfigOption{
			Path:      "config",
			Argument:  "config",
			Shorthand: "c",
			Default:   "",
			Usage:     "Configuration JSON string",
			Value:     &plugin.ConfigString,
		},
		*/
		&sensu.PluginConfigOption{
			Path:     "recently-active",
			Argument: "recently-active",
			Default:  false,
			Usage:    "Only include metrics recently active in aprox last 3 hours",
			Value:    &plugin.RecentlyActive,
		},
		&sensu.PluginConfigOption{
			Path:      "namespace",
			Argument:  "namespace",
			Shorthand: "N",
			Default:   "",
			Usage:     "Cloudwatch Metric Namespace",
			Value:     &plugin.Namespace,
		},
		&sensu.PluginConfigOption{
			Path:      "dimension-filter",
			Argument:  "dimension-filter",
			Shorthand: "D",
			Default:   []string{},
			Usage:     `Comma separated list of AWS Cloudwatch Dimension Filters Ex: "Name, SecondName=SecondValue"`,
			Value:     &plugin.DimensionFilterStrings,
		},
		&sensu.PluginConfigOption{
			Path:      "stats",
			Argument:  "stats",
			Shorthand: "S",
			Default:   []string{"Average", "Sum", "SampleCount", "Maximum", "Minimum"},
			Usage:     `Comma separated list of AWS Cloudwatch Status Ex: "Average, Sum"`,
			Value:     &plugin.StatsList,
		},
		&sensu.PluginConfigOption{
			Path:      "metric",
			Argument:  "metric",
			Shorthand: "M",
			Default:   "",
			Usage:     "Cloudwatch Metric Name",
			Value:     &plugin.MetricName,
		},
		&sensu.PluginConfigOption{
			Path:      "preset",
			Argument:  "preset",
			Shorthand: "P",
			Default:   "None",
			Usage:     "Preset Name",
			Value:     &plugin.PresetName,
		},
		&sensu.PluginConfigOption{
			Path:      "max-pages",
			Argument:  "max-pages",
			Shorthand: "m",
			Default:   1,
			Usage:     "Maximum number of result pages",
			Value:     &plugin.MaxPages,
		},
		&sensu.PluginConfigOption{
			Path:      "period-minutes",
			Argument:  "period-minutes",
			Shorthand: "p",
			Default:   1,
			Usage:     "Period in minutes for metrics statistic calculation",
			Value:     &plugin.PeriodMinutes,
		},
		&sensu.PluginConfigOption{
			Path:      "verbose",
			Argument:  "verbose",
			Shorthand: "v",
			Default:   false,
			Usage:     "Enable verbose output",
			Value:     &plugin.Verbose,
		},
		&sensu.PluginConfigOption{
			Path:      "dry-run",
			Argument:  "dry-run",
			Shorthand: "n",
			Default:   false,
			Usage:     "Dryrun only list metrics, do not get metrics data",
			Value:     &plugin.DryRun,
		},
	}
	// Keyed list of metricConfig lookup table
	// TODO: setup json config
	metricConfig = make(map[string]MetricConfig)
)

func init() {
	//append common AWS options to options list
	options = append(options, plugin.GetAWSOpts()...)
}

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

/* TODO: setup json config
func buildConfigKey(m types.MetricStat) string {
	return strings.ToLower(*m.Metric.Namespace + "::" + *m.Metric.MetricName + "::" + *m.Stat)
	return ""
}

func (m MetricConfig) buildConfigKey() string {
	if len(m.Namespace) > 0 && len(m.MetricName) > 0 && len(m.Stat) > 0 {
		return strings.ToLower(m.Namespace + "::" + m.MetricName + "::" + m.Stat)
	} else {
		fmt.Println(len(m.Namespace), len(m.MetricName), len(m.Stat))
		return ""
	}
}

func buildMetricConfig(jsonBlob []byte) (map[string]MetricConfig, error) {
	objs := []MetricConfig{}
	metricConfig := make(map[string]MetricConfig)
	err := json.Unmarshal(jsonBlob, &objs)
	for _, o := range objs {
		if len(o.Namespace) == 0 {
			o.Namespace = plugin.Namespace
		}
		if len(o.MetricName) == 0 {
			o.MetricName = plugin.MetricName
		}
		if len(o.Stat) == 0 {
			o.Stat = "Average"
		}
		key := o.buildConfigKey()
		if len(key) > 0 {
			metricConfig[key] = o
		} else {
			return nil, err
		}
	}
	return metricConfig, err
}
*/

func checkArgs(event *v2.Event) (int, error) {
	// Check for valid AWS credentials
	if plugin.Verbose {
		fmt.Println("Checking AWS Creds")
	}
	if state, err := plugin.CheckAWSCreds(); err != nil {
		return state, err
	}

	// Specific Argument Checking for this command
	if plugin.Verbose {
		fmt.Println("Checking Arguments")
	}
	if len(plugin.DimensionFilterStrings) > 0 {
		dimensionFilters, err := common.BuildDimensionFilters(plugin.DimensionFilterStrings)
		if err != nil {
			return sensu.CheckStateWarning, err
		}
		plugin.DimensionFilters = dimensionFilters
	}

	if len(strings.TrimSpace(plugin.PresetName)) > 0 {
		if p, ok := presets.Presets[strings.TrimSpace(plugin.PresetName)]; ok {
			plugin.Preset = p
		} else {
			keys := reflect.ValueOf(presets.Presets).MapKeys()
			strArr := []string{}
			for _, key := range keys {
				str := fmt.Sprintf(" %v : %v\n", key.String(), presets.Presets[key.String()].GetDescription())
				strArr = append(strArr, str)
			}
			err := fmt.Errorf("Preset %v not defined\nChoose from:\n%v", plugin.PresetName, strings.Join(strArr, ""))
			return sensu.CheckStateWarning, err
		}
	} else {
		err := fmt.Errorf("No Preset selected")
		return sensu.CheckStateWarning, err
	}
	if plugin.Preset == nil {
		err := fmt.Errorf("No Preset selected")
		return sensu.CheckStateWarning, err
	}
	if len(plugin.PresetName) == 0 || plugin.PresetName == "None" {
		// If haven't selected a cloudwatch filter argument switch to dryrun to avoid pulling data for all metrics
		if len(plugin.Namespace) == 0 && len(plugin.MetricName) == 0 && !plugin.DryRun {
			return sensu.CheckStateWarning, fmt.Errorf("Must select at least one of: --namespace, --metric, or --dry-run")
		}
	}
	/* TODO: setup json config
	if len(plugin.ConfigString) > 0 {
		result, err := buildMetricConfig([]byte(plugin.ConfigString))
		if err != nil {
			return sensu.CheckStateWarning, err
		}
		metricConfig = result
	}
	*/
	return sensu.CheckStateOK, nil
}

func executeCheck(event *v2.Event) (int, error) {
	//Make sure plugin.CheckAwsCreds() worked as expected
	if plugin.AWSConfig == nil {
		return sensu.CheckStateCritical, fmt.Errorf("AWS Config undefined, something went wrong in processing AWS configuration information")
	}
	//Start AWS Service specific client
	client := cloudwatch.NewFromConfig(*plugin.AWSConfig)
	//Run business logic for check
	state, err := checkFunction(client)
	return state, err
}

//Create service interface to help with mock testing
type ServiceAPI interface {
	ListMetrics(ctx context.Context,
		params *cloudwatch.ListMetricsInput,
		optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error)
	GetMetricData(ctx context.Context,
		params *cloudwatch.GetMetricDataInput,
		optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
}

func GetMetricsList(c context.Context, api ServiceAPI, input *cloudwatch.ListMetricsInput) (*cloudwatch.ListMetricsOutput, error) {
	return api.ListMetrics(c, input)
}

func GetMetricData(c context.Context, api ServiceAPI, input *cloudwatch.GetMetricDataInput) (*cloudwatch.GetMetricDataOutput, error) {
	return api.GetMetricData(c, input)
}

// Note: Use ServiceAPI interface definition to make function testable with mock API testing pattern
func buildListMetricsInput(preset presets.ServicePreset) (*cloudwatch.ListMetricsInput, error) {
	input := &cloudwatch.ListMetricsInput{}
	if plugin.RecentlyActive {
		input.RecentlyActive = "PT3H"
	}
	if namespace := preset.GetNamespace(); len(namespace) > 0 {
		input.Namespace = &namespace
	}
	if metricName := preset.GetMetricName(); len(metricName) > 0 {
		input.MetricName = &metricName
	}
	if filters := preset.GetDimensionFilters(); len(filters) > 0 {
		input.Dimensions = filters
	}
	return input, nil
}

func buildGetMetricDataInput(metricDataQueries []types.MetricDataQuery) (*cloudwatch.GetMetricDataInput, error) {
	input := &cloudwatch.GetMetricDataInput{}
	input.EndTime = aws.Time(time.Unix(time.Now().Unix(), 0))
	input.StartTime = aws.Time(time.Unix(time.Now().Add(time.Duration(-plugin.PeriodMinutes)*time.Minute).Unix(), 0))
	input.MetricDataQueries = metricDataQueries
	return input, nil
}

func checkFunction(client ServiceAPI) (int, error) {
	var err error
	var metricDataQueries []types.MetricDataQuery
	numMetrics := 0
	numResults := 0
	numPages := 0
	dataMessages := []types.MessageData{}
	outputStrings := []string{}
	if plugin.PresetName == "None" {
		none := &presets.None{}
		none.Namespace = plugin.Namespace
		none.AddStats(plugin.StatsList)
		plugin.Preset = none
	}
	plugin.Preset.AddDimensionFilters(plugin.DimensionFilters)
	plugin.Preset.SetMetricName(plugin.MetricName)
	plugin.Preset.Init(false)
	//List Metrics result page loop
	for getList := true; getList && numPages < plugin.MaxPages; {
		getList = false
		input, err := buildListMetricsInput(plugin.Preset)
		if err != nil {
			fmt.Println("Could not create ListMetricsInput")
			return sensu.CheckStateCritical, nil
		}
		listResult, err := GetMetricsList(context.TODO(), client, input)

		if err != nil {
			fmt.Println("Could not get metrics list")
			return sensu.CheckStateCritical, nil
		}
		if listResult.NextToken != nil {
			getList = true
			numPages++
			input.NextToken = listResult.NextToken
		}

		plugin.Preset.AddMetrics(listResult.Metrics)

		numMetrics += len(listResult.Metrics)

	}

	numPages++
	if plugin.Verbose {
		fmt.Println("Found " + strconv.Itoa(numMetrics) + " metrics")
		fmt.Println("Result Pages " + strconv.Itoa(numPages))
		fmt.Println("")

	}

	metricDataQueries, err = plugin.Preset.BuildMetricDataQueries(int32(plugin.PeriodMinutes))
	if err != nil {
		fmt.Println("Could not build DataQuery")
		return sensu.CheckStateCritical, nil
	}
	if len(metricDataQueries) == 0 {
		fmt.Println("No metricDataQueries to process")
		return sensu.CheckStateWarning, nil
	}

	metricQueryMap := make(map[string]MetricQueryMap)
	unusedQueryMap := make(map[string]MetricQueryMap)

	for _, d := range metricDataQueries {
		idString := *d.Id
		qMap := MetricQueryMap{
			Id:         *d.Id,
			Label:      *d.Label,
			Metric:     d.MetricStat.Metric,
			MetricName: *d.MetricStat.Metric.MetricName,
			Namespace:  *d.MetricStat.Metric.Namespace,
			Dimensions: d.MetricStat.Metric.Dimensions,
		}
		metricQueryMap[idString] = qMap
		unusedQueryMap[idString] = qMap
	}

	//Prepare the GetMetricData loop
	i := 0
	for i < len(metricDataQueries) {
		//Pack up to 500 data queries into GetMetricData call
		j := i + 500
		if j > len(metricDataQueries) {
			j = len(metricDataQueries)
		}
		dataQuerySlice := metricDataQueries[i:j]
		getMetricDataInput, err := buildGetMetricDataInput(dataQuerySlice)
		if err != nil {
			fmt.Println("Could not build GetMetricsDataInput")
			return sensu.CheckStateCritical, nil
		}
		i = j + 1

		if plugin.DryRun {
			for _, d := range dataQuerySlice {
				q, ok := metricQueryMap[*d.Id]
				if !ok {
					fmt.Printf("Could not look up MetricQuery: %v\n", *d.Id)
					return sensu.CheckStateCritical, nil
				}
				delete(unusedQueryMap, *d.Id)
				outputStrings = append(outputStrings,
					fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Dimensions:%v",
						q.Label, q.Namespace, q.MetricName, common.DimString(q.Dimensions, plugin.AWSRegion)))

			}
		} else {
			dataResult, err := GetMetricData(context.TODO(), client, getMetricDataInput)
			if err != nil {
				fmt.Printf("Could not get metrics: %v\n", err)
				return sensu.CheckStateCritical, nil
			}
			if dataResult.NextToken != nil {
				fmt.Printf("GetMetricData result too long")
				return sensu.CheckStateCritical, nil
			}
			if len(dataResult.Messages) > 0 {
				fmt.Printf("GetMetricData has DataMessage: %v\n", dataResult.Messages)
				dataMessages = append(dataMessages, dataResult.Messages...)
			}
			for _, d := range dataResult.MetricDataResults {
				numResults++
				q, ok := metricQueryMap[*d.Id]
				if !ok {
					fmt.Printf("Could not look up MetricQuery: %v\n", *d.Id)
					return sensu.CheckStateCritical, nil
				}
				if len(d.Timestamps) > 0 {
					delete(unusedQueryMap, *d.Id)
					//outputStrings = append(outputStrings,
					//	fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Dimensions:%v",
					//		q.Label, q.Namespace, q.MetricName, common.DimString(q.Dimensions, plugin.AWSRegion)))
					for i := range d.Timestamps {
						outputStrings = append(outputStrings,
							fmt.Sprintf("%v{%v} %v %v",
								q.Label, common.DimString(q.Dimensions, plugin.AWSRegion), d.Values[i], d.Timestamps[i].Unix()))
					}
				}
			}
			if plugin.Verbose {
				fmt.Printf("   NextToken: %+v\n", dataResult.NextToken)
				fmt.Printf("   Messages: %+v\n", dataResult.Messages)
				fmt.Printf("   Data Results:\n")
				for _, d := range dataResult.MetricDataResults {
					fmt.Printf("     Id: %v\n", *d.Id)
					fmt.Printf("     Label: %+v\n", *d.Label)
					fmt.Printf("     StatusCode: %+v\n", d.StatusCode)
					fmt.Printf("     Timestamps: %+v\n", d.Timestamps)
					fmt.Printf("     Values: %+v\n", d.Values)
				}
				fmt.Println("")
			}
		}

	}
	warnFlag := false
	if numPages > plugin.MaxPages {
		fmt.Printf("# Warning: max allowed ListMetrics result pages (%v) exceeded, either filter via --namespace or --metric option or increase --max-pages value",
			plugin.MaxPages)
		warnFlag = true
	}
	if len(dataMessages) > 0 {
		fmt.Println("# Warning: Some calls to GetMetricData resulted in error messages")
		for _, m := range dataMessages {
			fmt.Printf("# GetMetricData:: Code: %v Message: %v\n", *m.Code, *m.Value)
		}
		warnFlag = true
	}
	if warnFlag {
		return sensu.CheckStateWarning, nil
	}
	for _, s := range outputStrings {
		fmt.Println(s)
	}
	if plugin.Verbose {
		fmt.Println("Summary:")
		fmt.Printf("  Number of Metrics: %v\n  MetricDataQueries: %v\n  QueryMaps: %v\n", numMetrics, len(metricDataQueries), len(metricQueryMap))
		fmt.Printf("  Number of MetricDataResults: %v\n", numResults)
		if len(unusedQueryMap) > 0 {
			fmt.Printf("  MetricDataQueries with no results:\n")

			for _, q := range unusedQueryMap {
				fmt.Printf("    Label: %v\n      Namespace:%v MetricName:%v Dimensions:%v\n",
					q.Label, q.Namespace, q.MetricName, common.DimString(q.Dimensions, plugin.AWSRegion))
			}
		}
	}
	return sensu.CheckStateOK, nil
}
