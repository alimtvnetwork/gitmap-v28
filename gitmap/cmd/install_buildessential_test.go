package cmd

import (
	"strings"
	"testing"
)

func TestBuildEssentialPackagesExclusion(t *testing.T) {
	pkgs := getBuildEssentialPackages()

	for _, p := range pkgs {
		if strings.Contains(p, "libpcre3-dev") {
			t.Errorf("CRITICAL VIOLATION: libpcre3-dev must NOT be included in buildEssentialPackages (fails on Ubuntu 24.04+)")
		}
	}
}

func TestBuildEssentialAliases(t *testing.T) {
	aliases := []string{"build-essential", "buildessential", "be", "ubuntu-common", "ub-common"}
	for _, a := range aliases {
		if !isBuildEssentialAlias(a) {
			t.Errorf("Expected %q to be recognized as build-essential alias", a)
		}
	}

	if isBuildEssentialAlias("non-existent-tool") {
		t.Errorf("Expected non-existent tool to NOT be recognized as build-essential alias")
	}
}

func TestBuildEssentialProfileTreeResolution(t *testing.T) {
	profile, ok := resolveProfileTree("build-essential")
	if !ok {
		t.Fatalf("Expected profile 'build-essential' to resolve successfully")
	}

	if profile.Name != "build-essential" {
		t.Errorf("Expected profile name 'build-essential', got %q", profile.Name)
	}

	if len(profile.Tools) == 0 {
		t.Errorf("Expected non-empty tools list for build-essential profile")
	}
}

func TestBuildEssentialDryRun(t *testing.T) {
	opts := installOptions{
		Tool:   "build-essential",
		DryRun: true,
	}

	if err := runInstallBuildEssential(opts); err != nil {
		t.Errorf("Expected dry-run to succeed with nil error, got %v", err)
	}
}
