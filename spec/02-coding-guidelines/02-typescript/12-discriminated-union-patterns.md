# Discriminated Union & Action Type Patterns

> **Parent:** [TypeScript Standards](./00-overview.md)  
> **Version:** 1.0.0  
> **Updated:** 2026-04-05

---

## Rule: No Inline Types in Discriminated Unions

Discriminated unions must use **named types** for each variant — never inline object shapes. Each variant gets its own extracted interface. The discriminant field must reference a **PascalCase enum**, not string literals or `UPPER_SNAKE_CASE` constants.

---

## 1. Enum Values — PascalCase Only

Enum values must use **PascalCase**. `UPPER_SNAKE_CASE`, `camelCase`, and `_camelCase` are strictly prohibited.

```typescript
// ❌ PROHIBITED — UPPER_SNAKE_CASE values
enum ActionType {
  ADD_TOAST = "ADD_TOAST",
  UPDATE_TOAST = "UPDATE_TOAST",
  DISMISS_TOAST = "DISMISS_TOAST",
  REMOVE_TOAST = "REMOVE_TOAST",
}

// ✅ REQUIRED — PascalCase values
enum ActionType {
  AddToast = "AddToast",
  UpdateToast = "UpdateToast",
  DismissToast = "DismissToast",
  RemoveToast = "RemoveToast",
}
```

---

## 2. No Inline Object Types in Unions

Every variant in a discriminated union must be a **named interface**. Inline `{ type: ...; payload: ... }` blocks are prohibited — they cannot be reused, imported, or documented independently.

```typescript
// ❌ PROHIBITED — inline types, index access, complex nesting
type Action =
  | {
      type: ActionType["ADD_TOAST"];
      toast: ToasterToast;
    }
  | {
      type: ActionType["UPDATE_TOAST"];
      toast: Partial<ToasterToast>;
    }
  | {
      type: ActionType["DISMISS_TOAST"];
      toastId?: ToasterToast["id"];
    }
  | {
      type: ActionType["REMOVE_TOAST"];
      toastId?: ToasterToast["id"];
    };
```

```typescript
// ✅ REQUIRED — named interfaces, enum reference, reusable
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

## 3. Generic Action Pattern (Best Practice)

When many actions share the same shape, use a **generic action interface** to eliminate repetition:

```typescript
// ✅ BEST PRACTICE — generic discriminated action
interface TypedAction<T extends ActionType, P = void> {
  type: T;
  payload: P;
}

type ToastAction =
  | TypedAction<ActionType.AddToast, ToasterToast>
  | TypedAction<ActionType.UpdateToast, Partial<ToasterToast>>
  | TypedAction<ActionType.DismissToast, string | undefined>
  | TypedAction<ActionType.RemoveToast, string | undefined>;
```

### Reducer usage with exhaustive switch

```typescript
function toastReducer(state: ToastState, action: ToastAction): ToastState {
  switch (action.type) {
    case ActionType.AddToast:
      return { ...state, toasts: [action.payload, ...state.toasts] };

    case ActionType.UpdateToast:
      return {
        ...state,
        toasts: state.toasts.map((t) =>
          t.id === action.payload.id ? { ...t, ...action.payload } : t
        ),
      };

    case ActionType.DismissToast:
      return {
        ...state,
        toasts: state.toasts.map((t) =>
          t.id === action.payload ? { ...t, open: false } : t
        ),
      };

    case ActionType.RemoveToast:
      return {
        ...state,
        toasts: state.toasts.filter((t) => t.id !== action.payload),
      };
  }
}
```

---

## 4. Index Access on Enums — Prohibited

Never use bracket notation (`ActionType["AddToast"]`) to reference enum members. Use **dot notation** exclusively.

```typescript
// ❌ PROHIBITED — index access
type: ActionType["AddToast"]

// ✅ REQUIRED — dot notation
type: ActionType.AddToast
```

---

## Summary

| Rule | Prohibited | Required |
|------|-----------|----------|
| Enum values | `ADD_TOAST`, `REMOVE_TOAST` | `AddToast`, `RemoveToast` |
| Union variants | Inline `{ type: ...; }` | Named interface per variant |
| Enum reference | `ActionType["AddToast"]` | `ActionType.AddToast` |
| Repeated shapes | Copy-paste per variant | Generic `TypedAction<T, P>` |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| TypeScript Standards (§8) | [08-typescript-standards-reference.md § Discriminated Unions](./08-typescript-standards-reference.md) |
| Enum Conventions | [00-overview.md](./00-overview.md) |
| Consolidated Review Guide (Type Safety) | [../consolidated-review-guide.md](../consolidated-review-guide.md) |
| Condensed Review Guide (Types — CODE RED) | [../consolidated-review-guide-condensed.md](../consolidated-review-guide-condensed.md) |
| AI Quick-Reference Checklist | [../06-ai-optimization/02-ai-quick-reference-checklist.md](../06-ai-optimization/02-ai-quick-reference-checklist.md) |
