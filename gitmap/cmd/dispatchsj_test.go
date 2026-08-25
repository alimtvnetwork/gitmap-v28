package cmd

import (
	"context"
	"testing"
)

func TestDispatchSJ(t *testing.T) {
	ctx := context.Background()

	// Test with no args
	err := dispatchSJ(ctx, []string{}, nil)
	if err != nil {
		t.Logf("dispatchSJ empty args returned: %v", err)
	}

	// Test with prefix sj
	err = dispatchSJ(ctx, []string{"sj"}, nil)
	if err != nil {
		t.Logf("dispatchSJ 'sj' returned: %v", err)
	}

	// Test with prefix ssh-join
	err = dispatchSJ(ctx, []string{"ssh-join"}, nil)
	if err != nil {
		t.Logf("dispatchSJ 'ssh-join' returned: %v", err)
	}

	// Test with prefix ssh-joined
	err = dispatchSJ(ctx, []string{"ssh-joined"}, nil)
	if err != nil {
		t.Logf("dispatchSJ 'ssh-joined' returned: %v", err)
	}

	// Test routing to ls
	err = dispatchSJ(ctx, []string{"sj", "ls"}, nil)
	if err != nil {
		t.Logf("dispatchSJ 'sj ls' returned: %v", err)
	}

	// Test routing to rm with missing args
	err = dispatchSJ(ctx, []string{"sj", "rm"}, nil)
	if err == nil {
		t.Errorf("expected error for rm with missing args")
	}

	// Test routing to history
	err = dispatchSJ(ctx, []string{"sj", "history"}, nil)
	if err != nil {
		t.Logf("dispatchSJ 'sj history' returned: %v", err)
	}

	// Test routing to add-auth with missing args
	err = dispatchSJ(ctx, []string{"sj", "add-auth"}, nil)
	if err == nil {
		t.Errorf("expected error for add-auth with missing args")
	}
}
