package clonepick

// render.go: dry-run preview rendering. Mirrors clonenow.Render in
// shape (header + body + trailing newline) so users get a familiar
// look-and-feel across the clone-* family.

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Render writes the dry-run preview for a single Plan. Stable output
// (no map iteration, no time/now, no random ordering) so the goldens
// don't churn between runs.
func Render(w io.Writer, plan Plan) error {
	branch := plan.Branch
	if len(branch) == 0 {
		branch = "(default)"
	}
	sparseMode := "cone"
	if !plan.Cone {
		sparseMode = "non-cone"
	}

	header := fmt.Sprintf(constants.MsgClonePickDryHeader,
		plan.RepoUrl, plan.DestDir, plan.Mode, branch,
		plan.Depth, sparseMode, len(plan.Paths))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}

	for _, p := range plan.Paths {
		if _, err := fmt.Fprintf(w, "  + %s\n", p); err != nil {
			return err
		}
	}

	return renderCommands(w, plan)
}

// renderCommands prints the exact `git` invocations that would run.
// Useful for users who want to copy/paste the steps into their own
// scripts rather than calling gitmap.
func renderCommands(w io.Writer, plan Plan) error {
	cloneCmd := buildCloneCommandPreview(plan)
	sparseCmd := buildSparsePreview(plan)
	checkoutCmd := "git -C " + plan.DestDir + " checkout"

	body := "\nplanned commands:\n  $ " + cloneCmd +
		"\n  $ " + sparseCmd +
		"\n  $ " + checkoutCmd + "\n"
	if !plan.KeepGit {
		body += "  $ rm -rf " + plan.DestDir + "/.git\n"
	}
	_, err := io.WriteString(w, body)

	return err
}

// buildCloneCommandPreview assembles the printable `git clone` line.
// Kept tiny so the function-length rule is comfortably satisfied.
func buildCloneCommandPreview(plan Plan) string {
	parts := []string{"git", "clone", "--filter=blob:none", "--no-checkout"}
	if len(plan.Branch) > 0 {
		parts = append(parts, "--branch", plan.Branch)
	}
	if plan.Depth > 0 {
		parts = append(parts, "--depth", fmt.Sprintf("%d", plan.Depth))
	}
	parts = append(parts, plan.RepoUrl, plan.DestDir)

	return strings.Join(parts, " ")
}

// buildSparsePreview assembles the `git sparse-checkout set` preview
// in a single line so users can paste it verbatim.
func buildSparsePreview(plan Plan) string {
	mode := "--cone"
	if !plan.Cone {
		mode = "--no-cone"
	}
	parts := []string{"git", "-C", plan.DestDir, "sparse-checkout", "set", mode}
	parts = append(parts, plan.Paths...)

	return strings.Join(parts, " ")
}
