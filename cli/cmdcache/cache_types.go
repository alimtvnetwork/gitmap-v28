package cmdcache

// CacheFileRecord represents metadata of a cached file in the root split-db.
type CacheFileRecord struct {
	FileId       int64  `json:"fileId"`
	RepoUrl      string `json:"repoUrl"`
	RelativePath string `json:"relativePath"`
	AbsolutePath string `json:"absolutePath"`
	FileSize     int64  `json:"fileSize"`
	ModifiedTime int64  `json:"modifiedTime"`
	FolderSlug   string `json:"folderSlug"`
	IsKeep       bool   `json:"isKeep"`
}

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
