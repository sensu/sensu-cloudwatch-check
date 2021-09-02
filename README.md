[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/sensu/sensu-cloudwatch-check)
![Go Test](https://github.com/sensu/sensu-cloudwatch-check/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/sensu/sensu-cloudwatch-check/workflows/goreleaser/badge.svg)


# sensu-cloudwatch-check

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
  - [Help output](#help-output)
  - [Environment Variables](#environment-variables)
  - [Basic Usage](#basic-usage)
  - [Important Commandline Arguments](#important-commendline-arguments)
  - [Presets](#presets)
  - [Custom Presets](#custom-presets)
  - [Exporting Preset Configuration](#exporting-preset-configuration)

- [Configuration](#configuration)
  - [AWS credentials](#AWS-credentials)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The sensu-cloudwatch-check is a [Sensu Check][6] that generates AWS service metrics from the AWS Cloudwatch API.

## Usage examples

### Help output

```
Sensu Cloudwatch Check

Usage:
  sensu-cloudwatch-check [flags]
  sensu-cloudwatch-check [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
      --config-files strings        comma separated list of AWS config files
      --credentials-files strings   comma separated list of AWS Credential files
      --profile string              AWS Credential Profile (for security use envvar AWS_PROFILE)
  -c, --config string               Use measurement configuration JSON string
  -N, --namespace string            Cloudwatch Metric Namespace
  -D, --dimension-filters strings    Comma separated list of AWS Cloudwatch Dimension Filters Ex: "Name, SecondName=SecondValue"
  -M, --metric string               Cloudwatch Metric Name
  -S, --stats strings               Comma separated list of AWS Cloudwatch Status Ex: "Average, Sum" (default [Average,Sum,SampleCount,Maximum,Minimum])
  -m, --max-pages int               Maximum number of result pages. A zero value will disable the limit (default 1)
  -o, --output-config               Output measurement configuration JSON string
  -p, --period-minutes int          Period in minutes for metrics statistic calculation (default 1)
  -P, --preset string               Preset Name (default "None")
      --recently-active             Only include metrics recently active in aprox last 3 hours
      --region string               AWS Region to use, (or set envvar AWS_REGION)
  -v, --verbose                     Enable verbose output
  -n, --dry-run                     Dryrun only list metrics, do not get metrics data
  -h, --help                        help for sensu-cloudwatch-check

```

### Environment Variables

|Argument                       |Environment Variable                 |
|-------------------------------|-------------------------------------|
|--region                       |AWS_REGION                           |
|--profile                      |AWS_PROFILE                          |
|--namespace                    | CLOUDWATCH_CHECK_NAMESPACE          |
|--metric-filter                | CLOUDWATCH_CHECK_METRIC_FILTER      | 
|--dimension-filters            | CLOUDWATCH_CHECK_DIMENSION_FILTERS  |
|--stats                        | CLOUDWATCH_CHECK_STATS              |
|--config                       | CLOUDWATCH_CHECK_CONFIG             |
|--preset                       | CLOUDWATCH_CHECK_PRESET             |
|--max-pages                    | CLOUDWATCH_CHECK_MAX_PAGES          |
|--period-minutes               | CLOUDWATCH_CHECK_PERIOD_MINUTES     |

  
### Basic Usage
To retrieve all available metrics from a specific AWS service from a particular region is to specific the 
--namespace and --region cmdline arguments. Normally --region will be automatically detected as part of your 
AWS credentials profiles, but you may specify a different region if required.  Other arguments can be added
to optimize the metric response.

*Note:* The CloudWatch API uses a pagation strategy to limit the number of metrics returned in a single query. This check defaults to a limit of 1 page of results but this can be adjusted to meet your need. If the max number of result pages is too small, the check will return a warning status (return status 1) and will included a warning comment in the check output.

*Note:* This check enforces a restriction the cloudwatch API query to limit the size of the cloudwatch query. You must either include the `--namespace` or `--metric-filter` option. See the Cloudwatch ListMetrics API documentation for details.

### Important Commandline Arguments
####  Namespace
The `--namespace` argument limits the Cloudwatch query to the given namespace (ex: AWS/EC2).
*Note:* Either `--namespace` or `--metric` is required

####  Metric Filter
The `--metric-filter` argument limits the cloudwatch query to a given metric name (ex: CPUUtilization)
*Note:* Either `--namespace` or `--metric` is required

####  Dimension Filters
The dimension filters is an array of strings representing dimension key or key=value that must match.
Allowed dimension filters are specific to AWS Namespace and metric. 
You should refer to the AWS service documentation for a specific service when choosing the dimension filters to use.

### Example for AWS EC2 in region us-east-1 using stats and metric filter

```
sensu-cloudwatch-check --namespace "AWS/EC2" --region "us-east-1" --metric "StatusCheckFailed" --stats "Sum"
```
In this example the metric queries are limited to provide only the "Sum" statistica of "StatusCheckFailed" metric for the namespace "AWS/EC2"

#### Example for all metrics for specific AWS EC2 instance in region us-east-1 using dimension filter
```
sensu-cloudwatch-check --namespace "AWS/EC2" --region "us-east-1" --dimension-filters "InstanceId=i-0e302ffdcedaf34b1"
```

### AWS CloudWatch Metrics Presets
This check comes with several presets for specific AWS Services.  These presets provide a curated subset of possible Cloudwatch statistics following an opinionated naming scheme.  These preset configs can be exported as a starting for for your own custom preset configuration (see below.) 

The list of existing service presets includes:

| Preset Name      |Description                 |
|------------------|----------------------------|
| ALB              | Preset Metrics for AWS Application Load Balancer            |
| CLB              | Preset Metrics for AWS Classic Load Balancer                |
| EC2              | Preset Metrics for AWS EC2                                  |
| CloudFront       | Preset Metrics for AWS CloudFront. Note: requires --region us-east-1 |

*Note:* The --dimension-filters and --metric-filter arguments can be used to further narrow the results
from the service presets.

### Custom Presets

You can define your own service preset by passing a json preset config string into the check using the `--config` option 
or `CLOUDWATCH_CHECK_CONFIG` envvar.


### Exporting Preset Configuration

The `--output-config` option can be used to generate the json representation of a metric query. 
This json can be editted as needed to customize the metric query. 
#### Export examples 

Output a simple `AWS/EC2` namespace metrics configuration
```
sensu-cloudwatch-check --namespace "AWS/EC2" --region "us-east-1" --output-config
```

Output a basic `AWS/EC2` namespace metrics configuration with an active dimension-filter
```
sensu-cloudwatch-check --namespace "AWS/EC2" --region "us-east-1" --dimension-filters "InstanceId=i-0e302ffdcedaf34b1" --output-config
```

Output the EC2 preset configuration
```
sensu-cloudwatch-check -P "EC2" --region "us-east-1" --output-config
```


## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add sensu/sensu-cloudwatch-check
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/sensu/sensu-cloudwatch-check].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-cloudwatch-check 
  namespace: default
spec:
  command: sensu-cloudwatch-check --example example_arg
  subscriptions:
  - system
  runtime_assets:
  - sensu/sensu-cloudwatch-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-cloudwatch-check repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/aws-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/aws-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu-community/aws-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
