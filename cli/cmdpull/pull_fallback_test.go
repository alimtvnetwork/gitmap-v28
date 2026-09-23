package cmdpull

import (
	"testing"
)

func TestShouldFallbackToPullAll(t *testing.T) {
	// When slug or group or all is set, it shouldn't fallback
	optsWithSlug := pullOptions{slug: "my-repo"}
	if ShouldFallbackToPullAll(optsWithSlug) {
		t.Fatalf("expected false when slug is set")
	}

	optsWithGroup := pullOptions{group: "core"}
	if ShouldFallbackToPullAll(optsWithGroup) {
		t.Fatalf("expected false when group is set")
	}

	optsWithAll := pullOptions{all: true}
	if ShouldFallbackToPullAll(optsWithAll) {
		t.Fatalf("expected false when all is already true")
	}
}

func TestCheckEfficientSubcommand_Detection(t *testing.T) {
	if !isEfficientPullSubcmd("all-efficient") {
		t.Fatalf("expected all-efficient to be detected")
	}
	if !isEfficientPullSubcmd("ae") {
		t.Fatalf("expected ae to be detected")
	}
	if !isEfficientPullSubcmd("pae") {
		t.Fatalf("expected pae to be detected")
	}
	if !isEfficientPullSubcmd("pull-ae") {
		t.Fatalf("expected pull-ae to be detected")
	}

	if !isEfficientTableSubcmd("all-efficient-table") {
		t.Fatalf("expected all-efficient-table to be detected")
	}
	if !isEfficientTableSubcmd("paet") {
		t.Fatalf("expected paet to be detected")
	}

	if isEfficientPullSubcmd("unknown-cmd") {
		t.Fatalf("expected false for unknown-cmd")
	}
}
