package clonenow

// Public accessors for clone-now's input-validation surface, used
// by sibling packages (e.g. clonefrom's --emit-schema flag) that
// need to enumerate accepted field names without duplicating the
// authoritative knownScanFields map.
//
// Kept in its own file so parse_schema.go stays focused on
// parse-time validation logic and the per-file line cap is
// preserved (mem://core code-constraints).

import "sort"

// KnownScanFields returns the alphabetically-sorted list of JSON
// object keys / CSV header column names accepted by `gitmap clone
// <file>` (clone-now path). Returned slice is a fresh copy: callers
// may mutate it freely without affecting validation.
//
// Stable ordering is part of the contract: the JSON-Schema emitter
// embeds these as enum-of-properties and golden tests pin the
// resulting bytes.
func KnownScanFields() []string {
	out := make([]string, 0, len(knownScanFields))
	for k := range knownScanFields {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}

// RequiredScanURLFields returns the names of the URL-bearing fields
// at least one of which MUST be present on every input row. Mirrors
// validateJSONElement / validateCSVBody. Sorted; fresh copy per call.
func RequiredScanURLFields() []string {
	return []string{"httpsUrl", "sshUrl"}
}
