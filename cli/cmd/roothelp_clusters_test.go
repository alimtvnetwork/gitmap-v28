package cmd

import (
	"testing"
)

func TestGetSemanticClusterSections_Count(t *testing.T) {
	sections := GetSemanticClusterSections()
	if len(sections) != 5 {
		t.Fatalf("expected exactly 5 semantic clusters, got %d", len(sections))
	}
}

func TestBuildSemanticClustersMenu_Title(t *testing.T) {
	menu := BuildSemanticClustersMenu()
	if menu.Title != "gitmap — 5 Semantic Clusters" {
		t.Errorf("unexpected menu title: %q", menu.Title)
	}
	if len(menu.Sections) != 5 {
		t.Errorf("expected 5 cluster sections in menu, got %d", len(menu.Sections))
	}
}

func TestBuildSemanticClustersMenu_FooterAndTips(t *testing.T) {
	menu := BuildSemanticClustersMenu()
	if len(menu.FooterFlags) == 0 {
		t.Errorf("expected non-empty footer flags")
	}
	if len(menu.Tips) == 0 {
		t.Errorf("expected non-empty tips")
	}
}

func TestIsSemanticClusterName_Valid(t *testing.T) {
	valid := []string{"core", "release", "fleet", "ai", "system", "clusters"}
	for _, name := range valid {
		if !IsSemanticClusterName(name) {
			t.Errorf("expected %q to be recognized as semantic cluster", name)
		}
	}
}

func TestIsSemanticClusterName_Invalid(t *testing.T) {
	if IsSemanticClusterName("foobar") {
		t.Errorf("expected 'foobar' not to be recognized as semantic cluster")
	}
}
