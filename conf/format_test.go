package conf

import (
	"path/filepath"
	"testing"
)

func TestFormat(t *testing.T) {
	expert := "json"
	actual := JsonFormat.Name()
	if actual != expert {
		t.Errorf("test json format error expert:%s, actual:%s", expert, actual)
	}

	expert = "yaml"
	actual = YamlFormat.Name()
	if actual != expert {
		t.Errorf("test yaml format error expert:%s, actual:%s", expert, actual)
	}

	if !JsonFormat.MatchExt(filepath.Ext(jsonConfigPath)) {
		t.Errorf("match json confile error => %s", jsonConfigPath)
	}

	if !YamlFormat.MatchExt(filepath.Ext(yamlConfigPath)) {
		t.Errorf("match yarm confile error => %s", yamlConfigPath)
	}
}
