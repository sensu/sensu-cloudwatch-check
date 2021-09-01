package presets

import "fmt"

type ALB struct {
	Preset
}

// Overwrite the Preset Init function to enforce specific behavior
func (p *ALB) Ready() error {
	if p.verbose {
		fmt.Println("ALB::Ready Setting up presets")
	}

	// JSON Config String developed on 2021-08-18 from AWS Cloudwatch documentation
	//  Ref: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-cloudwatch-metrics.html
	measurementString :=
		`
{
  "namespace": "AWS/ApplicationELB",
  "dimension-filters": [],
  "measurements": [
    {
      "metric": "ActiveConnectionCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.active_connection_count"
        }
      ]
    },
    {
      "metric": "ClientTLSNegotiationErrorCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.client_tls_negotiation_error_count"
        }
      ]
    },
    {
      "metric": "TargetTLSNegotiationErrorCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.target_tls_negotiation_error_count"
        }
      ]
    },
    {
      "metric": "DesyncMitigationMode_NonCompliant_Request_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.desync_mitigation_mode_noncompliant_request_count"
        }
      ]
    },
    {
      "metric": "GrpcRequestCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.grpc_request_count"
        }
      ]
    },
    {
      "metric": "HTTP_Fixed_Response_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.http_fixed_response_count"
        }
      ]
    },
    {
      "metric": "HTTP_Redirect_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.http_redirect_count"
        }
      ]
    },
    {
      "metric": "HTTP_Redirect_Url_Limit_Exceeded_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.http_redirect_url_limit_exceeded_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_3XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_3xx_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_4XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_4xx_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_5XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_5xx_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_500_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_500_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_502_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_502_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_503_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_503_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_ELB_504_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_elb_504_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_Target_2XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_target_2xx_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_Target_3XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_target_3xx_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_Target_4XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_target_4xx_count"
        }
      ]
    },
    {
      "metric": "HTTPCode_Target_5XX_Count",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.httpcode_target_5xx_count"
        }
      ]
    },
    {
      "metric": "IPv6ProcessedBytes",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.ipv6_processed_bytes"
        }
      ]
    },
    {
      "metric": "IPv6RequestCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.ipv6_request_count"
        }
      ]
    },
    {
      "metric": "NewConnectionCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.new_connection_count"
        }
      ]
    },
    {
      "metric": "TargetConnectionErrorCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.target_connection_error_count"
        }
      ]
    },
    {
      "metric": "TargetResponseTime",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.alb.target_response_time.average"
        },
        {
          "stat": "p95",
          "measurement": "aws.alb.target_response_time.p95"
        },
        {
          "stat": "TM(:95)",
          "measurement": "aws.alb.target_response_time.tm95"
        }
      ]
    },
    {
      "metric": "NonStickyRequestCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.non_sticky_request_count"
        }
      ]
    },
    {
      "metric": "ProcessedBytes",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.processed_bytes"
        }
      ]
    },
    {
      "metric": "RejectedConnectionCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.rejected_connection_count"
        }
      ]
    },
    {
      "metric": "RequestCountPerTarget",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.request_count_per_target"
        }
      ]
    },
    {
      "metric": "ELBAuthError",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.elb_auth_error"
        }
      ]
    },
    {
      "metric": "ELBAuthFailure",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.elb_auth_failure"
        }
      ]
    },
    {
      "metric": "ELBAuthRefreshTokenSuccess",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.elb_auth_refresh_token_success"
        }
      ]
    },
    {
      "metric": "ELBAuthSuccess",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.elb_auth_success"
        }
      ]
    },
    {
      "metric": "ELBAuthUserClaimsSizeExceeded",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.elb_auth_user_claims_size_exceeded"
        }
      ]
    },
    {
      "metric": "ELBAuthLatency",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.elb_auth_latency"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.alb.elb_auth_latency.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.alb.elb_auth_latency.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.alb.elb_auth_latency.minimum"
        },
        {
          "stat": "Average",
          "measurement": "aws.alb.elb_auth_latency.average"
        }
      ]
    },
    {
      "metric": "LambdaInternalError",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.lambda_internal_error"
        }
      ]
    },
    {
      "metric": "LambdaTargetProcessedBytes",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.lambda_target_processed_bytes"
        }
      ]
    },
    {
      "metric": "LambdaUserError",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.lambda_user_error"
        }
      ]
    },
    {
      "metric": "RequestCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.request_count"
        }
      ]
    },
    {
      "metric": "RuleEvaluations",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.rule_evaluations"
        }
      ]
    },
    {
      "metric": "ConsumedLCUs",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.consumed_lcus.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.alb.consumed_lcus.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.alb.consumed_lcus.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.alb.consumed_lcus.minimum"
        },
        {
          "stat": "Average",
          "measurement": "aws.alb.consumed_lcus.average"
        }
      ]
    },
    {
      "metric": "HealthyHostCount",
      "config": [
        {
          "stat": "Maximum",
          "measurement": "aws.alb.healthy_host_count.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.alb.healthy_host_count.minimum"
        },
        {
          "stat": "Average",
          "measurement": "aws.alb.healthy_host_count.average"
        }
      ]
    },
    {
      "metric": "UnHealthyHostCount",
      "config": [
        {
          "stat": "Maximum",
          "measurement": "aws.alb.unhealthy_host_count.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.alb.unhealthy_host_count.minimum"
        },
        {
          "stat": "Average",
          "measurement": "aws.alb.unhealthy_host_count.average"
        }
      ]
    },
    {
      "metric": "DroppedInvalidHeaderRequestCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.dropped_invalid_header_requests.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.alb.dropped_invalid_header_requests.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.alb.dropped_invalid_header_requests.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.alb.dropped_invalid_header_requests.minimum"
        },
        {
          "stat": "Average",
          "measurement": "aws.alb.dropped_invalid_header_requests.average"
        }
      ]
    },
    {
      "metric": "ForwardedInvalidHeaderRequestCount",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.alb.forwarded_invalid_header_requests.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.alb.forwarded_invalid_header_requests.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.alb.forwarded_invalid_header_requests.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.alb.forwarded_invalid_header_requests.minimum"
        },
        {
          "stat": "Average",
          "measurement": "aws.alb.forwarded_invalid_header_requests.average"
        }
      ]
    }
  ]
}
`
	p.measurementString = measurementString
	p.BuildMeasurementConfig()

	return nil
}
