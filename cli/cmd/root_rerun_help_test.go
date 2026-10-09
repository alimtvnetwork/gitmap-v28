package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

func TestNormalizeHelpTopic_RerunAliases(t *testing.T) {
	aliases := []string{"rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue"}
	for _, alias := range aliases {
		normalized := normalizeHelpTopic(alias)
		if normalized != "rerun" {
			t.Errorf("normalizeHelpTopic(%q) = %q, expected \"rerun\"", alias, normalized)
		}
	}
}

func TestTryRenderRichTopic_Rerun(t *testing.T) {
	topics := []string{"rerun", "rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue"}
	for _, top := range topics {
		if !tryRenderRichTopic(top) {
			t.Errorf("expected tryRenderRichTopic(%q) to be true", top)
		}
	}
}

func TestAllHelpRows_ContainsRerun(t *testing.T) {
	rows := allHelpRows()
	var found bool
	for _, r := range rows {
		if r.Group == constants.HelpGroupIntegrations && strings.Contains(r.Line, "rerun") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("allHelpRows() did not contain rerun in %s", constants.HelpGroupIntegrations)
	}
}

func TestGetTopicDetailedSummary_Rerun(t *testing.T) {
	summary := helpdoc.GetTopicDetailedSummary("rerun")
	if !strings.Contains(summary, "Antigravity") {
		t.Errorf("expected summary to contain 'Antigravity', got %q", summary)
	}

	aliasSummary := helpdoc.GetTopicDetailedSummary("rr")
	if !strings.Contains(aliasSummary, "Antigravity") {
		t.Errorf("expected summary for 'rr' to contain 'Antigravity', got %q", aliasSummary)
	}
}
