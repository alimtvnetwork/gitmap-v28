// Package archive — source resolution helpers used by the cmd layer to
// turn user-supplied strings (local paths, HTTPS URLs, git URLs) into
// concrete on-disk paths the extract / create engines can consume.
//
// Network operations are deliberately kept small: we shell out to aria2c
// when available, fall back to net/http otherwise, and shell out to git
// for clone. The downloader package will replace this with the full
// engine in a later slice — until then this keeps `gitmap uzc <url>` and
// `gitmap zip <git-url>` functional.
package archive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// SourceKind classifies one entry on the `gitmap zip` / `gitmap uzc`
// command line. The cmd layer dispatches per-kind; the archive engine
// only ever sees concrete local paths.
type SourceKind int

const (
	SourceLocal SourceKind = iota
	SourceHTTP
	SourceGit
)

// ResolvedSource is the materialized form of one user-supplied input.
// LocalPath is always populated; CleanupDir, when non-empty, must be
// removed by the caller after the operation completes (this is how the
// HTTP and git branches signal they used a temp workspace).
type ResolvedSource struct {
	Original   string
	Kind       SourceKind
	LocalPath  string
	CleanupDir string
}

// ClassifySource is the cheap, pure-function classifier the command
// layer uses BEFORE doing any IO. Decision order matters: a path like
// `git@github.com:foo/bar.git` parses as a URL with no scheme, so we
// detect git first.
func ClassifySource(s string) SourceKind {
	if isGitURL(s) {
		return SourceGit
	}
	if isHTTPURL(s) {
		return SourceHTTP
	}

	return SourceLocal
}

func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}

// isGitURL covers SSH-style (git@host:owner/repo.git) and any HTTPS URL
// ending in .git. Plain HTTPS that happens to host a tarball is left to
// SourceHTTP.
func isGitURL(s string) bool {
	if strings.HasPrefix(s, "git@") && strings.Contains(s, ":") {
		return true
	}
	if strings.HasPrefix(s, "git://") {
		return true
	}
	if isHTTPURL(s) && strings.HasSuffix(strings.ToLower(s), ".git") {
		return true
	}

	return false
}

// ResolveSource turns one input string into a usable local path. The
// caller is responsible for invoking CleanupResolved afterwards.
func ResolveSource(ctx context.Context, raw string) (ResolvedSource, error) {
	switch ClassifySource(raw) {
	case SourceLocal:
		abs, err := filepath.Abs(raw)
		if err != nil {
			return ResolvedSource{Original: raw}, err
		}
		if _, err := os.Stat(abs); err != nil {
			return ResolvedSource{Original: raw}, fmt.Errorf("local source %q: %w", raw, err)
		}

		return ResolvedSource{Original: raw, Kind: SourceLocal, LocalPath: abs}, nil

	case SourceHTTP:
		return resolveHTTP(ctx, raw)

	case SourceGit:
		return resolveGit(ctx, raw)
	}

	return ResolvedSource{Original: raw}, errors.New("unsupported source kind")
}

// CleanupResolved removes any temp workspace recorded on the source.
// Always safe to call.
func CleanupResolved(r ResolvedSource) {
	if r.CleanupDir != "" {
		_ = os.RemoveAll(r.CleanupDir)
	}
}

// resolveHTTP downloads raw into a temp directory using aria2c when
// available, falling back to net/http. The returned LocalPath points at
// the downloaded file; CleanupDir is the temp workspace.
func resolveHTTP(ctx context.Context, raw string) (ResolvedSource, error) {
	dir, err := os.MkdirTemp("", "gitmap-fetch-*")
	if err != nil {
		return ResolvedSource{Original: raw}, err
	}

	name := filenameFromURL(raw)
	dst := filepath.Join(dir, name)

	if err := downloadWithAria2c(ctx, raw, dir, name); err == nil {
		return ResolvedSource{Original: raw, Kind: SourceHTTP, LocalPath: dst, CleanupDir: dir}, nil
	}

	if err := downloadWithHTTP(ctx, raw, dst); err != nil {
		_ = os.RemoveAll(dir)

		return ResolvedSource{Original: raw}, fmt.Errorf("download %q: %w", raw, err)
	}

	return ResolvedSource{Original: raw, Kind: SourceHTTP, LocalPath: dst, CleanupDir: dir}, nil
}

// downloadWithAria2c is a thin wrapper that returns nil only on a clean
// aria2c exit AND a non-empty file. Any other outcome causes the caller
// to fall back to net/http.
func downloadWithAria2c(ctx context.Context, rawURL, dir, name string) error {
	if _, err := exec.LookPath("aria2c"); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "aria2c",
		"--dir", dir,
		"--out", name,
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		"--console-log-level=warn",
		rawURL,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	info, err := os.Stat(filepath.Join(dir, name))
	if err != nil || info.Size() == 0 {
		return errors.New("aria2c produced empty file")
	}

	return nil
}

// downloadWithHTTP is the always-available fallback path.
func downloadWithHTTP(ctx context.Context, rawURL, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http %s", resp.Status)
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}

// filenameFromURL pulls the last path segment off a URL. Falls back to
// "download.bin" when the URL has no usable name (e.g. "https://x.com/").
func filenameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "download.bin"
	}
	base := path.Base(u.Path)
	if base == "" || base == "/" || base == "." {
		return "download.bin"
	}

	return base
}

// resolveGit shallow-clones the repo into a temp dir. The returned
// LocalPath is the cloned directory; CleanupDir matches it so it gets
// wiped after the calling command finishes.
func resolveGit(ctx context.Context, raw string) (ResolvedSource, error) {
	dir, err := os.MkdirTemp("", "gitmap-gitsrc-*")
	if err != nil {
		return ResolvedSource{Original: raw}, err
	}

	cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", raw, dir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)

		return ResolvedSource{Original: raw}, fmt.Errorf("git clone %q: %w", raw, err)
	}

	return ResolvedSource{Original: raw, Kind: SourceGit, LocalPath: dir, CleanupDir: dir}, nil
}

// AutoDetectSingleArchive scans dir for exactly one file with a
// recognized archive extension. Returns the absolute path on success, an
// error describing 0 or N>1 matches otherwise. Used by `gitmap uzc` when
// the user passes no explicit source.
func AutoDetectSingleArchive(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return "", err
	}

	var found []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if FormatFromPath(e.Name()) != FormatUnknown {
			found = append(found, filepath.Join(abs, e.Name()))
		}
	}
	switch len(found) {
	case 0:
		return "", fmt.Errorf("no archive in %s", abs)
	case 1:
		return found[0], nil
	default:
		return "", fmt.Errorf("found %d archives in %s", len(found), abs)
	}
}
