// Package e2e runs the changelog generator end-to-end against an
// ephemeral git repository and asserts that both regenerated outputs
// (CHANGELOG.md fragment and src/data/changelog.ts fragment) match
// committed golden fixtures byte-for-byte.
//
// This is the highest-value test in the changelog suite: it exercises
// gitlog → group → render → writer in one shot, the same way the
// `make changelog` target does, so any regression in the joint
// Markdown / TypeScript contract fails CI immediately.
//
// Fixtures live alongside this file in `testdata/` so updating them
// is a one-line `cp` away when the format intentionally changes.
package e2e
