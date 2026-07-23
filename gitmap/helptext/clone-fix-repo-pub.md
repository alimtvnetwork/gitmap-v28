# gitmap clone-fix-repo-pub

> 🚀 **One-shot**: `clone` → `cd` → `fix-repo --all` → `make-public --yes`.
> Same URL semantics as `gitmap clone`, including transport coercion
> (`--ssh` / `--https`) and versioned-URL auto-flatten.

Replaces the manual four-step dance:

```
gitmap clone <url>
cd <folder>
gitmap fix-repo --all
gitmap make-public --yes
```

## Aliases

- 🪄 `cfrp` — short form

## Synopsis

```
gitmap clone-fix-repo-pub [modifiers...] <url> [folder] [flags]
gitmap cfrp               [modifiers...] <url> [folder] [flags]
```

## Requirements

- `gh` or `glab` installed and authenticated (`gh auth login` /
  `glab auth login`). The `make-public` step wraps these CLIs.

## Modifiers (v6.76.0+)

Order-independent tokens that appear **before** the URL. Same
semantics as `cfr`; see `cfr --help` for the full modifier table.
`cfrp` implicitly behaves as `cfr p …`, so passing `p` explicitly is
a no-op.

| Token | Effect |
|-------|--------|
| 🧭 `cg` | After clone + fix-repo + make-public, run the OS-appropriate **Coding Guidelines v24** installer, then auto-commit + push. |
| 🌍 `p`  | No-op on `cfrp` (already public). Accepted for parity with `cfr`. |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| 🔐 `--ssh` / `-ssh` / `--sh` | false | Force the URL into `git@host:owner/repo.git` SSH-shorthand form before clone runs. Auto-converts `https://…` and `ssh://git@…` URLs. Mutually exclusive with `--https` (`--ssh` wins with a one-line stderr warning). |
| 🌐 `--https` / `-https` / `--ht` | false | Force the URL into `https://host/owner/repo.git` form. Converts SSH-shorthand and `ssh://…` URLs. Useful in CI where the SSH agent isn't unlocked. |
| 🚫 `--no-vscode-sync` | false | Forwarded to the `clone` step — skips writing the resolved folder into VS Code Project Manager `projects.json`. |
| 🤐 `--yes` / `-y` | false | Non-interactive: skip the prior-version privatize prompt and auto-confirm any chained `make-public` confirmation. |
| 🔒 `--require-version` | false | Strict mode: fail (exit 4) when the cloned repo identity has no `-vN` suffix. |
| 👁️ `--dry-run` / `-n` | false | Preview only — prints the chained sequence that *would* run without touching remote visibility. |
| 🙈 `--no-commit` | false | (Only with `cg`.) Skip the auto-commit after the guidelines install. Files stay staged. |
| 📵 `--no-push`   | false | (Only with `cg`.) Commit locally but do not push. Also implicit when no upstream is set. |

Path canonicalization (Clean + EvalSymlinks for Windows 8.3 short
names and symlinks, with soft-fail to the cleaned absolute path on
resolver error) is inherited from the forwarded `clone` step.

## Behavior

1. 📥 **Clone** — versioned URLs auto-flatten. `--ssh` / `--https` rewrite the URL before clone runs.
2. 📂 **cd** — chdirs into the resolved folder.
3. 🔧 **fix-repo** — re-execs `fix-repo --all`. Skipped when the repo identity has no `-vN` suffix, unless `--require-version` is set.
4. 🌍 **make-public** — re-execs `make-public --yes` (non-interactive).
5. 🧭 **coding-guidelines** (only when `cg` modifier is present) — dispatches the v24 installer, then stages + commits + pushes. `--no-commit` / `--no-push` opt out of the commit / push sub-steps.

> **v6.50.0+** — `cfrp` no longer scans sibling `-vN` repos nor prompts to privatize prior versions after `make-public`. Run `gitmap mapri <repo>` explicitly when you want bulk-privatize behavior.

Also (v5.61.0+) — if the user's shell cwd is already inside the
target folder, `cfrp` chdir's to the parent before re-cloning so the
Windows file-handle lock never blocks the remove step.

Each step's exit code is propagated as-is; the pipeline halts on
the first non-zero exit.

## Examples

```
# Clone, fix tokens, expose publicly
gitmap clone-fix-repo-pub https://github.com/acme/myrepo-v13.git

# 🔐 Coerce HTTPS URL to SSH transport, then fix + publish
gitmap cfrp https://github.com/acme/myrepo-v13.git --ssh

# 🌐 Coerce SSH URL to HTTPS (CI without SSH agent)
gitmap cfrp git@github.com:acme/myrepo-v13.git --https

# Explicit destination folder
gitmap cfrp git@github.com:acme/myrepo-v13.git myrepo-fresh

# 🧭 Publish + install coding-guidelines v24 (auto-commit + push)
gitmap cfrp cg https://github.com/acme/myrepo-v13.git

# 🧭 Publish + install guidelines, but skip the auto-push
gitmap cfrp cg https://github.com/acme/myrepo-v13.git --no-push
```

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | ✅ ok |
| `6`  | ❌ bad-flag (missing URL) |
| `9`  | ❌ chdir failed |
| `10` | ❌ chained step failed (forwards underlying step's exit code) |

## See also

- `gitmap clone-fix-repo` (`cfr`) — same pipeline, without the visibility flip.
- `gitmap clone` — the underlying clone step.
- `gitmap make-public` / `gitmap fix-repo` — the individual steps.

## Scripting (JSON)

`gitmap help --json --filter clone-fix-repo-pub` — schema at
`spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

