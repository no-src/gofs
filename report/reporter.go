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

var (
	// GlobalReporter the global reporter
	GlobalReporter *Reporter
)

// Reporter collect the report data
type Reporter struct {
	enabled bool
	report  Report
	mu      sync.RWMutex
}

// GetReport get current report data
func (r *Reporter) GetReport() Report {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.report.CurrentTime = timeutil.Now()
	r.report.UpTime = core.Duration(r.report.CurrentTime.Sub(r.report.StartTime))
	return r.report
}

// PutConnection put a new connection
func (r *Reporter) PutConnection(addr string) {
	if !r.enabled {
		return
	}
	go r.putConnection(addr)
}

func (r *Reporter) putConnection(addr string) {
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

func (r *Reporter) connAddr(addr string) string {
	return strings.ToLower(addr)
}

// DeleteConnection delete a closed connection
func (r *Reporter) DeleteConnection(addr string) {
	if !r.enabled {
		return
	}
	go r.deleteConnection(addr)
}

func (r *Reporter) deleteConnection(addr string) {
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

// PutAuth put an auth info to update connection status
func (r *Reporter) PutAuth(addr string, user *auth.HashUser) {
	if !r.enabled {
		return
	}
	go r.putAuth(addr, user)
}

func (r *Reporter) putAuth(addr string, user *auth.HashUser) {
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

// PutEvent put a file change event
func (r *Reporter) PutEvent(event eventlog.Event) {
	if !r.enabled {
		return
	}
	go r.putEvent(event)
}

func (r *Reporter) putEvent(event eventlog.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.report.Events.Add(event)
	r.report.EventStat[event.Op]++
}

// PutApiStat put an access log of api
func (r *Reporter) PutApiStat(ip string) {
	if !r.enabled {
		return
	}
	go r.putApiStat(ip)
}

func (r *Reporter) putApiStat(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.report.ApiStat.AccessCount++
	r.report.ApiStat.VisitorStat[ip]++
}

// Enable enable or disable the Reporter
func (r *Reporter) Enable(enabled bool) {
	r.enabled = enabled
}

func init() {
	initGlobalReporter()
}

func initGlobalReporter() {
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
	GlobalReporter = &Reporter{
		report: report,
	}
}
