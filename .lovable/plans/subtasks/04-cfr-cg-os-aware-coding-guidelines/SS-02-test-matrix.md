---
Slug: test-matrix
Parent: 04-cfr-cg-os-aware-coding-guidelines
Status: pending
Created: 2026-07-16
---

# SS-02 — Test matrix

## Modifier parser (`clonefixrepo_modifiers_test.go`)

Table-driven cases (input tokens → expected flags):

| Args                          | WantCG | WantPublic | RepoURL       |
|-------------------------------|--------|------------|---------------|
| `cg <url>`                    | true   | false      | `<url>`       |
| `<url> cg`                    | true   | false      | `<url>`       |
| `p cg <url>`                  | true   | true       | `<url>`       |
| `cg p <url>`                  | true   | true       | `<url>`       |
| `cfrp cg <url>` (via wrapper) | true   | true       | `<url>`       |
| `<url>` (no modifiers)        | false  | false      | `<url>`       |
| `cg` (no URL)                 | true   | false      | `""` → error  |
| `--no-push cg <url>`          | true   | false      | `<url>` + NoPush |
| `--no-commit p cg <url>`      | true   | true       | `<url>` + NoCommit |

Assertions verify order-independence and correct extraction of `--no-push` / `--no-commit`.

## Runner dispatch (`codingguidelines_test.go`)

Inject a fake `Runner` that records `(name, args)`:

- On `GOOS=windows` (guarded with `runtime.GOOS` check + a `dispatchForOS(goos string, opts)` seam so the test runs on any host): expect the recorded command to start with `powershell` or `pwsh`, contain `-NoProfile`, `-ExecutionPolicy`, `Bypass`, `-Command`, and the exact `irm <WindowsURL> | iex` script.
- On `GOOS=linux` / `darwin`: expect `bash -c "curl -fsSL <UnixURL> | bash"` recorded verbatim.
- Missing binary path: fake `LookPath` returns `exec.ErrNotFound` → runner returns `ErrCGShellNotFound`, no command executed.

Golden URL constants asserted against `constants.DefaultCodingGuidelinesURLWindows` / `...Unix` to catch accidental URL drift.

## No network in tests

The test suite MUST NOT hit `raw.githubusercontent.com`. Verified by the injectable `Runner` never being replaced with the real `exec.Command` in test paths.
