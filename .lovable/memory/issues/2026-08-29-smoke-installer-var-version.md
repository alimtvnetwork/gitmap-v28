# Memory Issue: 2026-08-29 Smoke Installer Var Version

## 1. Why it happened
`smoke-installer.sh` assumed `constants.Version` was declared with `const Version = "..."`, but in Go variables altered by `-ldflags -X` must be declared as `var Version = "..."`.

---

## 2. How it happened
When CI ran `smoke-installer.sh source`, `awk` searched for `^const Version` in `gitmap/constants/constants.go`. Because it found nothing, `$EXPECTED` was empty and the script threw `::error::Could not determine expected version`.

---

## 3. Root Cause
Rigid `^const Version` regex in `.github/scripts/smoke-installer.sh`.

---

## 4. Code Fix
Updated extraction to `grep -E '^(var|const) Version\b'` and piped to `sed -E 's/.*"([^"]+)".*/\1/'`. Added `Installer Smoke` check to local runner `03-cicd-local-runner.py`.
