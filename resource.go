package gofs

import "embed"

// Templates the web server templates
//
//go:embed server/template/*
var Templates embed.FS
