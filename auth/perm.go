package auth

import (
	"strings"
)

// Perm the basic permission info
type Perm string

const (
	// ReadPerm the permission of read
	ReadPerm = "r"
	// WritePerm the permission of write
	WritePerm = "w"
	// ExecutePerm the permission of execute
	ExecutePerm = "x"
	// DefaultPerm the default permission
	DefaultPerm = ReadPerm
	// FullPerm the permission of read, write and execute
	FullPerm = "rwx"
)

// String return perm string
func (p Perm) String() string {
	return string(p)
}

// R have the permission to read or not
func (p Perm) R() bool {
	return strings.Contains(p.String(), ReadPerm)
}

// W have the permission to write or not
func (p Perm) W() bool {
	return strings.Contains(p.String(), WritePerm)
}

// X have the permission to execute or not
func (p Perm) X() bool {
	return strings.Contains(p.String(), ExecutePerm)
}

// CheckTo check the target permission whether accord with current permission
// if the current permission is invalid, return false always
func (p Perm) CheckTo(t Perm) bool {
	if !p.IsValid() {
		return false
	}
	if !t.IsValid() || (p.R() && !t.R()) || (p.W() && !t.W()) || (p.X() && !t.X()) {
		return false
	}
	return true
}

// IsValid is a valid permission or not
func (p Perm) IsValid() bool {
	return len(p.String()) > 0 && (p.R() || p.W() || p.X())
}

// ToPermWithDefault convert a perm string to Perm
// defaultPerm if the perm is empty, replace with the defaultPerm
func ToPermWithDefault(perm string, defaultPerm string) (p Perm) {
	perm = strings.TrimSpace(perm)
	if len(perm) == 0 {
		perm = defaultPerm
	}
	return ToPerm(perm)
}

// ToPerm convert a perm string to Perm
func ToPerm(perm string) (p Perm) {
	perm = strings.TrimSpace(perm)
	permLen := len(perm)
	if permLen == 0 || permLen > 3 {
		return p
	}
	perm = strings.ToLower(perm)
	r, w, x := false, false, false
	for i := 0; i < permLen; i++ {
		c := perm[i : i+1]
		switch c {
		case ReadPerm:
			r = true
			break
		case WritePerm:
			w = true
			break
		case ExecutePerm:
			x = true
			break
		default:
			return p
		}
	}
	if r {
		p += ReadPerm
	}
	if w {
		p += WritePerm
	}
	if x {
		p += ExecutePerm
	}
	return p
}
