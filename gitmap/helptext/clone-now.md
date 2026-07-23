# gitmap clone-now (alias of `reclone`)

`clone-now` is a **legacy alias** of `gitmap reclone` — the canonical
command that re-runs `git clone` against `gitmap scan` artifacts.

The behavior, flags, exit codes, and arguments are identical. New
documentation, examples, and tab-completion all use the `reclone`
spelling. Existing scripts using `clone-now` (or `cnow`) will keep
working forever — there are no plans to remove this alias.

```
gitmap reclone <file> [flags]   # canonical
gitmap clone-now <file> [flags] # this alias
```

For the full reference, see:

```
gitmap help reclone
```

## Why the rename?

The CLI now has two clearly-distinct verbs:

- `gitmap clone <url>` — fresh clone of a single repo from a URL.
- `gitmap reclone <file>` — RE-clone an entire scanned tree from a
  `gitmap scan` artifact, restoring the recorded folder layout.

`reclone` makes the round-trip nature explicit, and avoids the older
ambiguity between "clone now" (sounds like immediate-mode `clone`)
and "re-clone from artifact". `clone-now` stays available so nothing
that already shipped breaks.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter clone-now
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap clone-now
```
