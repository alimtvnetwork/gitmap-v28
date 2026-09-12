// Package model defines the core data structures for gitmap.
package model

// ZipGroup represents a named collection of files/folders for archiving.
type ZipGroup struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ArchiveName string `json:"archiveName"`
	CreatedAt   string `json:"createdAt"`
}

// ZipGroupItem links a file or folder path to a zip group.
type ZipGroupItem struct {
	GroupID      int64  `json:"groupId"`
	RepoPath     string `json:"repoPath"`
	RelativePath string `json:"relativePath"`
	FullPath     string `json:"fullPath"`
	IsFolder     bool   `json:"isFolder"`
	// Path returns FullPath for backward compatibility with zip archive logic.
	Path string `json:"-"`
}

// ResolvePath sets the Path field from FullPath for archive operations.
func (z *ZipGroupItem) ResolvePath() {
	z.Path = z.FullPath
}
