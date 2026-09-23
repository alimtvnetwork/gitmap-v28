package cmdpipeline

import (
	"testing"
)

func TestGroupRunsByCommit_MultipleWorkflows(t *testing.T) {
	runs := buildMultiWorkflowTestRuns()
	groups := GroupRunsByCommit(runs)
	verifyGroupLengths(t, groups, 2)
	verifyFirstGroup(t, groups[0])
	verifySecondGroup(t, groups[1])
}

func buildMultiWorkflowTestRuns() []ghRunItem {
	runs := make([]ghRunItem, 0, 4)
	runs = append(runs, makeTestRun(101, "CI", "completed", "success", "aaa1111"))
	runs = append(runs, makeTestRun(102, "Release", "completed", "success", "aaa1111"))
	runs = append(runs, makeTestRun(103, "CI", "completed", "failure", "bbb2222"))
	runs = append(runs, makeTestRun(104, "Lint", "completed", "success", "bbb2222"))

	return runs
}

func makeTestRun(id uint64, name, status, conclusion, sha string) ghRunItem {
	return ghRunItem{
		DatabaseId: id,
		Name:       name,
		Status:     status,
		Conclusion: conclusion,
		HeadSha:    sha,
		CreatedAt:  "2026-09-11T10:00:00Z",
		UpdatedAt:  "2026-09-11T10:02:00Z",
	}
}

func verifyGroupLengths(t *testing.T, groups []CommitPipelineGroup, expected int) {
	if len(groups) != expected {
		t.Fatalf("expected %d groups, got %d", expected, len(groups))
	}
}

func verifyFirstGroup(t *testing.T, g CommitPipelineGroup) {
	if g.HeadSha != "aaa1111" || g.Conclusion != "success" {
		t.Fatalf("unexpected first group: sha=%s, conclusion=%s", g.HeadSha, g.Conclusion)
	}
	if len(g.Workflows) != 2 || g.PassedWorkflows != 2 {
		t.Fatalf("unexpected workflows in first group: count=%d", len(g.Workflows))
	}
}

func verifySecondGroup(t *testing.T, g CommitPipelineGroup) {
	if g.HeadSha != "bbb2222" || g.Conclusion != "failure" {
		t.Fatalf("unexpected second group: sha=%s, conclusion=%s", g.HeadSha, g.Conclusion)
	}
	if g.FailedWorkflows != 1 || g.PassedWorkflows != 1 {
		t.Fatalf("unexpected counts: failed=%d, passed=%d", g.FailedWorkflows, g.PassedWorkflows)
	}
}

func TestGroupRunsByCommit_CancelledWorkflows(t *testing.T) {
	runs := []ghRunItem{
		makeTestRun(301, "CI", "completed", "cancelled", "12e40b1"),
		makeTestRun(302, "Release", "completed", "success", "12e40b1"),
		makeTestRun(303, "CI", "completed", "cancelled", "5717700"),
		makeTestRun(304, "Release", "completed", "cancelled", "5717700"),
	}
	groups := GroupRunsByCommit(runs)
	verifyGroupLengths(t, groups, 2)
	verifyCancelledFirstGroup(t, groups[0])
	verifyCancelledSecondGroup(t, groups[1])
}

func verifyCancelledFirstGroup(t *testing.T, g CommitPipelineGroup) {
	if g.HeadSha != "12e40b1" || g.Conclusion != "failure" {
		t.Fatalf("unexpected first group conclusion: sha=%s, conclusion=%s", g.HeadSha, g.Conclusion)
	}
	if g.FailedWorkflows != 1 || g.PassedWorkflows != 1 {
		t.Fatalf("unexpected counts: failed=%d, passed=%d", g.FailedWorkflows, g.PassedWorkflows)
	}
}

func verifyCancelledSecondGroup(t *testing.T, g CommitPipelineGroup) {
	if g.HeadSha != "5717700" || g.Conclusion != "failure" {
		t.Fatalf("unexpected second group conclusion: sha=%s, conclusion=%s", g.HeadSha, g.Conclusion)
	}
	if g.FailedWorkflows != 2 || g.PassedWorkflows != 0 {
		t.Fatalf("unexpected counts: failed=%d, passed=%d", g.FailedWorkflows, g.PassedWorkflows)
	}
}

func TestResolveCommitGroupByOffset(t *testing.T) {
	groups := buildSampleGroups()
	verifyOffsetResolutions(t, groups)
	verifyOutOfRangeOffsets(t, groups)
}

func buildSampleGroups() []CommitPipelineGroup {
	return []CommitPipelineGroup{
		{HeadSha: "sha0"},
		{HeadSha: "sha1"},
		{HeadSha: "sha2"},
		{HeadSha: "sha3"},
	}
}

func verifyOffsetResolutions(t *testing.T, groups []CommitPipelineGroup) {
	assertOffsetMatch(t, groups, 0, "sha0")
	assertOffsetMatch(t, groups, -1, "sha1")
	assertOffsetMatch(t, groups, -2, "sha2")
	assertOffsetMatch(t, groups, -3, "sha3")
}

func assertOffsetMatch(t *testing.T, groups []CommitPipelineGroup, offset int, expectedSha string) {
	group, isFound := ResolveCommitGroupByOffset(groups, offset)
	if isFound && group.HeadSha == expectedSha {
		return
	}

	t.Fatalf("expected offset %d to resolve to %s", offset, expectedSha)
}

func verifyOutOfRangeOffsets(t *testing.T, groups []CommitPipelineGroup) {
	_, isFoundNegative := ResolveCommitGroupByOffset(groups, -4)
	if isFoundNegative {
		t.Fatalf("expected -4 to be out of range")
	}

	_, isFoundEmpty := ResolveCommitGroupByOffset(nil, 0)
	if isFoundEmpty {
		t.Fatalf("expected empty groups to return false")
	}
}

func TestResolveCommitGroupByTarget(t *testing.T) {
	groups := buildTargetGroups()
	testTargetKeywords(t, groups)
	testTargetShas(t, groups)
}

func buildTargetGroups() []CommitPipelineGroup {
	return []CommitPipelineGroup{
		{HeadSha: "abcdef123456"},
		{HeadSha: "1234567890ab"},
		{HeadSha: "fedcba987654"},
	}
}

func testTargetKeywords(t *testing.T, groups []CommitPipelineGroup) {
	assertTargetMatch(t, groups, "", "abcdef123456")
	assertTargetMatch(t, groups, "latest", "abcdef123456")
	assertTargetMatch(t, groups, "-1", "1234567890ab")
	assertTargetMatch(t, groups, "-2", "fedcba987654")
}

func testTargetShas(t *testing.T, groups []CommitPipelineGroup) {
	assertTargetMatch(t, groups, "abcdef1", "abcdef123456")
	assertTargetMatch(t, groups, "fedcba", "fedcba987654")
	_, isFound := ResolveCommitGroupByTarget(groups, "nonexistent")
	if isFound {
		t.Fatalf("expected nonexistent target to return false")
	}
}

func assertTargetMatch(t *testing.T, groups []CommitPipelineGroup, target, expectedSha string) {
	group, isFound := ResolveCommitGroupByTarget(groups, target)
	if isFound && group.HeadSha == expectedSha {
		return
	}

	t.Fatalf("expected target %s to resolve to %s", target, expectedSha)
}

func TestGroupRunsByCommit_InProgressConclusion(t *testing.T) {
	runs := buildInProgressRuns()
	groups := GroupRunsByCommit(runs)
	verifyInProgressGroup(t, groups)
	testFailurePrecedence(t)
}

func buildInProgressRuns() []ghRunItem {
	return []ghRunItem{
		{DatabaseId: 201, Name: "CI", Status: "completed", Conclusion: "success", HeadSha: "c1"},
		{DatabaseId: 202, Name: "Deploy", Status: "in_progress", Conclusion: "", HeadSha: "c1"},
	}
}

func verifyInProgressGroup(t *testing.T, groups []CommitPipelineGroup) {
	verifyGroupLengths(t, groups, 1)
	if groups[0].Status != "in_progress" || groups[0].Conclusion != "in_progress" {
		t.Fatalf("expected in_progress status/conclusion, got %s/%s", groups[0].Status, groups[0].Conclusion)
	}
	if groups[0].ActiveWorkflows != 1 {
		t.Fatalf("expected 1 active workflow, got %d", groups[0].ActiveWorkflows)
	}
}

func testFailurePrecedence(t *testing.T) {
	runs := buildFailurePrecedenceRuns()
	groups := GroupRunsByCommit(runs)
	if groups[0].Conclusion != "failure" {
		t.Fatalf("expected failure conclusion to take precedence, got %s", groups[0].Conclusion)
	}
}

func buildFailurePrecedenceRuns() []ghRunItem {
	return []ghRunItem{
		{DatabaseId: 203, Name: "CI", Status: "completed", Conclusion: "failure", HeadSha: "c2"},
		{DatabaseId: 204, Name: "Deploy", Status: "in_progress", Conclusion: "", HeadSha: "c2"},
	}
}

func TestCollapseConsecutiveEmptyLines(t *testing.T) {
	testCollapseRunawayLines(t)
	testCollapseWhitespaceLines(t)
	testCollapseSingleEmpty(t)
}

func testCollapseRunawayLines(t *testing.T) {
	input := "line1\n\n\n\nline2"
	expected := "line1\n\nline2"
	got := CollapseConsecutiveEmptyLines(input)
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func testCollapseWhitespaceLines(t *testing.T) {
	input := "a\n   \n\t\nb"
	expected := "a\n\nb"
	got := CollapseConsecutiveEmptyLines(input)
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func testCollapseSingleEmpty(t *testing.T) {
	input := "a\n\nb"
	expected := "a\n\nb"
	got := CollapseConsecutiveEmptyLines(input)
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
