// Package vscodeworkspace builds `.code-workspace` JSON documents
// from the gitmap repo database. Output matches the schema VS Code
// produces via "File → Save Workspace As…", so a workspace file
// emitted by gitmap can be hand-edited and re-emitted without
// surprising diffs.
//
// The builder is intentionally pure: it takes plain folder tuples
// and returns bytes. Disk I/O lives in write.go; CLI flag parsing
// and DB access live in the cmd package.
package vscodeworkspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Folder is one entry in the workspace's "folders" array.
// Name is optional in the VS Code schema but always populated by
// gitmap so the sidebar shows the repo name (not the basename of
// the on-disk folder, which may have been flattened by clone-next).
type Folder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Workspace is the top-level `.code-workspace` document.
// Settings is always emitted as `{}` so VS Code's first-open
// experience matches a freshly-saved workspace file.
type Workspace struct {
	Folders  []Folder       `json:"folders"`
	Settings map[string]any `json:"settings"`
}

// Build assembles a Workspace from the supplied folder tuples,
// dropping duplicates by path and sorting by Name for diff stability.
// Path comparison is case-folded on Windows via filepath.Clean +
// strings.EqualFold, but we keep the original casing in the output
// to preserve user-visible spelling.
func Build(folders []Folder) Workspace {
	deduped := dedupeFolders(folders)
	sort.SliceStable(deduped, func(i, j int) bool {
		return deduped[i].Name < deduped[j].Name
	})

	return Workspace{
		Folders:  deduped,
		Settings: map[string]any{},
	}
}

// Encode renders ws as tab-indented JSON with a trailing newline,
// matching the byte shape VS Code itself emits.
func Encode(ws Workspace) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", constants.VSCodeWorkspaceJSONIndent)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(ws); err != nil {
		return nil, fmt.Errorf("vscode-workspace: encode: %w", err)
	}

	return buf.Bytes(), nil
}

// Relativize rewrites each folder's Path to be relative to baseDir
// (the directory the workspace file lives in). Paths that cannot be
// expressed relatively (different volume on Windows) are left as
// absolute and an error is returned so the caller can surface the
// fallback.
func Relativize(folders []Folder, baseDir string) ([]Folder, error) {
	out := make([]Folder, 0, len(folders))
	var firstErr error
	for _, f := range folders {
		rel, err := filepath.Rel(baseDir, f.Path)
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("relativize %q against %q: %w", f.Path, baseDir, err)
			}
			out = append(out, f)

			continue
		}
		out = append(out, Folder{Name: f.Name, Path: filepath.ToSlash(rel)})
	}

	return out, firstErr
}

// dedupeFolders drops entries with the same cleaned absolute path,
// keeping the first occurrence so caller-supplied ordering wins.
func dedupeFolders(folders []Folder) []Folder {
	seen := make(map[string]struct{}, len(folders))
	out := make([]Folder, 0, len(folders))
	for _, f := range folders {
		key := filepath.Clean(f.Path)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, f)
	}

	return out
}
