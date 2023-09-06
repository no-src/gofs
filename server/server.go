package server

import (
	"fmt"
	"strings"
)

const (
	// RootGroupRoute the group route of root
	RootGroupRoute = "/"
	// DefaultRoute the route of default
	DefaultRoute = "/"
	// SourceRoutePrefix the route prefix of source
	SourceRoutePrefix = "/source/"
	// DestRoutePrefix the route prefix of dest
	DestRoutePrefix = "/dest/"
	// QueryRoute the route of query api
	QueryRoute = "/query"
	// LoginGroupRoute the group route of login
	LoginGroupRoute = "/login"
	// LoginIndexRoute the route of login index page
	LoginIndexRoute = "/index"
	// LoginIndexFullRoute the full route of login index page
	LoginIndexFullRoute = LoginGroupRoute + LoginIndexRoute
	// LoginSignInRoute the route of sign in api
	LoginSignInRoute = "/signin"
	// LoginSignInFullRoute the full route of sign in api
	LoginSignInFullRoute = LoginGroupRoute + LoginSignInRoute
	// WriteGroupRoute the group route of write api
	WriteGroupRoute = "/w"
	// PushRoute the route of push api
	PushRoute = "/push"
	// PushFullRoute the full route of push api
	PushFullRoute = WriteGroupRoute + PushRoute
	// ManageGroupRoute the group route of manage api
	ManageGroupRoute = "/manage"
	// ManageConfigRoute the route of manage config api
	ManageConfigRoute = "/config"
	// ManageReportRoute the route of report api
	ManageReportRoute = "/report"
	// PProfRoutePrefix the route prefix of pprof
	PProfRoutePrefix = "pprof"
)

const (
	// DefaultAddrHttps the default https address
	DefaultAddrHttps = ":443"
	// DefaultAddrHttp the default http address
	DefaultAddrHttp = ":80"
	// SchemeHttp the http scheme name
	SchemeHttp = "http"
	// SchemeHttps the https scheme name
	SchemeHttps = "https"
	// DefaultPortHttp the default port of http server
	DefaultPortHttp = 80
	// DefaultPortHttps the default port of https server
	DefaultPortHttps = 443
)

const (
	// SessionName the name of the session
	SessionName = "session_id"
	// SessionUser the key of the session user
	SessionUser = "user"
)

const (
	// ResourceTemplatePath the web server template resource path
	ResourceTemplatePath = "template/*"
)

// GenerateAddr generate http or https address
func GenerateAddr(scheme, host string, port int) string {
	addr := ""
	scheme = strings.ToLower(scheme)
	if scheme == SchemeHttp && port == DefaultPortHttp {
		addr = fmt.Sprintf("%s://%s", SchemeHttp, host)
	} else if scheme == SchemeHttps && port == DefaultPortHttps {
		addr = fmt.Sprintf("%s://%s", SchemeHttps, host)
	} else {
		addr = fmt.Sprintf("%s://%s:%d", scheme, host, port)
	}
	return addr
}
