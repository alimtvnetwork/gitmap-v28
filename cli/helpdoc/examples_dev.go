package helpdoc

// Curated help examples: IDE/tooling surface (agm, agy, install, update,
// prompt, llm-train, ai). See examples.go — one registry, single source.
func registerDevExamples(reg map[string]HelpExample) {
	reg["agm"] = HelpExample{
		WhatItDoes: "Manages the Antigravity Manager GUI desktop application: installs and updates its tools locally or across remote SSH fleet nodes concurrently.",
		Examples: []Example{
			{
				Invocation:     "gitmap agm update",
				ExpectedOutput: "antigravity-manager-tools: already at v4.200.1",
			},
		},
	}

	reg["agy"] = HelpExample{
		WhatItDoes: "Manages Google Antigravity workspaces: reruns historical prompts, inspects prompt diffs across projects, and orchestrates remote execution.",
		Examples: []Example{
			{
				Invocation: "gitmap agy lap",
				ExpectedOutput: `active projects (last 24h):
<path>  P1  <date>
<path>  P2  <date>`,
			},
		},
	}
	reg["antigravity"] = reg["agy"]

	reg["install"] = HelpExample{
		WhatItDoes: "Cross-platform installer for developer tools, runtimes, AI models, custom scripts, and workstation profiles. Use --tools for batch installs.",
		Examples: []Example{
			{
				Invocation: "gitmap install --tools chrome,vscode",
				ExpectedOutput: `installing chrome… ✓
installing vscode… ✓
2 installed, 0 failed`,
			},
			{
				Invocation:     "gitmap in antigravity chrome vscode --json",
				ExpectedOutput: `{"results":[{"tool":"antigravity","status":"installed"},{"tool":"chrome","status":"installed"}],"failed":[]}`,
			},
		},
	}
	reg["in"] = reg["install"]

	reg["update"] = HelpExample{
		WhatItDoes: "Updates gitmap itself (and with update-all, every tracked repo's tooling). Shows the update help screen with no arguments.",
		Examples: []Example{
			{
				Invocation:     "gitmap update",
				ExpectedOutput: "update help: self-update, update-all, update-cleanup, and runner options",
			},
		},
	}
	reg["ua"] = reg["update"]
	reg["update-all"] = reg["update"]

	reg["prompt"] = HelpExample{
		WhatItDoes: "Manages reusable Antigravity prompt prefix and verification templates: add, edit, list, remove, export, and import prompt templates stored as JSON.",
		Examples: []Example{
			{
				Invocation: "gitmap prompt list",
				ExpectedOutput: `templates (<count>):
verify-build    "verify the build before reporting done"
no-silent-error "surface the real error message with a fallback"`,
			},
		},
	}

	reg["llm-train"] = HelpExample{
		WhatItDoes: "Trains the local LLM helper index over the repository so natural-language repo questions resolve against real code. Part of the llm/aum automation surface.",
		Examples: []Example{
			{
				Invocation:     "gitmap llm-train",
				ExpectedOutput: "training index over <count> files… ✓ ready",
			},
			{
				Invocation: "gitmap llm train muse --days 7",
				ExpectedOutput: `GITMAP LLM TRAIN — MUSE MEMORY SYNC (last 7 days):
Digest root: repo-secrets/Muse
Commits: <count> | Files changed: <count> | Sequenced digests: <count>
AI TRAVERSAL CHECKLIST: gitmap cat repo-secrets/Muse/README.md …`,
			},
		},
	}
	reg["train"] = reg["llm-train"]
	reg["llmtrain"] = reg["llm-train"]

	reg["ai"] = HelpExample{
		WhatItDoes: "Native AI scripts automation catalog: lists, streams, and executes repository automation scripts, including the autofix engine.",
		Examples: []Example{
			{
				Invocation: "gitmap ai list",
				ExpectedOutput: `scripts (<count>):
01  seo-infographic        roll out article infographics
02  helptext-generator     regenerate helpdoc topics`,
			},
		},
	}
	reg["scripts"] = reg["ai"]

	reg["aum"] = HelpExample{
		WhatItDoes: "High-performance native Go automation: multi-core search, polyglot newline normalization, and side-by-side execution benchmarks.",
		Examples: []Example{
			{
				Invocation: "gitmap aum search TODO cli/",
				ExpectedOutput: `cli/cmd/root.go:88: // TODO: migrate to displayer
<count> hit(s) across <count> file(s)`,
			},
		},
	}
	reg["automation"] = reg["aum"]
	reg["auto"] = reg["aum"]

	reg["fix"] = HelpExample{
		WhatItDoes: "Runs the standardized repository autofix targets: coding-guideline fixes, newline normalization, naming, and path hygiene.",
		Examples: []Example{
			{
				Invocation:     "gitmap fix",
				ExpectedOutput: `autofix: <count> files normalized, 0 errors`,
			},
		},
	}

	reg["docs"] = HelpExample{
		WhatItDoes: "Opens the generated documentation surface for a topic: the same content as `gitmap help <topic>`, rendered for the docs viewer.",
		Examples: []Example{
			{
				Invocation:     "gitmap docs scan",
				ExpectedOutput: "the scan help topic rendered in the docs viewer",
			},
		},
	}

	reg["cg"] = HelpExample{
		WhatItDoes: "Orchestrates Coding Guidelines standards, Prompt Architect templates, compliance audits, and cross-repo synchronization.",
		Examples: []Example{
			{
				Invocation: "gitmap cg install --all --dry-run",
				ExpectedOutput: `coding-guidelines: dry-run install across <count> repositories
✓ scanned all target repositories
0 files modified (dry-run)`,
			},
			{
				Invocation: "gitmap cg audit",
				ExpectedOutput: `coding-guidelines audit:
✓ 02-spec compliance: PASS
✓ .ai-memory compliance: PASS
✓ naming & boolean conventions: PASS`,
			},
		},
	}
	reg["coding-guidelines"] = reg["cg"]

	reg["install-prompts"] = HelpExample{
		WhatItDoes: "Installs Prompt Architect v2 suite and canonical prompt library into target repositories.",
		Examples: []Example{
			{
				Invocation: "gitmap cg install-prompts --all",
				ExpectedOutput: `installing prompt architect templates across <count> repositories… ✓
prompts synchronized successfully`,
			},
		},
	}
}
