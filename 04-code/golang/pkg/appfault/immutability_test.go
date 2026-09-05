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

	// 2. Derive variant A via WithStatusCode(422)
	variantA := baseErr.WithStatusCode(422)
	if variantA.StatusCode() != 422 {
		t.Fatalf("expected variantA status 422, got %d", variantA.StatusCode())
	}

	// VERIFY BASE IS UNMUTATED!
	if baseErr.StatusCode() != 400 {
		t.Fatalf("VIOLATION: baseErr was mutated! Expected 400, got %d", baseErr.StatusCode())
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

	// 4. Derive variant C via WithCaller
	customCaller := appfault.CallerInfo{File: "auth/service.go", Line: 99, Function: "Test"}
	variantC := baseErr.WithCaller(customCaller)
	if variantC.Caller().Line != 99 {
		t.Fatalf("expected variantC caller line 99, got %d", variantC.Caller().Line)
	}

	// VERIFY BASE CALLER IS UNMUTATED!
	if baseErr.Caller().Line == 99 {
		t.Fatalf("VIOLATION: baseErr caller was mutated!")
	}
}

func TestAppBuilder_MutableStagingAndImmutableFreeze(t *testing.T) {
	// 1. Build via AppErrorBuilder (AppBuilder)
	builder := appfault.NewAppBuilder(errtype.Unauthorized, "access token expired").
		WithStatusCode(401).
		WithContext("tokenType", "Bearer").
		WithCaller(appfault.CallerInfo{File: "gateway/auth.go", Line: 120, Function: "Authenticate"})

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

	if immutableErr.StatusCode() != 401 {
		t.Fatalf("expected statusCode 401, got %d", immutableErr.StatusCode())
	}

	// 5. Convert immutable error back to builder via ToBuilder()
	stagedBuilder := immutableErr.ToBuilder().
		SetStatusCode(403).
		SetContext("scope", "admin")

	modifiedErr := stagedBuilder.Build()
	if modifiedErr.StatusCode() != 403 {
		t.Fatalf("expected modified status 403, got %d", modifiedErr.StatusCode())
	}

	// Original immutable error must remain unchanged!
	if immutableErr.StatusCode() != 401 {
		t.Fatalf("VIOLATION: immutableErr was mutated by builder modification!")
	}
}
