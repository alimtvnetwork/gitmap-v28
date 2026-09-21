package cmdpipeline

import (
	"strings"
	"testing"
)

func TestBuildFailedChecksTree_Empty(t *testing.T) {
	p := PipelineErrorLogsPayload{}
	tree := BuildFailedChecksTree(p, false)
	if len(tree) != 0 {
		t.Fatalf("expected empty tree for empty failed runs, got %q", tree)
	}
}

func TestBuildFailedChecksTree_SingleWorkflowAndJob(t *testing.T) {
	p := PipelineErrorLogsPayload{
		LastHash: "7d6e74b8b5d",
		FailedRuns: []FailedRunItem{
			{
				WorkflowName: "CI",
				RunId:        35576139471,
				FailedJobs: []FailedJobItem{
					{
						JobName:  "Relative Path Check",
						StepName: "check-relative-paths",
					},
				},
			},
		},
	}

	tree := BuildFailedChecksTree(p, false)
	if !strings.Contains(tree, "Commit 7d6e74b:") {
		t.Errorf("missing commit sha in tree: %s", tree)
	}
	if !strings.Contains(tree, "└── ✖ CI (#35576139471)") {
		t.Errorf("missing workflow node in tree: %s", tree)
	}
	if !strings.Contains(tree, "└── ✖ Relative Path Check [Step: check-relative-paths]") {
		t.Errorf("missing job step node in tree: %s", tree)
	}
}

func TestBuildFailedChecksTree_MultipleWorkflowsAndJobs(t *testing.T) {
	p := PipelineErrorLogsPayload{
		Sha: "abc1234def",
		FailedRuns: []FailedRunItem{
			{
				WorkflowName: "CI",
				RunId:        101,
				FailedJobs: []FailedJobItem{
					{JobName: "Lint", StepName: "golangci-lint"},
					{JobName: "Paths", StepName: "check-relative-paths"},
				},
			},
			{
				WorkflowName: "Cross-Platform Build",
				RunId:        102,
				FailedJobs: []FailedJobItem{
					{JobName: "ubuntu-latest / go build + test", StepName: "go test ./..."},
				},
			},
		},
	}

	tree := BuildFailedChecksTree(p, false)
	if !strings.Contains(tree, "├── ✖ CI (#101)") {
		t.Errorf("expected branch prefix for first workflow: %s", tree)
	}
	if !strings.Contains(tree, "└── ✖ Cross-Platform Build (#102)") {
		t.Errorf("expected terminating branch for last workflow: %s", tree)
	}
	if !strings.Contains(tree, "│   ├── ✖ Lint [Step: golangci-lint]") {
		t.Errorf("expected indented intermediate job branch: %s", tree)
	}
	if !strings.Contains(tree, "│   └── ✖ Paths [Step: check-relative-paths]") {
		t.Errorf("expected indented terminating job branch: %s", tree)
	}
}

func TestBuildFailedChecksTree_ColorFormatting(t *testing.T) {
	p := PipelineErrorLogsPayload{
		LastHash: "1234567",
		FailedRuns: []FailedRunItem{
			{
				WorkflowName: "CI",
				RunId:        999,
				FailedJobs: []FailedJobItem{
					{JobName: "Test", StepName: "unit"},
				},
			},
		},
	}

	colored := BuildFailedChecksTree(p, true)
	plain := BuildFailedChecksTree(p, false)

	if !strings.Contains(colored, "\033[") {
		t.Errorf("expected ANSI escape codes in colored tree: %s", colored)
	}
	if strings.Contains(plain, "\033[") {
		t.Errorf("expected no ANSI codes in plain tree: %s", plain)
	}
}

func TestAppendClipboardFailureTree(t *testing.T) {
	p := PipelineErrorLogsPayload{
		LastHash: "abcdef1",
		FailedRuns: []FailedRunItem{
			{
				WorkflowName: "CI",
				RunId:        123,
				FailedJobs: []FailedJobItem{
					{JobName: "Build", StepName: "compile"},
				},
			},
		},
	}

	var sb strings.Builder
	appendClipboardFailureTree(&sb, p)
	out := sb.String()

	if !strings.Contains(out, "FAILED PIPELINE CHECKS TREE:") {
		t.Errorf("expected tree header in clipboard content: %s", out)
	}
	if !strings.Contains(out, "└── ✖ CI (#123)") {
		t.Errorf("missing workflow in clipboard tree: %s", out)
	}
}

func TestBuildCommitGroupFailureTree_NilOrNoFailures(t *testing.T) {
	if res := BuildCommitGroupFailureTree("repo", nil, false); len(res) > 0 {
		t.Errorf("expected empty string for nil group, got: %s", res)
	}

	group := &CommitPipelineGroup{
		HeadSha: "abcdef123456",
		Workflows: []CommitWorkflowItem{
			{Name: "CI", Conclusion: "success"},
		},
	}
	if res := BuildCommitGroupFailureTree("repo", group, false); len(res) > 0 {
		t.Errorf("expected empty string for passing group, got: %s", res)
	}
}

func TestBuildCommitGroupFailureTree_WithFailingWorkflow(t *testing.T) {
	group := &CommitPipelineGroup{
		HeadSha: "abcdef123456",
		Workflows: []CommitWorkflowItem{
			{Name: "CI", DatabaseId: 999, Conclusion: "failure"},
			{Name: "Lint", DatabaseId: 1000, Conclusion: "success"},
		},
	}
	res := BuildCommitGroupFailureTree("repo", group, false)
	if !strings.Contains(res, "FAILED PIPELINE CHECKS TREE:") {
		t.Errorf("expected tree header, got: %s", res)
	}
	if !strings.Contains(res, "CI (#999)") {
		t.Errorf("expected failing workflow in tree, got: %s", res)
	}
	if strings.Contains(res, "Lint") {
		t.Errorf("expected passing workflow to be omitted from tree, got: %s", res)
	}
}
