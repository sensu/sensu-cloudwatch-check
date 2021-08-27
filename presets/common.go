package presets

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var (
	Presets = make(map[string]ServicePreset)
)

func init() {
	Presets["None"] = &None{ServicePresetStruct: ServicePresetStruct{Description: "No Service Presets Active, use cmdline --namespace --metric --dimension-filters to tailer cloudwatch results"}}
	Presets["CLB"] = &CLB{ServicePresetStruct: ServicePresetStruct{Description: "Preset Metrics for AWS Classic Load Balancer"}}
}

type ServicePresetStruct struct {
	Metrics           []types.Metric
	DimensionFilters  []types.DimensionFilter
	Namespace         string
	Description       string
	Name              string
	configMap         map[string][]StatConfig
	measurementString string
	verbose           bool
}

type ServicePreset interface {
	BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error)
	AddMetrics(metrics []types.Metric) error
	GetDescription() string
	GetNamespace() string
	GetMetricName() string
	SetMetricName(name string) error
	SetMeasurementString(name string) error
	GetDimensionFilters() []types.DimensionFilter
	AddDimensionFilters(filters []types.DimensionFilter) error
	Init(verbose bool) error
}

type StatConfig struct {
	Stat        string `json:"stat"`
	Measurement string `json:"measurement"`
}
type MetricConfig struct {
	MetricName string       `json:"metric"`
	Config     []StatConfig `json:"config"`
}
type MeasurementJSON struct {
	Metrics []MetricConfig `json:"metrics"`
}
