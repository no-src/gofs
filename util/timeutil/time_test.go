package timeutil

import (
	"bytes"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	unixTime := now.Unix()
	timeStr := now.Format(defaultTimeFormat)
	nt := NewTime(now)
	if nt.String() != timeStr {
		t.Errorf("test time error, expect:%s,actual:%s", timeStr, nt.String())
	}

	if nt.Unix() != unixTime {
		t.Errorf("test time error, expect:%d,actual:%d", unixTime, nt.Unix())
	}

	timeBytes, err := nt.MarshalText()
	if err != nil {
		t.Errorf("test time MarshalText error, get an error =>%v", err)
		return
	}
	if !bytes.Equal(timeBytes, []byte(timeStr)) {
		t.Errorf("test time MarshalText error, expect:%v, actual:%v", []byte(timeStr), timeBytes)
	}
}

func TestNow(t *testing.T) {
	now := Now()
	actual := now.String()
	expect := now.Time().Format(defaultTimeFormat)
	if actual != expect {
		t.Errorf("test Now error, expect:%s, actual:%s", expect, actual)
	}
}
