package utils

import (
	"os"
	"path/filepath"
)

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			ScanDirectories recursively scans a directory up to a specified depth.
//
//		@Description	Expands ~ to home directory and scans for subdirectories
//
//		@Param			basePath	string	Root directory to scan
//		@Param			maxDepth	int		Maximum recursion depth
//
//		@Return			[]string	List of discovered directories
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func ScanDirectories(basePath string, maxDepth int) []string {
	var dirs []string

	expanded := expandHome(basePath)

	err := scanDir(expanded, expanded, maxDepth, 0, &dirs)
	if err != nil {
		return dirs
	}

	return dirs
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			scanDir is the recursive helper for directory scanning.
//
//		@Param			basePath		string		Original base path
//		@Param			currentPath		string		Current directory being scanned
//		@Param			maxDepth		int			Maximum recursion depth
//		@Param			currentDepth	int			Current recursion depth
//		@Param			dirs			*[]string	Accumulator for discovered directories
//
//		@Return			error			Error if directory cannot be read
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			expandHome expands ~ to the user's home directory.
//
//		@Param			path	string	Path that may contain ~
//
//		@Return			string	Expanded path
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			isExcluded checks if a directory name should be excluded from scanning.
//
//		@Description	Excludes common build artifacts and dependency directories
//
//		@Param			name	string	Directory name
//
//		@Return			bool	True if directory should be excluded
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			GetProjectDirs scans multiple search paths and returns unique directories.
//
//		@Description	Deduplicates directories found across multiple search paths
//
//		@Param			searchPaths	[]string	List of paths to scan
//		@Param			depth		int			Maximum scanning depth
//
//		@Return			[]string	Deduplicated list of discovered directories
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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
