// Package enums hosts gitmap's closed-set value types.
//
// Every enum in this tree follows the Go enum pattern from the coding
// guidelines:
//
//	spec/02-coding-guidelines/03-coding-guidelines-spec/03-golang/01-enum-specification/01-enum-pattern.md
//	spec/02-coding-guidelines/03-coding-guidelines-spec/03-golang/01-enum-specification/03-folder-structure.md
//
// (in the coding-guidelines repo, ~/workspace/repos/coding-guidelines).
//
// Pattern rules applied here:
//   - one enum per package; package name ends with the `type` suffix
//     (e.g. ctxmodetype), file is always variant.go
//   - underlying type is byte; zero value is Invalid (always first)
//   - iota, PascalCase variants, documented variants
//   - single unexported variantLabels lookup table; String/Label/IsValid/
//     IsInvalid/Is{Value}/All/ByIndex/Parse/Values/MarshalJSON
//
// Enums migrate here out of cli/constants/ one closed set at a time,
// most self-contained first (spec 243, Wave D).
package enums
