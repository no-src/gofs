package core

import (
	"flag"
	"testing"
	"time"

	"github.com/no-src/gofs/util/jsonutil"
)

func TestDurationMarshalText(t *testing.T) {
	d := Duration(time.Second)
	data, err := jsonutil.Marshal(d)
	if err != nil {
		t.Errorf("test duration marshal error =>%s", err)
		return
	}
	expect := `"1s"`
	actual := string(data)
	if actual != expect {
		t.Errorf("test duration marshal error, expect:%s, actual:%s", expect, actual)
	}
}

func TestDurationUnmarshalText(t *testing.T) {
	var actual Duration
	data := []byte(`"1s"`)
	err := jsonutil.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("test duration unmarshal error =>%s", err)
		return
	}
	expect := time.Second
	if actual.Duration() != expect {
		t.Errorf("test duration unmarshal error, expect:%s, actual:%s", expect.String(), actual.Duration().String())
	}
}

func TestDurationUnmarshalTextError(t *testing.T) {
	var actual Duration
	data := []byte(`"1x"`)
	err := jsonutil.Unmarshal(data, &actual)
	if err == nil {
		t.Errorf("test duration unmarshal should be error")
		return
	}
}

func TestDurationVarDefaultValue(t *testing.T) {
	defaultValue := time.Minute * 2
	expect := defaultValue
	actual := Duration(time.Second)
	DurationVar(&actual, "core_test_duration_default", defaultValue, "test duration")
	parseFlag()
	if actual.Duration() != expect {
		t.Errorf("test DurationVar with default value error, expect:%s, actual:%s", expect.String(), actual.Duration().String())
	}
}

func TestDurationVar(t *testing.T) {
	expect := time.Second * 3
	actual := Duration(time.Second)
	defaultValue := time.Minute * 2
	DurationVar(&actual, "core_test_duration", defaultValue, "test duration")
	parseFlag("-core_test_duration=3s")
	if actual.Duration() != expect {
		t.Errorf("test DurationVar error, expect:%s, actual:%s", expect.String(), actual.Duration().String())
	}
}

// parseFlag parse the flags with specified arguments
func parseFlag(args ...string) {
	flag.CommandLine.Parse(args)
}
