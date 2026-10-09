package cmdautofix

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	gofmtPath     string
	gofmtPathOnce sync.Once
	gofmtPathErr  error
)

// resolveGofmt locates the gofmt binary once per process.
func resolveGofmt() error {
	gofmtPathOnce.Do(func() {
		gofmtPath, gofmtPathErr = exec.LookPath("gofmt")
		if gofmtPathErr != nil {
			gofmtPathErr = fmt.Errorf("gofmt not found in PATH — install the Go toolchain: %w", gofmtPathErr)
		}
	})
	return gofmtPathErr
}

// catToolCheck runs the pre-walk tool availability check for exec categories.
func catToolCheck(category string) error {
	if category == "gofmt" {
		if err := resolveGofmt(); err != nil {
			return fmt.Errorf("%w: %v", errToolMissing, err)
		}
	}
	return nil
}

// gofmtDiff runs `gofmt -d` on the file; any diff output is a violation.
// D8: exec gofmt (NOT go/format stdlib) so output matches the repo toolchain.
func gofmtDiff(absPath string) (string, error) {
	cmd := exec.Command(gofmtPath, "-d", absPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stdout.Len() == 0 {
			return "", fmt.Errorf("gofmt -d failed: %s", firstLine(stderr.String()))
		}
	}
	return stdout.String(), nil
}

// gofmtCheck is the Exec-category Check: it runs against the file on disk
// (the engine flushes byte-level fixes first, so disk == buffer here).
func gofmtCheck(relPath string, src []byte, opts Options) []Violation {
	absPath := absFromRel(opts.Root, relPath)
	diff, err := gofmtDiff(absPath)
	if err != nil {
		return []Violation{{Path: relPath, Category: "gofmt", Detail: err.Error()}}
	}
	if diff == "" {
		return nil
	}
	return []Violation{{Path: relPath, Category: "gofmt", Detail: "gofmt diff present (needs formatting)"}}
}

// gofmtFix runs `gofmt -w` only when a diff exists (no-op rewrites must not
// touch mtime), then re-reads the file so the engine sees the final bytes.
func gofmtFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	absPath := absFromRel(opts.Root, relPath)
	diff, err := gofmtDiff(absPath)
	if err != nil {
		return src, []Violation{{Path: relPath, Category: "gofmt", Detail: err.Error()}}
	}
	if diff == "" {
		return src, nil
	}
	cmd := exec.Command(gofmtPath, "-w", absPath)
	if out, werr := cmd.CombinedOutput(); werr != nil {
		return src, []Violation{{Path: relPath, Category: "gofmt", Detail: "gofmt -w failed: " + firstLine(string(out))}}
	}
	fixed, rerr := os.ReadFile(absPath)
	if rerr != nil {
		return src, []Violation{{Path: relPath, Category: "gofmt", Detail: "gofmt -w succeeded but re-read failed: " + rerr.Error()}}
	}
	return fixed, nil
}

func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	return s
}
