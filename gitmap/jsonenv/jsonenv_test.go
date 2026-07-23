package jsonenv

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteOKEnvelope(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteOK(&buf, "scan", map[string]int{"n": 3}); err != nil {
		t.Fatal(err)
	}
	var env Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if !env.OK || env.Schema != Schema || env.Command != "scan" {
		t.Fatalf("bad envelope: %+v", env)
	}
}

func TestWriteErrIncludesMessage(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteErr(&buf, "probe", "boom", nil)
	if !strings.Contains(buf.String(), `"error": "boom"`) {
		t.Fatalf("missing error field: %s", buf.String())
	}
}
