package errtype_test

import (
	"encoding/json"
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

func TestVariation_JSONRoundtrip(t *testing.T) {
	// 1. Marshal as uint16 integer
	v := errtype.Validation
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if string(bytes) != "2" {
		t.Fatalf("expected '2', got %s", string(bytes))
	}

	// 2. Unmarshal from integer
	var restored errtype.Variation
	if err := json.Unmarshal([]byte("3"), &restored); err != nil {
		t.Fatalf("unmarshal integer failed: %v", err)
	}

	if restored != errtype.NotFound {
		t.Fatalf("expected NotFound, got %v", restored)
	}

	// 3. Unmarshal from string name: "Unauthorized"
	var fromString errtype.Variation
	if err := json.Unmarshal([]byte(`"Unauthorized"`), &fromString); err != nil {
		t.Fatalf("unmarshal string failed: %v", err)
	}

	if fromString != errtype.Unauthorized {
		t.Fatalf("expected Unauthorized, got %v", fromString)
	}
}
