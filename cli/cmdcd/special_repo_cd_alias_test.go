package cmdcd

import (
	"testing"
)

func TestIsSpecialRepoCDAlias(t *testing.T) {
	if !isSpecialRepoCDAlias("rs") || !isSpecialRepoCDAlias("repo-secrets") {
		t.Fatalf("expected rs/repo-secrets to be special repo CD alias")
	}
	if !isSpecialRepoCDAlias("rc") || !isSpecialRepoCDAlias("repo-cache") || !isSpecialRepoCDAlias("repo-storage") {
		t.Fatalf("expected rc/repo-cache/repo-storage to be special repo CD alias")
	}
	if isSpecialRepoCDAlias("other-repo") {
		t.Fatalf("expected other-repo not to be special repo CD alias")
	}
}
