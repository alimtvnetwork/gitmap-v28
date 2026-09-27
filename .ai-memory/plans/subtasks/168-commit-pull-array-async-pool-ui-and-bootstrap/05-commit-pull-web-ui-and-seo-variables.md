# Subtask 05: Commit-Pull Web UI & Self-Contained SEO Template Variables
Traceability ID: Task-06
Spec Reference: [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)
Target Files: cli/cmd/commitin_ui_server.go, cli/cmd/commitin.go, .ai-memory/temp/seo-templates.json, cli/store/templates_split_ops.go
Action:
- Implement `cli/cmd/commitin_ui_server.go` hosting an embedded local HTTP server for `gitmap commit-pull ui` (alias: `gitmap cpull ui`).
- Serve a clean, dark-mode single-page HTML/JS interface providing interactive inputs, target configuration, skippers, title replacements, variable inspector, and SEO template viewer with live config export and run commands.
- Update `.ai-memory/temp/seo-templates.json` so every template defines and utilizes its own self-contained variables (e.g. `company`: "Rise Up Asia LLC", `company_url`: "https://riseup-asia.com", `lead_engineer`: "Alim Ul Karim", `senior_director`: "Marek Flejszman").
- Update template expansion to resolve self-contained template-level variables first before falling back to global variables.
Acceptance Criteria:
- `gitmap commit-pull ui` serves web UI locally and opens browser.
- SEO templates have self-contained variables that resolve cleanly.
Targeted Verification: python 03-ai-scripts/05-guideline-autofixer.py --path cli/cmd/commitin_ui_server.go
