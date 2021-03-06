package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sensu/sensu-cloudwatch-check/common"
	"github.com/sensu/sensu-cloudwatch-check/presets"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	sensuAWS "github.com/sensu/sensu-cloudwatch-check/aws"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-plugin-sdk/sensu/metric"
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
	ErrorOnMissing         bool
	DryRun                 bool
	RecentlyActive         bool
	MaxPages               int
	PeriodMinutes          int
	StatsList              []string
	PresetName             string
	Preset                 presets.PresetInterface
	OutputConfig           bool
	ConfigString           string
}

type MetricQueryMap struct {
	Id               string
	Label            string
	Namespace        string
	MetricName       string
	Dimensions       []types.Dimension
	Metric           *types.Metric
	MetricDataResult types.MetricDataResult
}

func (q MetricQueryMap) Points() ([]*v2.MetricPoint, error) {
	points := []*v2.MetricPoint{}
	metricTags := []*v2.MetricTag{}
	for _, d := range q.Dimensions {
		metricTags = append(metricTags, &v2.MetricTag{Name: *d.Name, Value: *d.Value})
	}

	for i := range q.MetricDataResult.Timestamps {

		point := v2.MetricPoint{Name: q.Label,
			Value:     q.MetricDataResult.Values[i],
			Timestamp: q.MetricDataResult.Timestamps[i].UnixNano() / 1000000,
			Tags:      metricTags,
		}
		points = append(points, &point)
	}
	return points, nil
}

func (q MetricQueryMap) Output(includeHelp bool, includeType bool, includeData bool) ([]string, error) {
	output := make([]string, 0)
	baseLabel := getBaseLabel(q.Label)
	if includeHelp {
		output = append(output,
			fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Region:%v",
				baseLabel, q.Namespace, q.MetricName, plugin.AWSConfig.Region))
	}
	if includeType {
		output = append(output,
			fmt.Sprintf("# TYPE %v gauge", baseLabel))
	}
	if includeData {
		for i := range q.MetricDataResult.Timestamps {
			output = append(output,
				fmt.Sprintf("%v{%v} %v %v",
					q.Label, common.DimString(q.Dimensions), q.MetricDataResult.Values[i], q.MetricDataResult.Timestamps[i].UnixNano()/1000000))
		}
	}
	return output, nil

}

func getBaseLabel(label string) string {
	last := strings.LastIndex(label, "_")
	if last > 0 {
		return label[0:last]
	}
	return label
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
	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "region",
			Env:       "AWS_REGION",
			Argument:  "region",
			Shorthand: "",
			Default:   "",
			Usage:     "AWS Region to use, (or set envvar AWS_REGION)",
			Value:     &plugin.AWSRegion,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "profile",
			Env:       "AWS_PROFILE",
			Argument:  "profile",
			Shorthand: "",
			Default:   "",
			Usage:     "AWS Credential Profile (for security use envvar AWS_PROFILE)",
			Secret:    false,
			Value:     &plugin.AWSProfile,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path:      "config-files",
			Env:       "",
			Argument:  "config-files",
			Shorthand: "",
			Default:   []string{},
			Usage:     "comma separated list of AWS config files",
			Secret:    false,
			Value:     &plugin.AWSConfigFiles,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path:      "credentials-files",
			Env:       "",
			Argument:  "credentials-files",
			Shorthand: "",
			Default:   []string{},
			Usage:     "comma separated list of AWS Credential files",
			Secret:    false,
			Value:     &plugin.AWSCredentialsFiles,
		},
		&sensu.PluginConfigOption[bool]{
			Value:     &plugin.OutputConfig,
			Path:      "output-config",
			Argument:  "output-config",
			Shorthand: "o",
			Default:   false,
			Usage:     "Output measurement configuration JSON string",
		},
		&sensu.PluginConfigOption[string]{
			Path:      "config",
			Argument:  "config",
			Env:       "CLOUDWATCH_CHECK_CONFIG",
			Shorthand: "c",
			Default:   "",
			Usage:     "The measurement configuration JSON string to use",
			Value:     &plugin.ConfigString,
		},
		&sensu.PluginConfigOption[bool]{
			Path:     "recently-active",
			Argument: "recently-active",
			Default:  false,
			Usage:    "Only include metrics recently active in aprox last 3 hours",
			Value:    &plugin.RecentlyActive,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "namespace",
			Argument:  "namespace",
			Env:       "CLOUDWATCH_CHECK_NAMESPACE",
			Shorthand: "N",
			Default:   "",
			Usage:     "Cloudwatch Metric Namespace",
			Value:     &plugin.Namespace,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path:      "dimension-filters",
			Argument:  "dimension-filters",
			Env:       "CLOUDWATCH_CHECK_DIMENSION_FILTERS",
			Shorthand: "D",
			Default:   []string{},
			Usage:     `Comma separated list of AWS Cloudwatch Dimension Filters Ex: "Name, SecondName=SecondValue"`,
			Value:     &plugin.DimensionFilterStrings,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path:      "stats",
			Argument:  "stats",
			Env:       "CLOUDWATCH_CHECK_STATS",
			Shorthand: "S",
			Default:   []string{"Average", "Sum", "SampleCount", "Maximum", "Minimum"},
			Usage:     `Comma separated list of AWS Cloudwatch Status Ex: "Average, Sum"`,
			Value:     &plugin.StatsList,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "metric-filter",
			Argument:  "metric-filter",
			Env:       "CLOUDWATCH_CHECK_METRIC_FILTER",
			Shorthand: "M",
			Default:   "",
			Usage:     "Cloudwatch Metric Filter, limit result to given Metric name",
			Value:     &plugin.MetricName,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "preset",
			Argument:  "preset",
			Env:       "CLOUDWATCH_CHECK_PRESET",
			Shorthand: "P",
			Default:   "None",
			Usage:     "The service preset to use",
			Value:     &plugin.PresetName,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "max-pages",
			Argument:  "max-pages",
			Env:       "CLOUDWATCH_CHECK_MAX_PAGES",
			Shorthand: "m",
			Default:   1,
			Usage:     "Maximum number of result pages. A zero value will disable the limit",
			Value:     &plugin.MaxPages,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "period-minutes",
			Argument:  "period-minutes",
			Env:       "CLOUDWATCH_CHECK_PERIOD_MINUTES",
			Shorthand: "p",
			Default:   1,
			Usage:     "Previous number of minutes to consider for metrics statistic calculation",
			Value:     &plugin.PeriodMinutes,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "verbose",
			Argument:  "verbose",
			Shorthand: "v",
			Default:   false,
			Usage:     "Enable verbose output",
			Value:     &plugin.Verbose,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "error-on-missing",
			Argument:  "error-on-missing",
			Shorthand: "",
			Default:   false,
			Env:       "CLOUDWATCH_CHECK_ERROR_ON_MISSING",
			Usage:     "Error if requested metrics configuration is missing a known metric from the AWS service metric list",
			Value:     &plugin.ErrorOnMissing,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "dry-run",
			Argument:  "dry-run",
			Shorthand: "n",
			Default:   false,
			Usage:     "Dryrun only list metrics, do not get metrics data",
			Value:     &plugin.DryRun,
		},
	}
)

func init() {
}

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(_ *v2.Event) (int, error) {

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
			strArr := make([]string, 0)
			for _, key := range keys {
				str := fmt.Sprintf(" %v : %v\n", key.String(), presets.Presets[key.String()].GetDescription())
				strArr = append(strArr, str)
			}
			err := fmt.Errorf("Preset %v not defined\nChoose from:\n%v", plugin.PresetName, strings.Join(strArr, ""))
			return sensu.CheckStateWarning, err
		}
	} else {
		err := fmt.Errorf("no preset selected")
		return sensu.CheckStateWarning, err
	}
	if plugin.Preset == nil {
		err := fmt.Errorf("no preset selected")
		return sensu.CheckStateWarning, err
	}
	if len(plugin.ConfigString) > 0 {
		if plugin.PresetName == "None" {
			plugin.PresetName = "Custom"
			p := presets.Preset{Description: "Custom Config"}

			err := p.SetMeasurementString(plugin.ConfigString)
			if err != nil {
				fmt.Println("Preset SetMeasurementString error")
				return sensu.CheckStateCritical, nil
			}
			plugin.Preset = &p
			err = plugin.Preset.BuildMeasurementConfig()
			if err != nil {
				fmt.Println("Preset BuildMeasurementConfig error")
				return sensu.CheckStateCritical, nil
			}

		} else {
			return sensu.CheckStateWarning, fmt.Errorf(`ConfigString not None Preset`)
		}
	}

	if len(plugin.PresetName) == 0 || plugin.PresetName == "None" {
		// If haven't selected a cloudwatch filter argument switch to dryrun to avoid pulling data for all metrics
		if len(plugin.ConfigString) == 0 && len(plugin.Namespace) == 0 && len(plugin.MetricName) == 0 && !plugin.DryRun {
			return sensu.CheckStateWarning, fmt.Errorf("must select at least one of: --config, --namespace, --metric, or --dry-run")
		}
	}
	if plugin.PresetName == "None" {
		none := &presets.None{}
		err := none.SetVerbose(plugin.Verbose)
		if err != nil {
			fmt.Println("Preset SetVerbose error")
			return sensu.CheckStateCritical, nil
		}
		err = none.SetPeriodMinutes(plugin.PeriodMinutes)
		if err != nil {
			fmt.Println("Preset SetPeriodMinutes error")
			return sensu.CheckStateCritical, nil
		}
		err = none.SetRegion(plugin.AWSRegion)
		if err != nil {
			fmt.Println("Preset SetRegion error")
			return sensu.CheckStateCritical, nil
		}
		err = none.Ready()
		if err != nil {
			fmt.Println("Preset Ready error")
			return sensu.CheckStateCritical, nil
		}
		none.Namespace = plugin.Namespace
		none.AddStats(plugin.StatsList)
		if len(plugin.ConfigString) > 0 {
			err = none.SetMeasurementString(plugin.ConfigString)
			if err != nil {
				fmt.Println("Preset SetMeasurementString error")
				return sensu.CheckStateCritical, nil
			}
			err = none.BuildMeasurementConfig()
			if err != nil {
				fmt.Println("Preset BuildMeasurementConfig error")
				return sensu.CheckStateCritical, nil
			}
		}
		plugin.Preset = none
	}
	// Use preset aws region if defined
	if region := plugin.Preset.GetRegion(); len(region) > 0 {
		plugin.AWSRegion = region
	}
	// Check for valid AWS credentials
	if plugin.Verbose {
		fmt.Println("Checking AWS Creds")
	}
	if state, err := plugin.CheckAWSCreds(); err != nil {
		return state, err
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(_ *v2.Event) (int, error) {
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

// ServiceAPI creates a service interface to help with mock testing
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
func buildListMetricsInput(preset presets.PresetInterface) (*cloudwatch.ListMetricsInput, error) {
	input := &cloudwatch.ListMetricsInput{}
	if plugin.RecentlyActive {
		input.RecentlyActive = "PT3H"
	}
	if namespace := preset.GetNamespace(); len(namespace) > 0 {
		input.Namespace = &namespace
	}
	if metricName := preset.GetMetricFilter(); len(metricName) > 0 {
		input.MetricName = &metricName
	}
	if filters := preset.GetDimensionFilters(); len(filters) > 0 {
		input.Dimensions = filters
	}
	return input, nil
}

func buildGetMetricDataInput(metricDataQueries []types.MetricDataQuery, periodMinutes int) (*cloudwatch.GetMetricDataInput, error) {
	input := &cloudwatch.GetMetricDataInput{}
	input.EndTime = aws.Time(time.Unix(time.Now().Unix(), 0))
	input.StartTime = aws.Time(time.Unix(time.Now().Add(time.Duration(-periodMinutes)*time.Minute).Unix(), 0))
	input.MetricDataQueries = metricDataQueries
	return input, nil
}

func getData(client ServiceAPI, metricDataQueries []types.MetricDataQuery, periodMinutes int) (int, error) {
	metricQueryMap := make(map[string]MetricQueryMap)
	unusedQueryMap := make(map[string]MetricQueryMap)
	dataMessages := make([]types.MessageData, 0)
	numResults := 0

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
	var results []*v2.MetricPoint
	//Prepare the GetMetricData loop
	i := 0
	for i < len(metricDataQueries) {
		//Pack up to 500 data queries into GetMetricData call
		j := i + 500
		if j > len(metricDataQueries) {
			j = len(metricDataQueries)
		}
		dataQuerySlice := metricDataQueries[i:j]
		getMetricDataInput, err := buildGetMetricDataInput(dataQuerySlice, periodMinutes)
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
				metricPoints, err := q.Points()
				if err == nil {
					results = append(results, metricPoints...)
				}

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
				q.MetricDataResult = d
				if !ok {
					fmt.Printf("Could not look up MetricQuery: %v\n", *d.Id)
					return sensu.CheckStateCritical, nil
				}
				if len(d.Timestamps) > 0 {
					delete(unusedQueryMap, *d.Id)
					metricPoints, err := q.Points()
					if err == nil {
						results = append(results, metricPoints...)
					}
				}
			}
		}

	}
	if plugin.Verbose {
		fmt.Println("\nExecution Summary:")
		fmt.Printf("  MetricDataQueries: %v\n", len(metricDataQueries))
		fmt.Printf("  Number of MetricDataResults: %v\n", numResults)
		if len(unusedQueryMap) > 0 {
			fmt.Printf("  MetricDataQueries with no results:\n")

			for _, q := range unusedQueryMap {
				fmt.Printf("    Label: %v\n      Namespace:%v MetricName:%v Region:%v Dimensions:%v\n",
					q.Label, q.Namespace, q.MetricName, plugin.AWSConfig.Region, common.DimString(q.Dimensions))
			}
		}
		fmt.Println("")
		fmt.Println("Normal Output:")
	}
	warnFlag := false
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
	if len(results) > 0 {
		writer := bufio.NewWriter(os.Stdout)
		err := metric.Points(results).ToProm(writer)
		if err != nil {
			return sensu.CheckStateCritical, err
		}
		writer.Flush()
	}
	return sensu.CheckStateOK, nil

}

func checkFunction(client ServiceAPI) (int, error) {
	var err error
	var metricDataQueries []types.MetricDataQuery
	numMetrics := 0
	numPages := 0
	err = plugin.Preset.AddDimensionFilters(plugin.DimensionFilters)
	if err != nil {
		fmt.Println("Preset AddDimensionFilters error")
		return sensu.CheckStateCritical, nil
	}
	err = plugin.Preset.SetMetricFilter(plugin.MetricName)
	if err != nil {
		fmt.Println("Preset SetMetricFilter error")
		return sensu.CheckStateCritical, nil
	}
	err = plugin.Preset.SetVerbose(plugin.Verbose)
	if err != nil {
		fmt.Println("Preset SetVerbose error")
		return sensu.CheckStateCritical, nil
	}
	err = plugin.Preset.SetErrorOnMissing(plugin.ErrorOnMissing)
	if err != nil {
		fmt.Println("Preset SetErrorOnMissing error")
		return sensu.CheckStateCritical, nil
	}
	err = plugin.Preset.Ready()
	if err != nil {
		fmt.Println("Preset Ready error")
		return sensu.CheckStateCritical, nil
	}
	//List Metrics result page loop

	for getList := true; getList && (plugin.MaxPages == 0 || numPages < plugin.MaxPages); {
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

		err = plugin.Preset.AddMetrics(listResult.Metrics)
		if err != nil {
			fmt.Println("Preset AddMetrics error")
			return sensu.CheckStateCritical, nil
		}

		numMetrics += len(listResult.Metrics)

	}
	numPages++
	if plugin.Verbose {
		fmt.Println("Found " + strconv.Itoa(numMetrics) + " metrics")
		fmt.Println("Result Pages " + strconv.Itoa(numPages))
		fmt.Println("")
	}
	if plugin.OutputConfig {
		if output, err := plugin.Preset.GetMeasurementString(true); err != nil {
			return sensu.CheckStateCritical, nil
		} else {
			if plugin.Verbose {
				fmt.Println("Measurement Configuration:")
			}
			fmt.Println(output)
		}
	} else {
		periodMinutes := plugin.PeriodMinutes
		if p := plugin.Preset.GetPeriodMinutes(); p > 0 {
			periodMinutes = p
		}
		metricDataQueries, err = plugin.Preset.BuildMetricDataQueries(int32(periodMinutes))
		if err != nil {
			fmt.Println("Could not build DataQuery")
			return sensu.CheckStateCritical, nil
		}
		if len(metricDataQueries) == 0 {
			fmt.Println("No metricDataQueries to process")
			return sensu.CheckStateWarning, nil
		}
		if state, err := getData(client, metricDataQueries, periodMinutes); state != sensu.CheckStateOK {
			return state, err
		}
		// Outputting Metrics
		if plugin.MaxPages > 0 && numPages > plugin.MaxPages {
			fmt.Printf("\n# Warning: max allowed ListMetrics result pages (%v) exceeded, either filter via --namespace or --metric option or increase --max-pages value\n",
				plugin.MaxPages)
			return sensu.CheckStateWarning, nil
		}

	}
	return sensu.CheckStateOK, nil
}
