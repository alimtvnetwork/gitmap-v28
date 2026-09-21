package cmdpull

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestInitPullArrayUIPart1(t *testing.T) {
	if err := InitPullArrayUIPart1(); err != nil {
		t.Fatal(err)
	}
}

func TestResolveSubArrow(t *testing.T) {
	orig := os.Getenv(constants.EnvGlyphs)
	defer os.Setenv(constants.EnvGlyphs, orig)

	os.Setenv(constants.EnvGlyphs, constants.GlyphsSafe)
	if arrow := resolveSubArrow(); arrow != "->" {
		t.Errorf("expected -> in safe mode, got %q", arrow)
	}

	os.Setenv(constants.EnvGlyphs, constants.GlyphsRich)
	if arrow := resolveSubArrow(); arrow != "→" {
		t.Errorf("expected → in rich mode, got %q", arrow)
	}
}
