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
	ServicePresetStruct
	MetricName string
	Stats      []string
}

func (p *None) Init(verbose bool) error {
	return nil
}

func (p *None) AddMetrics(metrics []types.Metric) error {
	for _, m := range metrics {
		p.Metrics = append(p.Metrics, m)
	}
	return nil
}

func (p *None) AddDimensionFilters(filters []types.DimensionFilter) error {
	for _, f := range filters {
		p.DimensionFilters = append(p.DimensionFilters, f)
	}
	return nil
}

func (p *None) AddStats(stats []string) {
	for i, _ := range stats {
		p.Stats = append(p.Stats, strings.TrimSpace(stats[i]))
	}
	return
}

func (p *None) GetDimensionFilters() []types.DimensionFilter {
	return p.DimensionFilters
}

func (p *None) GetDescription() string {
	return p.Description
}

func (p *None) GetNamespace() string {
	return p.Namespace
}

func (p *None) GetMetricName() string {
	return p.MetricName
}

func (p *None) SetMetricName(name string) error {
	p.MetricName = name
	return nil
}

func (p *None) SetMeasurementString(mstring string) error {
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

func (p *None) BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error) {
	dataQueries := []types.MetricDataQuery{}
	for i, _ := range p.Metrics {
		for j, _ := range p.Stats {
			id := uuid.New()
			idString := "aws_" + strings.ReplaceAll(id.String(), "-", "_")
			labelString := fmt.Sprintf("%v.%v", common.BuildLabelBase(p.Metrics[i]), common.ToSnakeCase(p.Stats[j]))
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
