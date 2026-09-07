package cmd

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var cgPostfixIncrementPattern = regexp.MustCompile(`\(\(([A-Za-z_][A-Za-z0-9_]*)\+\+\)\)`)

func writeCGCompatScript(url string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "gitmap-cg-*")
	if err != nil {
		return "", func() {}, err
	}

	path := filepath.Join(dir, "install.sh")
	if err := downloadCGScript(url, path); err != nil {
		_ = os.RemoveAll(dir)
		return "", func() {}, err
	}

	return path, func() { _ = os.RemoveAll(dir) }, patchCGScriptFile(path)
}

func writeCGCompatScriptWindows(url string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "gitmap-cg-win-*")
	if err != nil {
		return "", func() {}, err
	}

	path := filepath.Join(dir, "install.ps1")
	if err := downloadCGScript(url, path); err != nil {
		_ = os.RemoveAll(dir)
		return "", func() {}, err
	}

	return path, func() { _ = os.RemoveAll(dir) }, patchCGWindowsScriptFile(path)
}

func downloadCGScript(url, path string) error {
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		return copyHTTPBodyToFile(resp.Body, path)
	}
	if resp != nil {
		_ = resp.Body.Close()
	}

	cmd := exec.Command("curl", "-fsSL", url, "-o", path)

	return cmd.Run()
}

func copyHTTPBodyToFile(body io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, body)

	return err
}

func patchCGScriptFile(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	patched := patchCGArithmeticIncrements(string(body))
	patched = strings.ReplaceAll(patched, "local srchash", "")

	return os.WriteFile(path, []byte(patched), 0o700)
}

func patchCGWindowsScriptFile(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	patched := string(body)
	patched = strings.ReplaceAll(patched, "$oldFile:", "${oldFile}:")
	patched = strings.ReplaceAll(patched, "$destPath:", "${destPath}:")
	patched = strings.ReplaceAll(patched, "$targetVersionFile:", "${targetVersionFile}:")

	data := addUTF8BOM([]byte(patched))

	return os.WriteFile(path, data, 0644)
}

func addUTF8BOM(data []byte) []byte {
	if !bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		return append([]byte{0xef, 0xbb, 0xbf}, data...)
	}

	return data
}

func patchCGArithmeticIncrements(script string) string {
	return cgPostfixIncrementPattern.ReplaceAllString(script, `((${1}+=1))`)
}
