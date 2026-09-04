package appfaults_test

import (
	"errors"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
	"coding-guidelines/common/pkg/errtype"
)

func TestCollectionImmutableConcat(t *testing.T) {
	c1 := appfaults.New().AddTypeMsg(errtype.Validation, "err1")
	c2 := c1.Concat(appfault.New(errtype.Database, "err2"))

	if c1.Count() != 1 || c2.Count() != 2 {
		t.Fatalf("expected c1=1, c2=2 (immutable), got c1=%d, c2=%d", c1.Count(), c2.Count())
	}

	merged := c1.ConcatNew(c2)
	if merged.Count() != 3 {
		t.Fatalf("expected 3 merged items, got %d", merged.Count())
	}
}

func TestCollectionConcatErrors(t *testing.T) {
	c := appfaults.New()
	concatenated := c.ConcatErrors(errtype.IO, errors.New("e1"), errors.New("e2"))

	if c.Count() != 0 || concatenated.Count() != 2 {
		t.Fatalf("expected immutable concat with 2 items, got orig=%d, new=%d", c.Count(), concatenated.Count())
	}
}

func TestCollectionCompilation(t *testing.T) {
	c := appfaults.New().AddTypeMsg(errtype.Network, "socket closed")
	compiled := c.Compile()

	if !strings.Contains(compiled, "Network") || !strings.Contains(compiled, "socket closed") {
		t.Fatalf("unexpected collection compile: %s", compiled)
	}
}
