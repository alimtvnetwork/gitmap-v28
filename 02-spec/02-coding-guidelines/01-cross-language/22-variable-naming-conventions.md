# Cross-Language Variable Naming Conventions (AI Execution Prompt)

> **/goal** Standardize variable, collection, map, and loop variable naming conventions across PHP, TypeScript, and Go to maximize readability and eliminate ambiguity.
> **/learn** Understand singular-plural distinctions for collections, loop variable derivations, `Map`/`By` suffixes for lookup tables, affirmative boolean naming, prohibited abbreviations, and absolute bans on generic garbage names (`data`, `temp`, `obj`).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce singular nouns for single items and plural nouns for collections (arrays, lists, slices, sets) across all languages.
- [ ] `/learn` Enforce loop variables using the singular noun of the target collection being iterated.
- [ ] `/goal` Require `Map` suffix or `By[Key]` naming patterns for dictionary and map lookup variables.
- [ ] `/learn` Ban generic placeholder names (`data`, `temp`, `obj`, `val`, `item` outside loops) in favor of domain-explicit naming.
- [ ] `/learn` Ban arbitrary abbreviations (except universally accepted ones like `id`, `url`, `err`, `ctx`) and verify compliance via CI linters.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0
> **Updated:** 2026-03-31
> **Applies to:** PHP, TypeScript, Go

---

## Overview

Rules for naming variables, collections, and maps to maximize readability and reduce ambiguity. These complement [Function Naming](./10-function-naming.md) and [Boolean Principles](./02-boolean-principles/readme.md).

---

## Keywords

`variable-naming` · `singular-plural` · `collections` · `maps` · `dictionaries` · `naming-conventions` · `readability`

---

## Rule 1: Singular for Single Items, Plural for Collections

A variable holding **one item** uses the singular noun. A variable holding **multiple items** (array, list, slice, set) uses the plural.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Plural for single item
$users = $this->repository->findById($id);

// ❌ FORBIDDEN: Singular for collection
$user = $this->repository->findAll();

// ✅ REQUIRED
$user = $this->repository->findById($id);
$users = $this->repository->findAll();
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
const items = cart.getSelectedItem();
const item = await fetchAllItems();

// ✅ REQUIRED
const item = cart.getSelectedItem();
const items = await fetchAllItems();
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
users, err := repo.FindByID(ctx, id)  // returns one user

// ✅ REQUIRED
user, err := repo.FindByID(ctx, id)
users, err := repo.FindAll(ctx)
```

---

## Rule 2: Loop Variables Use Singular of the Collection Name

When iterating a collection, the loop variable must be the **singular form** of the collection name.

```php
// ❌ FORBIDDEN: Generic or mismatched loop variable
foreach ($users as $item) {
foreach ($orders as $o) {

// ✅ REQUIRED: Singular of collection name
foreach ($users as $user) {
foreach ($orders as $order) {
```

```typescript
// ❌ FORBIDDEN
for (const x of plugins) {
items.forEach((el) => {

// ✅ REQUIRED
for (const plugin of plugins) {
items.forEach((item) => {
```

```go
// ❌ FORBIDDEN
for _, v := range users {

// ✅ REQUIRED
for _, user := range users {
```

---

## Rule 3: Maps and Dictionaries — Use `Map` or `By` Suffix

Variables holding key-value lookups must indicate their nature with a `Map` suffix or `By[Key]` pattern.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Ambiguous — is this a list or a map?
$users = [];
$users[$id] = $user;

// ✅ REQUIRED: Clear map naming
$usersById = [];
$usersById[$id] = $user;

// ✅ ALSO OK
$userMap = [];
$userMap[$id] = $user;
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Unclear whether it's an array or record
const results: Record<string, Result> = {};

// ✅ REQUIRED
const resultsById: Record<string, Result> = {};
const resultMap = new Map<string, Result>();

// ✅ ALSO OK: Descriptive key type
const pricesByProductId: Record<string, number> = {};
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
users := make(map[string]*User)

// ✅ REQUIRED
usersByID := make(map[string]*User)
userMap := make(map[string]*User)
```

---

## Rule 4: Boolean Variables — `is`, `has` Prefix (99%), `should` Rare

Boolean variables **must** use a prefix that reads as a yes/no question. See [Boolean Principles](./02-boolean-principles/readme.md) for the full specification.

```typescript
// ❌ FORBIDDEN
const loading = true;
const admin = user.role === 'admin';

// ✅ REQUIRED
const isLoading = true;
const isAdmin = user.role === 'admin';
```

---

## Rule 5: Avoid Abbreviations — Spell Out Intent

Variable names must not use abbreviations unless universally understood (`id`, `url`, `ctx`, `err`, `req`, `res`).

```typescript
// ❌ FORBIDDEN
const usr = await fetchUser(id);
const btn = document.querySelector('.submit');
const mgr = new ConnectionManager();

// ✅ REQUIRED
const user = await fetchUser(id);
const submitButton = document.querySelector('.submit');
const connectionManager = new ConnectionManager();
```

### Allowed Abbreviations

| Abbreviation | Meaning |
|-------------|---------|
| `id` | Identifier |
| `url` | Uniform Resource Locator |
| `ctx` | Context (Go convention) |
| `err` | Error |
| `req` | Request (in handler scope) |
| `res` | Response (in handler scope) |
| `db` | Database |
| `fn` | Function (callback parameter) |
| `i`, `j`, `k` | Loop index (numeric loops only) |

---

## Rule 6: Absolute Ban on Generic Garbage Names (`data`, `temp`, `obj`, `val`)

Variable names **must** describe the specific domain entity they hold. Never use generic placeholder or throwaway names such as `data`, `temp`, `tmp`, `obj`, `val`, or `result` without domain qualification.

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Generic garbage names
const data = await fetchUser(id);
const temp = calculateSubtotal(items);
const obj = parseConfig(raw);

// ✅ REQUIRED: Domain-explicit names
const userProfile = await fetchUser(id);
const subtotal = calculateSubtotal(items);
const parsedConfig = parseConfig(raw);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Generic garbage names
data, err := repo.FindByID(ctx, id)
temp := computeTax(total)
obj := decodeSession(token)

// ✅ REQUIRED: Domain-explicit names
userRecord, err := repo.FindByID(ctx, id)
calculatedTax := computeTax(total)
sessionAuth := decodeSession(token)
```

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Generic garbage names
$data = $this->userRepo->findById($id);
$temp = $this->calculateDiscount($order);
$obj = $this->serializer->deserialize($payload);

// ✅ REQUIRED: Domain-explicit names
$user = $this->userRepo->findById($id);
$discountAmount = $this->calculateDiscount($order);
$sessionPayload = $this->serializer->deserialize($payload);
```

---

## Checklist Summary

```
[ ] Singular for single items, plural for collections
[ ] Loop variable is singular of collection name
[ ] Maps/dicts use `Map` suffix or `By[Key]` pattern
[ ] Booleans use is/has prefix (99%), should (rare)
[ ] No abbreviations except universally understood ones
[ ] Banned generic garbage names (`data`, `temp`, `obj`, `val`) in favor of domain entities
```

---

## Cross-References

- [Boolean Principles](./02-boolean-principles/readme.md) — P1–P6 boolean naming rules
- [Function Naming](./10-function-naming.md) — Function and method naming conventions
- [Key Naming PascalCase](./11-key-naming-pascalcase.md) — API/database key casing
- [Database Naming](./07-database-naming.md) — Table and column naming

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CROSS-022 / AC-CG-NAME-022: Cross-Language Variable and Collection Naming

**Given** Variable declarations, loop counters, collection identifiers, and map instances in PHP, TypeScript, or Go.
**When** Codebases are audited by coding guideline scanners and lint rules.
**Then** Singular/plural rules, loop singular derivations, `Map`/`By` suffixes, abbreviation bans, and strict prohibition of generic garbage names (`data`, `temp`, `obj`) are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.

---

*Cross-language variable naming conventions v1.0.0 — 2026-03-31*
