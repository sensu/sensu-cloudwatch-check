package presets

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestCloudFrontReady(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CloudFront{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	assert.Equal("AWS/CloudFront", preset.Namespace)
	assert.Equal(0, len(preset.DimensionFilters))
	allowed := []string{}
	for _, d := range preset.DimensionFilters {
		assert.Contains(allowed, *d.Name)
		assert.Nil(d.Value)
	}
}

func TestCloudFrontAddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CloudFront{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	metricNames := []string{
		"test",
	}
	metrics := []types.Metric{}
	namespace := "AWS/CloudFront"
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &namespace,
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

func TestCloudFrontBuildMetricDataQueries(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CloudFront{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	metricNames := []string{}
	metrics := []types.Metric{}
	namespace := "AWS/CloudFront"
	for i := range metricNames {
		m := types.Metric{
			MetricName: &metricNames[i],
			Namespace:  &namespace,
		}
		metrics = append(metrics, m)
	}
	err = preset.AddMetrics(metrics)
	assert.NoError(err)
	preset.BuildMetricDataQueries(int32(1))
}

func TestCloudFrontGetMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	preset := &CloudFront{}
	err := preset.SetVerbose(true)
	assert.NoError(err)
	err = preset.Ready()
	assert.NoError(err)
	output, err := preset.GetMeasurementString(true)
	assert.NoError(err)
	assert.Greater(len(output), 30)
	fmt.Println(output)

}
