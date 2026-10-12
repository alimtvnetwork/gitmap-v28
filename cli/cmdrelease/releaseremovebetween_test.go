package cmdrelease

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

func TestFilterTagsBetweenBounds(t *testing.T) {
	allTags := []string{
		"v4.16.0",
		"v4.160.0",
		"v4.161.0",
		"v4.162.0",
		"v4.162.0-tmp",
		"v4.162.1",
		"v4.162.5",
		"v4.163.0",
		"v4.175.0",
		"v4.189.0",
		"v4.190.0",
		"v4.191.0",
		"v5.0.0",
	}

	lo, err := release.Parse("4.162")
	if err != nil {
		t.Fatalf("Parse lo failed: %v", err)
	}
	hi, err := release.Parse("4.190")
	if err != nil {
		t.Fatalf("Parse hi failed: %v", err)
	}

	// 1. Without pins: intermediate tags (4.162.0-tmp, 4.162.1, 4.162.5, 4.163.0, 4.175.0, 4.189.0) should be deleted
	toDelete, kept := filterTagsBetweenBounds(allTags, lo, hi, nil)

	expectedDelete := map[string]bool{
		"v4.162.0-tmp": true,
		"v4.162.1":     true,
		"v4.162.5":     true,
		"v4.163.0":     true,
		"v4.175.0":     true,
		"v4.189.0":     true,
	}

	for _, d := range toDelete {
		if !expectedDelete[d] {
			t.Errorf("Unexpected tag marked for delete: %s", d)
		}
	}
	if len(toDelete) != len(expectedDelete) {
		t.Errorf("Expected %d tags to delete, got %d: %v", len(expectedDelete), len(toDelete), toDelete)
	}

	// Verify anchors v4.162.0 and v4.190.0 are in kept
	hasLoAnchor := false
	hasHiAnchor := false
	for _, k := range kept {
		if k == "v4.162.0" {
			hasLoAnchor = true
		}
		if k == "v4.190.0" {
			hasHiAnchor = true
		}
	}
	if !hasLoAnchor {
		t.Errorf("Lower bound v4.162.0 was not preserved in kept list!")
	}
	if !hasHiAnchor {
		t.Errorf("Upper bound v4.190.0 was not preserved in kept list!")
	}

	// 2. With pin for v4.175.0: it must be kept
	toDeleteWithPin, keptWithPin := filterTagsBetweenBounds(allTags, lo, hi, []string{"v4.175.0"})
	for _, d := range toDeleteWithPin {
		if d == "v4.175.0" {
			t.Errorf("Pinned tag v4.175.0 was marked for delete!")
		}
	}
	hasPinnedInKept := false
	for _, k := range keptWithPin {
		if k == "v4.175.0" {
			hasPinnedInKept = true
		}
	}
	if !hasPinnedInKept {
		t.Errorf("Pinned tag v4.175.0 not found in kept list!")
	}
}

func TestResolveReleaseBoundsInverted(t *testing.T) {
	lo, hi, err := resolveReleaseBounds("4.190", "4.162")
	if err != nil {
		t.Fatalf("resolveReleaseBounds failed: %v", err)
	}

	if lo.Minor != 162 || hi.Minor != 190 {
		t.Errorf("Expected bounds [4.162 .. 4.190], got [%s .. %s]", lo.String(), hi.String())
	}
}
