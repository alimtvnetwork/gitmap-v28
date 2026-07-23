package clonenow

// JSON half of the clone-now schema validator. Split out of
// parse_schema.go so each file stays under the 200-line code-style
// budget. Shared state (knownScanFields, knownFieldList) lives in
// parse_schema.go which both halves import implicitly via the
// package namespace.

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// validateJSONSchema ensures the JSON input is an array of objects
// whose keys are all known ScanRecord field names and where every
// object carries at least one URL. Two-stage decode (array first,
// then per-element object) so a non-object element is reported with
// its 1-based row number instead of a raw decoder offset that the
// user has to hand-translate to a row.
func validateJSONSchema(data []byte) error {
	var elems []json.RawMessage
	if err := json.Unmarshal(data, &elems); err != nil {
		return fmt.Errorf(constants.ErrCloneNowJSONShape, err)
	}
	for i, raw := range elems {
		if err := validateJSONElement(i, raw); err != nil {
			return err
		}
	}

	return nil
}

// validateJSONElement decodes one array element as an object and
// runs the per-row checks. A non-object element (string, number,
// null, nested array) is reported with its 1-based index and the
// observed JSON kind so the user can jump straight to the line.
func validateJSONElement(i int, raw json.RawMessage) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return fmt.Errorf(constants.ErrCloneNowJSONRowNotObject,
			i+1, jsonKind(raw))
	}

	return validateJSONRow(i, obj)
}

// validateJSONRow checks one decoded object: every key must be in
// knownScanFields and the row must carry at least one URL. The row
// index is 1-based in the error so it matches what a human reading
// the JSON file would count.
func validateJSONRow(i int, obj map[string]json.RawMessage) error {
	for k := range obj {
		if !knownScanFields[k] {
			return fmt.Errorf(constants.ErrCloneNowUnknownJSONField,
				i+1, k, knownFieldList())
		}
	}
	if !hasJSONURL(obj) {
		return fmt.Errorf(constants.ErrCloneNowMissingURL, i+1)
	}

	return nil
}

// jsonKind returns a short human label for a RawMessage's top-level
// JSON kind, used in error messages so users see "got string" rather
// than a parser offset. Falls back to "invalid" for malformed bytes.
func jsonKind(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if len(s) == 0 {
		return "empty"
	}
	switch s[0] {
	case '"':
		return "string"
	case '[':
		return "array"
	case '{':
		return "object"
	case 't', 'f':
		return "boolean"
	case 'n':
		return "null"
	}
	if (s[0] >= '0' && s[0] <= '9') || s[0] == '-' {
		return "number"
	}

	return "invalid"
}

// hasJSONURL reports whether the row carries a non-empty httpsUrl
// or sshUrl. Empty strings ("") count as missing -- the executor
// would skip the row anyway, and a clear pre-flight error is more
// useful than a silent drop.
func hasJSONURL(obj map[string]json.RawMessage) bool {
	return jsonStringNonEmpty(obj["httpsUrl"]) || jsonStringNonEmpty(obj["sshUrl"])
}

// jsonStringNonEmpty decodes a RawMessage as a string and reports
// whether the result is non-empty. Non-string values are treated as
// empty so a typo like `"httpsUrl": null` doesn't pass the URL gate.
func jsonStringNonEmpty(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return false
	}

	return len(strings.TrimSpace(s)) > 0
}
