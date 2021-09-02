package presets

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/google/uuid"
	"github.com/sensu/sensu-cloudwatch-check/common"
)

var (
	Presets = make(map[string]PresetInterface)
)

func init() {
	Presets["None"] = &None{Preset: Preset{Description: "No Service Presets Active, use cmdline --namespace --metric --dimension-filters to tailer cloudwatch results"}}
	Presets["CLB"] = &CLB{Preset: Preset{Description: "Preset Metrics for AWS Classic Load Balancer"}}
	Presets["ALB"] = &ALB{Preset: Preset{Description: "Preset Metrics for AWS Application Load Balancer"}}
	Presets["EC2"] = &EC2{Preset: Preset{Description: "Preset Metrics for AWS EC2"}}
	Presets["CloudFront"] = &CloudFront{Preset: Preset{Description: "Preset Metrics for AWS CloudFront. Note: requires --region us-east-1"}}
}

type Preset struct {
	Metrics           []types.Metric
	DimensionFilters  []types.DimensionFilter
	Namespace         string
	MetricFilter      string
	Region            string
	PeriodMinutes     int
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
	GetMetricFilter() string
	SetMetricFilter(name string) error
	GetPeriodMinutes() int
	SetPeriodMinutes(period int) error
	GetRegion() string
	SetRegion(region string) error
	SetVerbose(flag bool) error
	SetMeasurementString(config string) error
	BuildMeasurementConfig() error
	GetMeasurementString(pretty bool) (string, error)
	GetDimensionFilters() []types.DimensionFilter
	AddDimensionFilters(filters []types.DimensionFilter) error
	Ready() error
}

type StatConfig struct {
	Stat        string `json:"stat"`
	Measurement string `json:"measurement"`
}
type MeasurementConfig struct {
	MetricName string       `json:"metric"`
	Config     []StatConfig `json:"config"`
}

type MeasurementJSON struct {
	Namespace        string              `json:"namespace"`
	PeriodMinutes    int                 `json:"period-minutes,omitempty"`
	Region           string              `json:"region,omitempty"`
	MetricFilter     string              `json:"metric-filter,omitempty"`
	DimensionFilters []string            `json:"dimension-filters,omitempty"`
	Measurements     []MeasurementConfig `json:"measurements,omitempty"`
}

func (p *Preset) AddDimensionFilters(filters []types.DimensionFilter) error {
	for _, f := range filters {
		p.DimensionFilters = append(p.DimensionFilters, f)
	}
	return nil
}

func (p *Preset) GetMeasurementString(pretty bool) (string, error) {
	measurementConfig := MeasurementJSON{}
	measurementConfig.Measurements = []MeasurementConfig{}
	measurementConfig.Namespace = p.Namespace
	measurementConfig.PeriodMinutes = p.PeriodMinutes
	measurementConfig.Region = p.Region
	measurementConfig.MetricFilter = p.MetricFilter
	dimStrings := []string{}
	for _, d := range p.DimensionFilters {
		output := strings.TrimSpace(*d.Name)
		if d.Value != nil {
			output += "=" + strings.TrimSpace(*d.Value)
		}
		dimStrings = append(dimStrings, output)
	}
	measurementConfig.DimensionFilters = common.RemoveDuplicateStrings(dimStrings)

	for metricName := range p.configMap {
		config := MeasurementConfig{
			MetricName: metricName,
			Config:     p.configMap[metricName],
		}
		measurementConfig.Measurements = append(measurementConfig.Measurements, config)
	}
	prefix := ""
	indent := ""
	if pretty {
		indent = "  "
	}
	if output, err := json.MarshalIndent(measurementConfig, prefix, indent); err != nil {
		return "", err
	} else {
		return string(output), nil

	}
}

func (p *Preset) BuildMeasurementConfig() error {
	measurementConfig := MeasurementJSON{}
	if err := json.Unmarshal([]byte(p.measurementString), &measurementConfig); err != nil {
		return err
	}
	if len(measurementConfig.Namespace) > 0 {
		p.Namespace = measurementConfig.Namespace
	}
	if measurementConfig.PeriodMinutes > 0 {
		p.PeriodMinutes = measurementConfig.PeriodMinutes
	}
	if len(measurementConfig.Region) > 0 {
		p.Region = measurementConfig.Region
	}
	if len(measurementConfig.DimensionFilters) > 0 {
		if dimensionFilters, err := common.BuildDimensionFilters(measurementConfig.DimensionFilters); err == nil {
			p.AddDimensionFilters(dimensionFilters)
		} else {
			return err
		}
	}
	p.configMap = make(map[string][]StatConfig)
	for _, m := range measurementConfig.Measurements {
		p.configMap[m.MetricName] = []StatConfig{}
		for _, item := range m.Config {
			p.configMap[m.MetricName] = append(p.configMap[m.MetricName], item)
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

func (p *Preset) GetMetricFilter() string {
	return p.MetricFilter
}

func (p *Preset) SetMetricFilter(name string) error {
	p.MetricFilter = name
	return nil
}

func (p *Preset) GetPeriodMinutes() int {
	return p.PeriodMinutes
}

func (p *Preset) SetPeriodMinutes(period int) error {
	p.PeriodMinutes = period
	return nil
}

func (p *Preset) GetRegion() string {
	return p.Region
}

func (p *Preset) SetRegion(region string) error {
	p.Region = region
	return nil
}

func (p *Preset) SetVerbose(flag bool) error {
	p.verbose = flag
	return nil
}

func (p *Preset) SetMeasurementString(config string) error {
	p.measurementString = config
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

		if len(p.MetricFilter) > 0 {
			if p.MetricFilter != *m.MetricName {
				str := fmt.Sprintf("Preset.AddMetrics: MetricFilter: %v does not match Metric: %v \n", p.MetricFilter, *m.MetricName)
				if p.verbose {
					fmt.Println(str)
				}
				errStrings = append(errStrings, str)
				continue
			}
		}

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

// overwrite the Ready function when building a new preset to enforce specific behavior
func (p *Preset) Ready() error {
	return nil
}
