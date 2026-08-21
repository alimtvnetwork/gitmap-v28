package dashboard

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestNormalizeOptionsRecentEmptySince(t *testing.T) {
	opts := CollectOptions{
		Recent: true,
	}
	normalized := normalizeOptions(opts)
	if !normalized.Recent {
		t.Errorf("expected Recent to be true")
	}
	expected := time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02")
	if normalized.Since != expected {
		t.Errorf("expected Since to be %s, got %s", expected, normalized.Since)
	}
}

func TestNormalizeOptionsRecentExplicitSince(t *testing.T) {
	opts := CollectOptions{
		Recent: true,
		Since:  "2026-01-01",
	}
	normalized := normalizeOptions(opts)
	if normalized.Since != "2026-01-01" {
		t.Errorf("expected Since to remain 2026-01-01, got %s", normalized.Since)
	}
}

func TestNormalizeOptionsNonRecent(t *testing.T) {
	opts := CollectOptions{
		Recent: false,
	}
	normalized := normalizeOptions(opts)
	if len(normalized.Since) > 0 {
		t.Errorf("expected Since to remain empty, got %s", normalized.Since)
	}
}

func TestRecentSinceDateFormat(t *testing.T) {
	d := recentSinceDate()
	_, err := time.Parse("2006-01-02", d)
	if err != nil {
		t.Fatalf("recentSinceDate returned invalid format: %v", err)
	}
}

func TestBuildMeta(t *testing.T) {
	opts := CollectOptions{
		RepoPath: ".",
		Limit:    10,
		Since:    "2026-08-01",
		Recent:   true,
	}
	meta := buildMeta(opts, 5, 2, 1)
	if !meta.Recent {
		t.Errorf("expected meta.Recent to be true")
	}
	if meta.Limit != 10 || meta.Since != "2026-08-01" {
		t.Errorf("unexpected meta values: %+v", meta)
	}
	if meta.TotalCommits != 5 || meta.TotalBranches != 2 || meta.TotalTags != 1 {
		t.Errorf("unexpected totals in meta: %+v", meta)
	}
}

func TestAssembleDashboard(t *testing.T) {
	meta := model.DashboardMeta{RepoName: "test-repo", Recent: true}
	data := assembleDashboard(meta, nil, nil, nil, nil, model.FrequencyData{})
	if data.Meta.RepoName != "test-repo" || !data.Meta.Recent {
		t.Errorf("assembleDashboard did not preserve meta properly")
	}
}
