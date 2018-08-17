package file

import "os"

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsExist(err) {
		return true
	}

	return false
}
