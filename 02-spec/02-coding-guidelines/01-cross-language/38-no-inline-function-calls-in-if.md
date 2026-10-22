# Zero Inline Function Calls in `if` Statements — Coding Guidelines

## 1. Overview & Core Mission

A pervasive anti-pattern in modern programming—particularly in languages supporting inline statement initialization like Go (`if x := fn(); x != nil`), Python (`if (x := fn()) is not None`), and Rust—is cramming function execution, multi-parameter argument passing, and conditional branching into a single compound line.

This practice degrades readability, obstructs debugger step-over breakpoints, conceals complex error states, and clutters guard clauses.

This specification establishes the **Zero Inline Function Calls in `if` Standard** across all supported languages. Every function call MUST be executed on a dedicated statement line preceding the conditional check. The `if` expression MUST evaluate simple variables or direct boolean predicates only.

---

## 2. Rationale & Cognitive Load Reduction

### 2.1 The Problem with Inline Execution
1. **Hidden Side Effects:** Readers scanning code often mistake `if` statements for pure predicate evaluations. When side effects (such as disk I/O, network requests, DB transactions) occur inside an `if` header, their side effects are easily overlooked.
2. **Debugger Hostility:** Debuggers cannot stop between the completion of `fn()` and the evaluation of the boolean condition when they share a single line. Stepping over steps past the entire branch decision.
3. **Line Length Bloat:** A function with 3–4 arguments declared inline within an `if` statement quickly explodes past 100 characters or wraps illegibly.
4. **Violation of Function Sizing Limits:** Packing invocation and condition checking together leads developers to overlook bloated function execution boundaries.

---

## 3. Anti-Patterns & Canonical Modernization

### 3.1 Go: Inline Short-Variable Declaration

#### ❌ Anti-Pattern: Inlining Function Call Inside `if`
```go
// ❌ FORBIDDEN: Function call embedded inside if conditional
if err := setupSSHDirAndKey(keyPath, email, assumeYes, force); err != nil {
	return err
}

// ❌ FORBIDDEN: Value assignment and error tuple check crammed together
if sub, err := store.ClaimSubtask(dbPath, taskID, agentRole); err != nil || sub == nil {
	return handleClaimFailure(err)
}
```

#### ✅ MANDATORY: Statement Separation & Vertical Blank Line Rhythm
```go
// ✅ MANDATORY: Separate invocation line, vertical rhythm, clean guard clause
err := setupSSHDirAndKey(
	keyPath,
	email,
	assumeYes,
	force,
)
if err != nil {
	return err
}

sub, err := store.ClaimSubtask(dbPath, taskID, agentRole)
if err != nil {
	return handleClaimFailure(err)
}

if sub == nil {
	return handleEmptyQueue()
}
```

### 3.2 Python: Walrus Operator `:=` Abuse Inside `if`

#### ❌ Anti-Pattern: Walrus Assignment in `if`
```python
# ❌ FORBIDDEN: Walrus operator embedding side effects into branch expression
if (result := repository.find_active_lease(repo_id, timeout=30)) is not None:
    process_lease(result)
```

#### ✅ MANDATORY: Separate Statement Assignment
```python
# ✅ MANDATORY: Explicit assignment statement followed by clean predicate
result = repository.find_active_lease(repo_id, timeout=30)
if result is not None:
    process_lease(result)
```

### 3.3 Rust: `if let` with Inline Complex Function Calls

#### ❌ Anti-Pattern: Complex Inline Execution
```rust
// ❌ FORBIDDEN: Inline execution in if-let binding
if let Ok(config) = load_remote_configuration(&endpoint, auth_token, max_retries) {
    apply_config(config);
}
```

#### ✅ MANDATORY: Separate Execution
```rust
// ✅ MANDATORY: Execute function first, match on result cleanly
let config_result = load_remote_configuration(
    &endpoint,
    auth_token,
    max_retries,
);

if let Ok(config) = config_result {
    apply_config(config);
}
```

---

## 4. Coding Standard Enforcement Rules

1. **Separation Rule:** Any expression that invokes a function or method must be assigned to an explicit variable on its own line before any `if` statement.
2. **Spacing Rule:** There must be a blank line before the subsequent `if` condition when preceding assignments or setup logic exist.
3. **Predicate Simplicity:** The conditional expression of the `if` statement must only test:
   - Error comparisons (`err != nil`).
   - Boolean variables (`isReady`, `hasUpdates`).
   - Simple comparison operators between already-evaluated variables (`count > 0`, `status == StatusActive`).
4. **No Compound Evaluation with Invocations:** Compound conditions like `if isValid && checkRemote(token)` are forbidden; invoke `isRemoteValid := checkRemote(token)` first, then evaluate `if isValid && isRemoteValid`.

---

## 5. Summary Checklist

- [ ] Zero `if err := ...; err != nil` patterns in codebase.
- [ ] Zero walrus operators `:=` embedded directly inside conditional `if` headers in Python.
- [ ] Every function invocation sits on its own dedicated statement line.
- [ ] Clean vertical spacing: blank line before `if`, blank line after `}`, blank line before `return`.
