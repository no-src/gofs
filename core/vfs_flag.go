package core

import (
	"flag"
	"fmt"
)

func VFSVar(p *VFS, name string, value VFS, usage string) {
	flag.CommandLine.Var(newVFSValue(value, p), name, usage)
}

func VFSFlag(name string, value VFS, usage string) *VFS {
	p := new(VFS)
	VFSVar(p, name, value, usage)
	return p
}

type vfsValue VFS

func newVFSValue(val VFS, p *VFS) *vfsValue {
	*p = val
	return (*vfsValue)(p)
}

func (d *vfsValue) Set(s string) error {
	v := NewDiskVFS(s)
	*d = vfsValue(v)
	return nil
}

func (d *vfsValue) Get() interface{} { return VFS(*d) }

func (d *vfsValue) String() string {
	vfs := (*VFS)(d)
	return fmt.Sprintf("%s:%s", vfs.Type(), vfs.Path())
}
