package cmdscan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// NormalizeRepoURL trims whitespace, lowercases, strips trailing .git and /,
// and normalizes SSH and HTTPS URLs to a canonical host/path representation.
func NormalizeRepoURL(raw string) string {
	u := strings.TrimSpace(raw)
	if u == "" {
		return ""
	}
	u = strings.ToLower(u)

	// Strip protocol schemes
	u = strings.TrimPrefix(u, "ssh://")
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "http://")
	u = strings.TrimPrefix(u, "git://")

	// Strip user prefix if present (e.g., git@)
	if atIdx := strings.Index(u, "@"); atIdx != -1 {
		u = u[atIdx+1:]
	}

	// Normalize scp-style SSH syntax: host:owner/repo -> host/owner/repo
	slashIdx := strings.Index(u, "/")
	colonIdx := strings.Index(u, ":")
	if colonIdx != -1 && (slashIdx == -1 || colonIdx < slashIdx) {
		u = u[:colonIdx] + "/" + u[colonIdx+1:]
	}

	// Strip trailing slashes and .git extensions
	for {
		if strings.HasSuffix(u, "/") {
			u = strings.TrimSuffix(u, "/")
			continue
		}
		if strings.HasSuffix(u, ".git") {
			u = strings.TrimSuffix(u, ".git")
			continue
		}
		break
	}

	return u
}

// MergeScanRecordsWithExistingManifest merges scanned repositories into
// repo-cache/01-gitmap/gitmap.json with URL deduplication and metadata refreshing.
func MergeScanRecordsWithExistingManifest(root string, records []model.ScanRecord, isQuiet bool) error {
	target := filepath.Join(root, "01-gitmap", "gitmap.json")
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir repo-cache 01-gitmap directory")
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return writeInitialManifest(target, records, isQuiet)
	}

	return mergeExistingManifest(target, records, isQuiet)
}

func writeInitialManifest(target string, records []model.ScanRecord, isQuiet bool) error {
	exports := make([]exportRecord, 0, len(records))
	for _, r := range records {
		exports = append(exports, toExportRecord(r))
	}

	if err := writeExportRecordsJSON(target, exports); err != nil {
		return apperror.WrapSimple(err, "write initial repo-cache manifest")
	}

	if !isQuiet {
		fmt.Printf("  ✓ Created repo-cache manifest (%d repos): %s\n", len(exports), target)
	}

	return nil
}

func mergeExistingManifest(target string, records []model.ScanRecord, isQuiet bool) error {
	existing, _, err := readExportRecords(target)
	if err != nil {
		return apperror.WrapSimple(err, "read existing repo-cache manifest")
	}

	existingByURL := make(map[string]int, len(existing))
	for idx, exp := range existing {
		key := resolveRecordCanonicalKey(exp)
		if key != "" {
			existingByURL[key] = idx
		}
	}

	var updatedCount int
	var newCount int
	for _, rec := range records {
		incomingExp := toExportRecord(rec)
		key := resolveRecordCanonicalKey(incomingExp)

		if key != "" {
			if matchIdx, hasExisting := existingByURL[key]; hasExisting {
				existing[matchIdx].Branch = incomingExp.Branch
				existing[matchIdx].BranchSource = incomingExp.BranchSource
				existing[matchIdx].RelativePath = incomingExp.RelativePath
				existing[matchIdx].AbsolutePath = incomingExp.AbsolutePath
				existing[matchIdx].CloneInstruction = incomingExp.CloneInstruction
				if incomingExp.URL != "" {
					existing[matchIdx].URL = incomingExp.URL
				}
				updatedCount++
				continue
			}
		}

		existing = append(existing, incomingExp)
		if key != "" {
			existingByURL[key] = len(existing) - 1
		}
		newCount++
	}

	if err := writeExportRecordsJSON(target, existing); err != nil {
		return apperror.WrapSimple(err, "write merged repo-cache manifest")
	}

	if !isQuiet {
		fmt.Printf("  ✓ Repo-cache manifest merged (%d new, %d updated, %d total): %s\n",
			newCount, updatedCount, len(existing), target)
	}

	return nil
}

func resolveRecordCanonicalKey(rec exportRecord) string {
	if rec.URL != "" {
		if norm := NormalizeRepoURL(rec.URL); norm != "" {
			return norm
		}
	}
	if rec.HTTPSUrl != "" {
		if norm := NormalizeRepoURL(rec.HTTPSUrl); norm != "" {
			return norm
		}
	}
	if rec.SSHUrl != "" {
		if norm := NormalizeRepoURL(rec.SSHUrl); norm != "" {
			return norm
		}
	}
	if rec.DiscoveredURL != "" {
		if norm := NormalizeRepoURL(rec.DiscoveredURL); norm != "" {
			return norm
		}
	}
	return ""
}
