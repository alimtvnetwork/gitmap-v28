package cmdtoolchain

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ResolveLibclangDir discovers the libclang directory in deterministic tier order.
func ResolveLibclangDir() string {
	if dir := resolveLibclangFromEnv(); dir != "" {
		return dir
	}

	if dir := resolveLibclangFromClangPath(); dir != "" {
		return dir
	}

	return resolveLibclangFromPlatformCandidates()
}

func resolveLibclangFromEnv() string {
	val := os.Getenv("LIBCLANG_PATH")
	if val == "" {
		return ""
	}

	if isValidLibclangDir(val) {
		return val
	}

	if isValidLibclangDir(filepath.Dir(val)) {
		return filepath.Dir(val)
	}

	return ""
}

func resolveLibclangFromClangPath() string {
	p, err := exec.LookPath("clang")
	if err != nil {
		return ""
	}

	dir := filepath.Dir(p)
	if isValidLibclangDir(dir) {
		return dir
	}

	libSibling := filepath.Join(filepath.Dir(dir), "lib")
	if isValidLibclangDir(libSibling) {
		return libSibling
	}

	return ""
}

func resolveLibclangFromPlatformCandidates() string {
	candidates := getPlatformLibclangCandidates()
	for _, dir := range candidates {
		if isValidLibclangDir(dir) {
			return dir
		}
	}

	return ""
}

func getPlatformLibclangCandidates() []string {
	if runtime.GOOS == "windows" {
		return getWindowsLibclangCandidates()
	}

	return getPosixLibclangCandidates()
}

func getWindowsLibclangCandidates() []string {
	return []string{
		`D:\LLVM\LLVM\bin`,
		`D:\LLVM\bin`,
		`C:\Program Files\LLVM\bin`,
		`C:\LLVM\bin`,
		`C:\Program Files\LLVM\lib`,
	}
}

func getPosixLibclangCandidates() []string {
	return []string{
		"/usr/lib/llvm-18/lib",
		"/usr/lib/llvm-17/lib",
		"/usr/lib/llvm-16/lib",
		"/usr/lib/llvm-15/lib",
		"/opt/homebrew/opt/llvm/lib",
		"/usr/local/opt/llvm/lib",
		"/usr/lib",
	}
}

func isValidLibclangDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return hasLibclangFile(dir)
}

func hasLibclangFile(dir string) bool {
	fileNames := getLibclangFileNames()
	for _, name := range fileNames {
		p := filepath.Join(dir, name)
		if isRegularFile(p) {
			return true
		}
	}

	return false
}

func getLibclangFileNames() []string {
	if runtime.GOOS == "windows" {
		return []string{"libclang.dll", "libclang.lib"}
	}

	if runtime.GOOS == "darwin" {
		return []string{"libclang.dylib"}
	}

	return []string{"libclang.so", "libclang.so.1", "libclang.so.18", "libclang.so.17"}
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
