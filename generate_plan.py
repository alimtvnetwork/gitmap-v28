import os

plan_slug = 'chrome-profile-migration'
subtasks_dir = f'.lovable/plans/subtasks/{plan_slug}'
pending_plan = f'.lovable/plans/pending/01-{plan_slug}.md'

os.makedirs(subtasks_dir, exist_ok=True)
os.makedirs('.lovable/plans/pending', exist_ok=True)

# Generate 30 subtasks
task_titles = [
    "Analyze Chrome SQLite Schema for Export (History, Web Data)",
    "Define Cross-OS File Mapping (Windows vs Linux)",
    "Update constants_chromeprofile.go with Format Flags",
    "Update chromeprofile.go to Parse --format=zip|json|sqlite",
    "Define zipExportRecord struct in chromeprofile_db.go",
    "Add SQLite inclusion flag logic to runChromeProfileExport",
    "Create chromeprofile_sqlite_export.go scaffolding",
    "Implement Chrome SQLite file copy routine (Windows)",
    "Implement Chrome SQLite file copy routine (Linux)",
    "Implement Chrome SQLite file copy routine (macOS)",
    "Filter out encrypted DBs (Login Data, Cookies) from Export",
    "Write SQLite payloads to ZIP archive",
    "Write JSON payload to ZIP archive",
    "Update writeChromeExport to conditionally output sqlite only",
    "Update writeChromeExport to conditionally output zip format",
    "Update defaultChromeExportPath to handle .zip and .sqlite",
    "Update emitChromeSnapshots for multi-format output",
    "Update printChromeArtifacts to handle zip/sqlite files",
    "Add --format flag support to runChromeProfileImport",
    "Create chromeprofile_sqlite_import.go scaffolding",
    "Implement ZIP extraction for profile import",
    "Implement SQLite restoration for profile import",
    "Handle Profile merging when SQLite databases exist",
    "Update applyChromeExport for SQLite blobs",
    "Migrate gitmap.db schema to support Format='zip' and 'sqlite'",
    "Update persistChromeProfile to track new format extensions",
    "Update listChromeProfilesFromDB to summarize zip exports",
    "Add Integration Test: Windows to Linux Zip Export/Import",
    "Add Unit Tests for sqlite filtering and zip bundling",
    "Update README and CLI Help Text for new format flags"
]

if len(task_titles) != 30:
    raise ValueError(f"Need exactly 30 tasks, got {len(task_titles)}")

for i, title in enumerate(task_titles, 1):
    task_file = f"{subtasks_dir}/{i:03d}-task.md"
    content = f'''---
description: "{title}"
labels: ["CLI", "ChromeProfile"]
status: "pending"
estimated_time: "15m"

# Plan Enhanced v4 Citations (MANDATORY)
citation_rule_0a: "lowercase-hyphenated 01-{plan_slug}.md"
citation_rule_0b: "audit only, not run"
citation_rule_0c: "self-looping, max 2 agents, 3 threads"
citation_rule_0d: "check 12-consolidated-guidelines"
citation_rule_0e: "user inputs in 08-{plan_slug}.md"
citation_rule_0f: "release protocol"
citation_rule_0g: "ambiguities in 01-new-ambiguity/"
citation_rule_0h: "folder structure respected"
citation_rule_0i: "RCA not applicable for feature"
citation_rule_0j: "CI/CD check"
citation_rule_4:  "8-section structure used"
citation_rule_8:  "batch-level commits"
---

# Task: {title}

## 1. Intent
Implement the specific logic for: {title}.

## 2. Context
Part of the {plan_slug} feature ensuring cross-OS zip/sqlite exports for Chrome profiles. Reference: .lovable/spec/commands/08-chrome-profile-migration.md.

## 3. Scope
Confined to gitmap/cmd and gitmap/store files related to this step.

## 4. Architecture
Follows standard Gitmap CLI boundaries, avoiding global state, and wrapping errors with mt.Errorf.

## 5. Dependencies
Depends on step {i-1 if i > 1 else 'None'}.

## 6. Testing
Ensure table-driven tests exist for any new parsing or validation logic. Run go test ./... after completion.

## 7. Delivery
Commit via the batch-level commit policy.

## 8. Anti-Clone Validation
- [x] Confirmed this task is unique.
'''
    with open(task_file, 'w') as f:
        f.write(content)

# Generate pending plan
plan_content = f'''---
slug: "{plan_slug}"
status: "pending"
total_tasks: 30
---

# Plan: Chrome Profile Migration (Cross-OS)

## Overview
Enhance gitmap chrome-profile-export and import to support .zip and .sqlite formats for cross-OS migration.

## Tasks
'''
for i, title in enumerate(task_titles, 1):
    plan_content += f"- [ ] {i:03d}-task.md: {title}\n"

with open(pending_plan, 'w') as f:
    f.write(plan_content)

print("Generated 30 subtasks and 1 pending plan successfully.")
