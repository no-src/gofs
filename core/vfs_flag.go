package core

import (
	"flag"
	"fmt"
)

// VFSVar defines a VFS flag with specified name, default value, and usage string.
// The argument p points to a VFS variable in which to store the value of the flag.
func VFSVar(p *VFS, name string, value VFS, usage string) {
	flag.CommandLine.Var(newVFSValue(value, p), name, usage)
}

// VFSFlag defines a VFS flag with specified name, default value, and usage string.
// The return value is the address of a VFS variable that stores the value of the flag.
func VFSFlag(name string, value VFS, usage string) *VFS {
	p := new(VFS)
	VFSVar(p, name, value, usage)
	return p
}

// vfsValue implement the flag.Value and flag.Getter interface
type vfsValue VFS

func newVFSValue(val VFS, p *VFS) *vfsValue {
	*p = val
	return (*vfsValue)(p)
}

func (d *vfsValue) Set(s string) error {
	v := NewVFS(s)
	*d = vfsValue(v)
	return nil
}

func (d *vfsValue) Get() interface{} { return VFS(*d) }

func (d *vfsValue) String() string {
	vfs := (*VFS)(d)
	return fmt.Sprintf("%s:%s", vfs.Type(), vfs.Path())
}
