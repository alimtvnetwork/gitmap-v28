package store

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func modelKindRow(kind string) model.MakeAllVisibilityRunRecord {
	return model.MakeAllVisibilityRunRecord{
		CommandKind: kind, TargetVisibility: constants.VisibilityPublic,
		Provider: constants.ProviderGitHub, Owner: "acme", TargetRaw: "acme",
		PatternList: "*", StartedAt: "2026-06-06T11:00:00Z",
	}
}

func TestBuildRecentRunsQueryNoFilter(t *testing.T) {
	sql, args := BuildRecentRunsQuery(RecentRunsFilter{Limit: 20})
	if strings.Contains(sql, "WHERE") {
		t.Fatalf("no-filter must omit WHERE: %s", sql)
	}
	if len(args) != 1 || args[0] != 20 {
		t.Fatalf("args = %v, want [20]", args)
	}
}

func TestBuildRecentRunsQueryKindOnly(t *testing.T) {
	sql, args := BuildRecentRunsQuery(RecentRunsFilter{
		Kind: constants.CommandKindVisibilityUndo, Limit: 5,
	})
	if !strings.Contains(sql, "WHERE CommandKind = ?") {
		t.Fatalf("kind clause missing: %s", sql)
	}
	if strings.Contains(sql, " AND ") {
		t.Fatalf("must not emit AND with one clause: %s", sql)
	}
	if len(args) != 2 || args[0] != constants.CommandKindVisibilityUndo || args[1] != 5 {
		t.Fatalf("args = %v", args)
	}
}

func TestBuildRecentRunsQueryBothFilters(t *testing.T) {
	sql, args := BuildRecentRunsQuery(RecentRunsFilter{
		Kind: "MakeAllPublic", SinceISO: "2026-06-06T00:00:00Z", Limit: 50,
	})
	if !strings.Contains(sql, "WHERE CommandKind = ? AND StartedAt >= ?") {
		t.Fatalf("WHERE composition wrong: %s", sql)
	}
	if !strings.HasSuffix(sql, "ORDER BY MakeAllVisibilityRunId DESC LIMIT ?") {
		t.Fatalf("must end with ORDER+LIMIT: %s", sql)
	}
	want := []any{"MakeAllPublic", "2026-06-06T00:00:00Z", 50}
	if len(args) != 3 {
		t.Fatalf("args len = %d, want 3 (%v)", len(args), args)
	}
	for i, w := range want {
		if args[i] != w {
			t.Fatalf("args[%d] = %v, want %v", i, args[i], w)
		}
	}
}

func TestSelectRecentMakeAllVisibilityRunsFilteredKindPushdown(t *testing.T) {
	db := freshHistoryDB(t)
	_ = insertHistoryRun(t, db, "a") // MakeAllPublic
	id2, err := db.InsertMakeAllVisibilityRun(modelKindRow("VisibilityUndo"))
	if err != nil {
		t.Fatalf("insert undo: %v", err)
	}
	got, err := db.SelectRecentMakeAllVisibilityRunsFiltered(RecentRunsFilter{
		Kind: "VisibilityUndo", Limit: 10,
	})
	if err != nil || len(got) != 1 || got[0].ID != id2 {
		t.Fatalf("kind pushdown: got %+v err=%v", got, err)
	}
}
