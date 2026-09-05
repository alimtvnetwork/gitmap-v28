package examples_test

import (
	"bytes"
	"strings"
	"testing"

	"coding-guidelines/common/examples"
)

func TestStreamwriterLoggerExample(t *testing.T) {
	buf := &bytes.Buffer{}
	appErr := examples.RunLoggerExample(buf)
	if appErr != nil {
		t.Fatalf("RunLoggerExample failed: %v", appErr)
	}

	out := buf.String()
	if !strings.Contains(out, "User account authentication succeeded") {
		t.Errorf("expected info message in log output")
	}

	if !strings.Contains(out, "High memory utilization threshold reached") {
		t.Errorf("expected warn message in log output")
	}

	if !strings.Contains(out, "Payment gateway returned unexpected gateway timeout") {
		t.Errorf("expected error message in log output")
	}

	if !strings.Contains(out, "req-tx-8891") {
		t.Errorf("expected traceId in log output")
	}

	if !strings.Contains(out, "[audit-api-writer]") {
		t.Errorf("expected audit writer output")
	}
}

func TestStreamwriterJsonExample(t *testing.T) {
	buf := &bytes.Buffer{}
	appErr := examples.RunJsonExample(buf)
	if appErr != nil {
		t.Fatalf("RunJsonExample failed: %v", appErr)
	}

	out := buf.String()
	if !strings.Contains(out, "--- Pretty JSON Output ---") || !strings.Contains(out, "alim.karim") {
		t.Errorf("expected pretty JSON output")
	}

	if !strings.Contains(out, "--- Compact JSON Output ---") {
		t.Errorf("expected compact JSON output")
	}

	if !strings.Contains(out, "Unmarshaled Account: sarah.connor") {
		t.Errorf("expected unmarshaled account output")
	}

	if !strings.Contains(out, "Casted Public Profile: alim.karim [acc-901]") {
		t.Errorf("expected casted public profile output")
	}

	if !strings.Contains(out, "Extended JsonPayloadResult: alim.karim | IsValid: true | StatusCode: 200") {
		t.Errorf("expected extended JsonPayloadResult output")
	}

	if !strings.Contains(out, "Scoped Factory Result: alim@riseup.asia") {
		t.Errorf("expected scoped factory result output")
	}
}

func TestStreamwriterStreamerExample(t *testing.T) {
	buf := &bytes.Buffer{}
	appErr := examples.RunStreamerExample(buf)
	if appErr != nil {
		t.Fatalf("RunStreamerExample failed: %v", appErr)
	}

	out := buf.String()
	if !strings.Contains(out, "ord-1001") {
		t.Errorf("expected individual stream event")
	}

	if !strings.Contains(out, "ord-concurrent-") {
		t.Errorf("expected concurrent stream events")
	}

	if !strings.Contains(out, "=== BEGIN BATCH TRANSACTION ===") {
		t.Errorf("expected atomic batch start")
	}

	if !strings.Contains(out, "=== COMMIT BATCH TRANSACTION ===") {
		t.Errorf("expected atomic batch end")
	}

	if !strings.Contains(out, "[batch-writer][swapped] Message sent via runtime swapped method") {
		t.Errorf("expected hot-swapped write output")
	}
}

func TestDemonstrateAdvancedFileAndPayloadIntelligence(t *testing.T) {
	buf := &bytes.Buffer{}
	appErr := examples.DemonstrateAdvancedFileAndPayloadIntelligence(buf)
	if appErr != nil {
		t.Fatalf("DemonstrateAdvancedFileAndPayloadIntelligence failed: %v", appErr)
	}

	out := buf.String()
	if !strings.Contains(out, "Direct binary stream without Base64 encoding") {
		t.Errorf("expected direct binary stream in output")
	}

	if !strings.Contains(out, "❌ ERROR [Validation:2]") {
		t.Errorf("expected validation banner in output")
	}

	if !strings.Contains(out, `"Function": "HandleRegistration"`) {
		t.Errorf("expected caller function in json output")
	}

	if !strings.Contains(out, `caller="auth/service.go:55 (HandleRegistration)"`) {
		t.Errorf("expected caller in text log output")
	}

	if !strings.Contains(out, "Read") || !strings.Contains(out, "chunk(s)") {
		t.Errorf("expected chunked read output")
	}
}
