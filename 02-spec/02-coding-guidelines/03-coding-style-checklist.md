# AI Coding Style Checklist (Root Rule) (AI Execution Prompt)

> **/goal** Enforce strict parameter limits (max 3 parameters or ≤100 chars), PascalCase acronyms, zero magic strings/numbers, isolated temporary scripts, and Unix LF UTF-8 encoding.
> **/learn** Master the refactoring of multi-argument methods into options structs/classes, constant extraction for magic values, and strict git cleanliness for scratchpads.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Refactor any method requiring 4+ parameters or exceeding 100 characters into an options struct or class.
- [ ] `/learn` Never use uppercase abbreviations or raw magic values; enforce PascalCase (`UserId`, `HttpServer`) and extracted constants.
- [ ] `/goal` Keep all temporary debugging scripts confined to `.ai-memory/temp-scripts/` and never commit them to git.
- [ ] `/learn` Ensure all files are strictly UTF-8 encoded with Unix LF (`\n`) line endings and a terminating newline.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Method Arguments & Signatures (The "3 Parameter" Rule)

Functions and methods should have a maximum of **3 parameters**.

If a function requires 4 or more parameters, or if the method signature exceeds **100 characters** in length:

- You **MUST** refactor the parameters into an options struct/class/object (e.g., `UpdateUserOptions`).
- If you absolutely cannot use an options object (e.g., interfacing with a legacy system), you **MUST** split the signature to have **one parameter per line**.

### Example (Go)

❌ FORBIDDEN:
```go
func ProcessTransaction(userId int, amount float64, currency string, idempotencyKey string, retryCount int) error { ... }
```

✅ REQUIRED (Options Struct):
```go
type TransactionOptions struct {
    UserId         int
    Amount         float64
    Currency       string
    IdempotencyKey string
    RetryCount     int
}

func ProcessTransaction(opts TransactionOptions) error { ... }
```

✅ REQUIRED (One Per Line - Only if Struct is impossible):
```go
func ProcessTransaction(
    userId int,
    amount float64,
    currency string,
    idempotencyKey string,
    retryCount int,
) error { ... }
```

## 2. Acronyms & Magic Strings

- Acronyms must be PascalCase (`UserId` not `UserID`, `HttpServer` not `HTTPServer`).
- Magic strings and numbers must be extracted to constants at the top of the file or in a dedicated constants package.

## 3. Temporary Scripts

- Any temporary code, scratchpads, or debugging scripts you create must be written to the `.ai-memory/temp-scripts/` directory.
- **NEVER** commit temporary scripts to Git.

## 4. File Encoding & Line Endings

- **Encoding:** All files MUST be encoded in **UTF-8 without BOM**.
- **Line Endings:** All files MUST use strictly **Unix-style Line Feeds (LF / \n)**. Carriage returns (\r\n) are strictly prohibited.
- **EOF Newline:** All source files, markdown, and config files MUST end with a single empty newline (\n). This should be handled by .editorconfig (insert_final_newline = true).

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ROOT-003: Core Coding Style and Parameter Rules

**Given** Source files in the repository.
**When** Codebases are audited by coding guideline scanners and lint rules.
**Then** Function signatures (≤3 params), acronym casing (PascalCase), magic constant isolation, temporary script isolation, and Unix LF line endings are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
