package cmdinstall

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestProfileTreeConnectors(t *testing.T) {
	connLast := resolveTreeConnector(true)
	if !strings.Contains(connLast, constants.TreeCorner) {
		t.Errorf("expected corner connector for last child, got %s", connLast)
	}

	connBranch := resolveTreeConnector(false)
	if !strings.Contains(connBranch, constants.TreeBranch) {
		t.Errorf("expected branch connector for non-last child, got %s", connBranch)
	}
}

func TestProfileTreeNodeStatus(t *testing.T) {
	installed := map[string]string{
		"vscode": "1.85.0",
	}

	dotInstalled, verInstalled := resolveTreeNodeStatus("vscode", installed)
	if dotInstalled != installedDot || verInstalled != "1.85.0" {
		t.Errorf("expected installed dot and version, got %s, %s", dotInstalled, verInstalled)
	}

	dotMissing, verMissing := resolveTreeNodeStatus("nonexistent_tool_xyz", installed)
	if dotMissing != missingDot || verMissing != "—" {
		t.Errorf("expected missing dot and dash, got %s, %s", dotMissing, verMissing)
	}
}

func TestFormatProfileTreeNode(t *testing.T) {
	installed := map[string]string{"git": "2.40.0"}
	row := formatProfileTreeNode("git", false, installed)
	if !strings.Contains(row, "git") {
		t.Errorf("expected row to contain tool name git, got: %s", row)
	}

	if !strings.Contains(row, "2.40.0") {
		t.Errorf("expected row to contain version 2.40.0, got: %s", row)
	}
}

func TestProfileTreeKeywords(t *testing.T) {
	if !isProfileTreeKeyword("tree") {
		t.Errorf("expected 'tree' to be recognized as tree keyword")
	}

	if !isProfileTreeKeyword("--tree") {
		t.Errorf("expected '--tree' to be recognized as tree keyword")
	}

	if !isProfileTreeKeyword("-t") {
		t.Errorf("expected '-t' to be recognized as tree keyword")
	}

	if isProfileTreeKeyword("dev") {
		t.Errorf("did not expect 'dev' to be recognized as tree keyword")
	}
}

func TestLinuxJsHandlerResolution(t *testing.T) {
	hPnpm := resolveLinuxJsHandler(constants.ToolPnpm)
	if hPnpm == nil {
		t.Errorf("expected pnpm to resolve a linux JS handler")
	}

	hYarn := resolveLinuxJsHandler(constants.ToolYarn)
	if hYarn == nil {
		t.Errorf("expected yarn to resolve a linux JS handler")
	}

	hBun := resolveLinuxJsHandler(constants.ToolBun)
	if hBun == nil {
		t.Errorf("expected bun to resolve a linux JS handler")
	}

	hUnknown := resolveLinuxJsHandler("unknown_tool")
	if hUnknown != nil {
		t.Errorf("expected unknown tool to return nil handler")
	}
}

func TestResolveGlobalNpmCommand(t *testing.T) {
	cmd := resolveGlobalNpmCommand("pnpm")
	if len(cmd) < 4 {
		t.Fatalf("expected at least 4 args in npm command, got %v", cmd)
	}

	lastArg := cmd[len(cmd)-1]
	if lastArg != "pnpm" {
		t.Errorf("expected last arg to be pnpm, got %s", lastArg)
	}
}
