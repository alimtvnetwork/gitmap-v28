# Consolidated: App Issues — Complete Reference

**Version:** 3.3.0
**Updated:** 2026-04-16

---

## Purpose

Root-level folder for **application-specific issue analysis** — bug documentation, root cause analysis, diagnosis, and resolution guidance. Every bug gets a markdown file with structured analysis before a fix is implemented.

An AI reading only this file must be able to create, triage, diagnose, and document application issues correctly.

---

## Issue File Structure

Each issue file follows this template:

```markdown

# Issue #NN — [Title]

**Severity:** Critical | High | Medium | Low
**Status:** Open | Investigating | Resolved
**Created:** YYYY-MM-DD
**Resolved:** YYYY-MM-DD (if applicable)

## Symptom

What the user or developer observed. Include:
- Exact error message or unexpected behavior
- Steps to reproduce
- Environment (OS, browser, version)
- Frequency (always, intermittent, one-time)

## Diagnosis

Investigation steps taken:
1. What was checked first
2. What was ruled out
3. How the root cause was identified

## Root Cause

Technical explanation of why it happened. Include:
- Which component/module is responsible
- What condition triggers the bug
- Why it wasn't caught earlier

## Fix

What was changed and why:
- Files modified (with line-level detail)
- Why this approach was chosen over alternatives
- Any trade-offs accepted
