package globchanel

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
func Glob(pattern string) chan string {
	matches := make(chan string)
	go func() {
		glob(pattern, matches)
		close(matches)
	}()
	return matches
}

func glob(pattern string, matches chan string) {
	// When the pattern doesn't contain meta
	if !hasMeta(pattern) {
		if _, err := os.Lstat(pattern); err == nil {
			matches <- pattern
		}
		return
	}
	dir, file := filepath.Split(pattern)

	if onWindows {
		dir = cleanGlobPathWindows(dir)
	} else {
		dir = cleanGlobPath(dir)
	}
	if !hasMeta(dir) {
		// solid dir
		globpath(dir, file, matches)
		return
	}

	directories := Glob(dir)
	for d := range directories {
		globpath(d, file, matches)
	}
}

// cleanGlobPath prepares path for glob matching.
func cleanGlobPath(path string) string {
	switch path {
	case "":
		return "."
	case string(filepath.Separator):
		// do nothing to the path
		return path
	default:
		return path[0 : len(path)-1] // chop off trailing separator
	}
}

// cleanGlobPathWindows is windows version of cleanGlobPath.
func cleanGlobPathWindows(path string) string {
	vollen := volumeNameLen(path)
	switch {
	case path == "":
		return "."
	case vollen+1 == len(path) && os.IsPathSeparator(path[len(path)-1]): // /, \, C:\ and C:/
		// do nothing to the path
		return path
	case vollen == len(path) && len(path) == 2: // C:
		return path + "." // convert C: into C:.
	default:
		return path[0 : len(path)-1] // chop off trailing separator
	}
}

// glob searches for files matching pattern in the directory dir
// and appends them to matches. If the directory cannot be
// opened, it returns the existing matches. New matches are
// added in lexicographical order.
func globpath(dir, pattern string, matches chan string) {
	fi, err := os.Stat(dir)
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return
	}
	d, err := os.Open(dir)
	if err != nil {
		return
	}
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	sort.Strings(names)

	for _, n := range names {
		if matched, _ := filepath.Match(pattern, n); matched {
			matches <- filepath.Join(dir, n)
		}
	}
	return
}

// hasMeta reports whether path contains any of the magic characters
// recognized by Match.
func hasMeta(path string) bool {
	// TODO(niemeyer): Should other magic characters be added here?
	return strings.ContainsAny(path, "*?[")
}

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

var onWindows = runtime.GOOS == "windows"

// volumeNameLen returns length of the leading volume name on Windows.
// It returns 0 elsewhere. Copyer from filepath.path_windows.go
func volumeNameLen(path string) int {
	if !onWindows {
		return 0
	}
	if len(path) < 2 {
		return 0
	}
	// with drive letter
	c := path[0]
	if path[1] == ':' && ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') {
		return 2
	}
	// is it UNC
	if l := len(path); l >= 5 && isSlash(path[0]) && isSlash(path[1]) &&
		!isSlash(path[2]) && path[2] != '.' {
		// first, leading `\\` and next shouldn't be `\`. its server name.
		for n := 3; n < l-1; n++ {
			// second, next '\' shouldn't be repeated.
			if isSlash(path[n]) {
				n++
				// third, following something characters. its share name.
				if !isSlash(path[n]) {
					if path[n] == '.' {
						break
					}
					for ; n < l; n++ {
						if isSlash(path[n]) {
							break
						}
					}
					return n
				}
				break
			}
		}
	}
	return 0
}
