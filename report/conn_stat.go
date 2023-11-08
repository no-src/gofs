package report

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/nsgo/timeutil"
)

// ConnStat the client connection info
type ConnStat struct {
	// Addr the client connection address
	Addr string `json:"addr"`
	// UserName the username of client
	UserName string `json:"username"`
	// Perm the permission of client
	Perm string `json:"perm"`
	// ConnectTime the connected time of client
	ConnectTime timeutil.Time `json:"connect_time"`
	// DisconnectTime the disconnected time of client
	DisconnectTime timeutil.Time `json:"disconnect_time"`
	// LifeTime the lifetime of a client, it is 0s always that if the client is online
	LifeTime core.Duration `json:"life_time"`
}
