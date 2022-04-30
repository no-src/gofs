package gofs

import "embed"

// Templates the web server templates
//go:embed server/template/*
var Templates embed.FS

// Commit the commit file records the last commit hash value, used by release
//go:embed version/commit
var Commit string

// GoVersion the go_version file records the go version info, used by release
//go:embed version/go_version
var GoVersion string
