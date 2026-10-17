# Spec 277 — Muse Memory IO and GitMap Integration

**Status:** SPEC + IMPLEMENTATION IN PROGRESS (Worker 01: spec files, `gitmap llm train muse` command, new-command docs, changelog). This is the GitMap half of the "Muse Memory IO" program: a command that lets any LLM re-sync with the Muse 7-day memory digest authored into the repo-secrets vault by the memory-IO program.

- `01-scope-and-plan.md` — metadata, numbered IN SCOPE / OUT OF SCOPE, digest contract consumed, command design, file plan, coding constraints.
- `02-acceptance-criteria-and-evidence.md` — per-feature acceptance criteria; each criterion names the exact evidence (file paths, function names) the implementer must produce.
