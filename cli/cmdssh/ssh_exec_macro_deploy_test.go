package cmdssh

import (
	"testing"
)

func TestParseMacroAddTargetName(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantTarget string
		wantOk     bool
	}{
		{
			name:       "gitmap macro add with single target",
			args:       []string{"gitmap", "macro", "add", "alim1"},
			wantTarget: "alim1",
			wantOk:     true,
		},
		{
			name:       "macro add without gitmap prefix",
			args:       []string{"macro", "add", "deploy_task"},
			wantTarget: "deploy_task",
			wantOk:     true,
		},
		{
			name:       "combined quoted command string",
			args:       []string{"gitmap macro add alim1"},
			wantTarget: "alim1",
			wantOk:     true,
		},
		{
			name:       "macro add with flags only",
			args:       []string{"macro", "add", "alim1", "--desc", "test"},
			wantTarget: "alim1",
			wantOk:     true,
		},
		{
			name:       "macro add with real command steps",
			args:       []string{"macro", "add", "alim1", "echo hello"},
			wantTarget: "",
			wantOk:     false,
		},
		{
			name:       "unrelated gitmap command",
			args:       []string{"gitmap", "status"},
			wantTarget: "",
			wantOk:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotTarget, gotOk := parseMacroAddTargetName(tc.args)
			if gotTarget != tc.wantTarget || gotOk != tc.wantOk {
				t.Errorf("parseMacroAddTargetName(%v) = (%q, %v); want (%q, %v)",
					tc.args, gotTarget, gotOk, tc.wantTarget, tc.wantOk)
			}
		})
	}
}
