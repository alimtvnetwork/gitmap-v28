# SSH-identified repos cloned/pulled via HTTPS

**Status:** open
**Reported:** 2026-06-07

## Symptom
`gitmap scan` in D:\Work found 2 repos whose origin remotes are SSH, but the generated `clone` / direct-clone scripts and the background probe attempted HTTPS, triggering a browser auth prompt ("info: please complete authentication in your browser...").

## Expected
- Scan persists BOTH transport URLs (HTTPS + SSH) into SQLite for every repo.
- The repo's identified transport (the one actually configured on `origin`) is recorded as the canonical transport.
- All downstream operations (clone scripts, clone commands shown in terminal report, background probe, pull) use the identified transport — SSH repos must NOT silently fall back to HTTPS.
- Terminal report shows the identified transport AND both URL variants per repo.

## Actual
- Only one URL form surfaces in the report/scripts.
- Background probe and clone commands default to HTTPS even when origin is SSH.

## Related files (likely)
- `gitmap/mapper/transport.go`
- `gitmap/formatter/clonescript.go` (`cloneURL` picks HTTPS first)
- `gitmap/store/repo.go`, model `ScanRecord`
- `gitmap/probe/clone.go`
- `gitmap/cmd/scan*.go`, `gitmap/cmd/scanoutput.go`
