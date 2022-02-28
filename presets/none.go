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

type None struct {
	Preset
	Stats []string
}

func (p *None) GetMeasurementString(pretty bool) (string, error) {
	if err := p.BuildMeasurementString(); err != nil {
		return "", err
	}
	if err := p.BuildMeasurementConfig(); err != nil {
		return "", err
	}
	return p.Preset.GetMeasurementString(pretty)
}

func (p *None) BuildMeasurementString() error {
	if p.verbose {
		fmt.Println("None::BuildMeasurementString")
	}
	if len(p.Namespace) == 0 {
		return fmt.Errorf("Namespace is not set")
	}
	if len(p.Stats) == 0 {
		return fmt.Errorf("No stats set")
	}
	if len(p.Metrics) == 0 {
		return fmt.Errorf("No metrics set")
	}
	measurementConfig := MeasurementJSON{}
	measurementConfig.Measurements = []MeasurementConfig{}
	measurementConfig.Namespace = p.Namespace
	measurementConfig.PeriodMinutes = p.PeriodMinutes
	measurementConfig.MetricFilter = p.MetricFilter
	dimStrings := []string{}
	for _, d := range p.DimensionFilters {
		output := strings.TrimSpace(*d.Name)
		if d.Value != nil {
			output += "=" + strings.TrimSpace(*d.Value)
		}
		dimStrings = append(dimStrings, output)
	}
	measurementConfig.DimensionFilters = dimStrings
	for i := range p.Metrics {
		config := MeasurementConfig{
			MetricName: *p.Metrics[i].MetricName,
			Config:     []StatConfig{},
		}

		for j := range p.Stats {
			labelString := fmt.Sprintf("%v.%v", common.BuildLabelBase(p.Metrics[i]), common.ToSnakeCase(p.Stats[j]))
			s := StatConfig{
				Stat:        p.Stats[j],
				Measurement: labelString,
			}
			config.Config = append(config.Config, s)
		}

		measurementConfig.Measurements = append(measurementConfig.Measurements, config)
	}

	if output, err := json.Marshal(measurementConfig); err != nil {
		return err
	} else {
		p.measurementString = string(output)
		return nil
	}
}

func (p *None) AddStats(stats []string) {
	if p.verbose {
		fmt.Println("None::AddStats", stats)
	}
	for i := range stats {
		p.Stats = append(p.Stats, strings.TrimSpace(stats[i]))
	}
}

func (p *None) GetMetricFilter() string {
	if p.verbose {
		fmt.Println("None::GetMetricName", p.MetricFilter)
	}
	return p.MetricFilter
}

func (p *None) SetMetricFilter(name string) error {
	if p.verbose {
		fmt.Println("None::SetMetricFilter", name)
	}
	p.MetricFilter = name
	return nil
}

func (p *None) AddMetrics(metrics []types.Metric) error {
	if p.verbose {
		fmt.Println("None::AddMetrics", len(metrics))
	}
	p.Metrics = append(p.Metrics, metrics...)
	return nil
}

func (p *None) BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error) {
	if p.verbose {
		fmt.Println("None::BuildMetricDataQueries")
	}
	dataQueries := []types.MetricDataQuery{}
	for i := range p.Metrics {
		for j := range p.Stats {
			id := uuid.New()
			idString := "aws_" + strings.ReplaceAll(id.String(), "-", "_")
			labelString := fmt.Sprintf("%v_%v", common.BuildLabelBase(p.Metrics[i]), common.ToSnakeCase(p.Stats[j]))
			dataQuery := types.MetricDataQuery{
				Id:    &idString,
				Label: &labelString,
				MetricStat: &types.MetricStat{
					Metric: &p.Metrics[i],
					Period: aws.Int32(60 * period),
					Stat:   aws.String(p.Stats[j]),
				},
			}
			dataQueries = append(dataQueries, dataQuery)
		}
	}
	return dataQueries, nil
}
