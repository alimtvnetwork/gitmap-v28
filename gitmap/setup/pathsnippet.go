package setup

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RenderPathSnippet returns the canonical marker-block PATH snippet for
// the requested shell, with `dir` injected as the resolved deploy
// directory and `manager` shown in the header line.
//
// Output is byte-identical to what run.sh and gitmap/scripts/install.sh
// would produce — those scripts shell out to `gitmap setup
// print-path-snippet` and pipe the bytes into the user's rc file.
//
// Spec: spec/04-generic-cli/21-post-install-shell-activation/02-snippets.md
func RenderPathSnippet(shell, dir, manager string) (string, error) {
	if len(dir) == 0 {
		return "", fmt.Errorf(constants.ErrPathSnippetDirRequired)
	}
	tpl, err := snippetTemplate(shell)
	if err != nil {
		return "", err
	}
	if len(manager) == 0 {
		manager = "gitmap setup"
	}

	return fmt.Sprintf(tpl, manager, dir), nil
}

// snippetTemplate maps a shell identifier to its body template.
func snippetTemplate(shell string) (string, error) {
	switch shell {
	case constants.PathSnippetShellBash:
		return constants.PathSnippetBashFmt, nil
	case constants.PathSnippetShellZsh:
		return constants.PathSnippetZshFmt, nil
	case constants.PathSnippetShellFish:
		return constants.PathSnippetFishFmt, nil
	case constants.PathSnippetShellPwsh:
		return constants.PathSnippetPwshFmt, nil
	}

	return "", fmt.Errorf(constants.ErrPathSnippetUnknownShell, shell)
}

// MarkerOpenFor returns the rendered open-marker line for the manager
// string. Used by writers that need to detect/rewrite an existing block.
func MarkerOpenFor(manager string) string {
	if len(manager) == 0 {
		manager = "gitmap setup"
	}

	return fmt.Sprintf(constants.PathSnippetMarkerOpenFmt, manager)
}

// MarkerClose returns the literal close-marker line.
func MarkerClose() string {
	return constants.PathSnippetMarkerClose
}
