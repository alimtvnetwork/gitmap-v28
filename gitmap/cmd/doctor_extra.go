// Package cmd — additional `gitmap doctor` probes:
// config paths, GitHub token availability, and network/API connectivity.
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const doctorHTTPTimeoutSecs = 5

// doctorGitHubEndpoints is the list of hosts gh-api probes.
// Adding a URL here automatically extends the report.
var doctorGitHubEndpoints = []struct{ Name, URL string }{
	{"api.github.com", "https://api.github.com"},
	{"github.com", "https://github.com"},
	{"codeload.github.com", "https://codeload.github.com"},
	{"uploads.github.com", "https://uploads.github.com"},
	{"objects.githubusercontent.com", "https://objects.githubusercontent.com"},
}

// probeConfigPaths verifies the .gitmap/ working directory is reachable and writable.
func probeConfigPaths() DoctorCheck {
	return DoctorCheck{
		Name:    "config",
		FixHint: "Run gitmap from a writable directory; check perms on .gitmap/",
		Run: func() (bool, string) {
			wd, err := os.Getwd()
			if err != nil {
				return false, err.Error()
			}
			dir := filepath.Join(wd, ".gitmap")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return false, "cannot create " + dir + ": " + err.Error()
			}
			probe := filepath.Join(dir, ".doctor-probe")
			if err := os.WriteFile(probe, []byte("ok"), 0o644); err != nil {
				return false, "not writable: " + dir
			}
			_ = os.Remove(probe)
			return true, "writable: " + dir
		},
	}
}

// probeGitHubToken checks for GITHUB_TOKEN or GH_TOKEN in the environment.
// Missing token is a warning (ok=true with note) so anonymous flows still pass.
func probeGitHubToken() DoctorCheck {
	return DoctorCheck{
		Name:    "gh-token",
		FixHint: "export GITHUB_TOKEN=<pat>  # needed for private repos & higher rate limits",
		Run: func() (bool, string) {
			for _, k := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
				if v := os.Getenv(k); v != "" {
					return true, k + " set (" + maskToken(v) + ")"
				}
			}
			return false, "no GITHUB_TOKEN / GH_TOKEN in environment"
		},
	}
}

// probeGitHubAPI verifies network + GitHub API reachability with a short timeout.
func probeGitHubAPI() DoctorCheck {
	return DoctorCheck{
		Name:    "gh-api",
		FixHint: "Check internet / proxy / firewall; GitHub hosts must be reachable",
		Run: func() (bool, string) {
			client := &http.Client{Timeout: doctorHTTPTimeoutSecs * time.Second}
			tok := firstEnv("GITHUB_TOKEN", "GH_TOKEN")
			var oks, fails []string
			for _, ep := range doctorGitHubEndpoints {
				if ok, msg := probeGitHubEndpoint(client, ep.URL, tok); ok {
					oks = append(oks, ep.Name+" "+msg)
				} else {
					fails = append(fails, ep.Name+" "+msg)
				}
			}
			summary := fmt.Sprintf("%d/%d reachable", len(oks), len(doctorGitHubEndpoints))
			if len(fails) > 0 {
				return false, summary + "; failures: " + joinSemi(fails)
			}
			return true, summary + "; " + joinSemi(oks)
		},
	}
}

func probeGitHubEndpoint(client *http.Client, url, token string) (bool, string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err.Error()
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return false, resp.Status
	}
	return true, "(" + resp.Status + ")"
}

func joinSemi(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += "; "
		}
		out += p
	}
	return out
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func maskToken(t string) string {
	if len(t) <= 8 {
		return "****"
	}
	return t[:4] + "…" + t[len(t)-4:]
}
