---
description: "Implement Chrome SQLite file copy routine (macOS)"
labels: ["CLI", "ChromeProfile"]
status: "pending"
estimated_time: "15m"

# Plan Enhanced v4 Citations (MANDATORY)

citation_rule_0a: "lowercase-hyphenated 01-chrome-profile-migration.md"
citation_rule_0b: "audit only, not run"
citation_rule_0c: "self-looping, max 2 agents, 3 threads"
citation_rule_0d: "check 12-consolidated-guidelines"
citation_rule_0e: "user inputs in 08-chrome-profile-migration.md"
citation_rule_0f: "release protocol"
citation_rule_0g: "ambiguities in 01-new-ambiguity/"
citation_rule_0h: "folder structure respected"
citation_rule_0i: "RCA not applicable for feature"
citation_rule_0j: "CI/CD check"
citation_rule_4:  "8-section structure used"
citation_rule_8:  "batch-level commits"
---

# Task: Implement Chrome SQLite file copy routine (macOS)

## 1. Intent

Implement the specific logic for: Implement Chrome SQLite file copy routine (macOS).

## 2. Context

Part of the chrome-profile-migration feature ensuring cross-OS zip/sqlite exports for Chrome profiles. Reference: .lovable/spec/commands/08-chrome-profile-migration.md.

## 3. Scope

Confined to gitmap/cmd and gitmap/store files related to this step.

## 4. Architecture

Follows standard Gitmap CLI boundaries, avoiding global state, and wrapping errors with mt.Errorf.

## 5. Dependencies

Depends on step 9.

## 6. Testing

Ensure table-driven tests exist for any new parsing or validation logic. Run go test ./... after completion.

## 7. Delivery

Commit via the batch-level commit policy.

## 8. Anti-Clone Validation

- [x] Confirmed this task is unique.
