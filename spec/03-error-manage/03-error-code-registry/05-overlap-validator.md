# 02 — Error Code Registry Overlap Validator

> **Created:** 2026-02-28  
**Version:** 3.2.0  
> **Status:** Spec Ready  
> **Owner:** SM CLI tooling  
> **Cross-Reference:** `spec/03-error-manage/03-error-code-registry/01-registry.md`

---

## Purpose

Automated tool that parses the master error code registry (`01-registry.md`) and detects overlapping or duplicate range allocations. Prevents the class of collision bugs documented in `spec/23-how-app-issues-track/16-error-code-collision-remediation.md`.

---

## Input

The validator reads `spec/03-error-manage/03-error-code-registry/01-registry.md` and extracts all rows from the **Registered Project Prefixes** table (Section 1) and all **Standalone Specification Error Ranges** detail tables.

### Parsing Rules

1. Match markdown table rows with the pattern: `| {Prefix} | {Project} | {RangeStart}-{RangeEnd} | {SpecLocation} | {Status} |`
2. Skip rows where Status contains `Deprecated`
3. Parse `RangeStart` and `RangeEnd` as integers
4. Rows using `GEN-xxx-xx` format (e.g., `CAST`) are parsed as symbolic ranges and excluded from numeric overlap detection
5. Ignore table header and separator rows (`|---`)

---

## Overlap Detection Algorithm

### Parameters

```
Input:  registryPath string  (path to 01-registry.md)
Output: OverlapResult        (list of findings)
```

### Data Structure

```
RangeEntry {
    Prefix     string   // e.g., "SM-CG"
    Project    string   // e.g., "SM Code Generation"
    RangeStart int      // e.g., 16000
    RangeEnd   int      // e.g., 16799
    SpecPath   string   // e.g., "spec/02-spec-management-software/..."
    Line       int      // source line number in registry file
}

OverlapFinding {
    EntryA     RangeEntry
    EntryB     RangeEntry
    OverlapStart int    // first overlapping code
    OverlapEnd   int    // last overlapping code
    Severity   string   // "COLLISION" | "INTENTIONAL"
}
```

### Steps

1. **Parse:** Extract all `RangeEntry` records from the registry markdown
2. **Sort:** Order entries by `RangeStart` ascending, then by `RangeEnd` descending (wider ranges first)
3. **Sweep:** For each pair `(i, j)` where `j > i`:
   - If `entries[j].RangeStart <= entries[i].RangeEnd`:
     - Overlap detected: `OverlapStart = entries[j].RangeStart`, `OverlapEnd = min(entries[i].RangeEnd, entries[j].RangeEnd)`
     - Check intentional overlap allowlist (see below)
     - Emit `OverlapFinding`
4. **Report:** Return all findings sorted by severity (COLLISION first), then by `OverlapStart`

### Intentional Overlap Allowlist

Some overlaps are documented as intentional and format-separated (see Resolution 9 in the registry). The validator must maintain an allowlist:

| Entry A | Entry B | Reason |
|---------|---------|--------|
| `PS` 9500-9599 | `AB` 9500-9540 | Format-separated: PS uses `PS-9500-00` strings, AB uses flat integers |

Matches against the allowlist are reported as `Severity: INTENTIONAL` rather than `COLLISION`.

---

## Output Format

### Console Output

```
=== Error Code Registry Overlap Validator ===
Registry: spec/03-error-manage/03-error-code-registry/01-registry.md
Entries parsed: {N}
Active entries: {N} (excluded {N} deprecated)

--- Findings ---

[COLLISION] Range 12000-12599 overlap between:
  - WSP (WP SEO Publish) at line 48
  - SM-CG (SM Code Generation) at line 53
  Overlap: 12000-12599 (600 codes)

[INTENTIONAL] Range 9500-9540 overlap between:
  - PS (PowerShell Integration) at line 28
  - AB (AI Bridge SEO) at line 29
  Overlap: 9500-9540 (41 codes) — allowlisted

--- Summary ---
Total overlaps: {N}
  Collisions: {N}  ← MUST be 0 for registry to be valid
  Intentional: {N}
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No collisions detected (intentional overlaps OK) |
| 1 | One or more COLLISION findings |
| 2 | Parse error (registry file malformed or not found) |

---

## Integration Points

### When to Run

1. **Before any error code allocation:** Run validator after editing `01-registry.md`
2. **CI/CD gate:** Block merge if exit code is non-zero
3. **Remediation verification:** Run after any collision resolution to confirm zero remaining overlaps

### CLI Interface

```
spec-tools registry validate [--registry path/to/01-registry.md] [--format json|text]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--registry` | `spec/03-error-manage/03-error-code-registry/01-registry.md` | Path to registry file |
| `--format` | `text` | Output format (`text` for console, `json` for CI) |
| `--strict` | `false` | When true, intentional overlaps also cause exit code 1 |

### JSON Output Schema

```json
{
  "RegistryPath": "string",
  "EntriesParsed": "int",
  "ActiveEntries": "int",
  "Findings": [
    {
      "Severity": "COLLISION | INTENTIONAL",
      "OverlapStart": "int",
      "OverlapEnd": "int",
      "OverlapSize": "int",
      "EntryA": {
        "Prefix": "string",
        "Project": "string",
        "RangeStart": "int",
        "RangeEnd": "int",
        "SpecPath": "string",
        "Line": "int"
      },
      "EntryB": { "..." }
    }
  ],
  "Summary": {
    "TotalOverlaps": "int",
    "Collisions": "int",
    "Intentional": "int",
    "Valid": "bool"
  }
}
```

---

## Error Codes

> **Range:** SM 2850-2859 (within SM 2000-2999, adjacent to SM-RT 2800-2849)

| Code | Constant | Description |
|------|----------|-------------|
| 2850 | `ERR_REGISTRY_NOT_FOUND` | Registry file not found at specified path |
| 2851 | `ERR_REGISTRY_PARSE_FAILED` | Failed to parse markdown table structure |
| 2852 | `ERR_REGISTRY_RANGE_INVALID` | Range entry has start > end or non-numeric values |
| 2853 | `ERR_REGISTRY_DUPLICATE_PREFIX` | Same prefix registered with conflicting ranges (outside allowlist) |

---

## Acceptance Criteria

- **AC-1:** GIVEN a registry with zero overlaps, WHEN the validator runs, THEN exit code is 0 and no COLLISION findings are emitted
- **AC-2:** GIVEN a registry with an intentional PS/AB overlap, WHEN the validator runs without `--strict`, THEN exit code is 0 and one INTENTIONAL finding is emitted
- **AC-3:** GIVEN a registry where SM-CG still uses 12xxx (pre-fix state), WHEN the validator runs, THEN exit code is 1 and a COLLISION finding identifies SM-CG vs WSP
- **AC-4:** GIVEN a malformed registry (missing table headers), WHEN the validator runs, THEN exit code is 2 with `ERR_REGISTRY_PARSE_FAILED`
- **AC-5:** GIVEN `--format json`, WHEN findings exist, THEN output is valid JSON matching the schema above
- **AC-6:** GIVEN `--strict` flag, WHEN intentional overlaps exist, THEN exit code is 1

---

## Known Pitfalls and Prevention

- **Pitfall:** Regex parsing may fail on table rows with extra whitespace or inline formatting (bold, backticks). Use a permissive pattern that strips markdown formatting before extracting values.
- **Pitfall:** Symbolic ranges like `GEN-600-01` must not be parsed as numeric. Filter by format before integer conversion.
- **Prevention:** See `spec/23-how-app-issues-track/16-error-code-collision-remediation.md` for the class of bugs this tool prevents.
