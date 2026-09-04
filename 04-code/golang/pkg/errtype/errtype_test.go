package errtype_test

import (
	"testing"

	"coding-guidelines/common/pkg/errtype"
)

const CustomUserError errtype.Variation = 1001

func TestStandardVariations(t *testing.T) {
	if errtype.None.HasError() || !errtype.None.IsNoError() {
		t.Fatal("expected None to indicate no error")
	}

	if !errtype.Validation.HasError() || errtype.Validation.Code() != 2 {
		t.Fatalf("expected Validation to have code 2, got %d", errtype.Validation.Code())
	}

	if errtype.Database.String() != "Database" {
		t.Fatalf("expected Database name, got %s", errtype.Database.String())
	}
}

func TestCustomVariationExtension(t *testing.T) {
	if !CustomUserError.HasError() || CustomUserError.Code() != 1001 {
		t.Fatalf("expected custom error code 1001, got %d", CustomUserError.Code())
	}

	if CustomUserError.String() != "Custom(1001)" {
		t.Fatalf("expected Custom(1001), got %s", CustomUserError.String())
	}
}
