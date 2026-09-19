// Package render exposes shared output renderers used by multiple
// gitmap commands. repotermblock.go implements the "human-readable
// per-repo summary" emitted by --output terminal across the scan,
// clone-from, clone-next, and probe commands.
//
// One renderer in one place keeps the four commands byte-identical
// in their per-repo block — users learn the format once and rely on
// it being grep-able regardless of which command produced it.
//
// Format (intentionally fixed-width labels for column alignment):
//
//	N. <repo-name>
//	   branch:    <branch> (<source>)
//	   transport: <ssh|https|other>
//	   https:     <httpsURL>
//	   ssh:       <sshURL>
//	   from:      <originalURL>
//	   to:        <targetURL>
//	   command:   <cloneCommand>
//
// Any field that is empty is rendered as the literal string "(unknown)"
// so the block always has the same shape — readers don't have to
// special-case missing lines, and a diff between two runs of the same
// command stays line-by-line aligned.
package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RepoTermBlock is the input to RenderRepoTermBlock. All fields are
// strings so producers can pass through whatever they have without
// import cycles into model/probe/clonefrom packages.
type RepoTermBlock struct {
	// Index is the 1-based row number printed before the name.
	Index int
	// Name is the repo's short name (basename or RepoName).
	Name string
	// PrimaryAlias is the primary alias displayed in brackets.
	PrimaryAlias string
	// Aliases lists all registered aliases for multi-alias tree view.
	Aliases []string
	// Branch is the detected branch ("main", "develop", …).
	Branch string
	// BranchSource describes how the branch was chosen: "HEAD",
	// "config", "default", "manifest", "detached", "unknown".
	BranchSource string
	// Transport is the discovered origin bucket: "ssh", "https", or "other".
	Transport string
	// HTTPSUrl is the HTTPS clone URL when known.
	HTTPSUrl string
	// SSHUrl is the SSH clone URL when known.
	SSHUrl string
	// OriginalURL is the URL as discovered (HTTPS preferred,
	// SSH fallback) or as supplied by the user.
	OriginalURL string
	// TargetURL is the URL that will actually be passed to git
	// clone — may equal OriginalURL when no rewrite happens.
	TargetURL string
	// CloneCommand is the full `git clone …` invocation.
	CloneCommand string
}

// fieldUnknown is the placeholder used for empty fields. Exported as
// a package-level const so tests can pin the spelling without
// importing constants.
const fieldUnknown = "(unknown)"

// RenderRepoTermBlock writes one block to w. Returns the first write
// error so callers can surface broken pipes / closed stderr instead
// of silently dropping output.
//
// The function is pure with respect to its inputs: same RepoTermBlock
// in → same bytes out, regardless of TTY / color settings. Color is
// deliberately omitted here so the same renderer works for both
// interactive terminals and CI logs that strip ANSI sequences.
func RenderRepoTermBlock(w io.Writer, b RepoTermBlock) error {
	header := formatBlockHeader(b)
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}

	body := buildBlockBody(b)
	_, err := io.WriteString(w, body)

	return err
}

func formatBlockHeader(b RepoTermBlock) string {
	name := fallback(b.Name)
	primary := resolveBlockPrimaryAlias(b.PrimaryAlias, b.Aliases)
	if len(primary) > 0 {
		return fmt.Sprintf("  %d. %s [%s]\n", b.Index, name, primary)
	}

	return fmt.Sprintf("  %d. %s\n", b.Index, name)
}

func resolveBlockPrimaryAlias(primary string, aliases []string) string {
	if len(strings.TrimSpace(primary)) > 0 {
		return strings.TrimSpace(primary)
	}
	if len(aliases) > 0 {
		return strings.TrimSpace(aliases[0])
	}

	return ""
}

// RenderRepoTermBlocks renders a slice in order. Stops on first
// write error so a broken pipe doesn't cause us to keep formatting
// blocks the reader will never see.
func RenderRepoTermBlocks(w io.Writer, blocks []RepoTermBlock) error {
	for _, b := range blocks {
		if err := RenderRepoTermBlock(w, b); err != nil {
			return err
		}
	}

	return nil
}

// buildBlockBody is split out so the formatting is testable as a
// pure function without needing an io.Writer.
func buildBlockBody(b RepoTermBlock) string {
	var sb strings.Builder
	appendBlockAliasesTree(&sb, b)
	fmt.Fprintf(&sb, "     branch:    %s\n", formatBranch(b.Branch, b.BranchSource))
	fmt.Fprintf(&sb, "     transport: %s\n", fallback(blockTransport(b)))
	fmt.Fprintf(&sb, "     https:     %s\n", fallback(blockHTTPSUrl(b)))
	fmt.Fprintf(&sb, "     ssh:       %s\n", fallback(blockSSHUrl(b)))
	fmt.Fprintf(&sb, "     from:      %s\n", fallback(b.OriginalURL))
	fmt.Fprintf(&sb, "     to:        %s\n", fallback(b.TargetURL))
	fmt.Fprintf(&sb, "     command:   %s\n", fallback(b.CloneCommand))

	return sb.String()
}

func appendBlockAliasesTree(sb *strings.Builder, b RepoTermBlock) {
	if len(b.Aliases) <= 1 {
		return
	}
	primary := resolveBlockPrimaryAlias(b.PrimaryAlias, b.Aliases)
	sb.WriteString("     aliases:\n")
	for i, alias := range b.Aliases {
		formatOneAliasBranch(sb, alias, primary, i == len(b.Aliases)-1)
	}
}

func formatOneAliasBranch(sb *strings.Builder, alias, primary string, isLast bool) {
	prefix := "├──"
	if isLast {
		prefix = "└──"
	}
	label := "secondary"
	if alias == primary {
		label = "primary"
	}
	fmt.Fprintf(sb, "       %s %s (%s)\n", prefix, alias, label)
}

func blockTransport(b RepoTermBlock) string {
	if len(strings.TrimSpace(b.Transport)) > 0 {
		return b.Transport
	}

	return inferTransport(b.OriginalURL)
}

func blockHTTPSUrl(b RepoTermBlock) string {
	if len(strings.TrimSpace(b.HTTPSUrl)) > 0 {
		return b.HTTPSUrl
	}

	if inferTransport(b.OriginalURL) == constants.ScanTransportHTTPS {
		return b.OriginalURL
	}

	return ""
}

func blockSSHUrl(b RepoTermBlock) string {
	if len(strings.TrimSpace(b.SSHUrl)) > 0 {
		return b.SSHUrl
	}

	if inferTransport(b.OriginalURL) == constants.ScanTransportSSH {
		return b.OriginalURL
	}

	return ""
}

func inferTransport(url string) string {
	trimmed := strings.TrimSpace(strings.ToLower(url))
	if strings.HasPrefix(trimmed, constants.PrefixSSH) || strings.HasPrefix(trimmed, constants.CommitInUrlPrefixSsh) {
		return constants.ScanTransportSSH
	}

	if strings.HasPrefix(trimmed, constants.PrefixHTTPS) {
		return constants.ScanTransportHTTPS
	}

	return constants.ScanTransportOther
}

// formatBranch renders "<branch> (<source>)" or just "<branch>" when
// no source is known. Empty branch falls back to "(unknown)".
func formatBranch(branch, source string) string {
	branch = strings.TrimSpace(branch)
	source = strings.TrimSpace(source)
	if len(branch) == 0 {
		return fieldUnknown
	}

	if len(source) == 0 {
		return branch
	}

	return fmt.Sprintf("%s (%s)", branch, source)
}

// fallback returns "(unknown)" for empty/whitespace-only input.
func fallback(s string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return fieldUnknown
	}

	return s
}
