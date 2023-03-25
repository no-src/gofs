package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/no-src/gofs/util/jsonutil"
)

func TestDuration_MarshalText(t *testing.T) {
	testCases := []struct {
		d      Duration
		expect string
	}{
		{Duration(-1), `"-1ns"`},
		{Duration(0), `"0s"`},
		{Duration(time.Second), `"1s"`},
		{Duration(time.Second * 2), `"2s"`},
		{Duration(time.Second * 70), `"1m10s"`},
		{Duration(time.Minute), `"1m0s"`},
		{Duration(time.Minute * 70), `"1h10m0s"`},
		{Duration(time.Hour), `"1h0m0s"`},
		{Duration(time.Hour * 25), `"25h0m0s"`},
	}

	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			data, err := jsonutil.Marshal(tc.d)
			if err != nil {
				t.Errorf("test duration marshal error =>%s", err)
				return
			}
			actual := string(data)
			if actual != tc.expect {
				t.Errorf("test duration marshal error, expect:%s, actual:%s", tc.expect, actual)
			}
		})
	}
}

func TestDuration_UnmarshalText(t *testing.T) {
	testCases := []struct {
		s      string
		expect time.Duration
	}{
		{`"-1ns"`, -1},
		{`"0s"`, 0},
		{`"1s"`, time.Second},
		{`"2s"`, time.Second * 2},
		{`"1m10s"`, time.Second * 70},
		{`"1m0s"`, time.Minute},
		{`"1h10m0s"`, time.Minute * 70},
		{`"1h0m0s"`, time.Hour},
		{`"25h0m0s"`, time.Hour * 25},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			data := []byte(tc.s)
			var actual Duration
			err := jsonutil.Unmarshal(data, &actual)
			if err != nil {
				t.Errorf("test duration unmarshal error =>%s", err)
				return
			}
			if actual.Duration() != tc.expect {
				t.Errorf("test duration unmarshal error, expect:%s, actual:%s", tc.expect.String(), actual.Duration().String())
			}
		})
	}
}

func TestDuration_UnmarshalText_ReturnError(t *testing.T) {
	testCases := []struct {
		s string
	}{
		{"hms"},
		{`"abc"`},
		{`"1x"`},
		{"@#$"},
		{""},
		{"	"},
		{"\n"},
		{"\r"},
		{"\t"},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			var actual Duration
			data := []byte(tc.s)
			if err := jsonutil.Unmarshal(data, &actual); err == nil {
				t.Errorf("test duration unmarshal error, expect to get an error but get nil")
			}
		})
	}
}

func TestDurationVar_WithDefaultValue(t *testing.T) {
	testCases := []struct {
		s            string
		defaultValue time.Duration
	}{
		{"-1ns", -1},
		{"0s", 0},
		{"1s", time.Second},
		{"2s", time.Second * 2},
		{"1m10s", time.Second * 70},
		{"1m0s", time.Minute},
		{"1h10m0s", time.Minute * 70},
		{"1h0m0s", time.Hour},
		{"25h0m0s", time.Hour * 25},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			actual := Duration(time.Second)
			testCommandLine.DurationVar(&actual, "core_test_duration_default"+tc.s, tc.defaultValue, "test duration")
			parseFlag()
			if actual.Duration() != tc.defaultValue {
				t.Errorf("test DurationVar with default value error, expect:%s, actual:%s", tc.defaultValue.String(), actual.Duration().String())
			}
		})
	}
}

func TestDurationVar(t *testing.T) {
	testCases := []struct {
		s      string
		expect time.Duration
	}{
		{"-1ns", -1},
		{"0s", 0},
		{"1s", time.Second},
		{"2s", time.Second * 2},
		{"1m10s", time.Second * 70},
		{"1m0s", time.Minute},
		{"1h10m0s", time.Minute * 70},
		{"1h0m0s", time.Hour},
		{"25h0m0s", time.Hour * 25},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			actual := Duration(time.Second)
			defaultValue := time.Minute * 2
			flagName := "core_test_duration" + tc.s
			testCommandLine.DurationVar(&actual, flagName, defaultValue, "test duration")
			parseFlag(fmt.Sprintf("-%s=%s", flagName, tc.s))
			if actual.Duration() != tc.expect {
				t.Errorf("test DurationVar error, expect:%s, actual:%s", tc.expect.String(), actual.Duration().String())
			}
		})
	}
}

// parseFlag parse the flags with specified arguments
func parseFlag(args ...string) {
	testCommandLine.Parse(args)
}
