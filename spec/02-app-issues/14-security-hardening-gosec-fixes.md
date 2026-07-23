# Post-Mortem: Security Hardening — G305, G110, Format Verb Audit

## Issues Fixed (v2.54.1–v2.54.3)

### 1. G305 — Zip Path Traversal (`installnpp.go`)

**Risk**: A crafted zip archive could write files outside the target directory using `../` sequences in entry names.

**Fix**: `extractZipEntry` now resolves the destination path to an absolute path and validates it starts with the target directory prefix before extraction.

```go
// ✅ Required: validate zip entry destination stays within target
destPath := filepath.Join(targetDir, entry.Name)
absTarget, _ := filepath.Abs(targetDir)
absDest, _ := filepath.Abs(destPath)
if !strings.HasPrefix(absDest, absTarget+string(os.PathSeparator)) {
    return fmt.Errorf("path traversal detected: %s", entry.Name)
}
```

### 2. G110 — Decompression Bomb (`installnpp.go`)

**Risk**: A malicious zip entry could expand to consume all available memory/disk.

**Fix**: Replaced `io.Copy` with `io.LimitReader` capped at 10 MB per extracted file.

```go
// ✅ Required: cap extraction size to prevent decompression bombs
limited := io.LimitReader(rc, 10*1024*1024) // 10 MB max
_, err = io.Copy(outFile, limited)
```

### 3. Format Verb Mismatch (`cmd/tasksync.go`)

**Risk**: `fmt.Fprintf` with mismatched format arguments causes `go vet` failures and potential runtime panics.

**Fix**: Corrected argument count at `tasksync.go:138`.

**Audit**: All ~140 `fmt.Fprintf`, `fmt.Printf`, `fmt.Errorf` call sites across `cmd/`, `release/`, `store/` verified — 100% compliance.

### 4. Code Red Error Audit (v2.54.1)

**Risk**: Generic error messages like "file not found" without paths make debugging impossible in production.

**Fix**: Standardized format across 35+ constants and 36+ call sites:

```
Error: [message] at [path]: [error] (operation: [op], reason: [reason])
```

## Prevention Rules

1. **All zip extraction must validate paths** — no extracted file may escape the target directory. This is a mandatory security check.
2. **All `io.Copy` from untrusted sources must use `io.LimitReader`** — cap at a reasonable maximum (10 MB default for CLI tools).
3. **All `fmt.*f` calls must match argument count to format verbs** — enforce via `go vet` in CI.
4. **All file/path error messages must include the exact resolved path** — generic messages are prohibited.
5. **Run `gosec` on every CI build** — do not suppress findings without inline justification comments.

## Related

- `.golangci.yml` — gosec exclusion rules with inline comments
- `spec/05-coding-guidelines/08-security-secrets.md` — Security & Secrets
- `spec/02-app-issues/27-error-management-file-path-and-missing-file-code-red-rule.md` — Code Red Rule
- CHANGELOG.md v2.54.1, v2.54.2, v2.54.3
