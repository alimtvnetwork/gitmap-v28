---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_output_stream.go"]
depends_on: ["Task 014" if int(num) > 1 else "None"]
---
# Task 015 — Stream Output Collector

## 1. Goal
Capture and format stdout/stderr from installer script in `installer/prompt_output_stream.go`.
