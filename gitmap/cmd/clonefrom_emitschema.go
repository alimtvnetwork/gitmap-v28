package cmd

// Handler for `gitmap clone-from --emit-schema=<kind>`. Split out
// from clonefrom.go so the dispatcher file stays under the 200-line
// per-file cap (mem://core code-constraints).
//
// Exit-code mapping mirrors the rest of the clone-from surface:
//
//   0 — schema written successfully
//   1 — write to stdout failed (broken pipe, disk full, …)
//   2 — unknown --emit-schema kind (CLI-usage error class, same
//       bucket as a missing <file> positional)

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runCloneFromEmitSchema writes the requested JSON Schema to stdout
// and exits. Errors are routed to stderr per the standard CLI split
// (data on stdout so users can pipe `gitmap clone-from
// --emit-schema=report > schema.json` cleanly).
func runCloneFromEmitSchema(kind string) {
	body, err := clonefrom.EmitSchema(kind)
	if err != nil {
		cliexit.Fail(constants.CmdCloneFrom, "emit-schema", kind, err, 2)
	}
	if _, err := os.Stdout.Write(body); err != nil {
		cliexit.Fail(constants.CmdCloneFrom, "write-stdout", "emit-schema", err, 1)
	}
}
