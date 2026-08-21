package dashboard

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestBuildAuthors(t *testing.T) {
	commits := []model.CommitInfo{
		{Author: "Alice", Email: "alice@example.com", Date: "2026-08-20T10:00:00Z"},
		{Author: "Alice", Email: "alice@example.com", Date: "2026-08-21T10:00:00Z"},
		{Author: "Bob", Email: "bob@example.com", Date: "2026-08-21T10:00:00Z"},
	}
	authors := buildAuthors(commits)
	if len(authors) != 2 {
		t.Fatalf("expected 2 authors, got %d", len(authors))
	}
	if authors[0].Email != "alice@example.com" || authors[0].TotalCommits != 2 {
		t.Errorf("expected Alice to be top author with 2 commits, got %+v", authors[0])
	}
}

func TestBuildFrequency(t *testing.T) {
	commits := []model.CommitInfo{
		{Date: "2026-08-05T10:00:00Z"},
		{Date: "2026-08-05T11:00:00Z"},
		{Date: "2026-08-15T10:00:00Z"},
	}
	freq := buildFrequency(commits)
	if freq.Daily["2026-08-05"] != 2 || freq.Monthly["2026-08"] != 3 {
		t.Errorf("unexpected frequency counts: %+v", freq)
	}
}

func TestAttachTagsToCommits(t *testing.T) {
	commits := []model.CommitInfo{
		{SHA: "111122223333", ShortSHA: "1111222"},
	}
	tags := []model.TagInfo{
		{Name: "v1.0.0", SHA: "111122223333"},
	}
	tagged := attachTagsToCommits(commits, tags)
	if len(tagged[0].Tags) != 1 || tagged[0].Tags[0] != "v1.0.0" {
		t.Errorf("expected tag v1.0.0 to be attached, got %+v", tagged[0].Tags)
	}
}
