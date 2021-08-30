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
	s := strings.Split("Unknown/Namespace", "/")
	if m.Namespace != nil {
		s = strings.Split(*m.Namespace, "/")
	}
	for i, _ := range s {
		s[i] = ToSnakeCase(s[i])
	}

	labelString := ToSnakeCase(fmt.Sprintf("%v.%v", strings.Join(s, "."), ToSnakeCase(*m.MetricName)))
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

func BuildDimensionFilters(input []string) ([]types.DimensionFilter, error) {
	output := []types.DimensionFilter{}
	for _, item := range input {
		segments := strings.Split(strings.TrimSpace(item), "=")
		if len(segments) < 1 || len(segments) > 2 {
			return nil, fmt.Errorf("Error parsing dimension filters")
		}
		filter := types.DimensionFilter{Name: &segments[0]}
		if len(segments) > 1 {
			filter.Value = &segments[1]
		}
		output = append(output, filter)
	}
	return output, nil
}

func RemoveDuplicateStrings(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
