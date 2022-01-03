package gofs

import "embed"

//go:embed server/template/*
var Templates embed.FS

//go:embed version/commit
var Commit string
