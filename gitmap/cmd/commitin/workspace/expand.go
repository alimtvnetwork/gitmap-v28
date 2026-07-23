package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ResolvedInput is one entry in the post-expansion input list. Order
// in the returned slice IS the walk order (spec §3.4 determinism).
type ResolvedInput struct {
	OrderIndex int    // 1-based position; matches CommitInTempInputFormat
	Original   string // verbatim user token OR sibling basename
	Kind       string // CommitInInputKind* literal
	AbsPath    string // empty for GitUrl until staged; abs path otherwise
	URL        string // populated only when Kind == GitUrl
	Version    int    // -1 = not a versioned sibling; 0 = plain base; >=1 = -vN
}

// ExpandInputs implements spec §3.1 stage 07. When `keyword` is
// empty, each `inputs` token is classified individually. When the
// keyword is set, siblings of `source` are discovered and (optionally)
// truncated to the last `tail` entries.
func ExpandInputs(source string, inputs []string, keyword string, tail int) ([]ResolvedInput, error) {
	if keyword != "" {
		return expandKeyword(source, keyword, tail)
	}
	return expandExplicit(inputs)
}

// expandExplicit classifies each user-supplied token. Order preserved.
func expandExplicit(inputs []string) ([]ResolvedInput, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf(constants.CommitInErrBadArgs, "no inputs to expand")
	}
	out := make([]ResolvedInput, 0, len(inputs))
	for i, tok := range inputs {
		out = append(out, classifyExplicitInput(i+1, tok))
	}
	return out, nil
}

// classifyExplicitInput chooses GitUrl vs LocalFolder for one token.
func classifyExplicitInput(orderIndex int, tok string) ResolvedInput {
	if isGitURL(tok) {
		return ResolvedInput{
			OrderIndex: orderIndex,
			Original:   tok,
			Kind:       constants.CommitInInputKindGitUrl,
			URL:        tok,
			Version:    -1,
		}
	}
	abs, _ := filepath.Abs(tok)
	return ResolvedInput{
		OrderIndex: orderIndex,
		Original:   tok,
		Kind:       constants.CommitInInputKindLocalFolder,
		AbsPath:    abs,
		Version:    -1,
	}
}

// expandKeyword implements `all` and `-N` discovery (spec §2.4).
func expandKeyword(source, keyword string, tail int) ([]ResolvedInput, error) {
	siblings, err := discoverSiblings(source)
	if err != nil {
		return nil, err
	}
	if len(siblings) == 0 {
		return nil, fmt.Errorf(constants.CommitInErrBadArgs, fmt.Sprintf("keyword %q matched zero siblings of %s", keyword, source))
	}
	if keyword != constants.CommitInInputKeywordAll && tail > 0 && tail < len(siblings) {
		siblings = siblings[len(siblings)-tail:]
	}
	return reindex(siblings), nil
}

// reindex stamps a fresh 1-based OrderIndex onto each sibling so the
// caller's temp-folder layout matches CommitInTempInputFormat.
func reindex(in []ResolvedInput) []ResolvedInput {
	for i := range in {
		in[i].OrderIndex = i + 1
	}
	return in
}

// versionSuffix matches the canonical `-vN` suffix used everywhere in
// the project (clone-next, fix-repo, history-rewrite). N is captured
// as group 1.
var versionSuffix = regexp.MustCompile(`-v(\d+)$`)

// discoverSiblings reads source's parent directory and returns every
// versioned sibling sorted ascending by version (plain base = v0).
// `source` itself is excluded.
func discoverSiblings(source string) ([]ResolvedInput, error) {
	abs, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("absolutize source: %w", err)
	}
	parent := filepath.Dir(abs)
	base := stripVersionSuffix(filepath.Base(abs))
	entries, readErr := os.ReadDir(parent)
	if readErr != nil {
		return nil, fmt.Errorf("read parent %s: %w", parent, readErr)
	}
	return collectMatchingSiblings(parent, base, abs, entries), nil
}

// collectMatchingSiblings filters DirEntry slice down to versioned
// siblings, then sorts ascending. Pure — no filesystem calls.
func collectMatchingSiblings(parent, base, sourceAbs string, entries []os.DirEntry) []ResolvedInput {
	out := make([]ResolvedInput, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		full := filepath.Join(parent, e.Name())
		if full == sourceAbs {
			continue
		}
		if v, ok := matchSibling(base, e.Name()); ok {
			out = append(out, ResolvedInput{
				Original: e.Name(),
				Kind:     constants.CommitInInputKindVersionedSibling,
				AbsPath:  full,
				Version:  v,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out
}

// matchSibling returns (version, true) when `name` matches `base` or
// `base-vN`. Plain base resolves to v0.
func matchSibling(base, name string) (int, bool) {
	if name == base {
		return 0, true
	}
	if !strings.HasPrefix(name, base+"-v") {
		return 0, false
	}
	m := versionSuffix.FindStringSubmatch(name)
	if len(m) == 0 || strings.TrimSuffix(name, m[0]) != base {
		return 0, false
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, false
	}
	return n, true
}

// stripVersionSuffix removes a single trailing `-vN` from a basename.
// Public-named helper used by discoverSiblings AND the test suite.
func stripVersionSuffix(name string) string {
	m := versionSuffix.FindStringSubmatchIndex(name)
	if m == nil {
		return name
	}
	return name[:m[0]]
}
