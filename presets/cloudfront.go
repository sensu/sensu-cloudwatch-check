package presets

import "fmt"

type CloudFront struct {
	Preset
}

// Overwrite the Preset Ready function to enforce specific behavior
func (p *CloudFront) Ready() error {
	if p.verbose {
		fmt.Println("CloudFront::Ready Setting up presets")
	}

	// JSON Config String developed on 2021-08-18 from AWS Cloudwatch documentation
	//  Ref: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/programming-cloudwatch-metrics.html#cloudfront-metrics-global-values
	// NOTE: OriginLatency metric not supported as it required used of AWS API ExtendedStatistics
	measurementString :=
		`
{
  "namespace": "AWS/CloudFront",
  "measurements": [
    {
      "metric": "BytesUploaded",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.cloud_front.bytes_uploaded"
        }
      ]
    },
    {
      "metric": "5xxErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.5xx_error_rate"
        }
      ]
    },
    {
      "metric": "502ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.502_error_rate.average"
        }
      ]
    },
    {
      "metric": "503ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.503_error_rate.average"
        }
      ]
    },
    {
      "metric": "504ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.504_error_rate.average"
        }
      ]
    },
    {
      "metric": "CacheHitRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.cache_hit_rate"
        }
      ]
    },
    {
      "metric": "Requests",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.cloud_front.requests"
        }
      ]
    },
    {
      "metric": "TotalErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.total_error_rate"
        }
      ]
    },
    {
      "metric": "BytesDownloaded",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.cloud_front.bytes_downloaded.sum"
        }
      ]
    },
    {
      "metric": "Invocations",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.cloud_front.invocations"
        }
      ]
    },
    {
      "metric": "ValidationErrors",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.cloud_front.validation_errors"
        }
      ]
    },
    {
      "metric": "ExecutionErrors",
      "config": [
        {
          "stat": "Sum",
          "measurement": "aws.cloud_front.execution_errors"
        }
      ]
    },
    {
      "metric": "ExecutionTime",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.execution_time"
        }
      ]
    },
    {
      "metric": "4xxErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.4xx_error_rate.average"
        }
      ]
    },
    {
      "metric": "401ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.401_error_rate.average"
        }
      ]
    },
    {
      "metric": "402ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.402_error_rate.average"
        }
      ]
    },
    {
      "metric": "403ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.403_error_rate.average"
        }
      ]
    },
    {
      "metric": "404ErrorRate",
      "config": [
        {
          "stat": "Average",
          "measurement": "aws.cloud_front.404_error_rate.average"
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
