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

	registerMoveMergeExamples(reg)
}

func registerMoveMergeExamples(reg map[string]HelpExample) {
	reg["merge-right"] = HelpExample{
		WhatItDoes: "One-way file-level merge that writes only into RIGHT. Files missing on RIGHT are copied from LEFT; conflicts are resolved into RIGHT. LEFT is never modified. If RIGHT originated from a URL it is committed and pushed after the merge.",
		Examples: []Example{
			{
				Invocation: "gitmap merge-right ./local https://github.com/owner/repo -y",
				ExpectedOutput: `[merge-right] resolving RIGHT : https://github.com/owner/repo
[merge-right]   -> folder does not exist; cloning
[merge-right]   -> clone OK
[merge-right] diffing trees ...
[merge-right]   conflict src/app.ts -> took LEFT
[merge-right] committing in https://github.com/owner/repo ...
[merge-right]   commit <sha> "gitmap merge-right from ./local"
[merge-right] pushing https://github.com/owner/repo ...
[merge-right]   push OK
[merge-right] done`,
			},
			{
				Invocation: "gitmap mr ./local https://github.com/owner/repo:develop -y",
				ExpectedOutput: `[merge-right] resolving RIGHT : https://github.com/owner/repo:develop
[merge-right]   -> clone OK
[merge-right] pushing https://github.com/owner/repo ...
[merge-right]   push OK
[merge-right] done`,
			},
		},
	}
	reg["mrg"] = reg["merge-right"]

	reg["merge-left"] = HelpExample{
		WhatItDoes: "One-way file-level merge that writes only into LEFT. Files missing on LEFT are copied from RIGHT; conflicts are resolved into LEFT. RIGHT is never modified. If LEFT originated from a URL it is committed and pushed after the merge.",
		Examples: []Example{
			{
				Invocation: "gitmap merge-left ./local https://github.com/owner/upstream -y",
				ExpectedOutput: `[merge-left] resolving RIGHT : https://github.com/owner/upstream
[merge-left]   -> clone OK
[merge-left] diffing trees ...
[merge-left]   conflict main.go -> took RIGHT
[merge-left] done`,
			},
			{
				Invocation: "gitmap ml ./mine ./theirs -y --prefer-left",
				ExpectedOutput: `[merge-left] diffing trees ...
[merge-left]   conflict app.ts -> took LEFT
[merge-left] done`,
			},
		},
	}
	reg["ml"] = reg["merge-left"]

	reg["merge-both"] = HelpExample{
		WhatItDoes: "Two-way bidirectional merge synchronizing changes between LEFT and RIGHT endpoints.",
		Examples: []Example{
			{
				Invocation: "gitmap merge-both ./repo-a ./repo-b -y",
				ExpectedOutput: `[merge-both] diffing trees ...
[merge-both] synchronized changes between endpoints
[merge-both] done`,
			},
		},
	}
	reg["mb"] = reg["merge-both"]
}
