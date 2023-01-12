package report

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/internal/toplist"
	"github.com/no-src/gofs/util/timeutil"
)

// Report the program report data
type Report struct {
	// CurrentTime returns the current server time
	CurrentTime timeutil.Time `json:"current_time"`
	// StartTime returns the server start time
	StartTime timeutil.Time `json:"start_time"`
	// UpTime returns the server up time
	UpTime core.Duration `json:"up_time"`
	// Pid returns the process id of the caller
	Pid int `json:"pid"`
	// PPid returns the process id of the caller's parent
	PPid int `json:"ppid"`
	// Hostname returns the host name reported by the kernel
	Hostname string `json:"hostname"`
	// GOOS is the running program's operating system target
	GOOS string `json:"go_os"`
	// GOARCH is the running program's architecture target
	GOARCH string `json:"go_arch"`
	// GOVersion returns the Go tree's version string
	GOVersion string `json:"go_version"`
	// Version returns the version info of the gofs
	Version string `json:"version"`
	// Commit returns last commit hash value of the gofs
	Commit string `json:"commit"`
	// Online returns the client connection info that is online
	Online map[string]*ConnStat `json:"online"`
	// Offline returns the client connection info that is offline
	Offline []*ConnStat `json:"offline"`
	// Events returns some latest file change events
	Events *toplist.TopList `json:"events"`
	// EventStat returns the statistical data of file change events
	EventStat EventStat `json:"event_stat"`
	// ApiStat returns the statistical data of api access info
	ApiStat ApiStat `json:"api_stat"`
}
