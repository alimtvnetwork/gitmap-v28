package helpdoc

// Curated help examples: operations surface (task, schedule, db,
// pipeline, tree, summary, and discovery helpers). See examples.go —
// one registry, single source.
func registerOpsExamples(reg map[string]HelpExample) {
	reg["task"] = HelpExample{
		WhatItDoes: "Manages pending and completed task execution queues: inspect history, undo/redo operations, and track what the automation fleet is doing.",
		Examples: []Example{
			{
				Invocation: "gitmap task list",
				ExpectedOutput: `pending (<count>):
#12  release-asset-check   queued
done (<count>):  #11 theme-registry`,
			},
		},
	}
	reg["tasks"] = reg["task"]

	reg["schedule"] = HelpExample{
		WhatItDoes: "Manages scheduled tasks, macro schedules, OS startup tasks, and their run history and execution logs, backed by isolated split SQLite databases.",
		Examples: []Example{
			{
				Invocation: "gitmap schedule list",
				ExpectedOutput: `schedules (<count>):
nightly-pull   02:00 daily   enabled
db-backup      03:00 daily   enabled`,
			},
		},
	}
	reg["schedules"] = reg["schedule"]

	reg["db"] = HelpExample{
		WhatItDoes: "Split-DB operations: migrate, reset, and inspect the three-tier SQLite engine (gitmap.db, installation.db, repodb/pipeline.db).",
		Examples: []Example{
			{
				Invocation:     "gitmap db --help",
				ExpectedOutput: "the db help screen: migrate, reset, and inspection options",
			},
		},
	}
	reg["db-migrate"] = HelpExample{
		WhatItDoes: "Runs pending Split-DB SQLite migrations with rollback protection and history tracking.",
		Examples: []Example{
			{
				Invocation:     "gitmap db-migrate",
				ExpectedOutput: "migrations: <count> applied, 0 pending",
			},
		},
	}
	reg["db-reset"] = HelpExample{
		WhatItDoes: "Resets the Split-DB databases to a clean state. Destructive — snapshots first when available.",
		Examples: []Example{
			{
				Invocation:     "gitmap db-reset --yes",
				ExpectedOutput: "databases reset — snapshot kept at <path>",
			},
		},
	}

	reg["pipeline"] = HelpExample{
		WhatItDoes: "Queries live CI/CD pipeline runs: remaining wait times (ETA), step failure logs, telemetry export, and error handoff to Antigravity.",
		Examples: []Example{
			{
				Invocation: "gitmap pipeline",
				ExpectedOutput: `run <count>: Release — in progress (ETA 4m)
run <count>: CI — success`,
			},
		},
	}
	reg["pl"] = reg["pipeline"]
	reg["pe"] = HelpExample{
		WhatItDoes: "Pipeline errors inspector: shows failing CI/CD runs with step failure logs, then feeds the errors to Antigravity with a fix prompt.",
		Examples: []Example{
			{
				Invocation: "gitmap pe",
				ExpectedOutput: `failing runs (<count>):
Release / windows-build — exit 1
  step log: <path>

fix prompt copied to clipboard`,
			},
		},
	}
	reg["pea"] = HelpExample{
		WhatItDoes: "Pipeline errors for all repos: the fleet-wide variant of pe — failing runs across every tracked repository with fix prompts.",
		Examples: []Example{
			{
				Invocation: "gitmap pea",
				ExpectedOutput: `failing runs (<count>) across <count> repos:
repos/gitmap-v28  Release / windows-build — exit 1`,
			},
		},
	}
	reg["peat"] = reg["pea"]
	reg["te"] = HelpExample{
		WhatItDoes: "Telemetry errors view: the pipeline-errors family's telemetry surface — structured error telemetry for failing runs.",
		Examples: []Example{
			{
				Invocation:     "gitmap te",
				ExpectedOutput: "telemetry errors: <count> failing steps across <count> runs",
			},
		},
	}

	reg["tree"] = HelpExample{
		WhatItDoes: "Renders the repository folder tree with Split-DB persistence: fast hierarchical views of large workspaces.",
		Examples: []Example{
			{
				Invocation: "gitmap tree ~/workspace/repos --depth 2",
				ExpectedOutput: `repos/
  gitmap-v28/
  story-telling-v1/
<count> dirs, <count> files`,
			},
		},
	}

	reg["summary"] = HelpExample{
		WhatItDoes: "Prints the repository status summary: branch, dirty state, and pending changes across the workspace.",
		Examples: []Example{
			{
				Invocation: "gitmap summary",
				ExpectedOutput: `repos/gitmap-v28        main  clean
repos/story-telling-v1  main  dirty (3 files)`,
			},
		},
	}

	reg["new-commands"] = HelpExample{
		WhatItDoes: "Catalogs and searches CLI commands added in the last 100 commits, with syntax and usage examples — the fastest way to discover what's new.",
		Examples: []Example{
			{
				Invocation: "gitmap nc",
				ExpectedOutput: `recent commands (<count>):
help example   canonical what-it-does + invocation + output
--theme        design-system terminal themes`,
			},
		},
	}
	reg["nc"] = reg["new-commands"]

	reg["pending-commits"] = HelpExample{
		WhatItDoes: "Discovers uncommitted modifications and unpushed commits with a tree view, remediation hints, and batch fix options.",
		Examples: []Example{
			{
				Invocation: "gitmap pc",
				ExpectedOutput: `repos/story-telling-v1  3 uncommitted files
repos/gitmap-v28        1 unpushed commit`,
			},
		},
	}
	reg["pc"] = reg["pending-commits"]

	reg["find"] = HelpExample{
		WhatItDoes: "High-speed file finder across the workspace: find by name, prefix, suffix, or regex, with JSON output for scripting.",
		Examples: []Example{
			{
				Invocation:     "gitmap find themes.go cli/",
				ExpectedOutput: `cli/termout/themes.go`,
			},
		},
	}

	reg["space"] = HelpExample{
		WhatItDoes: "Workspace space management: disk usage, backup branches, and workspace hygiene operations.",
		Examples: []Example{
			{
				Invocation:     "gitmap space",
				ExpectedOutput: "space help: usage, backup-branch, and cleanup options",
			},
		},
	}

	reg["help-dashboard"] = HelpExample{
		WhatItDoes: "Opens the help dashboard: a browsable overview of the help system, topics, and recent commands.",
		Examples: []Example{
			{
				Invocation:     "gitmap help-dashboard",
				ExpectedOutput: "help dashboard: <count> topics indexed",
			},
		},
	}
}
