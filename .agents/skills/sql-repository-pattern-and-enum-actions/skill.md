---
name: sql-repository-pattern-and-enum-actions
description: Enforces the SQL repository pattern in every language and enum-driven action dispatch - all SQL behind repository interfaces, every action an enum member dispatched by exhaustive switch, enums authored with String, Parse/TryParse, IsValid, and All/Values utils, typed errors for unknown values.
---

# SQL Repository Pattern & Enum-Driven Actions (`sql-repository-pattern-and-enum-actions`)

Source spec: 02-spec/21-app/03-execute-n-steps-letterly-prompt-formatter/04-sql-repository-pattern-and-enum-actions-spec.md

## 1. Core Mandate

These rules are non-negotiable in every repository and every language.

- **Rule A — Repository pattern:** all SQL and data access lives behind a
  repository interface. Services, handlers, and use-cases never touch SQL,
  query builders, or ORMs directly.
- **Rule B — Enum-driven actions:** every action, operation, command kind,
  and event kind is a member of a named enum, dispatched by an exhaustive
  switch or match — never by raw strings or magic numbers.
- **Rule C — Enum authoring:** enums are defined once, in a shared module,
  with explicit stable values, and ship the full util-method set:
  `String()`, `Parse` / `TryParse`, `IsValid`, and `All` / `Values`.
- **Rule D — Typed errors:** unknown enum values, failed parses, unhandled
  switch cases, and repository failures all surface as typed errors with
  context — never a panic, never a silent default, never a swallowed error.

If a change touches SQL, an action, or an enum, this skill applies in full.
Partial compliance is a defect, not a style choice.

## 2. Rule A - SQL Repository Pattern (Every Language)

### A.1 One repository interface per aggregate

- Declare one repository interface per aggregate or table group.
- The interface is the only data-access surface services may depend on.
- Concrete implementations hold the driver, connection, and dialect.

### A.2 No SQL outside repository implementations

- Services, handlers, and use-cases must contain zero SQL strings.
- Query builders and ORM calls are equally banned outside implementations.
- No caller builds, concatenates, or parameterizes a query by hand.

### A.3 Intention-named methods only

- Name methods by intent: `FindByID`, `ListActive`, `Insert`, `Update`,
  `SoftDelete`.
- Never expose `RunQuery`, `ExecuteRaw`, or generic execution methods.
- Never accept stringly-typed table or column names from callers.
- Parameters are typed domain values; return values are domain types.

### A.4 Dialect stays inside the implementation

- SQL dialect, quoting, placeholders, and driver specifics live only in
  the repository implementation.
- Swapping the database means swapping the implementation only.
- Callers, services, and tests against the interface do not change.

### A.5 Unit-of-work for transactions

- Multi-repository transactions go through a unit-of-work abstraction.
- No ad-hoc transaction handling scattered across services or handlers.
- The unit-of-work owns begin, commit, and rollback for the whole case.

### Language-agnostic shape

The shape is identical in Go, TypeScript, Python, Rust, and PHP:

- an interface (or protocol / trait / abstract class) per aggregate,
- one implementation per database or dialect,
- services depending only on the interface,
- a unit-of-work wherever one use-case spans multiple repositories.

Only the idiom changes; the boundary never does.

## 3. Rule B - Every Action From an Enum

### B.1 Every action is an enum member

- Every action, operation, command kind, and event kind is a member of a
  named enum.
- Dispatch never compares raw strings or magic numbers.

### B.2 Exhaustive switch or match dispatch

- Dispatch happens by switching or matching on the enum value.
- Switches over action enums are exhaustive: the compiler or linter must
  flag a missing case.
- The default branch raises a typed error (see Rule D) — never a silent
  no-op, never a fall-through success.

### B.3 Validate before dispatch

- Run `IsValid` / the membership check before the switch runs.
- Unchecked external input never reaches the dispatch.
- An invalid value fails at validation with a typed error, not mid-case.

### B.4 Adding an action is two steps

- Adding an action = adding one enum member + one switch case.
- No other registration step, registry edit, or string table may be
  required.
- If a second registration exists, remove it — the enum is the registry.

### Scope clarifier

Pure data lookups are not action dispatch. A static table such as an
extension-to-MIME mapping holds data, not behaviour, and may stay
map-based. The moment a value selects behaviour, a handler, or a side
effect, it is an action and must be an enum member under this rule.

## 4. Rule C - Enum Authoring & Util-Methods Checkup

### Definition rules

- One enum per concept, defined in the shared types / constants module
  of the repo — never inline in a handler, use-case, or test.
- Use the language's idiomatic casing for the type and its members.
- Serialized enums carry explicit, stable string values.
- A serialized value never changes once shipped; renaming a member keeps
  its value, and a retired value is never reused.
- Where the language needs it, declare an explicit unknown / zero
  sentinel — and never treat that sentinel as a valid action.

### Usage rules

- No raw literals anywhere an enum exists: comparisons, map keys, logs,
  and payloads all use the member or its canonical string value.
- External input enters only through the enum's parse util, which
  returns the member or a typed error.
- Never cast, coerce, or string-compare external input into an enum.

### 4.1 Util-Methods Checkup

Every enum ships all four utils, or carries a documented exception:

1. **`String()` / label** — the canonical string form of a member, used
   for logs, payloads, and persistence instead of ad-hoc formatting.
2. **`Parse` / `TryParse`** — string to member; unknown input returns a
   typed error carrying the offending value (see Rule D).
3. **`IsValid`** — the membership check called before every dispatch
   (see Rule B.3) and on any value reconstructed from storage.
4. **`All` / `Values`** — the full member list; it drives exhaustive
   tests and UI lists so neither hand-maintains its own copy.

Checkup rule: when a change touches an enum, audit all four utils. A
missing util is added in the same change — never deferred to a follow-up.
If a util genuinely does not apply, record the exception next to the
enum definition in the same change.

## 5. Rule D - Error-Management Tie-Ins

### D.1 Unknown values are typed errors

- Unknown enum value, failed parse, or unhandled switch case produces a
  typed error under the repo's AppError / result convention.
- The error message includes the offending value verbatim.
- Never panic on an unknown value; never fall back to a silent default.

### D.2 Repository errors are wrapped at the boundary

- Wrap repository errors with context at the repository boundary,
  before they cross into services.
- The context names the repository and the operation that failed.
- Services receive the wrapped typed error, never a raw driver error.

### D.3 Canonical examples

- Enum-parse failure and repository-wrap are the canonical examples in
  error prompts and error-handling skills.
- New error examples follow those two patterns: typed error, offending
  value or operation in the message, context added once at the boundary.

## 6. Verification Checklist

- [ ] Zero SQL, query-builder, or ORM calls outside repository implementations.
- [ ] One repository interface exists per aggregate / table group.
- [ ] Repository methods are intention-named only; no `RunQuery` / `ExecuteRaw`.
- [ ] No stringly-typed table or column parameters cross the repository boundary.
- [ ] Multi-repository transactions go through a unit-of-work abstraction.
- [ ] Every action / operation / command kind / event kind is an enum member.
- [ ] Action dispatch uses exhaustive switch / match; missing cases fail checks.
- [ ] Every switch default raises a typed error; no silent no-op defaults.
- [ ] No raw literals remain where an enum exists (comparisons, keys, logs, payloads).
- [ ] External input reaches an enum only through its parse util.
- [ ] Util-methods checkup - String()/label present for every enum (or documented exception).
- [ ] Util-methods checkup - Parse/TryParse present; unknown input returns a typed error carrying the offending value.
- [ ] Util-methods checkup - IsValid present and called before every action dispatch.
- [ ] Util-methods checkup - All/Values present and used for exhaustive tests/UI lists.
- [ ] Unknown enum values and failed parses never panic and never default silently.
- [ ] Repository errors are wrapped with repository + operation context at the boundary.
