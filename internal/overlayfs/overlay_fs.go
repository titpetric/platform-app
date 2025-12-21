package overlayfs

import (
	"io/fs"
	"sort"
)

type OverlayFS struct {
	Upper fs.FS
	Lower fs.FS
}

func NewOverlayFS(upper, lower fs.FS) *OverlayFS {
	return &OverlayFS{
		Upper: upper,
		Lower: lower,
	}
}

func (o *OverlayFS) Open(name string) (fs.File, error) {
	if o.Upper != nil {
		f, err := o.Upper.Open(name)
		if err == nil {
			return f, nil
		}
	}
	if o.Lower != nil {
		return o.Lower.Open(name)
	}
	return nil, fs.ErrNotExist
}

func (o *OverlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var upperDir []fs.DirEntry
	var upperErr error
	var lowerDir []fs.DirEntry
	var lowerErr error

	if o.Upper != nil {
		upperDir, upperErr = fs.ReadDir(o.Upper, name)
	} else {
		upperErr = fs.ErrNotExist
	}

	if o.Lower != nil {
		lowerDir, lowerErr = fs.ReadDir(o.Lower, name)
	} else {
		lowerErr = fs.ErrNotExist
	}

	if upperErr != nil && lowerErr != nil {
		return nil, upperErr
	}

	merged := make(map[string]fs.DirEntry)
	for _, e := range lowerDir {
		merged[e.Name()] = e
	}
	for _, e := range upperDir {
		merged[e.Name()] = e // upper layer overrides
	}

	entries := make([]fs.DirEntry, 0, len(merged))
	for _, e := range merged {
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil
}

func (o *OverlayFS) Glob(pattern string) ([]string, error) {
	var upperMatches []string
	if o.Upper != nil {
		upperMatches, _ = fs.Glob(o.Upper, pattern)
	}

	var lowerMatches []string
	if o.Lower != nil {
		lowerMatches, _ = fs.Glob(o.Lower, pattern)
	}

	matchMap := make(map[string]struct{})
	for _, m := range lowerMatches {
		matchMap[m] = struct{}{}
	}
	for _, m := range upperMatches {
		matchMap[m] = struct{}{}
	}

	results := make([]string, 0, len(matchMap))
	for m := range matchMap {
		results = append(results, m)
	}

	sort.Strings(results)
	return results, nil
}
