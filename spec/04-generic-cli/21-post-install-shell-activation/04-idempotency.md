# Post-Install Shell Activation — Idempotency

> **Parent spec:** [../21-post-install-shell-activation.md](../21-post-install-shell-activation.md)
> **Sibling files:**
> - [01-contract.md](01-contract.md) — Required behaviours and activation flow
> - [02-snippets.md](02-snippets.md) — Per-shell profile snippet bodies
> - [03-doctor.md](03-doctor.md) — `doctor` wrapper status detection

## Idempotency Rules

- Profile rewrites MUST be safe to run repeatedly (`setup`, `update`,
  CI provisioning).
- The marker comment is the **only** legal anchor for rewrites. Tools
  MUST NOT use line counts, absolute offsets, or content matching of
  the wrapper body.
- Bumping the version (`v2` → `v3`) MUST first remove every previous
  marker block before injecting the new one.
- Removing the wrapper (`<tool> uninstall --shell-wrapper`) MUST
  delete the marker block in full and leave surrounding content
  byte-identical.

---

## Rewrite Algorithm

1. Read profile file as UTF-8 text (or create if missing).
2. Locate every line matching `# <tool> shell wrapper v<N>` for any
   `<N>`.
3. For each match, find the matching `# <tool> shell wrapper v<N> end`
   line. Delete the inclusive range.
4. Append the new snippet block (current version) to the end of the
   file with a single leading blank line for readability.
5. Write the file atomically (write to temp + rename).

If step 3 cannot find a matching `end` marker, abort with a clear
error: `Profile contains an unterminated <tool> wrapper block. Edit
manually and re-run setup.` Do NOT attempt to guess the boundary.

---

## Removal Algorithm

Identical to step 1–3 of the rewrite algorithm, but skip step 4 and 5
write the file unconditionally. Surrounding lines (blank lines,
unrelated content) MUST be preserved byte-for-byte.

---

## Testing Requirements

Every CLI implementing this contract MUST cover:

- Snippet injection on a **fresh** profile (file does not exist).
- Snippet injection on a profile that already contains other content.
- Snippet **re-injection** (running `setup` twice produces no
  duplicates and a byte-identical file).
- Version bump (`v2` → `v3`) removes the old block and adds the new
  one in the same write.
- Marker-based **removal** leaves surrounding content byte-identical.
- `doctor` returns the correct status for each of the three states
  (LOADED, INSTALLED_BUT_NOT_LOADED, NOT_INSTALLED).

See [12-testing.md](../12-testing.md) for the project-wide testing
conventions.
