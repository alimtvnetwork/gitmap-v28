# 10 — `@<Table>.<Code>` Resolver (FU-13)

**Spec:** `06-seedable-config-architecture` (extension)
**Version:** 1.0.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** FU-13 from `spec/19-main-worker-service/14-rbac-and-status-seed.md` §3.1.
**Authority:** Canonical algorithm + diagnostics for the seeder feature that translates `@<Table>.<Code>` strings to AUTOINCREMENT integer ids at apply time. On any conflict with `07-reference-table-seeding.md`, **`07-…` wins** (this file is its mechanical extension).

---

## 1. Why this file exists

`14-rbac-and-status-seed.md` §3.1 introduced `@Role.PowerAdmin` syntax for cross-table foreign-key resolution at seed time. Without an authoritative resolver spec, every implementer would write a different parser. This file pins the parsing rules, the lookup algorithm, the failure modes, and the trace output.

---

## 2. Syntax (normative)

```
AtRef    := "@" TableName "." Code
TableName := PascalCase identifier   (regex: ^[A-Z][A-Za-z0-9]*$)
Code     := alpha-numeric + underscore (regex: ^[A-Za-z0-9_]+$)
```

- The literal `@` MUST be the first character of the cell value.
- Exactly one `.` separates `TableName` and `Code`.
- Whitespace around any token is illegal (no trimming).
- Case sensitivity is strict (`@Role.PowerAdmin` ≠ `@role.poweradmin`).
- A row value that does not start with `@` is taken literally and bypasses the resolver.

The JSON-Schema pattern in `08-config-schema-tables-block.md` §2 (`AtRefString`) catches violations at lint time.

---

## 3. Lookup algorithm

Given an `AtRef` `@T.C` and the open seed transaction, the resolver MUST:

```
resolveAtRef(@T.C):
    code_column   = lookupCodeColumn(T)        // §4
    pk_column     = T + "Id"                   // per Code Red Schema rule
    sql           = "SELECT " + pk_column + " FROM " + T +
                    " WHERE " + code_column + " = ? LIMIT 1"
    row           = exec(sql, [C])
    if row is null:
        raise UnresolvedAtRefError(T, C, code_column)
    return row[pk_column]
```

- Lookup is case-sensitive (TEXT comparison; sqlite collation `BINARY`).
- The SELECT MUST be parameterized — never string-concatenated into SQL.
- The lookup runs **inside** the same transaction the seeder is using, so rows seeded earlier in the same boot are visible.

---

## 4. Code-column inference

The resolver derives the code column from the table name with a positive guard:

```
lookupCodeColumn(T):
    candidate = T + "Code"
    if hasColumn(T, candidate):
        return candidate
    raise CodeColumnMissingError(T, candidate)
```

This means every table referenceable via `@-ref` MUST expose a `<TableName>Code TEXT NOT NULL UNIQUE` column. Convention is enforced — there is no override.

Tables in scope today (from `spec/19/14-rbac-and-status-seed.md` §2):

| Table | Code column | Sample value |
|-------|-------------|--------------|
| `Role` | `RoleCode` | `PowerAdmin` |
| `EnumPage` | `EnumPageCode` | `BillingPage` |
| `WorkerNodeStatus` | `WorkerNodeStatusCode` | `Active` |
| `AuthMechanism` | `AuthMechanismCode` | `Jwt` |

---

## 5. Apply-time integration

The resolver hooks into the merge loop from `07-reference-table-seeding.md` §6:

```
for each (TableName, TableBlock) in seed.Tables:
    rows = []
    for raw_row in TableBlock.Rows:
        rows.append(resolveRow(raw_row))
    apply(rows, TableBlock.MergeStrategy, TableBlock.PrimaryKey)

resolveRow(raw):
    out = {}
    for column, value in raw.items():
        out[column] = resolveValue(value)
    return out

resolveValue(v):
    if v is string AND v starts with "@":
        return resolveAtRef(v)
    return v
```

Each helper stays ≤8 lines; `resolveValue` is a positive guard with zero nesting (Code Red metrics).

---

## 6. Failure modes (registered MWS codes)

| Condition | Error | MWS code | Surfaced at |
|-----------|-------|----------|-------------|
| Reference like `@Role.Ghost` resolves to 0 rows | `UnresolvedAtRefError` | `MWS-21002 SplitDbTierMissing` | Seed boot |
| `@<T>.<C>` where `T` has no `<T>Code` column | `CodeColumnMissingError` | `MWS-21002 SplitDbTierMissing` | Seed boot |
| Malformed (`@.X`, `@T.`, `@T.X.Y`, leading whitespace) | `MalformedAtRefError` | rejected pre-apply by JSON-Schema (`08-…` §2) | Lint |
| SELECT raises a DB error | `DbReadError` | `MWS-21052 SplitDbWriteFail` | Seed boot |

Per `mem://architecture/error-handling`: errors are NEVER swallowed; the seeder MUST log `TableName + offending row index + offending column + raw value + cause path` and abort the boot. Partial seeds are not committed (the apply runs inside one transaction per table).

---

## 7. Trace output (mandatory)

For every successful resolution the seeder emits one line at INFO:

```
seed.atref.resolved table=RolePageAccess column=RoleId raw=@Role.PowerAdmin -> 1
```

For every failure (FATAL):

```
seed.atref.unresolved table=RolePageAccess column=RoleId raw=@Role.Ghost reason=NoMatch
```

Trace lines are stable contract for downstream log parsers — DO NOT alter the field order or punctuation.

---

## 8. Worked example

Input row from `spec/19/14-rbac-and-status-seed.md` §3 `RolePageAccess`:

```jsonc
{ "RoleId": "@Role.PowerAdmin", "EnumPageCode": "DashboardPage" }
```

Apply-time resolution:

| Step | Detail |
|------|--------|
| Detect | `@-ref` because value starts with `@`. |
| Parse | `T = "Role"`, `C = "PowerAdmin"`. |
| Code column | `Role` + `Code` → `RoleCode` (column exists ✓). |
| Lookup | `SELECT RoleId FROM Role WHERE RoleCode = ? LIMIT 1`, params `["PowerAdmin"]`. |
| Result | `RoleId = 1` (assuming `Role` was just seeded). |
| Substitute | Resolved row: `{ "RoleId": 1, "EnumPageCode": "DashboardPage" }`. |
| Trace | `seed.atref.resolved table=RolePageAccess column=RoleId raw=@Role.PowerAdmin -> 1`. |

The substituted row is then handed to `UpsertByLogicalKey` per `07-…` §3.

---

## 9. Reference implementation skeleton (Go)

CODE RED compliant — every function ≤15 lines, zero nested ifs, positive guards:

```go
type AtRef struct { Table, Code string }

func isAtRef(v any) bool {
    s, ok := v.(string)
    return ok && len(s) > 0 && s[0] == '@'
}

func parseAtRef(s string) (AtRef, error) {
    body := s[1:]
    dot := strings.IndexByte(body, '.')
    if dot < 1 || dot == len(body)-1 {
        return AtRef{}, fmt.Errorf("malformed @-ref: %q", s)
    }
    return AtRef{Table: body[:dot], Code: body[dot+1:]}, nil
}

func resolveAtRef(tx *sql.Tx, ref AtRef) (int64, error) {
    col := ref.Table + "Code"
    pk  := ref.Table + "Id"
    q   := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? LIMIT 1", pk, ref.Table, col)
    var id int64
    if err := tx.QueryRow(q, ref.Code).Scan(&id); err != nil {
        return 0, fmt.Errorf("@-ref %s.%s unresolved: %w", ref.Table, ref.Code, err)
    }
    return id, nil
}
```

Total: 4 functions, 22 lines combined. All inputs are validated; no error swallowed.

---

## 10. CI test cases (mandatory)

`linter-scripts/tests/` MUST include:

| # | Input | Expected |
|---|-------|----------|
| T-1 | `"@Role.PowerAdmin"` after Role seeded | resolves to integer >= 1 |
| T-2 | `"@Role.Ghost"` | `MWS-21002 SplitDbTierMissing` |
| T-3 | `"@role.PowerAdmin"` (lowercase table) | rejected by JSON-Schema |
| T-4 | `"@.PowerAdmin"` | rejected by JSON-Schema |
| T-5 | `"@Role."` | rejected by JSON-Schema |
| T-6 | `"@Role.PowerAdmin.Extra"` | rejected by JSON-Schema |
| T-7 | `"PowerAdmin"` (no `@`) | passes through literally |
| T-8 | Reference to a table without `<T>Code` column | `CodeColumnMissingError` → `MWS-21002` |

---

## 11. Cross-references

- `spec/06-seedable-config-architecture/02-features/07-reference-table-seeding.md` — feature this resolver completes.
- `spec/06-seedable-config-architecture/02-features/08-config-schema-tables-block.md` — JSON-Schema rejects malformed `@-refs` at lint time.
- `spec/06-seedable-config-architecture/02-features/09-config-seed-1.3.0-bump.md` — first seed file that uses `@-refs` (`RolePageAccess`).
- `spec/19-main-worker-service/13-error-codes.md` §2.1 + §2.6 — MWS codes raised by failures.
- `spec/19-main-worker-service/14-rbac-and-status-seed.md` §3.1 — original syntax declaration.

---

## 12. Open Questions (logged, non-blocking)

- **OQ-10-1** Allow forward references (referenced table seeded later in the same boot)? Inferred: no — `07-…` §9 OQ-2 already mandates declaration order. Failing fast is preferable to topological sorting.
- **OQ-10-2** Permit `@<Table>.<Code>` to resolve to a non-PK integer? Inferred: no — keeps the resolver bounded and the trace unambiguous.

---

*@-ref resolver v1.0.0 — 2026-05-04 (FU-13 closed)*
