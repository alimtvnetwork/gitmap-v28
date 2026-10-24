package cmdhealth

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	cmdconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth/constants"
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
	if len(facts.dirtyNames) >= cmdconstants.MaxDirtyReposDisplay {
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
	res := gitutil.RunGitInDir(repoPath, constants.GitStatus, constants.GitPorcelainFlag)
	if res.IsFailure() {
		facts.diagErrors = append(facts.diagErrors, res.Err)

		return
	}

	if strings.TrimSpace(string(res.Data)) == "" {
		return
	}

	recordDirtyRepo(name, facts)
}

func isRepoDir(repoPath string) bool {
	_, statErr := os.Stat(filepath.Join(repoPath, constants.DotGit))

	return statErr == nil
}

func inspectRepoEntry(workdir string, e os.DirEntry, facts *healthFacts) {
	repoPath := filepath.Join(workdir, e.Name())
	if !e.IsDir() || !isRepoDir(repoPath) {
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
	if len(facts.redNames) >= cmdconstants.MaxRedPipelinesDisplay {
		return
	}

	facts.redNames = append(facts.redNames, name)
}

func recordPipelineRed(name string, facts *healthFacts) {
	facts.pipelineRed++
	appendRedName(facts, name)
}

func probePipelineSlug(slug, name string, facts *healthFacts) {
	cmd := exec.Command(cmdconstants.CmdGh, cmdconstants.ArgGhRun, cmdconstants.ArgGhList,
		cmdconstants.FlagGhRepo, slug, cmdconstants.FlagGhLimit, cmdconstants.ValGhLimitOne,
		cmdconstants.FlagGhJSON, cmdconstants.ValGhConclusion)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		facts.diagErrors = append(facts.diagErrors, runErr)

		return
	}

	if strings.Contains(stdout.String(), cmdconstants.JsonConclusionFailure) {
		recordPipelineRed(name, facts)
	}
}

func inspectPipelineEntry(workdir string, e os.DirEntry, facts *healthFacts) bool {
	repoPath := filepath.Join(workdir, e.Name())
	if !e.IsDir() || !isRepoDir(repoPath) {
		return false
	}

	slug, hasSlug := healthRepoSlug(repoPath)
	if !hasSlug {
		return false
	}

	probePipelineSlug(slug, e.Name(), facts)

	return true
}

func scanPipelineEntries(workdir string, entries []os.DirEntry, facts *healthFacts) {
	probed := 0
	for _, e := range entries {
		if probed >= cmdconstants.MaxPipelineProbedRepos {
			break
		}

		if hasProbed := inspectPipelineEntry(workdir, e, facts); hasProbed {
			probed++
		}
	}
}

// collectPipelineFacts counts repos whose latest run is red (gh only).
func collectPipelineFacts(workdir string, facts *healthFacts) {
	if _, lookErr := exec.LookPath(cmdconstants.CmdGh); lookErr != nil {
		facts.pipelineNote = cmdconstants.NoteGhNotInstalled

		return
	}

	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		facts.pipelineNote = cmdconstants.NoteWorkdirUnreadable

		return
	}

	scanPipelineEntries(workdir, entries, facts)
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

	parts := strings.Split(path, cmdconstants.PathSeparator)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0] + cmdconstants.PathSeparator + parts[1], true
}

func matchGitHubPrefix(lower, raw string) string {
	switch {
	case strings.HasPrefix(lower, string(constants.GitHubHTTPSPrefix)):
		return raw[len(constants.GitHubHTTPSPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubHTTPPrefix)):
		return raw[len(constants.GitHubHTTPPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubSSHPrefix)):
		return raw[len(constants.GitHubSSHPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubSSHSchemePrefix)):
		return raw[len(constants.GitHubSSHSchemePrefix):]
	default:
		return ""
	}
}

func extractHealthSlugPath(raw string) string {
	path := matchGitHubPrefix(strings.ToLower(raw), raw)
	if path == "" {
		return ""
	}

	return strings.TrimSuffix(path, constants.GitDotGitSuffix)
}

func inspectFleetDb(dbConn *store.DB, facts *healthFacts) {
	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		facts.fleetNote = cmdconstants.NoteNodeRegistryUnreachable

		return
	}

	facts.fleetNodes = len(res.Data)
	facts.fleetNote = cmdconstants.NoteFleetSkewHint
}

// collectFleetFacts counts registered fleet nodes.
func collectFleetFacts(facts *healthFacts) {
	dbConn, openErr := store.OpenDefault()
	if openErr != nil {
		facts.fleetNote = cmdconstants.NoteNodeRegistryUnreachable

		return
	}
	defer dbConn.Close()

	inspectFleetDb(dbConn, facts)
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
	cutoff := time.Duration(cmdconstants.StaleDbDaysCutoff*cmdconstants.HoursPerDay) * time.Hour
	dbNames := []string{
		string(constants.DBFileMaster),
		string(constants.DBFileInstallation),
		string(constants.DBFilePipeline),
	}
	for _, name := range dbNames {
		checkStaleDb(binDir, name, cutoff, facts)
	}
}
