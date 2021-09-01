package presets

import "fmt"

type CLB struct {
	Preset
}

// Overwrite the Preset Init function to enforce specific behavior
func (p *CLB) Ready() error {
	if p.verbose {
		fmt.Println("CLB::Init Setting up clb preset")
	}

	// JSON Config String developed on 2021-08-18 from AWS Cloudwatch documentation
	//  Ref: https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/elb-cloudwatch-metrics.html#loadbalancing-metrics-clb
	measurementString := `{ "namespace" : "AWS/ELB",
                                "dimension-filters" : [ "LoadBalancerName", "AvailabilityZone" ],
                                "measurements" : 
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
	p.measurementString = measurementString
	p.BuildMeasurementConfig()

	return nil
}
