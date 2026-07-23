package helptext

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/fixtureversion"
)

// TestValidateStamp_UnstampedPasses confirms the opt-in contract:
// help files without a fixture-stamp comment must validate cleanly
// so pre-existing help pages never regress into failing tests.
func TestValidateStamp_UnstampedPasses(t *testing.T) {
	t.Parallel()
	body := "# clone\n\nUsage: gitmap clone <url>\n"
	if err := ValidateStamp("clone", body, fixtureversion.Expectation{MinGeneration: 1}); err != nil {
		t.Fatalf("unstamped help must pass, got: %v", err)
	}
}

// TestValidateStamp_StaleFails proves the guard fires when a
// stamped file's generation is behind the caller's expectation.
// This is the actionable failure mode the extension was built for.
func TestValidateStamp_StaleFails(t *testing.T) {
	t.Parallel()
	body := "<!-- fixture-stamp: name=version generation=1 min-current=6 for=v6-help -->\n" +
		"# version\n"
	err := ValidateStamp("version", body, fixtureversion.Expectation{
		MinGeneration:    2,
		CurrentVersion:   6,
		RegenerateRecipe: "regenerate helptext/version.md",
	})
	if err == nil {
		t.Fatal("stale stamp must fail validation")
	}
	if !strings.Contains(err.Error(), "generation 1") {
		t.Fatalf("error should surface stale generation, got: %v", err)
	}
}

// TestValidateStamp_FreshPasses is the happy path — a stamped help
// file whose generation meets the caller's floor validates cleanly.
func TestValidateStamp_FreshPasses(t *testing.T) {
	t.Parallel()
	body := "<!-- fixture-stamp: name=version generation=3 min-current=6 for=v6-help -->\n" +
		"# version\n"
	if err := ValidateStamp("version", body, fixtureversion.Expectation{
		MinGeneration:  2,
		CurrentVersion: 6,
	}); err != nil {
		t.Fatalf("fresh stamp should pass, got: %v", err)
	}
}

// TestValidateStamp_MalformedFails guards against a half-written
// stamp silently passing — if the comment prefix is present but
// the marker cannot be parsed, that's a bug in the producer, not
// an unstamped file.
func TestValidateStamp_MalformedFails(t *testing.T) {
	t.Parallel()
	body := "<!-- fixture-stamp: garbage -->\n# version\n"
	err := ValidateStamp("version", body, fixtureversion.Expectation{MinGeneration: 1})
	if err == nil {
		t.Fatal("malformed stamp must fail validation")
	}
}
