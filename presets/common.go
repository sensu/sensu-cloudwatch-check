package presets

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/google/uuid"
)

var (
	Presets = make(map[string]PresetInterface)
)

func init() {
	Presets["None"] = &None{Preset: Preset{Description: "No Service Presets Active, use cmdline --namespace --metric --dimension-filters to tailer cloudwatch results"}}
	Presets["CLB"] = &CLB{Preset: Preset{Description: "Preset Metrics for AWS Classic Load Balancer"}}
}

type Preset struct {
	Metrics           []types.Metric
	DimensionFilters  []types.DimensionFilter
	Namespace         string
	MetricName        string
	Description       string
	Name              string
	configMap         map[string][]StatConfig
	measurementString string
	verbose           bool
}

type PresetInterface interface {
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

func (p *Preset) AddDimensionFilters(filters []types.DimensionFilter) error {
	for _, f := range filters {
		p.DimensionFilters = append(p.DimensionFilters, f)
	}
	return nil
}

func (p *Preset) SetMeasurementString(mstring string) error {
	p.measurementString = mstring
	measurementConfig := MeasurementJSON{}
	if err := json.Unmarshal([]byte(p.measurementString), &measurementConfig); err != nil {
		return err
	}
	p.configMap = make(map[string][]StatConfig)
	for _, metric := range measurementConfig.Metrics {
		p.configMap[metric.MetricName] = []StatConfig{}
		for _, item := range metric.Config {
			p.configMap[metric.MetricName] = append(p.configMap[metric.MetricName], item)
		}

	}

	return nil
}

func (p *Preset) GetDimensionFilters() []types.DimensionFilter {
	return p.DimensionFilters
}

func (p *Preset) GetNamespace() string {
	return p.Namespace
}

func (p *Preset) GetMetricName() string {
	return ""
}

func (p *Preset) SetMetricName(name string) error {
	return nil
}

func (p *Preset) GetDescription() string {
	return p.Description
}

func (p *Preset) AddMetrics(metrics []types.Metric) error {
	if p.verbose {
		fmt.Println("Preset::AddMetrics", len(metrics))
	}
	errStrings := []string{}
	for _, m := range metrics {
		if p.verbose {
			fmt.Printf("Preset.AddMetrics: Metric: %v\n", *m.MetricName)
		}
		if _, ok := p.configMap[*m.MetricName]; ok {
			if p.verbose {
				fmt.Printf("Preset.AddMetrics: Found config for Metric: %v\n", *m.MetricName)
			}
			p.Metrics = append(p.Metrics, m)
		} else {
			str := fmt.Sprintf("Preset.AddMetrics: No config for Metric: %v\n", *m.MetricName)
			if p.verbose {
				fmt.Println(str)
			}
			errStrings = append(errStrings, str)
		}
	}
	if len(errStrings) > 0 {
		return fmt.Errorf("%v", strings.Join(errStrings, ""))
	} else {
		return nil
	}
}

func (p *Preset) BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error) {
	if p.verbose {
		fmt.Println("Preset::BuildMetricDataQueries")
	}
	dataQueries := []types.MetricDataQuery{}
	for _, m := range p.Metrics {
		if statConfigs, ok := p.configMap[*m.MetricName]; ok {
			for _, config := range statConfigs {
				stat := config.Stat
				measurement := config.Measurement
				id := uuid.New()
				idString := "aws_" + strings.ReplaceAll(id.String(), "-", "_")
				if p.verbose {
					fmt.Printf("Preset.BuildMetricDataQueries: %v %v %v %v\n", *m.MetricName, idString, stat, measurement)
				}
				labelString := measurement
				dataQuery := types.MetricDataQuery{
					Id:    &idString,
					Label: &labelString,
					MetricStat: &types.MetricStat{
						Metric: &m,
						Period: aws.Int32(60 * period),
						Stat:   aws.String(stat),
					},
				}
				dataQueries = append(dataQueries, dataQuery)
			}
		} else {
			fmt.Printf("Preset.BuildMetricDataQueries no config for: %v\n", *m.MetricName)
		}
	}
	return dataQueries, nil
}
