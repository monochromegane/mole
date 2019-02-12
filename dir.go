package mole

import (
	"go/build"
	"path/filepath"
)

func modDirs() []string {
	return dirs(filepath.Join("pkg", "mod"))
}

func dirs(path string) []string {
	gopaths := filepath.SplitList(build.Default.GOPATH)

	dirs := make([]string, len(gopaths))
	for i, gopath := range gopaths {
		dirs[i] = filepath.Join(gopath, path)
	}
	return dirs
}
