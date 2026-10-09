# Subtask 04 — AGM scope tiers + unwrap reduction (task 251)

## Scope
`~/workspace/repos/Antigravity-Manager` — improvements only, NO deletions
(hard rule), NO merges.

## Build (spec 01 §B4)
1. **Scope tier map**: inventory `src-tauri/src/modules/` (54 files). Classify
   each as CORE (protocol gateway), ADJACENT (email, telegram, supabase, ssh,
   cloudflared, etc.), or SUPPORT (config, db, logging). Write the map to
   `docs/module-tiers.md` (new file). Add a guardrail: a `cargo` doc comment
   in `proxy/mod.rs` stating adjacent modules must not be imported by core
   proxy code (documentary guardrail, not a compile gate).
2. **unwrap/expect reduction**: in `src-tauri/src/proxy/handlers/` only,
   replace `unwrap()`/`expect()` on untrusted-input paths with `?` or
   explicit error returns. Count before/after per file. Do NOT touch logic;
   only the panic paths. If a replacement is non-trivial, leave it and list
   it in the report.
3. Respect the repo's AGENTS.md (read it first): `cargo fmt -- --check` and
   `cargo clippy` must pass for touched files. Do NOT run `cargo test`.

## Rules
GitMap tools for search; no rg/grep. Never `git add`/`git commit`.
Never delete files.

## Deliverable
Diff stat + unwrap counts before/after + the tier map file. Then stop.
