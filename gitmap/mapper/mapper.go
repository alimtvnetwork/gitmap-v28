// Package mapper converts raw scan data into ScanRecord structs.
package mapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/scanner"
)

// BuildOptions and resolveDefaultBranch live in mapper_options.go to
// keep this file under the project's 200-line per-file budget.

// BuildRecords converts a list of RepoInfo into ScanRecords using the
// per-repo RelativePath the scanner already computed (against the scan
// dir). Kept as a thin wrapper so legacy callers (cmd/as.go,
// cmd/releaseautoregister.go) don't need to thread an unused root.
func BuildRecords(repos []scanner.RepoInfo, mode, defaultNote string) []model.ScanRecord {
	return BuildRecordsWithOptions(repos, BuildOptions{Mode: mode, DefaultNote: defaultNote})
}

// BuildRecordsWithRoot is like BuildRecords but rewrites every
// RelativePath against `relRoot` (must be an absolute, cleaned path)
// when non-empty. This is what powers `gitmap scan --relative-root`:
// it pins the base used for all output artifacts so running the same
// scan from different cwds yields byte-identical CSV/JSON/scripts.
//
// When a repo lives outside relRoot, filepath.Rel returns a "../"-prefixed
// path. We refuse that and emit a clear stderr message naming the offending
// repo, falling back to the scanner-computed RelativePath for THAT row so
// one bad ancestor doesn't drop the entire record.
func BuildRecordsWithRoot(repos []scanner.RepoInfo, mode, defaultNote, relRoot string) []model.ScanRecord {
	return BuildRecordsWithOptions(repos, BuildOptions{
		Mode: mode, DefaultNote: defaultNote, RelRoot: relRoot,
	})
}

// BuildRecordsWithOptions is the full-fat entry point. New options
// (e.g. DefaultBranch) ride on the BuildOptions struct so the helper
// signatures don't keep growing. The two wrappers above just forward.
func BuildRecordsWithOptions(repos []scanner.RepoInfo, opts BuildOptions) []model.ScanRecord {
	records := make([]model.ScanRecord, 0, len(repos))
	for _, repo := range repos {
		repo.RelativePath = relativePathFor(repo, opts.RelRoot)
		rec := buildOneRecord(repo, opts)
		records = append(records, rec)
	}
	// Pin (RelativePath, HTTPSUrl, SSHUrl, AbsolutePath) order so
	// terminal/CSV/JSON exports are byte-identical across runs even
	// when the scanner ordering changes upstream. See sort.go.
	SortRecords(records)

	return records
}

// relativePathFor returns the RelativePath to record for `repo`. With an
// empty relRoot we keep the scanner's value verbatim. Otherwise we recompute
// against relRoot and reject "../" escapes.
func relativePathFor(repo scanner.RepoInfo, relRoot string) string {
	if relRoot == "" {
		return repo.RelativePath
	}
	rel, err := filepath.Rel(relRoot, repo.AbsolutePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanRelativeRootNotAncestor,
			relRoot, repo.AbsolutePath, err.Error())

		return repo.RelativePath
	}
	if strings.HasPrefix(rel, "..") {
		fmt.Fprintf(os.Stderr, constants.ErrScanRelativeRootNotAncestor,
			relRoot, repo.AbsolutePath, rel)

		return repo.RelativePath
	}

	return rel
}

// buildOneRecord creates a single ScanRecord from a RepoInfo. The
// fallback branch name is taken from opts.DefaultBranch when set and
// from constants.DefaultBranch otherwise — see resolveDefaultBranch.
func buildOneRecord(repo scanner.RepoInfo, opts BuildOptions) model.ScanRecord {
	remoteURL, _ := gitutil.RemoteURL(repo.AbsolutePath)
	branch, branchSource := gitutil.DetectBranchWithDefault(
		repo.AbsolutePath, resolveDefaultBranch(opts.DefaultBranch))
	httpsURL := toHTTPS(remoteURL)
	sshURL := toSSH(remoteURL)
	transport := classifyTransport(remoteURL)
	// Per-repo transport drives clone-URL selection so an SSH-origin
	// repo never gets a silent HTTPS clone command (which triggers
	// browser auth on private remotes). The scan-wide --mode flag is
	// the tiebreaker only when transport is "other".
	cloneURL := selectCloneURLForTransport(httpsURL, sshURL, transport, opts.Mode)
	repoName := extractRepoName(remoteURL)
	noteText := buildNote(remoteURL, opts.DefaultNote)
	instruction := buildInstruction(cloneURL, branch, repo.RelativePath)
	repoID := gitutil.CanonicalRepoID(remoteURL)

	return model.ScanRecord{
		Slug:     buildSlug(httpsURL, repoName),
		RepoID:   repoID,
		RepoName: repoName, HTTPSUrl: httpsURL, SSHUrl: sshURL,
		DiscoveredURL: remoteURL,
		Branch:        branch, BranchSource: branchSource,
		RelativePath: repo.RelativePath, AbsolutePath: repo.AbsolutePath,
		CloneInstruction: instruction, Notes: noteText,
		Depth:     repo.Depth,
		Transport: transport,
	}
}

// selectCloneURLForTransport picks the URL whose transport matches the
// repo's identified origin transport. Falls through to the other URL
// if the preferred one is empty, then to the user-mode default.
func selectCloneURLForTransport(httpsURL, sshURL, transport, mode string) string {
	switch transport {
	case constants.ScanTransportSSH:
		if len(sshURL) > 0 {
			return sshURL
		}

		return httpsURL
	case constants.ScanTransportHTTPS:
		if len(httpsURL) > 0 {
			return httpsURL
		}

		return sshURL
	}

	return selectCloneURL(httpsURL, sshURL, mode)
}

// toHTTPS converts a remote URL to HTTPS format.
func toHTTPS(raw string) string {
	if strings.HasPrefix(raw, constants.PrefixHTTPS) {
		return raw
	}
	if strings.HasPrefix(raw, constants.PrefixSSH) {
		host, path := splitSSH(raw)

		return fmt.Sprintf(constants.HTTPSFromSSHFmt, host, path)
	}

	return raw
}

// toSSH converts a remote URL to SSH format.
func toSSH(raw string) string {
	if strings.HasPrefix(raw, constants.PrefixSSH) {
		return raw
	}
	if strings.HasPrefix(raw, constants.PrefixHTTPS) {
		trimmed := strings.TrimPrefix(raw, constants.PrefixHTTPS)
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) == 2 {
			return fmt.Sprintf(constants.SSHFromHTTPSFmt, parts[0], parts[1])
		}
	}

	return raw
}

// splitSSH splits a git@host:path URL into host and path.
func splitSSH(raw string) (string, string) {
	trimmed := strings.TrimPrefix(raw, constants.PrefixSSH)
	parts := strings.SplitN(trimmed, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return trimmed, ""
}

// selectCloneURL picks HTTPS or SSH URL based on mode.
func selectCloneURL(httpsURL, sshURL, mode string) string {
	if mode == constants.ModeSSH {
		return sshURL
	}

	return httpsURL
}

// extractRepoName derives the repository name from a remote URL.
func extractRepoName(raw string) string {
	if len(raw) == 0 {
		return constants.UnknownRepoName
	}
	base := filepath.Base(raw)

	return strings.TrimSuffix(base, constants.ExtGit)
}

// buildNote generates the notes field for a record.
func buildNote(remoteURL, defaultNote string) string {
	if len(remoteURL) == 0 {
		return constants.NoteNoRemote
	}

	return defaultNote
}

// buildInstruction creates the full git clone command string.
func buildInstruction(url, branch, relPath string) string {
	if len(url) == 0 {
		return ""
	}

	return fmt.Sprintf(constants.CloneInstructionFmt, branch, url, relPath)
}

// buildSlug derives a lowercase slug from the HTTPS URL.
// Falls back to repoName when the URL is empty.
func buildSlug(httpsURL, repoName string) string {
	if len(httpsURL) == 0 {
		return strings.ToLower(repoName)
	}
	base := filepath.Base(httpsURL)
	trimmed := strings.TrimSuffix(base, constants.ExtGit)

	return strings.ToLower(trimmed)
}
