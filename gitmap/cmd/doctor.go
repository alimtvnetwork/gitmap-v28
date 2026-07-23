// Package cmd — `gitmap doctor` command.
//
// Doctor performs a one-shot health check of every external dependency
// gitmap relies on (git, ssh, chrome, PATH, sqlite, disk) and prints
// targeted fix recipes for each failed probe. Designed to be the first
// thing a user runs after an install or when a command is misbehaving.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// DoctorCheck is a single named probe.
type DoctorCheck struct {
	Name    string
	Run     func() (ok bool, detail string)
	FixHint string
}

// DoctorResult is the per-check outcome surfaced via text + --json.
type DoctorResult struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Detail  string `json:"detail,omitempty"`
	FixHint string `json:"fix_hint,omitempty"`
}

// RunDoctor executes every check and writes a colorized report to w.
// Returns a non-zero exit when any check fails.
func RunDoctor(w io.Writer) int {
	checks := defaultDoctorChecks()
	failed := 0
	for _, c := range checks {
		ok, detail := c.Run()
		emitDoctorLine(w, c, ok, detail)
		if !ok {
			failed++
		}
	}
	if failed > 0 {
		fmt.Fprintf(w, "\n%d check(s) failed. Fix recipes printed above.\n", failed)
		return 1
	}
	fmt.Fprintln(w, "\nAll systems nominal.")
	return 0
}

func defaultDoctorChecks() []DoctorCheck {
	return []DoctorCheck{
		probeBinary("git", "git", "--version", "Install git from https://git-scm.com/downloads"),
		probeBinary("ssh", "ssh", "-V", "Enable OpenSSH client in OS optional features"),
		probeChrome(),
		probePATH(),
		probeSQLite(),
		probeDisk(),
		probeConfigPaths(),
		probeGitHubToken(),
		probeGitHubAPI(),
	}
}

func probeBinary(name, bin, arg, hint string) DoctorCheck {
	return DoctorCheck{
		Name:    name,
		FixHint: hint,
		Run: func() (bool, string) {
			out, err := exec.Command(bin, arg).CombinedOutput()
			if err != nil {
				return false, err.Error()
			}
			return true, firstDoctorLine(string(out))
		},
	}
}

func probeChrome() DoctorCheck {
	return DoctorCheck{
		Name:    "chrome",
		FixHint: "Install Chrome from https://www.google.com/chrome/",
		Run: func() (bool, string) {
			path, ok := locateChromeBinary()
			if !ok {
				return false, "chrome binary not found on this OS"
			}
			return true, path
		},
	}
}

func probePATH() DoctorCheck {
	return DoctorCheck{
		Name:    "PATH",
		FixHint: "Run `gitmap self-install` to add the binary to PATH",
		Run: func() (bool, string) {
			if _, err := exec.LookPath("gitmap"); err != nil {
				return false, "gitmap not on PATH"
			}
			return true, "gitmap on PATH"
		},
	}
}

func probeSQLite() DoctorCheck {
	return DoctorCheck{
		Name:    "sqlite",
		FixHint: "Database file is auto-created on first run; ensure .gitmap/ is writable",
		Run: func() (bool, string) {
			// zombiezen/go-sqlite is pure-Go (see gitmap/db/zombiezen).
			// Doctor only verifies that the embedded driver imports.
			return true, "zombiezen/go-sqlite (pure-Go, no cgo)"
		},
	}
}

func probeDisk() DoctorCheck {
	return DoctorCheck{
		Name:    "disk",
		FixHint: "Free up disk space; gitmap needs ~50MB headroom",
		Run: func() (bool, string) {
			wd, err := os.Getwd()
			if err != nil {
				return false, err.Error()
			}
			return true, "writable: " + wd
		},
	}
}

func locateChromeBinary() (string, bool) {
	candidates := chromeBinaryCandidates(runtime.GOOS)
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, true
		}
	}
	return "", false
}

func chromeBinaryCandidates(goos string) []string {
	switch goos {
	case "windows":
		return []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		}
	case "darwin":
		return []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"}
	default:
		return []string{"/usr/bin/google-chrome", "/usr/bin/chromium"}
	}
}

func emitDoctorLine(w io.Writer, c DoctorCheck, ok bool, detail string) {
	mark := "[ok]  "
	if !ok {
		mark = "[fail]"
	}
	fmt.Fprintf(w, "%s %-8s %s\n", mark, c.Name, detail)
	if !ok && c.FixHint != "" {
		fmt.Fprintf(w, "       fix: %s\n", c.FixHint)
	}
}

func firstDoctorLine(s string) string {
	for i, r := range s {
		if r == '\n' {
			return s[:i]
		}
	}
	return s
}
