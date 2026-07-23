// Package cmd — chrome_backup.go: snapshot all Chrome profiles into a
// tar.gz under .gitmap/chrome/backup/, and restore from one.
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func runChromeBackup(args []string) {
	out := chromeBackupDefaultPath()
	for i := 0; i < len(args); i++ {
		if args[i] == "-o" || args[i] == "--out" {
			if i+1 < len(args) {
				out = args[i+1]
				i++
			}
		}
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "chrome backup: ERROR mkdir: %v\n", err)
		os.Exit(1)
	}
	srcRoot := chromeUserDataDir()
	n, err := writeChromeBackup(srcRoot, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chrome backup: ERROR %v\n", err)
		os.Exit(1)
	}
	manifestPath, mErr := writeChromeManifestWithSource(out, srcRoot)
	if mErr != nil {
		fmt.Fprintf(os.Stderr, "chrome backup: WARN manifest write failed: %v\n", mErr)
	}

	fmt.Printf("\033[1;92m✓ chrome backup\033[0m  %d files → \033[1;96m%s\033[0m\n", n, out)
	if manifestPath != "" {
		fmt.Printf("  manifest: \033[2;37m%s\033[0m\n", manifestPath)
	}
	fmt.Printf("  restore: \033[1;96mgitmap chrome restore %s\033[0m\n", out)
}

func runChromeRestore(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "chrome restore: ERROR usage: gitmap chrome restore <tarball> [--into <dir>] [--force|-f] [--yes|-y] [--dry-run] [--no-verify]")
		os.Exit(2)
	}
	src := args[0]
	dst := ""
	intoSet := false
	force, yes, dryRun, skipVerify := false, false, false, false
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--into":
			if i+1 < len(args) {
				dst = args[i+1]
				intoSet = true
				i++
			}
		case "--force", "-f":
			force = true
		case "--yes", "-y":
			yes = true
		case "--dry-run":
			dryRun = true
		case "--no-verify":
			skipVerify = true
		}
	}
	if !intoSet {
		if recorded := readChromeManifestSource(src); recorded != "" {
			dst = recorded
			fmt.Printf("\033[2;37m• restoring to recorded source profile: %s\033[0m\n", dst)
		} else {
			dst = chromeUserDataDir()
		}
	}

	if !skipVerify {
		ok, miss, err := verifyChromeManifest(src)
		switch {
		case err != nil:
			fmt.Fprintf(os.Stderr, "chrome restore: WARN checksum verify skipped: %v\n", err)
		case !ok:
			fmt.Fprintf(os.Stderr, "chrome restore: ERROR checksum mismatch in %d file(s):\n  - %s\n  rerun with --no-verify to override\n", len(miss), strings.Join(miss, "\n  - "))
			os.Exit(1)
		default:
			fmt.Printf("\033[1;92m✓ checksum verified\033[0m  %s\n", src+chromeManifestSuffix)
		}
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "chrome restore: ERROR mkdir: %v\n", err)
		os.Exit(1)
	}
	existing := countChromeProfileFiles(dst)
	if existing > 0 && !force {
		fmt.Fprintf(os.Stderr, "chrome restore: REFUSED %s already contains %d file(s); pass --force to overwrite\n", dst, existing)
		os.Exit(1)
	}
	if existing > 0 && force && !yes {
		fmt.Fprintf(os.Stderr, "\033[1;31m! chrome restore --force\033[0m will overwrite %d existing file(s) under %s\n", existing, dst)
		if !confirmYesNo("proceed?") {
			fmt.Fprintln(os.Stderr, "chrome restore: aborted")
			os.Exit(1)
		}
	}
	if dryRun {
		n, err := previewChromeBackup(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "chrome restore: ERROR %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\033[1;93m✓ chrome restore (dry-run)\033[0m  %d file(s) would land in \033[1;96m%s\033[0m\n", n, dst)
		return
	}
	n, err := readChromeBackup(src, dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chrome restore: ERROR %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[1;92m✓ chrome restore\033[0m  %d files → \033[1;96m%s\033[0m\n", n, dst)
}

// countChromeProfileFiles returns the number of regular files already
// living under dst (used to decide whether --force is required).
func countChromeProfileFiles(dst string) int {
	n := 0
	_ = filepath.Walk(dst, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && info.Mode().IsRegular() {
			n++
		}
		return nil
	})
	return n
}

// previewChromeBackup walks the tarball without writing and returns the
// number of regular file entries that would be restored.
func previewChromeBackup(src string) (int, error) {
	f, err := os.Open(src) //nolint:gosec
	if err != nil {
		return 0, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return 0, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	n := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
		if !hdr.FileInfo().IsDir() {
			n++
		}
	}
}

func chromeBackupDefaultPath() string {
	ts := time.Now().UTC().Format("20060102-150405")
	return filepath.Join(".gitmap", "chrome", "backup", "chrome-"+ts+".tar.gz")
}

func writeChromeBackup(srcRoot, outPath string) (int, error) {
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
	err = filepath.Walk(srcRoot, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil // skip unreadable
		}
		rel, _ := filepath.Rel(srcRoot, p)
		if rel == "." {
			return nil
		}
		if strings.HasSuffix(p, "LOCK") || strings.HasSuffix(p, "lockfile") {
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

func readChromeBackup(src, dstRoot string) (int, error) {
	f, err := os.Open(src) //nolint:gosec
	if err != nil {
		return 0, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return 0, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	count := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}
		clean := filepath.Clean(hdr.Name)
		if strings.HasPrefix(clean, "..") || filepath.IsAbs(clean) {
			continue // path-traversal guard
		}
		target := filepath.Join(dstRoot, clean)
		if hdr.FileInfo().IsDir() {
			_ = os.MkdirAll(target, 0o755)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(target), 0o755)
		out, err := os.Create(target) //nolint:gosec
		if err != nil {
			continue
		}
		if _, err := io.Copy(out, tr); err != nil { //nolint:gosec
			out.Close()
			return count, err
		}
		out.Close()
		count++
	}
	return count, nil
}
