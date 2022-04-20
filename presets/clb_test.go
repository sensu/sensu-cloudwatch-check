package presets

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestCLBReady(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CLB{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	assert.Equal("AWS/ELB", preset.Namespace)
	assert.Equal(2, len(preset.DimensionFilters))
	allowed := []string{"LoadBalancerName", "AvailabilityZone"}
	for _, d := range preset.DimensionFilters {
		assert.Contains(allowed, *d.Name)
		assert.Nil(d.Value)
	}
}

func TestCLBAddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CLB{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
        err = preset.SetErrorOnMissing(true)
        assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	metricNames := []string{
		"test",
	}
	metrics := []types.Metric{}
	namespace := "AWS/ELB"
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &namespace,
		}
		metrics = append(metrics, m)
	}
	err = preset.AddMetrics(metrics)
	assert.Error(err)
	metricNames = []string{
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
	metrics = []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &namespace,
		}
		metrics = append(metrics, m)
	}
	err = preset.AddMetrics(metrics)
	assert.NoError(err)
}

func TestCLBBuildMetricDataQueries(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CLB{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
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
	namespace := "AWS/ELB"
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &namespace,
		}
		metrics = append(metrics, m)
	}
	err = preset.AddMetrics(metrics)
	assert.NoError(err)
	_, err = preset.BuildMetricDataQueries(int32(1))
	assert.NoError(err)
}

func TestCLBGetMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CLB{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	output, err := preset.GetMeasurementString(true)
	assert.NoError(err)
	assert.Greater(len(output), 30)
	fmt.Println(output)

}
