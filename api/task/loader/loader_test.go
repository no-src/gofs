package loader

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/no-src/gofs/conf"
)

func TestLoader(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"memory:"},
		{"buntdb://buntdb.db"},
		{"buntdb://:memory:"},
		{"file://test-task-file.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			loader, err := NewLoader(tc.path)
			if err != nil {
				t.Errorf("create loader error => %v", err)
				return
			}

			defer closeLoader(t, loader)

			wc, err := getTaskConfig()
			if err != nil {
				t.Errorf("get test task config error => %v", err)
				return
			}
			err = loader.SaveConfig(wc)
			if err != nil {
				t.Errorf("save config error => %v", err)
				return
			}

			conf := "local-disk-sync.yaml"
			wContent := "source: ./source\ndest: ./dest"
			err = loader.SaveContent(conf, wContent)
			if err != nil {
				t.Errorf("save content error => %v", err)
				return
			}

			rc, err := loader.LoadConfig()
			if err != nil {
				t.Errorf("load config error => %v", err)
				return
			}
			if len(rc.Tasks) != len(wc.Tasks) {
				t.Errorf("load config expect to get tasks %d, but actual get tasks %d", len(wc.Tasks), len(rc.Tasks))
				return
			}

			rContent, err := loader.LoadContent(conf)
			if err != nil {
				t.Errorf("load content error => %v", err)
				return
			}
			if rContent != wContent {
				t.Errorf("load content expect to get %s, but actual get %s", wContent, rContent)
				return
			}
		})
	}

	// clear testdata
	os.Remove("test-task-file.yaml")
	os.Remove("local-disk-sync.yaml")
	os.Remove("buntdb.db")
}

func TestLoader_SaveConfig_WithDuplicateTasks(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"memory:"},
		{"file://test-task-file.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			loader, err := NewLoader(tc.path)
			if err != nil {
				t.Errorf("create loader error => %v", err)
				return
			}

			defer closeLoader(t, loader)

			wc, err := getTaskConfigWithDuplicateTasks()
			if err != nil {
				t.Errorf("get test task config error => %v", err)
				return
			}
			err = loader.SaveConfig(wc)
			if !errors.Is(err, errDuplicateTask) {
				t.Errorf("save config expect get error %v, but actual get %v", errDuplicateTask, err)
			}
		})
	}
}

func TestLoader_SaveConfig_WithNilConfig(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"memory:"},
		{"file://test-task-file.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			loader, err := NewLoader(tc.path)
			if err != nil {
				t.Errorf("create loader error => %v", err)
				return
			}

			defer closeLoader(t, loader)

			err = loader.SaveConfig(nil)
			if !errors.Is(err, errNilTaskConfig) {
				t.Errorf("save config expect get error %v, but actual get %v", errNilTaskConfig, err)
			}
		})
	}
}

func TestLoader_SaveConfig_WithUnsupportedFormat(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"file://test-task-file.x"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			loader, err := NewLoader(tc.path)
			if err != nil {
				t.Errorf("create loader error => %v", err)
				return
			}

			defer closeLoader(t, loader)

			expect := "unsupported config format"
			err = loader.SaveConfig(&TaskConfig{})
			if err == nil || err.Error() != expect {
				t.Errorf("save config expect get error %s, but actual get %v", expect, err)
			}
		})
	}
}

func TestLoader_LoadConfig_ReturnError(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"memory:"},
		{"file://not-exist.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			loader, err := NewLoader(tc.path)
			if err != nil {
				t.Errorf("create loader error => %v", err)
				return
			}

			defer closeLoader(t, loader)

			_, err = loader.LoadConfig()
			if err == nil {
				t.Errorf("load config expect to get an error, but actual get nil")
			}
		})
	}
}

func TestCacheLoader_WithUnsupportedDriver(t *testing.T) {
	_, err := NewLoader("unsupported://config.json")
	if err == nil {
		t.Errorf("new loader expect to get an error, but actual get nil")
	}
}

func TestFileLoader_LoadConfig_WithDuplicateTasks(t *testing.T) {
	wc, err := getTaskConfigWithDuplicateTasks()
	if err != nil {
		t.Errorf("get test task config error => %v", err)
		return
	}
	data, err := conf.ToString(".yaml", wc)
	if err != nil {
		t.Errorf("get task config string error => %v", err)
		return
	}
	duplicateTaskFile := "duplicate-tasks.yaml"
	err = os.WriteFile(duplicateTaskFile, []byte(data), fs.ModePerm)
	if err != nil {
		t.Errorf("create duplicate task config error => %v", err)
		return
	}
	defer os.Remove(duplicateTaskFile)
	loader, err := NewLoader("file://" + duplicateTaskFile)
	if err != nil {
		t.Errorf("create loader error => %v", err)
		return
	}

	defer closeLoader(t, loader)

	_, err = loader.LoadConfig()
	if !errors.Is(err, errDuplicateTask) {
		t.Errorf("load config expect get error %v, but actual get %v", errDuplicateTask, err)
	}
}

func TestCacheLoader_LoadConfig_WithDuplicateTasks(t *testing.T) {
	wc, err := getTaskConfigWithDuplicateTasks()
	if err != nil {
		t.Errorf("get test task config error => %v", err)
		return
	}
	loader, err := NewLoader("memory:")
	if err != nil {
		t.Errorf("create loader error => %v", err)
		return
	}

	defer closeLoader(t, loader)

	realLoader := loader.(*cacheLoader)
	err = realLoader.cache.Set(realLoader.confKey, wc, 0)
	if err != nil {
		t.Errorf("write cache data error => %v", err)
		return
	}
	_, err = loader.LoadConfig()
	if !errors.Is(err, errDuplicateTask) {
		t.Errorf("load config expect get error %v, but actual get %v", errDuplicateTask, err)
	}
}

func getTaskConfig() (c *TaskConfig, err error) {
	err = conf.Parse("../../testdata/tasks.yaml", &c)
	return
}

func getTaskConfigWithDuplicateTasks() (c *TaskConfig, err error) {
	c, err = getTaskConfig()
	if err != nil {
		return
	}
	c.Tasks = append(c.Tasks, c.Tasks...)
	return
}

func closeLoader(t *testing.T, loader Loader) {
	if err := loader.Close(); err != nil {
		t.Errorf("close the loader error => %v", err)
	}
}
