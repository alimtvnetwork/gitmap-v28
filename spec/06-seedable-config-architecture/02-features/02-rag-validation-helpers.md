# RAG Configuration Validation Helpers

**Version:** 3.2.0  
**Created:** 2026-02-02  
**Updated:** 2026-04-16  
**Status:** Active  
**Parent:** [02-rag-chunk-settings.md](./01-rag-chunk-settings.md)

---

## Overview

This document specifies validation helper patterns for RAG configuration settings in Go, including validators, error handling, and API integration.

---

## Error Code Mapping

| Error Code | Name | Validation Rule |
|------------|------|-----------------|
| AB-9301 | `RagChunkSizeInvalid` | ChunkSize outside 256-8192 range |
| AB-9302 | `RagChunkSizeNotMultiple` | ChunkSize not multiple of 256 |
| AB-9303 | `RagOverlapTooLarge` | ChunkOverlap > 25% of ChunkSize |
| AB-9304 | `RagContextBudgetInvalid` | ContextTokenBudget outside 512-16384 |
| AB-9305 | `RagEmbeddingModelInvalid` | Unsupported embedding model |
| AB-9306 | `RagSimilarityThresholdInvalid` | Threshold outside 0.0-1.0 |
| AB-9307 | `RagTopkInvalid` | TopK outside 1-50 range |
| AB-9308 | `RagConfigLoadFailed` | Failed to load configuration |
| AB-9309 | `RagConfigSaveFailed` | Failed to save configuration |
| AB-9310 | `RagConfigSourceConflict` | Conflicting multi-source config |

---

## Go Implementation

### Core Types

```go
package rag

import (
    "fmt"
)

// RagConfig represents the complete RAG configuration
type RagConfig struct {
    ChunkSize           int
    ChunkOverlap        int
    ContextTokenBudget  int
    EmbeddingModel      string
    SimilarityThreshold float64
    TopK                int
    Source              string // "seed", "root", "app"
}

// Supported embedding models (PascalCase identifiers mapped to external strings)
var SupportedEmbeddingModels = map[string]bool{
    "NomicEmbedText":       true,
    "TextEmbedding3Small":  true,
    "TextEmbedding3Large":  true,
    "AllMiniLmL6V2":        true,
}
```

### Validator Interface

```go
// ConfigValidator defines the validation interface
type ConfigValidator interface {
    Validate(config *RagConfig) []*apperror.AppError
}

// DefaultValidator implements ConfigValidator
type DefaultValidator struct{}

func (v *DefaultValidator) Validate(config *RagConfig) []*apperror.AppError {
    var errors []*apperror.AppError
    
    // ChunkSize validation
    if err := v.validateChunkSize(config.ChunkSize); err != nil {
        errors = append(errors, *err)
    }
    
    // ChunkOverlap validation (depends on ChunkSize)
    if err := v.validateChunkOverlap(config.ChunkOverlap, config.ChunkSize); err != nil {
        errors = append(errors, *err)
    }
    
    // ContextTokenBudget validation
    if err := v.validateContextBudget(config.ContextTokenBudget); err != nil {
        errors = append(errors, *err)
    }
    
    // EmbeddingModel validation
    if err := v.validateEmbeddingModel(config.EmbeddingModel); err != nil {
        errors = append(errors, *err)
    }
    
    // SimilarityThreshold validation
    if err := v.validateSimilarityThreshold(config.SimilarityThreshold); err != nil {
        errors = append(errors, *err)
    }
    
    // TopK validation
    if err := v.validateTopK(config.TopK); err != nil {
        errors = append(errors, *err)
    }
    
    return errors
}
```

### Individual Validators

```go
func (v *DefaultValidator) validateChunkSize(size int) *apperror.AppError {
    if size < 256 || size > 8192 {
        return apperror.New(
            ErrRagChunkSizeInvalid,
            "Chunk size outside valid range",
        ).WithContext("Field", "ChunkSize").
            WithContext("Value", size).
            WithContext("Expected", "256-8192")
    }

    if size%256 != 0 {
        return apperror.New(
            ErrRagChunkSizeNotMultiple,
            "Chunk size must be multiple of 256",
        ).WithContext("Field", "ChunkSize").
            WithContext("Value", size).
            WithContext("Expected", "Multiple of 256")
    }

    return nil
}

func (v *DefaultValidator) validateChunkOverlap(overlap, chunkSize int) *apperror.AppError {
    if overlap < 0 || overlap > 512 {
        return apperror.New(
            ErrRagOverlapTooLarge,
            "Chunk overlap outside valid range",
        ).WithContext("Field", "ChunkOverlap").
            WithContext("Value", overlap).
            WithContext("Expected", "0-512")
    }

    // Overlap cannot exceed 25% of chunk size
    maxOverlap := chunkSize / 4
    if overlap > maxOverlap {
        return apperror.New(
            ErrRagOverlapTooLarge,
            fmt.Sprintf("Chunk overlap cannot exceed 25%% of chunk size (%d)", maxOverlap),
        ).WithContext("Field", "ChunkOverlap").
            WithContext("Value", overlap).
            WithContext("MaxOverlap", maxOverlap)
    }

    return nil
}

func (v *DefaultValidator) validateContextBudget(budget int) *apperror.AppError {
    if budget < 512 || budget > 16384 {
        return apperror.New(
            ErrRagContextBudgetInvalid,
            "Context token budget outside valid range",
        ).WithContext("Field", "ContextTokenBudget").
            WithContext("Value", budget).
            WithContext("Expected", "512-16384")
    }

    return nil
}

func (v *DefaultValidator) validateEmbeddingModel(model string) *apperror.AppError {
    if !SupportedEmbeddingModels[model] {
        return apperror.New(
            ErrRagEmbeddingModelInvalid,
            "Embedding model not supported",
        ).WithContext("Field", "EmbeddingModel").
            WithContext("Value", model).
            WithContext("Expected", "NomicEmbedText, TextEmbedding3Small, TextEmbedding3Large, AllMiniLmL6V2")
    }

    return nil
}

func (v *DefaultValidator) validateSimilarityThreshold(threshold float64) *apperror.AppError {
    if threshold < 0.0 || threshold > 1.0 {
        return apperror.New(
            ErrRagSimilarityThresholdInvalid,
            "Similarity threshold outside valid range",
        ).WithContext("Field", "SimilarityThreshold").
            WithContext("Value", threshold).
            WithContext("Expected", "0.0-1.0")
    }

    return nil
}

func (v *DefaultValidator) validateTopK(topK int) *apperror.AppError {
    if topK < 1 || topK > 50 {
        return apperror.New(
            ErrRagTopkInvalid,
            "TopK outside valid range",
        ).WithContext("Field", "TopK").
            WithContext("Value", topK).
            WithContext("Expected", "1-50")
    }

    return nil
}
```

### Validation Service

```go
// RagConfigService handles loading, validation, and saving of RAG config
type RagConfigService struct {
    validator ConfigValidator
    seedPath  string
    rootDb    string
}

func NewRagConfigService(seedPath, rootDb string) *RagConfigService {
    return &RagConfigService{
        validator: &DefaultValidator{},
        seedPath:  seedPath,
        rootDb:    rootDb,
    }
}

// Load retrieves config with priority resolution
func (s *RagConfigService) Load(appName string) apperror.Result[*RagConfig] {
    config := &RagConfig{}

    // 1. Load seed defaults
    seedConfig, err := s.loadSeed()
    if err != nil {
        return apperror.Fail[*RagConfig](apperror.Wrap(
            err,
            ErrRagConfigLoadFailed,
            "Failed to load seed configuration",
        ).WithContext("SeedPath", s.seedPath))
    }
    *config = *seedConfig
    config.Source = "seed"

    // 2. Override with root DB settings
    if rootConfig, err := s.loadRootDb(); err == nil {
        s.mergeConfig(config, rootConfig, "root")
    }

    // 3. Override with app-level settings
    if appName != "" {
        if appConfig, err := s.loadAppDb(appName); err == nil {
            s.mergeConfig(config, appConfig, "app")
        }
    }

    // Validate final config
    if errors := s.validator.Validate(config); len(errors) > 0 {
        return apperror.Fail[*RagConfig](errors[0])
    }

    return apperror.Ok(config)
}

// Save persists config with validation
func (s *RagConfigService) Save(appName string, config *RagConfig) apperror.Result[bool] {
    // Validate before saving
    if errors := s.validator.Validate(config); len(errors) > 0 {
        return apperror.Fail[bool](errors[0])
    }

    // Save to appropriate location
    if appName != "" {
        return s.saveAppDb(appName, config)
    }

    return s.saveRootDb(config)
}

// RagConfigUpdate holds typed partial update fields
type RagConfigUpdate struct {
    ChunkSize           *int     `json:"ChunkSize,omitempty"`
    ChunkOverlap        *int     `json:"ChunkOverlap,omitempty"`
    ContextTokenBudget  *int     `json:"ContextTokenBudget,omitempty"`
    EmbeddingModel      *string  `json:"EmbeddingModel,omitempty"`
    SimilarityThreshold *float64 `json:"SimilarityThreshold,omitempty"`
    TopK                *int     `json:"TopK,omitempty"`
}

// ValidatePartial validates only provided fields
func (s *RagConfigService) ValidatePartial(updates *RagConfigUpdate) []*apperror.AppError {
    var errors []*apperror.AppError
    
    if updates.ChunkSize != nil {
        if err := (&DefaultValidator{}).validateChunkSize(*updates.ChunkSize); err != nil {
            errors = append(errors, *err)
        }
    }
    
    // Check overlap with size if both provided
    if updates.ChunkOverlap != nil {
        chunkSize := 2048 // default
        if updates.ChunkSize != nil {
            chunkSize = *updates.ChunkSize
        }
        if err := (&DefaultValidator{}).validateChunkOverlap(*updates.ChunkOverlap, chunkSize); err != nil {
            errors = append(errors, *err)
        }
    }
    
    // Continue for other fields...
    return errors
}

func (s *RagConfigService) mergeConfig(base *RagConfig, override *RagConfig, source string) {
    if override.ChunkSize != 0 {
        base.ChunkSize = override.ChunkSize
        base.Source = source
    }
    if override.ChunkOverlap != 0 {
        base.ChunkOverlap = override.ChunkOverlap
    }
    if override.ContextTokenBudget != 0 {
        base.ContextTokenBudget = override.ContextTokenBudget
    }
    if override.EmbeddingModel != "" {
        base.EmbeddingModel = override.EmbeddingModel
    }
    if override.SimilarityThreshold != 0 {
        base.SimilarityThreshold = override.SimilarityThreshold
    }
    if override.TopK != 0 {
        base.TopK = override.TopK
    }
}
```

---

## API Integration

### Validation Middleware

```go
// ValidationMiddleware validates RAG config in requests
func ValidationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "PUT" || r.Method == "POST" {
            var config RagConfig
            if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
            }
            
            validator := &DefaultValidator{}
            if errors := validator.Validate(&config); len(errors) > 0 {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(RagValidationErrorResponse{
                    Errors: errors,
                })
                return
            }
            
            // Re-encode for next handler
            encoded, _ := json.Marshal(config)
            r.Body = io.NopCloser(bytes.NewReader(encoded))
        }
        next.ServeHTTP(w, r)
    })
}
```

### API Handler

```go
// UpdateRagSettings handles PUT /api/v1/settings/rag
func (h *SettingsHandler) UpdateRagSettings(w http.ResponseWriter, r *http.Request) {
    var updates RagConfigUpdate
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        h.respondError(w, 9308, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate partial updates
    if errors := h.configService.ValidatePartial(&updates); len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(RagValidationErrorResponse{
            Error:  &errors[0],
            Errors: errors,
        })
        return
    }
    
    // Load current config
    appName := r.URL.Query().Get("app")
    config, err := h.configService.Load(appName)
    if err != nil {
        h.respondError(w, 9308, "Failed to load config", http.StatusInternalServerError)
        return
    }
    
    // Apply typed updates
    if updates.ChunkSize != nil {
        config.ChunkSize = *updates.ChunkSize
    }
    if updates.ChunkOverlap != nil {
        config.ChunkOverlap = *updates.ChunkOverlap
    }
    if updates.ContextTokenBudget != nil {
        config.ContextTokenBudget = *updates.ContextTokenBudget
    }
    if updates.EmbeddingModel != nil {
        config.EmbeddingModel = *updates.EmbeddingModel
    }
    if updates.SimilarityThreshold != nil {
        config.SimilarityThreshold = *updates.SimilarityThreshold
    }
    if updates.TopK != nil {
        config.TopK = *updates.TopK
    }
    
    // Save updated config
    if err := h.configService.Save(appName, config); err != nil {
        var valErr *apperror.AppError
        if errors.As(err, &valErr) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(valErr)
            return
        }
        h.respondError(w, 9309, "Failed to save config", http.StatusInternalServerError)
        return
    }
    
    // Return updated config
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(RagConfigUpdateResponse{
        Config:  *config,
        Updated: getUpdatedFieldsFromUpdate(&updates),
        Source:  config.Source,
    })
}
```

---

## Response Formats

### Validation Success

```json
{
  "Config": {
    "ChunkSize": 4096,
    "ChunkOverlap": 200,
    "ContextTokenBudget": 4096,
    "EmbeddingModel": "NomicEmbedText",
    "SimilarityThreshold": 0.7,
    "TopK": 10,
    "Source": "app"
  },
  "Updated": ["ChunkSize", "ChunkOverlap"],
  "Source": "app"
}
```

### Validation Error

```json
{
  "Error": {
    "Code": 9301,
    "Name": "RagChunkSizeInvalid",
    "Message": "Chunk size outside valid range",
    "Field": "ChunkSize",
    "Value": 100,
    "Expected": "256-8192",
    "Timestamp": "2026-02-02T10:30:00Z"
  },
  "Errors": [
    {
      "Code": 9301,
      "Name": "RagChunkSizeInvalid",
      "Message": "Chunk size outside valid range",
      "Field": "ChunkSize",
      "Value": 100,
      "Expected": "256-8192"
    }
  ]
}
```

### Multi-Error Response

```json
{
  "Errors": [
    {
      "Code": 9301,
      "Name": "RagChunkSizeInvalid",
      "Field": "ChunkSize",
      "Value": 100
    },
    {
      "Code": 9303,
      "Name": "RagOverlapTooLarge",
      "Field": "ChunkOverlap",
      "Value": 600
    }
  ]
}
```

---

## Unit Tests

```go
func TestChunkSizeValidation(t *testing.T) {
    v := &DefaultValidator{}
    
    tests := []struct {
        name     string
        size     int
        wantCode int
    }{
        {"valid_2048", 2048, 0},
        {"valid_min", 256, 0},
        {"valid_max", 8192, 0},
        {"too_small", 100, 9301},
        {"too_large", 10000, 9301},
        {"not_multiple", 1000, 9302},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.validateChunkSize(tt.size)
            if tt.wantCode == 0 {
                if err != nil {
                    t.Errorf("expected no error, got %v", err)
                }
            } else {
                if err == nil {
                    t.Errorf("expected error code %d, got nil", tt.wantCode)
                } else if err.Code != tt.wantCode {
                    t.Errorf("expected error code %d, got %d", tt.wantCode, err.Code)
                }
            }
        })
    }
}

func TestOverlapPercentageValidation(t *testing.T) {
    v := &DefaultValidator{}
    
    tests := []struct {
        name      string
        chunkSize int
        overlap   int
        wantErr   bool
    }{
        {"25_percent_ok", 2048, 512, false},
        {"20_percent_ok", 2048, 400, false},
        {"26_percent_fail", 2048, 600, true},
        {"50_percent_fail", 2048, 1024, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.validateChunkOverlap(tt.overlap, tt.chunkSize)
            if tt.wantErr && err == nil {
                t.Error("expected error, got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error, got %v", err)
            }
        })
    }
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| RAG Chunk Settings | `./01-rag-chunk-settings.md` |
| Error Code Registry | `../03-error-manage/03-error-code-registry/00-overview.md` |
| AI Bridge Error Codes | `../03-error-manage/03-error-code-registry/00-overview.md` |
| Seedable Config Overview | `../00-overview.md` |
