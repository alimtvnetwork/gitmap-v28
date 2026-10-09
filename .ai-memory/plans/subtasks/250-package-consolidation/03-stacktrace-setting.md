# Subtask 03 — stack-trace setting

> Spec: `02-spec/21-app/250-package-consolidation/03-stdout-and-stacktrace.md`
> ("Stack-trace setting" section). Sequencing: independent of the package
> merges; can run any time.

## Checkboxes

- [ ] Add `ShowStackTrace bool \`json:"showStackTrace"\`` to `model.Config`
      in `cli/model/record.go`, directly after the `ErrorDisplay` field
      (naming precedent: `CommitReplayKeepUrl`).
- [ ] Default it to `true` in `DefaultConfig()` — default ON preserves current
      behavior; missing keys keep the default via the existing
      unmarshal-onto-defaults in `parseConfig`.
- [ ] In `handleGlobalError` (`cli/cmd/root.go`, ~lines 261-264): resolve a
      local `showStack` defaulting to `true`, set it from
      `cfg.ShowStackTrace` only when `cfgErr == nil` (mirror the existing
      `display` resolution in the same function), and gate the trace print:

      ```go
      stack := resolveErrorStackTrace(err)
      if stack != "" && showStack {
          fmt.Fprintf(os.Stderr, "Stack Trace:%s\n", stack)
      }
      ```

- [ ] Leave the `errorDisplay == "simple"` early exit untouched — the two
      controls stay orthogonal (format vs. trace detail).
- [ ] No `validate.go` change (a bool needs no validation).

## Manual verification (lead-run binary)

1. Build a fresh binary from this branch.
2. From a directory containing `./data/config.json` (the config path is
   cwd-relative — run from the repo root), trigger a global error that reaches
   `handleGlobalError` with the default config: the `Stack Trace:` block MUST
   print (current behavior preserved).
3. Set `"showStackTrace": false` in `./data/config.json`, re-run: the error
   message prints but the `Stack Trace:` block MUST NOT.
4. Set `"showStackTrace": true`, re-run: the trace prints again.
5. Set `"errorDisplay": "simple"`, re-run: the simple-path behavior is
   unchanged (early exit as before).
6. Run the existing tests scoped to `cli/cmd` and `cli/model` only (no routine
   full test runs).
