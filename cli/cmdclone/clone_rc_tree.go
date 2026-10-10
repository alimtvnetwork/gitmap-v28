package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// DiscoveredManifest holds metadata and entries for a discovered repo-cache manifest file.
type DiscoveredManifest struct {
	Index        int              `json:"index"`
	RelativePath string           `json:"relativePath"`
	FullPath     string           `json:"fullPath"`
	Entries      []RepoCacheEntry `json:"entries"`
	TotalRepos   int              `json:"totalRepos"`
	SSHCount     int              `json:"sshCount"`
	HTTPSCount   int              `json:"httpsCount"`
}

// FormatOwnerRepo extracts a concise "owner/repo" pair from any raw Git URL or path.
func FormatOwnerRepo(rawURL string) (string, string) {
	trimmed := strings.TrimSpace(rawURL)
	trimmed = strings.TrimRight(trimmed, "/\\")
	trimmed = strings.TrimSuffix(trimmed, ".git")
	trimmed = strings.TrimRight(trimmed, "/\\")

	if len(trimmed) == 0 {
		return "unknown", "unknown"
	}

	normalized := strings.ReplaceAll(trimmed, "\\", "/")

	// SSH shorthand: git@host:owner/repo
	if strings.HasPrefix(strings.ToLower(normalized), "git@") {
		atIdx := strings.Index(normalized, "@")
		colonIdx := strings.Index(normalized[atIdx:], ":")
		if colonIdx > 0 {
			pathPart := normalized[atIdx+colonIdx+1:]
			pathPart = strings.Trim(pathPart, "/")
			parts := strings.Split(pathPart, "/")
			if len(parts) >= 2 {
				return parts[len(parts)-2], parts[len(parts)-1]
			} else if len(parts) == 1 && len(parts[0]) > 0 {
				return "unknown", parts[0]
			}
		}
	}

	// URL with scheme: https://, http://, ssh://, etc.
	if schemeIdx := strings.Index(normalized, "://"); schemeIdx > 0 {
		withoutScheme := normalized[schemeIdx+3:]
		firstSlash := strings.Index(withoutScheme, "/")
		if firstSlash >= 0 {
			pathPart := strings.Trim(withoutScheme[firstSlash+1:], "/")
			parts := strings.Split(pathPart, "/")
			if len(parts) >= 2 {
				return parts[len(parts)-2], parts[len(parts)-1]
			} else if len(parts) == 1 && len(parts[0]) > 0 {
				return "unknown", parts[0]
			}
		}
	}

	// Local filesystem path or shorthand (e.g. C:/repos/local-repo, ./local-repo, local-repo)
	cleanPath := strings.TrimPrefix(normalized, "./")
	parts := strings.Split(cleanPath, "/")
	if len(parts) >= 2 {
		last := parts[len(parts)-1]
		secondLast := parts[len(parts)-2]
		if strings.HasSuffix(secondLast, ":") || secondLast == "repos" || strings.HasPrefix(secondLast, ".") {
			return "local", last
		}
		return secondLast, last
	} else if len(parts) == 1 {
		return "local", parts[0]
	}

	return "local", "local-repo"
}

// ClassifyProtocol inspects a URL and returns "[SSH]", "[Public HTTPS]", or "[Local Path]".
func ClassifyProtocol(rawURL string) string {
	low := strings.ToLower(strings.TrimSpace(rawURL))
	if strings.HasPrefix(low, "git@") || strings.HasPrefix(low, "ssh://") {
		return "[SSH]"
	}
	if strings.HasPrefix(low, "https://") || strings.HasPrefix(low, "http://") {
		return "[Public HTTPS]"
	}
	return "[Local Path]"
}

// DiscoverRCManifests scans repoCacheRoot in priority sequence:
// 1. Primary manifest: 01-gitmap/gitmap.json or 01-gitmap.json
// 2. Sequential manifests: 02-gitmap.json, 03-gitmap.json, etc. (and 02-gitmap/gitmap.json)
// 3. Auxiliary JSON files (*.json)
func DiscoverRCManifests(repoCacheRoot string) ([]DiscoveredManifest, error) {
	if len(repoCacheRoot) == 0 {
		repoCacheRoot = ResolveRepoCacheRoot()
	}

	db, err := OpenRepoCacheDB()
	if err != nil {
		return nil, fmt.Errorf("failed to open repo cache db: %w", err)
	}
	defer db.Close()

	var discoveredPaths []string
	seenPaths := make(map[string]bool)

	addPath := func(fullPath string) {
		clean := filepath.Clean(fullPath)
		if seenPaths[clean] {
			return
		}
		if info, statErr := os.Stat(clean); statErr == nil && !info.IsDir() {
			seenPaths[clean] = true
			discoveredPaths = append(discoveredPaths, clean)
		}
	}

	// Priority 1: Primary Target 01-gitmap/gitmap.json or 01-gitmap.json
	primaryFolder := filepath.Join(repoCacheRoot, "01-gitmap", "gitmap.json")
	primaryFile := filepath.Join(repoCacheRoot, "01-gitmap.json")
	if _, err := os.Stat(primaryFolder); err == nil {
		addPath(primaryFolder)
	} else if _, err := os.Stat(primaryFile); err == nil {
		addPath(primaryFile)
	}

	// Priority 2: Sequential Manifests (02-gitmap.json, 03-gitmap.json, ...)
	type seqEntry struct {
		seq  int
		path string
	}
	var seqList []seqEntry

	entries, readErr := os.ReadDir(repoCacheRoot)
	if readErr == nil {
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				// e.g. 02-gitmap/gitmap.json
				if num, hasSeq := extractTwoDigitSeqPrefix(name); hasSeq && num > 1 {
					subPath := filepath.Join(repoCacheRoot, name, "gitmap.json")
					if _, err := os.Stat(subPath); err == nil {
						seqList = append(seqList, seqEntry{seq: num, path: subPath})
					}
				}
			} else {
				// e.g. 02-gitmap.json
				if strings.HasSuffix(strings.ToLower(name), ".json") {
					if num, hasSeq := extractTwoDigitSeqPrefix(name); hasSeq && num > 1 {
						seqList = append(seqList, seqEntry{seq: num, path: filepath.Join(repoCacheRoot, name)})
					}
				}
			}
		}
	}

	sort.Slice(seqList, func(i, j int) bool {
		return seqList[i].seq < seqList[j].seq
	})
	for _, item := range seqList {
		addPath(item.path)
	}

	// Priority 3: Auxiliary JSON Manifests
	if readErr == nil {
		var auxFiles []string
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
				full := filepath.Join(repoCacheRoot, entry.Name())
				if !seenPaths[filepath.Clean(full)] {
					auxFiles = append(auxFiles, full)
				}
			}
		}
		sort.Strings(auxFiles)
		for _, aux := range auxFiles {
			addPath(aux)
		}
	}

	// Load entries and build DiscoveredManifest structs
	var results []DiscoveredManifest
	for _, path := range discoveredPaths {
		entries, syncErr := GetOrSyncManifestEntries(db, path)
		if syncErr != nil || len(entries) == 0 {
			continue
		}

		sshCount := 0
		httpsCount := 0
		for _, e := range entries {
			if e.IsSSH {
				sshCount++
			}
			if e.IsHTTPS {
				httpsCount++
			}
		}

		rel, relErr := filepath.Rel(repoCacheRoot, path)
		if relErr != nil {
			rel = filepath.Base(path)
		}
		rel = filepath.ToSlash(rel)

		results = append(results, DiscoveredManifest{
			Index:        len(results) + 1,
			RelativePath: rel,
			FullPath:     path,
			Entries:      entries,
			TotalRepos:   len(entries),
			SSHCount:     sshCount,
			HTTPSCount:   httpsCount,
		})
	}

	return results, nil
}

func extractTwoDigitSeqPrefix(name string) (int, bool) {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 || len(parts[0]) < 2 {
		return 0, false
	}
	num, err := strconv.Atoi(parts[0])
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}

// RenderShortTreeView formats discovered manifests into the concise TreeView terminal layout.
func RenderShortTreeView(manifests []DiscoveredManifest) string {
	if len(manifests) == 0 {
		return "No manifests found in repo-cache.\n"
	}

	var b strings.Builder
	b.WriteString("Available Repository Manifests in repo-cache:\n\n")

	for i, m := range manifests {
		repoWord := "repositories"
		if m.TotalRepos == 1 {
			repoWord = "repository"
		}
		fmt.Fprintf(&b, "[%d] %s (%d %s)\n", m.Index, m.RelativePath, m.TotalRepos, repoWord)

		previewLimit := 3
		if len(m.Entries) <= previewLimit {
			for idx, e := range m.Entries {
				connector := "├── "
				if idx == len(m.Entries)-1 {
					connector = "└── "
				}
				proto := ClassifyProtocol(e.CloneUrl)
				fmt.Fprintf(&b, "    %s%s/%s %s\n", connector, e.Owner, e.RepoName, proto)
			}
		} else {
			for idx := 0; idx < previewLimit; idx++ {
				e := m.Entries[idx]
				proto := ClassifyProtocol(e.CloneUrl)
				fmt.Fprintf(&b, "    ├── %s/%s %s\n", e.Owner, e.RepoName, proto)
			}
			remaining := len(m.Entries) - previewLimit
			remWord := "repositories"
			if remaining == 1 {
				remWord = "repository"
			}
			fmt.Fprintf(&b, "    └── ... (%d more %s)\n", remaining, remWord)
		}

		if i < len(manifests)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
