# Subtask 294.3: Refactor cmdfixauth Error Handling, Control Flow & Formatting

- **Parent Plan:** [.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md](../../pending/294-code-refactoring-and-coding-guideline-update.md)
- **Status:** Completed
- **Spec Reference:** [02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md](../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md)
- **File Boundaries:**
  - `cli/cmdfixauth/fixauth.go`

---

## 1. Objectives & Scope

Refactor `cli/cmdfixauth/fixauth.go` to strictly adhere to repository error sovereignty, clean control flow, and parameter formatting guidelines:

1. **Leaf Error Sovereignty (`checkGitRepo`):**
   - Refactor `checkGitRepo() error` to return `*apperror.AppError`.
   - Propagate strongly-typed error directly in `RunFixAuth` without loss of type fidelity.
2. **Statement Decoupling from `if`:**
   - In `executeFixAuthPipeline`, separate the initialization of `setupSSHDirAndKey` from the `if` guard condition.
   - Insert required vertical blank line before `if` and after `}`.
3. **Multiline Argument & Parameter Formatting:**
   - Reformat all function signatures taking more than two arguments across discrete lines with trailing commas.
   - Apply formatting to `executeFixAuthPipeline`, `setupSSHDirAndKey`, `fixAuthGenerate`, and `validateKeyGeneration`.
4. **Style & Rhythm Verification:**
   - Ensure every function maintains $\le 15$ lines.
   - Verify zero nested conditionals and proper blank line spacing before `return` statements.

---

## 2. Step-by-Step Execution Plan

### Step 1: Update `checkGitRepo` Signature & Call Site
- Update `func checkGitRepo() error` to `func checkGitRepo() *apperror.AppError`.
- Ensure `return apperror.NewSimple(...)` returns the typed error.
- In `RunFixAuth`, update caller handling:
  ```go
  repoErr := checkGitRepo()
  if repoErr != nil {
  	return repoErr
  }
  ```

### Step 2: Separate Statement from `if` in `executeFixAuthPipeline`
- Target:
  ```go
  // BEFORE:
  if err := setupSSHDirAndKey(keyPath, email, assumeYes, force); err != nil {
  	return err
  }

  // AFTER:
  err := setupSSHDirAndKey(keyPath, email, assumeYes, force)
  if err != nil {
  	return err
  }
  ```

### Step 3: Reformat Function Signatures ($>2$ Parameters)
Format the following four functions with multiline parameter definitions:

1. `executeFixAuthPipeline`:
   ```go
   func executeFixAuthPipeline(
   	user string,
   	email string,
   	assumeYes bool,
   	force bool,
   ) error
   ```
2. `setupSSHDirAndKey`:
   ```go
   func setupSSHDirAndKey(
   	keyPath string,
   	email string,
   	assumeYes bool,
   	force bool,
   ) error
   ```
3. `fixAuthGenerate`:
   ```go
   func fixAuthGenerate(
   	keyPath string,
   	email string,
   	assumeYes bool,
   	force bool,
   ) error
   ```
4. `validateKeyGeneration`:
   ```go
   func validateKeyGeneration(
   	hasKey bool,
   	keyPath string,
   	assumeYes bool,
   	force bool,
   ) bool
   ```

### Step 4: Vertical Blank Line & Guard Clause Audit
- Audit all conditionals in `cli/cmdfixauth/fixauth.go` to ensure:
  - Blank line before every `if`.
  - Blank line after every closing brace `}`.
  - Blank line before every `return` statement.
  - All booleans strictly follow `is` or `has` prefixes (e.g., `isGitRepo`, `hasKey`, `isConfirmed`).

---

## 3. Verification & Acceptance Criteria

- [x] `checkGitRepo()` returns `*apperror.AppError`.
- [x] No inline assignment statements inside `if` conditionals in `cli/cmdfixauth/fixauth.go`.
- [x] All 4 functions with $>2$ arguments formatted multiline.
- [x] Zero compiler build or test commands executed during planning or authoring.
- [x] Relative Git paths strictly used throughout all documentation.
