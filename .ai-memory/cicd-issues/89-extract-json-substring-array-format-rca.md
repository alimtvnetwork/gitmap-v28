# Root Cause Analysis: Extract JSON Substring Array Format Truncation

**Issue ID:** CI-89
**Target Workflow:** CI (`#36333276195`)
**Target Job:** `Full Suite Guard`
**Target Step:** `Run full test suite`
**Date:** 2026-09-28
**Status:** Resolved

---

## 1. Reproduction

Running `go test ./...` in `cli/cmdupdate` triggered:
```text
--- FAIL: TestParseFleetUpdateTelemetry (0.00s)
    --- FAIL: TestParseFleetUpdateTelemetry/array_format (0.00s)
        update_fleet_test.go:290: expected 2 updated, got 0
FAIL
FAIL github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate 0.081s
```

Input payload:
```json
[{"name": "gitmap", "status": "updated"}, {"name": "node", "status": "ok"}]
```

---

## 2. Root Cause Analysis

In `cli/cmdupdate/update_fleet.go`, `extractJSONSubstring` was implemented strictly searching for `{` and `}` delimiters:
```go
func extractJSONSubstring(s string) string {
    start := strings.Index(s, "{")
    end := strings.LastIndex(s, "}")
    if start >= 0 && end > start {
        return s[start : end+1]
    }
    return s
}
```

When an array of telemetry items was passed (starting with `[` and ending with `]`), `extractJSONSubstring` found the inner `{` and `}` and sliced out:
`{"name": "gitmap", "status": "updated"}, {"name": "node", "status": "ok"}`

This stripped the enclosing array brackets `[` and `]`, making it invalid JSON. Consequently:
- `json.Unmarshal([]byte(jsonStr), &parsed)` failed.
- `json.Unmarshal([]byte(jsonStr), &items)` failed with syntax error.
- Telemetry fell back to `buildFallbackTelemetry`, resulting in `Updated` slice length of 0 instead of 2.

---

## 3. Code Fix

Refactored `extractJSONSubstring` and introduced helper `extractDelimitedRange` adhering to <= 8 line function guidelines:
1. Check `startArr := strings.Index(s, "[")` and `startObj := strings.Index(s, "{")`.
2. If `startArr >= 0` and it appears before `startObj` (or no `startObj`), extract between `[` and `]`.
3. Otherwise, if `startObj >= 0`, extract between `{` and `}`.
4. Preserved both single JSON objects and JSON array payloads without stripping root brackets.

---

## 4. Prevention

- Maintain unit test coverage for both object and array formats in `cli/cmdupdate/update_fleet_test.go`.
- Avoid premature substring trimming when the input string is already valid JSON.
