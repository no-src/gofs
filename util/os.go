package util

import "runtime"

// IsWindows is the Windows operating system
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
