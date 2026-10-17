package cmdhealth

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// healthFacts is the aggregated workstation state.
type healthFacts struct {
	repoTotal    int
	repoDirty    int
	dirtyNames   []string
	pipelineRed  int
	redNames     []string
	pipelineNote string
	fleetNodes   int
	fleetNote    string
	staleDBs     int
	staleNames   []string
}

// collectHealthFacts gathers the read-only signals health reports on.
func collectHealthFacts(workdir string) healthFacts {
	var facts healthFacts
	collectRepoFacts(workdir, &facts)
	collectPipelineFacts(workdir, &facts)
	collectFleetFacts(&facts)
	collectDbFacts(&facts)

	return facts
}

// collectRepoFacts counts repos and dirty working trees.
func collectRepoFacts(workdir string, facts *healthFacts) {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return
	}
	if _, lookErr := exec.LookPath("git"); lookErr != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		repoPath := filepath.Join(workdir, e.Name())
		if _, statErr := os.Stat(filepath.Join(repoPath, ".git")); statErr != nil {
			continue
		}
		facts.repoTotal++
		cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		if runErr := cmd.Run(); runErr != nil {
			continue
		}
		if strings.TrimSpace(stdout.String()) != "" {
			facts.repoDirty++
			if len(facts.dirtyNames) < 5 {
				facts.dirtyNames = append(facts.dirtyNames, e.Name())
			}
		}
	}
}

// collectPipelineFacts counts repos whose latest run is red (gh only).
func collectPipelineFacts(workdir string, facts *healthFacts) {
	if _, lookErr := exec.LookPath("gh"); lookErr != nil {
		facts.pipelineNote = "gh not installed; pipeline check skipped"
		return
	}
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		facts.pipelineNote = "workdir unreadable; pipeline check skipped"
		return
	}
	probed := 0
	for _, e := range entries {
		if !e.IsDir() || probed >= 50 {
			continue
		}
		repoPath := filepath.Join(workdir, e.Name())
		if _, statErr := os.Stat(filepath.Join(repoPath, ".git")); statErr != nil {
			continue
		}
		slug, hasSlug := healthRepoSlug(repoPath)
		if !hasSlug {
			continue
		}
		probed++
		cmd := exec.Command("gh", "run", "list", "-R", slug, "--limit", "1", "--json", "conclusion")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		if runErr := cmd.Run(); runErr != nil {
			continue
		}
		if strings.Contains(stdout.String(), `"failure"`) {
			facts.pipelineRed++
			if len(facts.redNames) < 5 {
				facts.redNames = append(facts.redNames, e.Name())
			}
		}
	}
}

// healthRepoSlug resolves owner/repo from the origin remote.
func healthRepoSlug(dir string) (string, bool) {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		return "", false
	}
	raw := strings.TrimSpace(stdout.String())
	lower := strings.ToLower(raw)
	var path string
	switch {
	case strings.HasPrefix(lower, "https://github.com/"):
		path = raw[len("https://github.com/"):]
	case strings.HasPrefix(lower, "git@github.com:"):
		path = raw[len("git@github.com:"):]
	default:
		return "", false
	}
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0] + "/" + parts[1], true
}

// collectFleetFacts counts registered fleet nodes.
func collectFleetFacts(facts *healthFacts) {
	dbConn, openErr := store.OpenDefault()
	if openErr != nil {
		facts.fleetNote = "node registry unreachable"
		return
	}
	defer dbConn.Close()
	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		facts.fleetNote = "node registry unreachable"
		return
	}
	facts.fleetNodes = len(res.Data)
	facts.fleetNote = "run `gitmap nodes version` for skew"
}

// collectDbFacts flags database files untouched for 30+ days.
func collectDbFacts(facts *healthFacts) {
	binDir := store.BinaryDataDir()
	cutoff := 30 * 24 * time.Hour
	for _, name := range []string{"gitmap.db", "installation.db", "pipeline.db"} {
		full := filepath.Join(binDir, name)
		info, statErr := os.Stat(full)
		if statErr != nil {
			continue
		}
		if time.Since(info.ModTime()) > cutoff {
			facts.staleDBs++
			facts.staleNames = append(facts.staleNames, name)
		}
	}
}
