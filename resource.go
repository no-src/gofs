package gofs

import "embed"

// Templates the web server templates
//
//go:embed server/template/*
var Templates embed.FS

// Commit the commit file records the last commit hash value, used by release
//
//go:embed internal/version/commit
var Commit string
