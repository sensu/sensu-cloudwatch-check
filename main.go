package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/google/uuid"
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
	ConfigString           string
	DimensionFilterStrings []string
	DimensionFilters       []types.DimensionFilter
	Verbose                bool
	DryRun                 bool
	RecentlyActive         bool
	MaxPages               int
	PeriodMinutes          int
}

type MetricQueryMap struct {
	Id     *string
	Metric *types.Metric
}

var (
	//initialize Sensu plugin Config object
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "Sensu Cloudwatch Check",
			Short:    "Sensu Cloudwatch Check",
			Keyspace: "sensu.io/plugins/sensu-cloudwatch-check/config",
		},
	}
	//initialize options list with custom options
	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "config",
			Argument:  "config",
			Shorthand: "c",
			Default:   "",
			Usage:     "Configuration JSON string",
			Value:     &plugin.ConfigString,
		},
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
			Path:      "metric",
			Argument:  "metric",
			Shorthand: "M",
			Default:   "",
			Usage:     "Cloudwatch Metric Name",
			Value:     &plugin.MetricName,
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
	// Setup regexp for use with toSnakeCase
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func init() {
	//append common AWS options to options list
	options = append(options, plugin.GetAWSOpts()...)
}

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func buildDimensionFilters(input []string) ([]types.DimensionFilter, error) {
	output := []types.DimensionFilter{}
	for _, item := range input {
		segments := strings.Split(strings.TrimSpace(item), "=")
		if len(segments) < 1 || len(segments) > 2 {
			return nil, fmt.Errorf("Error parsing dimension filters")
		}
		filter := types.DimensionFilter{Name: &segments[0]}
		if len(segments) > 1 {
			filter.Value = &segments[1]
		}
		output = append(output, filter)
	}
	return output, nil
}

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
		dimensionFilters, err := buildDimensionFilters(plugin.DimensionFilterStrings)
		if err != nil {
			return sensu.CheckStateWarning, err
		}
		plugin.DimensionFilters = dimensionFilters
	}
	// If haven't selected a cloudwatch filter argument switch to dryrun to avoid pulling data for all metrics
	if len(plugin.Namespace) == 0 && len(plugin.MetricName) == 0 && !plugin.DryRun {
		return sensu.CheckStateWarning, fmt.Errorf("Must select at least one of: --namespace, --metric, or --dry-run")
	}
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
func buildListMetricsInput() (*cloudwatch.ListMetricsInput, error) {
	input := &cloudwatch.ListMetricsInput{}
	if plugin.RecentlyActive {
		input.RecentlyActive = "PT3H"
	}
	if len(plugin.Namespace) > 0 {
		input.Namespace = &plugin.Namespace
	}
	if len(plugin.MetricName) > 0 {
		input.MetricName = &plugin.MetricName
	}
	if len(plugin.DimensionFilters) > 0 {
		input.Dimensions = plugin.DimensionFilters
	}
	return input, nil
}

func buildLabelBase(m types.Metric) string {
	s := strings.Split(*m.Namespace, "/")
	labelString := toSnakeCase(fmt.Sprintf("%v.%v.%v", toSnakeCase(s[0]), toSnakeCase(s[1]), toSnakeCase(*m.MetricName)))
	return labelString
}

func buildMetricDataQueries(m types.Metric, stats []string) ([]types.MetricDataQuery, map[string]MetricQueryMap, error) {
	dataQueries := []types.MetricDataQuery{}
	queryMap := make(map[string]MetricQueryMap)

	for _, stat := range stats {
		id := uuid.New()
		idString := "aws_" + strings.ReplaceAll(id.String(), "-", "_")
		labelString := fmt.Sprintf("%v.%v", buildLabelBase(m), toSnakeCase(stat))
		dataQuery := types.MetricDataQuery{
			Id:    &idString,
			Label: &labelString,
			MetricStat: &types.MetricStat{
				Metric: &m,
				Period: aws.Int32(60 * int32(plugin.PeriodMinutes)),
				Stat:   aws.String(stat),
			},
		}
		queryMap[idString] = MetricQueryMap{
			Id:     &idString,
			Metric: &m,
		}
		dataQueries = append(dataQueries, dataQuery)
	}
	return dataQueries, queryMap, nil
}

func buildGetMetricDataInput(metricDataQueries []types.MetricDataQuery) (*cloudwatch.GetMetricDataInput, error) {
	input := &cloudwatch.GetMetricDataInput{}
	input.EndTime = aws.Time(time.Unix(time.Now().Unix(), 0))
	input.StartTime = aws.Time(time.Unix(time.Now().Add(time.Duration(-plugin.PeriodMinutes)*time.Minute).Unix(), 0))
	input.MetricDataQueries = metricDataQueries
	return input, nil
}

func dimString(m *types.Metric) string {
	if m == nil {
		return ""
	}
	dimStrings := []string{}
	if len(m.Dimensions) > 0 {
		for _, d := range m.Dimensions {
			dimStrings = append(dimStrings, fmt.Sprintf(`%v="%v"`, *d.Name, *d.Value))
		}
	}
	dimStr := strings.Join(dimStrings, ",")
	return dimStr
}

func checkFunction(client ServiceAPI) (int, error) {
	numMetrics := 0
	numPages := 0
	dataMessages := []types.MessageData{}
	metricDataQueries := []types.MetricDataQuery{}
	metricQueryMap := make(map[string]MetricQueryMap)

	outputStrings := []string{}
	stats := []string{"SampleCount", "Average", "Maximum", "Minimum", "Sum"}

	if plugin.Verbose {
		fmt.Println("Metrics:")
	}
	for getList := true; getList && numPages < plugin.MaxPages; {
		getList = false
		input, err := buildListMetricsInput()
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
		for _, m := range listResult.Metrics {
			if plugin.Verbose {
				fmt.Println("   Metric Name: " + *m.MetricName)
				fmt.Println("   Namespace:   " + *m.Namespace)
				fmt.Println("   Dimensions:")
				for _, d := range m.Dimensions {
					fmt.Println("      " + *d.Name + ": " + *d.Value)
				}

			}
			numMetrics++
			//Prepare data queries based on metric
			metricQueries, metricMap, err := buildMetricDataQueries(m, stats)
			if err != nil {
				fmt.Println("Could not build DataQuery")
				return sensu.CheckStateCritical, nil
			}
			// append data queries to global array
			metricDataQueries = append(metricDataQueries, metricQueries...)
			// add query metadata to global map query id map
			for k, v := range metricMap {
				metricQueryMap[k] = v
			}
		}
	}

	//Prepare the GetMetricData loop
	i := 0
	for i < len(metricDataQueries) {
		//Pack up to 500 data queries into GetMetricData call
		j := i + 500
		if j >= len(metricDataQueries) {
			j = len(metricDataQueries) - 1
		}
		dataQuerySlice := metricDataQueries[i:j]
		i = j + 1
		getMetricDataInput, err := buildGetMetricDataInput(dataQuerySlice)
		if err != nil {
			fmt.Println("Could not build GetMetricsDataInput")
			return sensu.CheckStateCritical, nil
		}

		if plugin.DryRun {
			for _, q := range dataQuerySlice {
				m := metricQueryMap[*q.Id].Metric
				if m != nil {
					outputStrings = append(outputStrings, fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Dimensions:%v", buildLabelBase(*m), *m.Namespace, *m.MetricName, dimString(m)))
				} else {
					fmt.Printf("Could not look up MetricQuery: %v\n", *q.Id)
					return sensu.CheckStateCritical, nil
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
			for _, d := range dataResult.MetricDataResults {
				m := metricQueryMap[*d.Id].Metric
				if m != nil {
					if len(d.Timestamps) > 0 {
						outputStrings = append(outputStrings, fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Dimensions:%v", *d.Label, *m.Namespace, *m.MetricName, dimString(m)))
						for i := range d.Timestamps {
							outputStrings = append(outputStrings, fmt.Sprintf("%v{%v} %v %v", *d.Label, dimString(m), d.Values[i], d.Timestamps[i].Unix()))
						}
						outputStrings = append(outputStrings, "")
					}
				} else {
					fmt.Printf("Could not look up MetricQuery: %v\n", *d.Id)
					return sensu.CheckStateCritical, nil
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
			if len(dataResult.Messages) > 0 {
				dataMessages = append(dataMessages, dataResult.Messages...)
			}
		}

	}
	numPages++
	if plugin.Verbose {
		fmt.Println("Found " + strconv.Itoa(numMetrics) + " metrics")
		fmt.Println("Result Pages " + strconv.Itoa(numPages))
		fmt.Println("")

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
	return sensu.CheckStateOK, nil
}
