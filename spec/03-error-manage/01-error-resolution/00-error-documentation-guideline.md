# Error Documentation Guideline

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

When any error is encountered and fixed in code, it **MUST** be documented in the spec folder to prevent future AI hallucination and repeated mistakes.

---

## Mandatory Process

### Step 1: Create an Issue File

Location: `spec/03-error-manage/01-error-resolution/app-issues/`

File naming: `YYYY-MM-DD-short-description.md` (e.g., `2026-04-02-cors-header-missing.md`)

### Step 2: Document Using This Template

```markdown

# Issue: [Short Title]

**Date:** YYYY-MM-DD  
**Severity:** Critical | High | Medium | Low  
**Status:** Resolved

---

## Error Description

What happened? Include exact error messages, HTTP status codes, stack traces, or UI symptoms.

## Root Cause

Why did it happen? Be specific — name the file, function, misconfiguration, or logic flaw.

## Solution

What was done to fix it? Include the exact code change, config update, or architectural decision.

## Prevention

How to prevent recurrence? Name the rule, lint check, or test that should catch this in the future.

## Related

- Link to any related spec files, error codes, or PRs.
```

### Step 3: Update the Error Resolution Index

Add the new issue to `01-error-resolution/00-overview.md` inventory.

---

## Why This Matters

| Without documentation | With documentation |
|-----------------------|-------------------|
| AI re-introduces the same bug | AI checks known issues before generating code |
| Debugging time wasted on solved problems | Instant lookup of root cause + fix |
| No institutional memory | Searchable knowledge base grows over time |

---

## Rules

1. **Every fix gets documented** — no exceptions, even for "trivial" bugs
2. **Root cause is mandatory** — "it just works now" is NOT acceptable
3. **AI must check `app-issues/` before proposing fixes** — to avoid hallucinating solutions that were already tried and failed
4. **Keep files under 100 lines** — one issue per file, concise and scannable

---

*This guideline is mandatory for all AI-assisted development sessions.*
