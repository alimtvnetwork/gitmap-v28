// Package probe implements the hybrid HEAD-then-clone version probe.
//
// The probe inspects a single repo's remote, looking for the highest
// semver-style tag (vN.N.N or N.N.N). Strategy order:
//
//  1. `git ls-remote --tags --sort=-v:refname <url>` — cheap, no working
//     copy required. Most servers return tags in a single round trip.
//  2. Fallback `git clone --depth 1 --filter=blob:none --no-checkout <url>`
//     into a temp dir, then `git tag --sort=-v:refname` — used when ls-remote
//     fails outright (auth quirks, smart-http rejection, etc.) or returns
//     zero tags despite the remote actually having some.
//
// The fallback is intentionally treeless and depth-1: we only need the
// refs database, not the worktree. The temp dir is removed before return.
package probe

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// Result is what RunOne returns. Mirrors model.VersionProbe but without the
// DB-only fields (ID, RepoID, ProbedAt are filled by the caller).
type Result struct {
	NextVersionTag string
	NextVersionNum int64
	Method         string
	IsAvailable    bool
	Error          string
}

// RunOne probes a single repo for the highest tag using the default
// shallow-clone depth (1). Kept for source compatibility with callers
// that don't expose a depth knob; new code should call RunOneWithDepth.
func RunOne(cloneURL string) Result {
	return RunOneWithDepth(cloneURL, constants.ProbeDefaultDepth)
}

// RunOneWithDepth probes a single repo for the highest tag, passing
// `depth` to the `git clone --depth N` shallow-clone fallback. cloneURL
// must be a usable HTTPS or SSH URL. Returns a Result with Error
// populated on failure (never returns a non-nil error — failures are
// recorded, not bubbled). depth<1 is coerced to 1 inside tryShallowClone.
func RunOneWithDepth(cloneURL string, depth int) Result {
	if strings.TrimSpace(cloneURL) == "" {
		return Result{Method: constants.ProbeMethodNone, Error: "empty clone url"}
	}

	if tag, ok := tryLsRemote(cloneURL); ok {
		return makeResult(tag, constants.ProbeMethodLsRemote)
	}

	if tag, err := tryShallowClone(cloneURL, depth); err == nil {
		return makeResult(tag, constants.ProbeMethodShallowClone)
	} else {
		return Result{
			Method: constants.ProbeMethodShallowClone,
			Error:  err.Error(),
		}
	}
}

// makeResult builds a Result for a successfully-probed tag.
func makeResult(tag, method string) Result {
	if tag == "" {
		return Result{Method: method, IsAvailable: false}
	}

	return Result{
		NextVersionTag: tag,
		NextVersionNum: parseSemverInt(tag),
		Method:         method,
		IsAvailable:    true,
	}
}

// AsModel converts a probe.Result into a model.VersionProbe ready for
// store.RecordVersionProbe.
func (r Result) AsModel(repoID int64) model.VersionProbe {
	return model.VersionProbe{
		RepoID:         repoID,
		NextVersionTag: r.NextVersionTag,
		NextVersionNum: r.NextVersionNum,
		Method:         r.Method,
		IsAvailable:    r.IsAvailable,
		Error:          r.Error,
	}
}

// tryLsRemote runs `git ls-remote --tags --sort=-v:refname <url>` and
// returns the first tag it parses. ok=false means no usable tag found.
func tryLsRemote(url string) (string, bool) {
	cmd := exec.Command("git", "ls-remote", "--tags", "--sort=-v:refname", url)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}

	return parseFirstTag(string(out)), true && parseFirstTag(string(out)) != ""
}

// parseFirstTag pulls the first refs/tags/<name> line from ls-remote output,
// stripping the dereference suffix `^{}` that annotated tags emit.
func parseFirstTag(out string) string {
	for _, line := range strings.Split(out, "\n") {
		idx := strings.Index(line, "refs/tags/")
		if idx < 0 {
			continue
		}
		tag := strings.TrimSpace(line[idx+len("refs/tags/"):])
		tag = strings.TrimSuffix(tag, "^{}")
		if tag != "" {
			return tag
		}
	}

	return ""
}

// parseSemverInt converts vMAJOR.MINOR.PATCH (or MAJOR.MINOR.PATCH) into a
// monotonically-increasing int64 suitable for ORDER BY: MAJOR*1e6 + MINOR*1e3 + PATCH.
// Anything non-numeric collapses to 0 — used only for sort keys, never for display.
func parseSemverInt(tag string) int64 {
	clean := strings.TrimPrefix(tag, "v")
	parts := strings.SplitN(clean, ".", 3)
	if len(parts) < 1 {
		return 0
	}

	var major, minor, patch int64
	_, _ = fmt.Sscanf(parts[0], "%d", &major)
	if len(parts) > 1 {
		_, _ = fmt.Sscanf(parts[1], "%d", &minor)
	}
	if len(parts) > 2 {
		// Trim any pre-release suffix (e.g. "1-rc1").
		head := parts[2]
		for i, c := range head {
			if c < '0' || c > '9' {
				head = head[:i]
				break
			}
		}
		_, _ = fmt.Sscanf(head, "%d", &patch)
	}

	return major*1_000_000 + minor*1_000 + patch
}
