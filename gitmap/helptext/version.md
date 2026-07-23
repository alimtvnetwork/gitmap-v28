# gitmap version

Show the current gitmap version number.

## Alias

v

## Usage

    gitmap version

## Flags

None.

## Prerequisites

- None

## Examples

### Example 1: Show version

    gitmap version

**Output:**

    gitmap v2.22.0
    Built:  2025-03-10
    Go:     go1.22.1
    OS:     windows/amd64

### Example 2: Using alias

    gitmap v

**Output:**

    gitmap v2.22.0

## See Also

- [update](update.md) — Update gitmap to the latest version
- [doctor](doctor.md) — Diagnose installation issues

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter version
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
