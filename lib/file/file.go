package file

import "os"

// IsExists is used to determine whether the file exists
func IsExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
