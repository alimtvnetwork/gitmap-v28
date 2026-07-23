package stablejson

// Internal writers for the public encoders in stablejson.go. Split
// into a sibling file purely to keep stablejson.go under the
// 200-line code-style budget — the package surface (WriteArray,
// WriteArrayIndent, WriteJSONLines, Field) all lives in stablejson.go
// so godoc and import-time discovery stay obvious.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// writeArrayPretty emits the multi-line indented form. Split out so
// WriteArrayIndent stays under the 15-line code-style budget and
// the minified path doesn't pay for branch noise on every line.
func writeArrayPretty(w io.Writer, buf *bytes.Buffer, items [][]Field, indent string) error {
	buf.WriteString("[\n")
	for i, obj := range items {
		if err := writeObject(buf, obj, indent); err != nil {
			return err
		}
		if i < len(items)-1 {
			buf.WriteString(",\n")
		} else {
			buf.WriteString("\n")
		}
	}
	buf.WriteString("]\n")
	_, err := w.Write(buf.Bytes())

	return err
}

// writeArrayMinified emits `[{...},{...}]\n` on a single line. The
// per-object encoder is the same one JSONL uses, ensuring a value
// rendered minified inside an array is byte-identical to the same
// value rendered as one JSONL line — important for hash-based
// integrity checks that re-encode and compare.
func writeArrayMinified(w io.Writer, buf *bytes.Buffer, items [][]Field) error {
	buf.WriteByte('[')
	for i, obj := range items {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := writeCompactObject(buf, obj); err != nil {
			return err
		}
	}
	buf.WriteString("]\n")
	_, err := w.Write(buf.Bytes())

	return err
}

// writeCompactObject writes a single `{"k":v,"k2":v2}` block (no
// whitespace between tokens) into buf. Key order follows the slice.
// Each value is JSON-marshaled in isolation so a malformed value
// fails the WHOLE call rather than emitting half a corrupt line —
// critical for JSONL because a half-written line would desync every
// downstream parser that splits on `\n`.
func writeCompactObject(buf *bytes.Buffer, fields []Field) error {
	buf.WriteByte('{')
	for i, f := range fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := writeKeyValue(buf, f, ""); err != nil {
			return err
		}
	}
	buf.WriteByte('}')

	return nil
}

// writeObject writes a single `{ ... }` block into buf using
// caller-controlled indentation. The outer brace sits at one level
// of `indent`; each key/value line sits at two levels. Keys appear
// in the exact order given. Each value is JSON-marshaled in
// isolation so a malformed value fails the WHOLE call rather than
// emitting half a corrupt object.
func writeObject(buf *bytes.Buffer, fields []Field, indent string) error {
	outer := indent
	inner := indent + indent
	buf.WriteString(outer + "{\n")
	for i, f := range fields {
		buf.WriteString(inner)
		if err := writeKeyValue(buf, f, " "); err != nil {
			return err
		}
		if i < len(fields)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.WriteString(outer + "}")

	return nil
}

// writeKeyValue marshals one Field as `"key":value` (when
// `colonSpace==""`) or `"key": value` (when `colonSpace==" "`) and
// appends to buf. Centralized so the compact and pretty paths share
// identical key/value JSON marshaling — a malformed value triggers
// the same wrapped error from both call sites.
func writeKeyValue(buf *bytes.Buffer, f Field, colonSpace string) error {
	keyBytes, err := json.Marshal(f.Key)
	if err != nil {
		return fmt.Errorf("stablejson: encode key %q: %w", f.Key, err)
	}
	buf.Write(keyBytes)
	buf.WriteByte(':')
	buf.WriteString(colonSpace)
	valBytes, err := json.Marshal(f.Value)
	if err != nil {
		return fmt.Errorf("stablejson: encode value for key %q: %w", f.Key, err)
	}
	buf.Write(valBytes)

	return nil
}
