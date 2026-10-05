package cmdpipeline

import (
	"testing"
)

func TestParsePipelineErrorFlags_All(t *testing.T) {
	flags1 := ParsePipelineErrorFlags([]string{"all"})
	if !flags1.IsAll {
		t.Errorf("expected IsAll to be true for 'all', got false")
	}

	flags2 := ParsePipelineErrorFlags([]string{"--all"})
	if !flags2.IsAll {
		t.Errorf("expected IsAll to be true for '--all', got false")
	}

	flags3 := ParsePipelineErrorFlags([]string{"pe", "all", "--json"})
	if !flags3.IsAll {
		t.Errorf("expected IsAll to be true for 'pe all --json', got false")
	}

	if !flags3.IsJSON {
		t.Errorf("expected IsJSON to be true for 'pe all --json', got false")
	}
}

func TestExecuteAllPipelineErrorLogs(t *testing.T) {
	flags := PipelineErrorFlags{
		IsAll:  true,
		IsJSON: true,
	}

	err := executeAllPipelineErrorLogs(flags, []string{"all", "--json"})
	if err != nil {
		t.Errorf("expected executeAllPipelineErrorLogs to succeed, got %v", err)
	}
}
