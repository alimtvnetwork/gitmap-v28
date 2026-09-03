# AI Coding Style Checklist (Root Rule)

> **CRITICAL AI INSTRUCTION:** Apply these stylistic rules to every single file you generate or modify.

## 1. Method Arguments & Signatures (The "3 Parameter" Rule)

Functions and methods should have a maximum of **3 parameters**.

If a function requires 4 or more parameters, or if the method signature exceeds **100 characters** in length:

- You **MUST** refactor the parameters into an options struct/class/object (e.g., \UpdateUserOptions\).
- If you absolutely cannot use an options object (e.g., interfacing with a legacy system), you **MUST** split the signature to have **one parameter per line**.

### Example (Go)

❌ BAD:
\\\go
func ProcessTransaction(userId int, amount float64, currency string, idempotencyKey string, retryCount int) error { ... }
\\\

✅ GOOD (Options Struct):
\\\go
type TransactionOptions struct {
    UserID         int
    Amount         float64
    Currency       string
    IdempotencyKey string
    RetryCount     int
}

func ProcessTransaction(opts TransactionOptions) error { ... }
\\\

✅ GOOD (One Per Line - Only if Struct is impossible):
\\\go
func ProcessTransaction(
    userId int,
    amount float64,
    currency string,
    idempotencyKey string,
    retryCount int,
) error { ... }
\\\

## 2. Acronyms & Magic Strings

- Acronyms must be PascalCase (\UserId\ not \UserID\, \HttpServer\ not \HTTPServer\).
- Magic strings and numbers must be extracted to constants at the top of the file or in a dedicated constants package.

## 3. Temporary Scripts

- Any temporary code, scratchpads, or debugging scripts you create must be written to the \.lovable/temp-scripts/\ directory.
- **NEVER** commit temporary scripts to Git.

## 4. File Encoding & Line Endings

- **Encoding:** All files MUST be encoded in **UTF-8 without BOM**.
- **Line Endings:** All files MUST use strictly **Unix-style Line Feeds (LF / \n)**. Carriage returns (\r\n) are strictly prohibited.
- **EOF Newline:** All source files, markdown, and config files MUST end with a single empty newline (\n). This should be handled by .editorconfig (insert_final_newline = true).
