package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/no-src/gofs/api/task/loader"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/flag"
)

// Dispatcher the task dispatcher interface
type Dispatcher interface {
	// Acquire try to acquire a task
	Acquire(client *ClientInfo, ip string) (task *TaskInfo, err error)
}

type dispatcher struct {
	loader loader.Loader
	m      sync.Map
}

func newDispatcher(taskConf string) (Dispatcher, error) {
	loader, err := loader.NewLoader(taskConf)
	if err != nil {
		return nil, err
	}
	return &dispatcher{
		loader: loader,
	}, nil
}

func (d *dispatcher) Acquire(client *ClientInfo, ip string) (task *TaskInfo, err error) {
	if client == nil {
		return nil, errors.New("invalid client info")
	}
	tasks, err := d.loadTasks()
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		if d.acquired(client, t) {
			continue
		}

		if !d.checkAllowIP(t.AllowIP, ip) {
			continue
		}

		if !d.checkLabels(t.Labels, client.Labels) {
			continue
		}

		content, err := d.loader.LoadContent(t.Conf)
		if err != nil {
			return nil, err
		}
		ext := filepath.Ext(t.Conf)
		// get default config
		c := flag.ParseFlags([]string{os.Args[0], "-conf="})
		// override config
		if err = conf.ParseContent([]byte(content), ext, &c); err != nil {
			return nil, err
		}
		if content, err = conf.ToString(ext, c); err != nil {
			return nil, err
		}
		d.markAcquired(client, t)
		return &TaskInfo{
			Name:    t.Name,
			Ext:     ext,
			Content: content,
		}, nil
	}

	return nil, nil
}

func (d *dispatcher) loadTasks() ([]*loader.TaskItem, error) {
	c, err := d.loader.LoadConfig()
	if err != nil {
		return nil, err
	}
	return c.Tasks, nil
}

func (d *dispatcher) checkAllowIP(allowIP []string, clientIP string) bool {
	if len(allowIP) == 0 {
		return true
	}
	for _, ip := range allowIP {
		if ip == clientIP {
			return true
		}
	}
	return false
}

func (d *dispatcher) checkLabels(taskLabels []string, clientLabels []string) bool {
	if len(taskLabels) == 0 {
		return true
	}
	if len(taskLabels) > len(clientLabels) {
		return false
	}
	for _, serverLabel := range taskLabels {
		if !d.contain(clientLabels, serverLabel) {
			return false
		}
	}
	return true
}

func (d *dispatcher) contain(list []string, s string) bool {
	for _, str := range list {
		str = strings.TrimSpace(str)
		if len(str) > 0 && str == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}

func (d *dispatcher) acquired(c *ClientInfo, t *loader.TaskItem) bool {
	k := fmt.Sprintf("%s:%s", c.GetClientId(), t.Name)
	_, ok := d.m.Load(k)
	return ok
}

func (d *dispatcher) markAcquired(c *ClientInfo, t *loader.TaskItem) {
	k := fmt.Sprintf("%s:%s", c.GetClientId(), t.Name)
	d.m.Store(k, t)
}
