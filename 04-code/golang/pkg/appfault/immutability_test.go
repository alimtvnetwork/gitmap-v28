package appfault_test

import (
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

func TestAppError_StrictImmutability(t *testing.T) {
	// 1. Create base error
	baseErr := appfault.New(errtype.Validation, "base validation failure")
	if baseErr.StatusCode() != 400 { // Default status for Validation is 400
		t.Fatalf("expected default status 400, got %d", baseErr.StatusCode())
	}

	// 3. Derive variant B via WithContext
	variantB := baseErr.WithContext("userId", "usr-101")
	if !variantB.Context().Has("userId") {
		t.Fatalf("expected variantB to have userId")
	}

	// VERIFY BASE CONTEXT IS UNMUTATED!
	if baseErr.Context().Has("userId") {
		t.Fatalf("VIOLATION: baseErr context was mutated!")
	}
}

func TestAppBuilder_MutableStagingAndImmutableFreeze(t *testing.T) {
	// 1. Build via AppErrorBuilder (AppBuilder)
	builder := appfault.NewAppBuilder(errtype.Unauthorized, "access token expired").
		WithContext("tokenType", "Bearer")

	// 2. Marshaling effect on builder
	jsonBytes, err := builder.MarshalJSON()
	if err != nil {
		t.Fatalf("builder.MarshalJSON failed: %v", err)
	}

	if !strings.Contains(string(jsonBytes), `"Type":10`) {
		t.Fatalf("expected Type 10 in builder json: %s", string(jsonBytes))
	}

	// 3. Unmarshaling effect on builder
	var unmarshaledBuilder appfault.AppBuilder
	if err := unmarshaledBuilder.UnmarshalJSON(jsonBytes); err != nil {
		t.Fatalf("builder.UnmarshalJSON failed: %v", err)
	}

	// 4. Freeze into immutable *AppError
	immutableErr := unmarshaledBuilder.Build()
	if immutableErr.Type() != errtype.Unauthorized {
		t.Fatalf("expected type Unauthorized, got %v", immutableErr.Type())
	}

	// 5. Convert immutable error back to builder via ToBuilder()
	stagedBuilder := immutableErr.ToBuilder().
		SetContext("scope", "admin")

	modifiedErr := stagedBuilder.Build()
	if !modifiedErr.Context().Has("scope") {
		t.Fatalf("expected modified scope context")
	}

	// Original immutable error must remain unchanged!
	if immutableErr.Context().Has("scope") {
		t.Fatalf("VIOLATION: immutableErr was mutated by builder modification!")
	}
}
