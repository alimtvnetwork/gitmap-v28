// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import "time"

// SourceRepoConfig models an incoming repository to be merged.
type SourceRepoConfig struct {
	RepoURL        string    `json:"repoUrl"`
	Branch         string    `json:"branch,omitempty"`
	CommitRange    string    `json:"commitRange,omitempty"`
	LastCommitDate time.Time `json:"lastCommitDate"`
	Description    string    `json:"description,omitempty"`
	ReleaseTag     string    `json:"releaseTag,omitempty"`
	LocalStaging   string    `json:"localStaging,omitempty"`
	FolderTree     []string  `json:"folderTree,omitempty"`
}

// FileCollisionVariant models one sequenced instance of a conflicting file.
type FileCollisionVariant struct {
	Sequence   string `json:"sequence"`
	File       string `json:"file"`
	SourceRepo string `json:"sourceRepo"`
	Commit     string `json:"commit"`
}

// FileCollisionRecord models a conflicting path requiring AI consolidation.
type FileCollisionRecord struct {
	CanonicalPath string                 `json:"canonicalPath"`
	Variants      []FileCollisionVariant `json:"variants"`
	Status        string                 `json:"status"` // "pending_ai_consolidation", "consolidated"
}

// MergeAIManifest models the root merge-ai-manifest.json artifact.
type MergeAIManifest struct {
	Attributes struct {
		GeneratedAt        string `json:"generatedAt"`
		GitMapVersion      string `json:"gitMapVersion"`
		Tool               string `json:"tool"`
		CollisionStrategy  string `json:"collisionStrategy"`
	} `json:"attributes"`
	Data struct {
		Destination struct {
			TargetType   string `json:"targetType"`
			Slug         string `json:"slug"`
			ResolvedPath string `json:"resolvedPath"`
			OriginURL    string `json:"originUrl,omitempty"`
		} `json:"destination"`
		RepoSequence               []SourceRepoSequenceEntry `json:"repoSequence"`
		FileCollisions             []FileCollisionRecord     `json:"fileCollisions"`
		UniqueFilesDirectlyPlaced  int                       `json:"uniqueFilesDirectlyPlaced"`
	} `json:"data"`
}

// SourceRepoSequenceEntry models one entry in data.repoSequence of the manifest.
type SourceRepoSequenceEntry struct {
	SequenceOrder  int      `json:"sequenceOrder"`
	RepoURL        string   `json:"repoUrl"`
	Branch         string   `json:"branch"`
	CommitRange    string   `json:"commitRange"`
	LastCommitDate string   `json:"lastCommitDate"`
	Description    string   `json:"description"`
	ReleaseInfo    struct {
		LastTag string `json:"lastTag,omitempty"`
		TagHash string `json:"tagHash,omitempty"`
	} `json:"releaseInfo"`
	FolderTree []string `json:"folderTree"`
}

// MergeAIConfigFile models optional input configuration JSON (gitmap ma config.json).
type MergeAIConfigFile struct {
	Destination string   `json:"destination"`
	Sources     []string `json:"sources"`
}
