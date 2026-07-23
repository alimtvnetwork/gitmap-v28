package message

import "github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/profile"

// Inputs collects everything the pipeline needs to build one final
// commit message. All fields are inputs; nothing is mutated.
type Inputs struct {
	OriginalMessage string
	FunctionIntel   string // pre-rendered §6.3 block, may be ""
	Resolved        profile.Resolved
	// PickIndex returns a value in [0,n). Inject deterministic picker
	// from tests; production wires a per-run PRNG seed (spec §3.4).
	PickIndex func(n int) int
}

// Result holds the built message and a flag for the §6.1 step-6 empty
// check; callers map IsEmpty=true to SkipReason EmptyAfterMessageRules.
type Result struct {
	Message string
	IsEmpty bool
}
