package cmd

// probereport.go — output, JSON shaping, and persistence helpers for
// `gitmap probe`. Split out of probe.go to honor the 200-line per-file
// budget. The dispatcher (runProbe) and the worker pool live in
// probe.go; everything here is single-repo concerns that the workers
// call into. tallyProbe is the only function that mutates shared
// state — its caller (probeWorker in probe.go) holds a mutex around
// each call, so this file does not import sync.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/probe"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// probeJSONEntry is a single repo-level result emitted under `--json`.
// Embeds the result + repo identity so a CI consumer can join on either.
type probeJSONEntry struct {
	RepoID         int64  `json:"repoId"`
	Slug           string `json:"slug"`
	AbsolutePath   string `json:"absolutePath"`
	NextVersionTag string `json:"nextVersionTag"`
	NextVersionNum int64  `json:"nextVersionNum"`
	Method         string `json:"method"`
	IsAvailable    bool   `json:"isAvailable"`
	Error          string `json:"error,omitempty"`
}

// executeOneProbe runs a single probe and persists it, mirroring the
// missing-URL handling that the sequential loop used. depth is forwarded
// to the shallow-clone fallback (probe.RunOneWithDepth).
func executeOneProbe(db *store.DB, repo model.ScanRecord, depth int) probe.Result {
	url := pickProbeURL(repo)
	if url == "" {
		result := probe.Result{Method: constants.ProbeMethodNone, Error: fmt.Sprintf(constants.ErrProbeMissingURL, repo.Slug)}
		recordProbeResult(db, repo, result)

		return result
	}
	result := probe.RunOneWithDepth(url, depth)
	recordProbeResult(db, repo, result)

	return result
}

// makeProbeEntry converts a probe.Result + repo into a JSON-friendly row.
func makeProbeEntry(repo model.ScanRecord, r probe.Result) probeJSONEntry {
	return probeJSONEntry{
		RepoID:         repo.ID,
		Slug:           repo.Slug,
		AbsolutePath:   repo.AbsolutePath,
		NextVersionTag: r.NextVersionTag,
		NextVersionNum: r.NextVersionNum,
		Method:         r.Method,
		IsAvailable:    r.IsAvailable,
		Error:          r.Error,
	}
}

// emitProbeJSON dumps the collected entries as stable JSON to stdout.
func emitProbeJSON(entries []probeJSONEntry) {
	if err := encodeProbeJSON(os.Stdout, entries); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// pickProbeURL honors the repo's identified Transport so an SSH-origin
// repo is probed over SSH (avoids interactive browser-auth prompts on
// private remotes). Falls back to the other transport only when the
// preferred URL is empty; "other"/unknown still prefers HTTPS for CI
// auth-friction reasons.
func pickProbeURL(r model.ScanRecord) string {
	if r.Transport == constants.ScanTransportSSH && r.SSHUrl != "" {
		return r.SSHUrl
	}
	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}

	return r.SSHUrl
}

// recordProbeResult persists the probe row, logging-but-not-exiting on error.
func recordProbeResult(db *store.DB, repo model.ScanRecord, result probe.Result) {
	if err := db.RecordVersionProbe(result.AsModel(repo.ID)); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// tallyProbe updates the running counters and (unless jsonOut) prints the
// per-repo summary line. Caller is responsible for serializing access to
// the counters; with the worker pool that's `counterMu` in runProbePool.
func tallyProbe(repo model.ScanRecord, r probe.Result, ok, none, fail int, jsonOut bool) (int, int, int) {
	if r.Error != "" {
		if !jsonOut {
			fmt.Printf(constants.MsgProbeFailFmt, repo.Slug, r.Error)
		}
		return ok, none, fail + 1
	}
	if r.IsAvailable {
		if !jsonOut {
			fmt.Printf(constants.MsgProbeOkFmt, repo.Slug, r.NextVersionTag, r.Method)
		}
		return ok + 1, none, fail
	}
	if !jsonOut {
		fmt.Printf(constants.MsgProbeNoneFmt, repo.Slug, r.Method)
	}

	return ok, none + 1, fail
}

// emitProbeTermBlock prints the standardized RepoTermBlock for one
// probed repo. Called by the worker (under counterMu) when
// `--output terminal` is active so the per-probe summary matches the
// shape used by scan, clone-from, and clone-next.
//
// The block surfaces the probe outcome in CloneCommand: a successful
// probe shows the would-be `git clone -b <nextTag> <url>` invocation
// the user can copy/paste; a failed probe shows the trimmed error.
func emitProbeTermBlock(idx int, repo model.ScanRecord, r probe.Result) {
	url := pickProbeURL(repo)
	cmd := probeCloneCommandFor(url, r)
	block := render.RepoTermBlock{
		Index:        idx,
		Name:         repo.Slug,
		Branch:       repo.Branch,
		BranchSource: probeBranchSource(repo),
		OriginalURL:  url,
		TargetURL:    url,
		CloneCommand: cmd,
	}
	_ = render.RenderRepoTermBlock(os.Stdout, block)
}

// probeCloneCommandFor renders the next-version git clone command
// when the probe surfaced a tag, falls back to the bare clone when
// no new tag exists, and surfaces the probe error verbatim when the
// probe failed (so the user immediately sees why no command is
// suggested).
func probeCloneCommandFor(url string, r probe.Result) string {
	if r.Error != "" {
		return fmt.Sprintf("(probe failed: %s)", r.Error)
	}
	if !r.IsAvailable || r.NextVersionTag == "" {
		return fmt.Sprintf("%s %s %s", constants.GitBin, constants.GitClone, url)
	}

	return fmt.Sprintf("%s %s -b %s %s",
		constants.GitBin, constants.GitClone, r.NextVersionTag, url)
}

// probeBranchSource returns "manifest" when the repo has a branch
// recorded in the scan DB, "" otherwise (the renderer drops the
// parenthesized source segment when source is empty).
func probeBranchSource(repo model.ScanRecord) string {
	if repo.Branch == "" {
		return ""
	}

	return "manifest"
}
