package cmdagy

import (
	"testing"
)

func TestAgyRprpSubcommandNormalization(t *testing.T) {
	aliases := []string{"rprp", "rapwrp", "read-all-projects-with-read-prompts", "read-all-with-prompts"}
	for _, alias := range aliases {
		if !isRprpAlias(alias) {
			t.Errorf("expected alias %q to be recognized", alias)
		}
		if norm := normalizeAgySubcommand(alias); norm != "read-all-projects-with-read-prompts" {
			t.Errorf("expected %q to normalize to read-all-projects-with-read-prompts, got %q", alias, norm)
		}
	}
}
