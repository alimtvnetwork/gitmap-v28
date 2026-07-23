package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// SourceHandle is the resolved <source> from spec §2.3. Path is the
// absolute on-disk location after the four-case resolution. Kind
// identifies which spec case fired so the caller can log the right
// banner. IsFreshlyInit reports whether `git init` was run as part of
// resolution (cases 3 and 4).
type SourceHandle struct {
	Path          string
	Kind          SourceKind
	IsFreshlyInit bool
}

// SourceKind enumerates spec §2.3 cases. PascalCase per Core memory
// rules; literal strings live in this file (single use site).
type SourceKind uint8

const (
	// SourceKindCloned is spec §2.3 case 1: URL → git clone.
	SourceKindCloned SourceKind = iota + 1
	// SourceKindExistingRepo is spec §2.3 case 2: dir w/ .git → reuse.
	SourceKindExistingRepo
	// SourceKindInitInPlace is spec §2.3 case 3: dir w/o .git → git init.
	SourceKindInitInPlace
	// SourceKindCreatedAndInit is spec §2.3 case 4: missing → mkdir + init.
	SourceKindCreatedAndInit
)

// EnsureSource implements spec §3.1 stage 06. Pure routing — heavy
// work is delegated to one helper per spec case so each function stays
// well under the 15-line cap.
func EnsureSource(rawSource string) (*SourceHandle, error) {
	if isGitURL(rawSource) {
		return resolveByClone(rawSource)
	}
	abs, err := filepath.Abs(rawSource)
	if err != nil {
		return nil, fmt.Errorf("absolutize source: %w", err)
	}
	info, statErr := os.Stat(abs)
	if statErr == nil && info.IsDir() {
		return resolveExistingDir(abs)
	}
	if statErr != nil && os.IsNotExist(statErr) {
		return resolveMissingDir(abs)
	}
	if statErr != nil {
		return nil, fmt.Errorf(constants.CommitInErrSourceMkdir, statErr)
	}
	return nil, fmt.Errorf("commit-in: source: %q is not a directory", abs)
}

// isGitURL matches spec §2.3 case 1 prefixes. Trailing `.git` is
// optional; the caller strips it when computing the clone target.
func isGitURL(s string) bool {
	prefixes := []string{
		constants.CommitInUrlPrefixHttps,
		constants.CommitInUrlPrefixHttp,
		constants.CommitInUrlPrefixSshAt,
		constants.CommitInUrlPrefixSsh,
		constants.CommitInUrlPrefixGit,
	}
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// resolveByClone handles spec §2.3 case 1.
func resolveByClone(url string) (*SourceHandle, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getwd: %w", err)
	}
	target := filepath.Join(cwd, cloneBasename(url))
	if err := runGitClone(url, target); err != nil {
		return nil, fmt.Errorf(constants.CommitInErrSourceClone, err)
	}
	return &SourceHandle{Path: target, Kind: SourceKindCloned}, nil
}

// resolveExistingDir handles spec §2.3 cases 2 and 3.
func resolveExistingDir(abs string) (*SourceHandle, error) {
	if hasGitMetadata(abs) {
		return &SourceHandle{Path: abs, Kind: SourceKindExistingRepo}, nil
	}
	if err := runGitInit(abs); err != nil {
		return nil, fmt.Errorf(constants.CommitInErrSourceInit, err)
	}
	return &SourceHandle{Path: abs, Kind: SourceKindInitInPlace, IsFreshlyInit: true}, nil
}

// resolveMissingDir handles spec §2.3 case 4.
func resolveMissingDir(abs string) (*SourceHandle, error) {
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf(constants.CommitInErrSourceMkdir, err)
	}
	if err := runGitInit(abs); err != nil {
		return nil, fmt.Errorf(constants.CommitInErrSourceInit, err)
	}
	return &SourceHandle{Path: abs, Kind: SourceKindCreatedAndInit, IsFreshlyInit: true}, nil
}

// hasGitMetadata returns true when the directory is a working tree
// (has `.git/`) or a bare repo (has `HEAD` + `objects/`).
func hasGitMetadata(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return true
	}
	_, headErr := os.Stat(filepath.Join(dir, "HEAD"))
	_, objErr := os.Stat(filepath.Join(dir, "objects"))
	return headErr == nil && objErr == nil
}

// cloneBasename derives the on-disk folder name from a git URL by
// taking the final path segment and trimming a single trailing `.git`.
func cloneBasename(url string) string {
	trimmed := strings.TrimSuffix(url, "/")
	idx := strings.LastIndexAny(trimmed, "/:")
	last := trimmed[idx+1:]
	return strings.TrimSuffix(last, constants.CommitInUrlSuffixGit)
}

// runGitClone shells out to `git clone <url> <target>`. Replaceable
// in tests via gitRunner (see workspace_testhooks.go).
func runGitClone(url, target string) error {
	return gitRunner("clone", url, target)
}

// runGitInit invokes `git -C <dir> init` via the swappable runner so
// tests can intercept it without spawning a real git process.
func runGitInit(dir string) error {
	return gitRunner("-C", dir, "init")
}
