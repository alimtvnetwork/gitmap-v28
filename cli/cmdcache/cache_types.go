package cmdcache

// CacheFileRecord moved to store package
// SearchMatch represents a single search match in a cached file.
type SearchMatch struct {
	RelativePath string `json:"relativePath"`
	LineNumber   int    `json:"lineNumber"`
	Content      string `json:"content"`
}

// CacheSearchOptions holds filtering and formatting flags for search operations.
type CacheSearchOptions struct {
	Patterns    []string `json:"patterns"`
	FileGlobs   []string `json:"fileGlobs"`
	LinesToShow int      `json:"linesToShow"`
	ResultLimit int      `json:"resultLimit"`
	IsRegex     bool     `json:"isRegex"`
}
