// Package config — schema validation for user-supplied JSON configs.
//
// Wired into LoadFromFile so the CLI fails fast at startup with a
// clear, multi-line error message when the config file on disk is
// missing required keys or contains invalid enum values. The check
// only runs when the user actually provided a config file — a
// missing file falls back to model.DefaultConfig() which is correct
// by construction and needs no validation.
//
// Two classes of failure are surfaced:
//
//  1. Missing required keys at the top level of the JSON object.
//     Detected by parsing into map[string]json.RawMessage BEFORE
//     unmarshaling into model.Config, since a typed unmarshal
//     silently retains struct defaults for absent keys and the user
//     would never know their config file was broken.
//
//  2. Invalid enum values for defaultMode / defaultOutput. Detected
//     against the typed struct AFTER unmarshal — covers both the
//     "explicit empty string" case (`"defaultMode": ""`) and the
//     "typo" case (`"defaultMode": "htps"`).
//
// Error messages list ALL violations at once (not first-fail) so the
// user can fix a broken config in a single edit cycle. Each line is
// prefixed with the offending JSON key path for grep-ability.
package config

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// requiredConfigKeys are the top-level JSON keys that MUST be
// present in every user-supplied config file. Kept as a sorted slice
// (rather than a set) so the missing-key error message is stable and
// diffable across runs. Add to this list as new fields graduate from
// "optional with default" to "required and explicit".
var requiredConfigKeys = []string{
	"defaultMode",
	"defaultOutput",
	"outputDir",
}

// validModeValues is the closed set of accepted defaultMode values.
// Sourced from constants so the enum stays in lockstep with the rest
// of the codebase — a future ModeSFTP would need adding here too.
var validModeValues = []string{constants.ModeHTTPS, constants.ModeSSH}

// validOutputValues is the closed set of accepted defaultOutput
// values. Same lockstep-with-constants rationale as validModeValues.
var validOutputValues = []string{
	constants.OutputTerminal,
	constants.OutputCSV,
	constants.OutputJSON,
}

// ValidateRawConfig inspects the raw JSON bytes of a user-supplied
// config file and returns a single aggregated error listing every
// missing required key. Called BEFORE typed unmarshal so absent keys
// can be distinguished from defaulted ones.
//
// Returns nil when every required key is present. The caller is
// expected to follow up with ValidateConfig once the typed struct
// is populated, to catch invalid enum values too.
func ValidateRawConfig(data []byte) error {
	var raw map[string]json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {

		return fmt.Errorf("config: invalid JSON: %w", err)
	}
	missing := findMissingKeys(raw, requiredConfigKeys)
	if len(missing) == 0 {

		return nil
	}

	return fmt.Errorf(
		"config: missing required key(s): %s",
		strings.Join(missing, ", "),
	)
}

// ValidateConfig checks the populated struct for invalid enum
// values. Aggregates ALL violations into one error message so users
// can fix a broken config in a single edit instead of fix-rerun-
// fix-rerun. Each violation line names the JSON key, the bad value,
// and the accepted set — copy-pasteable into the config file.
func ValidateConfig(cfg model.Config) error {
	violations := collectEnumViolations(cfg)
	if len(violations) == 0 {

		return nil
	}
	sort.Strings(violations)

	return fmt.Errorf(
		"config: invalid value(s):\n  - %s",
		strings.Join(violations, "\n  - "),
	)
}

// findMissingKeys returns the subset of `required` not present in
// `raw`, preserving the input order so the error message ordering
// is deterministic across runs.
func findMissingKeys(raw map[string]json.RawMessage, required []string) []string {
	missing := make([]string, 0, len(required))
	for _, key := range required {
		_, ok := raw[key]
		if !ok {
			missing = append(missing, key)
		}
	}

	return missing
}

// collectEnumViolations gathers every enum-mismatch from the typed
// config. Split out from ValidateConfig so the per-field check stays
// under the 15-line function budget.
func collectEnumViolations(cfg model.Config) []string {
	violations := make([]string, 0, 2)
	violations = appendEnumViolation(
		violations, "defaultMode", cfg.DefaultMode, validModeValues,
	)
	violations = appendEnumViolation(
		violations, "defaultOutput", cfg.DefaultOutput, validOutputValues,
	)

	return violations
}

// appendEnumViolation appends a formatted violation line if `value`
// is not in `allowed`. Empty values count as violations — an
// explicit `"defaultMode": ""` in the JSON should fail loudly, not
// silently fall back to a default.
func appendEnumViolation(violations []string, key, value string, allowed []string) []string {
	if isAllowedValue(value, allowed) {

		return violations
	}

	return append(violations, fmt.Sprintf(
		"%s: %q is not one of [%s]",
		key, value, strings.Join(allowed, ", "),
	))
}

// isAllowedValue reports whether value is in allowed. Linear scan is
// fine — both enum sets are tiny and constant-sized.
func isAllowedValue(value string, allowed []string) bool {
	for _, candidate := range allowed {
		if value == candidate {

			return true
		}
	}

	return false
}
