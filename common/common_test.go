package common

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
)

var (
	enableQuiet = false
)

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	if enableQuiet {
		os.Stdout = null
		os.Stderr = null
		log.SetOutput(null)
	}
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
func TestToSnakeCase(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	input := "ThisTestIsFun"
	expected := "this_test_is_fun"
	output := ToSnakeCase(input)
	assert.Equal(expected, output)
	assert.NotEqual(input, output)
}

func TestBuildLabelBase(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	namespace := "Test"
	name := "MetricName"
	expected := "test.metric_name"
	metric := types.Metric{
		Namespace:  &namespace,
		MetricName: &name,
	}
	label := BuildLabelBase(metric)
	namespace = "test/test"
	name = "MetricName"
	expected = "test.test.metric_name"
	metric = types.Metric{
		Namespace:  &namespace,
		MetricName: &name,
	}
	label = BuildLabelBase(metric)
	assert.Equal(expected, label)
}

func TestDimString(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	name := "hey"
	value := "you"
	expected := `"hey=you"`
	input := []types.Dimension{
		types.Dimension{
			Name:  &name,
			Value: &value,
		},
	}
	output := DimString(input)
	assert.Equal(expected, output)
}

func TestBuildDimensionFilters(t *testing.T) {
	defer quiet()()
	assert := assert.New(t)
	input := []string{
		"hey=you",
		"what",
	}
	output, err := BuildDimensionFilters(input)
	assert.NoError(err)
	assert.Equal(2, len(output))
}
