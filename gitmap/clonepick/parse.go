package clonepick

// parse.go: turns raw CLI flag values into a validated Plan.
//
// All validation happens here (one place) so cmd/clonepick.go can
// stay a thin flag-binding shim and Execute / Render can assume their
// input Plan is well-formed.

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// ParseArgs builds a Plan from the user-supplied URL, the
// comma-separated path list, and the parsed flag values.
//
// Returns an error (not os.Exit) so the cmd layer owns the exit
// policy and tests can assert on the error value directly.
func ParseArgs(rawURL, rawPaths string, flags Flags) (Plan, error) {
	plan := Plan{
		Name:    flags.Name,
		Mode:    flags.Mode,
		Branch:  flags.Branch,
		Depth:   flags.Depth,
		Cone:    flags.Cone,
		KeepGit: flags.KeepGit,
		DestDir: flags.Dest,
		UsedAsk: flags.Ask,
		DryRun:  flags.DryRun,
		Quiet:   flags.Quiet,
		Force:   flags.Force,
	}

	if err := validateMode(flags.Mode); err != nil {
		return plan, err
	}
	if err := validateDepth(flags.Depth); err != nil {
		return plan, err
	}

	url, err := resolveURL(rawURL, flags.Mode)
	if err != nil {
		return plan, err
	}
	plan.RepoUrl = url
	plan.RepoCanonicalId = gitutil.CanonicalRepoID(url)

	paths, err := normalisePaths(rawPaths, flags.Ask)
	if err != nil {
		return plan, err
	}
	plan.Paths = paths

	// Auto-disable cone mode when the path list contains glob chars
	// or file-shaped entries -- cone mode only matches at directory
	// granularity. Done after path validation so the heuristic only
	// runs over already-clean strings.
	if flags.Cone && hasNonConeShape(paths) {
		plan.Cone = false
	}

	return plan, nil
}

// Flags is the cmd-side flag bundle, decoupled from the flag.FlagSet
// so the parser is testable without spinning up a full CLI.
type Flags struct {
	Ask     bool
	Name    string
	Mode    string
	Branch  string
	Depth   int
	Cone    bool
	Dest    string
	KeepGit bool
	DryRun  bool
	Quiet   bool
	Force   bool
	// Replay is the raw value of --replay (numeric SelectionId or
	// Name). Empty means "no replay -- parse positional args
	// normally". Non-empty short-circuits ParseArgs in the cmd
	// dispatcher: the Plan is loaded from the DB instead.
	Replay string
}

// DefaultFlags returns the flag bundle with the spec'd defaults.
// Centralized so cmd/clonepick.go and tests stay in sync.
func DefaultFlags() Flags {
	return Flags{
		Ask:     false,
		Name:    "",
		Mode:    constants.ClonePickModeHTTPS,
		Branch:  "",
		Depth:   1,
		Cone:    true,
		Dest:    ".",
		KeepGit: true,
		DryRun:  false,
		Quiet:   false,
		Force:   false,
	}
}

// validateMode rejects anything that isn't "https" or "ssh". Done
// before URL resolution so a typo can't silently pick the wrong
// transport when the input was shorthand.
func validateMode(mode string) error {
	if mode != constants.ClonePickModeHTTPS && mode != constants.ClonePickModeSSH {
		return fmt.Errorf(constants.ErrClonePickBadMode, mode)
	}

	return nil
}

// validateDepth blocks negative depths (--depth=-1 is meaningless to
// git and would silently fall through to a confusing remote error).
func validateDepth(depth int) error {
	if depth < 0 {
		return fmt.Errorf(constants.ErrClonePickBadDepth, depth)
	}

	return nil
}

// resolveURL expands `owner/repo` shorthand using the chosen mode and
// returns full URLs verbatim (after a TrimSpace). Anything that looks
// like an existing scheme is passed through.
func resolveURL(raw, mode string) (string, error) {
	s := strings.TrimSpace(raw)
	if len(s) == 0 {
		return "", fmt.Errorf("%s", constants.MsgClonePickMissingURL)
	}
	if strings.Contains(s, "://") || strings.HasPrefix(s, "git@") {
		return s, nil
	}
	// Shorthand: owner/repo OR host/owner/repo.
	parts := strings.Split(s, "/")
	switch len(parts) {
	case 2:
		return shorthandToURL("github.com", parts[0], parts[1], mode), nil
	case 3:
		return shorthandToURL(parts[0], parts[1], parts[2], mode), nil
	default:
		// Anything else, hand it to git verbatim. Git will reject
		// nonsense with its own (more accurate) error message.
		return s, nil
	}
}

// shorthandToURL builds the canonical https or ssh URL for the host
// triple. Kept tiny so it stays under the 15-line function rule.
func shorthandToURL(host, owner, repo, mode string) string {
	repo = strings.TrimSuffix(repo, ".git")
	if mode == constants.ClonePickModeSSH {
		return fmt.Sprintf("git@%s:%s/%s.git", host, owner, repo)
	}

	return fmt.Sprintf("https://%s/%s/%s.git", host, owner, repo)
}

// normalisePaths splits, validates, deduplicates, and sorts the
// comma-separated path list. askMode==true permits an empty list:
// the picker will fill it in interactively, so requiring `<paths>`
// up front would be a usability bug.
func normalisePaths(raw string, askMode bool) ([]string, error) {
	if len(strings.TrimSpace(raw)) == 0 {
		if askMode {
			return nil, nil
		}

		return nil, fmt.Errorf("%s", constants.MsgClonePickMissingPaths)
	}
	seen := make(map[string]struct{})
	out := make([]string, 0, 4)
	for _, p := range strings.Split(raw, ",") {
		clean, err := cleanPath(p)
		if err != nil {
			return nil, err
		}
		if _, dup := seen[clean]; dup {
			continue
		}
		seen[clean] = struct{}{}
		out = append(out, clean)
	}
	sort.Strings(out)

	return out, nil
}

// cleanPath enforces the spec's path rules and strips the cosmetic
// `./` prefix and trailing `/` so two equivalent forms ("docs" and
// "./docs/") collapse into one canonical entry.
func cleanPath(raw string) (string, error) {
	p := strings.TrimSpace(raw)
	if len(p) == 0 {
		return "", fmt.Errorf("%s", constants.MsgClonePickPathEmpty)
	}
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") {
		return "", fmt.Errorf(constants.MsgClonePickPathAbsolute, p)
	}
	p = strings.TrimPrefix(p, "./")
	p = strings.TrimSuffix(p, "/")
	if strings.Contains(p, "..") {
		return "", fmt.Errorf(constants.MsgClonePickPathTraversal, p)
	}
	if len(p) > constants.ClonePickPathMaxBytes {
		return "", fmt.Errorf(constants.MsgClonePickPathTooLong, p)
	}

	return p, nil
}

// hasNonConeShape reports whether any path looks like a glob or a
// file (extension after the last /). Used to auto-flip --cone off so
// the user doesn't have to remember the cone-mode constraint.
func hasNonConeShape(paths []string) bool {
	for _, p := range paths {
		if strings.ContainsAny(p, "*?[") {
			return true
		}
		base := p
		if i := strings.LastIndex(p, "/"); i >= 0 {
			base = p[i+1:]
		}
		if strings.Contains(base, ".") {
			return true
		}
	}

	return false
}
