package appfault_test

import (
	"errors"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

func createSampleAppError() (*appfault.AppError, error) {
	rawErr := errors.New("underlying socket closed")
	orig := appfault.Wrap(errtype.Network, rawErr, "dial timeout").
		WithOp("net.dial").
		WithStatusCode(504).
		WithContext("port", 8080)

	return orig, rawErr
}

func assertRestoredMatches(t *testing.T, restored *appfault.AppError, orig *appfault.AppError, rawErr error) {
	if restored.Type() != orig.Type() || restored.StatusCode() != 504 {
		t.Fatalf("JSON restore mismatch: %+v", restored)
	}

	if restored.Cause() == nil || restored.Cause().Error() != rawErr.Error() {
		t.Fatalf("expected cause '%s', got '%v'", rawErr.Error(), restored.Cause())
	}
}

func TestAppErrorSerializationRoundtrip(t *testing.T) {
	orig, rawErr := createSampleAppError()
	restored, err := appfault.FromJSON([]byte(orig.ToJSONString()))
	if err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	assertRestoredMatches(t, restored, orig, rawErr)
	if restored.Caller().IsEmpty() {
		t.Fatal("expected non-empty caller")
	}
}

func TestAppErrorYAMLSerialization(t *testing.T) {
	appErr := appfault.New(errtype.Unauthorized, "invalid credentials").WithStatusCode(401)
	yamlStr := appErr.ToYAMLString()

	if !strings.Contains(yamlStr, "Type: 10") || !strings.Contains(yamlStr, "StatusCode: 401") {
		t.Fatalf("unexpected YAML string: %s", yamlStr)
	}
}

type sampleUser struct {
	Name string
	Age  int
}

func TestGenericJSONHelpers(t *testing.T) {
	u := sampleUser{Name: "Alice", Age: 30}
	jsonStr, appErr := appfault.SerializeToJSONString(u)
	if appErr != nil || len(jsonStr) == 0 {
		t.Fatalf("failed to serialize sampleUser: %v", appErr)
	}

	restored, err := appfault.DeserializeFromJSONString[sampleUser](jsonStr)
	if err != nil || restored.Name != "Alice" || restored.Age != 30 {
		t.Fatalf("failed to deserialize sampleUser: %v", err)
	}
}

func TestAppErrorCompilation(t *testing.T) {
	appErr := appfault.New(errtype.Database, "query failed").WithOp("db.exec")
	compiled := appErr.Compile()

	if !strings.Contains(compiled, "Database") || !strings.Contains(compiled, "query failed") {
		t.Fatalf("unexpected compiled output: %s", compiled)
	}
}
