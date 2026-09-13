package cmdinstall

import (
	"context"
	"testing"
)

func TestAgmCommandFlags(t *testing.T) {
	dryFlag := agmInstallCmd.Flags().Lookup("dry-run")
	if dryFlag == nil {
		t.Errorf("expected --dry-run flag on agmInstallCmd")
	}

	verFlag := agmInstallCmd.Flags().Lookup("version")
	if verFlag == nil {
		t.Errorf("expected --version flag on agmInstallCmd")
	}

	yesFlag := agmInstallCmd.Flags().Lookup("yes")
	if yesFlag == nil {
		t.Errorf("expected --yes flag on agmInstallCmd")
	}

	verbFlag := agmInstallCmd.Flags().Lookup("verbose")
	if verbFlag == nil {
		t.Errorf("expected --verbose flag on agmInstallCmd")
	}
}

func TestIsAgmCommand(t *testing.T) {
	tests := []struct {
		arg  string
		want bool
	}{
		{"agm", true},
		{"AGM", true},
		{"ag-manager", true},
		{"antigravity-manager", true},
		{"agy", false},
		{"gitmap", false},
	}

	for _, tc := range tests {
		if got := isAgmCommand(tc.arg); got != tc.want {
			t.Errorf("isAgmCommand(%q) = %v; want %v", tc.arg, got, tc.want)
		}
	}
}

func TestDispatchAgmDryRun(t *testing.T) {
	args := []string{"agm", "install", "--dry-run"}
	err := DispatchAgm(context.Background(), args, nil)
	if err != nil {
		t.Errorf("expected nil error for agm dry-run dispatch, got %v", err)
	}
}
