// Package model — envregistry.go defines the env variable registry record.
package model

// EnvVariable represents a managed environment variable.
type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EnvPathEntry represents a managed PATH directory entry.
type EnvPathEntry struct {
	Path string `json:"path"`
}

// EnvRegistry represents the top-level env-registry.json structure.
type EnvRegistry struct {
	Variables []EnvVariable  `json:"variables"`
	Paths     []EnvPathEntry `json:"paths"`
}
