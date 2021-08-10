package main

import (
	"context"
	"fmt"

	// FIXME: Replace s3 with the correct aws service subpackage from aws-sdk-go-v2
	"github.com/aws/aws-sdk-go-v2/service/s3"

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
	if state, err := plugin.CheckAWSCreds(); err != nil {
		return state, err
	}

	// Specific Argument Checking for this command
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
	// FIXME: replace s3 with correct service from AWS SDK
	client := s3.NewFromConfig(*plugin.AWSConfig)
	//Run business logic for check
	state, err := checkFunction(client)
	return state, err
}

//Create service interface to help with mock testing
// FIXME: replace s3 functions with correct service functions from AWS SDK
type ServiceAPI interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
}

// Note: Use ServiceAPI interface definition to make function testable with mock API testing pattern
// FIXME: replace s3 with correct service from AWS SDK
func checkFunction(client ServiceAPI) (int, error) {
	inputs := &s3.ListBucketsInput{}
	output, err := client.ListBuckets(context.Background(), inputs)
	if err != nil {
		return sensu.CheckStateCritical, err
	}
	if output != nil && output.Buckets != nil && len(output.Buckets) > 0 {
		for _, bucket := range output.Buckets {
			bucketInput := &s3.GetBucketTaggingInput{Bucket: bucket.Name}
			bucketOutput, err := client.GetBucketTagging(context.Background(), bucketInput)
			if err != nil {
				continue
			}
			if bucketOutput == nil {
				continue
			}

		}
	}
	return sensu.CheckStateOK, nil
}
