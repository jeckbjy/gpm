package gpm

import "os"

// Exists check dir of file exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
