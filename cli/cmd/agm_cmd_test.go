package cmd

import (
	"context"
	"testing"
)

func TestDispatchAgm(t *testing.T) {
	err := dispatchAgm(context.Background(), []string{"agm", "install", "--dry-run"}, nil)
	if err != nil {
		t.Errorf("expected nil error for dispatchAgm dry-run, got %v", err)
	}
}

func TestDispatchAgmUpdate(t *testing.T) {
	err := dispatchAgm(context.Background(), []string{"agm", "update", "--dry-run"}, nil)
	if err != nil {
		t.Errorf("expected nil error for dispatchAgm update dry-run, got %v", err)
	}
}

func TestIsAgmUpdateTarget(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{[]string{"agm"}, true},
		{[]string{"ag-manager"}, true},
		{[]string{"antigravity-manager"}, true},
		{[]string{"--dry-run", "agm"}, true},
		{[]string{"--force"}, false},
		{[]string{"gitmap"}, false},
	}

	for _, tc := range tests {
		if got := isAgmUpdateTarget(tc.args); got != tc.want {
			t.Errorf("isAgmUpdateTarget(%v) = %v; want %v", tc.args, got, tc.want)
		}
	}
}

func TestRunUpdateAgManagerTargetDryRun(t *testing.T) {
	args := []string{"agm", "--dry-run"}
	err := runUpdateAgManagerTarget(args)
	if err != nil {
		t.Errorf("expected nil error for runUpdateAgManagerTarget dry-run, got %v", err)
	}
}
