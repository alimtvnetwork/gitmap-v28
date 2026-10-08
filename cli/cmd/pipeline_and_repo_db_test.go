package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
)

func TestPipelineDBDispatch(t *testing.T) {
	// 1. Status
	if err := cmdpipeline.HandlePipelineDB([]string{"status"}); err != nil {
		t.Fatalf("pipeline db status failed: %v", err)
	}

	// 2. Status JSON
	if err := cmdpipeline.HandlePipelineDB([]string{"status", "--json"}); err != nil {
		t.Fatalf("pipeline db status --json failed: %v", err)
	}

	// 3. Error logs
	if err := cmdpipeline.HandlePipelineDB([]string{"error-logs"}); err != nil {
		t.Fatalf("pipeline db error-logs failed: %v", err)
	}

	// 4. Optimize
	if err := cmdpipeline.HandlePipelineDB([]string{"optimize"}); err != nil {
		t.Fatalf("pipeline db optimize failed: %v", err)
	}

	// 5. Help
	if err := cmdpipeline.HandlePipelineDB([]string{"help"}); err != nil {
		t.Fatalf("pipeline db help failed: %v", err)
	}

	// 6. Unknown subcommand error check
	err := cmdpipeline.HandlePipelineDB([]string{"unknown-sub"})
	if err == nil || !strings.Contains(err.Error(), "unknown pipeline db subcommand") {
		t.Fatalf("expected unknown subcommand error, got %v", err)
	}
}

func TestRepoDBDispatch(t *testing.T) {
	// 1. Status
	if err := runRepoCommand([]string{"db", "status"}); err != nil {
		t.Fatalf("repo db status failed: %v", err)
	}

	// 2. Status JSON
	if err := runRepoCommand([]string{"db", "status", "--json"}); err != nil {
		t.Fatalf("repo db status --json failed: %v", err)
	}

	// 3. Log
	if err := runRepoCommand([]string{"db", "log"}); err != nil {
		t.Fatalf("repo db log failed: %v", err)
	}

	// 4. Error logs
	if err := runRepoCommand([]string{"db", "error-logs"}); err != nil {
		t.Fatalf("repo db error-logs failed: %v", err)
	}

	// 5. Optimize
	if err := runRepoCommand([]string{"db", "optimize"}); err != nil {
		t.Fatalf("repo db optimize failed: %v", err)
	}

	// 6. Help
	if err := runRepoCommand([]string{"help"}); err != nil {
		t.Fatalf("repo help failed: %v", err)
	}
}

func TestDBStatusAndOptimizeDispatch(t *testing.T) {
	if err := cmddb.RunDB([]string{"status"}); err != nil {
		t.Fatalf("gitmap db status failed: %v", err)
	}

	if err := cmddb.RunDB([]string{"status", "--json"}); err != nil {
		t.Fatalf("gitmap db status --json failed: %v", err)
	}

	if err := cmddb.RunDB([]string{"optimize"}); err != nil {
		t.Fatalf("gitmap db optimize failed: %v", err)
	}
}
