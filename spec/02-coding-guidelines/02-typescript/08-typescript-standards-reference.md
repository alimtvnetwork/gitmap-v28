# TypeScript Coding Standards

> **Version:** 3.0.0  
> **Updated:** 2026-03-09  
> **Applies to:** All frontend TypeScript/React code  
> **Priority:** CRITICAL — These rules override all other conventions

---

## 1. Generics First — The Cardinal Rule

**Generics are the ONLY acceptable approach for parameterized types.** Never use `any`, `unknown`, `Record<string, unknown>`, or loose interfaces where a generic can express the constraint.

### Rule 1.1: All reusable functions MUST be generic

```typescript
// ❌ FORBIDDEN — loose typing
function fetchData(endpoint: string): Promise<unknown> { ... }
function parseResponse(data: unknown): Record<string, unknown> { ... }

// ✅ REQUIRED — generic with constraints
function fetchData<T>(endpoint: string): Promise<T> { ... }
function parseResponse<T extends object>(data: string): T { ... }
```

### Rule 1.2: API response types MUST use generic envelope

```typescript
// ❌ FORBIDDEN
interface RawEnvelope {
  Results: unknown[];
}

// ✅ REQUIRED
interface RawEnvelope<T = never> {
  Status: EnvelopeStatus;
  Attributes: EnvelopeAttributes;
  Results: T[];
  Navigation?: EnvelopeNavigation;
  Errors?: EnvelopeErrors;
}
```

### Rule 1.3: Collection utilities MUST be generic

```typescript
// ❌ FORBIDDEN
function buildQuery(params: Record<string, string | number | undefined>): string { ... }

// ✅ REQUIRED
function buildQuery<T extends Record<string, string | number | undefined | null>>(params: T): string { ... }
```

### Rule 1.4: Hook factories MUST propagate generics

```typescript
// ❌ FORBIDDEN
function useApiQuery(key: string): { data: unknown } { ... }

// ✅ REQUIRED
function useApiQuery<T>(key: string[]): { data: T | undefined; isLoading: boolean } { ... }
```

---

## 2. Zero Tolerance for `any` and Untyped Patterns

### Rule 2.1: `any` is PROHIBITED everywhere

No exceptions. Not in catch blocks, not in type assertions, not in generic defaults.

```typescript
// ❌ FORBIDDEN — all of these
catch (err: any) { ... }
const x = value as any;
getQueryData<any>(key);
(result as any)?.deleted;
v as any;

// ✅ REQUIRED
catch (err) {
  const message = err instanceof Error ? err.message : String(err);
}
const x = value as SpecificType;
getQueryData<DashboardStats>(key);
```

### Rule 2.2: `unknown` is acceptable ONLY at parse boundaries

`unknown` may appear in:
- JSON parsing entry points (immediately narrowed via type guard)
- Error catch blocks (without `: any` annotation — bare `catch (err)`)
- Internal type narrowing functions (e.g., `isEnvelope(obj: unknown)`)

`unknown` MUST NOT appear in:
- Component props, hook return types, store state
- API method return types (use generics instead)
- Exported function signatures

### Rule 2.3: `Record<string, unknown>` is PROHIBITED in API signatures

```typescript
// ❌ FORBIDDEN
createRemoteSnapshot: (siteId: number, opts?: Record<string, unknown>) => ...
updateSettings: (settings: Record<string, unknown>) => ...

// ✅ REQUIRED — use specific interfaces
interface CreateSnapshotOptions {
  name?: string;
  scope?: SnapshotScope;
  snapshotType?: SnapshotType;
  parentId?: number;
  tables?: string[];
}
createRemoteSnapshot: (siteId: number, opts?: CreateSnapshotOptions) => ...
```

---

## 3. No Magic Strings, No Magic Numbers

### Rule 3.1: All string literals used as identifiers MUST come from constants or enums

```typescript
// ❌ FORBIDDEN — magic strings
if (status === "connected") { ... }
if (action === "self-update") { ... }
toast.success("Cleanup complete");

// ✅ REQUIRED — constants or enums
const enum ConnectionStatus {
  Connected = "connected",
  Disconnected = "disconnected",
  Unknown = "unknown",
}

const enum SnapshotAction {
  Create = "create",
  Restore = "restore",
  Delete = "delete",
  Export = "export",
  Import = "import",
  Cleanup = "cleanup",
}

if (status === ConnectionStatus.Connected) { ... }
if (action === SnapshotAction.Create) { ... }
```

### Rule 3.2: All numeric literals with semantic meaning MUST be named constants

```typescript
// ❌ FORBIDDEN — magic numbers
staleTime: 60_000,
const limit = 25;
setTimeout(fn, 5000);

// ✅ REQUIRED
const STALE_TIME_MS = 60_000 as const;
const DEFAULT_PAGE_SIZE = 25 as const;
const POLLING_INTERVAL_MS = 5_000 as const;

staleTime: STALE_TIME_MS,
const limit = DEFAULT_PAGE_SIZE;
setTimeout(fn, POLLING_INTERVAL_MS);
```

### Rule 3.3: String unions MUST be proper enums with PascalCase values

```typescript
// ❌ FORBIDDEN — inline string unions
status: "success" | "failed" | "partial";
type: "publish" | "snapshot" | "plugin" | "config" | "connection";

// ❌ FORBIDDEN — named type aliases (still string unions)
type PublishStatus = "success" | "failed" | "partial";

// ✅ REQUIRED — proper enums with Type suffix and PascalCase values
enum PublishStatusType {
  Success = "Success",
  Failed = "Failed",
  Partial = "Partial",
}

enum ActivityType {
  Publish = "Publish",
  Snapshot = "Snapshot",
  Plugin = "Plugin",
  Config = "Config",
  Connection = "Connection",
}
```

**Rule:** Never use string union types (`"a" | "b"`) — always use proper `enum` with PascalCase values and a `Type` suffix on the enum name.

---

## 4. Specific Type Rules for Common Patterns

### 4.1: Error handling — structured catch blocks

```typescript
// ✅ The ONLY acceptable catch pattern
try {
  await apiCall();
} catch (err) {
  const message = err instanceof Error ? err.message : String(err);
  // OR for ApiClientError:
  if (isApiClientError(err)) {
    handleApiError(err.apiError);
  } else {
    handleGenericError(message);
  }
}
```

### 4.2: React Query — typed query data

```typescript
// ❌ FORBIDDEN
queryClient.getQueryData<any>(["dashboard-stats"]);
(data as any).entries;

// ✅ REQUIRED
queryClient.getQueryData<DashboardStats>(["dashboard-stats"]);
```

### 4.3: Component prop drilling — avoid `Record<string, unknown>`

```typescript
// ❌ FORBIDDEN
metadata: Record<string, unknown>;

// ✅ REQUIRED — specific metadata types per domain
interface PublishMetadata {
  pluginName: string;
  version: string;
  filesUpdated: number;
  sessionId?: string;
}

interface SnapshotMetadata {
  snapshotType: SnapshotType;
  tables?: number;
  size?: number;
}

type ActivityMetadata = PublishMetadata | SnapshotMetadata | PluginMetadata | ConfigMetadata | ConnectionMetadata;
```

---

## 5. Function Size — Max 15 Lines

> **Canonical source:** [Cross-Language Code Style](../01-cross-language/04-code-style/00-overview.md) — Rule 6

Every function/method body must be **15 lines or fewer**. Extract logic into small, well-named helper functions.

```typescript
// ❌ FORBIDDEN: 20+ line function
const handleSubmit = async (data: FormData) => {
    // validation, API call, state update, toast... all inline
};

// ✅ REQUIRED: Decomposed
const handleSubmit = async (data: FormData) => {
    const validated = validateFormData(data);
    const result = await submitToApi(validated);
    updateLocalState(result);
    showSuccessToast(result.message);
};
```

---

## 6. Zero Nested `if` — Absolute Ban

> **Canonical source:** [Cross-Language Code Style](../01-cross-language/04-code-style/00-overview.md) — Rule 2 & 7

Nested `if` blocks are **absolutely forbidden** — zero tolerance, no exceptions. Flatten with early returns or combined conditions.

```typescript
// ❌ FORBIDDEN: Nested if
if (response) {
    if (response.status >= 400) {
        handleError(response);
    }
}

// ✅ REQUIRED: Early return
if (!response) {
    return;
}

if (response.status >= 400) {
    handleError(response);
}
```

---

## 7. Enforcement

- **TypeScript strict mode:** Must be enabled (`strict: true` in tsconfig)
- **ESLint rules (REQUIRED):**
  - `@typescript-eslint/no-explicit-any`: `error`
  - `@typescript-eslint/no-unsafe-assignment`: `error`
  - `@typescript-eslint/no-unsafe-member-access`: `error`
  - `@typescript-eslint/no-unsafe-call`: `error`
  - `@typescript-eslint/no-unsafe-return`: `error`
- **Code review:** Any PR introducing `any`, bare `unknown` in public APIs, or magic strings/numbers must be rejected
- **Exceptions:** Must include a `// SAFETY:` comment explaining why and a `// TODO:` for removal with ticket reference

---

## 9. Generics Reference — When to Use What

| Scenario | Pattern |
|----------|---------|
| API response parsing | `parseEnvelope<T>(env: RawEnvelope<T>): ApiResponse<T>` |
| Data fetching hooks | `useApiQuery<T>(key: string[]): QueryResult<T>` |
| Form state | `useForm<TFormValues extends FieldValues>()` |
| List rendering | `function DataTable<T extends { id: string }>(props: { data: T[] })` |
| Store slices | `createSlice<TState>(initialState: TState)` |
| Utility functions | `function groupBy<T, K extends keyof T>(items: T[], key: K): Map<T[K], T[]>` |

---

## 8. No Raw Negations — Use Positive Guard Functions

> **Canonical source:** [No Raw Negations](../01-cross-language/12-no-negatives.md)

**Never use `!` on a function call in a condition.** Wrap every negative check in a positively named guard function.

```typescript
// ❌ FORBIDDEN
if (!response.ok) { handleError(response); }
if (!array.includes(item)) { array.push(item); }
if (!fs.existsSync(path)) { throw new Error('Missing'); }

// ✅ REQUIRED
if (isResponseFailed(response)) { handleError(response); }
if (isItemMissing(array, item)) { array.push(item); }
if (isFileMissing(path)) { throw new Error('Missing'); }
```

**Utility location:** `src/utils/guards.ts` — see canonical spec for full guard function table.

---

## 9. `isDefined()` and `isDefinedAndValid()` — Positive Null/Existence Guards

> **Cross-language parity:** [Go Standards §IsDefined](../03-golang/04-golang-standards-reference/00-overview.md)

Raw `!== null` / `!== undefined` combined with negation or nested validity checks creates cognitive overhead. Use positive guard functions that express intent clearly.

### `isDefined()` — Value Existence Check

Returns `true` when the value is not `null` or `undefined`. Replaces `!= null` checks.

```typescript
// ❌ FORBIDDEN: Raw null/undefined checks
if (config !== null && config !== undefined) {
  applyConfig(config);
}

// ❌ FORBIDDEN: Double negation
if (!!config) {
  applyConfig(config);
}

// ✅ REQUIRED: Positive existence check
if (isDefined(config)) {
  applyConfig(config);
}
```

### `isDefinedAndValid()` — Existence + Validity Combined

Returns `true` when the value exists AND passes validation. Replaces nested null+validity checks.

```typescript
// ❌ FORBIDDEN: Nested null + validity check
if (config !== undefined) {
  if (config.isValid()) {
    applyConfig(config);
  }
}

// ❌ FORBIDDEN: Compound with null check
if (config != null && config.isValid()) {
  applyConfig(config);
}

// ✅ REQUIRED: Single positive guard
if (isDefinedAndValid(config)) {
  applyConfig(config);
}
```

### Implementation — `src/utils/guards.ts`

```typescript
/** Value exists (not null/undefined) */
export function isDefined<T>(value: T | null | undefined): value is T {
  return value !== null && value !== undefined;
}

/** Value exists AND passes validation */
export function isDefinedAndValid<T extends { isValid(): boolean }>(
  value: T | null | undefined
): value is T {
  return value !== null && value !== undefined && value.isValid();
}

/** Value is absent (null or undefined) */
export function isEmpty<T>(value: T | null | undefined): value is null | undefined {
  return value === null || value === undefined;
}
```

### Guard Function Table

| Guard Function | Replaces | Description |
|---------------|----------|-------------|
| `isDefined(x)` | `x !== null && x !== undefined` | Value exists |
| `isDefinedAndValid(x)` | `x != null && x.isValid()` | Value exists AND passes validation |
| `isEmpty(x)` | `x === null \|\| x === undefined` | No value (absent) |

### Type Narrowing Benefit

These guards use TypeScript type predicates (`value is T`), providing automatic type narrowing:

```typescript
function processConfig(config: Config | null) {
  // ❌ Without guard — config is still Config | null inside the block
  if (config !== null) {
    config.apply(); // works but no semantic intent
  }

  // ✅ With guard — config is narrowed to Config, AND intent is clear
  if (isDefined(config)) {
    config.apply(); // TypeScript knows config is Config
  }
}
```

### Real-World Example

```typescript
// Service layer — checking optional input
function updateSite(input: UpdateSiteInput): Result<Site> {
  if (isDefined(input.config)) {
    if (isDefinedAndValid(input.config)) {
      applyConfig(input.config);
    } else {
      return Result.fail("E3010", "invalid site config");
    }
  }
  // ...
}
```

---

## 8. Discriminated Unions — Named Interfaces Required

**All union-type action/event objects MUST use extracted named interfaces with enum discriminators.** Never use inline type literals or string-literal discriminators.

### Rule 8.1: Every variant MUST be a named interface

```typescript
// ❌ FORBIDDEN — inline variants with string literals
type Action =
  | { type: "ADD_TOAST"; toast: ToasterToast }
  | { type: "DISMISS"; toastId?: string };

// ✅ REQUIRED — named interfaces with enum discriminator
enum ActionType {
  AddToast = "AddToast",
  DismissToast = "DismissToast",
}

interface AddToastAction {
  type: ActionType.AddToast;
  toast: ToasterToast;
}

interface DismissToastAction {
  type: ActionType.DismissToast;
  toastId?: string;
}

type ToastAction = AddToastAction | DismissToastAction;
```

### Rule 8.2: Use dot notation for enum access

```typescript
// ❌ FORBIDDEN
if (action.type === "ADD_TOAST") { ... }
if (action.type === ActionType["AddToast"]) { ... }

// ✅ REQUIRED
if (action.type === ActionType.AddToast) { ... }
```

### Rule 8.3: Enum values MUST be PascalCase

```typescript
// ❌ FORBIDDEN — UPPER_SNAKE_CASE, camelCase, or _camelCase
enum ActionType {
  ADD_TOAST = "ADD_TOAST",
  addToast = "addToast",
  _addToast = "_addToast",
}

// ✅ REQUIRED — PascalCase
enum ActionType {
  AddToast = "AddToast",
}
```

### Full Example

```typescript
enum ActionType {
  AddToast = "AddToast",
  UpdateToast = "UpdateToast",
  DismissToast = "DismissToast",
  RemoveToast = "RemoveToast",
}

interface AddToastAction {
  type: ActionType.AddToast;
  toast: ToasterToast;
}

interface UpdateToastAction {
  type: ActionType.UpdateToast;
  toast: Partial<ToasterToast>;
}

interface DismissToastAction {
  type: ActionType.DismissToast;
  toastId?: string;
}

interface RemoveToastAction {
  type: ActionType.RemoveToast;
  toastId?: string;
}

type ToastAction =
  | AddToastAction
  | UpdateToastAction
  | DismissToastAction
  | RemoveToastAction;
```

---

## Enum Inventory

All standardized TypeScript enums for the frontend. Each enum has a dedicated spec file with values, usage examples, and Go parity mapping.

| Enum | Values | Spec File | Go Equivalent |
|------|--------|-----------|---------------|
| `HttpMethod` | `Get`, `Head`, `Post`, `Put`, `Patch`, `Delete`, `Options` | [http-method-enum.md](./05-http-method-enum.md) | `pkg/enums/httpmethodtype` |
| `ExecutionStatus` | `Idle`, `Running`, `Paused`, `Completed`, `Failed`, `Cancelled` | [execution-status-enum.md](./03-execution-status-enum.md) | `pkg/enums/executionstatustype` |
| `ConnectionStatus` | `Connected`, `Disconnected`, `Connecting`, `Reconnecting`, `Error` | [connection-status-enum.md](./01-connection-status-enum.md) | `pkg/enums/connectionstatustype` |
| `ExportStatus` | `Pending`, `Processing`, `Completed`, `Failed` | [export-status-enum.md](./04-export-status-enum.md) | `pkg/enums/exportstatustype` |
| `MessageStatus` | `Pending`, `Streaming`, `Completed`, `Error` | [message-status-enum.md](./06-message-status-enum.md) | `pkg/enums/messagestatustype` |
| `EntityStatus` | `Active`, `Inactive`, `Draft`, `Archived` | [entity-status-enum.md](./02-entity-status-enum.md) | `pkg/enums/entitystatustype` |

**Location:** All frontend enums live in `src/lib/enums/` (e.g., `src/lib/enums/http-method-type.ts`).

---

## Cross-References

- [No Raw Negations](../01-cross-language/12-no-negatives.md) — Positive guard functions (all languages)
- [Cross-Language Code Style](../01-cross-language/04-code-style/00-overview.md) — Braces, nesting & spacing rules (canonical)
- [Function Naming](../01-cross-language/10-function-naming.md) — No boolean flag parameters (all languages)
- [Boolean Flag Method Splitting](../01-cross-language/24-boolean-flag-methods.md) — Split bool-flag methods into two named methods
- [Strict Typing](../01-cross-language/13-strict-typing.md) — Type declarations & docblock rules (all languages)
- [DRY Principles](../01-cross-language/08-dry-principles.md)
- [Golang Standards](../03-golang/04-golang-standards-reference/00-overview.md)
- [Response Envelope Spec](../../03-error-manage/02-error-architecture/05-response-envelope/envelope.schema.json)

---

*TypeScript standards v3.2.0 — generics-first, zero-any, no-magic-strings, max-15-lines, zero-nesting, isDefined-guards — 2026-02-25*
