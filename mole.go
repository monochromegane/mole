package mole

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

var regPath = regexp.MustCompile("@v[0-9]+\\.[0-9]+\\.[0-9].*?(/|$)")
var regVersion = regexp.MustCompile("^v[0-9]+")

func Run(all bool) ([]string, error) {
	modLibs := map[string]map[string]string{}
	for _, modDir := range modDirs() {
		err := filepath.Walk(modDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}
			if isCacheDir(modDir, path) {
				return filepath.SkipDir
			}
			if isSkipDir(info) {
				return filepath.SkipDir
			}

			importPath, version, err := extractImportPathAndVersion(modDir, path, info)
			if err != nil {
				return nil
			}

			_, ok := modLibs[importPath]
			if !ok {
				versions := make(map[string]string)
				modLibs[importPath] = versions
			}
			modLibs[importPath][version] = path
			return nil
		})
		if err != nil {
			return []string{}, err
		}
	}

	results := []string{}
	for _, versions := range modLibs {
		var semvers semver.Versions
		for version, _ := range versions {
			sv, err := semver.ParseTolerant(version)
			if err != nil {
				panic(err)
			}
			semvers = append(semvers, sv)
		}
		semver.Sort(semvers)

		if all {
			for _, sv := range semvers {
				results = append(results, versions["v"+sv.String()])
			}
		} else {
			latest := semvers[len(semvers)-1]
			results = append(results, versions["v"+latest.String()])
		}
	}
	return results, nil
}

func extractImportPathAndVersion(base, path string, info os.FileInfo) (string, string, error) {
	version := regPath.FindString(info.Name())
	if version == "" {
		return "", "", fmt.Errorf("No version directory")
	}
	pkg := regPath.ReplaceAllString(info.Name(), "")
	if regVersion.MatchString(pkg) {
		rel, _ := filepath.Rel(base, filepath.Dir(path))
		return rel, strings.Trim(version, "@"), nil
	}
	rel, _ := filepath.Rel(base, strings.TrimSuffix(path, version))
	return rel, strings.Trim(version, "@"), nil
}

func isCacheDir(base, path string) bool {
	rel, _ := filepath.Rel(base, path)
	return rel == "cache"
}

func isSkipDir(fi os.FileInfo) bool {
	name := fi.Name()
	switch name {
	case "", "internal", "testdata", "vendor":
		return true
	}
	switch name[0] {
	case '.', '_':
		return true
	}
	return false
}
