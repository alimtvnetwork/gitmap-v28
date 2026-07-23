# 05 — Release Assets

## Purpose

Define the naming conventions, compression formats, and packaging rules
for release artifacts published to GitHub Releases.

---

## Asset Naming Convention

```
<binary>-<os>-<arch>.<ext>
```

| Component | Values |
|-----------|--------|
| `<binary>` | The CLI tool name (lowercase, hyphenated) |
| `<os>` | `windows`, `linux`, `darwin` |
| `<arch>` | `amd64`, `arm64` |
| `<ext>` | `.zip` (Windows), `.tar.gz` (Unix) |

### Examples

```
mytool-linux-amd64.tar.gz
mytool-darwin-arm64.tar.gz
mytool-windows-amd64.zip
```

---

## Compression

### Windows → `.zip`

```bash
zip "<binary>-${os}-${arch}.zip" "<binary>-${os}-${arch}.exe"
```

**Why zip**: Natively supported by PowerShell (`Expand-Archive`) and
Windows Explorer. No additional tools required.

### Linux / macOS → `.tar.gz`

```bash
tar czf "<binary>-${os}-${arch}.tar.gz" "<binary>-${os}-${arch}"
```

**Why tar.gz**: Preserves Unix file permissions (executable bit).
Universally available on all Unix systems.

---

## Complete Asset List Per Release

A standard release publishes these files:

```
dist/
├── <binary>-windows-amd64.zip
├── <binary>-windows-arm64.zip
├── <binary>-linux-amd64.tar.gz
├── <binary>-linux-arm64.tar.gz
├── <binary>-darwin-amd64.tar.gz
├── <binary>-darwin-arm64.tar.gz
├── checksums.txt
├── install.ps1
└── install.sh
```

Total: **9 files** (6 archives + checksums + 2 install scripts).

---

## Raw Binary Cleanup

After compression, remove the raw binaries from `dist/` so they are
not accidentally published:

```bash
cd dist
for f in <binary>-*; do
    [[ "$f" != *.zip && "$f" != *.tar.gz && "$f" != *.txt && "$f" != *.ps1 && "$f" != *.sh ]] && rm "$f"
done
```

---

## Archive Contents

Each archive contains a **single file** — the binary itself. No nested
directories, no READMEs, no LICENSE files inside the archive.

```bash
# Correct: flat binary in archive
tar czf mytool-linux-amd64.tar.gz mytool-linux-amd64

# Wrong: nested directory structure
tar czf mytool-linux-amd64.tar.gz mytool/mytool
```

This simplifies extraction and ensures install scripts can find the
binary at a predictable path.

---

## Multiple Binaries Per Release

If the project produces multiple CLI tools (e.g., a main tool and an
updater), build and package each independently with the same naming
convention:

```
dist/
├── mytool-linux-amd64.tar.gz
├── mytool-windows-amd64.zip
├── mytool-updater-linux-amd64.tar.gz
├── mytool-updater-windows-amd64.zip
├── checksums.txt           ← covers ALL archives
├── install.ps1
└── install.sh
```

The `checksums.txt` file contains hashes for **all** archives, not
just the primary binary.

---

## File Size Reporting

After building, print human-readable file sizes for the release summary:

```bash
for f in dist/*; do
    size=$(stat -c%s "$f" 2>/dev/null || stat -f%z "$f")
    size_mb=$(echo "scale=2; $size / 1048576" | bc)
    echo "  $(basename "$f"): ${size_mb} MB"
done
```

---

## Constraints

- Asset names must be lowercase with hyphens (no underscores or spaces).
- Archives contain a single flat binary — no nested directories.
- Raw binaries must be removed from `dist/` after compression.
- One `checksums.txt` covers all archives in the release.
- Install scripts are version-pinned at generation time.
- All assets are uploaded in a single `softprops/action-gh-release` step.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
