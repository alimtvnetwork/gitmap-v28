# Blank Lines & Spacing

> **Version:** 4.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go  
> **Rules covered:** 4, 5, 10

---

## Rule 4: Blank Line Before `return` or `throw` When Preceded by Other Statements

If a block contains statements before `return` or `throw`, insert **one blank line** before the `return`/`throw`. If `return`/`throw` is the **only statement** in the block, no blank line is needed.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: No blank line before return
if (ErrorChecker::isInvalidPdoExtension()) {
    $this->logger->error('PDO/SQLite not available');

    return $this->envelope->error('SQLite support not available', 500);
}

// ✅ REQUIRED: Blank line separates logic from exit
if (ErrorChecker::isInvalidPdoExtension()) {
    $this->logger->error('PDO/SQLite not available');

    return $this->envelope->error('SQLite support not available', 500);
}

// ❌ FORBIDDEN: No blank line before throw
if (PathHelper::isFileMissing($path)) {
    $this->logger->error('File not found: ' . $path);

    throw new RuntimeException('File not found: ' . $path);
}

// ✅ REQUIRED: Blank line before throw
if (PathHelper::isFileMissing($path)) {
    $this->logger->error('File not found: ' . $path);

    throw new RuntimeException('File not found: ' . $path);
}

// ✅ OK: Return is the only statement — no blank line needed
if ($error === null) {
    return false;
}

// ✅ OK: Throw is the only statement — no blank line needed
if ($error === null) {
    throw new InvalidArgumentException('Error required');
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
const processData = (data: unknown[]) => {
    const filtered = data.filter(isValid);

    return filtered.map(transform);
};

// ✅ REQUIRED
const processData = (data: unknown[]) => {
    const filtered = data.filter(isValid);

    return filtered.map(transform);
};

// ❌ FORBIDDEN: No blank line before throw
const validate = (input: string) => {
    const trimmed = input.trim();

    throw new Error(`Invalid input: ${trimmed}`);
};

// ✅ REQUIRED
const validate = (input: string) => {
    const trimmed = input.trim();

    throw new Error(`Invalid input: ${trimmed}`);
};

// ✅ OK: Return is the only statement
if (!data) {
    return null;
}
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
func process(data []Item) apperror.Result[[]Item] {
    filtered := filter(data)

    return apperror.Ok(filtered)
}

// ✅ REQUIRED
func process(data []Item) apperror.Result[[]Item] {
    filtered := filter(data)

    return apperror.Ok(filtered)
}
```

---

## Rule 5: Blank Line After Closing `}` When Followed by More Code

If code continues after a closing `}` (i.e., not followed by another `}`, `else`, `catch`, or end of function), insert **one blank line** after it. This applies to **all block types**: `if`, `foreach`/`for`/`for...of`, `while`, `switch`, `try`, and any other brace-delimited block.

### 5a — After `if` Blocks

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: No blank line after block when code follows
if ($this->initialized) {
    return;
}

$this->initialized = true;
add_action(HookType::Init->value, [$this, 'setup']);

// ✅ REQUIRED: Blank line after block when code follows
if ($this->initialized) {
    return;
}

$this->initialized = true;
add_action(HookType::Init->value, [$this, 'setup']);
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
if (!user) {
    return;
}

const profile = await fetchProfile(user.id);

// ✅ REQUIRED
if (!user) {
    return;
}

const profile = await fetchProfile(user.id);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
if err != nil {
    return err
}

result := compute()

// ✅ REQUIRED
if err != nil {
    return err
}

result := compute()
```

### 5b — After Loop Blocks (`foreach`, `for`, `while`, `for...of`)

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: No blank line after foreach when code follows
foreach (array_keys($data) as $col) {
    $setParts[] = "{$col} = ?";
}

$setClause = implode(', ', $setParts);
$sql       = "UPDATE {$table} SET {$setClause} WHERE {$where}";

// ✅ REQUIRED: Blank line separates the loop from subsequent logic
foreach (array_keys($data) as $col) {
    $setParts[] = "{$col} = ?";
}

$setClause = implode(', ', $setParts);
$sql       = "UPDATE {$table} SET {$setClause} WHERE {$where}";
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: No blank line after for...of when code follows
for (const item of items) {
    processed.push(transform(item));
}

const result = merge(processed);

// ✅ REQUIRED
for (const item of items) {
    processed.push(transform(item));
}

const result = merge(processed);

// ❌ FORBIDDEN: No blank line after while when code follows
while (queue.length > 0) {
    const task = queue.shift()!;
    execute(task);
}

logCompletion(queue);

// ✅ REQUIRED
while (queue.length > 0) {
    const task = queue.shift()!;
    execute(task);
}

logCompletion(queue);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: No blank line after for range when code follows
for _, item := range items {
    results = append(results, process(item))
}

total := len(results)

// ✅ REQUIRED
for _, item := range items {
    results = append(results, process(item))
}

total := len(results)

// ❌ FORBIDDEN: No blank line after for loop when code follows
for i := 0; i < retries; i++ {
    if err = attempt(ctx); err == nil {
        break
    }
}

logger.Info("retries exhausted", "attempts", retries)

// ✅ REQUIRED
for i := 0; i < retries; i++ {
    if err = attempt(ctx); err == nil {
        break
    }
}

logger.Info("retries exhausted", "attempts", retries)
```

### 5c — After `switch` / `try` Blocks

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
try {
    $result = $this->execute($sql);
} catch (Throwable $e) {
    $this->logger->error($e->getMessage());
}

$this->cleanup();

// ✅ REQUIRED
try {
    $result = $this->execute($sql);
} catch (Throwable $e) {
    $this->logger->error($e->getMessage());
}

$this->cleanup();
```

### Exception: Consecutive Closing Braces, `else`, `catch`, `finally`

No blank line is needed when a `}` is immediately followed by another `}`, `else`, `catch`, or `finally`:

```php

if (ErrorChecker::isFatalError($error)) {
    $this->logger->fatal($error);
}
// ✅ No blank line — next line is another closing brace
```

```go

if err != nil {
    return err
} // ✅ No blank line — function ends here (outer })
```

---

## Rule 10: Blank Line Before Control Structures When Preceded by Statements

When an `if`, `for`, `foreach`/`for...of`, or `while` block is preceded by **one or more non-brace statements** (assignments, function calls, etc.), insert **one blank line** before the control structure. This visually separates "setup" from "decision" logic.

**Exception:** No blank line is needed when the control structure is the first statement in a block or immediately follows another closing `}` (already covered by Rule 5).

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: No blank line between statement and if
$result = $this->apiRequest($agentId, HttpMethodType::Post->value, $endpoint);

if (is_wp_error($result)) {
    return $result;
}

// ✅ REQUIRED: Blank line before if when preceded by a statement
$result = $this->apiRequest($agentId, HttpMethodType::Post->value, $endpoint);

if (is_wp_error($result)) {
    return $result;
}

// ❌ FORBIDDEN: No blank line between statement and foreach
$items = $this->fetchItems();

foreach ($items as $item) {
    $this->process($item);
}

// ✅ REQUIRED
$items = $this->fetchItems();

foreach ($items as $item) {
    $this->process($item);
}

// ✅ OK: if is the first statement — no blank line needed
public function handle(): void {
    if ($this->isDone()) {
        return;
    }
}

// ✅ OK: if follows a closing brace — Rule 5 applies instead
if ($guardA) {
    return;
}

if ($guardB) {
    return;
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
const data = await fetchData(url);

if (!data) {
    return null;
}

// ✅ REQUIRED
const data = await fetchData(url);

if (!data) {
    return null;
}
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
result, err := doWork(ctx)

if err != nil {
    return err
}

// ✅ REQUIRED
result, err := doWork(ctx)

if err != nil {
    return err
}
```

---

*Part of [Code Style](./00-overview.md) — Rules 4, 5, 10*
