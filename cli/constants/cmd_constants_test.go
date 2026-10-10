package constants

import "testing"

// topLevelCmds enumerates every top-level Cmd* constant exposed to the CLI
// dispatcher. Entries marked with the `// gitmap:cmd skip` comment in
// constants_cli.go (subcommand verbs like "create" / "add" that are reused
// across subcommand groups) are intentionally omitted — duplicates of those
// values are expected and safe.
//
// When you add or remove a top-level Cmd* constant in constants_cli.go,
// update this slice. CI enforces parity via TestTopLevelCmd*.
func topLevelCmds() map[string]string {
	m := make(map[string]string, 600)
	for k, v := range topLevelCmdsPart1() {
		m[k] = v
	}
	for k, v := range topLevelCmdsPart2() {
		m[k] = v
	}
	return m
}

// TestTopLevelCmdConstantsAreUnique asserts that every top-level Cmd*
// constant has a distinct value, so CI rejects accidental redeclarations
// or value collisions (e.g. two constants both equal to "cd") before they
// reach the runtime dispatcher.
func TestTopLevelCmdConstantsAreUnique(t *testing.T) {
	seen := make(map[string]string, len(topLevelCmds()))
	for name, value := range topLevelCmds() {
		if prev, exists := seen[value]; exists {
			t.Errorf("duplicate top-level Cmd constant value %q: %s collides with %s", value, name, prev)
			continue
		}

		seen[value] = name
	}
}

// TestTopLevelCmdAliasesAreUnique asserts that every short alias (any
// top-level Cmd* value of length <= 2) is unique across the entire CLI
// surface. A future CmdFooAlias = "ls" would collide with CmdListAlias and
// be rejected here. Long-form command names are covered by the broader
// TestTopLevelCmdConstantsAreUnique check above; this test focuses
// specifically on the short-alias namespace where collisions are easiest
// to introduce by accident and hardest to spot in code review.
func TestTopLevelCmdAliasesAreUnique(t *testing.T) {
	const maxAliasLen = 2
	seen := make(map[string]string)
	for name, value := range topLevelCmds() {
		if len(value) == 0 || len(value) > maxAliasLen {
			continue
		}

		if prev, exists := seen[value]; exists {
			t.Errorf("duplicate short alias %q: %s collides with %s", value, name, prev)
			continue
		}

		seen[value] = name
	}
}
