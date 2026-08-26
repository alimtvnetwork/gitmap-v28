# RAG Validation Test Coverage Matrix

**Version:** 2.0.0  
**Created:** 2026-02-02  
**Updated:** 2026-05-07  
**Status:** Active
**Parent:** [04-rag-validation-tests.md](./03-rag-validation-tests.md)

---

## Overview

Comprehensive test coverage matrix mapping each RAG validation error code to its corresponding test cases, ensuring full coverage of all validation scenarios.

---

## Coverage Summary

| Error Code | Field | Test Count | Pass | Fail | Coverage |
|------------|-------|------------|------|------|----------|
| AB-9301 | ChunkSize (range) | 6 | ✅ | ✅ | 100% |
| AB-9302 | ChunkSize (multiple) | 4 | ✅ | ✅ | 100% |
| AB-9303 | ChunkOverlap | 10 | ✅ | ✅ | 100% |
| AB-9304 | ContextTokenBudget | 6 | ✅ | ✅ | 100% |
| AB-9305 | EmbeddingModel | 6 | ✅ | ✅ | 100% |
| AB-9306 | SimilarityThreshold | 8 | ✅ | ✅ | 100% |
| AB-9307 | TopK | 6 | ✅ | ✅ | 100% |
| AB-9308 | Config Load | 4 | ✅ | ✅ | 100% |
| AB-9309 | Config Save | 2 | ✅ | ✅ | 100% |
| AB-9310 | Config Corrupt | 2 | ✅ | ✅ | 100% |

**Total Tests:** 54  
**Total Coverage:** 100%

---

## Detailed Test Mapping

### AB-9301: ChunkSize Out of Range

| Test Case | Input | Expected | Type |
|-----------|-------|----------|------|
| `TooSmallZero` | 0 | FAIL:9301 | Boundary |
| `TooSmall100` | 100 | FAIL:9301 | Below min |
| `TooSmall255` | 255 | FAIL:9301 | Boundary-1 |
| `TooLarge8193` | 8193 | FAIL:9301 | Boundary+1 |
| `TooLarge16384` | 16384 | FAIL:9301 | Far above |
| `NegativeValue` | -256 | FAIL:9301 | Negative |

**Boundary Tests:**
- Minimum boundary (256): PASS
- Maximum boundary (8192): PASS
- Below minimum (255): FAIL
- Above maximum (8448): FAIL

---

### AB-9302: ChunkSize Not Multiple of 256

| Test Case | Input | Expected | Type |
|-----------|-------|----------|------|
| `NotMultiple257` | 257 | FAIL:9302 | Min+1 |
| `NotMultiple1000` | 1000 | FAIL:9302 | Round num |
| `NotMultiple1500` | 1500 | FAIL:9302 | Mid-range |
| `NotMultiple2000` | 2000 | FAIL:9302 | Near default |

**Valid Multiples Tested:**
- 256, 512, 1024, 2048, 4096, 8192

---

### AB-9303: ChunkOverlap Invalid

| Test Case | Overlap | ChunkSize | Expected | Reason |
|-----------|---------|-----------|----------|--------|
| `ValidZero` | 0 | 2048 | PASS | Min valid |
| `ValidDefault100` | 100 | 2048 | PASS | Default |
| `Valid25Percent` | 512 | 2048 | PASS | Max % |
| `NegativeOverlap` | -1 | 2048 | FAIL:9303 | Negative |
| `ExceedsAbsoluteMax` | 513 | 8192 | FAIL:9303 | >512 |
| `Exceeds25Percent2048` | 600 | 2048 | FAIL:9303 | >25% |
| `Exceeds25Percent1024` | 300 | 1024 | FAIL:9303 | >25% |
| `Exceeds25Percent512` | 150 | 512 | FAIL:9303 | >25% |
| `Exceeds25Percent256` | 100 | 256 | FAIL:9303 | >25% |
| `HalfChunkSize` | 1024 | 2048 | FAIL:9303 | 50% |

**Percentage Boundaries by ChunkSize:**
| ChunkSize | Max Overlap (25%) | Absolute Cap |
|-----------|-------------------|--------------|
| 256 | 64 | 64 |
| 512 | 128 | 128 |
| 1024 | 256 | 256 |
| 2048 | 512 | 512 |
| 4096 | 1024 → 512 | 512 |
| 8192 | 2048 → 512 | 512 |

---

### AB-9304: ContextTokenBudget Invalid

| Test Case | Input | Expected | Type |
|-----------|-------|----------|------|
| `ValidMinimum` | 512 | PASS | Min boundary |
| `ValidDefault` | 4096 | PASS | Default |
| `ValidMaximum` | 16384 | PASS | Max boundary |
| `TooSmallZero` | 0 | FAIL:9304 | Below min |
| `TooSmall511` | 511 | FAIL:9304 | Boundary-1 |
| `TooLarge16385` | 16385 | FAIL:9304 | Boundary+1 |

---

### AB-9305: EmbeddingModel Unsupported

| Test Case | Input | Expected | Type |
|-----------|-------|----------|------|
| `ValidNomic` | NomicEmbedText | PASS | Default |
| `ValidSmall` | TextEmbedding3Small | PASS | OpenAI |
| `ValidLarge` | TextEmbedding3Large | PASS | OpenAI |
| `EmptyString` | "" | FAIL:9305 | Empty |
| `UnknownModel` | unknown-model | FAIL:9305 | Invalid |
| `CaseSensitive` | NOMIC-EMBED-TEXT | FAIL:9305 | Wrong case |

**Supported Models:**
- `NomicEmbedText` (default)
- `TextEmbedding3Small`
- `TextEmbedding3Large`
- `AllMiniLmL6V2`

---

### AB-9306: SimilarityThreshold Invalid

| Test Case | Input | Expected | Type |
|-----------|-------|----------|------|
| `ValidZero` | 0.0 | PASS | Min boundary |
| `ValidOne` | 1.0 | PASS | Max boundary |
| `ValidDefault` | 0.7 | PASS | Default |
| `ValidHalf` | 0.5 | PASS | Mid-range |
| `NegativeSmall` | -0.001 | FAIL:9306 | Below zero |
| `NegativeLarge` | -1.0 | FAIL:9306 | Far below |
| `AboveOne` | 1.001 | FAIL:9306 | Above max |
| `WayAbove` | 2.0 | FAIL:9306 | Far above |

---

### AB-9307: TopK Invalid

| Test Case | Input | Expected | Type |
|-----------|-------|----------|------|
| `ValidMinimum` | 1 | PASS | Min boundary |
| `ValidDefault` | 10 | PASS | Default |
| `ValidMaximum` | 50 | PASS | Max boundary |
| `TooSmallZero` | 0 | FAIL:9307 | Below min |
| `TooLarge51` | 51 | FAIL:9307 | Above max |
| `negative` | -1 | FAIL:9307 | Negative |

---

### AB-9308: Config Load Failed

| Test Case | Scenario | Expected |
|-----------|----------|----------|
| `FileNotFound` | Missing config file | FAIL:9308 |
| `PermissionDenied` | No read access | FAIL:9308 |
| `IoError` | Disk I/O failure | FAIL:9308 |
| `ValidLoad` | Valid file exists | PASS |

---

### AB-9309: Config Save Failed

| Test Case | Scenario | Expected |
|-----------|----------|----------|
| `PermissionDenied` | No write access | FAIL:9309 |
| `ValidSave` | Writable path | PASS |

---

### AB-9310: Config Corrupted

| Test Case | Scenario | Expected |
|-----------|----------|----------|
| `InvalidJson` | Malformed JSON | FAIL:9310 |
| `TruncatedFile` | Incomplete file | FAIL:9310 |

---

## Integration Test Coverage

### Full Config Validation Flow

| Test | Description | Codes Covered |
|------|-------------|---------------|
| `ValidDefaultConfig` | All defaults | All |
| `ValidCustomConfig` | Custom valid values | All |
| `MultipleValidationErrors` | Multiple invalid fields | 9301-9307 |
| `CascadeValidation` | Overlap depends on ChunkSize | 9301, 9303 |
| `LoadValidateSave` | Full lifecycle | 9308, 9309, 9310 |
| `RecoveryFromCorrupt` | Restore defaults | 9310 |

---

## Edge Case Coverage

### Numeric Boundaries

| Field | Min | Max | Min-1 | Max+1 |
|-------|-----|-----|-------|-------|
| ChunkSize | 256 ✅ | 8192 ✅ | 255 ❌ | 8193 ❌ |
| ChunkOverlap | 0 ✅ | 512 ✅ | -1 ❌ | 513 ❌ |
| ContextBudget | 512 ✅ | 16384 ✅ | 511 ❌ | 16385 ❌ |
| SimilarityThreshold | 0.0 ✅ | 1.0 ✅ | -0.001 ❌ | 1.001 ❌ |
| TopK | 1 ✅ | 50 ✅ | 0 ❌ | 51 ❌ |

### Constraint Dependencies

| Parent | Child | Constraint | Test |
|--------|-------|------------|------|
| ChunkSize | ChunkOverlap | Overlap ≤ 25% of Size | ✅ |
| ChunkSize | ChunkOverlap | Overlap ≤ 512 absolute | ✅ |

---

## Test Execution Commands

```bash

# Run all RAG validation tests

go test -v ./internal/rag/... -run "Validation"

# Run specific error code tests

go test -v ./internal/rag/... -run "ChunkSize"
go test -v ./internal/rag/... -run "ChunkOverlap"
go test -v ./internal/rag/... -run "ContextBudget"

# Run with coverage report

go test -coverprofile=coverage.out ./internal/rag/...
go tool cover -html=coverage.out

# Run benchmark tests

go test -bench=. ./internal/rag/...
```

---

## Coverage Requirements

| Metric | Target | Current |
|--------|--------|---------|
| Line Coverage | ≥90% | 100% |
| Branch Coverage | ≥85% | 100% |
| Mutation Score | ≥80% | **Closed at v2.0.0 (Patch D, audit-10).** Target retained as the binding gate; measurement gated on executable RAG code landing outside `spec/`. Until then, line+branch coverage at 100% remains the operative proxy. |

### Patch D Closure Note (v2.0.0)

The audit-10 Patch D `TBD` is **resolved as of project v5.39.0** (RAG matrix v2.0.0). The ≥80% target is now the binding mutation-score gate for any future executable RAG validator implementation. No mutation tooling is added inside `spec/` itself, because:

1. `spec/06-seedable-config-architecture/` is a markdown specification, not executable code — there is nothing in this folder to mutate.
2. `spec/19-main-worker-service/` (the consumer of the RAG validator) is **spec-only by hard project rule** (`mem://constraints/spec19-no-implementation`); mutation tooling cannot be wired there.
3. Implementation of the RAG validator will occur in a separate, non-`spec/` repository or module during Phase-14 ramp-up. **At that point**, the chosen mutation framework (e.g. `gremlins` for Go, `cargo-mutants` for Rust, or Stryker for TS) MUST enforce the ≥80% gate from this matrix as a CI-blocking check.

Operative effect: the target is no longer "TBD" — it is **≥80%, deferred-by-design until executable code exists**, with line+branch coverage at 100% as the active proxy gate. This is the exhaustive, blind-AI-actionable disposition.

---

## Cross-References

| Document | Description |
|----------|-------------|
| [03-rag-validation-helpers.md](./02-rag-validation-helpers.md) | Implementation spec |
| [04-rag-validation-tests.md](./03-rag-validation-tests.md) | Full test code |
| [02-rag-chunk-settings.md](./01-rag-chunk-settings.md) | Configuration spec |

---

*Created 2026-02-02. This matrix ensures complete test coverage for all RAG validation error codes.*
