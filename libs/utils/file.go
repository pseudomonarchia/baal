package utils

import "os"

// FileIsExists is used to determine whether the file exists
func FileIsExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
