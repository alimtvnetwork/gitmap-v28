# CI Issue: `go-winres` Icon Size Limit

## Error

```
2026/04/16 16:26:46 image size too big, must fit in 256x256
Error: Process completed with exit code 1
```

## Root Cause

`go-winres` embeds icons as Windows `.ico` resources inside the compiled `.exe`.
The Windows `.ico` format has a **hard limit of 256x256 pixels** per image frame.
The project's `icon.png` was **512x512**, which exceeds this limit.

The error originates from the `go-winres make` step in the CI release workflow
(`.github/workflows/release.yml`, line 54), which reads `gitmap-v28/winres/winres.json`
and attempts to convert the referenced PNG into an `.ico` resource.

## Why This Wasn't Caught Earlier

- Local builds on Windows using `run.ps1` do **not** run `go-winres` — the icon
  embedding only happens in the CI release pipeline.
- The `go build` command succeeds without `go-winres` — it simply produces a
  binary without embedded Windows metadata (icon, version info, manifest).

## Solution

1. Created a 256x256 resized copy: `gitmap-v28/assets/icon-256.png`.
2. Updated `gitmap-v28/winres/winres.json` to reference `icon-256.png` instead.
3. Kept the original 512x512 `icon.png` for other uses (docs site, README, etc.).

## Prevention

- Any icon referenced in `winres.json` **must** be ≤256x256.
- If a higher-resolution source icon is needed for other purposes, maintain
  separate files (e.g., `icon.png` for web, `icon-256.png` for Windows resources).
- Consider adding a pre-check in CI before `go-winres make`:
  ```bash
  python3 -c "from PIL import Image; img=Image.open('assets/icon-256.png'); assert max(img.size)<=256, f'Icon too large: {img.size}'"
  ```

## Related Files

| File | Purpose |
|------|---------|
| `gitmap-v28/winres/winres.json` | Windows resource manifest for `go-winres` |
| `gitmap-v28/assets/icon-256.png` | 256x256 icon for `.exe` embedding |
| `gitmap-v28/assets/icon.png` | 512x512 original icon (web/docs use) |
| `.github/workflows/release.yml` | CI pipeline that runs `go-winres make` |

## Cross-References

- [08-repo-path-sync.md](08-repo-path-sync.md) — Related CI pipeline spec
- [04-build-scripts.md](04-build-scripts.md) — Build pipeline documentation
