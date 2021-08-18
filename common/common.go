package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var (
	// Setup regexp for use with toSnakeCase
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func BuildLabelBase(m types.Metric) string {
	s := strings.Split(*m.Namespace, "/")
	labelString := ToSnakeCase(fmt.Sprintf("%v.%v.%v", ToSnakeCase(s[0]), ToSnakeCase(s[1]), ToSnakeCase(*m.MetricName)))
	return labelString
}

func DimString(dims []types.Dimension, region string) string {
	dimStrings := []string{}
	if len(dims) > 0 {
		for _, d := range dims {
			dimStrings = append(dimStrings, fmt.Sprintf(`%v="%v"`, *d.Name, *d.Value))
		}
	}
	if len(region) > 0 {
		dimStrings = append(dimStrings, fmt.Sprintf(`Region="%v"`, strings.TrimSpace(region)))
	}
	dimStr := strings.Join(dimStrings, ",")
	return dimStr
}
