package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/aws"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	//Base Sensu plugin configs
	sensu.PluginConfig
	//AWS specific Sensu plugin configs
	aws.AWSPluginConfig
	//Additional configs for this check command
	Example string
	Verbose bool
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
			Path:      "example",
			Env:       "CHECK_EXAMPLE",
			Argument:  "example",
			Shorthand: "e",
			Default:   "",
			Usage:     "An example string configuration option",
			Value:     &plugin.Example,
		},
		&sensu.PluginConfigOption{
			Path:      "verbose",
			Argument:  "verbose",
			Shorthand: "v",
			Default:   false,
			Usage:     "Enable verbose output",
			Value:     &plugin.Verbose,
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
		fmt.Println("  Checking AWS Creds")
	}
	if state, err := plugin.CheckAWSCreds(); err != nil {
		return state, err
	}

	// Specific Argument Checking for this command
	if plugin.Verbose {
		fmt.Println("Checking Arguments")
	}

	if len(plugin.Example) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--example or CHECK_EXAMPLE environment variable is required")
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
}

func GetMetrics(c context.Context, api ServiceAPI, input *cloudwatch.ListMetricsInput) (*cloudwatch.ListMetricsOutput, error) {
	return api.ListMetrics(c, input)
}

// Note: Use ServiceAPI interface definition to make function testable with mock API testing pattern
// FIXME: replace s3 with correct service from AWS SDK
func checkFunction(client ServiceAPI) (int, error) {
	input := &cloudwatch.ListMetricsInput{}
	result, err := GetMetrics(context.TODO(), client, input)

	if err != nil {
		fmt.Println("Could not get metrics")
		return sensu.CheckStateCritical, nil
	}

	fmt.Println("Metrics:")
	numMetrics := 0

	for _, m := range result.Metrics {
		fmt.Println("   Metric Name: " + *m.MetricName)
		fmt.Println("   Namespace:   " + *m.Namespace)
		fmt.Println("   Dimensions:")
		for _, d := range m.Dimensions {
			fmt.Println("      " + *d.Name + ": " + *d.Value)
		}

		fmt.Println("")
		numMetrics++
	}

	fmt.Println("Found " + strconv.Itoa(numMetrics) + " metrics")

	return sensu.CheckStateOK, nil
}
