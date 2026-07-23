package vscodepm

// Entry mirrors one object in projects.json. Field order and JSON tags
// match the alefragnani.project-manager schema exactly so encoded files
// stay diff-friendly with manual edits and the sample fixture
// (spec/01-vscode-project-manager-sync/sample-projects.json).
type Entry struct {
	Name     string   `json:"name"`
	RootPath string   `json:"rootPath"`
	Paths    []string `json:"paths"`
	Tags     []string `json:"tags"`
	Enabled  bool     `json:"enabled"`
	Profile  string   `json:"profile"`
}

// SyncSummary is returned from Sync to describe what changed.
type SyncSummary struct {
	Added     int
	Updated   int
	Unchanged int
	Total     int
}

// newEntry builds a default Entry for a freshly registered tuple.
// Tags and Paths are always emitted as non-nil empty slices so the encoded
// JSON contains `[]` rather than `null` (matches the sample fixture).
// Auto-detected tags arrive via the `tags` parameter (v3.40.0+).
func newEntry(rootPath, name string, paths, tags []string) Entry {
	if paths == nil {
		paths = []string{}
	}
	if tags == nil {
		tags = []string{}
	}

	return Entry{
		Name:     name,
		RootPath: rootPath,
		Paths:    paths,
		Tags:     tags,
		Enabled:  true,
		Profile:  "",
	}
}

// ensureSlices replaces nil slices with empty ones so re-encoded JSON
// preserves `[]` even on entries that arrived with `null` from disk.
func ensureSlices(e Entry) Entry {
	if e.Paths == nil {
		e.Paths = []string{}
	}

	if e.Tags == nil {
		e.Tags = []string{}
	}

	return e
}

// unionPaths returns the union of `existing` and `incoming`, preserving
// the order of `existing` first then appending any new entries from
// `incoming`. Path comparison is OS-aware (case-insensitive on Windows).
func unionPaths(existing, incoming []string) []string {
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	out := make([]string, 0, len(existing)+len(incoming))

	for _, p := range existing {
		key := normalizePath(p)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}

	for _, p := range incoming {
		key := normalizePath(p)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}

	return out
}
