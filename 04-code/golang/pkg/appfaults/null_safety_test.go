package appfaults_test

import (
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
	"coding-guidelines/common/pkg/errtype"
)

func TestCollection_NilReceiverSafety(t *testing.T) {
	var nilColl *appfaults.Collection

	if !nilColl.IsNull() {
		t.Fatal("expected IsNull() to be true on nil *Collection")
	}

	if !nilColl.IsEmpty() {
		t.Fatal("expected IsEmpty() to be true on nil *Collection")
	}

	if !nilColl.HasZero() {
		t.Fatal("expected HasZero() to be true on nil *Collection")
	}

	if !nilColl.IsZero() {
		t.Fatal("expected IsZero() to be true on nil *Collection")
	}

	if !nilColl.HasNull() {
		t.Fatal("expected HasNull() to be true on nil *Collection")
	}

	if nilColl.HasError() {
		t.Fatal("expected HasError() to be false on nil *Collection")
	}

	if nilColl.Count() != 0 {
		t.Fatalf("expected Count() == 0 on nil, got %d", nilColl.Count())
	}

	// Clone on nil collection
	cloned := nilColl.Clone()
	if cloned == nil || cloned.Count() != 0 {
		t.Fatal("expected cloned nil collection to produce empty collection")
	}

	// Concat on nil collection
	err := appfault.New(errtype.Validation, "bad input")
	concatenated := nilColl.Concat(err)
	if concatenated.Count() != 1 {
		t.Fatalf("expected Concat to have 1 error, got %d", concatenated.Count())
	}

	// ConcatNew on nil collection
	combined := nilColl.ConcatNew(nilColl, nilColl)
	if combined.Count() != 0 {
		t.Fatalf("expected 0 count from ConcatNew on nil, got %d", combined.Count())
	}
}
