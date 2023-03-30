package report

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/eventlog"
	"github.com/no-src/gofs/internal/toplist"
	"github.com/no-src/gofs/internal/version"
	"github.com/no-src/gofs/util/timeutil"
)

// Reporter collect the report data
type Reporter interface {
	// GetReport get current report data
	GetReport() Report
	// PutConnection put a new connection
	PutConnection(addr string)
	// DeleteConnection delete a closed connection
	DeleteConnection(addr string)
	// PutAuth put an auth info to update connection status
	PutAuth(addr string, user *auth.HashUser)
	// PutEvent put a file change event
	PutEvent(event eventlog.Event)
	// PutApiStat put an access log of api
	PutApiStat(ip string)
	// Enable enable or disable the Reporter
	Enable(enabled bool)
}

type reporter struct {
	enabled bool
	report  Report
	mu      sync.Mutex
}

// NewReporter create an instance of the Reporter component
func NewReporter() Reporter {
	report := Report{
		StartTime: timeutil.Now(),
		Pid:       os.Getpid(),
		PPid:      os.Getppid(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
		GOVersion: runtime.Version(),
		Version:   version.VERSION,
		Commit:    version.Commit,
		Online:    make(map[string]*ConnStat),
		EventStat: make(map[string]uint64),
		ApiStat: ApiStat{
			VisitorStat: make(map[string]uint64),
		},
	}
	report.Events, _ = toplist.New(100)
	report.Hostname, _ = os.Hostname()
	return &reporter{
		report: report,
	}
}

func (r *reporter) GetReport() Report {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.report.CurrentTime = timeutil.Now()
	r.report.UpTime = core.Duration(r.report.CurrentTime.Sub(r.report.StartTime))
	return r.report
}

func (r *reporter) PutConnection(addr string) {
	if !r.enabled {
		return
	}
	go r.putConnection(addr)
}

func (r *reporter) putConnection(addr string) {
	now := timeutil.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	addr = r.connAddr(addr)
	stat := &ConnStat{
		Addr:        addr,
		ConnectTime: now,
	}
	r.report.Online[addr] = stat
}

func (r *reporter) connAddr(addr string) string {
	return strings.ToLower(addr)
}

func (r *reporter) DeleteConnection(addr string) {
	if !r.enabled {
		return
	}
	go r.deleteConnection(addr)
}

func (r *reporter) deleteConnection(addr string) {
	now := timeutil.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	addr = r.connAddr(addr)
	stat := r.report.Online[addr]
	if stat != nil {
		delete(r.report.Online, addr)
		stat.DisconnectTime = now
		stat.LifeTime = core.Duration(stat.DisconnectTime.Time().Sub(stat.ConnectTime.Time()))
		r.report.Offline = append(r.report.Offline, stat)
	}
}

func (r *reporter) PutAuth(addr string, user *auth.HashUser) {
	if !r.enabled {
		return
	}
	go r.putAuth(addr, user)
}

func (r *reporter) putAuth(addr string, user *auth.HashUser) {
	if len(addr) == 0 || user == nil {
		return
	}
	now := timeutil.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	addr = r.connAddr(addr)
	stat := r.report.Online[addr]
	if stat != nil {
		stat.IsAuth = true
		stat.UserName = user.UserNameHash
		stat.Perm = user.Perm.String()
		stat.AuthTime = now
	}
}

func (r *reporter) PutEvent(event eventlog.Event) {
	if !r.enabled {
		return
	}
	go r.putEvent(event)
}

func (r *reporter) putEvent(event eventlog.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.report.Events.Add(event)
	r.report.EventStat[event.Op]++
}

func (r *reporter) PutApiStat(ip string) {
	if !r.enabled {
		return
	}
	go r.putApiStat(ip)
}

func (r *reporter) putApiStat(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.report.ApiStat.AccessCount++
	r.report.ApiStat.VisitorStat[ip]++
}

func (r *reporter) Enable(enabled bool) {
	r.enabled = enabled
}
