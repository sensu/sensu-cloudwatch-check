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
	Preset
	Stats []string
}

func (p *None) Init(verbose bool) error {
	p.verbose = verbose
	if p.verbose {
		fmt.Println("None::Init Setting up none preset")
	}
	p.addMetricsFunc = p.AddMetrics
	p.buildMetricDataQueriesFunc = p.BuildMetricDataQueries
	p.setMetricNameFunc = p.SetMetricName
	p.getMetricNameFunc = p.GetMetricName
	return nil
}

func (p *None) AddStats(stats []string) {
	if p.verbose {
		fmt.Println("None::AddStats", stats)
	}
	for i, _ := range stats {
		p.Stats = append(p.Stats, strings.TrimSpace(stats[i]))
	}
	return
}

func (p *None) GetMetricName() string {
	if p.verbose {
		fmt.Println("None::GetMetricName", p.MetricName)
	}
	return p.MetricName
}

func (p *None) SetMetricName(name string) error {
	if p.verbose {
		fmt.Println("None::SetMetricName", name)
	}
	p.MetricName = name
	return nil
}

func (p *None) AddMetrics(metrics []types.Metric) error {
	if p.verbose {
		fmt.Println("None::AddMetrics", len(metrics))
	}
	for _, m := range metrics {
		p.Metrics = append(p.Metrics, m)
	}
	return nil
}

func (p *None) BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error) {
	if p.verbose {
		fmt.Println("None::BuildMetricDataQueries")
	}
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
