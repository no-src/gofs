package core

// UnmarshalText implement interface encoding.TextUnmarshaler
func (vfs *VFS) UnmarshalText(data []byte) error {
	return vfs.unmarshal(data)
}

// MarshalText implement interface encoding.TextMarshaler
func (vfs VFS) MarshalText() (text []byte, err error) {
	return []byte(vfs.original), nil
}

func (vfs *VFS) unmarshal(data []byte) error {
	v := NewVFS(string(data))
	*vfs = v
	return nil
}
