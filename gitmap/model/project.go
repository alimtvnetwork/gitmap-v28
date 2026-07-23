// Package model — project.go defines the DetectedProject struct.
package model

// DetectedProject represents a project detected inside a Git repository.
type DetectedProject struct {
	ID               int64  `json:"id"`
	RepoID           int64  `json:"repoId"`
	RepoName         string `json:"repoName"`
	ProjectTypeID    int64  `json:"projectTypeId"`
	ProjectType      string `json:"projectType"`
	ProjectName      string `json:"projectName"`
	AbsolutePath     string `json:"absolutePath"`
	RepoPath         string `json:"repoPath"`
	RelativePath     string `json:"relativePath"`
	PrimaryIndicator string `json:"primaryIndicator"`
	DetectedAt       string `json:"detectedAt"`
}
