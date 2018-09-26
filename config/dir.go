package config

import (
	"runtime"
	"path/filepath"
)

func Path(n int) string {
	_, filename, _, _ := runtime.Caller(1)
	if n == 0 {
		return filepath.Dir(filename)
	}

	n--
	return filepath.Dir(Path(n))
}