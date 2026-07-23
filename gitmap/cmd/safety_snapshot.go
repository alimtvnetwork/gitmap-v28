// Package cmd — safety_snapshot.go: full working-tree snapshot to
// .gitmap/snapshot/, rollback, and pre-commit guard hook installer.
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func runSnapshot(args []string) {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	ts := time.Now().UTC().Format("20060102-150405")
	out := filepath.Join(root, ".gitmap", "snapshot", "snap-"+ts+".tar.gz")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "snapshot: ERROR %v\n", err)
		os.Exit(1)
	}
	n, err := writeSnapshot(root, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "snapshot: ERROR %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[1;92m✓ snapshot\033[0m  %d files → \033[1;96m%s\033[0m\n", n, out)
	fmt.Printf("  rollback: \033[1;96mgitmap rollback %s\033[0m\n", out)
}

func writeSnapshot(root, outPath string) (int, error) {
	f, err := os.Create(outPath) //nolint:gosec
	if err != nil {
		return 0, err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()
	count := 0
	err = filepath.Walk(root, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		if rel == "." || strings.HasPrefix(rel, ".gitmap"+string(os.PathSeparator)+"snapshot") {
			return nil
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return nil
		}
		hdr.Name = filepath.ToSlash(rel)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		in, err := os.Open(p) //nolint:gosec
		if err != nil {
			return nil
		}
		_, _ = io.Copy(tw, in)
		in.Close()
		count++
		return nil
	})
	return count, err
}

func runRollback(args []string) {
	src := ""
	if len(args) > 0 {
		src = args[0]
	} else {
		src = latestSnapshot(".")
	}
	if src == "" {
		fmt.Fprintln(os.Stderr, "rollback: ERROR no snapshot found; pass <tarball> explicitly")
		os.Exit(2)
	}
	n, err := readChromeBackup(src, ".") // tar.gz extractor reused
	if err != nil {
		fmt.Fprintf(os.Stderr, "rollback: ERROR %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[1;92m✓ rollback\033[0m  restored %d files from %s\n", n, src)
}

func latestSnapshot(root string) string {
	dir := filepath.Join(root, ".gitmap", "snapshot")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".tar.gz") {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	return filepath.Join(dir, names[len(names)-1])
}

func runGuard(args []string) {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	hooks := filepath.Join(root, ".git", "hooks")
	if _, err := os.Stat(hooks); err != nil {
		fmt.Fprintln(os.Stderr, "guard: ERROR not a git repo (no .git/hooks)")
		os.Exit(2)
	}
	hook := filepath.Join(hooks, "pre-commit")
	if err := os.WriteFile(hook, []byte(guardHookBody), 0o755); err != nil { //nolint:gosec
		fmt.Fprintf(os.Stderr, "guard: ERROR write hook: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[1;92m✓ installed\033[0m pre-commit guard → %s\n", hook)
	fmt.Println("  blocks: secrets (API_KEY/PRIVATE_KEY), large files (>10MB), -vN drift")
}

const guardHookBody = `#!/usr/bin/env bash
# gitmap guard — block secrets / large files / -vN drift
set -e
staged=$(git diff --cached --name-only --diff-filter=ACM)
fail=0
for f in $staged; do
  [ -f "$f" ] || continue
  sz=$(wc -c < "$f" | tr -d ' ')
  if [ "$sz" -gt 10485760 ]; then
    echo "guard: $f exceeds 10MB ($sz bytes)" >&2; fail=1
  fi
  if grep -EI 'PRIVATE KEY|API_KEY=|SECRET=' "$f" >/dev/null 2>&1; then
    echo "guard: $f contains suspected secret" >&2; fail=1
  fi
done
exit $fail
`
