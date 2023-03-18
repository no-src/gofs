package server

import "embed"

// Templates the web server templates
//
//go:embed template/*
var Templates embed.FS
