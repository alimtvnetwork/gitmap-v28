package cmdagy

import (
	"testing"
)

func TestComputeBaseMatchScore(t *testing.T) {
	isBoth := computeBaseMatchScore(true, true)
	if isBoth != 1000 {
		t.Errorf("expected 1000, got %d", isBoth)
	}

	isPath := computeBaseMatchScore(true, false)
	if isPath != 500 {
		t.Errorf("expected 500, got %d", isPath)
	}

	isName := computeBaseMatchScore(false, true)
	if isName != 200 {
		t.Errorf("expected 200, got %d", isName)
	}
}

func TestComputeSummaryScore(t *testing.T) {
	rRoot := summaryRow{
		id:           "conv-root",
		title:        "Antigravity-Manager",
		workspaceURI: `["file:///d%3A/work/Antigravity-Manager"]`,
		projectID:    "proj-123",
		depth:        0,
		stepCount:    100,
	}

	scoreRoot := computeSummaryScore(rRoot, "d:\\work\\antigravity-manager", "antigravity-manager", []string{"proj-123"})
	if scoreRoot != 1300 {
		t.Errorf("expected 1300 for root exact match, got %d", scoreRoot)
	}

	rSub := summaryRow{
		id:           "conv-sub",
		title:        "",
		workspaceURI: `["file:///d%3A/work/Antigravity-Manager"]`,
		projectID:    "proj-123",
		depth:        1,
		stepCount:    50,
	}
	scoreSub := computeSummaryScore(rSub, "d:\\work\\antigravity-manager", "antigravity-manager", []string{"proj-123"})
	if scoreSub != 1000 {
		t.Errorf("expected 1000 for subagent, got %d", scoreSub)
	}

	rUnrelated := summaryRow{
		id:           "conv-other",
		title:        "other-project",
		workspaceURI: `["file:///d%3A/work/other"]`,
		projectID:    "proj-999",
		depth:        0,
		stepCount:    10,
	}
	scoreOther := computeSummaryScore(rUnrelated, "d:\\work\\antigravity-manager", "antigravity-manager", []string{"proj-123"})
	if scoreOther != 0 {
		t.Errorf("expected 0 for unrelated project, got %d", scoreOther)
	}
}

func TestCheckProjectNameMatch(t *testing.T) {
	r := summaryRow{
		id:        "conv-1",
		title:     "AGM - Antigravity Manager",
		projectID: "proj-123",
	}

	isIDMatch := checkProjectNameMatch(r, "other", []string{"proj-123"})
	if isIDMatch == false {
		t.Errorf("expected true for project ID match")
	}

	isTitleMatch := checkProjectNameMatch(r, "Antigravity", nil)
	if isTitleMatch == false {
		t.Errorf("expected true for title match")
	}

	isMismatch := checkProjectNameMatch(r, "nomatch", nil)
	if isMismatch {
		t.Errorf("expected false for mismatch")
	}
}

func TestIsProjectIDInList(t *testing.T) {
	list := []string{"id-1", "id-2", "id-3"}
	hasID := isProjectIDInList("id-2", list)
	if hasID == false {
		t.Errorf("expected true for existing ID")
	}

	hasMissing := isProjectIDInList("id-9", list)
	if hasMissing {
		t.Errorf("expected false for missing ID")
	}

	hasEmpty := isProjectIDInList("", list)
	if hasEmpty {
		t.Errorf("expected false for empty ID")
	}
}

func TestSortScoredSummaries(t *testing.T) {
	items := []scoredSummary{
		{info: AgyConvInfo{ID: "low-score", StepCount: 100}, score: 500},
		{info: AgyConvInfo{ID: "high-score", StepCount: 10}, score: 1300},
		{info: AgyConvInfo{ID: "mid-score-more-steps", StepCount: 50}, score: 1000},
		{info: AgyConvInfo{ID: "mid-score-fewer-steps", StepCount: 20}, score: 1000},
	}

	sortScoredSummaries(items)
	if items[0].info.ID != "high-score" {
		t.Errorf("expected high-score first, got %s", items[0].info.ID)
	}
	if items[1].info.ID != "mid-score-more-steps" {
		t.Errorf("expected mid-score-more-steps second, got %s", items[1].info.ID)
	}
	if items[2].info.ID != "mid-score-fewer-steps" {
		t.Errorf("expected mid-score-fewer-steps third, got %s", items[2].info.ID)
	}
	if items[3].info.ID != "low-score" {
		t.Errorf("expected low-score fourth, got %s", items[3].info.ID)
	}

	infos := extractConvInfos(items)
	if len(infos) != 4 {
		t.Errorf("expected 4 infos, got %d", len(infos))
	}
}
