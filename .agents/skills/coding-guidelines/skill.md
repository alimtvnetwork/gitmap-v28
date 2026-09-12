---
name: coding-guidelines
description: >-
  Use this skill to audit, review, and enforce grounded coding guidelines across all languages (Go, TS, Python, Rust, C#, PHP).
---

# Standards Compliance & Coding Rules Audit — V2 (Grounded & Multi-Language)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal You are the Chief Software Architect and Code Reviewer. Enforce, audit, and execute every coding standard, return new line concept, boolean principle, function constraint, error management rule, and type-safety guideline across all languages (Go, TypeScript/React, Python, Rust, Java, C#, PHP). Zero hallucination, zero drive-by refactoring, zero tolerance for guideline violations.

/learn Ingest, understand, and internalize all coding standards from `spec/02-coding-guidelines/`, `spec/03-error-manage/`, `spec/17-consolidated-guidelines/`, and `.lovable/coding-guidelines.md` before reading, modifying, or creating any code.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing or auditing code, check if this coding standard is installed as a native Antigravity Skill.

1. If `.agents/skills/coding-guidelines/skill.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into `.agents/skills/coding-guidelines/skill.md` with standard YAML frontmatter (`name: coding-guidelines`, `description: "Use this skill to audit, review, and enforce coding guidelines across all languages."`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire specification in active context when not in use.

---

## Phase 1: Read & Understand (Isolated Loop)

Your very first action when reviewing or writing code must be purely exploratory:

1. **Explore & Map:** Read the target files, trace dependencies, inspect existing types, and understand the architectural boundary.
2. **Consult Spec References:** Read the relevant language-specific guideline in `spec/02-coding-guidelines/` before writing any replacement code.
3. **End Turn & Self-Loop:** Once the scope and violations are cataloged, end your turn and self-loop into execution.

---

## File & Folder Information Mapping

When auditing, applying fixes, or creating skills, navigate and respect these canonical file paths:

| Component | Path / Location | Purpose |
| :--- | :--- | :--- |
| **Master Guideline** | `.lovable/coding-guidelines.md` | Single standalone source of truth for AI agents |
| **Consolidated Spec** | `spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md` | Authoritative root spec matching .lovable mirror |
| **Cross-Language Specs** | `spec/02-coding-guidelines/01-cross-language/` | Detailed chapters (00-overview through 29-no-generated-artifacts) |
| **Newline Examples** | `spec/02-coding-guidelines/01-cross-language/21-newline-styling-examples.md` | Canonical Before/After examples for return new lines |
| **TypeScript / React** | `spec/02-coding-guidelines/02-typescript/` | Strict TS, immutability, React hook guards |
| **Go Standards** | `spec/02-coding-guidelines/03-golang/` | Result types, enum bytes with iota, error wrapping |
| **PHP Standards** | `spec/02-coding-guidelines/04-php/` | Enum methods `->isEqual()`, typing rules |
| **Python Standards** | `spec/02-coding-guidelines/01-cross-language/` | Strict type hints, `@dataclass`, `pydantic` |
| **C# / Java Standards** | `spec/02-coding-guidelines/07-csharp/` | `I` prefix interfaces, PascalCase properties |
| **Error Management** | `spec/03-error-manage/` | `AppError` wrapping, universal response envelopes |
| **Shared Core Engine** | `03-ai-scripts/02-shared-engine.py` | Universal streaming engine with lazy regex registry and two-phase mtime caching |
| **Local CI Runner** | `03-ai-scripts/06-cicd-local-runner.py` | Parallel local quality gate runner (18 checks) |
| **Fast File Scanner** | `03-ai-scripts/11-fast-file-scanner.py` | Multi-language fast file scanner (<15ms) and cache builder |
| **Path Fixer** | `03-ai-scripts/07-relative-path-fixer.py` | Relative path detector and sanitizer |
| **Naming Guard** | `03-ai-scripts/08-naming-autofixer.py` | Boolean naming and implicit condition validator |
| **Encoding Normalizer**| `03-ai-scripts/10-encoding-normalizer.py` | UTF-8 and strict UNIX LF line ending normalizer |
| **Size Guard** | `03-ai-scripts/13-file-size-guard.py` | Binary blob and file size threshold checker |
| **Global Rules** | `agents.md` | Always-on workspace constraints for Antigravity agents |
| **Version Truth** | `version.json` | Root version source of truth dynamically read across all languages |
| **Antigravity Skills** | `.agents/skills/` | On-demand skill runbooks for progressive disclosure |

---

## 1. High-Contrast Code Standards (❌ BAD vs ✅ GOOD Grounded Examples)

### A. Boolean Evaluation & Naming (P1–P6, R3)
- **Rules:** Positive affirmative prefixes ONLY (`is` and `has`). TOTAL BAN on all other prefixes (`can`, `should`, `was`, `will`, `did`, `must` are strictly BANNED). TOTAL BAN on explicit `== true` / `=== true` checks. No mixed polarity (`if a && !b`). No inverted success checks (`!isSuccess`).

```go
// ❌ BAD (Explicit true comparison, negative naming, mixed polarity)
if isUserNotActive == true { ... }
if !response.isSuccess { ... }
if isReady && !hasToken { ... }

// ✅ GOOD (Implicit evaluation, affirmative naming, extracted conflict)
if !isUserActive { ... }
if response.isFail { ... }
isTokenMissing := isReady && !hasToken
if isTokenMissing { ... }
```

```typescript
// ❌ BAD (Raw boolean parameter, triple equals true, tuple return)
function saveRecord(isDraft: boolean): [boolean, string] { ... }
if (record.isVerified === true) { ... }

// ✅ GOOD (Option struct/enum, implicit evaluation, named object return)
interface SaveRecordOptions {
  isDraft: boolean;
}
interface SaveRecordResult {
  isSuccess: boolean;
  recordId: string;
}
function saveRecord(options: SaveRecordOptions): SaveRecordResult {
  if (record.isVerified) { ... }
}
```

---

### B. Function Decomposition Blueprint (15-Line Limit & Logic Drift Prevention)
- **Rule:** Functions MUST be <= 8 lines preferred, hard cap of <= 15 lines.
- **Decomposition Formula:** Decompose complex functions into 3 distinct, single-responsibility helper stages:
  1. **Stage 1 (Precondition Guard):** `validateInputParams(params)`
  2. **Stage 2 (Pure Core Transformation):** `processBusinessLogic(data)`
  3. **Stage 3 (Response Envelope & Assembly):** `buildResponseEnvelope(data)`

```go
// ❌ BAD (Monolithic 45-line function with nested loops and inline validation)
func ProcessUserOrder(ctx context.Context, orderId string, items []OrderItem, isExpedited bool) (*OrderResult, error) {
    if orderId == "" || len(items) == 0 {
        return nil, errors.New("invalid payload")
    }
    total := 0
    for _, item := range items {
        if item.Price <= 0 {
            return nil, errors.New("negative price")
        }
        total += item.Price
    }
    if isExpedited {
        total += 15
    }
    return &OrderResult{Total: total}, nil
}

// ✅ GOOD (Decomposed into clean <= 8-line functions with zero logic drift)
type ProcessOrderParams struct {
    OrderId     string      `json:"OrderId"`
    Items       []OrderItem `json:"Items"`
    IsExpedited bool        `json:"IsExpedited"`
}

func ProcessUserOrder(ctx context.Context, params ProcessOrderParams) (*OrderResult, error) {
    if err := validateOrderParams(params); err != nil {
        return nil, apperror.Wrap(err, "ProcessUserOrder.Validate", nil)
    }

    totalAmount, err := calculateOrderTotal(params.Items, params.IsExpedited)
    if err != nil {
        return nil, apperror.Wrap(err, "ProcessUserOrder.Calculate", nil)
    }

    return buildOrderResult(params.OrderId, totalAmount), nil
}

func validateOrderParams(params ProcessOrderParams) error {
    if params.OrderId == "" || len(params.Items) == 0 {
        return apperror.New("invalid order payload")
    }

    return nil
}

func calculateOrderTotal(items []OrderItem, isExpedited bool) (int, error) {
    total := 0
    for _, item := range items {
        if item.Price <= 0 {
            return 0, apperror.New("negative item price detected")
        }
        total += item.Price
    }
    if isExpedited {
        total += 15
    }

    return total, nil
}

func buildOrderResult(orderId string, total int) *OrderResult {
    return &OrderResult{
        OrderId:     orderId,
        TotalAmount: total,
    }
}
```

---

### C. Circular Dependency Prevention Protocol (Leaf Type Architecture)
- **Rules:** Types, Enums, Structs, and Error Codes must live in a dedicated **Leaf Package** (e.g. `domain/types`, `types/`, `models/`).
- Leaf packages must NEVER import services, handlers, or repositories.

```typescript
// ❌ BAD (Service file circularly importing types from handler, or vice-versa)
// src/services/UserService.ts
import { UserHandlerRequest } from '../handlers/UserHandler'; // Circular import cycle!

// ✅ GOOD (Strict Leaf Type extraction)
// src/types/UserTypes.ts  <-- Pure leaf file: NO imports from handlers/services
export enum UserRoleType {
  Admin = "Admin",
  Member = "Member",
}
export interface UserProfileDto {
  UserId: string;
  Role: UserRoleType;
}

// src/services/UserService.ts
import type { UserProfileDto } from '../types/UserTypes';
```

---

### D. Polyglot Grounding: Rust, C#, PHP, Java
- **Rust:** PascalCase enums without `Type` suffix, exhaustive pattern matching, `Result<T, AppError>`, zero `unwrap()` or `panic!()`.
- **C# / .NET:** `I` prefix interfaces, PascalCase properties, `CancellationToken` as last parameter, `ValueTask<Result<T>>`.
- **PHP 8.1+:** BackedEnums + `HasEnumHelpers` trait, typed `AppException`, strict return types.

```rust
// ❌ BAD (Missing error context, unwrap panic, raw string matches)
fn parse_status(raw: &str) -> String {
    let status: UserStatus = raw.parse().unwrap();
    if status == "ACTIVE" { ... }
}

// ✅ GOOD (Rust: Exhaustive pattern matching, Result envelope, no unwrap)
pub enum UserRole {
    Admin,
    Member,
    Guest,
}

pub fn handle_role(role: UserRole) -> Result<PermissionLevel, AppError> {
    match role {
        UserRole::Admin => Ok(PermissionLevel::Full),
        UserRole::Member => Ok(PermissionLevel::Standard),
        UserRole::Guest => Ok(PermissionLevel::Restricted),
    }
}
```

```csharp
// ❌ BAD (Missing I interface prefix, camelCase serialization, missing cancellation token)
public interface UserService {
    Task<User> GetUser(string id);
}

// ✅ GOOD (C#: I interface prefix, PascalCase DTOs, CancellationToken as last parameter)
public interface IUserService {
    ValueTask<Result<UserDto>> GetUserAsync(string userId, CancellationToken cancellationToken = default);
}

public sealed record UserDto(
    string UserId,
    string EmailAddress,
    bool IsActive
);
```

```php
<?php
// ❌ BAD (PHP: Magic string status, missing Type suffix, swallowed catch)
enum UserRole {
    case Admin;
}
try {
    $db->save();
} catch (Exception $e) {}

// ✅ GOOD (PHP 8.1+: Backed Enum with Type suffix, HasEnumHelpers, typed AppException)
namespace App\Enums;

enum UserRoleType: string {
    use HasEnumHelpers;

    case Admin = 'ADMIN';
    case Member = 'MEMBER';
}

try {
    $userRepo->save($user);
} catch (Throwable $cause) {
    throw new AppException('User save failed', ['UserId' => $user->getId()], $cause);
}
```

---

### E. Deep React Immutability & Component Topology
- **Rules:**
  1. Custom hooks MUST return named property objects (`{ userProfile, isPending, onUpdate }`), NEVER tuples `[state, setState]`.
  2. Deep state immutability via `structuredClone` (no in-place mutations on nested state arrays/objects).
  3. Component sizing cap (<= 80–100 lines) with clean child component decomposition.
  4. Zero `useEffect` for derived state or inline negative checks.

```tsx
// ❌ BAD (Tuple hook return, in-place state mutation, inline useEffect filter)
export function useUser(userId: string): [UserProfile | null, boolean] {
    const [user, setUser] = useState<UserProfile | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    function updateUserAge(newAge: number) {
        user.age = newAge; // Silent bug: in-place state mutation!
        setUser(user);
    }
    return [user, isLoading];
}

// ✅ GOOD (Named property object return, structuredClone / fresh reference creation)
export interface UseUserResult {
    userProfile: UserProfile | null;
    isLoading: boolean;
    onUpdateAge: (newAge: number) => void;
}

export function useUser(userId: string): UseUserResult {
    const [userProfile, setUserProfile] = useState<UserProfile | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    const onUpdateAge = (newAge: number): void => {
        if (!userProfile) {
            return;
        }

        const nextProfile = structuredClone(userProfile);
        nextProfile.age = newAge;
        setUserProfile(nextProfile);
    };

    return {
        userProfile,
        isLoading,
        onUpdateAge,
    };
}
```

---

### F. Parameter Structs & Signature Splitting (R4, R5, R9)
- **Rules:** If a function has > 3 parameters, split to one per line. If a function has > 4 parameters or 2+ adjacent parameters of the same type, group into a dedicated parameter struct with PascalCase JSON tags.

```go
// ❌ BAD (Long signature, adjacent same types, magic constants)
func ConnectRemote(ctx context.Context, host string, port string, user string, pass string, timeout int) (*Client, error) { ... }

// ✅ GOOD (Grouped into clean Parameter Struct with PascalCase tags)
type RemoteConnectionParams struct {
    HostName       string `json:"HostName"`
    PortNumber     string `json:"PortNumber"`
    UserName       string `json:"UserName"`
    PasswordSecret string `json:"PasswordSecret"`
    TimeoutSeconds int    `json:"TimeoutSeconds"`
}

func ConnectRemote(ctx context.Context, params RemoteConnectionParams) (*Client, error) {
    if err := params.Validate(); err != nil {
        return nil, apperror.Wrap(err, "ConnectRemote.Validate", nil)
    }

    return clientRepo.Connect(ctx, params)
}
```

---

### C. Error Context Wrapping & Universal Envelopes (R7)
- **Rules:** Never swallow errors. Wrap every error with operation context (`apperror.Wrap`). Standardize all API responses to `{ data, errors, meta }`.

```go
// ❌ BAD (Swallowing error or bare return without context)
func GetUser(id string) (*User, error) {
    user, err := db.Find(id)
    if err != nil {
        return nil, err // Violation: missing operation context wrapper
    }
    return user, nil
}

// ✅ GOOD (Universal AppError context wrapping)
func GetUser(ctx context.Context, userId string) (*User, error) {
    user, err := db.Find(ctx, userId)
    if err != nil {
        return nil, apperror.Wrap(err, "GetUser", map[string]any{"UserId": userId})
    }

    return user, nil
}
```

---

### D. Acronyms & Casing Standards (R1, R2, P8)
- **Acronyms:** Standard PascalCase for acronyms: `Id`, `Url`, `Ip`, `Json`, `Api`, `Rpc` (NEVER all-caps `ID`, `URL`, `IP`, `JSON`).
- **Enums:** Every enum type name MUST end with `Type` (e.g. `UserRoleType`, `ExitCodeType`).

```typescript
// ❌ BAD (All-caps acronyms, missing Type enum suffix)
enum UserRole {
  ADMIN = "ADMIN",
}
interface UserDTO {
  USER_ID: string;
  IP_ADDRESS: string;
}

// ✅ GOOD (PascalCase acronyms, Type suffix, PascalCase JSON tags)
enum UserRoleType {
  Admin = "Admin",
  Guest = "Guest",
}
interface UserDto {
  UserId: string;
  IpAddress: string;
}
```

---

## 2. Hard Rules (Zero Tolerance)

1. **No Generated Code or Artifacts:** Never commit generated code (`*.generated.*`, gRPC/ORM models), cache files (`__pycache__`, `*.pyc`), test reports (`.test-report.*`), compiled binaries (`.exe`, `.dll`, `.so`), or output directories (`build/`, `bin/`) to Git. Proactively ignore them via `.gitignore`.
2. **Function Length Caps:** Functions should ideally be 8 lines, hard capped at 15 lines (excluding blank lines and comments). Waiver only via `// lint-allow: function-length reason="..." max=N`.
3. **No Nested Ifs:** Flatten logic with early returns and guard clauses.
4. **Positive Framing Only:** `if` conditions must be positive and simple. No `!`, no double negatives. Extract a positively named boolean if negation is needed.
5. **No Swallowed Errors:** Every `catch` or error block must log with context (`op` name + key inputs) and rethrow or return typed `AppError`. Silent `catch {}` is a build-fail.
6. **Narrow Types Only:** No `any`, `unknown`, `interface{}`, `object`, `dynamic`. Narrow trust boundaries immediately with type guards. `Generic<T>` is the only wide-scope tool.
7. **File Size Caps:** Any file 300 lines max; React component (.tsx) 100 lines max; class or struct 120 lines max.
8. **No Magic Strings or Numbers:** Use enums or typed constants. Every comparison must be against a named symbol.
9. **Definitions in Dedicated Files:** Types, enums, constants, and interfaces get their own files (e.g., `src/types/`, `src/enums/`), never defined inline next to first use.
10. **DRY is Priority One:** Duplicate logic across two sites must be extracted immediately.
11. **Component Modularity:** Small, reusable components. For features with 3+ components, produce a Mermaid diagram first.
12. **Immutable-First:** Assign every variable once at declaration. Never reassign except loop indices. Prefer `const`, `let`, `final`, `val` over mutable variables. Build result objects with spread or copy.
13. **Asset Naming:** Assets go to `assets/<NN-folder>/<NN-file>.<ext>` with two-digit sequence prefixes (e.g., `assets/01-icons/03-logo.svg`).
14. **No Inverted Complex Conditions:** Never use `!` on complex conditions containing AND/OR. Simplify or extract into named booleans.
15. **Boolean Return Wrapper:** Functions returning multiple values including a boolean must return a named struct/object (e.g. `{ data, isSuccess }`), never bare tuples `(int, bool)`.
16. **Strict Conditional Joins:** Never mix logical operators (e.g., OR with AND) and keep `if` conditions to at most one join (two operands).
17. **No Mixed Polarity:** Never combine positive and negative conditions in the same `if` statement (e.g., `if isA && !isB` is banned; extract `isConflict := isA && !isB`).
18. **No Explicit True Checks (TOTAL BAN):** NEVER evaluate a boolean explicitly against `true` or `false` (e.g., `if isReady == true` is FORBIDDEN; write `if isReady`).
19. **Enum Naming:** Every enum name MUST end with the suffix `Type` (e.g. `UserRoleType`), except in Rust where PascalCase is used without suffix. In Python, Enum classes use `PascalCase`, variable members use `UPPER_CASE` with underscores, and string values mirror member names exactly (e.g. `RegexPatternType.UPPERCASE = "UPPERCASE"`, `ExitCodeType.SUCCESS = 0`).
20. **Version Source of Truth:** `version.json` at root is the sole version authority. All languages import or read this file dynamically.

---

## 2. The Return New Line & Whitespace Concept (Mandatory)

The return new line and whitespace standard governs readability and clean code structure. This is mechanically checked by Rule R13–R20 and auto-fixed by `03-ai-scripts/05-guideline-autofixer.py`.

### Rule R13: Blank Line Before `return` / `throw` / `raise`

- **Rule:** Exactly ONE blank line before every `return`, `throw`, `raise`, or early exit statement when it is preceded by other statements in the block.
- **Exception (Single-Statement Blocks):** If `return` or `throw` is the **ONLY statement** in the block or function body, **NO** blank line is placed before it.

#### Multi-Language Examples:

```go
// Go: Single-statement function -> NO blank line
func GetDefaultPort() int {
    return 8080
}

// Go: Multi-statement function -> Blank line REQUIRED before return
func CalculateTotal(price int, tax int) int {
    subtotal := price + tax
    discount := calculateDiscount(subtotal)

    return subtotal - discount
}

// Go: Inside conditional blocks
func FindUser(ctx context.Context, params UserSearchParams) (*User, error) {
    if params.UserId == "" {
        return nil, apperror.New("empty user id") // Single-statement block: no blank line
    }

    user, err := repo.GetById(ctx, params.UserId)
    if err != nil {
        return nil, apperror.Wrap(err, "FindUser", map[string]any{"UserId": params.UserId})
    }

    return user, nil
}
```

```typescript
// TypeScript / React
export function calculateDiscount(price: number, isVip: boolean): number {
  if (isVip) {
    return price * 0.2; // Single statement: tight against brace
  }

  const standardRate = getStandardRate();
  const adjustedPrice = applyBaseDiscount(price, standardRate);

  return adjustedPrice; // Preceded by statements: blank line required
}
```

```python
# Python
def calculate_metrics(data_points: list[int]) -> int:
    if not data_points:
        return 0

    total_sum = sum(data_points)
    average_val = total_sum // len(data_points)

    return average_val
```

### Rule R14: Blank Line After Closing `}` Brace

- **Rule:** Exactly ONE blank line after every closing brace `}`, **unless** the next line is another `}`, `else`, `catch`, `finally`, or `case`.
- **In Python:** Exactly one blank line after a dedent ending a block, unless followed by `else`, `elif`, `except`, or `finally`.

```go
// CORRECT (Rule R14):
if err := validate(params); err != nil {
    return err
}

cmd := buildCommand(params)
if err := cmd.Run(); err != nil {
    return err
}

return nil
```

### Rule R15: Never Two Blank Lines in a Row

- Never place two consecutive blank lines anywhere in any file.

### Rule R16: No Padded Braces

- No blank line immediately after an opening `{` brace, and no blank line immediately before a closing `}` brace.

---

## 3. Boolean Principles (P1–P9)

1. **Prefixes:** Every boolean variable, function, parameter, or struct field MUST start with `is` or `has` ONLY (e.g. `isValid`, `hasAccess`, `isReady`, `hasData`); all other prefixes (`can`, `should`, `was`, `will`, `did`, `must`, etc.) are strictly BANNED.
2. **Positive Framing:** Never use negative names (`isNotReady`, `disableCache` are banned). Invert to positive equivalents (`isReady`, `isCacheEnabled`).
3. **No Inverted Success:** Never check `!response.isSuccess`. Use `response.isFail`.
4. **No Explicit True Checks (TOTAL BAN):** Never write `if isReady == true` or `if (hasMatch === true)`. Positive booleans MUST be implicit: `if isReady { ... }`.
5. **No Mixed Polarity:** Never combine positive and negative checks in the same condition (`if isA && !isB`). Extract to `isConflict := isA && !isB; if isConflict { ... }`.
6. **No Boolean Flag Parameters:** Never write `save(true)`. Split into `saveDraft()` / `savePublished()` or pass a configuration struct `save(SaveOptions{ IsDraft: true })`.

---

## 4. Function Signatures & Parameter Grouping

- **Rule R4 (Splitting):** If a function has > 3 parameters OR the signature line exceeds 100 characters, split to **one parameter per line**.
- **Rule R5 (Grouping):** If a function has > 4 parameters OR 2+ adjacent parameters of the same type, group them into a dedicated parameters struct/object (e.g. `SwapIpParams`).
- **Rule R9 (Call Site Mirror):** Any function call exceeding 100 characters or 4 arguments must be split to **one argument per line**.

```go
// Param Struct Example (Go)
type SwapIpParams struct {
    InterfaceName string `json:"InterfaceName"`
    OldIp         string `json:"OldIp"`
    NewIp         string `json:"NewIp"`
    SubnetMask    string `json:"SubnetMask"`
    IsDryRun      bool   `json:"IsDryRun"`
}

func SwapIp(ctx context.Context, params SwapIpParams) error { ... }
```

---

## 5. Error Management (`spec/03-error-manage/`)

- **Never Swallow Errors:** Every `catch` or error check must log with context and rethrow/return.
- **Wrap with Context:** Use `apperror.Wrap(err, "operationName", contextMap)` in Go, or `new AppError("message", { cause, op, context })` in TypeScript.
- **Universal Response Envelope:** APIs return `{ data, errors[], meta }`.
- **No Generic Errors:** Never throw base `Error` or `Exception`. Use domain-specific `AppError` classes with registered error codes.

---

## 6. R1–R21 AI Code Review & Linter Index

| Rule | Machine-Checkable Rule | Severity | Auto-Fixable Script |
| :--- | :--- | :--- | :--- |
| **R1** | Acronyms PascalCase (`Id`, `Url`, `Ip`, `Json`), never all-caps | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R2** | JSON / serialization keys are PascalCase (`{"UserId": "...", "IsActive": true}`) | Must Fix | Yes |
| **R3** | Boolean naming starts with is or has only (all other prefixes banned) | Must Fix | Yes |
| **R4** | Signature > 3 params or > 100 chars -> one param per line | Must Fix | Yes |
| **R5** | > 4 params or adjacent same-type params -> group into param struct | Must Fix | No (AI manual refactor) |
| **R6** | Every parameter is used, or discarded as `_` with explanatory comment | Must Fix | No (AI manual refactor) |
| **R7** | Every error propagated with context; no swallowed errors | Must Fix | No (AI manual refactor) |
| **R8** | No magic literals passed as arguments — extract named constants | Must Fix | No (AI manual refactor) |
| **R9** | Call > 100 chars or > 4 args -> one argument per line | Must Fix | Yes |
| **R10** | No boolean positional parameters -> use named struct or enum | Suggestion | No |
| **R11** | Go: `ctx` first parameter. C#: `cancellationToken` last parameter | Must Fix | Yes |
| **R12** | Return values documented; multi-return meanings unambiguous | Suggestion | No |
| **R13** | One blank line before `return`/`throw`, unless sole statement in block | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R14** | One blank line after closing `}`, unless followed by `}`, `else`, `catch`, etc. | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R15** | Never two blank lines in a row anywhere | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R16** | No blank line immediately after `{` or before `}` | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R17** | Exactly one blank line between top-level declarations | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R18** | Import grouping: stdlib -> third-party -> first-party absolute -> relative | Must Fix | Yes |
| **R19** | Trailing newline at EOF; no trailing whitespace | Must Fix | Yes (`05-guideline-autofixer.py`) |
| **R20** | No section-separator blank lines inside functions (refactor instead) | Suggestion | No |
| **R21** | No comments on struct fields unless explaining units or non-obvious defaults | Must Fix | No |

---

## 7. Step-by-Step AI Execution Workflow

When tasked with auditing, reviewing, or fixing coding guidelines across a codebase, follow these sequential steps:

1. **Step 1: Automated Pre-Pass (Deterministic AST Autofixers):**
   - Run the autofixers to instantly eliminate 85% of mechanical violations:
     ```bash
     python 03-ai-scripts/05-guideline-autofixer.py <target-dir>
     python 03-ai-scripts/08-naming-autofixer.py <target-dir>
     python 03-ai-scripts/04-newline-fixer.py <target-dir>
     python 03-ai-scripts/07-relative-path-fixer.py <target-dir>
     ```
2. **Step 2: Run Linters for Violations:**
   - Execute `python linter-scripts/validate-guidelines.py` and `python linter-scripts/check-boolean-guidelines.py` to identify remaining violations.
3. **Step 3: Sequential Manual Fixes (Bounded Micro-Tasks):**
   - Address remaining non-autofixable violations (R5 param structs, R6 dead params, R7 error context, R8 magic constants) file by file.
   - Respect the 15-line function cap and flatten all nested conditionals.
4. **Step 4: Local CI/CD Pipeline Quality Gate:**
   - Execute targeted file linters / autofixers on modified files ensuring `exit 0`. DO NOT run the full `06-cicd-local-runner.py` during routine guideline turns.
5. **Step 5: File Change Summary:**
   - Output a detailed summary in chat listing exactly which files changed, what changed, and why.

---

## 8. Anti-Hallucination & Checklist Execution

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> Operate using the 3-Phase Bounded Execution Model:
>
> 1. **Phase 1: Read & Understand (Isolated Loop):** First turn is purely exploratory. Do NOT write code.
> 2. **Phase 2: Bounded Micro-Tasking:** Fix one file or section at a time.
> 3. **Phase 3: Multi-Agent Parallelization:** Spawn sub-agents with micro-boundaries (single-file bounding box).

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## 9. Non-Negotiable Review Checklist

- [ ] **Return New Line Concept (R13-R16):** I verified that every `return`, `throw`, and `raise` has a blank line before it (unless sole statement in block), every closing brace `}` has a blank line after it, no two blank lines exist in a row, and no padded braces exist.
- [ ] **No Explicit True Checks (P4):** Absolutely zero `== true`, `=== true`, `!= false`, `!== false` comparisons exist.
- [ ] **No Mixed Polarity (P5):** No mixed positive and negative conditions in `if` statements.
- [ ] **Acronyms & PascalCase (R1, R2):** All acronyms (`Id`, `Url`, `Ip`, `Json`) and serialization keys use PascalCase.
- [ ] **Boolean Prefixes (R3):** All booleans start with is or has only (all other prefixes banned). No negative boolean names.
- [ ] **Function Decomposition & Signatures (R4, R5):** All functions <= 15 lines decomposed via 3-Stage Blueprint (Guard -> Core Logic -> Envelope) without logic drift; parameter structs for > 3 arguments.
- [ ] **Circular Dependency Prevention:** All extracted types/enums reside in leaf packages (`domain/types` or `types/`) with zero circular dependency cycles.
- [ ] **Polyglot & React Compliance:** Rust match expressions, C# Task/records, PHP BackedEnums, React structuredClone & object hook returns.
- [ ] **Error Handling (R7):** All errors are wrapped with context (`apperror.Wrap`) and not swallowed.
- [ ] **No Magic Constants (R8):** All magic strings/numbers are extracted to named constants.
- [ ] **Strict Lowercase Filenames:** All generated or modified files use strictly lowercase naming (`readme.md`, `agents.md`, `skill.md`).
- [ ] **Tooling Execution:** I ran `03-ai-scripts/05-guideline-autofixer.py` and verified clean output with `python linter-scripts/validate-guidelines.py`.
- [ ] **Targeted Verification:** All modified files pass targeted linters / autofixers cleanly with exit 0. (Full CI runner `06-cicd-local-runner.py` is banned in routine turns).
- [ ] **Temp Storage & Pre-Build Clean (R17):** All OS/user temporary files are scoped under `<temp>/gitmap/<category>/` (build, test, purge, downloads). Before running any build, previous build artifacts in the target directory are purged to respect storage and prevent disk bloat.
- [ ] **File Change Summary:** I provided a detailed summary in chat of what files changed, what changed inside them, and why.

---

## Metadata

- slug: coding-guidelines
- status: active
