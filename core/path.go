package core

import (
	"path/filepath"
	"strings"
)

type Path struct {
	origin string
	fsType VFSType
	bucket string
	base   string
}

func newPath(path string, fsType VFSType) Path {
	p := Path{
		origin: path,
		fsType: fsType,
	}
	p.parse()
	return p
}

// Bucket returns bucket name
func (p Path) Bucket() string {
	return p.bucket
}

// Base returns base path
func (p Path) Base() string {
	return p.base
}

// String return the origin path
func (p Path) String() string {
	return p.origin
}

func (p *Path) parse() {
	// maybe the remote os is different from the current os, force convert remote path to slash
	if p.fsType != MinIO {
		p.origin = filepath.ToSlash(filepath.Clean(p.origin))
	}
	p.base = p.origin

	if p.fsType == MinIO {
		// protocol => bucket:path
		// example => mybucket:/workspace
		if strings.Contains(p.origin, ":") && !strings.HasPrefix(p.origin, "/") {
			list := strings.Split(p.origin, ":")
			p.bucket = list[0]
			p.base = list[1]
		} else {
			p.bucket = p.origin
			p.base = ""
		}
	}
}
