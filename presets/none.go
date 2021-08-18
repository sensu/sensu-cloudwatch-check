package presets

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/google/uuid"
	"github.com/sensu/sensu-cloudwatch-check/common"
)

type None struct {
	Metrics           []types.Metric
	Stats             []string
	DimensionFilters  []types.DimensionFilter
	Namespace         string
	MetricName        string
	MeasurementString string
}

func (p *None) Init() error {
	return nil
}

func (p *None) AddMetrics(metrics []types.Metric) {
	for _, m := range metrics {
		p.Metrics = append(p.Metrics, m)
	}
	return
}

func (p *None) AddDimensionFilters(filters []types.DimensionFilter) {
	for _, f := range filters {
		p.DimensionFilters = append(p.DimensionFilters, f)
	}
	return
}

func (p *None) AddStats(stats []string) {
	for _, s := range stats {
		p.Stats = append(p.Stats, s)
	}
	return
}

func (p *None) GetDimensionFilters() []types.DimensionFilter {
	return p.DimensionFilters
}

func (p *None) GetNamespace() string {
	return p.Namespace
}

func (p *None) GetMetricName() string {
	return p.MetricName
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