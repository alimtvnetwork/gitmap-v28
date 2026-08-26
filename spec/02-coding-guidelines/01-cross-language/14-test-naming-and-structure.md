# Cross-Language Rule: Test Naming & Structure

> **Version:** 1.0.0  
> **Updated:** 2026-03-11  
> **Applies to:** Go, TypeScript, PHP

---

## Overview

This spec defines mandatory conventions for test file organization, test function naming, table-driven test structure, and test helper placement. All tests must be self-documenting — the test name alone must communicate **what is being tested**, **under what conditions**, and **what the expected outcome is**.

---

## Rule 1: Test File Naming

Test files must mirror the source file they test, with a language-appropriate test suffix.

| Language | Source File | Test File |
|----------|------------|-----------|
| Go | `SiteManager.go` | `SiteManager_test.go` |
| TypeScript | `UserProfile.tsx` | `UserProfile.test.tsx` |
| PHP | `FileOperations.php` | `FileOperationsTest.php` |

### Constraints

- One test file per source file — no multi-source test files
- Test file must reside in the same directory as the source file (Go: same package; TypeScript: colocated; PHP: mirrored `tests/` tree)
- Integration test files use the `_integration_test.go` / `.integration.test.tsx` suffix

```
// ❌ FORBIDDEN — generic test file covering multiple sources
tests/
  all_services_test.go

// ✅ REQUIRED — one test file per source
internal/session/
  SessionService.go
  SessionService_test.go
  SessionStore.go
  SessionStore_test.go
```

---

## Rule 2: Test Function Naming — Three-Part Convention

Every test function name must follow the pattern:

```
Test{Unit}_{Scenario}_{ExpectedOutcome}
```

| Segment | Purpose | Example |
|---------|---------|---------|
| `{Unit}` | The function, method, or component under test | `CreateSession` |
| `{Scenario}` | The condition or input variation | `WithExpiredToken` |
| `{ExpectedOutcome}` | What should happen | `ReturnsAuthError` |

### Go

```go
// ❌ FORBIDDEN — vague, no scenario, no outcome
func TestCreateSession(t *testing.T) { ... }
func TestSession(t *testing.T) { ... }
func Test_create(t *testing.T) { ... }

// ✅ REQUIRED — three-part naming
func TestCreateSession_WithValidCredentials_ReturnsSessionId(t *testing.T) { ... }
func TestCreateSession_WithExpiredToken_ReturnsAuthError(t *testing.T) { ... }
func TestCreateSession_WithMissingUserId_ReturnsValidationError(t *testing.T) { ... }
```

### TypeScript

```typescript
// ❌ FORBIDDEN — vague description
describe('UserProfile', () => {
    it('works', () => { ... });
    it('handles error', () => { ... });
});

// ✅ REQUIRED — three-part naming in describe/it
describe('UserProfile', () => {
    it('renders_WithValidUser_ShowsDisplayName', () => { ... });
    it('renders_WithMissingAvatar_ShowsDefaultIcon', () => { ... });
    it('onSubmit_WithInvalidEmail_ShowsValidationError', () => { ... });
});
```

### PHP

```php
// ❌ FORBIDDEN — vague
public function testProcess(): void { ... }

// ✅ REQUIRED — three-part naming
public function testProcessUpload_WithEmptyFile_ReturnsError(): void { ... }
public function testProcessUpload_WithValidPdf_SavesAndReturnsPath(): void { ... }
```

---

## Rule 3: Table-Driven Tests

When a function has **3 or more** test scenarios with the same setup/assertion structure, use table-driven tests. Each case must have a `name` field following the three-part convention (without the `Test` prefix).

### Go

```go
// ❌ FORBIDDEN — repetitive individual tests
func TestValidateEmail_WithEmptyString_ReturnsError(t *testing.T) { ... }
func TestValidateEmail_WithNoAtSign_ReturnsError(t *testing.T) { ... }
func TestValidateEmail_WithValidFormat_ReturnsNil(t *testing.T) { ... }

// ✅ REQUIRED — table-driven when 3+ similar scenarios
func TestValidateEmail(t *testing.T) {
    cases := []struct {
        name     string
        input    string
        hasError bool
    }{
        {
            name:     "WithEmptyString_ReturnsError",
            input:    "",
            hasError: true,
        },
        {
            name:     "WithNoAtSign_ReturnsError",
            input:    "invalid",
            hasError: true,
        },
        {
            name:     "WithValidFormat_ReturnsNil",
            input:    "user@example.com",
            hasError: false,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            err := ValidateEmail(tc.input)

            if tc.hasError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### TypeScript

```typescript
// ✅ REQUIRED — table-driven with it.each or similar
describe('validateEmail', () => {
    const cases = [
        {
            name: 'WithEmptyString_ReturnsError',
            input: '',
            isValid: false,
        },
        {
            name: 'WithNoAtSign_ReturnsError',
            input: 'invalid',
            isValid: false,
        },
        {
            name: 'WithValidFormat_ReturnsTrue',
            input: 'user@example.com',
            isValid: true,
        },
    ];

    it.each(cases)('$name', ({ input, isValid }) => {
        expect(validateEmail(input)).toBe(isValid);
    });
});
```

### Constraints

- Each table case must have a `name` field — no anonymous cases
- Table case structs must not exceed 6 fields; extract a builder or factory if more are needed
- The loop body must not contain conditional logic — if a case needs different assertions, it belongs in a separate test function

---

## Rule 4: Test Helper Placement & Naming

### Placement

| Helper Scope | Location |
|-------------|----------|
| Used by one test file | Bottom of the same `_test.go` / `.test.tsx` file |
| Used across files in one package | `testutil_test.go` (Go) / `__testutils__.ts` (TypeScript) in the same directory |
| Used across packages | `internal/testutil/` (Go) / `src/test-utils/` (TypeScript) |

### Naming

Test helpers must start with a descriptive verb — never `helper` or `util` as a prefix.

```go
// ❌ FORBIDDEN — vague helper names
func helperSetup(t *testing.T) { ... }
func newTestThing() *Thing { ... }

// ✅ REQUIRED — descriptive verb prefix
func createTestSession(t *testing.T) *Session { ... }
func seedDatabaseWithUsers(t *testing.T, db *gorm.DB, count int) { ... }
func assertErrorCode(t *testing.T, err *apperror.AppError, code string) { ... }
```

### t.Helper() Mandate (Go)

Every Go test helper function that calls `t.Fatal`, `t.Error`, or assertion methods **must** call `t.Helper()` as its first statement.

```go
// ❌ FORBIDDEN — missing t.Helper()
func assertSessionValid(t *testing.T, session *Session) {
    assert.NotNil(t, session)
    assert.NotEmpty(t, session.ID)
}

// ✅ REQUIRED — t.Helper() as first statement
func assertSessionValid(t *testing.T, session *Session) {
    t.Helper()

    assert.NotNil(t, session)
    assert.NotEmpty(t, session.ID)
}
```

---

## Rule 5: Test Body Structure — Arrange / Act / Assert

Every test body must follow the AAA pattern with **blank line separators** between sections. Comments marking each section are optional but encouraged for complex tests.

```go
func TestCreateSession_WithValidCredentials_ReturnsSessionId(t *testing.T) {
    // Arrange
    store := NewMockSessionStore()
    service := NewSessionService(store)
    credentials := createValidCredentials(t)

    // Act
    result := service.CreateSession(credentials)

    // Assert
    assert.NoError(t, result.Error)
    assert.NotEmpty(t, result.Value.ID)
}
```

### Constraints

- No logic in the Assert section — only assertions
- No assertions in the Arrange section
- One Act call per test — if testing multiple operations, write separate tests
- `t.Cleanup()` (Go) / `afterEach` (TypeScript) for teardown — never defer in tests when `t.Cleanup` is available

---

## Rule 6: Test Isolation

- Each test must be independently runnable — no dependency on execution order
- Tests must not share mutable state — create fresh instances in each test or table case
- Database tests must use transactions with rollback, or per-test temporary databases
- File system tests must use `t.TempDir()` (Go) — never write to the project tree

```go
// ❌ FORBIDDEN — shared mutable state
var globalStore = NewMockStore()

func TestA(t *testing.T) {
    globalStore.Add("key", "value") // Mutates shared state
}

func TestB(t *testing.T) {
    val := globalStore.Get("key") // Depends on TestA running first
}

// ✅ REQUIRED — isolated state per test
func TestA(t *testing.T) {
    store := NewMockStore()
    store.Add("key", "value")
    // ...
}

func TestB(t *testing.T) {
    store := NewMockStore()
    // ...
}
```

---

## Rule 7: Integration Test Boundaries

Integration tests must clearly declare their external dependencies and be skippable in fast-test mode.

### Go

```go
func TestDatabaseSync_WithLiveConnection_SyncsAllTables(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    db := setupTestDatabase(t)

    result := SyncAllTables(db)

    assert.NoError(t, result.Error)
    assert.Equal(t, 5, result.Value.TablesProcessed)
}
```

### TypeScript

```typescript
describe.skipIf(process.env.CI !== 'true')('DatabaseSync', () => {
    it('WithLiveConnection_SyncsAllTables', async () => {
        const db = await setupTestDatabase();

        const result = await syncAllTables(db);

        expect(result.tablesProcessed).toBe(5);
    });
});
```

---

## Complexity Targets

| Metric | Target |
|--------|--------|
| Test function length | ≤ 20 lines (excluding table cases) |
| Table case fields | ≤ 6 |
| Assertions per test | ≤ 5 (split if more needed) |
| Cyclomatic complexity per test | 0–1 (no branching in tests) |

---

## Checklist Summary (Copy for PRs)

```
[ ] Test file mirrors source file name with test suffix
[ ] Test functions use three-part naming: Unit_Scenario_Outcome
[ ] 3+ similar scenarios use table-driven tests with named cases
[ ] Test helpers use descriptive verb prefixes and call t.Helper() (Go)
[ ] Test body follows Arrange / Act / Assert with blank line separators
[ ] Each test is independently runnable with no shared mutable state
[ ] Integration tests are skippable in short/fast mode
[ ] No branching logic inside test bodies
```

---

## Cross-References

- [Code Style — Rule 6: 15-Line Limit](./04-code-style/00-overview.md) — Test functions have a relaxed 20-line limit
- [Cyclomatic Complexity](./06-cyclomatic-complexity.md) — Tests must have 0–1 complexity
- [Function Naming](./10-function-naming.md) — Test helpers follow the same no-boolean-flag rule
- [Boolean Principles](./02-boolean-principles/00-overview.md) — Table case fields use `is`/`has` prefixes for boolean columns
- [Master Coding Guidelines — §13](./15-master-coding-guidelines/00-overview.md) — Summary table in the master reference

---

*Test naming & structure specification v1.0.0 — 2026-03-11*
