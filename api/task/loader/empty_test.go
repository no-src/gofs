package loader

import "testing"

func TestEmptyLoader(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Errorf("create empty loader error => %v", err)
		return
	}

	defer closeLoader(t, loader)

	err = loader.SaveConfig(nil)
	if err != nil {
		t.Errorf("save config error => %v", err)
		return
	}
	err = loader.SaveContent("", "")
	if err != nil {
		t.Errorf("save content error => %v", err)
		return
	}
	_, err = loader.LoadConfig()
	if err != nil {
		t.Errorf("load config error => %v", err)
		return
	}
	_, err = loader.LoadContent("")
	if err != nil {
		t.Errorf("load content error => %v", err)
		return
	}
}
