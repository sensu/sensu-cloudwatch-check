package presets

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestNoneReady(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.SetVerbose(true)
	assert.NoError(err)
	err = none.Ready()
	assert.NoError(err)
}

func TestNoneAddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.SetVerbose(true)
	assert.NoError(err)
	err = none.Ready()
	assert.NoError(err)
	metricNames := []string{
		"test",
	}
	metrics := []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
		}
		metrics = append(metrics, m)
	}
	err = none.AddMetrics(metrics)
	assert.NoError(err)
}

func TestNoneBuildMetricDataQueries(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.SetVerbose(true)
	assert.NoError(err)
	err = none.Ready()
	assert.NoError(err)
	metricNames := []string{
		"BackendConnectionErrors",
		"DesyncMitigationMode_NonCompliant_Request_Count",
		"HealthyHostCount",
		"HTTPCode_Backend_5XX",
		"HTTPCode_Backend_4XX",
		"HTTPCode_Backend_3XX",
		"HTTPCode_Backend_2XX",
		"HTTPCode_ELB_4XX",
		"HTTPCode_ELB_5XX",
		"Latency",
		"RequestCount",
		"SpilloverCount",
		"SurgeQueueLength",
		"UnHealthyHostCount",
	}
	metrics := []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  aws.String("Test"),
		}
		metrics = append(metrics, m)
	}
	stats := []string{"Average", "Sum"}
	none.AddStats(stats)
	err = none.AddMetrics(metrics)
	assert.NoError(err)
	_, err = none.BuildMetricDataQueries(int32(1))
	assert.NoError(err)

}

func TestNoneBuildMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.SetVerbose(true)
	assert.NoError(err)
	err = none.Ready()
	none.Namespace = "Test"
	assert.NoError(err)
	metricNames := []string{
		"BackendConnectionErrors",
		"DesyncMitigationMode_NonCompliant_Request_Count",
		"HealthyHostCount",
		"HTTPCode_Backend_5XX",
		"HTTPCode_Backend_4XX",
		"HTTPCode_Backend_3XX",
		"HTTPCode_Backend_2XX",
		"HTTPCode_ELB_4XX",
		"HTTPCode_ELB_5XX",
		"Latency",
		"RequestCount",
		"SpilloverCount",
		"SurgeQueueLength",
		"UnHealthyHostCount",
	}
	metrics := []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &none.Namespace,
		}
		metrics = append(metrics, m)
	}
	err = none.AddMetrics(metrics)
	stats := []string{"Average", "Sum"}
	none.AddStats(stats)
	assert.NoError(err)
	err = none.BuildMeasurementString()
	assert.NoError(err)
}
func TestNoneGetMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.SetVerbose(true)
	assert.NoError(err)
	err = none.Ready()
	assert.NoError(err)
	none.Namespace = "Test"
	stats := []string{"Average", "Sum"}
	dims := []types.DimensionFilter{
		types.DimensionFilter{
			Name:  aws.String("Dim1_Name"),
			Value: aws.String("Dim1_Value"),
		},
		types.DimensionFilter{
			Name: aws.String("Dim2_Name"),
		},
	}
	err = none.AddDimensionFilters(dims)
	assert.NoError(err)
	metricNames := []string{
		"BackendConnectionErrors",
		"DesyncMitigationMode_NonCompliant_Request_Count",
		"HealthyHostCount",
		"HTTPCode_Backend_5XX",
		"HTTPCode_Backend_4XX",
		"HTTPCode_Backend_3XX",
		"HTTPCode_Backend_2XX",
		"HTTPCode_ELB_4XX",
		"HTTPCode_ELB_5XX",
		"Latency",
		"RequestCount",
		"SpilloverCount",
		"SurgeQueueLength",
		"UnHealthyHostCount",
	}
	metrics := []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &none.Namespace,
		}
		metrics = append(metrics, m)
	}
	err = none.AddMetrics(metrics)
	assert.NoError(err)
	none.AddStats(stats)
	output, err := none.GetMeasurementString(false)
	assert.NoError(err)
	assert.Greater(len(output), 30)
}
