# `JsonResult` Multi-Source Creation & AUK Go CoreJSON Architecture

> **Document:** `research/11-jsonresult-multi-source-creation-and-aukgo-architecture.md`  
> **Status:** Implemented & Verified  
> **Package Reference:** `04-code/golang/pkg/streamwriter` & `03-aukgo/core/coredata/corejson`  
> **Date:** 2026-09-04  

---

## 1. Executive Summary & Research Context

Analysis of the reference codebase `03-aukgo/core/coredata/corejson` demonstrates a comprehensive, industrial-grade JSON data manipulation pipeline. In `corejson`, a JSON result is not merely a static serialized byte slice; it is an active monadic container capable of originating from diverse sources (memory, streams, interfaces, strings, errors, or generic reflection) and navigating round-trip casting.

This document abstracts the architectural principles of `corejson` into a generic blueprint designed for AI agents and Go engineers implementing high-performance, type-safe streaming systems.

---

## 2. Deep Dive: Key Patterns Discovered in `corejson`

### 2.1 The Struct-As-Namespace Pattern
Rather than exposing dozens of unorganized top-level package functions, `corejson` organizes capabilities under package-level singleton structs acting as namespaces:
- `corejson.Serialize.*`: High-level serialization engines (`ToString`, `Raw`, `UsingAny`, `Apply`).
- `corejson.Deserialize.*`: High-level deserialization engines (`UsingBytes`, `UsingString`, `Apply`, `FromTo`).
- `corejson.NewResult.*`: Specialized factory methods for building `Result` objects.
- `corejson.AnyTo.*`: Polymorphic type converters.
- `corejson.CastAny.*`: Arbitrary type casting via JSON round-trip serialization.

### 2.2 Multi-Source Ingestion Hierarchy
In `newResultCreator.go` and `serializerLogic.go`, a `Result` can be created from at least 10 distinct sources:
1. **Raw Byte Slices (`[]byte`):** From existing buffers or I/O reads (`UsingBytes`).
2. **Raw Strings (`string`):** From API payloads, logs, or CLI arguments (`UsingString`, `UsingTypePlusString`).
3. **Structured Objects (`any` / `T`):** Structs, maps, slices marshaled via standard encoder (`Any`, `Serialize`).
4. **Custom Serializer Interfaces (`bytesSerializer`):** Any object implementing `Serialize() ([]byte, error)`.
5. **Dynamic Serializer Closures (`func() ([]byte, error)`):** On-demand lazy execution (`UsingSerializerFunc`).
6. **Domain Model Interfaces (`Jsoner`):** Objects providing `Json() Result` or `JsonPtr() *Result`.
7. **Explicit Error Envelopes (`error` / `*appfault.AppError`):** Creating failed results where the error is preserved (`Error`, `ErrorPtr`).
8. **Pre-existing Result Envelopes:** Copying or re-wrapping results (`DeserializeUsingResult`).
9. **Polymorphic Universal Casting (`AnyTo`):** Dynamic runtime type switches inspecting `fromAny` and delegating to the appropriate constructor.
10. **I/O Streams (`io.Reader`):** Streaming incoming network/file data into a validated JSON envelope.

### 2.3 Safe Accessors vs. Error Handling
`corejson.Result` maintains three core fields:
```go
type Result struct {
    Bytes    []byte
    Error    error
    TypeName string
}
```
Methods are cleanly segregated by error semantics:
- **Non-panicking safe accessors:** `SafeBytes()`, `JsonString()`, `PrettyJsonString()`, `HasError()`, `HasIssuesOrEmpty()`.
- **Structured Error extraction:** `MeaningfulError()`, `ErrorString()`.
- **Enforced safety variants:** `MustBeSafe()`, `RawMust()`, `RawStringMust()`.

---

## 3. Modern Generic Architecture: Upgrades for `streamwriter`

While `corejson` relies on Go 1.18 reflection and standard library `error`, our `streamwriter` modernization introduces three major enhancements:

1. **Non-Generic `JsonResult` Envelope (Without `T`):** Eliminates unnecessary generic type parameter pollution across call sites. Ingestion envelopes should represent dynamic JSON payloads without forcing callers to specify type parameters (`JsonResult[any]`). Callers can extract types safely via `.Unmarshal(&dest)` or `UnmarshalAs[Target](res)`.
2. **Universal Error Wrapper `*appfault.AppError`:** Total elimination of untyped Go `error`. All failures carry classification (`errtype.Variation`), stack traces, and status codes.
3. **`sync.Locker` & Status Flag Integration:** Every `JsonResult` carries an explicit `status bool` and `statusCode int`, conforming directly to `WrappedBytes[any]` and `WrappedJson`.

---

## 4. Multi-Source Factory Design: Why `JsonResult` is Without `T`

### 4.1 Rationale for Eliminating `T` from `JsonResult`

In production workflows, callers ingesting raw JSON from external sources (`[]byte`, `string`, `io.Reader`) do not yet know or have an instance of the structured type before parsing and validation. Forcing a generic signature like `JsonResult[T]` requires every consumer to either propagate `T` everywhere or deal with clumsy `JsonResult[any]` declarations.

By making `JsonResult` a non-generic struct (`type JsonResult struct` with `payload any`), the design achieves maximum ergonomics:

```go
type JsonResult struct {
    data       []byte
    payload    any
    status     bool
    statusCode int
    appError   *appfault.AppError
}
```

1. The global namespace singleton `JsonSource` returns clean `JsonResult`:
   - `JsonSource.FromBytes(data []byte, payload ...any) JsonResult`
   - `JsonSource.FromString(str string, payload ...any) JsonResult`
   - `JsonSource.FromReader(r io.Reader, payload ...any) JsonResult`
   - `JsonSource.FromPayload(payload any) JsonResult`
   - `JsonSource.FromSerializer(fn, payload ...any) JsonResult`
   - `JsonSource.FromBytesEnvelope(wb any) JsonResult`
   - `JsonSource.FromError(appErr *appfault.AppError) JsonResult`
   - `JsonSource.FromAny(source any) JsonResult`
   - `JsonSource.Cast(source any, targetPtr any) *appfault.AppError`
2. Strongly-typed extraction is provided cleanly through:
   - **Method Unmarshal:** `res.Unmarshal(&dest)` returning `*appfault.AppError`.
   - **Generic Helper Unmarshal:** `dest, err := streamwriter.UnmarshalAs[Target](res)`.
   - **Type-Casting Helpers:** `streamwriter.Cast[Target](source)` returning `JsonResult`.
   - **Scoped Factory:** `JsonSourceOf[T]()` for workflows that want to attach typed payload `T` to the resulting `JsonResult`.

```go
// 1. Untyped / Dynamic Ingestion
res := streamwriter.JsonSource.FromBytes(rawBytes)

// 2. Direct Pointer Casting
var dest Account
err := streamwriter.JsonSource.Cast(rawBytes, &dest)

// 3. Typed Extraction
account, appErr := streamwriter.UnmarshalAs[Account](res)

// 4. Standalone Helpers
resFromPayload := streamwriter.FromPayload(myAccount)
resCast := streamwriter.CastTo[Account](rawBytes)
```

---

## 5. Standard AI Guidelines for Implementing Multi-Source Results

When any AI agent implements or extends JSON result containers across this meta-repository, they MUST adhere to these architectural standards:

1. **Non-Generic Container:** `JsonResult` and `WrappedJson` MUST be non-generic (without `T`). Internal payload is stored as `any` and satisfies `WrappedBytes[any]`.
2. **Use `Json` Naming Convention:** All types and identifiers MUST use camel/pascal casing `Json` (e.g. `JsonResult`, `WrappedJson`, `JsonSource`) rather than all-caps `JSON`. Backwards-compatible aliases (`JSONResult`, `WrappedJSON`, `JSONSource`) may be maintained for transition periods.
3. **Never Panic in Constructors:** If input bytes or strings are malformed JSON, return a valid `JsonResult` with `status: false`, `statusCode: 400`, and an attached `*appfault.AppError`.
4. **Preserve Payload Along with Bytes:** Maintain the `payload any` inside the result container so callers can inspect structured fields without redundant unmarshaling.
5. **Unified Interface Conformance:** Every result container MUST satisfy `WrappedBytes[any]` and `WrappedJson`, exposing `Raw()`, `String()`, `Len()`, `Value()`, `Error()`, `Status()`, `StatusCode()`, `IsSuccess()`, and `Unwrap()`.
6. **Deterministic Formatting:** Indented formatting (`Pretty()`) and minification (`Compact()`) must produce stable outputs.
7. **No Mixed Polarity:** Boolean evaluation methods (`IsSuccess()`, `HasError()`, `IsValid()`) must never evaluate `isA && !isB` in the same condition.
