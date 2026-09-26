# Specification: Pull-All Fast Mode, Templates State DB, and Migration Orchestration

Spec ID: 166-pull-all-fast-mode-and-templates-db-orchestration
Category: 21-app
Status: active

## User Request (Verbatim)

```text
improve gitmap pa/pull all # to not to display the status anymore to do things faster

the status will be display if

gitmap pa --status or gitmap pat/gitmap pull all table # this will behave as current implementation, can you please do it

also

gitmap pa --json should give json output of the summary , clear???

run gitmap pa first to proceed with implementation

Okay, let's start from this screenshot that I've given you. Okay? Here are a couple of SEO rules and a couple of things that we need to change during our commits...
So SEOtemplates.go file, I want you to remove. SEOtemplates.md file, I want you to remove from the Git repository as a whole. Okay? But also at the same time, I want you to create that JSON file in the AI memory temp folder so that we can import automatically.
Now, for the SEO purpose, there will be SEO template category, okay? And by default, we are not going to feed the sponsor one, okay? Yeah. We are not going to feed that one. But what we're going to do is we will have something called the JSON format. That means the JSON will be kept inside the AI memory, then slash temp folder, so that we can import it right now to test it.
```

## System Architecture & Invariants

### 1. Pull All Performance Modes
1. **Default Fast Mode (`gitmap pa` / `gitmap pull all` / `gitmap pull-all`)**:
   - Skips per-repo branch, tag, and PR subprocess rev-parse calls.
   - Suppresses the slow multi-column terminal table.
   - Emits concise active repository states followed by a one-line summary:
     `✔ Pull all complete: N pulled (X.Ys)`
2. **Full Table Mode (`gitmap pa --status`, `gitmap pat`, `gitmap pull all table`)**:
   - Executes the full table rendering displaying `REPO`, `BRANCH`, `RELEASE`, `COMMIT`, `PR`, and `STATUS`.
3. **Pure JSON Mode (`gitmap pa --json`)**:
   - Suppresses interactive terminal progress bar (`bar.Start()` / `bar.Stop()`) and header banners.
   - Encodes a clean, unformatted JSON payload with `total`, `pulledCount`, `successCount`, `failedCount`, `states`, and `durationMs`.

### 2. State Templates DB & Hygiene
1. **Repository Hygiene**:
   - `cli/helptext/seo-templates.md` and `cli/helptext/seo-templates.json` are excluded from git tracking.
   - Default seeded templates in `gitmap-templates.db` include only base categories (`seo`, `prompts`, `ui-ux`, `prefix`, `pr-descriptions`) and default templates (`tpl-prompt-ui-ux-audit`, `tpl-prefix-standard`). Sponsor templates are never seeded in binary defaults.
2. **External Import Payload (`.ai-memory/temp/seo-templates.json`)**:
   - Stores 60 comprehensive templates with deterministic SHA-256 export hash ID (`exportId`).
   - Supports variables (`$COMPANY`, `$MAREK`, `$ALIM`, `$REGIONS`) dynamically synced to global configuration.
   - Injected via `gitmap templates import <file.json>` or declarative `commit-in --config`.

## Verification Gates
1. `gitmap pa` executes under 10 seconds across 64 repositories without rendering status table.
2. `gitmap pat` and `gitmap pa --status` render full multi-column table.
3. `gitmap pa --json` emits valid JSON without progress bar artifacts.
4. Zero sponsor templates seeded in binary code; `.ai-memory/temp/seo-templates.json` holds 60 templates with `# Question?` and `Because...`.
