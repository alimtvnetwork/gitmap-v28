package cmdagy

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
)

func resolveRepoSlugFromDir(repoDir string) string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return ""
	}

	return parseRepoSlug(strings.TrimSpace(string(out)))
}

func parseRepoSlug(rawURL string) string {
	clean := strings.TrimSuffix(rawURL, ".git")
	clean = strings.TrimPrefix(clean, "git@github.com:")
	clean = strings.TrimPrefix(clean, "https://github.com/")

	return strings.TrimSpace(clean)
}

func scanSingleProjectCandidate(repo localRepo, isDetailed bool) (ProjectFixCandidate, bool) {
	slug := resolveRepoSlugFromDir(repo.Path)
	if len(slug) == 0 {
		slug = repo.Name
	}

	payload, report, hasFail := cmdpipeline.FetchPipelineErrorReportWithMeta(slug, isDetailed)
	cand := ProjectFixCandidate{
		Name: repo.Name, Path: repo.Path, RepoSlug: slug,
		HasFailures: hasFail, RunID: payload.RunId, SHA: payload.Sha, ErrorReport: report,
	}

	return cand, hasFail
}

func probeCandidateWorker(r localRepo, isDetailed bool, mu *sync.Mutex, failing *[]ProjectFixCandidate) {
	cand, isFail := scanSingleProjectCandidate(r, isDetailed)
	if !isFail {
		return
	}
	mu.Lock()
	*failing = append(*failing, cand)
	mu.Unlock()
}

func launchCandidateScanWorkers(repos []localRepo, isDetailed bool, mu *sync.Mutex, failing *[]ProjectFixCandidate) {
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(r localRepo) {
			defer wg.Done()
			probeCandidateWorker(r, isDetailed, mu, failing)
		}(repo)
	}
	wg.Wait()
}

// ScanFailingProjectCandidates discovers local repos and probes pipeline failures concurrently.
func ScanFailingProjectCandidates(isDetailed bool) []ProjectFixCandidate {
	var (
		mu      sync.Mutex
		failing []ProjectFixCandidate
	)
	launchCandidateScanWorkers(discoverLocalRepos(), isDetailed, &mu, &failing)

	return failing
}
