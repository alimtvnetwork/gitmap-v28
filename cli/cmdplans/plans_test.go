package cmdplans

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"
)

func TestParsePlansOptions_Defaults(t *testing.T) {
	opts, isProceed := parsePlansOptions([]string{})
	if !isProceed {
		t.Fatal("expected isProceed true")
	}
	if opts.Limit != constants.DefaultLimit {
		t.Errorf("expected limit %d, got %d", constants.DefaultLimit, opts.Limit)
	}
	if opts.RepoPath != constants.DefaultRepoPath {
		t.Errorf("expected repo path %q, got %q", constants.DefaultRepoPath, opts.RepoPath)
	}
	if opts.Format != enums.PlansOutputFormatTable {
		t.Errorf("expected default table format, got %q", opts.Format)
	}
}

func TestParsePlansOptions_Help(t *testing.T) {
	_, isProceed := parsePlansOptions([]string{"-h"})
	if isProceed {
		t.Error("expected isProceed false for -h")
	}

	_, isProceedHelp := parsePlansOptions([]string{"--help"})
	if isProceedHelp {
		t.Error("expected isProceed false for --help")
	}
}

func TestParsePlansOptions_Flags(t *testing.T) {
	opts, isProceed := parsePlansOptions([]string{"--json", "--pending", "-n", "25"})
	if !isProceed {
		t.Fatal("expected isProceed true")
	}
	if !opts.IsJSON || !opts.Format.IsJSON() {
		t.Error("expected JSON format set")
	}
	if !opts.IsPendingOnly || !opts.FilterMode.IsPending() {
		t.Error("expected pending filter mode set")
	}
	if opts.Limit != 25 {
		t.Errorf("expected limit 25, got %d", opts.Limit)
	}
}

func TestParsePlansOptions_YamlAndCompleted(t *testing.T) {
	opts, isProceed := parsePlansOptions([]string{"--yaml", "-c", "30"})
	if !isProceed {
		t.Fatal("expected isProceed true")
	}
	if !opts.IsYAML || !opts.Format.IsYAML() {
		t.Error("expected YAML format set")
	}
	if !opts.IsCompletedOnly || !opts.FilterMode.IsCompleted() {
		t.Error("expected completed filter mode set")
	}
	if opts.Limit != 30 {
		t.Errorf("expected positional limit 30, got %d", opts.Limit)
	}
}

func TestParsePlansOptions_AllFlag(t *testing.T) {
	opts, isProceed := parsePlansOptions([]string{"-a"})
	if !isProceed {
		t.Fatal("expected isProceed true")
	}
	if !opts.IsAll || !opts.FilterMode.IsAll() {
		t.Error("expected All filter mode set")
	}
}

func TestSlicePlans(t *testing.T) {
	items := []PlanItem{
		{ID: 1, Title: "Item 1"},
		{ID: 2, Title: "Item 2"},
		{ID: 3, Title: "Item 3"},
	}

	sliced := slicePlans(items, 2, false)
	if len(sliced) != 2 {
		t.Errorf("expected 2 items, got %d", len(sliced))
	}

	all := slicePlans(items, 2, true)
	if len(all) != 3 {
		t.Errorf("expected 3 items when isAll true, got %d", len(all))
	}
}

func TestBuildSummary(t *testing.T) {
	summary := buildSummary(15, 25, 10, 10, 10, false)
	if summary.CompletedCount != 15 || summary.PendingCount != 25 {
		t.Errorf("unexpected counts in summary")
	}
	if summary.TotalCount != 40 {
		t.Errorf("expected total count 40, got %d", summary.TotalCount)
	}
	if !summary.HasMorePending || !summary.HasMoreCompleted {
		t.Error("expected overflow flags to be true")
	}
}

func TestFilterPlanReport(t *testing.T) {
	comp := []PlanItem{{ID: 1, Status: StatusCompleted}}
	pend := []PlanItem{{ID: 2, Status: StatusPending}}

	opts := PlansOptions{Limit: 10, IsPendingOnly: true}
	c, p := filterPlanReport(opts, comp, pend)
	if len(c) != 0 || len(p) != 1 {
		t.Errorf("expected only pending items")
	}

	optsComp := PlansOptions{Limit: 10, IsCompletedOnly: true}
	c2, p2 := filterPlanReport(optsComp, comp, pend)
	if len(c2) != 1 || len(p2) != 0 {
		t.Errorf("expected only completed items")
	}
}
