package presets

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestCLBInit(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &CLB{}
	err := elb.Init(true)
	assert.NoError(err)
	assert.Equal("AWS/ELB", elb.Namespace)
	assert.Equal(2, len(elb.DimensionFilters))
	allowed := []string{"LoadBalancerName", "AvailabilityZone"}
	for _, d := range elb.DimensionFilters {
		assert.Contains(allowed, *d.Name)
		assert.Nil(d.Value)
	}
}

func TestCLBAddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &CLB{}
	err := elb.Init(false)
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
	err = elb.AddMetrics(metrics)
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
		}
		metrics = append(metrics, m)
	}
	err = elb.AddMetrics(metrics)
	assert.NoError(err)
}

func TestCLBBuildMetricDataQueries(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &CLB{}
	err := elb.Init(false)
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
	err = elb.AddMetrics(metrics)
	assert.NoError(err)
	elb.BuildMetricDataQueries(int32(1))
}

func TestCLBGetMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &CLB{}
	err := elb.Init(true)
	assert.NoError(err)
	output, err := elb.GetMeasurementString(true)
	assert.NoError(err)
	assert.Greater(len(output), 30)
	fmt.Println(output)

}
