package cmdautomation

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ProbeRuntime probes the runtime with two-tier caching and system detection.
func ProbeRuntime(name string) result.Result[RuntimeRecord] {
	dbRes := OpenAutomationDB("")
	if dbRes.IsSuccess() {
		defer func() { _ = dbRes.Value.Close() }()
		return ProbeRuntimeWithDB(dbRes.Value, name)
	}
	return ProbeRuntimeWithDB(nil, name)
}

// ProbeRuntimeWithDB executes the two-tier runtime discovery protocol.
func ProbeRuntimeWithDB(db *sql.DB, name string) result.Result[RuntimeRecord] {
	norm := normalizeRuntimeKey(name)
	if cached := checkCachedRuntime(db, norm); cached.IsSuccess() {
		return cached
	}
	probed := probeSystemRuntime(norm)
	persistProbeResult(db, norm, probed)
	return probed
}

func checkCachedRuntime(db *sql.DB, norm string) result.Result[RuntimeRecord] {
	if db == nil {
		return result.Fail[RuntimeRecord](apperror.NewSimple("nil db", "E7102"))
	}
	res := GetCachedRuntime(db, norm)
	if res.IsSuccess() && isRuntimeCacheValid(res.Value) {
		return res
	}
	return result.Fail[RuntimeRecord](apperror.NewSimple("cache miss", "E7103"))
}

func isRuntimeCacheValid(rec RuntimeRecord) bool {
	if !rec.IsValid || rec.BinaryPath == "" || !isFileExisting(rec.BinaryPath) {
		return false
	}
	t, err := time.Parse(time.RFC3339, rec.LastVerifiedAt)
	return err == nil && time.Since(t) <= 24*time.Hour
}

func probeSystemRuntime(norm string) result.Result[RuntimeRecord] {
	bin := detectBinaryPath(norm)
	if bin == "" {
		return result.Fail[RuntimeRecord](NewMissingRuntimeError(norm))
	}
	return result.Ok(buildRuntimeRecord(norm, bin, "active", true))
}

func buildRuntimeRecord(normName, binPath, status string, isValid bool) RuntimeRecord {
	meta := runtimeMetaFor(normName)
	ver, binName := "missing", normName
	if isValid {
		ver, binName = queryRuntimeVersion(binPath), filepath.Base(binPath)
	}
	return RuntimeRecord{
		Name: normName, BinaryName: binName, BinaryPath: binPath, Version: ver,
		Status: status, InstallCmd: meta.InstallCmd, ProfileSuggestion: meta.ProfileSuggestion,
		FallbackCmd: meta.FallbackCmd, IsValid: isValid,
	}
}

func persistProbeResult(db *sql.DB, normName string, probed result.Result[RuntimeRecord]) {
	if db == nil {
		return
	}
	rec := buildRuntimeRecord(normName, "", "missing", false)
	if probed.IsSuccess() {
		rec = probed.Value
	}
	_ = SaveCachedRuntime(db, rec)
}

func detectBinaryPath(normName string) string {
	cands := candidateBinaryNames(normName)
	for _, cand := range cands {
		if path, err := exec.LookPath(cand); err == nil && path != "" {
			return filepath.ToSlash(path)
		}
	}
	return checkStandardLocations(normName, cands)
}

func candidateBinaryNames(norm string) []string {
	cands := map[string][]string{
		"python": {"python.exe", "python3.exe", "py.exe", "python3", "python"},
		"node":   {"node.exe", "node"},
		"go":     {"go.exe", "go"},
		"rust":   {"cargo.exe", "rustc.exe", "cargo", "rustc"},
		"pwsh":   {"pwsh.exe", "powershell.exe", "pwsh"},
		"bash":   {"bash.exe", "bash"},
	}
	if list, ok := cands[norm]; ok {
		return list
	}
	return []string{norm}
}

func checkStandardLocations(normName string, cands []string) string {
	for _, dir := range standardDirs() {
		for _, cand := range cands {
			if full := filepath.Join(dir, cand); isFileExisting(full) {
				return filepath.ToSlash(full)
			}
		}
	}
	return ""
}

func standardDirs() []string {
	if runtime.GOOS != "windows" {
		return []string{"/usr/local/bin", "/usr/bin", "/bin", "/opt/homebrew/bin"}
	}
	u, l, p := os.Getenv("USERPROFILE"), os.Getenv("LOCALAPPDATA"), os.Getenv("ProgramFiles")
	return []string{
		filepath.Join(p, "PowerShell", "7"), filepath.Join(p, "nodejs"),
		filepath.Join(p, "Go", "bin"), filepath.Join(p, "Git", "bin"),
		filepath.Join(u, ".cargo", "bin"),
		filepath.Join(l, "Programs", "Python", "Python312"),
	}
}

func isFileExisting(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func queryRuntimeVersion(binPath string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binPath, "--version")
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return "detected"
	}
	return strings.TrimSpace(string(out))
}

// FormatMissingRuntimeMessage formats the missing runtime suggestion UI per Spec 124 Section 3.2.
func FormatMissingRuntimeMessage(name string) string {
	norm := normalizeRuntimeKey(name)
	meta := runtimeMetaFor(norm)
	sugg := fmt.Sprintf("Suggested Installation via GitMap:\n  %s\n  gitmap install profile dev-full\n", meta.InstallCmd)
	if meta.ProfileSuggestion != "" && meta.ProfileSuggestion != "gitmap install profile dev-full" {
		sugg += fmt.Sprintf("  %s\n", meta.ProfileSuggestion)
	}
	tpl := "[E7100:RUNTIME_MISSING] %s runtime not detected on host system.\n\n%s\nAlternative System Package Managers:\n%s\nTip: Configure automatic bootstrap in config:\n  gitmap automation config set autoInstallMissingRuntimes true\n"
	return fmt.Sprintf(tpl, titleCase(norm), sugg, meta.FallbackCmd)
}

// NewMissingRuntimeError creates an AppError with code E7100:RUNTIME_MISSING.
func NewMissingRuntimeError(name string) *apperror.AppError {
	return apperror.NewSimple(FormatMissingRuntimeMessage(name), "E7100:RUNTIME_MISSING")
}

func titleCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func runtimeMetaFor(norm string) runtimeMeta {
	cmd, prof := "gitmap install "+norm, "gitmap install profile dev-full"
	if norm == "bash" {
		cmd = "gitmap install git"
	}
	if norm == "python" {
		prof = "gitmap install profile python-data"
	}
	if norm == "node" {
		prof = "gitmap install profile web-full"
	}
	return runtimeMeta{InstallCmd: cmd, ProfileSuggestion: prof, FallbackCmd: formatFallbackCmd(norm)}
}

// ListRuntimes probes or reads all supported polyglot runtimes.
func ListRuntimes() RuntimeListResultMonad {
	known := []string{"python", "node", "go", "rust", "pwsh", "bash"}
	var list []RuntimeRecord
	for _, k := range known {
		res := ProbeRuntime(k)
		if res.IsSuccess() {
			list = append(list, res.Value)
		} else {
			list = append(list, buildRuntimeRecord(k, "", "missing", false))
		}
	}
	return result.Ok(list)
}

func reprobeRuntime(db *sql.DB, norm string) RuntimeRecord {
	probed := probeSystemRuntime(norm)
	persistProbeResult(db, norm, probed)
	if probed.IsSuccess() {
		return probed.Value
	}
	return buildRuntimeRecord(norm, "", "missing", false)
}

// RefreshRuntimes forces re-probing of all supported polyglot runtimes.
func RefreshRuntimes() RuntimeListResultMonad {
	known := []string{"python", "node", "go", "rust", "pwsh", "bash"}
	var list []RuntimeRecord
	dbRes := OpenAutomationDB("")
	var db *sql.DB
	if dbRes.IsSuccess() {
		db = dbRes.Value
		defer func() { _ = db.Close() }()
	}
	for _, k := range known {
		list = append(list, reprobeRuntime(db, normalizeRuntimeKey(k)))
	}
	return result.Ok(list)
}
