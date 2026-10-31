package helpdoc

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
)

// registerMergeWriteExamples pulls merge write and related examples into registry.
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
	registerCommitInExamples(reg)
	registerCommitTransferExamples(reg)
	registerHelpAndSettingsExamples(reg)
}

func registerMoveMergeExamples(reg map[string]HelpExample) {
	reg["merge"] = buildMergeRootExample()
	reg["merge-right"] = buildMergeRightExample()
	reg["mrg"] = reg["merge-right"]
	reg["merge-left"] = buildMergeLeftExample()
	reg["ml"] = reg["merge-left"]
	reg["merge-both"] = buildMergeBothExample()
	reg["mb"] = reg["merge-both"]
}

func buildMergeRootExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Multi-repository merge suite combining one-way file merges (merge-right, merge-left), two-way sync (merge-both), chronological commit replay (merge-in), and pull requests (merge write).",
		Examples: []Example{
			{
				Invocation:     "gitmap merge right ./local https://github.com/owner/repo -y",
				ExpectedOutput: "[merge-right] diffing trees ...\n[merge-right] commit <sha> \"gitmap merge-right from ./local\"\n[merge-right] done",
			},
			{
				Invocation:     "gitmap merge in ./canonical all --default",
				ExpectedOutput: "[commit-in] replaying commits from sibling repos ...\n[commit-in] appended 14 commits\n[commit-in] done",
			},
		},
	}
}

func buildMergeRightExample() HelpExample {
	return HelpExample{
		WhatItDoes: "One-way file-level merge that writes only into RIGHT. Files missing on RIGHT are copied from LEFT; conflicts are resolved into RIGHT. LEFT is never modified. If RIGHT originated from a URL it is committed and pushed after the merge.",
		Examples: []Example{
			{
				Invocation:     "gitmap merge-right ./local https://github.com/owner/repo -y",
				ExpectedOutput: "[merge-right] resolving RIGHT : https://github.com/owner/repo\n[merge-right]   -> clone OK\n[merge-right] diffing trees ...\n[merge-right]   conflict src/app.ts -> took LEFT\n[merge-right] commit <sha> \"gitmap merge-right from ./local\"\n[merge-right] push OK\n[merge-right] done",
			},
			{
				Invocation:     "gitmap mr ./local https://github.com/owner/repo:develop -y",
				ExpectedOutput: "[merge-right] resolving RIGHT : https://github.com/owner/repo:develop\n[merge-right]   -> clone OK\n[merge-right] push OK\n[merge-right] done",
			},
		},
	}
}

func buildMergeLeftExample() HelpExample {
	return HelpExample{
		WhatItDoes: "One-way file-level merge that writes only into LEFT. Files missing on LEFT are copied from RIGHT; conflicts are resolved into LEFT. RIGHT is never modified. If LEFT originated from a URL it is committed and pushed after the merge.",
		Examples: []Example{
			{
				Invocation:     "gitmap merge-left ./local https://github.com/owner/upstream -y",
				ExpectedOutput: "[merge-left] resolving RIGHT : https://github.com/owner/upstream\n[merge-left]   -> clone OK\n[merge-left] diffing trees ...\n[merge-left]   conflict main.go -> took RIGHT\n[merge-left] done",
			},
			{
				Invocation:     "gitmap ml ./mine ./theirs -y --prefer-left",
				ExpectedOutput: "[merge-left] diffing trees ...\n[merge-left]   conflict app.ts -> took LEFT\n[merge-left] done",
			},
		},
	}
}

func buildMergeBothExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Two-way bidirectional merge synchronizing changes between LEFT and RIGHT endpoints.",
		Examples: []Example{
			{
				Invocation:     "gitmap merge-both ./repo-a ./repo-b -y",
				ExpectedOutput: "[merge-both] diffing trees ...\n[merge-both] synchronized changes between endpoints\n[merge-both] done",
			},
		},
	}
}

func registerCommitInExamples(reg map[string]HelpExample) {
	ex := buildCommitInExample()
	reg["commit-in"] = ex
	reg["cin"] = ex
	reg["merge-in"] = buildMergeInExample()
	reg["mergein"] = reg["merge-in"]
}

func buildCommitInExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Chronological multi-source commit replay engine that walks repositories in author-date order, appends commits into the target repo, and tracks SHAs in SQLite.",
		Examples: []Example{
			{
				Invocation:     "gitmap commit-in ./canonical all --save-profile Default --set-default",
				ExpectedOutput: "[commit-in] walked all sibling repos\n[commit-in] profile 'Default' saved\n[commit-in] done",
			},
			{
				Invocation:     "gitmap cin ./canonical ./fork1,./fork2 --dry-run",
				ExpectedOutput: "[commit-in] dry-run: 12 commits analyzed\n[commit-in] done",
			},
		},
	}
}

func buildMergeInExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Chronological multi-source commit replay that walks one or more source repositories in author-date order and appends each commit into the target repository without rewriting history (alias for commit-in / cin).",
		Examples: []Example{
			{
				Invocation:     "gitmap merge-in ./canonical all --default",
				ExpectedOutput: "[commit-in] walked 3 repos: canonical-v1, canonical-v2, canonical-v3\n[commit-in] appended 42 commits\n[commit-in] done",
			},
			{
				Invocation:     "gitmap merge-in ./target https://github.com/org/v1.git,https://github.com/org/v2.git --dry-run",
				ExpectedOutput: "[commit-in] dry-run: would replay 18 commits into ./target\n[commit-in] done",
			},
		},
	}
}

func registerCommitTransferExamples(reg map[string]HelpExample) {
	reg["commit-right"] = buildCommitRightExample()
	reg["cmr"] = reg["commit-right"]
	reg["commit-left"] = buildCommitLeftExample()
	reg["cml"] = reg["commit-left"]
	reg["commit-both"] = buildCommitBothExample()
	reg["cmb"] = reg["commit-both"]
}

func buildCommitRightExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Replay commits from left repository onto right repository.",
		Examples: []Example{
			{
				Invocation:     "gitmap commit-right ./repo-a ./repo-b",
				ExpectedOutput: "[commit-transfer] replaying commits from ./repo-a onto ./repo-b ...\n[commit-transfer] done",
			},
		},
	}
}

func buildCommitLeftExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Replay commits from right repository onto left repository.",
		Examples: []Example{
			{
				Invocation:     "gitmap commit-left ./repo-a ./repo-b",
				ExpectedOutput: "[commit-transfer] replaying commits from ./repo-b onto ./repo-a ...\n[commit-transfer] done",
			},
		},
	}
}

func buildCommitBothExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Bidirectional commit replay between two repositories.",
		Examples: []Example{
			{
				Invocation:     "gitmap commit-both ./repo-a ./repo-b --interleave",
				ExpectedOutput: "[commit-both] interleaving commits by author date ...\n[commit-both] synchronized 8 commits\n[commit-both] done",
			},
		},
	}
}

func registerHelpAndSettingsExamples(reg map[string]HelpExample) {
	thelpEx := buildThelpExample()
	reg["thelp"] = thelpEx
	reg["th"] = thelpEx
	reg["specific"] = buildSpecificHelpExample()
	reg["settings"] = buildSettingsExample()
	reg["login"] = buildLoginExample()
}

func buildThelpExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Terminal help dispatcher for discovering commands, options, and concrete examples (alias for gitmap help).",
		Examples: []Example{
			{
				Invocation:     "gitmap thelp merge",
				ExpectedOutput: "Merge & Transfer Suite (gitmap merge / thelp merge)\nUsage synopsis, options, and commands",
			},
			{
				Invocation:     "gitmap thelp specific",
				ExpectedOutput: "Command-specific terminal help query system",
			},
		},
	}
}

func buildSpecificHelpExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Guide and command registry for querying topic-specific terminal help screens.",
		Examples: []Example{
			{
				Invocation:     "gitmap thelp specific",
				ExpectedOutput: "Topic table: merge, merge-right, merge-left, merge-in, commit-in, login, settings",
			},
		},
	}
}

func buildSettingsExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Manage unified speed settings, AI Auto bot configurations, thresholds, and fleet communication.",
		Examples: []Example{
			{
				Invocation:     "gitmap settings",
				ExpectedOutput: "● GitMap & Antigravity Unified Speed Settings\n  AI Auto Bot Name: gitmap automation bot - runs autonomously\n  AI Auto Bot Email: info@rgitmap.com\n  AI Auto Default Agent: gemini",
			},
			{
				Invocation:     "gitmap settings set ai_auto.default_agent antigravity",
				ExpectedOutput: "setting 'ai_auto.default_agent' updated to 'antigravity'",
			},
		},
	}
}

func buildLoginExample() HelpExample {
	return HelpExample{
		WhatItDoes: "Authenticate with Git providers (SSH, token, browser) and configure AI automation bot identity.",
		Examples: []Example{
			{
				Invocation:     "gitmap login set ai-auto gemini",
				ExpectedOutput: "Switched Git identity to AI bot:\n  Name: gitmap automation bot - (gemini) - runs autonomously\n  Email: info@rgitmap.com",
			},
			{
				Invocation:     "gitmap login set global ai-auto",
				ExpectedOutput: "Switched global Git identity to AI bot:\n  Name: gitmap automation bot - runs autonomously\n  Email: info@rgitmap.com",
			},
		},
	}
}
