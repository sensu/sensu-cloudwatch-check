module github.com/sensu/sensu-cloudwatch-check

go 1.15

require (
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.7.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.11.1
	github.com/sensu/sensu-go/api/core/v2 v2.3.0
	github.com/sensu/sensu-plugin-sdk v0.14.0
)
