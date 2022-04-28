## Intro
Want to contribute to the available AWS service presets? 
Or you just want develop your own service configuration and get a different combination of metrics than the plugin presets offer?
Then this document is for you.



### Make sure your local environment has the sensu-cloudwatch-check binary
You cant find the binaries for this command in the GitHub releases, or just checkout the repository and run "go build" 
to produce a local build of the binary you can use.

### Make sure your local environment has AWS auth ready to go.
For any development environment, I would expect you already have the aws cmdline client authenticated 
have generated a `.aws` directory with a default profile and associated credentials that let you access the cloudwatch API.
The `sensu-cloydwatch-check` executable uses the supported AWS SDK and will automatically detect your default profile and creds.


### The --output-config option is your friend
The `--output-config` option for the `sensu-cloudwatch-check` executabe will output a skeleton json configuration string that can be used 
with either `--config` or "CLOUDWATCH_CHECK_CONFIG" environment variable.  You can use this as a starting point for 
developing your own custom metrics for a given AWS namespace 


## Walkthrough example for AWS/ElastiCache
At the time of writing this, this plugin doesn't yet have an `AWS/ElastiCache` preset yet, so I'm going to use this as an example to walk-through

### Make sure you have an instances running in a region

The Cloudwatch API (afaict) only returns information about metrics for which there is data available. So for this method to work, 
its best to have an instance running already so we can build a configuration skeleton that has some metrics defined.

I have an instance of Redis ElasticCache in region `us-west-2` I'm going to use that to produce my skeleton measurement configuration
```
$ ./sensu-cloudwatch-check --namespace "AWS/ElastiCache" --region "us-west-2" -o >  elasticache_config.json
```

The `elasticache_config.json` file now has a valid measurement configuration for the `AWS/ElastiCache` namespace for region `us-west-2` that 
includes all the Cloudwatch metrics I was able to see using the CloudWatch API's ListMetrics call.  Each Cloudwatch metrics is actually 5 measurements, 
one each for the CloudWatch Statistics "Sum","SampleCount","Minimum","Maximum" and "Average."  

Usually you won't need each of those stats for each metric, so you'll want to edit the json file down based on what you need.  

### Configuration editting checklist
1. Remove the region, usually you want to specify the region via cmdline or envvar
2. remove unwanated measurements and/or measurement configs
3. rename the output measurements as desired. 

Here's an editted down version of my `elasticache_config.json` that only keeps a single measurement I'm interested in right now
```
{
  "namespace": "AWS/ElastiCache",
  "period-minutes": 1,
  "measurements": [
    {
      "metric": "FreeableMemory",
      "config": [
        {
          "stat": "Minimum",
          "measurement": "aws_elasticache_freeable_memory_minimum"
        }
      ]
    }
  ]
}

```
### Test editted measurement configuration
Let's save this json string as the CLOUDWATCH_CHECK_CONFIG envvar and then pull some metrics.
```
$  export CLOUDWATCH_CHECK_CONFIG=$(cat elasticache_config.json  | jq -cM .)
$ ./sensu-cloudwatch-check --region "us-west-2"
# HELP aws_elasticache_freeable_memory Namespace:AWS/ElastiCache MetricName:FreeableMemory Region:us-west-2
# TYPE aws_elasticache_freeable_memory gauge
aws_elasticache_freeable_memory_minimum{CacheClusterId="jspaleta-poc-redis"} 5.23587584e+08 1651185600000
aws_elasticache_freeable_memory_minimum{CacheClusterId="jspaleta-poc-redis",CacheNodeId="0001"} 5.23587584e+08 1651185600000
aws_elasticache_freeable_memory_minimum{} 5.23587584e+08 1651185600000
```

The output now only includes the measurements I have defined in the configuration json string.

### Writing your own custom checks.
Once you have the measurement configuration you like, you can set the CLOUDWATCH_CHECK_CONFIG envvar in a 
Sensu Check resource that uses the sensu-cloudwatch-check command,  and assuming the AWS authentication is correct, the agent
you'll be collecting the metrics in Prometheus exposition format.

### Submitting a preset.
We want to also have AWS service presets for the most commonly actionable metrics.  Presets are intended to be highly opinionated 
and to provide metrics that most AWS infra operators will want to build alerting around. 

But building a preset is a matter of encoding the json measurement into a new Preset object definition and calling the BuildMeasurementConfig function. 
Take a look a the presets/cloudfront.go as an example to follow
