// Package model — export.go defines the full database export structure.
package model

// DatabaseExport holds the complete database state for portable backup.
type DatabaseExport struct {
	Version    string                 `json:"version"`
	ExportedAt string                 `json:"exportedAt"`
	Repos      []ScanRecord           `json:"repos"`
	Groups     []GroupExport          `json:"groups"`
	Releases   []ReleaseRecord        `json:"releases"`
	History    []CommandHistoryRecord `json:"history"`
	Bookmarks  []BookmarkRecord       `json:"bookmarks"`
}

// GroupExport extends Group with its member repo slugs.
type GroupExport struct {
	Group
	RepoSlugs []string `json:"repoSlugs"`
}
