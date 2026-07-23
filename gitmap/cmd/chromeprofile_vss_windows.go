//go:build windows

// chromeprofile_vss_windows.go — Volume Shadow Copy snapshot helper
// for Chrome profile copy on Windows (#8).
//
// When Chrome is open it holds an exclusive lock on `LOCK`, `Cookies`,
// `History`, and friends. The existing copy path skips such files; VSS
// sidesteps the lock entirely by reading from a point-in-time snapshot
// of the volume.
//
// Strategy:
//  1. `vssadmin create shadow /for=<volume>` (requires Administrator).
//  2. Parse `GLOBALROOT\Device\HarddiskVolumeShadowCopyN` from output.
//  3. Translate the source path under the snapshot.
//  4. Caller copies from the snapshot path; deferred cleanup deletes
//     the shadow with `vssadmin delete shadows /shadow=<id>`.
//
// Graceful degradation: any failure (not-admin, vssadmin missing,
// non-NTFS volume) returns ok=false so the caller falls back to the
// existing skip-list path. No user-visible error in that case.
package cmd

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Snapshot represents an active VSS shadow copy.
type Snapshot struct {
	ID         string // {GUID} for `vssadmin delete shadows /shadow=`
	DevicePath string // \\?\GLOBALROOT\Device\HarddiskVolumeShadowCopyN
	Volume     string // e.g. "C:\"
}

// CreateSnapshot creates a VSS shadow copy for the volume hosting
// `anyPath` and returns the snapshot handle. ok=false on any error
// — caller MUST treat that as "fall back to skip-list copy".
func CreateSnapshot(anyPath string) (Snapshot, bool) {
	volume := filepath.VolumeName(anyPath) + `\`
	if volume == `\` {
		return Snapshot{}, false
	}
	out, err := exec.Command("vssadmin", "create", "shadow", "/for="+volume).CombinedOutput()
	if err != nil {
		return Snapshot{}, false
	}
	id := parseVSSField(string(out), `Shadow Copy ID:\s*(\{[A-F0-9-]+\})`)
	dev := parseVSSField(string(out), `Shadow Copy Volume Name:\s*(\\\\\?\\GLOBALROOT\\[^\r\n]+)`)
	if id == "" || dev == "" {
		return Snapshot{}, false
	}
	return Snapshot{ID: id, DevicePath: strings.TrimSpace(dev), Volume: volume}, true
}

// TranslatePath rewrites a source path so it points inside the
// snapshot. e.g. C:\Users\x → \\?\GLOBALROOT\...\Users\x.
func (s Snapshot) TranslatePath(srcAbs string) string {
	if s.DevicePath == "" {
		return srcAbs
	}
	rel := strings.TrimPrefix(srcAbs, filepath.VolumeName(srcAbs))
	rel = strings.TrimPrefix(rel, `\`)
	return s.DevicePath + `\` + rel
}

// Delete removes the snapshot. Best-effort; errors are swallowed
// (caller already has the data it needed).
func (s Snapshot) Delete() {
	if s.ID == "" {
		return
	}
	_ = exec.Command("vssadmin", "delete", "shadows", "/shadow="+s.ID, "/quiet").Run()
}

func parseVSSField(out, pattern string) string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(out)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}
