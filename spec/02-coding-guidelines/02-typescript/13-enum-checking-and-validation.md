# 13. Enum Checking and Validation Guidelines

> **Version:** 2.0.0 (V2)  
> **Applies to:** TypeScript Codebase  
> **Priority:** CRITICAL

This document details the **zero-tolerance policy for magic strings** and provides explicit examples for defining, checking, and validating Enums.

---

## 1. Zero Magic Strings Policy

Raw strings must **never** be used for domain logic, event names, statuses, or arbitrary flags. Any string that represents a state, type, or identifier with a bounded set of values MUST be extracted into a `const` or `enum`.

### ❌ FORBIDDEN: Raw Magic Strings

```typescript
// ❌ WRONG: Hardcoded strings scattered in code
if (status === 'success') {
  console.log('info', 'Operation succeeded');
}

// ❌ WRONG: Raw string types in interfaces
interface ButtonProps {
  variant: 'primary' | 'secondary' | 'danger';
}
```

### ✅ REQUIRED: Strict Constant Maps / Enums

```typescript
// ✅ RIGHT: Centralized constants
export const Status = {
  Success: 'success',
  Error: 'error',
} as const;

export type StatusType = typeof Status[keyof typeof Status];

if (status === Status.Success) {
  Logger.log(LogLevel.Info, 'Operation succeeded');
}
```

---

## 2. Strict Enum Validation (Type Guards)

When receiving a string from an API, user input, or legacy data, you MUST validate it against the Enum values. **Do not use `any` or cast blindly.**

### ❌ FORBIDDEN: Blind Casting

```typescript
// ❌ WRONG: Casting raw string to Enum type blindly
const backendStatus = response.status as StatusType; // Dangerous! What if response.status is 'foo'?
```

### ✅ REQUIRED: Safe Type Guard Validation

You must implement a type guard to validate if an arbitrary string belongs to the Enum.

```typescript
// ✅ RIGHT: Type Guard Function
export const Status = {
  Success: 'success',
  Error: 'error',
  Pending: 'pending',
} as const;
export type StatusType = typeof Status[keyof typeof Status];

// Type Guard
export function isValidStatus(value: unknown): value is StatusType {
  return Object.values(Status).includes(value as StatusType);
}

// Usage Example
if (isValidStatus(response.status)) {
  // TypeScript now strictly understands response.status is StatusType
  handleStatus(response.status); 
} else {
  // Fallback or error logging
  Logger.error(LogLevel.Error, `Invalid status received: ${response.status}`);
}
```

---

## 3. Iterating Over Enums

When you need to render UI elements (like dropdowns) based on an Enum, you should iterate over the object values.

### ✅ REQUIRED: Iterating with `Object.values`

```typescript
import { LogLevel } from '@/constants/log-levels';

function LogLevelSelector() {
  return (
    <select>
      {Object.values(LogLevel).map((level) => (
        <option key={level} value={level}>
          {level.toUpperCase()}
        </option>
      ))}
    </select>
  );
}
```

---

## 4. UI Events & React Properties

React `onChange`, `onClick`, and other native event listeners should map internal strings to domain-specific constants if they are passed through generic handlers.

### ❌ FORBIDDEN: Magic String Event Handlers

```typescript
// ❌ WRONG
function handleEvent(eventName: string) {
  if (eventName === 'mouseover') { ... }
}
```

### ✅ REQUIRED: Using Event Constants

```typescript
// ✅ RIGHT
export const UIEvents = {
  MouseOver: 'mouseover',
  Click: 'click',
  Change: 'change'
} as const;

function handleEvent(eventName: typeof UIEvents[keyof typeof UIEvents]) {
  if (eventName === UIEvents.MouseOver) { ... }
}
```
