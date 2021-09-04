package presets

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

func TestALBReady(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &ALB{}
	err := elb.SetVerbose(true)
	assert.NoError(err)
	err = elb.Ready()
	assert.NoError(err)
	assert.Equal("AWS/ApplicationELB", elb.Namespace)
	assert.Equal(0, len(elb.DimensionFilters))
	allowed := []string{"LoadBalancerName", "AvailabilityZone"}
	for _, d := range elb.DimensionFilters {
		assert.Contains(allowed, *d.Name)
		assert.Nil(d.Value)
	}
}

func TestALBAddMetrics(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &ALB{}
	err := elb.SetVerbose(true)
	assert.NoError(err)
	err = elb.Ready()
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
	err = elb.AddMetrics(metrics)
	assert.Error(err)
	metricNames = []string{
		"DesyncMitigationMode_NonCompliant_Request_Count",
		"HTTPCode_ELB_502_Count",
		"HealthyHostCount",
		"RequestCount",
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
	err = elb.AddMetrics(metrics)
	assert.NoError(err)
}

func TestALBBuildMetricDataQueries(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &ALB{}
	err := elb.SetVerbose(true)
	assert.NoError(err)
	err = elb.Ready()
	assert.NoError(err)
	metricNames := []string{
		"DesyncMitigationMode_NonCompliant_Request_Count",
		"HealthyHostCount",
		"HTTPCode_ELB_4XX_Count",
		"HTTPCode_ELB_5XX_Count",
		"RequestCount",
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
	err = elb.AddMetrics(metrics)
	assert.NoError(err)
	elb.BuildMetricDataQueries(int32(1))
}

func TestALBGetMeasurementString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &ALB{}
	err := elb.SetVerbose(true)
	assert.NoError(err)
	err = elb.Ready()
	assert.NoError(err)
	output, err := elb.GetMeasurementString(true)
	assert.NoError(err)
	assert.Greater(len(output), 30)
	fmt.Println(output)

}
