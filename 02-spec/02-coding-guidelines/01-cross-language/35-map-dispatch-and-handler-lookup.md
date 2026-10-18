# Map Dispatch & Handler Table Patterns — Coding Guidelines

## 1. Overview & Core Mission
In event-driven, CLI, UI, and protocol routing codebases, developers frequently construct long sequences of chained conditional branches (`if ... else if ... else if`) or verbose `switch-case` statements to inspect input keys, tokens, action names, or status codes and invoke corresponding handler functions.

This specification establishes the **Map Dispatch & Handler Table Pattern** as a strict coding standard across all supported languages (TypeScript, Go, Python, Rust, PHP). Repetitive conditional branching matching literal strings, tokens, or event codes MUST be replaced with declarative, immutable lookup tables that execute discrete handler functions in a single unified line.

---

## 2. Anti-Pattern: Chained Branching Dispatch Table

### ❌ Anti-Pattern: Chained If-Else Keyboard / Event Handlers
```typescript
// ❌ ANTI-PATTERN: High cyclomatic complexity, repetitive boilerplate, hard to extend
useEffect(() => {
  if (!isOpen) return;

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
    } else if (e.key === '+' || e.key === '=') {
      e.preventDefault();
      handleZoomIn();
    } else if (e.key === '-') {
      e.preventDefault();
      handleZoomOut();
    } else if (e.key === '0') {
      e.preventDefault();
      handleResetZoom();
    }
  };

  window.addEventListener('keydown', handleKeyDown);
  return () => window.removeEventListener('keydown', handleKeyDown);
}, [isOpen]);
```

### ❌ Why This Code Smell Fails:
1. **Excessive Branching Complexity:** Every new shortcut or key requires adding another `else if` branch, driving cyclomatic complexity higher and inflating function lines past the 8–15 line limit.
2. **Repetitive Boilerplate:** Statements like `e.preventDefault()` are repeated in every branch instead of being handled centrally upon matching.
3. **Fragile Extensibility:** Adding, removing, or testing keyboard bindings requires modifying the nested logic flow rather than editing a declarative data structure.
4. **Poor Testability:** Discrete keyboard behaviors cannot be unit-tested without mocking synthetic DOM keyboard events and triggering the entire component lifecycle.

---

## 3. The Canonical Pattern: Declarative Map Dispatch

### ✅ Refactored Solution: TypeScript / React
```typescript
// ✅ CANONICAL PATTERN: Declarative handler map with unified dispatch
type KeyHandler = () => void;

const KEY_ACTIONS: Record<string, KeyHandler> = {
  Escape: onClose,
  '+': handleZoomIn,
  '=': handleZoomIn,
  '-': handleZoomOut,
  '0': handleResetZoom,
};

const handleKeyDown = (e: KeyboardEvent) => {
  const action = KEY_ACTIONS[e.key];
  if (action) {
    e.preventDefault();
    action();
  }
};
```

### ✅ Clean Component Isolation Pattern:
```typescript
// ✅ In React components where handlers close over state/props:
const getKeyActionMap = (actions: {
  onClose: () => void;
  onZoomIn: () => void;
  onZoomOut: () => void;
  onResetZoom: () => void;
}): Record<string, () => void> => ({
  Escape: actions.onClose,
  '+': actions.onZoomIn,
  '=': actions.onZoomIn,
  '-': actions.onZoomOut,
  '0': actions.onResetZoom,
});

useEffect(() => {
  if (!isOpen) return;

  const keyMap = getKeyActionMap({
    onClose,
    onZoomIn: handleZoomIn,
    onZoomOut: handleZoomOut,
    onResetZoom: handleResetZoom,
  });

  const handleKeyDown = (e: KeyboardEvent) => {
    const handler = keyMap[e.key];
    if (handler) {
      e.preventDefault();
      handler();
    }
  };

  window.addEventListener('keydown', handleKeyDown);
  return () => window.removeEventListener('keydown', handleKeyDown);
}, [isOpen, onClose, handleZoomIn, handleZoomOut, handleResetZoom]);
```

---

## 4. Polyglot Implementations

### 4.1 Go (Golang)
In Go, replace switch or if-chains with a package-level or local map of handler functions returning `*appfault.AppError`:

```go
// ❌ BAD: Chained branching command router
func dispatchCommand(cmd string, ctx *CommandContext) *appfault.AppError {
    if cmd == "init" {
        return handleInit(ctx)
    } else if cmd == "status" {
        return handleStatus(ctx)
    } else if cmd == "sync" {
        return handleSync(ctx)
    }
    return appfault.NewNotFoundError("command not recognized: " + cmd)
}

// ✅ GOOD: Map dispatch table
type CommandHandler func(*CommandContext) *appfault.AppError

var commandRegistry = map[string]CommandHandler{
    "init":   handleInit,
    "status": handleStatus,
    "sync":   handleSync,
}

func dispatchCommand(cmd string, ctx *CommandContext) *appfault.AppError {
    handler, isFound := commandRegistry[cmd]
    if !isFound {
        return appfault.NewNotFoundError("command not recognized: " + cmd)
    }

    return handler(ctx)
}
```

### 4.2 Python
In Python, declare a typed dictionary mapping command/event names to callable functions:

```python
# ❌ BAD: Chained if-elif ladder
def route_event(event_type: str, payload: dict) -> None:
    if event_type == "user_created":
        handle_user_created(payload)
    elif event_type == "user_updated":
        handle_user_updated(payload)
    elif event_type == "user_deleted":
        handle_user_deleted(payload)

# ✅ GOOD: Typed action map dispatch
from collections.abc import Callable

EVENT_DISPATCH: dict[str, Callable[[dict], None]] = {
    "user_created": handle_user_created,
    "user_updated": handle_user_updated,
    "user_deleted": handle_user_deleted,
}

def route_event(event_type: str, payload: dict) -> None:
    handler = EVENT_DISPATCH.get(event_type)
    if handler:
        handler(payload)
```

### 4.3 Rust
In Rust, use compile-time static slices or `phf` / `HashMap` for zero-allocation dispatch:

```rust
// ✅ GOOD: Function pointer table in Rust
type ActionHandler = fn(&mut Context);

const ACTION_MAP: &[(&str, ActionHandler)] = &[
    ("Escape", handle_close),
    ("+", handle_zoom_in),
    ("=", handle_zoom_in),
    ("-", handle_zoom_out),
    ("0", handle_reset_zoom),
];

pub fn handle_key(key: &str, ctx: &mut Context) {
    if let Some((_, handler)) = ACTION_MAP.iter().find(|(k, _)| *k == key) {
        handler(ctx);
    }
}
```

### 4.4 PHP
In PHP 8+, use an associative array mapping strings to callables:

```php
// ✅ GOOD: Associative array handler dispatch
$handlers = [
    'create' => fn(array $data) => $this->handleCreate($data),
    'update' => fn(array $data) => $this->handleUpdate($data),
    'delete' => fn(array $data) => $this->handleDelete($data),
];

$handler = $handlers[$action] ?? null;
if ($handler) {
    $handler($payload);
}
```

---

## 5. Non-Negotiable Coding Rules

1. **Threshold Rule:** Any chained `if ... else if ... else if` sequence or `switch` block evaluating 3 or more literal string/token values to invoke corresponding functions MUST be refactored into a map dispatch table.
2. **Immutable Maps:** Wherever possible, declare handler maps as constants (`const KEY_ACTIONS`, package-level `var`, or module-level frozen dicts) outside hot loops or render cycles to avoid per-invocation allocation.
3. **Single-Line Lookup & Dispatch:** The lookup and execution site must be compact and affirmative:
   ```typescript
   const action = MAP[key];
   if (action) {
     action();
   }
   ```
4. **Shared Pre/Post Hooks:** Common logic (e.g. `e.preventDefault()`, telemetry logging, error catching) must wrap the dispatch call ONCE rather than being duplicated in every handler branch.
5. **Positive Booleans Only:** In presence testing, never use explicit true checks (`if isFound == true` is FORBIDDEN); use implicit evaluation (`if isFound { ... }` or `if (action) { ... }`).
6. **Strict Function Sizing:** Discrete handlers mapped inside the table must be small, bounded functions (<= 8–15 lines).
7. **Strict Relative Paths:** All documentation, specifications, and comments referencing this rule must strictly use relative paths from the git repository root.
