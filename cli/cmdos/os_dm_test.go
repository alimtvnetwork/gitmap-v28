package cmdos

import (
	"strings"
	"testing"
)

func TestUpdateWaylandInConfig_Enable(t *testing.T) {
	orig := "[daemon]\nWaylandEnable=false\nAutomaticLoginEnable=true\n"
	updated := updateWaylandInConfig(orig, true)

	if strings.Contains(updated, "WaylandEnable=false") {
		t.Errorf("expected WaylandEnable=false to be removed/updated, got:\n%s", updated)
	}

	if !strings.Contains(updated, "WaylandEnable=true") {
		t.Errorf("expected WaylandEnable=true, got:\n%s", updated)
	}
}

func TestUpdateWaylandInConfig_Disable(t *testing.T) {
	orig := "[daemon]\n#WaylandEnable=false\nAutomaticLoginEnable=true\n"
	updated := updateWaylandInConfig(orig, false)

	if !strings.Contains(updated, "WaylandEnable=false") {
		t.Errorf("expected WaylandEnable=false, got:\n%s", updated)
	}
}

func TestInsertWaylandUnderDaemon_Missing(t *testing.T) {
	orig := []string{"[daemon]", "AutomaticLoginEnable=true"}
	res := insertWaylandUnderDaemon(orig, false)

	if !strings.Contains(res, "WaylandEnable=false") {
		t.Errorf("expected WaylandEnable=false to be inserted under daemon, got:\n%s", res)
	}
}

func TestRunOSDMRouting(t *testing.T) {
	// help should exit cleanly with nil error
	if err := runOSDMCommand([]string{"help"}); err != nil {
		t.Errorf("expected nil error for help, got: %v", err)
	}
}
