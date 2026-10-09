// Package cmdspec implements the `gitmap spec` namespace (spec 252):
// concurrency-safe spec-number issuance for the current repository.
//
// `gitmap spec next` issues the next free spec number and prints it as a
// bare value on stdout (script-friendly); `--json` prints the issuance
// record instead. The claim is atomic (INSERT OR IGNORE over a
// UNIQUE(number) column), so parallel callers never receive the same
// number.
package cmdspec
