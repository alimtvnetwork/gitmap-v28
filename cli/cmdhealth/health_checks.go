package cmdhealth

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
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
	diagErrors   []error
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

func appendDirtyName(facts *healthFacts, name string) {
	if len(facts.dirtyNames) >= 5 {
		return
	}

	facts.dirtyNames = append(facts.dirtyNames, name)
}

func recordDirtyRepo(name string, facts *healthFacts) {
	facts.repoDirty++
	appendDirtyName(facts, name)
}

func probeRepoDirty(repoPath, name string, facts *healthFacts) {
	facts.repoTotal++
	out, runErr := gitutil.RunGitInDir(repoPath, constants.GitStatus, constants.GitPorcelainFlag)
	if runErr != nil {
		facts.diagErrors = append(facts.diagErrors, runErr)

		return
	}

	if strings.TrimSpace(string(out)) == "" {
		return
	}

	recordDirtyRepo(name, facts)
}

func inspectRepoEntry(workdir string, e os.DirEntry, facts *healthFacts) {
	if !e.IsDir() {
		return
	}

	repoPath := filepath.Join(workdir, e.Name())
	if _, statErr := os.Stat(filepath.Join(repoPath, constants.DotGit)); statErr != nil {
		return
	}

	probeRepoDirty(repoPath, e.Name(), facts)
}

// collectRepoFacts counts repos and dirty working trees.
func collectRepoFacts(workdir string, facts *healthFacts) {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return
	}

	if _, lookErr := exec.LookPath(constants.GitBin); lookErr != nil {
		return
	}

	for _, e := range entries {
		inspectRepoEntry(workdir, e, facts)
	}
}

func appendRedName(facts *healthFacts, name string) {
	if len(facts.redNames) >= 5 {
		return
	}

	facts.redNames = append(facts.redNames, name)
}

func recordPipelineRed(name string, facts *healthFacts) {
	facts.pipelineRed++
	appendRedName(facts, name)
}

func probePipelineSlug(slug, name string, facts *healthFacts) {
	cmd := exec.Command(constants.GHBin, constants.GHRun, constants.GHList, "-R", slug, "--limit", "1", "--json", "conclusion")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		facts.diagErrors = append(facts.diagErrors, runErr)

		return
	}

	if strings.Contains(stdout.String(), `"failure"`) {
		recordPipelineRed(name, facts)
	}
}

func inspectPipelineEntry(workdir string, e os.DirEntry, facts *healthFacts) bool {
	if !e.IsDir() {
		return false
	}

	repoPath := filepath.Join(workdir, e.Name())
	if _, statErr := os.Stat(filepath.Join(repoPath, constants.DotGit)); statErr != nil {
		return false
	}

	slug, hasSlug := healthRepoSlug(repoPath)
	if !hasSlug {
		return false
	}

	probePipelineSlug(slug, e.Name(), facts)

	return true
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
		if probed >= 50 {
			break
		}

		if hasProbed := inspectPipelineEntry(workdir, e, facts); hasProbed {
			probed++
		}
	}
}

// healthRepoSlug resolves owner/repo from the origin remote.
func healthRepoSlug(dir string) (string, bool) {
	origin := gitutil.GetRemoteOriginURL(dir)
	if origin == "" {
		return "", false
	}

	return parseHealthSlug(origin)
}

func parseHealthSlug(raw string) (string, bool) {
	path := extractHealthSlugPath(raw)
	if path == "" {
		return "", false
	}

	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0] + "/" + parts[1], true
}

func extractHealthSlugPath(raw string) string {
	lower := strings.ToLower(raw)
	var path string
	switch {
	case strings.HasPrefix(lower, string(constants.GitHubHTTPSPrefix)):
		path = raw[len(constants.GitHubHTTPSPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubHTTPPrefix)):
		path = raw[len(constants.GitHubHTTPPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubSSHPrefix)):
		path = raw[len(constants.GitHubSSHPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubSSHSchemePrefix)):
		path = raw[len(constants.GitHubSSHSchemePrefix):]
	default:
		return ""
	}

	return strings.TrimSuffix(path, constants.GitDotGitSuffix)
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

func checkStaleDb(binDir, name string, cutoff time.Duration, facts *healthFacts) {
	full := filepath.Join(binDir, name)
	info, statErr := os.Stat(full)
	if statErr != nil {
		return
	}

	if time.Since(info.ModTime()) <= cutoff {
		return
	}

	facts.staleDBs++
	facts.staleNames = append(facts.staleNames, name)
}

// collectDbFacts flags database files untouched for 30+ days.
func collectDbFacts(facts *healthFacts) {
	binDir := store.BinaryDataDir()
	cutoff := 30 * 24 * time.Hour
	for _, name := range []string{string(constants.DBFileMaster), string(constants.DBFileInstallation), string(constants.DBFilePipeline)} {
		checkStaleDb(binDir, name, cutoff, facts)
	}
}
