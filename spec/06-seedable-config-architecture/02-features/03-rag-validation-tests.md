# RAG Validation Helpers: Unit Test Specification

**Version:** 3.2.0  
**Created:** 2026-02-02  
**Updated:** 2026-04-16  
**Status:** Active  
**Parent:** [03-rag-validation-helpers.md](./02-rag-validation-helpers.md)

---

## Overview

Comprehensive unit test specifications for RAG configuration validation helpers, covering all error codes (9301-9310) with edge cases, boundary conditions, and integration scenarios.

---

## Test Categories

| Category | Tests | Coverage |
|----------|-------|----------|
| ChunkSize Validation | 12 | AB-9301, AB-9302 |
| ChunkOverlap Validation | 10 | AB-9303 |
| ContextBudget Validation | 6 | AB-9304 |
| EmbeddingModel Validation | 6 | AB-9305 |
| SimilarityThreshold Validation | 8 | AB-9306 |
| TopK Validation | 6 | AB-9307 |
| Config Load/Save | 8 | AB-9308, AB-9309, AB-9310 |
| Integration Tests | 6 | Full config validation |

---

## ChunkSize Validation Tests (AB-9301, AB-9302)

### Test Suite: ChunkSize

```go
package rag_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "aibridge/internal/rag"
)

func TestChunkSizeValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name         string
        size         int
        expectError  bool
        expectedCode int
        description  string
    }{
        // Valid cases
        {
            name:        "ValidMinimum",
            size:        256,
            expectError: false,
            description: "Minimum valid chunk size",
        },
        {
            name:        "ValidDefault",
            size:        2048,
            expectError: false,
            description: "Default chunk size",
        },
        {
            name:        "ValidMaximum",
            size:        8192,
            expectError: false,
            description: "Maximum valid chunk size",
        },
        {
            name:        "ValidMultiple512",
            size:        512,
            expectError: false,
            description: "Valid multiple of 256",
        },
        {
            name:        "ValidMultiple1024",
            size:        1024,
            expectError: false,
            description: "Valid multiple of 256",
        },
        {
            name:        "ValidMultiple4096",
            size:        4096,
            expectError: false,
            description: "Large valid chunk size",
        },
        
        // Invalid range cases (AB-9301)
        {
            name:         "TooSmallZero",
            size:         0,
            expectError:  true,
            expectedCode: 9301,
            description:  "Zero is below minimum",
        },
        {
            name:         "TooSmall100",
            size:         100,
            expectError:  true,
            expectedCode: 9301,
            description:  "100 is below minimum 256",
        },
        {
            name:         "TooSmall255",
            size:         255,
            expectError:  true,
            expectedCode: 9301,
            description:  "Just below minimum boundary",
        },
        {
            name:         "TooLarge8193",
            size:         8193,
            expectError:  true,
            expectedCode: 9301,
            description:  "Just above maximum boundary",
        },
        {
            name:         "TooLarge16384",
            size:         16384,
            expectError:  true,
            expectedCode: 9301,
            description:  "Far above maximum",
        },
        {
            name:         "NegativeValue",
            size:         -256,
            expectError:  true,
            expectedCode: 9301,
            description:  "Negative values invalid",
        },
        
        // Invalid multiple cases (AB-9302)
        {
            name:         "NotMultiple257",
            size:         257,
            expectError:  true,
            expectedCode: 9302,
            description:  "Just above valid minimum but not multiple",
        },
        {
            name:         "NotMultiple1000",
            size:         1000,
            expectError:  true,
            expectedCode: 9302,
            description:  "Round number but not multiple of 256",
        },
        {
            name:         "NotMultiple1500",
            size:         1500,
            expectError:  true,
            expectedCode: 9302,
            description:  "In range but not multiple",
        },
        {
            name:         "NotMultiple2000",
            size:         2000,
            expectError:  true,
            expectedCode: 9302,
            description:  "Close to 2048 but not multiple",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateChunkSize(tt.size)
            
            if tt.expectError {
                require.NotNil(t, err, "expected error for %s", tt.description)
                assert.Equal(t, tt.expectedCode, err.Code, "wrong error code")
                assert.Equal(t, "ChunkSize", err.Field)
                assert.Equal(t, tt.size, err.Value)
            } else {
                assert.Nil(t, err, "expected no error for %s", tt.description)
            }
        })
    }
}

func TestChunkSizeBoundaryConditions(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    // Test exact boundaries
    assert.Nil(t, v.ValidateChunkSize(256), "minimum boundary should pass")
    assert.Nil(t, v.ValidateChunkSize(8192), "maximum boundary should pass")
    
    // Test just outside boundaries
    assert.NotNil(t, v.ValidateChunkSize(255), "just below minimum should fail")
    assert.NotNil(t, v.ValidateChunkSize(8448), "just above maximum should fail")
}
```

---

## ChunkOverlap Validation Tests (AB-9303)

### Test Suite: ChunkOverlap

```go
func TestChunkOverlapValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name         string
        overlap      int
        chunkSize    int
        expectError  bool
        expectedCode int
        description  string
    }{
        // Valid cases
        {
            name:        "ValidZero",
            overlap:     0,
            chunkSize:   2048,
            expectError: false,
            description: "Zero overlap is valid",
        },
        {
            name:        "ValidDefault100",
            overlap:     100,
            chunkSize:   2048,
            expectError: false,
            description: "Default 100 overlap with 2048 chunk",
        },
        {
            name:        "Valid25Percent",
            overlap:     512,
            chunkSize:   2048,
            expectError: false,
            description: "Exactly 25% of chunk size",
        },
        {
            name:        "Valid10Percent",
            overlap:     204,
            chunkSize:   2048,
            expectError: false,
            description: "10% of chunk size",
        },
        {
            name:        "ValidMaxAbsolute",
            overlap:     512,
            chunkSize:   8192,
            expectError: false,
            description: "Maximum absolute overlap",
        },
        {
            name:        "ValidSmallChunk",
            overlap:     64,
            chunkSize:   256,
            expectError: false,
            description: "25% of minimum chunk size",
        },
        
        // Invalid absolute range cases
        {
            name:         "NegativeOverlap",
            overlap:      -1,
            chunkSize:    2048,
            expectError:  true,
            expectedCode: 9303,
            description:  "Negative overlap invalid",
        },
        {
            name:         "ExceedsAbsoluteMax",
            overlap:      513,
            chunkSize:    8192,
            expectError:  true,
            expectedCode: 9303,
            description:  "Just above 512 absolute max",
        },
        
        // Invalid percentage cases (>25%)
        {
            name:         "Exceeds25Percent2048",
            overlap:      600,
            chunkSize:    2048,
            expectError:  true,
            expectedCode: 9303,
            description:  "600 > 512 (25% of 2048)",
        },
        {
            name:         "Exceeds25Percent1024",
            overlap:      300,
            chunkSize:    1024,
            expectError:  true,
            expectedCode: 9303,
            description:  "300 > 256 (25% of 1024)",
        },
        {
            name:         "Exceeds25Percent512",
            overlap:      150,
            chunkSize:    512,
            expectError:  true,
            expectedCode: 9303,
            description:  "150 > 128 (25% of 512)",
        },
        {
            name:         "Exceeds25Percent256",
            overlap:      100,
            chunkSize:    256,
            expectError:  true,
            expectedCode: 9303,
            description:  "100 > 64 (25% of 256)",
        },
        {
            name:         "HalfChunkSize",
            overlap:      1024,
            chunkSize:    2048,
            expectError:  true,
            expectedCode: 9303,
            description:  "50% overlap is invalid",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateChunkOverlap(tt.overlap, tt.chunkSize)
            
            if tt.expectError {
                require.NotNil(t, err, "expected error for %s", tt.description)
                assert.Equal(t, tt.expectedCode, err.Code)
                assert.Equal(t, "ChunkOverlap", err.Field)
            } else {
                assert.Nil(t, err, "expected no error for %s", tt.description)
            }
        })
    }
}

func TestChunkOverlapPercentageEdgeCases(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    // Test exact 25% boundaries for various chunk sizes
    testCases := []struct {
        chunkSize  int
        maxOverlap int
    }{
        {256, 64},
        {512, 128},
        {1024, 256},
        {2048, 512},
        {4096, 512}, // Capped at 512 max
        {8192, 512}, // Capped at 512 max
    }
    
    for _, tc := range testCases {
        t.Run(fmt.Sprintf("chunk_%d", tc.chunkSize), func(t *testing.T) {
            // Exact boundary should pass
            err := v.ValidateChunkOverlap(tc.maxOverlap, tc.chunkSize)
            assert.Nil(t, err, "exact 25%% boundary should pass for chunk %d", tc.chunkSize)
            
            // One above should fail (unless at absolute max)
            if tc.maxOverlap < 512 {
                err = v.ValidateChunkOverlap(tc.maxOverlap+1, tc.chunkSize)
                assert.NotNil(t, err, "just above 25%% should fail for chunk %d", tc.chunkSize)
            }
        })
    }
}
```

---

## ContextBudget Validation Tests (AB-9304)

```go
func TestContextBudgetValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name         string
        budget       int
        expectError  bool
        expectedCode int
    }{
        // Valid cases
        {"ValidMinimum", 512, false, 0},
        {"ValidDefault", 4096, false, 0},
        {"ValidMaximum", 16384, false, 0},
        {"ValidMidRange", 8192, false, 0},
        {"Valid1024", 1024, false, 0},
        
        // Invalid cases
        {"TooSmallZero", 0, true, 9304},
        {"TooSmall511", 511, true, 9304},
        {"TooSmall256", 256, true, 9304},
        {"TooLarge16385", 16385, true, 9304},
        {"TooLarge32768", 32768, true, 9304},
        {"negative", -1, true, 9304},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateContextBudget(tt.budget)
            
            if tt.expectError {
                require.NotNil(t, err)
                assert.Equal(t, tt.expectedCode, err.Code)
                assert.Equal(t, "ContextTokenBudget", err.Field)
            } else {
                assert.Nil(t, err)
            }
        })
    }
}
```

---

## EmbeddingModel Validation Tests (AB-9305)

```go
func TestEmbeddingModelValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name         string
        model        string
        expectError  bool
        expectedCode int
    }{
        // Valid models
        {"ValidNomic", "NomicEmbedText", false, 0},
        {"ValidSmall", "TextEmbedding3Small", false, 0},
        {"ValidLarge", "TextEmbedding3Large", false, 0},
        {"ValidMinilm", "AllMiniLmL6V2", false, 0},
        
        // Invalid models
        {"EmptyString", "", true, 9305},
        {"UnknownModel", "unknown-model", true, 9305},
        {"TypoNomic", "nomic-embed", true, 9305},
        {"CaseSensitive", "NOMIC-EMBED-TEXT", true, 9305},
        {"spaces", "nomic embed text", true, 9305},
        {"PartialMatch", "text-embedding", true, 9305},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateEmbeddingModel(tt.model)
            
            if tt.expectError {
                require.NotNil(t, err)
                assert.Equal(t, tt.expectedCode, err.Code)
                assert.Equal(t, "EmbeddingModel", err.Field)
                assert.Equal(t, tt.model, err.Value)
            } else {
                assert.Nil(t, err)
            }
        })
    }
}
```

---

## SimilarityThreshold Validation Tests (AB-9306)

```go
func TestSimilarityThresholdValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name         string
        threshold    float64
        expectError  bool
        expectedCode int
    }{
        // Valid cases
        {"ValidZero", 0.0, false, 0},
        {"ValidOne", 1.0, false, 0},
        {"ValidDefault", 0.7, false, 0},
        {"ValidHalf", 0.5, false, 0},
        {"ValidLow", 0.1, false, 0},
        {"ValidHigh", 0.9, false, 0},
        {"ValidSmall", 0.001, false, 0},
        {"ValidNearOne", 0.999, false, 0},
        
        // Invalid cases
        {"NegativeSmall", -0.001, true, 9306},
        {"NegativeLarge", -1.0, true, 9306},
        {"AboveOne", 1.001, true, 9306},
        {"WayAbove", 2.0, true, 9306},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateSimilarityThreshold(tt.threshold)
            
            if tt.expectError {
                require.NotNil(t, err)
                assert.Equal(t, tt.expectedCode, err.Code)
                assert.Equal(t, "SimilarityThreshold", err.Field)
            } else {
                assert.Nil(t, err)
            }
        })
    }
}

func TestSimilarityThresholdFloatPrecision(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    // Test float precision edge cases
    edgeCases := []float64{
        0.0000001,           // Very small positive
        0.9999999,           // Very close to 1
        1.0 - 1e-10,         // Float precision near 1
        0.0 + 1e-10,         // Float precision near 0
    }
    
    for _, threshold := range edgeCases {
        err := v.ValidateSimilarityThreshold(threshold)
        assert.Nil(t, err, "threshold %v should be valid", threshold)
    }
}
```

---

## TopK Validation Tests (AB-9307)

```go
func TestTopKValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name         string
        topK         int
        expectError  bool
        expectedCode int
    }{
        // Valid cases
        {"ValidMinimum", 1, false, 0},
        {"ValidDefault", 10, false, 0},
        {"ValidMaximum", 50, false, 0},
        {"ValidMid", 25, false, 0},
        
        // Invalid cases
        {"zero", 0, true, 9307},
        {"negative", -1, true, 9307},
        {"AboveMax", 51, true, 9307},
        {"WayAbove", 100, true, 9307},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateTopK(tt.topK)
            
            if tt.expectError {
                require.NotNil(t, err)
                assert.Equal(t, tt.expectedCode, err.Code)
                assert.Equal(t, "TopK", err.Field)
            } else {
                assert.Nil(t, err)
            }
        })
    }
}
```

---

## Full Config Validation Tests

```go
func TestFullConfigValidation(t *testing.T) {
    v := &rag.DefaultValidator{}
    
    tests := []struct {
        name        string
        config      rag.RagConfig
        errorCount  int
        errorCodes  []int
    }{
        {
            name: "ValidConfig",
            config: rag.RagConfig{
                ChunkSize:           2048,
                ChunkOverlap:        100,
                ContextTokenBudget:  4096,
                EmbeddingModel:      "NomicEmbedText",
                SimilarityThreshold: 0.7,
                TopK:                10,
            },
            errorCount: 0,
        },
        {
            name: "SingleErrorChunkSize",
            config: rag.RagConfig{
                ChunkSize:           100, // Invalid
                ChunkOverlap:        100,
                ContextTokenBudget:  4096,
                EmbeddingModel:      "NomicEmbedText",
                SimilarityThreshold: 0.7,
                TopK:                10,
            },
            errorCount: 1,
            errorCodes: []int{9301},
        },
        {
            name: "MultipleErrors",
            config: rag.RagConfig{
                ChunkSize:           100,      // Invalid range
                ChunkOverlap:        1000,     // Too large
                ContextTokenBudget:  100,      // Too small
                EmbeddingModel:      "invalid",
                SimilarityThreshold: 2.0,      // Out of range
                TopK:                100,      // Too large
            },
            errorCount: 6,
            errorCodes: []int{9301, 9303, 9304, 9305, 9306, 9307},
        },
        {
            name: "OverlapPercentageError",
            config: rag.RagConfig{
                ChunkSize:           1024,
                ChunkOverlap:        300, // > 25% of 1024
                ContextTokenBudget:  4096,
                EmbeddingModel:      "NomicEmbedText",
                SimilarityThreshold: 0.7,
                TopK:                10,
            },
            errorCount: 1,
            errorCodes: []int{9303},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            errors := v.Validate(&tt.config)
            
            assert.Len(t, errors, tt.errorCount, "wrong number of errors")
            
            if len(tt.errorCodes) > 0 {
                for i, expectedCode := range tt.errorCodes {
                    if i < len(errors) {
                        assert.Equal(t, expectedCode, errors[i].Code)
                    }
                }
            }
        })
    }
}
```

---

## Config Load/Save Tests (AB-9308, AB-9309, AB-9310)

```go
func TestConfigLoad(t *testing.T) {
    t.Run("LoadSuccessFromSeed", func(t *testing.T) {
        service := rag.NewRagConfigService("testdata/config.seed.json", ":memory:")
        config, err := service.Load("")
        
        require.NoError(t, err)
        assert.Equal(t, 2048, config.ChunkSize)
        assert.Equal(t, "seed", config.Source)
    })
    
    t.Run("LoadFailureMissingSeed", func(t *testing.T) {
        service := rag.NewRagConfigService("nonexistent.json", ":memory:")
        _, err := service.Load("")
        
        require.Error(t, err)
        ragErr, ok := err.(*apperror.AppError)
        require.True(t, ok)
        assert.Equal(t, 9308, ragErr.Code)
    })
    
    t.Run("LoadAppOverride", func(t *testing.T) {
        service := setupTestService(t)
        
        // Set app-level override
        err := service.SaveAppSetting("testapp", "ChunkSize", 4096)
        require.NoError(t, err)
        
        config, err := service.Load("testapp")
        require.NoError(t, err)
        assert.Equal(t, 4096, config.ChunkSize)
        assert.Equal(t, "app", config.Source)
    })
}

func TestConfigSave(t *testing.T) {
    t.Run("SaveSuccess", func(t *testing.T) {
        service := setupTestService(t)
        
        config := &rag.RagConfig{
            ChunkSize:           4096,
            ChunkOverlap:        200,
            ContextTokenBudget:  8192,
            EmbeddingModel:      "NomicEmbedText",
            SimilarityThreshold: 0.8,
            TopK:                20,
        }
        
        err := service.Save("testapp", config)
        assert.NoError(t, err)
    })
    
    t.Run("SaveFailureValidation", func(t *testing.T) {
        service := setupTestService(t)
        
        config := &rag.RagConfig{
            ChunkSize: 100, // Invalid
        }
        
        err := service.Save("testapp", config)
        require.Error(t, err)
        
        ragErr, ok := err.(*apperror.AppError)
        require.True(t, ok)
        assert.Equal(t, 9301, ragErr.Code)
    })
    
    t.Run("SaveFailureDbError", func(t *testing.T) {
        service := rag.NewRagConfigService("testdata/config.seed.json", "/invalid/path/db.sqlite")
        
        config := &rag.RagConfig{
            ChunkSize:           2048,
            ChunkOverlap:        100,
            ContextTokenBudget:  4096,
            EmbeddingModel:      "NomicEmbedText",
            SimilarityThreshold: 0.7,
            TopK:                10,
        }
        
        err := service.Save("", config)
        require.Error(t, err)
        
        ragErr, ok := err.(*apperror.AppError)
        require.True(t, ok)
        assert.Equal(t, 9309, ragErr.Code)
    })
}

func TestConfigSourceConflict(t *testing.T) {
    t.Run("ConflictDetection", func(t *testing.T) {
        // Test when same setting has conflicting values from different sources
        service := setupTestService(t)
        
        // Set conflicting values
        service.SetRootSetting("ChunkSize", 2048)
        service.SetAppSetting("testapp", "ChunkSize", 4096)
        service.SetSeedDefault("ChunkSize", 1024)
        
        // Should resolve by priority (app > root > seed)
        config, err := service.Load("testapp")
        require.NoError(t, err)
        assert.Equal(t, 4096, config.ChunkSize)
        assert.Equal(t, "app", config.Source)
    })
}
```

---

## Table-Driven Test Helpers

```go
// TestHelper provides common setup for RAG validation tests
type TestHelper struct {
    t       *testing.T
    service *rag.RagConfigService
}

func NewTestHelper(t *testing.T) *TestHelper {
    service := rag.NewRagConfigService(
        "testdata/config.seed.json",
        ":memory:",
    )
    return &TestHelper{t: t, service: service}
}

func (h *TestHelper) AssertValidConfig(config *rag.RagConfig) {
    h.t.Helper()
    v := &rag.DefaultValidator{}
    errors := v.Validate(config)
    assert.Empty(h.t, errors, "expected valid config, got errors: %v", errors)
}

func (h *TestHelper) AssertErrorCode(err error, expectedCode int) {
    h.t.Helper()
    ragErr, ok := err.(*apperror.AppError)
    require.True(h.t, ok, "expected RagValidationError, got %T", err)
    assert.Equal(h.t, expectedCode, ragErr.Code)
}
```

---

## Benchmark Tests

```go
func BenchmarkValidateFullConfig(b *testing.B) {
    v := &rag.DefaultValidator{}
    config := &rag.RagConfig{
        ChunkSize:           2048,
        ChunkOverlap:        100,
        ContextTokenBudget:  4096,
        EmbeddingModel:      "NomicEmbedText",
        SimilarityThreshold: 0.7,
        TopK:                10,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        v.Validate(config)
    }
}

func BenchmarkValidateChunkSize(b *testing.B) {
    v := &rag.DefaultValidator{}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        v.ValidateChunkSize(2048)
    }
}
```

---

## Test Data Files

### testdata/config.seed.json

```json
{
  "$schema": "./config.schema.json",
  "version": "1.3.0",
  "categories": {
    "rag": {
      "DisplayName": "RAG Configuration",
      "settings": {
        "ChunkSize": {
          "type": "number",
          "label": "Chunk Size",
          "default": 2048,
          "min": 256,
          "max": 8192
        },
        "ChunkOverlap": {
          "type": "number",
          "label": "Chunk Overlap",
          "default": 100,
          "min": 0,
          "max": 512
        },
        "ContextTokenBudget": {
          "type": "number",
          "label": "Context Budget",
          "default": 4096
        },
        "EmbeddingModel": {
          "type": "string",
          "label": "Embedding Model",
          "default": "NomicEmbedText"
        },
        "SimilarityThreshold": {
          "type": "number",
          "label": "Similarity",
          "default": 0.7
        },
        "TopK": {
          "type": "number",
          "label": "Top K",
          "default": 10
        }
      }
    }
  }
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| RAG Validation Helpers | `./02-rag-validation-helpers.md` |
| RAG Chunk Settings | `./02-rag-chunk-settings.md` |
| Error Code Registry | `../03-error-code-registry/01-registry.md` |
| AI Bridge Error Codes | `../22-ai-bridge-cli/01-backend/05-error-codes.md` |
