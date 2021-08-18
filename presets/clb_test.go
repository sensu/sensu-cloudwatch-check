package presets

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

var (
	enableQuiet = false
)

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

func TestCLBInit(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	elb := &CLB{}
	err := elb.Init(true)
	assert.NoError(err)
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
