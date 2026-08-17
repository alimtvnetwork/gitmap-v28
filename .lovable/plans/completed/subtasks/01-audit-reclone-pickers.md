# Audit of Reclone-class URL Pickers

- **`gitmap clone`** (`gitmap/cmd/clone.go:L548` in `applyURLSchemeFlags`): Hard-codes HTTPS fallback logic if `--ssh` is not passed. Doesn't read `IdentifiedTransport`.
- **`gitmap clone-fix-repo` / `cfr`** (`gitmap/cmd/clonefixrepo.go:L185` in `applyCloneFixRepoScheme`): Doesn't know about `IdentifiedTransport` yet, forces SSH/HTTPS only if flags are provided.
- **`gitmap clone-now` / `cn`** (`gitmap/clonenow/clonenow.go:L87` in `Row.PickURL`): Respects `--ssh` or `--https` mode flags, but ignores the row's `Transport` entirely.
- **`gitmap scan`** (`gitmap/cloner/summary.go:L41` in `pickURL`): **Honors** `rec.Transport == "ssh"`! This is the only one doing it right today.

Verdict: `cloner/summary.go` honors `IdentifiedTransport`, while `clone.go`, `clonefixrepo.go`, and `clonenow.go` still hard-code HTTPS-first logic because they don't read `IdentifiedTransport` from the DB row.
