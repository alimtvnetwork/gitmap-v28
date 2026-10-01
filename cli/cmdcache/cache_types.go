package cmdcache

import "github.com/alimtvnetwork/gitmap-v28/cli/store"

// SearchMatch represents a single search match in a cached file.
type SearchMatch struct {
	RelativePath string                   `json:"relativePath"`
	LineNumber   int                      `json:"lineNumber"`
	Content      string                   `json:"content"`
	ContextLines []store.CacheContextLine `json:"contextLines,omitempty"`
}

// CacheSearchOptions holds filtering and formatting flags for search operations.
type CacheSearchOptions struct {
	Patterns    []string `json:"patterns"`
	FileGlobs   []string `json:"fileGlobs"`
	LinesToShow int      `json:"linesToShow"`
	ResultLimit int      `json:"resultLimit"`
	IsRegex     bool     `json:"isRegex"`
}

// CacheCreateOptions holds options for cache creation.
type CacheCreateOptions struct {
	Targets []string `json:"targets"`
	IsKeep  bool     `json:"isKeep"`
}
