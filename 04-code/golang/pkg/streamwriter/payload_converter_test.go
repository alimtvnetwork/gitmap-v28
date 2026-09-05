package streamwriter_test

import (
	"context"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/streamwriter"
)

type sampleUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestPayloadConverter_InspectPayload(t *testing.T) {
	if streamwriter.InspectPayload(nil) != streamwriter.PayloadNil {
		t.Fatalf("expected PayloadNil for nil")
	}

	rawBytes := []byte("raw binary buffer")
	if streamwriter.InspectPayload(rawBytes) != streamwriter.PayloadBytes {
		t.Fatalf("expected PayloadBytes for []byte")
	}

	str := "test string"
	if streamwriter.InspectPayload(str) != streamwriter.PayloadString {
		t.Fatalf("expected PayloadString for string")
	}

	appErr := appfault.New(errtype.Validation, "bad payload")
	if streamwriter.InspectPayload(appErr) != streamwriter.PayloadError {
		t.Fatalf("expected PayloadError for *AppError")
	}

	m := map[string]any{"key": "value"}
	if streamwriter.InspectPayload(m) != streamwriter.PayloadMap {
		t.Fatalf("expected PayloadMap for map")
	}

	u := &sampleUser{Name: "Alice", Email: "alice@example.com"}
	if streamwriter.InspectPayload(u) != streamwriter.PayloadStruct {
		t.Fatalf("expected PayloadStruct for *sampleUser")
	}
}

func TestPayloadConverter_ExtractBytes_NoBase64Mangle(t *testing.T) {
	// CRITICAL TEST: Ensure []byte is NOT converted into a base64-encoded string!
	raw := []byte("plain binary payload without base64")
	extracted := streamwriter.ExtractBytes(raw)

	if string(extracted) != "plain binary payload without base64" {
		t.Fatalf("unexpected extracted bytes: %s", string(extracted))
	}

	// Verify that it is NOT base64 encoded
	if strings.Contains(string(extracted), "cGxhaW4") {
		t.Fatalf("FATAL: payload was base64 encoded!")
	}
}

func TestPayloadConverter_ExtractBytes_Types(t *testing.T) {
	// String
	strBytes := streamwriter.ExtractBytes("hello world")
	if string(strBytes) != "hello world" {
		t.Fatalf("unexpected string extraction: %s", string(strBytes))
	}

	// Nil
	if streamwriter.ExtractBytes(nil) != nil {
		t.Fatalf("expected nil for nil payload")
	}

	// AppError
	appErr := appfault.New(errtype.NotFound, "item missing")
	errBytes := streamwriter.ExtractBytes(appErr)
	if !strings.Contains(string(errBytes), "item missing") {
		t.Fatalf("expected AppError message in extracted bytes: %s", string(errBytes))
	}
}

func TestPayloadConverter_ExtractJSONBytes_ValidJSON(t *testing.T) {
	rawJSON := []byte(`{"status":"success","code":200}`)
	extracted, appErr := streamwriter.ExtractJSONBytes(rawJSON)
	if appErr != nil {
		t.Fatalf("ExtractJSONBytes failed: %v", appErr)
	}

	if string(extracted) != string(rawJSON) {
		t.Fatalf("expected valid JSON bytes preserved without double-encoding: %s", string(extracted))
	}
}

func TestNewAnyWriter(t *testing.T) {
	buf := &strings.Builder{}
	writer := streamwriter.NewAnyWriter(streamwriter.WriterOptions[any]{
		Name:        "any-test-writer",
		Destination: &writerAdapter{builder: buf},
	})

	if writer.Name() != "any-test-writer" {
		t.Fatalf("unexpected name: %s", writer.Name())
	}

	appErr := writer.Write(context.Background(), "first task")
	if appErr != nil {
		t.Fatalf("write failed: %v", appErr)
	}

	if !strings.Contains(buf.String(), "[any-test-writer] first task") {
		t.Fatalf("unexpected buffer: %s", buf.String())
	}
}

type writerAdapter struct {
	builder *strings.Builder
}

func (w *writerAdapter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p)
}
