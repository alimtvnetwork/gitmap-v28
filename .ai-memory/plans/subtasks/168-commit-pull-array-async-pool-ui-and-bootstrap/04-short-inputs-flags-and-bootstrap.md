# Subtask 04: Shortened Input Syntax, Terminal Flags & Bootstrap Command
Traceability ID: Task-05
Spec Reference: [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)
Target Files: cli/cmd/commitin.go, cli/cmd/commitin/config_json.go, cli/helptext/commit-pull.md, readme.md
Action:
- Add support for short comma-separated input syntax: e.g. `gitmap commit-pull "D:\target" "git-repo-navigator,gitmap-v{2..28}"` where base remote URL (e.g. `https://github.com/alimtvnetwork/` or detected user remote) is prepended automatically when not a full URL or absolute path.
- Add `gitmap commit-pull bootstrap [-file <path>]` (and aliases `gitmap cpull bootstrap`, `gitmap commit-pull init`) to generate a fully documented `commit-pull-config.json` template.
- Update `cli/helptext/commit-pull.md` and root `readme.md` with rich examples demonstrating short syntax, range expansion, and the bootstrap command.
Acceptance Criteria:
- `gitmap commit-pull bootstrap` generates a valid JSON file.
- Short comma-separated input names expand with the inferred owner/base URL.
- Help docs and root readme updated with examples.
Targeted Verification: python 03-ai-scripts/09-cli-help-auditor.py
