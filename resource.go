package gofs

import "embed"

//go:embed server/template/*
var Templates embed.FS

//go:embed version
var Version embed.FS
