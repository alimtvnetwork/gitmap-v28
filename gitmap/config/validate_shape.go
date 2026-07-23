package config

// Shape-level validation: catches user-supplied configs whose JSON
// keys are present (so ValidateRawConfig is happy) but carry the
// wrong TYPE or out-of-range numeric values that a permissive
// json.Unmarshal would silently coerce.
//
// Examples the basic raw-key + enum checks would miss:
//
//   - {"defaultMode": 42, ...}   -- number, not string. Unmarshal
//     errors out with a confusing "cannot unmarshal number into
//     Go struct field" message that doesn't tell the user which
//     of their fields is wrong if they have several. This file
//     surfaces the per-key culprit BEFORE Unmarshal runs.
//   - {"outputDir": ""}          -- empty string. Survives the
//     "key present" check and the "is a string" check, but writing
//     scan output to "" silently drops the artifacts in cwd, which
//     is almost never what the user wants. Treated as invalid.
//   - {"dashboardRefresh": -5}   -- negative refresh interval would
//     either hang the TUI or panic depending on the consumer; we
//     reject it up-front so the failure mode is a startup error,
//     not a runtime crash 30 seconds in.
//   - {"excludeDirs": ["", "vendor"]}  -- an empty entry would
//     match every directory in the walk and effectively disable
//     scanning. Reject the whole config rather than auto-pruning.
//   - {"release": {"targets": [{"goos": "linux"}]}}  -- missing
//     goarch. The release builder would later fail with a less
//     contextual error; surface it at config-load time instead.
//
// All shape checks aggregate violations the same way ValidateConfig
// does so users see every problem in one go.

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// expectedRequiredKeyTypes maps each required JSON key to the JSON
// kind ("string" / "number" / "object" / "array") it must hold. We
// model "kind" as a tiny enum local to this file because Go's
// reflect.Kind doesn't map cleanly onto JSON types and we'd rather
// have one obvious switch statement than indirection through reflect.
var expectedRequiredKeyTypes = map[string]jsonKind{
	"defaultMode":   kindString,
	"defaultOutput": kindString,
	"outputDir":     kindString,
}

// jsonKind enumerates the JSON value types we care about. Stringer-
// style String() method exists so error messages can name the
// expected type without a separate lookup table.
type jsonKind int

const (
	kindString jsonKind = iota
	kindNumber
	kindBool
	kindObject
	kindArray
)

// String renders the human-readable JSON type name used in error
// messages. Spelled out (not "str"/"num") because the audience is a
// user reading their own config file.
func (k jsonKind) String() string {
	switch k {
	case kindString:
		return "string"
	case kindNumber:
		return "number"
	case kindBool:
		return "boolean"
	case kindObject:
		return "object"
	case kindArray:
		return "array"
	}

	return "unknown"
}

// ValidateRawShape inspects raw JSON bytes for type mismatches on
// required keys. Runs AFTER ValidateRawConfig has confirmed every
// required key is present, so we can assume the lookup will hit and
// only need to verify the value's JSON kind. Returns a single
// aggregated error listing every type mismatch — never first-fail.
func ValidateRawShape(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		// Same error class ValidateRawConfig already surfaces; return
		// nil here and let the existing path emit the message so we
		// don't double-report on malformed JSON.
		return nil
	}
	violations := collectShapeViolations(raw)
	if len(violations) == 0 {
		return nil
	}

	return fmt.Errorf(
		"config: type mismatch(es):\n  - %s",
		strings.Join(violations, "\n  - "),
	)
}

// collectShapeViolations walks the required-key type table and
// records one violation per key whose JSON value is the wrong kind.
// Centralized so ValidateRawShape stays under the function-line
// budget.
func collectShapeViolations(raw map[string]json.RawMessage) []string {
	violations := make([]string, 0, len(expectedRequiredKeyTypes))
	for key, want := range expectedRequiredKeyTypes {
		val, ok := raw[key]
		if !ok {
			continue // ValidateRawConfig already reported it.
		}
		got := detectKind(val)
		if got == want {
			continue
		}
		violations = append(violations, fmt.Sprintf(
			"%s: expected %s, got %s",
			key, want, got,
		))
	}

	return violations
}

// detectKind classifies a raw JSON value by its first non-space
// byte. Cheap, allocation-free, and sufficient for the coarse
// "is this a string vs a number vs an object" check we need —
// json.Unmarshal would do the same dispatch internally.
func detectKind(raw json.RawMessage) jsonKind {
	for _, b := range raw {
		switch b {
		case ' ', '\t', '\n', '\r':
			continue
		case '"':
			return kindString
		case '{':
			return kindObject
		case '[':
			return kindArray
		case 't', 'f':
			return kindBool
		}

		return kindNumber
	}

	return kindString
}

// ValidateConfigStruct runs the post-unmarshal struct-level checks
// that don't fit the enum-only ValidateConfig contract: empty
// outputDir, negative dashboardRefresh, empty excludeDirs entries,
// and incomplete release targets. Returns a single aggregated error
// so users see every problem at once, mirroring ValidateConfig.
//
// Kept as a separate exported function (rather than folded into
// ValidateConfig) so callers that already hand-build a model.Config
// for tests can opt into the stricter checks without forcing every
// existing test fixture to satisfy them.
func ValidateConfigStruct(cfg model.Config) error {
	violations := make([]string, 0, 4)
	violations = checkOutputDir(violations, cfg)
	violations = checkDashboardRefresh(violations, cfg)
	violations = checkExcludeDirs(violations, cfg)
	violations = checkReleaseTargets(violations, cfg)
	if len(violations) == 0 {
		return nil
	}

	return fmt.Errorf(
		"config: invalid value(s):\n  - %s",
		strings.Join(violations, "\n  - "),
	)
}

// checkOutputDir rejects an explicit empty string. A whitespace-only
// value (e.g. "   ") is also rejected because it's almost certainly
// a typo and writing artifacts to a directory named " " is
// indistinguishable from a config bug.
func checkOutputDir(violations []string, cfg model.Config) []string {
	if len(strings.TrimSpace(cfg.OutputDir)) > 0 {
		return violations
	}

	return append(violations, "outputDir: must be a non-empty path")
}

// checkDashboardRefresh rejects negative values. Zero is allowed —
// callers interpret 0 as "use built-in default" and that contract
// shouldn't be re-policed here.
func checkDashboardRefresh(violations []string, cfg model.Config) []string {
	if cfg.DashboardRefresh >= 0 {
		return violations
	}

	return append(violations, fmt.Sprintf(
		"dashboardRefresh: %d is negative (must be >= 0)",
		cfg.DashboardRefresh,
	))
}

// checkExcludeDirs rejects any entry that's empty after trimming
// whitespace. An empty pattern would match every directory in the
// scan walk and silently disable scanning, which is exactly the
// "limps along" failure mode this validator family exists to
// prevent.
func checkExcludeDirs(violations []string, cfg model.Config) []string {
	for i, entry := range cfg.ExcludeDirs {
		if len(strings.TrimSpace(entry)) > 0 {
			continue
		}
		violations = append(violations, fmt.Sprintf(
			"excludeDirs[%d]: empty entry not allowed", i,
		))
	}

	return violations
}

// checkReleaseTargets rejects targets with missing goos or goarch.
// We do NOT validate the specific values against a known list —
// Go gains new GOOS/GOARCH pairs over time and a hard-coded
// allowlist would block legitimate configs the day after a Go
// release. Presence is the right contract here.
func checkReleaseTargets(violations []string, cfg model.Config) []string {
	for i, target := range cfg.Release.Targets {
		if len(strings.TrimSpace(target.GOOS)) == 0 {
			violations = append(violations, fmt.Sprintf(
				"release.targets[%d].goos: must be non-empty", i,
			))
		}
		if len(strings.TrimSpace(target.GOARCH)) == 0 {
			violations = append(violations, fmt.Sprintf(
				"release.targets[%d].goarch: must be non-empty", i,
			))
		}
	}

	return violations
}
