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
|--region                       |AWS_REGION                |
|--profile                      |AWS_PROFILE                |


### Basic Usage
To retrieve all available metrics from a specific AWS service from a particular region is to specific the 
--namespace and --region cmdline arguments. Normally --region will be automatically detected as part of your 
AWS credentials profiles, but you may specify a different region if required.  Other arguments can be added
to optimize the metric response.

*Note:* The CloudWatch API uses a pagation strategy to limit the number of metrics returned in a single query. This check defaults to a limit of 1 page of results but this can be adjusted to meet your need. If the max number of result pages is too small, the check will return a warning status (return status 1) and will included a warning comment in the check output.

#### Example for AWS EC2 in region us-east-1 using stats and metric filter

```
sensu-cloudwatch-check --namespace "AWS/EC2" --region "us-east-1" --metric "StatusCheckFailed" --stats "Sum"
```
In this example the metric queries are limited to provide only the "Sum" statistica of "StatusCheckFailed" metric for the namespace "AWS/EC2"

#### Example for all metrics for specific AWS EC2 instance in region us-east-1 using dimension filter
```
sensu-cloudwatch-check --namespace "AWS/EC2" --region "us-east-1" --dimension-filters "InstanceId=i-0e302ffdcedaf34b1"
```

### Presets

### Custom Presets

### Exporting Preset Configuration
 

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
