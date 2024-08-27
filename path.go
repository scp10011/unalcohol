package unalcohol

import (
	"path"
	"slices"
)

func JoinPath(ptr API, url string) string {
	paths := []string{url}
	for {
		paths = append(paths, ptr.GetPath())
		if p := ptr.GetPtr(); p != nil {
			ptr = p
		} else {
			break
		}
	}
	slices.Reverse(paths)
	return path.Join(paths...)
}
