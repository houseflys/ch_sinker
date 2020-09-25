package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGjsonInt(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	var expected int64 = 1536813227
	result := metric.GetInt("its")
	assert.Equal(t, result, expected)
}

func TestGjsonArrayInt(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	arr := metric.GetArray("mp.a", "int").([]int64)
	expected := []int64{1, 2, 3}
	for i := range arr {
		assert.Equal(t, arr[i], expected[i])
	}
}

func TestGjsonStr(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	var expected string = "ws"
	result := metric.GetString("channel")
	assert.Equal(t, result, expected)
}

func TestGjsonArrayString(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	arr := metric.GetArray("mps.a", "string").([]string)
	expected := []string{"aa", "bb", "cc"}
	for i := range arr {
		assert.Equal(t, arr[i], expected[i])
	}
}

func TestGjsonFloat(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	var expected float64 = 0.11
	result := metric.GetFloat("percent")
	assert.Equal(t, result, expected)
}

func TestGjsonArrayFloat(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	arr := metric.GetArray("mp.f", "float").([]float64)
	expected := []float64{1.11, 2.22, 3.33}
	for i := range arr {
		assert.Equal(t, arr[i], expected[i])
	}
}

func TestGjsonElasticDateTime(t *testing.T) {
	parser := NewParser("gjson", nil, "", DefaultTSLayout)
	metric, _ := parser.Parse(jsonSample)

	// {"date": "2019-12-16T12:10:30Z"}
	// TZ=UTC date -d @1576498230 => Mon 16 Dec 2019 12:10:30 PM UTC
	var expected int64 = 1576498230
	result := metric.GetElasticDateTime("date")
	assert.Equal(t, result, expected)
}
