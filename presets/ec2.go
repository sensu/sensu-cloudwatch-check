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
          "stat": "Sum",
          "measurement": "aws.ec2.metadata_no_token"
        }
      ]
    },
    {
      "metric": "StatusCheckFailed_System",
      "config": [
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.status_check_failed_system"
        }
      ]
    },
    {
      "metric": "StatusCheckFailed_Instance",
      "config": [
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.status_check_failed__instance"
        }
      ]
    },
    {
      "metric": "StatusCheckFailed",
      "config": [
        {
          "stat": "Maximum",
          "measurement": "aws.ec2.status_check_failed"
        }
      ]
    },
    {
      "metric": "DiskWriteBytes",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_write_bytes"
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
          "stat": "Sum",
          "measurement": "aws.ec2.disk_read_ops"
        }
      ]
    },
    {
      "metric": "DiskWriteOps",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_write_ops"
        }
      ]
    },
    {
      "metric": "DiskReadBytes",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.ec2.disk_read_bytes"
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
