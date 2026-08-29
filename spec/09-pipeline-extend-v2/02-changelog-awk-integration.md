# Changelog `awk` Integration

This repository relies on `.github/workflows/release.yml` to automatically construct GitHub Release objects. The release body (the "Notes") is extracted natively from `changelog.md` using `awk`.

## Extraction Mechanism

The workflow uses the following `awk` script to scrape `changelog.md`:

```bash
changelog=$(awk -v tag="[${GITHUB_REF_NAME}]" '
  $0 ~ "^## " tag {flag=1; next}
  /^## \[v/ && flag {flag=0; exit}
  flag {print}
' changelog.md)
```

### How it Works:

1. `GITHUB_REF_NAME` will evaluate to the pushed tag (e.g., `v6.87.1`).
2. The `awk` script searches for a line matching `## [v6.87.1]`.
3. It toggles a `flag=1` and begins printing every subsequent line.
4. If it encounters the next version block (`/^## \[v/`), it toggles `flag=0` and exits.

## AI Implementation Mandate

If an AI agent bumps a version, they MUST inject `## [vX.Y.Z] YYYY-MM-DD` at the very top of the `changelog.md` file. If the tag name does not match the header precisely, the `awk` script will quietly output `"No changelog entry found"` and the release will contain zero patch notes.
