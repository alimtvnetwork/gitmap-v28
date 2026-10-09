package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

func TestNormalizeHelpTopic_ShutdownUntilAliases(t *testing.T) {
	aliases := []string{"sug", "shutdown-until", "shutdown-until-green"}
	for _, alias := range aliases {
		normalized := normalizeHelpTopic(alias)
		if normalized != constants.CmdShutdownUntil {
			t.Errorf("normalizeHelpTopic(%q) = %q, expected %q", alias, normalized, constants.CmdShutdownUntil)
		}
	}
}

func TestTryRenderRichTopic_ShutdownUntil(t *testing.T) {
	topics := []string{"shutdown-until", "shutdown-until-green", "sug"}
	for _, top := range topics {
		if !tryRenderRichTopic(top) {
			t.Errorf("expected tryRenderRichTopic(%q) to be true", top)
		}
	}
}

func TestAllHelpRows_ContainsShutdownUntil(t *testing.T) {
	rows := allHelpRows()
	var found bool
	for _, r := range rows {
		if r.Group == constants.HelpGroupIntegrations && strings.Contains(r.Line, "shutdown-until") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("allHelpRows() did not contain shutdown-until in %s", constants.HelpGroupIntegrations)
	}
}

func TestHelpText_ShutdownUntilReadRaw(t *testing.T) {
	keys := []string{"shutdown-until", "sug", "shutdown-until-green"}
	for _, k := range keys {
		data, err := helpdoc.ReadRaw(k)
		if err != nil {
			t.Fatalf("helpdoc.ReadRaw(%q) failed: %v", k, err)
		}
		if !strings.Contains(string(data), "shutdown-until") {
			t.Errorf("expected helptext for %q to contain 'shutdown-until'", k)
		}
	}
}
