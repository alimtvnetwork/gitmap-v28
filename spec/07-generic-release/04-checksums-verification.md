# 04 — Checksums & Verification

## Purpose

Ensure binary integrity between the CI build environment and the user's
machine. Every release artifact must have a corresponding SHA-256 hash
that is verified before installation.

---

## Checksum Generation

After compressing all binaries, generate a single `checksums.txt` file:

```bash
cd dist
sha256sum *.zip *.tar.gz > checksums.txt
```

### Output Format

```
a1b2c3d4...  <binary>-linux-amd64.tar.gz
e5f6a7b8...  <binary>-linux-arm64.tar.gz
c9d0e1f2...  <binary>-darwin-amd64.tar.gz
f3a4b5c6...  <binary>-darwin-arm64.tar.gz
d7e8f9a0...  <binary>-windows-amd64.zip
b1c2d3e4...  <binary>-windows-arm64.zip
```

The format is `<hash>  <filename>` (two spaces between hash and name),
which is the standard output of `sha256sum`.

---

## Checksum Verification — Installer Side

### PowerShell

```powershell
$expectedHash = (Get-Content checksums.txt |
    Where-Object { $_ -match $archiveName } |
    ForEach-Object { ($_ -split '\s+')[0] })

$actualHash = (Get-FileHash $archivePath -Algorithm SHA256).Hash.ToLower()

if ($actualHash -ne $expectedHash.ToLower()) {
    Write-Host "  XX Checksum verification failed!" -ForegroundColor Red
    Write-Host "     Expected: $expectedHash" -ForegroundColor Red
    Write-Host "     Got:      $actualHash" -ForegroundColor Red
    exit 1
}

Write-Host "  OK Checksum verified" -ForegroundColor Green
```

### Bash

```bash
expected=$(grep "$archive_name" checksums.txt | awk '{print $1}')
actual=$(sha256sum "$archive_path" | awk '{print $1}')

if [[ "$actual" != "$expected" ]]; then
    echo "  XX Checksum verification failed!" >&2
    echo "     Expected: $expected" >&2
    echo "     Got:      $actual" >&2
    exit 1
fi

echo "  OK Checksum verified"
```

On macOS, use `shasum -a 256` instead of `sha256sum`:

```bash
if command -v sha256sum &>/dev/null; then
    actual=$(sha256sum "$archive_path" | awk '{print $1}')
else
    actual=$(shasum -a 256 "$archive_path" | awk '{print $1}')
fi
```

---

## Checksum Download

The `checksums.txt` file is downloaded from the same release as the
binary archive:

```
https://github.com/<org>/<repo>/releases/download/v<version>/checksums.txt
```

Both the archive and checksums must be downloaded before verification
begins. If `checksums.txt` fails to download, the installer should
warn but NOT skip verification — it must exit with an error.

---

## Publishing Checksums

`checksums.txt` is published as a release asset alongside the archives:

```yaml
- uses: softprops/action-gh-release@v2
  with:
    files: |
      dist/*.zip
      dist/*.tar.gz
      dist/checksums.txt
      dist/install.ps1
      dist/install.sh
```

---

## Release Body Checksum Table

Include a formatted checksum table in the release description:

```markdown
## Checksums (SHA-256)

| File | Hash |
|------|------|
| `<binary>-linux-amd64.tar.gz` | `a1b2c3d4...` |
| `<binary>-windows-amd64.zip` | `d7e8f9a0...` |
```

---

## Constraints

- Use SHA-256 exclusively — not MD5 or SHA-1.
- Checksum verification must never be optional or skippable.
- If the checksum file cannot be downloaded, the install must fail.
- Generate checksums from the **compressed archives**, not the raw
  binaries — users verify what they download.
- Case-insensitive comparison (normalize to lowercase before comparing).
- The checksum file must be generated in the same CI job as the
  archives to prevent TOCTOU (time-of-check-time-of-use) issues.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
