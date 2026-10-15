package helpdoc

// Curated help examples: core repo workflows (scan, clone, pull, commit
// suite) plus version/help/theme-flag. One Go registry is the single
// source of truth — terminal help, `help example`, and generated docs
// all render from here.
func registerCoreExamples(reg map[string]HelpExample) {
	reg["scan"] = HelpExample{
		WhatItDoes: "Scans a directory tree for git repositories and caches them in the Split-DB. Run it before any repo-wide command so gitmap knows what exists on disk.",
		Examples: []Example{
			{
				Invocation: "gitmap scan ~/workspace/repos",
				ExpectedOutput: `scanned: 14 repos in 0.83s (cache: gitmap.db)
path                                branch  status
repos/gitmap-v28                    main     clean
repos/story-telling-v1              main     dirty`,
			},
			{
				Invocation:     "gitmap scan --json",
				ExpectedOutput: `{"repos":[{"path":"<path>","branch":"main","status":"clean","count":<count>}],"scanned":<count>}`,
			},
		},
	}
	reg["s"] = reg["scan"]

	reg["clone"] = HelpExample{
		WhatItDoes: "Clones one or more GitHub repositories into the workspace, preserving the owner/repo folder hierarchy. Accepts owner/repo slugs, full URLs, or a comma-separated list.",
		Examples: []Example{
			{
				Invocation: "gitmap clone alimtvnetwork/story-telling-v1",
				ExpectedOutput: `cloning alimtvnetwork/story-telling-v1 → repos/story-telling-v1
✓ cloned <sha> (main) in 4.2s`,
			},
			{
				Invocation:     "gitmap clone alimtvnetwork/gitmap-v28,alimtvnetwork/alim-cv-v8 --json",
				ExpectedOutput: `{"cloned":[{"repo":"alimtvnetwork/gitmap-v28","path":"<path>","sha":"<sha>"}],"failed":[]}`,
			},
		},
	}

	reg["pull"] = HelpExample{
		WhatItDoes: "Pulls the latest changes for a single repository, a group, or every tracked repository. Use pull-all (pa) for the fast batch pull across the whole catalog.",
		Examples: []Example{
			{
				Invocation: "gitmap pull",
				ExpectedOutput: `pulling repos/gitmap-v28 (main)…
✓ up to date — <sha>`,
			},
		},
	}
	reg["pull-all"] = HelpExample{
		WhatItDoes: "Batch-pulls every tracked repository in the catalog with a fast concise progress view. Add --status for the full post-pull table or --ssh to delegate the pull to fleet nodes.",
		Examples: []Example{
			{
				Invocation: "gitmap pa",
				ExpectedOutput: `pull-all: 14 repos
✓ 12 up to date   ⟳ 2 updated   ✗ 0 failed   (3.1s)`,
			},
			{
				Invocation: "gitmap pa --status",
				ExpectedOutput: `repo                   branch  before    after     result
repos/gitmap-v28         main    <sha>     <sha>     updated
repos/story-telling-v1   main    <sha>     <sha>     up-to-date`,
			},
		},
	}
	reg["pa"] = reg["pull-all"]
	reg["pat"] = reg["pull-all"]

	reg["commit"] = HelpExample{
		WhatItDoes: "Stages all changes and creates a flat commit with your message. Add --push to push immediately or --dry-run to preview the commit without writing it.",
		Examples: []Example{
			{
				Invocation: `gitmap commit "theme - design-system theme registry"`,
				ExpectedOutput: `staged 6 files
✓ committed <sha> — "theme - design-system theme registry"`,
			},
		},
	}
	reg["c"] = reg["commit"]

	reg["cm"] = HelpExample{
		WhatItDoes: "Semantic commit helper: stages all changes, commits with the message you give it, and stays local. The cpf/cpb/cpr siblings add the type prefix and push.",
		Examples: []Example{
			{
				Invocation: `gitmap cm "help - canonical help-example registry"`,
				ExpectedOutput: `staged 4 files
✓ committed <sha> — "help - canonical help-example registry"`,
			},
		},
	}

	reg["cpf"] = HelpExample{
		WhatItDoes: "Stages all changes, commits with a \"Feature: \" prefix, and pushes to origin. The atomic commit+push flow used for every gitmap change.",
		Examples: []Example{
			{
				Invocation: `gitmap cpf "help - canonical help-example registry"`,
				ExpectedOutput: `staged 4 files
✓ committed <sha> — "Feature: help - canonical help-example registry"
✓ pushed to origin/main`,
			},
		},
	}
	reg["commit-push-feature"] = reg["cpf"]

	reg["cpb"] = HelpExample{
		WhatItDoes: "Stages all changes, commits with a \"Bug: \" prefix, and pushes to origin. Use for fixes; cpf is for features and cpr for releases.",
		Examples: []Example{
			{
				Invocation: `gitmap cpb "theme - fix mono strip on piped help"`,
				ExpectedOutput: `staged 2 files
✓ committed <sha> — "Bug: theme - fix mono strip on piped help"
✓ pushed to origin/main`,
			},
		},
	}
	reg["commit-push-bug"] = reg["cpb"]

	reg["cpr"] = HelpExample{
		WhatItDoes: "Stages all changes, commits with a \"Release: \" prefix, and pushes to origin. Reserved for version bumps and release commits.",
		Examples: []Example{
			{
				Invocation: `gitmap cpr "v6.554.0"`,
				ExpectedOutput: `staged 3 files
✓ committed <sha> — "Release: v6.554.0"
✓ pushed to origin/main`,
			},
		},
	}
	reg["commit-push-release"] = reg["cpr"]

	reg["pcp"] = HelpExample{
		WhatItDoes: "Pulls latest with rebase, stages all changes, commits with your message, and pushes. The full sync-commit-push cycle in one command.",
		Examples: []Example{
			{
				Invocation: `gitmap pcp "docs - refresh help topics"`,
				ExpectedOutput: `✓ pulled — already up to date
staged 5 files
✓ committed <sha> — "docs - refresh help topics"
✓ pushed to origin/main`,
			},
		},
	}
	reg["pull-commit-push"] = reg["pcp"]

	reg["commit-in"] = HelpExample{
		WhatItDoes: "Comprehensive commit automation engine with JSON author rotation, SEO templates, and deduplication heuristics. For everyday commits prefer cm/cpf.",
		Examples: []Example{
			{
				Invocation: `gitmap commit-in --dry-run`,
				ExpectedOutput: `dry-run: would commit 3 staged files as <author>
message: "Chore: <summary>"`,
			},
		},
	}
	reg["commitin"] = reg["commit-in"]
	reg["commit-write"] = reg["commit-in"]

	reg["amend"] = HelpExample{
		WhatItDoes: "Amends the last commit: folds staged changes into it without creating a new commit. Lists recent amendable commits with amend-list.",
		Examples: []Example{
			{
				Invocation:     "gitmap amend",
				ExpectedOutput: `amended <sha> — 2 files folded into the previous commit`,
			},
		},
	}
	reg["amend-list"] = HelpExample{
		WhatItDoes: "Lists recent commits that can be amended, newest first.",
		Examples: []Example{
			{
				Invocation: "gitmap amend-list",
				ExpectedOutput: `<sha>  <date>  theme - design-system theme registry
<sha>  <date>  help - canonical help-example registry`,
			},
		},
	}

	reg["version"] = HelpExample{
		WhatItDoes: "Prints the gitmap binary version. Use version-history or list-versions to see the release timeline.",
		Examples: []Example{
			{
				Invocation:     "gitmap version",
				ExpectedOutput: "gitmap version 6.553.0",
			},
		},
	}

	reg["help"] = HelpExample{
		WhatItDoes: "Shows help: bare `gitmap help` prints the grouped command index, `gitmap help <command>` shows that command's help screen, and `gitmap help example <command>` shows its canonical what-it-does + invocation + expected output.",
		Examples: []Example{
			{
				Invocation:     "gitmap help scan",
				ExpectedOutput: "the scan help screen: usage synopsis, flags, and examples",
			},
			{
				Invocation: "gitmap help example clone",
				ExpectedOutput: `gitmap clone — help example

WHAT IT DOES
  Clones one or more GitHub repositories into the workspace…

INVOCATION
  $ gitmap clone alimtvnetwork/story-telling-v1

EXPECTED OUTPUT
  cloning alimtvnetwork/story-telling-v1 → repos/story-telling-v1`,
			},
		},
	}

	reg["theme"] = HelpExample{
		WhatItDoes: "Global terminal palette switch. Pick one of the four design-system themes (navy-purple default, tokyo-night, vscode-dark-plus, monokai-pro); legacy bright/standard/monochrome aliases still work. Piped output never carries ANSI — colors render on a TTY only.",
		Examples: []Example{
			{
				Invocation:     "gitmap --theme tokyo-night scan",
				ExpectedOutput: "scan results rendered in the Tokyo Night palette (tty only; piped output is plain)",
			},
			{
				Invocation:     "gitmap --theme mono pa | cat",
				ExpectedOutput: "pull-all progress with every SGR escape stripped — byte-clean for logs and diff",
			},
		},
	}
	reg["os-theme"] = HelpExample{
		WhatItDoes: "Switches the desktop appearance between Dark Mode and Light Mode across Windows and Linux. (For gitmap's own terminal palette see `gitmap help example theme`.)",
		Examples: []Example{
			{
				Invocation:     "gitmap os-theme dark",
				ExpectedOutput: "desktop appearance → dark",
			},
		},
	}
}
