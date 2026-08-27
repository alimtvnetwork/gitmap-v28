package cmd

import (
	"testing"
)

func TestResolveProfileTree_Hierarchy(t *testing.T) {
	t.Parallel()

	// 1. Basic profile
	basicProfile, hasBasic := resolveProfileTree("ubuntu-basic")
	if !hasBasic {
		t.Fatal("expected ubuntu-basic to resolve")
	}
	if basicProfile.Name != "ubuntu-basic" || basicProfile.Base != nil {
		t.Errorf("unexpected basic profile structure: %+v", basicProfile)
	}
	if len(basicProfile.Tools) != 4 {
		t.Errorf("expected 4 tools in ubuntu-basic, got %d", len(basicProfile.Tools))
	}

	// 2. VSCode profile
	vscodeProfile, hasVscode := resolveProfileTree("ubuntu+vscode")
	if !hasVscode {
		t.Fatal("expected ubuntu+vscode to resolve")
	}
	if vscodeProfile.Base == nil || vscodeProfile.Base.Name != "ubuntu-basic" {
		t.Errorf("expected base ubuntu-basic, got %+v", vscodeProfile.Base)
	}

	// 3. Small Dev profile
	smallDevProfile, hasSmallDev := resolveProfileTree("ubuntu+small-dev")
	if !hasSmallDev {
		t.Fatal("expected ubuntu+small-dev to resolve")
	}
	if smallDevProfile.Base == nil || smallDevProfile.Base.Name != "ubuntu+vscode" {
		t.Errorf("expected base ubuntu+vscode, got %+v", smallDevProfile.Base)
	}

	// 4. Dev profile
	devProfile, hasDev := resolveProfileTree("ubuntu+dev")
	if !hasDev {
		t.Fatal("expected ubuntu+dev to resolve")
	}
	if devProfile.Base == nil || devProfile.Base.Name != "ubuntu+small-dev" {
		t.Errorf("expected base ubuntu+small-dev, got %+v", devProfile.Base)
	}
}

func TestResolveProfileTree_Aliases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		slug         string
		expectedName string
	}{
		{"ub", "ubuntu-basic"},
		{"ubuntu+basic", "ubuntu-basic"},
		{"UBUNTU-BASIC", "ubuntu-basic"},
		{"  ub+code  ", "ubuntu+vscode"},
		{"vscode", "vscode-settings"},
		{"ubuntu-vscode", "ubuntu+vscode"},
		{"ub+sdev", "ubuntu+small-dev"},
		{"small-dev", "ubuntu+small-dev"},
		{"ubuntu-small-dev", "ubuntu+small-dev"},
		{"ub+dev", "ubuntu+dev"},
		{"dev", "ubuntu+dev"},
		{"ubuntu-dev", "ubuntu+dev"},
	}

	for _, testCase := range testCases {
		profile, hasProfile := resolveProfileTree(testCase.slug)
		if !hasProfile {
			t.Errorf("slug %q failed to resolve", testCase.slug)
			continue
		}
		if profile.Name != testCase.expectedName {
			t.Errorf("slug %q resolved to %q, expected %q", testCase.slug, profile.Name, testCase.expectedName)
		}
	}
}

func TestResolveProfileTree_NotFound(t *testing.T) {
	t.Parallel()

	_, hasUnknown := resolveProfileTree("non-existent-profile")
	if hasUnknown {
		t.Error("expected non-existent-profile to return false")
	}
}

func TestProfileToTreeNode(t *testing.T) {
	t.Parallel()

	devProfile, hasDev := resolveProfileTree("ubuntu+dev")
	if !hasDev {
		t.Fatal("expected ubuntu+dev to resolve")
	}

	treeNode := profileToTreeNode(devProfile)
	if treeNode.Title != "ubuntu+dev" {
		t.Errorf("expected root title 'ubuntu+dev', got %q", treeNode.Title)
	}
	// Children should have BaseProfile node + 2 tools = 3 children
	if len(treeNode.Children) != 3 {
		t.Errorf("expected 3 children on root node, got %d", len(treeNode.Children))
	}
}

func TestPrintProfileInstallSummary(t *testing.T) {
	t.Parallel()

	// Ensure calling these does not panic
	printProfileInstallSummary("ubuntu+dev")
	printProfileInstallSummary("ubuntu-basic")
	printProfileInstallSummary("unknown-slug")
}
