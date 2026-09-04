# Transaction Log: Swappable Writer Methods & Functional Injection Research

> **Directory:** `05-changes-history/03-swappable-writer-methods-research/`  
> **Date:** 2026-09-03  
> **Topic:** Swappable Write Methods, Higher-Order Function Injection via Options, and Log-Agnostic Payloads  
> **Status:** Completed  

---

## 1. Context & Objectives

1. **User Feedback:**
   - Appreciated the fluent writer chaining and options concept.
   - Pointed out that the previous `BaseWriter` code was "fixated" (rigidly hardcoded write execution).
   - Requested that the `Write` method itself be swappable on the fly or injected via `Options.WriteMethod`.
   - Emphasized that the writer must support both log-based records (`LogRecord`) and non-log payloads (arbitrary data, raw bytes, domain events).
   - Referenced the AUK Go `core` architecture (`AllLogWriter`, `LogDefinerWriter`) as a design influence.
   - Requested 3-4 distinct architectural patterns to be presented for review and selection.

---

## 2. Files Created & Modified

### Research Documents
- `research/01-index.md`: Updated with research topic 03.
- `research/03-swappable-writer-methods-and-functional-injection.md`: Comprehensive breakdown of 4 distinct swappable write method patterns.

### Transaction History Updates
- `05-changes-history/01-index.md`: Registered Task 03.
- `05-changes-history/03-swappable-writer-methods-research/01-transaction-log.md`: This file.

---

## 3. Four Architectural Patterns Explored

1. **Pattern 1: Functional Delegate Injection (`WriteFunc` in Options):**
   - The writer struct holds a `writeMethod WriteFunc` pointer.
   - `Write(ctx, payload)` delegates to `w.writeMethod(ctx, payload)`.
   - Users can supply a custom function at initialization via `Options.WriteMethod`, or hot-swap it dynamically via `SetWriteMethod()`.
2. **Pattern 2: Composable Pipeline (`FormatFunc` + `EmitFunc`):**
   - Deconstructs the writing process into two swappable stages: formatting (payload -> bytes) and emission (bytes -> destination).
   - Allows changing the output format without changing the destination, and vice versa.
3. **Pattern 3: Polymorphic Method Slot Overrides (AUK Go Style):**
   - Mirrors AUK Go's `LogDefinerWriter` design.
   - Provides dedicated slots for `WriteLog`, `WriteRaw`, and `WriteAny`, each independently swappable via `Options`.
4. **Pattern 4: Generic Dual-Channel Streamer (`Writer[T]`):**
   - Employs Go 1.18+ generics for compile-time type safety with zero interface boxing.

---

## 4. Verification & Status

- All documents verified against repo rules (strict lowercase filenames, relative git paths, no em dashes).
- Committed and pushed to both `coding-guidelines` and `gitmap`.
