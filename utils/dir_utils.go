package utils

import (
	"os"
	"path/filepath"
)

func ScanDirectories(basePath string, maxDepth int) []string {
	var dirs []string

	expanded := expandHome(basePath)

	err := scanDir(expanded, expanded, maxDepth, 0, &dirs)
	if err != nil {
		return dirs
	}

	return dirs
}

func scanDir(basePath, currentPath string, maxDepth, currentDepth int, dirs *[]string) error {
	if currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if len(entry.Name()) > 0 && entry.Name()[0] == '.' {
			continue
		}

		if isExcluded(entry.Name()) {
			continue
		}

		fullPath := filepath.Join(currentPath, entry.Name())
		*dirs = append(*dirs, fullPath)

		if currentDepth < maxDepth {
			scanDir(basePath, fullPath, maxDepth, currentDepth+1, dirs)
		}
	}

	return nil
}

func expandHome(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if len(path) == 1 {
		return home
	}

	return filepath.Join(home, path[1:])
}

func isExcluded(name string) bool {
	excluded := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		"build":        true,
		"dist":         true,
		"target":       true,
		".git":         true,
		".cache":       true,
		"__pycache__":  true,
	}

	return excluded[name]
}

func GetProjectDirs(searchPaths []string, depth int) []string {
	var allDirs []string
	seen := make(map[string]bool)

	for _, path := range searchPaths {
		dirs := ScanDirectories(path, depth)
		for _, dir := range dirs {
			if !seen[dir] {
				seen[dir] = true
				allDirs = append(allDirs, dir)
			}
		}
	}

	return allDirs
}
