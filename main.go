package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	Namespace      string
	MetricName     string
	Verbose        bool
	DryRun         bool
	RecentlyActive bool
	MaxPages       int
	PeriodMinutes  int
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
			Default:   60,
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
)

func init() {
	//append common AWS options to options list
	options = append(options, plugin.GetAWSOpts()...)
}

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
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
	if len(plugin.Namespace) == 0 {
		plugin.DryRun = true
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
// FIXME: replace s3 functions with correct service functions from AWS SDK
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
// FIXME: replace s3 with correct service from AWS SDK
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
	return input, nil
}

func buildLabelBase(m types.Metric) string {
	s := strings.Split(*m.Namespace, "/")
	labelString := strings.ToLower(fmt.Sprintf("%v.%v.%v", s[0], s[1], *m.MetricName))
	return labelString
}

func buildGetMetricDataInput(m types.Metric) (*cloudwatch.GetMetricDataInput, error) {
	metricDataQueries := []types.MetricDataQuery{}
	s := strings.Split(*m.Namespace, "/")

	for _, stat := range []string{"SampleCount", "Average", "Maximum", "Minimum", "Sum"} {
		idString := strings.ToLower(fmt.Sprintf("%v_%v_%v_%v", s[0], s[1], *m.MetricName, stat))
		labelString := strings.ToLower(fmt.Sprintf("%v.%v", buildLabelBase(m), stat))
		dataQuery := types.MetricDataQuery{
			Id:    &idString,
			Label: &labelString,
			MetricStat: &types.MetricStat{
				Metric: &m,
				Period: aws.Int32(60 * int32(plugin.PeriodMinutes)),
				Stat:   aws.String(stat),
			},
		}
		metricDataQueries = append(metricDataQueries, dataQuery)
	}
	input := &cloudwatch.GetMetricDataInput{}
	input.EndTime = aws.Time(time.Unix(time.Now().Unix(), 0))
	input.StartTime = aws.Time(time.Unix(time.Now().Add(time.Duration(-plugin.PeriodMinutes)*time.Minute).Unix(), 0))
	input.MetricDataQueries = metricDataQueries
	return input, nil
}

func checkFunction(client ServiceAPI) (int, error) {
	numMetrics := 0
	numPages := 0
	dataMessages := []types.MessageData{}
	outputStrings := []string{}

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
			dimStrings := []string{}
			for _, d := range m.Dimensions {
				dimStrings = append(dimStrings, fmt.Sprintf(`%v="%v"`, *d.Name, *d.Value))
			}
			dimString := strings.Join(dimStrings, ",")
			if plugin.Verbose {
				fmt.Println("   Metric Name: " + *m.MetricName)
				fmt.Println("   Namespace:   " + *m.Namespace)
				fmt.Println("   Dimensions:")
				for _, d := range m.Dimensions {
					fmt.Println("      " + *d.Name + ": " + *d.Value)
				}

			}
			numMetrics++

			getMetricDataInput, err := buildGetMetricDataInput(m)
			if err != nil {
				fmt.Println("Could not build GetMetricsDataInput")
				return sensu.CheckStateCritical, nil
			}
			if plugin.DryRun {
				outputStrings = append(outputStrings, fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Dimensions:%v", buildLabelBase(m), *m.Namespace, *m.MetricName, dimString))
			} else {
				dataResult, err := GetMetricData(context.TODO(), client, getMetricDataInput)
				if err != nil {
					fmt.Printf("Could not get metrics: %v\n", err)
					return sensu.CheckStateCritical, nil
				}
				for _, d := range dataResult.MetricDataResults {
					if len(d.Timestamps) > 0 {
						outputStrings = append(outputStrings, fmt.Sprintf("# HELP %v Namespace:%v MetricName:%v Dimensions:%v", *d.Label, *m.Namespace, *m.MetricName, dimString))
						for i := range d.Timestamps {
							outputStrings = append(outputStrings, fmt.Sprintf("%v{%v} %v %v", *d.Label, dimString, d.Values[i], d.Timestamps[i].Unix()))
						}
						outputStrings = append(outputStrings, "")
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
			fmt.Printf("# GetMetricData:: Code: %v Message: %v\n", m.Code, m.Value)
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
