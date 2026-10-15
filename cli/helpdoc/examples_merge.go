package helpdoc

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
)

// registerMergeWriteExamples pulls the `merge write` / `mr` help examples
// from cmdmerge's canonical map (single source of truth lives in
// cli/cmdmerge/merge_write_help.go) into the help-example registry.
func registerMergeWriteExamples(reg map[string]HelpExample) {
	for verb, src := range cmdmerge.MergeWriteHelpExamples {
		examples := make([]Example, 0, len(src.Examples))
		for _, e := range src.Examples {
			examples = append(examples, Example{
				Invocation:     e.Invocation,
				ExpectedOutput: e.ExpectedOutput,
			})
		}
		reg[normalizeExampleKey(verb)] = HelpExample{
			WhatItDoes: src.WhatItDoes,
			Examples:   examples,
		}
	}
}
