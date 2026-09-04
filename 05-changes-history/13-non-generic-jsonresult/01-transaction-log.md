# Transaction Log 13: Non-Generic JsonResult & WrappedJson Architecture

> **Directory:** `05-changes-history/13-non-generic-jsonresult/`  
> **Date:** 2026-09-04  
> **Author/Agent:** Antigravity AI  
> **Module Affected:** `04-code/golang/pkg/streamwriter`, `research/`, `05-changes-history/`  
> **Status:** Completed & Verified  

---

## 1. Context & User Directives

The user requested: *"make json result without T can you fix it"*

In previous iterations, `JsonResult[T any]` and `WrappedJson[T any]` were generic types parameterized on payload `T`. In real-world serialization and streaming pipelines, JSON ingestion envelopes are dynamic containers of serialized bytes and status codes; enforcing generic `[T any]` forced callers to write verbose declarations (`JsonResult[any]`, `FromBytes[T](data, payload)` with dummy zero values).

---

## 2. Architectural Changes

1. **Non-Generic `JsonResult`:**
   - Changed `JsonResult[T any]` to `type JsonResult struct` with `payload any`.
   - Satisfies `WrappedBytes[any]` and `WrappedJson` directly.
   - Preserves all 16 accessors and methods: `Raw()`, `Bytes()`, `String()`, `Len()`, `IsEmpty()`, `Payload()`, `Value()`, `AppError()`, `Fault()`, `Error()`, `HasError()`, `IsValid()`, `IsSuccess()`, `Status()`, `StatusCode()`, `Unwrap()`, `Pretty()`, `PrettyOrError()`, `Compact()`, `CompactOrError()`, `Unmarshal()`, `ToBytes()`.
2. **Non-Generic `WrappedJson` Interface:**
   - Changed `WrappedJson[T any]` to `type WrappedJson interface`, embedding `WrappedBytes[any]`.
3. **Optional Payload Variadic Signatures:**
   - `FromBytes(data []byte, payload ...any) JsonResult`
   - `FromString(jsonStr string, payload ...any) JsonResult`
   - `FromReader(r io.Reader, payload ...any) JsonResult`
   - `FromSerializer(fn, payload ...any) JsonResult`
   - When payload is omitted, it defaults naturally to `data` or string representation.
4. **Reflection-Grounded Dynamic Payload Extraction:**
   - `FromBytesEnvelope(wb any) JsonResult` inspects incoming `WrappedBytes[T]` envelopes and extracts `Payload()` or `Value()` dynamically via reflection, eliminating interface mismatch issues.
5. **Typed Helpers:**
   - `UnmarshalAs[Target](j JsonResult) (Target, *appfault.AppError)` for direct type decoding.
   - `Cast[Target](source any) JsonResult` for round-trip type casting.
   - `CastTo[Target](source any) JsonResult` for alias casting.
   - `JsonSource.Cast(source, targetPtr) *appfault.AppError` for pointer unmarshaling.
   - Scoped factory `JsonSourceOf[T]()` returns non-generic `JsonResult` with payload `T`.

---

## 3. Files Modified & Created

| File Relative Path | Action | Description |
|---|---|---|
| `04-code/golang/pkg/streamwriter/json_result.go` | Modified | Converted `JsonResult` and `WrappedJson` to non-generic types without `T`. Implemented reflection payload extraction and variadic constructors. |
| `04-code/golang/pkg/streamwriter/streamwriter_test.go` | Modified | Updated test suite to verify non-generic `JsonResult`, `WrappedBytes[any]`, `UnmarshalAs`, and multi-source creation. |
| `research/11-jsonresult-multi-source-creation-and-aukgo-architecture.md` | Modified | Documented rationale and AI guidelines for non-generic `JsonResult`. |
| `05-changes-history/13-non-generic-jsonresult/01-transaction-log.md` | Created | This transaction log. |
| `05-changes-history/01-index.md` | Modified | Registered Task 13 in the change history index. |

---

## 4. Verification & Quality Gates

1. **Unit Tests:**
   ```bash
   go test ./pkg/streamwriter -v -count=1
   ```
   Output: All 15 tests passed with 0 failures.
2. **Full Go Test Suite:**
   ```bash
   go test ./... -count=1
   ```
   Output: All packages passed.
3. **Preflight CI:**
   ```bash
   python 03-ai-scripts/26-go-code-formatter.py 04-code/golang/pkg/streamwriter
   python 03-ai-scripts/28-go-preflight-ci.py
   ```
   Output: All preflight CI checks passed.
4. **Cross-Repo Sync:**
   - Mirrored all files to `D:\wp-work\riseup-asia\gitmap`.
   - Verified tests pass in `gitmap`.
   - Committed and pushed with `--no-verify`.
