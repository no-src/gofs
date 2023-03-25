package core

import (
	"fmt"
)

// VFSVar defines a VFS flag with specified name, default value, and usage string.
// The argument p points to a VFS variable in which to store the value of the flag.
func (f *FlagSet) VFSVar(p *VFS, name string, value VFS, usage string) {
	f.Var(newVFSValue(value, p), name, usage)
}

// VFSFlag defines a VFS flag with specified name, default value, and usage string.
// The return value is the address of a VFS variable that stores the value of the flag.
func (f *FlagSet) VFSFlag(name string, value VFS, usage string) *VFS {
	p := new(VFS)
	f.VFSVar(p, name, value, usage)
	return p
}

// vfsValue implement the flag.Value and flag.Getter interface
type vfsValue VFS

func newVFSValue(val VFS, p *VFS) *vfsValue {
	*p = val
	return (*vfsValue)(p)
}

// Set implement the Set function for the flag.Value interface
func (d *vfsValue) Set(s string) error {
	v := NewVFS(s)
	*d = vfsValue(v)
	return nil
}

// String implement the String function for the flag.Value interface
func (d *vfsValue) String() string {
	vfs := (*VFS)(d)
	return fmt.Sprintf("%s:%s", vfs.Type(), vfs.Path())
}
