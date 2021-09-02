package presets

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestEC2Ready(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &EC2{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	assert.Equal("AWS/EC2", preset.Namespace)
	assert.Equal(0, len(preset.DimensionFilters))
	allowed := []string{"LoadBalancerName", "AvailabilityZone"}
	for _, d := range preset.DimensionFilters {
		assert.Contains(allowed, *d.Name)
		assert.Nil(d.Value)
	}
}

func TestEC2AddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &EC2{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
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
	err = preset.AddMetrics(metrics)
	assert.Error(err)
	metricNames = []string{}
	metrics = []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
		}
		metrics = append(metrics, m)
	}
	err = preset.AddMetrics(metrics)
	assert.NoError(err)
}

func TestEC2BuildMetricDataQueries(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &EC2{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	metricNames := []string{}
	metrics := []types.Metric{}
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
		}
		metrics = append(metrics, m)
	}
	err = preset.AddMetrics(metrics)
	assert.NoError(err)
	preset.BuildMetricDataQueries(int32(1))
}

func TestEC2GetMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &EC2{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	output, err := preset.GetMeasurementString(true)
	assert.NoError(err)
	assert.Greater(len(output), 30)
	fmt.Println(output)

}
