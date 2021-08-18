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

type ELB struct {
	Metrics           []types.Metric
	Stats             []string
	DimensionFilters  []types.DimensionFilter
	Namespace         string
	MetricName        string
	MeasurementString string
	configMap         map[string][]StatConfig
}

func (p *ELB) Init() error {
	p.Namespace = "AWS/ELB"
	p.Stats = []string{"Average"}
	if filters, err := common.BuildDimensionFilters([]string{
		"LoadBalancerName", "AvailabilityZone"}); err == nil {
		p.DimensionFilters = filters
	} else {
		return err
	}
	p.MeasurementString = `{"metrics" : 
                                  [{"metric" : "test" , "config" : 
				      [{"stat" : "Average"     , "measurement" : "aws.elb.test.ave" }, 
                                       {"stat" : "Sum"         , "measurement" : "aws.elb.test.sum" }, 
                                       {"stat" : "SampleCount" , "measurement" : "aws.elb.test.count" } 
                                      ]	
			           }
			          ]
		     }
		     `
	measurementConfig := MeasurementJSON{}
	if err := json.Unmarshal([]byte(p.MeasurementString), &measurementConfig); err != nil {
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

func (p *ELB) AddMetrics(metrics []types.Metric) {
	for _, m := range metrics {
		p.Metrics = append(p.Metrics, m)
	}
	return
}

func (p *ELB) AddDimensionFilters(filters []types.DimensionFilter) {
	for _, f := range filters {
		p.DimensionFilters = append(p.DimensionFilters, f)
	}
	return
}

func (p *ELB) AddStats(stats []string) {
	for _, s := range stats {
		p.Stats = append(p.Stats, s)
	}
	return
}

func (p *ELB) GetDimensionFilters() []types.DimensionFilter {
	return p.DimensionFilters
}

func (p *ELB) GetNamespace() string {
	return p.Namespace
}

func (p *ELB) GetMetricName() string {
	return p.MetricName
}

func (p *ELB) BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error) {
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
