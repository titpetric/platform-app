package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// check if folder contains at least one *.go file
func hasGoFiles(folder string) bool {
	found := false
	filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".go") {
			found = true
			return filepath.SkipDir // stop walking further
		}
		return nil
	})
	return found
}

// cellValue returns the display value for a given app and column.
// For wildcard columns (e.g. "service/*"), it returns a CSV of matching subfolder names.
// For regular columns, it returns "[x]" or "[ ]".
func cellValue(app, col string) string {
	if col == catchAllCol {
		return catchAllValue(app)
	}
	if strings.HasSuffix(col, "/*") {
		prefix := strings.TrimSuffix(col, "/*")
		dir := filepath.Join(app, prefix)
		subs, err := os.ReadDir(dir)
		if err != nil {
			return ""
		}
		var names []string
		for _, sub := range subs {
			if !sub.IsDir() {
				continue
			}
			names = append(names, sub.Name())
		}
		return strings.Join(names, ", ")
	}
	full := filepath.Join(app, col)
	_, err := os.Stat(full)
	if err != nil {
		return "[ ]"
	}
	// file or directory exists
	return "[x]"
}

// compute max width for padding
func columnWidths(apps []string, columns []string) map[string]int {
	widths := make(map[string]int)
	widths["folder"] = len("folder")
	for _, col := range columns {
		widths[col] = len(col)
	}
	for _, app := range apps {
		if len(app) > widths["folder"] {
			widths["folder"] = len(app)
		}
		for _, col := range columns {
			val := cellValue(app, col)
			if len(val) > widths[col] {
				widths[col] = len(val)
			}
		}
	}
	return widths
}

// pad string to width
func pad(s string, width int) string {
	format := fmt.Sprintf(" %%-%ds ", width)
	return fmt.Sprintf(format, s)
}

// catchAllValue returns a CSV of subdirectories with Go files not covered by expected patterns.
func catchAllValue(app string) string {
	// build set of covered top-level subfolder names
	covered := make(map[string]bool)
	for _, col := range expected {
		name := col
		if strings.HasSuffix(name, "/*") {
			name = strings.TrimSuffix(name, "/*")
		}
		// only the top-level part matters
		if i := strings.Index(name, "/"); i >= 0 {
			name = name[:i]
		}
		covered[name] = true
	}

	ignore := map[string]bool{"bin": true}

	subs, err := os.ReadDir(app)
	if err != nil {
		return ""
	}
	var names []string
	for _, sub := range subs {
		if !sub.IsDir() || covered[sub.Name()] || ignore[sub.Name()] {
			continue
		}
		names = append(names, sub.Name())
	}
	return strings.Join(names, ", ")
}

const catchAllCol = "*"

func main() {
	skip := map[string]bool{"cmd": true, "autoload": true}

	entries, _ := os.ReadDir(".")
	var apps []string
	for _, e := range entries {
		if e.IsDir() && e.Name()[0] != '.' && !skip[e.Name()] {
			apps = append(apps, e.Name())
		}
	}

	// filter columns: only keep those with at least one non-empty cell
	var filtered []string
	for _, col := range expected {
		for _, app := range apps {
			if cellValue(app, col) != "" && cellValue(app, col) != "[ ]" {
				filtered = append(filtered, col)
				break
			}
		}
	}

	// add catch-all column if any app has uncovered subdirs
	hasCatchAll := false
	for _, app := range apps {
		if catchAllValue(app) != "" {
			hasCatchAll = true
			break
		}
	}
	if hasCatchAll {
		filtered = append(filtered, catchAllCol)
	}

	// filter apps: only keep those that have Go files
	var filteredApps []string
	for _, app := range apps {
		if hasGoFiles(app) {
			filteredApps = append(filteredApps, app)
		}
	}
	apps = filteredApps

	widths := columnWidths(apps, filtered)

	// header
	fmt.Print("|")
	fmt.Print(pad("folder", widths["folder"]))
	for _, col := range filtered {
		fmt.Print("|")
		fmt.Print(pad(col, widths[col]))
	}
	fmt.Println("|")

	// separator
	fmt.Print("|")
	fmt.Print(pad(strings.Repeat("-", widths["folder"]), widths["folder"]))
	for _, col := range filtered {
		fmt.Print("|")
		fmt.Print(pad(strings.Repeat("-", widths[col]), widths[col]))
	}
	fmt.Println("|")

	// rows
	for _, app := range apps {
		fmt.Print("|")
		fmt.Print(pad(app, widths["folder"]))
		for _, col := range filtered {
			fmt.Print("|")
			fmt.Print(pad(cellValue(app, col), widths[col]))
		}
		fmt.Println("|")
	}
}
