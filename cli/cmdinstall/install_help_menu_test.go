package cmdinstall

import (
	"strings"
	"testing"
)

func TestBuildInstallOSIntegrationsSection_ContainsContextMenu(t *testing.T) {
	sec := buildInstallOSIntegrationsSection()
	hasContextMenu := false

	for _, e := range sec.Entries {
		if strings.Contains(e.Command, "context-menu") {
			hasContextMenu = true
			break
		}
	}

	if !hasContextMenu {
		t.Errorf("expected OS Integrations section to contain 'context-menu', got %+v", sec.Entries)
	}
}

func TestResolveInstallHelpSections_ContainsOSIntegrations(t *testing.T) {
	sections := resolveInstallHelpSections()
	hasOSSection := false

	for _, sec := range sections {
		if strings.Contains(sec.Title, "OS") {
			hasOSSection = true
			break
		}
	}

	if !hasOSSection {
		t.Errorf("expected sections to contain OS section, got %+v", sections)
	}
}
