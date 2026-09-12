package logging

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestNewLoggerDefaults(t *testing.T) {
	l := NewLogger(nil, "test", false)
	if l.w != os.Stderr {
		t.Fatalf("expected stderr default, got %v", l.w)
	}

	if l.enabled {
		t.Fatal("expected logger to be disabled")
	}
}

func TestLoggerDisabledNoOp(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, "test", false)
	l.Log(LevelInfo, "ignored", nil)
	if buf.Len() > 0 {
		t.Fatalf("expected 0 bytes written when disabled, got %d", buf.Len())
	}
}

func TestLoggerEnabledLevels(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, "mycmd", true)

	l.Info("info-msg", map[string]interface{}{"k1": "v1"})
	l.Warn("warn-msg", nil)
	l.Error("err-msg", map[string]interface{}{"code": 42})

	dec := json.NewDecoder(&buf)
	var entries []Entry
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			t.Fatalf("decode entry: %v", err)
		}

		entries = append(entries, e)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	if entries[0].Level != LevelInfo || entries[0].Message != "info-msg" {
		t.Fatalf("unexpected entry 0: %+v", entries[0])
	}
}
