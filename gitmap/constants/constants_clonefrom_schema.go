package constants

// JSON-Schema emit surface for `gitmap clone-from --emit-schema=<kind>`.
//
// Split into its own file so the main constants_clonefrom.go stays
// under the per-file line cap (mem://core code-constraints) and so a
// future schema family (e.g. release-report) can add its own constants
// here without churning the parent file.
//
// The flag is a CLI-only convenience: it does NOT touch parse / execute
// surfaces. When --emit-schema is set, clone-from short-circuits BEFORE
// requiring the positional <file> argument so users can run e.g.
// `gitmap clone-from --emit-schema=report` from any directory.

const (
	// FlagCloneFromEmitSchema is the long-form flag; values are
	// validated against EmitSchemaKind* below. Empty string = flag
	// absent = normal clone-from operation.
	FlagCloneFromEmitSchema     = "emit-schema"
	FlagDescCloneFromEmitSchema = "Emit a JSON Schema (draft 2020-12) " +
		"to stdout and exit 0. Kinds: 'report' (clone-from JSON " +
		"report envelope) or 'input' (accepted clone-now scan-record " +
		"input array). Use to validate exported manifests in CI."

	// EmitSchemaKindReport documents the JSON Schema for the
	// .gitmap/clone-from-report-<unixts>.json envelope produced by
	// clonefrom.WriteReportJSON. Tracks CloneFromReportSchemaVersion.
	EmitSchemaKindReport = "report"

	// EmitSchemaKindInput documents the JSON Schema for the array of
	// scan records accepted by `gitmap clone <file>` (the clone-now
	// path). Mirrors clonenow.knownScanFields.
	EmitSchemaKindInput = "input"

	// MsgCloneFromEmitSchemaUnknown is the user-facing error when an
	// unrecognized --emit-schema value is passed. %q = bad value.
	MsgCloneFromEmitSchemaUnknown = "clone-from: --emit-schema %q is not one of 'report', 'input'"

	// $id base for emitted schemas. Versioned URL scheme so a future
	// breaking change publishes under a new path. Not network-fetched
	// by gitmap itself; consumers may resolve it for documentation.
	CloneFromSchemaIDReport = "https://gitmap.dev/schema/clone-from-report-v3.json"
	CloneFromSchemaIDInput  = "https://gitmap.dev/schema/clone-now-input-v1.json"

	// JSONSchemaDialect2020_12 is the canonical $schema URI emitted at
	// the top of every gitmap-produced JSON Schema. Centralized so a
	// future dialect bump is a one-line change.
	JSONSchemaDialect2020_12 = "https://json-schema.org/draft/2020-12/schema"
)
