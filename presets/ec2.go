package presets

import "fmt"

type EC2 struct {
	Preset
}

// Overwrite the Preset Ready function to enforce specific behavior
func (p *EC2) Ready() error {
	if p.verbose {
		fmt.Println("EC2::Ready Setting up presets")
	}

	// JSON Config String developed on 2021-08-18 from AWS Cloudwatch documentation
	//  Ref: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html
	measurementString :=
		`
{
  "namespace": "AWS/EC2",
  "measurements": [
    {
      "metric": "CPUSurplusCreditsCharged",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.cpu_surplus_credits_charged.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.cpu_surplus_credits_charged.sum"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.cpu_surplus_credits_charged.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.cpu_surplus_credits_charged.minimum"
        }
      ]
    },
    {
      "metric": "CPUCreditUsage",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.cpu_credit_usage.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.cpu_credit_usage.sum"
        }
      ]
    },
    {
      "metric": "CPUSurplusCreditBalance",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.cpu_surplus_credit_balance.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.cpu_surplus_credit_balance.sum"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.cpu_surplus_credit_balance.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.cpu_surplus_credit_balance.minimum"
        }
      ]
    },
    {
      "metric": "CPUCreditBalance",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.cpu_credit_balance.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.cpu_credit_balance.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.cpu_credit_balance.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.cpu_credit_balance.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.cpu_credit_balance.minimum"
        }
      ]
    },
    {
      "metric": "NetworkIn",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.network_in"
        }
      ]
    },
    {
      "metric": "NetworkPacketsIn",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.network_packets_in"
        }
      ]
    },
    {
      "metric": "CPUUtilization",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.cpu_utilization.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.cpu_utilization.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.cpu_utilization.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.cpu_utilization.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.cpu_utilization.minimum"
        }
      ]
    },
    {
      "metric": "MetadataNoToken",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.metadata_no_token.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.metadata_no_token.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.metadata_no_token.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.metadata_no_token.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.metadata_no_token.minimum"
        }
      ]
    },
    {
      "metric": "StatusCheckFailed_System",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.status_check_failed_system.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.status_check_failed_system.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.status_check_failed_system.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.status_check_failed_system.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.status_check_failed_system.minimum"
        }
      ]
    },
    {
      "metric": "StatusCheckFailed_Instance",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.status_check_failed__instance.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.status_check_failed__instance.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.status_check_failed__instance.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.status_check_failed__instance.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.status_check_failed__instance.minimum"
        }
      ]
    },
    {
      "metric": "StatusCheckFailed",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.status_check_failed.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.status_check_failed.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.status_check_failed.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.status_check_failed.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.status_check_failed.minimum"
        }
      ]
    },
    {
      "metric": "DiskWriteBytes",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.disk_write_bytes.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_write_bytes.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.disk_write_bytes.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.disk_write_bytes.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.disk_write_bytes.minimum"
        }
      ]
    },
    {
      "metric": "NetworkPacketsOut",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.network_packets_out"
        }
      ]
    },
    {
      "metric": "DiskReadOps",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.disk_read_ops.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_read_ops.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.disk_read_ops.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.disk_read_ops.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.disk_read_ops.minimum"
        }
      ]
    },
    {
      "metric": "DiskWriteOps",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.disk_write_ops.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_write_ops.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.disk_write_ops.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.disk_write_ops.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.disk_write_ops.minimum"
        }
      ]
    },
    {
      "metric": "DiskReadBytes",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.ec2.disk_read_bytes.average"
        },
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_read_bytes.sum"
        },
        {
          "stat": "SampleCount",
          "measurement": "aws.ec2.disk_read_bytes.sample_count"
        },
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.disk_read_bytes.maximum"
        },
        {
          "stat": "Minimum",
          "measurement": "aws.ec2.disk_read_bytes.minimum"
        }
      ]
    },
    {
      "metric": "NetworkOut",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.network_out"
        }
      ]
    }
  ]
}
`
	p.measurementString = measurementString
	err := p.BuildMeasurementConfig()

	return err
}
