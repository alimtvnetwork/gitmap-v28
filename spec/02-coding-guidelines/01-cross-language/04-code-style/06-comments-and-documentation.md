# Comments, Documentation & Dead Code

> **Version:** 4.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go  
> **Rules covered:** 8, 14, 15, 16

---

## Rule 8: No Leading Backslash on Global Types

In catch blocks and type hints, use `Throwable` without the leading backslash, even in namespaced files. The same applies to other global types used in catch blocks or parameter hints.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
catch (\Throwable $e)
function foo(\Throwable $e): array

// ✅ REQUIRED
catch (Throwable $e)
function foo(Throwable $e): array
```

```typescript
// ── TypeScript / Go ─────────────────────────────────────────
// Not applicable — these languages don't have leading-backslash syntax.
```

---

## Rule 14: No Commented-Out or Dead Code

Commented-out code and unreachable (dead) code **must be deleted**, not left in the codebase. Version control preserves history — comments are not an archive tool.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Commented-out code
// $oldValue = $this->legacyLookup($key);
// if ($oldValue !== null) {
//     return $oldValue;
// }
$value = $this->lookup($key);

// ❌ FORBIDDEN: Dead code after unconditional return
return $result;
$this->cleanup();   // never executes

// ✅ REQUIRED: Remove dead code entirely
$value = $this->lookup($key);
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Commented-out code
// const legacyResult = await fetchLegacy(id);
// return legacyResult;
const result = await fetchData(id);

// ✅ REQUIRED: Clean — no commented-out code
const result = await fetchData(id);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Commented-out code
// oldVal, err := legacyFetch(ctx, key)
// if err != nil {
//     return err
// }
val, err := fetch(ctx, key)

// ✅ REQUIRED: Clean
val, err := fetch(ctx, key)
```

### Exceptions

- **TODO/FIXME comments** with a ticket reference are allowed: `// TODO(PROJ-123): migrate to new API`
- **Intentional stubs** for future implementation are allowed when marked: `// STUB: placeholder for upcoming feature PROJ-456`

---

## Rule 15: Comment Formatting — Space After `//`

Every line comment **must** have a space between the comment marker and the text. Block comments must also have consistent spacing.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
//this is a comment
//$value = 5;

// ✅ REQUIRED
// This is a comment
// $value = 5;
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
//calculate total
const total = items.reduce((sum, i) => sum + i.price, 0);

// ✅ REQUIRED
// Calculate total
const total = items.reduce((sum, i) => sum + i.price, 0);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
//fetchData retrieves data from the store
func fetchData(ctx context.Context) error {

// ✅ REQUIRED
// fetchData retrieves data from the store.
func fetchData(ctx context.Context) error {
```

---

## Rule 16: Method and Function Documentation (When To Write, When Not To)

Core rule: simple methods do NOT require documentation. Do not write verbose comments. Comments lie, code does not. The name and the signature are the primary documentation. If you feel the need to explain what a method does in prose, first try to rename it or split it until the code explains itself.

Method documentation is required ONLY when at least one of these is true, and even then the preferred fix is to refactor so the doc is no longer needed:

1. The method is doing multiple non-obvious things that could not be expressed in the name (which itself is a smell, refactor first).
2. The method processes or transforms data in a way where a short example clarifies the contract. Example: `path.Clean` performs path cleaning and normalization, so a one-line example is worth more than prose.
3. The code is adapted or copied from another source. Citation is mandatory (link plus license note).
4. The team runs an automated doc generator (godoc, TypeDoc, phpDocumentor) against the codebase. In that case public/exported APIs get one-line doc comments so the generated docs are usable.

Reference for the Go convention (doc comment starts with the identifier name, no blank line between doc and declaration): https://go.dev/src/go/doc/example.go

```go
// ── Go (canonical example) ───────────────────────────────────

// ❌ AVOID: verbose prose that repeats the code
// GetUser gets a user by id. It takes an id and returns the user
// object from the database, and if there is an error it returns
// the error.
func GetUser(id int64) (User, error) { ... }

// ❌ AVOID: doc on a trivially named simple method
// Add adds a and b and returns the sum.
func Add(a, b int) int { return a + b }

// ✅ OK: exported, non-trivial behavior, includes a usage example.
// Clean returns the shortest path name equivalent to path by purely
// lexical processing. It applies the following rules iteratively
// until no further processing can be done:
//   1. Replace multiple slashes with a single slash.
//   2. Eliminate each . path name element.
//   3. Eliminate each inner .. path name element.
//   4. Eliminate .. elements that begin a rooted path.
func Clean(path string) string { ... }
```

The same principle applies to every language: TypeScript, PHP, Rust, C#, PowerShell, Python. Only the syntax of the comment changes.

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ AVOID: doc on a self-explanatory function
/** Returns the total price. */
export const getTotalPrice = (items: Item[]): number =>
  items.reduce((sum, i) => sum + i.price, 0);

// ✅ OK: non-obvious behavior worth documenting once
/**
 * Applies the tier discount then the promo code, in that order.
 * Order matters: promo codes stack on the already-discounted price.
 */
export const applyPricing = (price: number, tier: Tier, promo?: Promo): number => { ... }
```

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ AVOID: restating the signature
/** Returns the user by id. */
public function getUser(int $id): User { ... }

// ✅ OK: adapted code with a citation
/**
 * Levenshtein distance with an early-exit threshold.
 * Adapted from https://en.wikipedia.org/wiki/Levenshtein_distance (CC-BY-SA).
 */
public function distance(string $a, string $b, int $maxDistance): int { ... }
```

### Decision Checklist (apply before writing any doc comment)

1. Can I rename the method so the doc becomes redundant? If yes, rename and delete the doc.
2. Can I split the method so each piece is trivially named? If yes, split and delete the doc.
3. Does the doc restate the signature or parameter names? If yes, delete it.
4. Is the doc explaining WHY (business rule, ordering constraint, cited source) or a non-obvious example? If yes, keep it, short.
5. Does the team run automated doc generation? If yes, one-liner on exported APIs is fine.

Trivial getters, setters, single-expression helpers, and any function whose name fully describes its behavior get no doc comment. Adding one is a review-blocking violation of "comments lie, code does not".

---

*Part of [Code Style](./00-overview.md) — Rules 8, 14, 15, 16*
