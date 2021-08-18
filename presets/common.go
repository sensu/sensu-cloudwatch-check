package presets

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var (
	Presets = make(map[string]ServicePreset)
)

func init() {
	Presets["None"] = &None{}
	Presets["ELB"] = &ELB{}
}

type ServicePreset interface {
	BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error)
	AddMetrics(metrics []types.Metric)
	GetNamespace() string
	GetMetricName() string
	GetDimensionFilters() []types.DimensionFilter
	Init() error
}

type MeasurementJSON struct {
	Metrics []struct {
		MetricName string `json:"metric"`
		Config     []struct {
			Stat        string `json:"stat"`
			Measurement string `json:"measurement"`
		} `json:"config"`
	}
}
