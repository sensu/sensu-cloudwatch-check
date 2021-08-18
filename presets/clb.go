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

type CLB struct {
	Metrics           []types.Metric
	Stats             []string
	DimensionFilters  []types.DimensionFilter
	Namespace         string
	MeasurementString string
	configMap         map[string][]StatConfig
	verbose           bool
	Description       string
	Name              string
}

func (p *CLB) AddDimensionFilters(filters []types.DimensionFilter) {
	for _, f := range filters {
		p.DimensionFilters = append(p.DimensionFilters, f)
	}
	return
}

func (p *CLB) AddStats(stats []string) {
	for _, s := range stats {
		p.Stats = append(p.Stats, s)
	}
	return
}

func (p *CLB) GetDimensionFilters() []types.DimensionFilter {
	return p.DimensionFilters
}

func (p *CLB) GetNamespace() string {
	return p.Namespace
}
func (p *CLB) GetMetricName() string {
	return ""
}

func (p *CLB) GetDescription() string {
	return p.Description
}

func (p *CLB) AddMetrics(metrics []types.Metric) error {
	errStrings := []string{}
	for _, m := range metrics {
		if p.verbose {
			fmt.Printf("CLB.AddMetrics: Metric: %v\n", *m.MetricName)
		}
		if _, ok := p.configMap[*m.MetricName]; ok {
			if p.verbose {
				fmt.Printf("CLB.AddMetrics: Found config for Metric: %v\n", *m.MetricName)
			}
			p.Metrics = append(p.Metrics, m)
		} else {
			str := fmt.Sprintf("CLB.AddMetrics: No config for Metric: %v\n", *m.MetricName)
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

func (p *CLB) BuildMetricDataQueries(period int32) ([]types.MetricDataQuery, error) {
	dataQueries := []types.MetricDataQuery{}
	for _, m := range p.Metrics {
		if statConfigs, ok := p.configMap[*m.MetricName]; ok {
			for _, config := range statConfigs {
				stat := config.Stat
				measurement := config.Measurement
				id := uuid.New()
				idString := "aws_" + strings.ReplaceAll(id.String(), "-", "_")
				if p.verbose {
					fmt.Printf("CLB.BuildMetricDataQueries: %v %v %v %v\n", *m.MetricName, idString, stat, measurement)
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
			fmt.Printf("CLB.BuildMetricDataQueries no config for: %v\n", *m.MetricName)
		}
	}
	return dataQueries, nil
}

func (p *CLB) Init(verbose bool) error {
	p.verbose = verbose
	p.Namespace = "AWS/ELB"
	p.Stats = []string{"Average"}
	if filters, err := common.BuildDimensionFilters([]string{
		"LoadBalancerName", "AvailabilityZone"}); err == nil {
		p.DimensionFilters = filters
	} else {
		return err
	}

	// JSON Config String developed on 2021-08-18 from AWS Cloudwatch documentation
	//  Ref: https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/elb-cloudwatch-metrics.html#loadbalancing-metrics-clb
	p.MeasurementString = `{"metrics" : 
                                  [
				   {"metric":"HTTPCode_ELB_4XX" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.httpcode_elb_4xx"} 
                                      ]	
			           },
				   {"metric":"HTTPCode_ELB_5XX" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.httpcode_elb_5xx"} 
                                      ]	
			           },
				   {"metric":"HTTPCode_Backend_5XX" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.httpcode_backend_5xx"} 
                                      ]	
			           },
				   {"metric":"HTTPCode_Backend_4XX" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.httpcode_backend_4xx"} 
                                      ]	
			           },
				   {"metric":"HTTPCode_Backend_3XX" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.httpcode_backend_3xx"} 
                                      ]	
			           },
				   {"metric":"HTTPCode_Backend_2XX" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.httpcode_backend_2xx"} 
                                      ]	
			           },
				   {"metric":"BackendConnectionErrors" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.backend_connection_errors"} 
                                      ]	
			           },
				   {"metric":"DesyncMitigationMode_NonCompliant_Request_Count" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.noncompliant_requests"} 
                                      ]	
			           },
				   {"metric":"RequestCount" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.request_count"} 
                                      ]	
			           },
				   {"metric":"SpilloverCount" , "config": 
				      [
                                       {"stat":"Sum" , "measurement":"aws.elb.spillover_count"} 
                                      ]	
			           },
				   {"metric":"Latency" , "config": 
				      [
                                       {"stat":"Maximum" , "measurement":"aws.elb.latency.maximum"},
                                       {"stat":"Average" , "measurement":"aws.elb.latency.average"} 
                                      ]	
			           },
				   {"metric":"SurgeQueueLength" , "config": 
				      [
                                       {"stat":"Maximum" , "measurement":"aws.elb.surge_queue_length.maximum"},
                                       {"stat":"Minimum" , "measurement":"aws.elb.surge_queue_length.minimum"},
                                       {"stat":"Average" , "measurement":"aws.elb.surge_queue_length.average"} 
                                      ]	
			           },
				   {"metric":"HealthyHostCount" , "config": 
				      [
                                       {"stat":"Maximum" , "measurement":"aws.elb.healthy_host_count.maximum"},
                                       {"stat":"Minimum" , "measurement":"aws.elb.healthy_host_count.minimum"},
                                       {"stat":"Average" , "measurement":"aws.elb.healthy_host_count.average"} 
                                      ]	
			           },
				   {"metric":"UnHealthyHostCount" , "config": 
				      [
                                       {"stat":"Maximum" , "measurement":"aws.elb.unhealthy_host_count.maximum"},
                                       {"stat":"Minimum" , "measurement":"aws.elb.unhealthy_host_count.minimum"},
                                       {"stat":"Average" , "measurement":"aws.elb.unhealthy_host_count.average"} 
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
