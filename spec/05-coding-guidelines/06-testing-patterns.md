# Testing Patterns — Unit Test Conventions

## Overview

Standards for writing, naming, and organizing tests across TypeScript
and Go projects. Tests are first-class code — the same quality rules
apply.

---

## 1. Test File Placement

Tests live next to the code they verify.

| Language | Convention | Example |
|----------|-----------|---------|
| TypeScript | `ComponentName.test.tsx` | `TypeBadge.test.tsx` |
| TypeScript | `utilName.test.ts` | `specData.test.ts` |
| Go | `filename_test.go` | `scanner_test.go` |

Never place tests in a separate top-level `tests/` directory unless
they are integration or end-to-end tests.

---

## 2. Test Naming — Describe What, Not How

### TypeScript (Vitest)

Use `describe` for the unit and `it` for the behavior.

```ts
describe("ProjectTypes", () => {
  it("returns correct label for Go projects", () => {
    const config = ProjectTypes[ProjectType.Go];

    expect(config.label).toBe("Go");
  });

  it("includes an icon for every project type", () => {
    const types = Object.values(ProjectTypes);
    const hasIcons = types.every((type) => type.icon !== undefined);

    expect(hasIcons).toBe(true);
  });
});
```

### Go

Use `Test<Function>_<scenario>` naming.

```go
func TestParseSemver_ValidVersion(t *testing.T) {
    version, err := Parse("1.2.3")

    require.NoError(t, err)
    assert.Equal(t, 1, version.Major)
}

func TestParseSemver_PreReleaseTag(t *testing.T) {
    version, err := Parse("1.0.0-rc.1")

    require.NoError(t, err)
    assert.Equal(t, "rc.1", version.PreRelease)
}
```

---

## 3. Test Structure — Arrange, Act, Assert

Every test follows three phases separated by blank lines.

```ts
it("filters active repos", () => {
  // Arrange
  const repos = [
    { name: "api", isActive: true },
    { name: "legacy", isActive: false },
  ];

  // Act
  const result = filterActiveRepos(repos);

  // Assert
  expect(result).toHaveLength(1);
  expect(result[0].name).toBe("api");
});
```

```go
func TestFilterActiveRepos(t *testing.T) {
    repos := []Repo{
        {Name: "api", IsActive: true},
        {Name: "legacy", IsActive: false},
    }

    result := filterActiveRepos(repos)

    assert.Len(t, result, 1)
    assert.Equal(t, "api", result[0].Name)
}
```

---

## 4. Table-Driven Tests (Go)

Use table-driven tests for functions with multiple input/output
combinations.

```go
func TestIsEligible(t *testing.T) {
    tests := []struct {
        name     string
        status   string
        role     string
        expected bool
    }{
        {"active admin", "active", "admin", true},
        {"inactive user", "inactive", "user", false},
        {"active user", "active", "user", true},
    }

    for _, testCase := range tests {
        t.Run(testCase.name, func(t *testing.T) {
            result := isEligible(testCase.status, testCase.role)

            assert.Equal(t, testCase.expected, result)
        })
    }
}
```

---

## 5. Parameterized Tests (TypeScript)

Use `it.each` or loops for similar test cases.

```ts
describe("statusColor", () => {
  const cases = [
    { status: RepoStatus.Dirty, expected: "text-yellow-400" },
    { status: RepoStatus.Clean, expected: "text-primary" },
  ];

  cases.forEach(({ status, expected }) => {
    it(`returns ${expected} for ${status}`, () => {
      const result = statusColor(status);

      expect(result).toBe(expected);
    });
  });
});
```

---

## 6. What to Test

| Category | Test | Priority |
|----------|------|----------|
| Pure functions | Input → output | High |
| Constants/enums | Completeness checks | Medium |
| Edge cases | Empty input, nulls, boundaries | High |
| Error paths | Invalid input, missing data | High |
| Component render | Key elements visible | Medium |
| User interaction | Click, type, submit | Medium |
| Integration | Multi-module workflows | High |

### What NOT to Unit Test

- Implementation details (private methods, internal state).
- Third-party library internals.
- Trivial getters/setters with no logic.

---

## 7. Test Isolation

- Each test must be independent — no shared mutable state.
- Use `beforeEach` (TypeScript) or `t.TempDir()` (Go) for setup.
- Never rely on test execution order.
- Clean up side effects (files, database rows) after each test.

---

## 8. Assertion Style

### TypeScript

Prefer specific matchers over generic `toBe(true)`.

```ts
// ✅ Correct
expect(result).toHaveLength(3);
expect(element).toBeInTheDocument();
expect(value).toContain("expected");

// ❌ Wrong
expect(result.length === 3).toBe(true);
expect(!!element).toBe(true);
```

### Go

Use `assert` for soft checks, `require` for hard stops.

```go
require.NoError(t, err)          // stop if error
assert.Equal(t, expected, actual) // continue on failure
assert.Contains(t, output, "success")
```

---

## 9. Coverage Expectations

| Scope | Target |
|-------|--------|
| Pure utility functions | 90%+ |
| Core business logic | 80%+ |
| UI components | Key renders and interactions |
| CLI commands | Happy path + error paths |
| Integration tests | Critical workflows |

Coverage is a guide, not a goal. A well-tested 70% codebase is
better than a poorly-tested 100% codebase.

---

## 10. Test Quality Rules

All coding guidelines apply to test code:

- Meaningful variable names (no `x`, `s`, `d`).
- Functions under 25 lines — extract test helpers.
- No magic strings — use constants for expected values.
- Boolean naming — `isValid`, `hasError`.
- Blank line before `return` and before assertions.

---

## References

- Code Quality Improvement: `spec/05-coding-guidelines/01-code-quality-improvement.md`
- Go Code Style: `spec/05-coding-guidelines/02-go-code-style.md`
- Naming Conventions: `spec/05-coding-guidelines/03-naming-conventions.md`
- Golden Fixture Regeneration: `spec/05-coding-guidelines/21-golden-fixture-regeneration.md`
