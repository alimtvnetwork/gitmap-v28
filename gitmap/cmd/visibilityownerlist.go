// Package cmd — visibilityownerlist.go: enumerates every repo under a
// given owner/org via the host provider CLI.
//
// GitHub:  gh   repo list <owner> --limit <N> --json name
// GitLab:  glab repo list --group <owner> -P <N> -F json
//
// Output JSON is parsed with a permissive `[{"name":"…"}]` scanner —
// no full JSON unmarshal — so we tolerate extra fields without coupling
// to the provider's schema version.
//
// Pagination: the provider CLIs cap at --limit; when the returned slice
// length equals the cap we log a stderr WARNING per spec §plan step 26
// so users with >cap repos are not silently truncated.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan step 9.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// listOwnerRepos returns the bare repo names (no owner prefix) under
// the given owner. Surfaces provider CLI errors verbatim with Code Red
// context for the operator.
func listOwnerRepos(provider, owner string, verbose bool) ([]string, error) {
	args := ownerRepoListArgs(provider, owner)
	out, err := runProviderCLI(provider, args, verbose)
	if err != nil {
		return nil, fmt.Errorf("Error: provider repo list failed for %s/%s: %v (operation: %s %s, reason: %s, stderr: %s)",
			provider, owner, err, providerCLI(provider), strings.Join(args, " "), err.Error(), strings.TrimSpace(out))
	}

	names, err := parseOwnerRepoNames(out)
	if err != nil {
		return nil, fmt.Errorf("Error: cannot parse provider repo list for %s/%s: %v (operation: json decode, reason: %s)",
			provider, owner, err, err.Error())
	}

	if len(names) >= constants.OwnerRepoListLimit {
		fmt.Fprintf(os.Stderr, constants.WarnOwnerRepoListCapFmt,
			constants.OwnerRepoListLimit, owner)
	}

	return names, nil
}

// ownerRepoListArgs builds the argv for "list every repo under owner".
func ownerRepoListArgs(provider, owner string) []string {
	limit := fmt.Sprintf("%d", constants.OwnerRepoListLimit)
	if provider == constants.ProviderGitLab {
		return []string{"repo", "list", "--group", owner, "-P", limit, "-F", "json"}
	}

	return []string{"repo", "list", owner, "--limit", limit, "--json", "name"}
}

// parseOwnerRepoNames decodes the JSON array emitted by gh / glab.
// Both CLIs share a `[{"name":"<bare-repo-name>", ...}]` envelope, so
// one scanner covers both — extra fields are ignored.
func parseOwnerRepoNames(raw string) ([]string, error) {
	var rows []struct {
		Name string `json:"name"`
		Path string `json:"path"` // glab fallback (older versions)
	}
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(rows))
	for _, r := range rows {
		name := r.Name
		if len(name) == 0 {
			name = r.Path
		}
		if len(name) > 0 {
			out = append(out, name)
		}
	}

	return out, nil
}
