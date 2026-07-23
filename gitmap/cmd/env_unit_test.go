package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

func TestEnvNamePatternValidates(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{"FOO", true},
		{"_FOO", true},
		{"foo_bar_1", true},
		{"GITMAP_HOME", true},
		{"1FOO", false},
		{"FOO-BAR", false},
		{"FOO BAR", false},
		{"FOO=BAR", false},
		{"", false},
	}
	for _, c := range cases {
		got := envNamePattern.MatchString(c.name)
		if got != c.ok {
			t.Errorf("envNamePattern.MatchString(%q) = %v, want %v", c.name, got, c.ok)
		}
	}
}

func TestUpsertEnvVariableInsertsAndUpdates(t *testing.T) {
	reg := model.EnvRegistry{Variables: []model.EnvVariable{}}
	reg = upsertEnvVariable(reg, "FOO", "1")
	if len(reg.Variables) != 1 || reg.Variables[0].Value != "1" {
		t.Fatalf("insert failed: %+v", reg.Variables)
	}
	reg = upsertEnvVariable(reg, "FOO", "2")
	if len(reg.Variables) != 1 || reg.Variables[0].Value != "2" {
		t.Fatalf("update failed: %+v", reg.Variables)
	}
	reg = upsertEnvVariable(reg, "BAR", "x")
	if len(reg.Variables) != 2 {
		t.Fatalf("second insert failed: %+v", reg.Variables)
	}
}

func TestRemoveEnvVariableDropsMatch(t *testing.T) {
	reg := model.EnvRegistry{Variables: []model.EnvVariable{
		{Name: "A", Value: "1"},
		{Name: "B", Value: "2"},
		{Name: "C", Value: "3"},
	}}
	out := removeEnvVariable(reg, "B")
	if len(out.Variables) != 2 {
		t.Fatalf("len = %d, want 2", len(out.Variables))
	}
	for _, v := range out.Variables {
		if v.Name == "B" {
			t.Errorf("B should have been removed: %+v", out.Variables)
		}
	}
}

func TestRemoveEnvVariableNoopOnMiss(t *testing.T) {
	reg := model.EnvRegistry{Variables: []model.EnvVariable{{Name: "A", Value: "1"}}}
	out := removeEnvVariable(reg, "Z")
	if len(out.Variables) != 1 {
		t.Errorf("len = %d, want 1", len(out.Variables))
	}
}

func TestRemoveEnvPathDropsMatch(t *testing.T) {
	reg := model.EnvRegistry{Paths: []model.EnvPathEntry{
		{Path: `C:\a`}, {Path: `C:\b`}, {Path: `C:\c`},
	}}
	out := removeEnvPath(reg, `C:\b`)
	if len(out.Paths) != 2 {
		t.Fatalf("len = %d, want 2", len(out.Paths))
	}
	for _, p := range out.Paths {
		if p.Path == `C:\b` {
			t.Errorf("path should have been removed: %+v", out.Paths)
		}
	}
}
