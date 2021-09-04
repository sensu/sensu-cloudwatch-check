package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/sensu/sensu-cloudwatch-check/common"
	"github.com/sensu/sensu-cloudwatch-check/presets"
	"github.com/stretchr/testify/assert"
)

// Create mockService Object to use in testing.
// FIXME: replace s3 specific items with correct AWS service items
var (
	nextToken   = false
	enableQuiet = false
)

type mockService struct {
	statusCode      types.StatusCode
	dataResultId    string
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
	results := []types.MetricDataResult{}
	for _, d := range params.MetricDataQueries {
		// Create a list of two dummy metrics
		result := types.MetricDataResult{
			Id:         d.Id,
			Label:      d.Label,
			StatusCode: m.statusCode,
			Timestamps: []time.Time{
				time.Now(),
			},
			Values: []float64{
				0.0,
			},
		}
		results = append(results, result)
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
			types.MessageData{
				Code:  aws.String("400"),
				Value: aws.String("test message"),
			},
		}
	}
	return output, nil
}

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	if enableQuiet {
		os.Stdout = null
		os.Stderr = null
		log.SetOutput(null)
	}
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}

func cleanPluginValues() {
	config := aws.Config{}

	plugin.StatsList = []string{}
	plugin.Verbose = false
	plugin.RecentlyActive = false
	plugin.DryRun = false
	plugin.ConfigString = ""
	plugin.MetricName = ""
	plugin.Namespace = ""
	plugin.MaxPages = 0
	plugin.PeriodMinutes = 0
	plugin.PresetName = ""
	plugin.AWSConfig = &config
}

func TestQueryMapOutput(t *testing.T) {
	defer quiet()()
	cleanPluginValues()
	plugin.AWSConfig.Region = "us-east-1"
	cases := []struct {
		expectedId    string
		expectedLabel string
		namespace     string
		metricName    string
	}{ //start of array
		{ //start of struct
			expectedId:    "test",
			expectedLabel: "test_label",
			namespace:     "AWS/Test",
			metricName:    "TestMetric",
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert := assert.New(t)
			q := MetricQueryMap{
				Id:         tt.expectedId,
				Label:      tt.expectedLabel,
				Namespace:  tt.namespace,
				MetricName: tt.metricName,
				Dimensions: []types.Dimension{},
				Metric: &types.Metric{
					Namespace:  aws.String(tt.namespace),
					MetricName: aws.String(tt.metricName),
				},
			}
			output, err := q.Output(true, true, true)
			for _, o := range output {
				fmt.Println(o)
			}
			assert.NoError(err)
		})

	}
	cleanPluginValues()
}
func TestGetMetricData(t *testing.T) {
	defer quiet()()
	cleanPluginValues()

	cases := []struct {
		client             mockService
		expectedStatusCode types.StatusCode
		expectedId         string
	}{ //start of array
		{ //start of struct
			client:             mockService{},
			expectedStatusCode: "Complete",
			expectedId:         "test",
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			client.statusCode = tt.expectedStatusCode
			client.dataResultId = tt.expectedId
			input := &cloudwatch.GetMetricDataInput{}
			output, err := client.GetMetricData(context.TODO(), input)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if len(output.MetricDataResults) != 0 {
				t.Fatalf("expected 0 data results")
			}
			id := "test"
			input = &cloudwatch.GetMetricDataInput{
				MetricDataQueries: []types.MetricDataQuery{
					types.MetricDataQuery{
						Id: &id,
					},
				},
			}
			output, err = client.GetMetricData(context.TODO(), input)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if len(output.MetricDataResults) != 1 {
				t.Fatalf("expected 1 data result")
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

/* TODO: setup json config
func TestBuildMetricConfig(t *testing.T) {
	assert := assert.New(t)
	jsonStr := []byte(`[{"measaurement": "aws.test.measurement", "namespace" : "AWS/TEST", "metric-name": "Test", "stats" : ["Average","SampleCount"] }]`)
	key := "aws/test::test::average"
	result, err := buildMetricConfig(jsonStr)
	assert.NoError(err)
	assert.NotNil(result[key])
}
*/
func TestToSnakeCase(t *testing.T) {
	defer quiet()()
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"already_snake", "already_snake"},
		{"A", "a"},
		{"AA", "aa"},
		{"AaAa", "aa_aa"},
		{"HTTPRequest", "http_request"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}
	for _, test := range tests {
		have := common.ToSnakeCase(test.input)
		if have != test.want {
			t.Errorf("input=%q:\nhave: %q\nwant: %q", test.input, have, test.want)
		}
	}
}
func TestCheckArgs(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	cleanPluginValues()
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		assert.Error(err)
		assert.Equal(state, 1)
	})
	plugin.PresetName = "test"
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		assert.Error(err)
		assert.Equal(state, 1)
	})
	plugin.PresetName = "None"
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		assert.Error(err)
		assert.Equal(state, 1)
	})
	plugin.DryRun = true
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		assert.NoError(err)
		assert.Equal(state, 0)
	})
	plugin.DimensionFilterStrings = []string{"what=now", "brown=cow"}
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		assert.NoError(err)
		assert.Equal(state, 0)
	})
	plugin.DimensionFilterStrings = []string{"what=now=brown=cow", "when=where=what=why"}
	t.Run("CheckArgs", func(t *testing.T) {
		state, err := checkArgs(nil)
		assert.Error(err)
		assert.Equal(state, 1)
	})
	cleanPluginValues()
}

func TestCheckFunction(t *testing.T) {
	defer quiet()()
	cleanPluginValues()
	plugin.PresetName = "None"
	plugin.StatsList = []string{"Average"}
	none := presets.None{}
	none.AddStats(plugin.StatsList)
	plugin.Preset = &none
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
		expectedId      string
	}{ //start of array
		{ //start of struct
			client:          mockService{},
			maxPages:        2,
			nextToken:       true,
			expectedState:   0,
			includeMessages: false,
			expectedId:      "test",
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        1,
			nextToken:       true,
			expectedState:   1,
			includeMessages: false,
			expectedId:      "test",
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        2,
			nextToken:       true,
			expectedState:   1,
			includeMessages: true,
			expectedId:      "test",
		},
	}
	for i, tt := range cases {
		t.Run("CheckFunction Run: "+strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			client.includeMessages = tt.includeMessages
			client.dataResultId = tt.expectedId
			nextToken = tt.nextToken
			plugin.MaxPages = tt.maxPages
			fmt.Printf("plugin:: %+v\n", plugin)
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

func TestCheckFunctionDryRun(t *testing.T) {
	defer quiet()()
	cleanPluginValues()
	plugin.PresetName = "None"
	plugin.DryRun = true
	plugin.RecentlyActive = true
	plugin.MetricName = "test"
	plugin.Namespace = "test"
	plugin.Verbose = true
	plugin.StatsList = []string{"Average"}
	cases := []struct {
		client          mockService
		expectedState   int
		nextToken       bool
		maxPages        int
		includeMessages bool
		expectedId      string
	}{ //start of array
		{ //start of struct
			client:          mockService{},
			maxPages:        2,
			nextToken:       true,
			expectedState:   0,
			includeMessages: false,
			expectedId:      "test",
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        1,
			nextToken:       true,
			expectedState:   1,
			includeMessages: false,
			expectedId:      "test",
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        0,
			nextToken:       true,
			expectedState:   0,
			includeMessages: false,
			expectedId:      "test",
		},
		{ //start of struct
			client:          mockService{},
			maxPages:        2,
			nextToken:       true,
			expectedState:   0,
			includeMessages: true,
			expectedId:      "test",
		},
	}
	for i, tt := range cases {
		t.Run("CheckFunction Run: "+strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			client.includeMessages = tt.includeMessages
			client.dataResultId = tt.expectedId
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
