package cmd

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func mkRun(kind, started string) model.MakeAllVisibilityRunRecord {
	return model.MakeAllVisibilityRunRecord{CommandKind: kind, StartedAt: started}
}

func TestParseHistoryFilters(t *testing.T) {
	now := time.Now()
	f := parseHistoryFilters([]string{"--kind", "VisibilityUndo", "--since", "24h"}, now)
	if f.Kind != "VisibilityUndo" || f.Since != 24*time.Hour {
		t.Fatalf("parse: got %+v", f)
	}
	if got := parseHistoryFilters([]string{"--since", "garbage"}, now); got.Since != 0 {
		t.Fatalf("bad --since must be ignored, got %+v", got)
	}
	if got := parseHistoryFilters(nil, now); got.Kind != "" || got.Since != 0 {
		t.Fatalf("empty must be zero, got %+v", got)
	}
}

func TestApplyHistoryFilters(t *testing.T) {
	now := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	runs := []model.MakeAllVisibilityRunRecord{
		mkRun(constants.CommandKindVisibilityRedo, "2026-06-06T11:30:00Z"),
		mkRun(constants.CommandKindVisibilityUndo, "2026-06-06T10:00:00Z"),
		mkRun(constants.CommandKindMakeAllPublic, "2026-06-05T08:00:00Z"),
		mkRun(constants.CommandKindMakeAllPrivate, "bogus-ts"),
	}
	if got := applyHistoryFilters(runs, historyFilters{}, now); len(got) != 4 {
		t.Fatalf("no-op must keep all, got %d", len(got))
	}
	kind := applyHistoryFilters(runs, historyFilters{Kind: constants.CommandKindVisibilityUndo}, now)
	if len(kind) != 1 || kind[0].CommandKind != constants.CommandKindVisibilityUndo {
		t.Fatalf("kind filter: %+v", kind)
	}
	since := applyHistoryFilters(runs, historyFilters{Since: 6 * time.Hour}, now)
	if len(since) != 2 {
		t.Fatalf("--since 6h must keep 2 (drops 28h-old + bogus-ts), got %d", len(since))
	}
	both := applyHistoryFilters(runs, historyFilters{
		Kind: constants.CommandKindVisibilityRedo, Since: 6 * time.Hour}, now)
	if len(both) != 1 {
		t.Fatalf("combined filter: %+v", both)
	}
}
