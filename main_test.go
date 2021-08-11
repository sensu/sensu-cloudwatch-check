package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// Create mockService Object to use in testing.
// FIXME: replace s3 specific items with correct AWS service items
type mockService struct {
	statusCode types.StatusCode
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
		},
	}
	output := &cloudwatch.ListMetricsOutput{
		Metrics: metrics,
	}

	return output, nil
}

func (m mockService) GetMetricData(ctx context.Context,
	params *cloudwatch.GetMetricDataInput,
	optFns ...func(cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
	name := "test"
	namespace := "AWS/test"
	// Create a list of two dummy metrics
	results := []types.MetricDataResult{
		types.MetricDataResult{
			Id:         &name,
			Label:      &namespace,
			StatusCode: m.statusCode,
		},
	}
	output := &cloudwatch.GetMetricDataOutput{
		MetricDataResults: results,
	}

	return output, nil
}

func TestGetMetricData(t *testing.T) {
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
}

func TestCheckFunction(t *testing.T) {
	cases := []struct {
		client             mockService
		tags               []string
		expectedState      int
		expectedStatusCode types.StatusCode
	}{ //start of array
		{ //start of struct
			client:             mockService{},
			expectedState:      0,
			expectedStatusCode: "Complete",
			tags: []string{
				"test_key",
			},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			state, err := checkFunction(client)
			client.statusCode = tt.expectedStatusCode
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if state != tt.expectedState {
				t.Errorf("expect state: %v, got %v", tt.expectedState, state)
			}
		})

	}
}
