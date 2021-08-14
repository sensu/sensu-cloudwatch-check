package main

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// Create mockService Object to use in testing.
// FIXME: replace s3 specific items with correct AWS service items
var (
	nextToken = false
)

type mockService struct {
	statusCode      types.StatusCode
	includeMessages bool
}

// Create mockService Functions that match functions defined in ServiceAPI interface in main.go
func (m mockService) ListMetrics(ctx context.Context,
	params *cloudwatch.ListMetricsInput,
	optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error) {
	name := "test"
	namespace := "AWS/test"
	// Create a list of two dummy metrics
	metrics := []types.Metric{
		types.Metric{
			MetricName: &name,
			Namespace:  &namespace,
			Dimensions: []types.Dimension{
				types.Dimension{
					Name:  aws.String("test_name"),
					Value: aws.String("test_value"),
				},
			},
		},
	}
	output := &cloudwatch.ListMetricsOutput{
		Metrics: metrics,
	}
	if nextToken {
		output.NextToken = aws.String("yes")
		nextToken = false
	}
	return output, nil
}

func (m mockService) GetMetricData(ctx context.Context,
	params *cloudwatch.GetMetricDataInput,
	optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
	name := "test"
	namespace := "AWS/test"
	// Create a list of two dummy metrics
	results := []types.MetricDataResult{
		types.MetricDataResult{
			Id:         &name,
			Label:      &namespace,
			StatusCode: m.statusCode,
			Timestamps: []time.Time{
				time.Now(),
			},
			Values: []float64{
				0.0,
			},
		},
	}
	output := &cloudwatch.GetMetricDataOutput{
		MetricDataResults: results,
	}
	if m.includeMessages {
		output.Messages = []types.MessageData{
			types.MessageData{
				Code:  aws.String("400"),
				Value: aws.String("test message"),
			},
		}
	}
	return output, nil
}

func TestGetMetricData(t *testing.T) {
	cleanPluginValues()
	cases := []struct {
		client             mockService
		expectedStatusCode types.StatusCode
	}{ //start of array
		{ //start of struct
			client:             mockService{},
			expectedStatusCode: "Complete",
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			client.statusCode = tt.expectedStatusCode
			input := &cloudwatch.GetMetricDataInput{}
			output, err := client.GetMetricData(context.TODO(), input)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if len(output.MetricDataResults) == 0 {
				t.Fatalf("expected at least 1 data return")
			}
			for _, o := range output.MetricDataResults {
				if o.StatusCode != tt.expectedStatusCode {
					t.Errorf("expect status: %v, got %v", tt.expectedStatusCode, o.StatusCode)
				}
			}
		})

	}
	cleanPluginValues()
}
func cleanPluginValues() {
	plugin.Verbose = false
	plugin.RecentlyActive = false
	plugin.DryRun = false
	plugin.ConfigString = ""
	plugin.MetricName = ""
	plugin.Namespace = ""

}

func TestCheckArgs(t *testing.T) {
	cleanPluginValues()
	plugin.Verbose = true
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		if err == nil {
			t.Fatalf("expect error, got %v", err)
		}
		if state != 1 {
			t.Errorf("expect state: %v, got %v", 1, state)
		}
	})
	plugin.DryRun = true
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
		if state != 0 {
			t.Errorf("expect state: %v, got %v", 0, state)
		}
	})
	cleanPluginValues()
}

func TestCheckFunction(t *testing.T) {
	cleanPluginValues()
	plugin.RecentlyActive = true
	plugin.MetricName = "test"
	plugin.Namespace = "test"
	plugin.Verbose = true
	cases := []struct {
		client          mockService
		expectedState   int
		nextToken       bool
		maxPages        int
		includeMessages bool
	}{ //start of array
		{ //start of struct
			client:          mockService{},
			maxPages:        2,
			nextToken:       true,
			expectedState:   0,
			includeMessages: false,
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        1,
			nextToken:       true,
			expectedState:   1,
			includeMessages: false,
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        2,
			nextToken:       true,
			expectedState:   1,
			includeMessages: true,
		},
	}
	for i, tt := range cases {
		t.Run("CheckFunction Run: "+strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			client.includeMessages = tt.includeMessages
			nextToken = tt.nextToken
			plugin.MaxPages = tt.maxPages
			state, err := checkFunction(client)
			if err != nil {
				t.Fatalf("expect nil error, got %v", err)
			}
			if state != tt.expectedState {
				t.Errorf("expect state: %v, got %v", tt.expectedState, state)
			}
		})

	}
	cleanPluginValues()
}
