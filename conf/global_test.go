package conf

import (
	"testing"
)

func TestSetGlobalConfig(t *testing.T) {
	testCases := []struct {
		name   string
		config *Config
		err    error
	}{
		{"normal address", &Config{FileServerAddr: ":8080"}, nil},
		{"empty address", &Config{FileServerAddr: ""}, nil},
		{"nil config", nil, errConfigIsNil},
		{"config exists", &Config{FileServerAddr: ":8080"}, errConfigExist},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetGlobalConfig(tc.config)
			if err != tc.err {
				t.Errorf("expect to get error %v, but get %v", tc.err, err)
			}
		})
	}
}

func TestGetGlobalConfig(t *testing.T) {
	if err := SetGlobalConfig(&Config{FileServerAddr: ":8088"}); err != nil {
		t.Errorf("call SetGlobalConfig error")
		return
	}
	testCases := []struct {
		name    string
		address string
		exist   bool
	}{
		{"normal address", ":8088", true},
		{"not exist config", ":8000", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := GetGlobalConfig(tc.address)
			exist := config != nil
			if exist != tc.exist {
				t.Errorf("expect to get config %v, but get %v", tc.exist, exist)
			}
		})
	}
}
