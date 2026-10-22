# Multiline Parameter & Argument Formatting — Coding Guidelines

## 1. Overview & Core Mission

Horizontal code scrolling and unwieldy, packed parameter lists represent a significant obstacle to code review velocity, git diff readability, and bug detection. When functions declaring or passing $>2$ arguments are formatted on a single line, git diffs flag the entire line on any single argument change, concealing the actual modification.

This specification establishes the **Multiline Parameter & Argument Formatting Standard** across all supported languages (Go, TypeScript, Python, Rust, PHP).

---

## 2. Threshold Rules & Alignment Rules

### 2.1 The Cardinal Threshold Rules
1. **$\le 2$ Arguments / Parameters:** May remain inline on a single line, provided the entire line does not exceed the 100-character column cap.
2. **$\ge 3$ Arguments / Parameters ($> 2$ arguments):** **MUST** be formatted vertically across multiple lines. Every parameter/argument gets its own dedicated line.

### 2.2 Formatting & Alignment Invariants
1. **One Argument Per Line:** Never place multiple arguments on the same continuation line.
2. **Dedicated Opening & Closing Lines:**
   - The opening parenthesis remains on the function declaration/call line.
   - The first argument begins on the immediately following line, indented by 1 tab / standard indentation.
   - The closing parenthesis sits on its own dedicated line at the same indentation level as the calling statement or declaration.
3. **Trailing Commas:**
   - In Go: Trailing comma is mandatory on the final parameter/argument line.
   - In TypeScript / JavaScript: Trailing comma is mandatory (enables clean single-line git diffs).
   - In Python: Trailing comma is mandatory.
   - In Rust: Trailing comma is mandatory.
   - In PHP: Trailing comma is mandatory.

---

## 3. Polyglot Examples

### 3.1 Go Declarations & Invocations

#### ❌ Anti-Pattern: Crammed Signatures and Call Sites
```go
// ❌ FORBIDDEN: 4 parameters crammed on a single line
func setupSSHDirAndKey(keyPath, email string, assumeYes, force bool) error { ... }

// ❌ FORBIDDEN: 4 arguments crammed on a single line invocation
err := setupSSHDirAndKey(keyPath, email, assumeYes, force)
```

#### ✅ MANDATORY: Clean Multiline Formatting
```go
// ✅ MANDATORY: Multiline declaration with dedicated lines and trailing comma
func setupSSHDirAndKey(
	keyPath string,
	email string,
	assumeYes bool,
	force bool,
) error {
	// body
}

// ✅ MANDATORY: Multiline invocation with trailing comma
err := setupSSHDirAndKey(
	keyPath,
	email,
	assumeYes,
	force,
)
if err != nil {
	return err
}
```

### 3.2 TypeScript / React

#### ❌ Anti-Pattern: Crammed Arguments
```typescript
// ❌ FORBIDDEN: >2 parameters crammed inline
function renderSubtaskCard(subtask: Subtask, agentRole: string, isSelected: boolean, onSelect: () => void) { ... }

// ❌ FORBIDDEN: Inline call site
renderSubtaskCard(item, "worker-01", true, handleSelect);
```

#### ✅ MANDATORY: Multiline Formatting
```typescript
// ✅ MANDATORY: Multiline declaration
function renderSubtaskCard(
  subtask: Subtask,
  agentRole: string,
  isSelected: boolean,
  onSelect: () => void,
): React.ReactElement {
  // body
}

// ✅ MANDATORY: Multiline invocation
renderSubtaskCard(
  item,
  "worker-01",
  true,
  handleSelect,
);
```

### 3.3 Python

#### ❌ Anti-Pattern: Crammed Arguments
```python
# ❌ FORBIDDEN: >2 arguments in declaration or call
def create_subtask_entry(db_path, parent_id, task_code, title, owned_files):
    pass

result = create_subtask_entry(db, pid, code, title, files)
```

#### ✅ MANDATORY: Multiline Formatting
```python
# ✅ MANDATORY: Multiline declaration with trailing comma
def create_subtask_entry(
    db_path: str,
    parent_id: str,
    task_code: str,
    title: str,
    owned_files: list[str],
) -> SubtaskResult:
    pass

# ✅ MANDATORY: Multiline call site
result = create_subtask_entry(
    db,
    pid,
    code,
    title,
    files,
)
```

---

## 4. Git Diff Benefits

Compare the git diffs when adding an argument or modifying a flag:

### ❌ Single-Line Diff (Opaque & High Conflict Risk)
```diff
- err := executePipeline(ctx, repoPath, branchName, true)
+ err := executePipeline(ctx, repoPath, branchName, true, timeout)
```
*The entire line is flagged as replaced. In parallel worker swarms, two agents modifying different arguments will generate a git merge conflict.*

### ✅ Multiline Diff (Isolated, Atomic & Conflict-Free)
```diff
  err := executePipeline(
      ctx,
      repoPath,
      branchName,
      true,
+     timeout,
  )
```
*Only the new line is added. Zero conflicts with other edits.*

---

## 5. Summary Checklist

- [ ] All function declarations with $>2$ parameters formatted one parameter per line.
- [ ] All function call sites with $>2$ arguments formatted one argument per line.
- [ ] Closing parenthesis placed on its own line aligned with declaration/statement indentation.
- [ ] Trailing commas present on the final argument in languages supporting it.
