package conf

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/util/yamlutil"
)

func TestConfig_ToArgs(t *testing.T) {
	c := Config{
		SyncOnce:      true,
		ChunkSize:     1024,
		SyncDelayTime: core.Duration(time.Second * 3),
		Source:        core.NewVFS("rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1"),
	}
	args, err := c.ToArgs()
	if err != nil {
		t.Errorf("parse config to arguments error, %v", err)
		return
	}
	testCases := []struct {
		k string
		v string
	}{
		{"sync_once", "true"},
		{"chunk_size", "1024"},
		{"sync_delay_time", "3s"},
		{"source", "rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1"},
		{"retry_async", "false"},
		{"retry_count", "0"},
		{"retry_wait", "0s"},
		{"dest", "\"\""},
	}
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		k := kv[0]
		v := kv[1]
		for _, tc := range testCases {
			if k == tc.k && v != tc.v {
				t.Errorf("k=%s expect to get [%s],but get [%s]", k, tc.v, v)
			}
		}
	}
}

func TestConfig_ToArgs_MarshalError(t *testing.T) {
	var c Config
	errMarshal := errors.New("yaml marshal error mock")
	m := yamlutil.Marshal
	defer func() {
		yamlutil.Marshal = m
	}()
	yamlutil.Marshal = func(v any) ([]byte, error) {
		return nil, errMarshal
	}
	_, err := c.ToArgs()
	if err != errMarshal {
		t.Errorf("expect to get error %v, but actual get %v", errMarshal, err)
	}
}
