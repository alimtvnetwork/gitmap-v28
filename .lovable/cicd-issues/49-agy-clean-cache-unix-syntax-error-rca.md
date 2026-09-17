# CI/CD Issue 49: `agy_clean_cache_unix.go` Syntax Error RCA

## 1. Symptom
In GitHub Actions runs triggered by commit `6dadbd96` (Release `v6.257.1`), including:
- **Release (#35249825042)**: Step `Build release binaries` failed.
- **History Rewrite Smoke (#35249812579)**: Step `Build gitmap binary` failed.
- **race-detector (#35249812546)**: Step `go test -race` failed.

### Error Output
```text
cmdagy/agy_clean_cache_unix.go:168:2: syntax error: non-declaration statement outside function body
FAIL    github.com/alimtvnetwork/gitmap-v28/cli/cmd [build failed]
Process completed with exit code 1.
```

---

## 2. Root Cause
In `cli/cmdagy/agy_clean_cache_unix.go`, an incomplete refactor left orphaned trailing lines (`return killed, warnings` and a closing brace `}`) outside of any function body after line 164.

---

## 3. Resolution
Removed the orphaned statements from the tail of `cli/cmdagy/agy_clean_cache_unix.go`, ensuring all statements reside strictly within valid function declarations. Verified clean compilation across `GOOS=linux`, `GOOS=darwin`, and `GOOS=windows`.

---

## 4. Prevention & Learnings
- **Cross-Platform Pre-Release Build Gate**: Always run `GOOS=linux go build ./...` and `GOOS=darwin go build ./...` prior to cutting release tags to catch OS-specific build failures (`_unix.go`, `_windows.go`, `_darwin.go`) before CI workflows trigger.
- **Hermetic Linter Validation**: Ensure `golangci-lint` or `go vet` runs with explicit `GOOS=linux` during local quality sweeps.
