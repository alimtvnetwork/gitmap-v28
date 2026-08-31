package cmd

import (
	"os"
	"testing"
)

func TestParsePipelineAIDelay(t *testing.T) {
	tests := []struct {
		args        []string
		wantDelay   int
		wantSubArgs []string
	}{
		{
			args:        []string{"status"},
			wantDelay:   20,
			wantSubArgs: []string{"status"},
		},
		{
			args:        []string{"status", "-t", "35", "--json"},
			wantDelay:   35,
			wantSubArgs: []string{"status", "--json"},
		},
		{
			args:        []string{"status", "-t", "5"},
			wantDelay:   20,
			wantSubArgs: []string{"status"},
		},
		{
			args:        []string{"eta", "--delay", "50"},
			wantDelay:   50,
			wantSubArgs: []string{"eta"},
		},
		{
			args:        []string{"etc", "--time", "65"},
			wantDelay:   65,
			wantSubArgs: []string{"etc"},
		},
	}

	for _, tt := range tests {
		gotDelay, gotSubArgs := parsePipelineAIDelay(tt.args)
		if gotDelay != tt.wantDelay {
			t.Errorf("parsePipelineAIDelay(%v) delay = %d, want %d", tt.args, gotDelay, tt.wantDelay)
		}
		if len(gotSubArgs) != len(tt.wantSubArgs) {
			t.Errorf("parsePipelineAIDelay(%v) subArgs = %v, want %v", tt.args, gotSubArgs, tt.wantSubArgs)
		}
	}
}

func TestRunPipelineAI_SkipDelay(t *testing.T) {
	_ = os.Setenv("GITMAP_SKIP_DELAY", "1")
	defer os.Unsetenv("GITMAP_SKIP_DELAY")

	err := runPipelineAI([]string{"status", "-t", "25", "--json"})
	if err != nil {
		t.Fatalf("runPipelineAI failed with GITMAP_SKIP_DELAY=1: %v", err)
	}
}
