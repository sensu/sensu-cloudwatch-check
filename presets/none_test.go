package presets

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestNoneInit(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.Init(true)
	assert.NoError(err)
}

func TestNoneAddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	none := &None{}
	err := none.Init(false)
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
	err := none.Init(false)
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
		}
		metrics = append(metrics, m)
	}
	err = none.AddMetrics(metrics)
	assert.NoError(err)
	none.BuildMetricDataQueries(int32(1))
}