# RAG Chunk Configuration Settings

**Version:** 3.2.0  
**Created:** 2026-02-02  
**Updated:** 2026-04-16  
**Status:** Active  
**Parent:** [00-overview.md](../00-overview.md)

---

## Overview

This document specifies the RAG (Retrieval-Augmented Generation) chunking configuration settings, including defaults, validation rules, and override mechanisms.

---

## Configuration Settings

### Chunk Settings

| Setting Key | Type | Default | Min | Max | Description |
|-------------|------|---------|-----|-----|-------------|
| `Rag.ChunkSize` | int | 2048 | 256 | 8192 | Tokens per chunk |
| `Rag.ChunkOverlap` | int | 100 | 0 | 512 | Overlap between chunks |
| `Rag.ContextTokenBudget` | int | 4096 | 512 | 16384 | Max tokens for RAG context |
| `Rag.EmbeddingModel` | string | `NomicEmbedText` | - | - | Embedding model |
| `Rag.SimilarityThreshold` | float | 0.7 | 0.0 | 1.0 | Min similarity for inclusion |
| `Rag.TopK` | int | 10 | 1 | 50 | Max chunks returned |

---

## Seed Configuration

### config.seed.json Structure

```json
{
  "$schema": "./config.schema.json",
  "Version": "1.3.0",
  "Changelog": "Added configurable RAG chunk settings",
  "Categories": {
    "Rag": {
      "DisplayName": "RAG Configuration",
      "Description": "Retrieval-Augmented Generation settings",
      "Version": "1.3.0",
      "AddedIn": "1.3.0",
      "Settings": {
        "ChunkSize": {
          "Type": "number",
          "Label": "Chunk Size (tokens)",
          "Description": "Number of tokens per chunk. Larger chunks provide more context but reduce granularity.",
          "Default": 2048,
          "Min": 256,
          "Max": 8192
        },
        "ChunkOverlap": {
          "Type": "number",
          "Label": "Chunk Overlap (tokens)",
          "Description": "Token overlap between adjacent chunks to maintain context continuity.",
          "Default": 100,
          "Min": 0,
          "Max": 512
        },
        "ContextTokenBudget": {
          "Type": "number",
          "Label": "Context Token Budget",
          "Description": "Maximum tokens for RAG context injection into prompts.",
          "Default": 4096,
          "Min": 512,
          "Max": 16384
        },
        "EmbeddingModel": {
          "Type": "select",
          "Label": "Embedding Model",
          "Description": "Model used for generating vector embeddings.",
          "Default": "NomicEmbedText",
          "Options": [
            "NomicEmbedText",
            "TextEmbedding3Small",
            "TextEmbedding3Large",
            "AllMiniLmL6V2"
          ]
        },
        "SimilarityThreshold": {
          "Type": "number",
          "Label": "Similarity Threshold",
          "Description": "Minimum cosine similarity for chunk inclusion (0.0-1.0).",
          "Default": 0.7,
          "Min": 0.0,
          "Max": 1.0
        },
        "TopK": {
          "Type": "number",
          "Label": "Top-K Chunks",
          "Description": "Maximum number of chunks to return from similarity search.",
          "Default": 10,
          "Min": 1,
          "Max": 50
        }
      }
    }
  }
}
```

---

## Validation Rules

### ChunkSize Validation

```go
type ChunkSizeValidator struct{}

func (v ChunkSizeValidator) Validate(value int) error {
    if value < 256 {
        return apperror.New(
            ErrChunkSizeTooSmall,
            "ChunkSize must be >= 256",
        ).WithContext("Value", value)
    }
    if value > 8192 {
        return apperror.New(
            ErrChunkSizeTooLarge,
            "ChunkSize must be <= 8192",
        ).WithContext("Value", value)
    }
    // Must be power of 2 or multiple of 256
    if value%256 != 0 {
        return apperror.New(
            ErrChunkSizeAlignment,
            "ChunkSize must be multiple of 256",
        ).WithContext("Value", value)
    }
    return nil
}
```

### ChunkOverlap Validation

```go
type ChunkOverlapValidator struct {
    ChunkSize int
}

func (v ChunkOverlapValidator) Validate(value int) error {
    if value < 0 {
        return apperror.New(
            ErrChunkOverlapNegative,
            "ChunkOverlap must be >= 0",
        ).WithContext("Value", value)
    }
    if value > 512 {
        return apperror.New(
            ErrChunkOverlapTooLarge,
            "ChunkOverlap must be <= 512",
        ).WithContext("Value", value)
    }
    // Overlap cannot exceed 25% of chunk size
    maxOverlap := v.ChunkSize / 4
    if value > maxOverlap {
        return apperror.New(
            ErrChunkOverlapExceedsRatio,
            "ChunkOverlap cannot exceed 25% of ChunkSize",
        ).WithContext("Value", value).
            WithContext("MaxOverlap", maxOverlap)
    }
    return nil
}
```

### Context Budget Validation

```go
type ContextBudgetValidator struct{}

func (v ContextBudgetValidator) Validate(value int) error {
    if value < 512 {
        return apperror.New(
            ErrContextBudgetTooSmall,
            "ContextTokenBudget must be >= 512",
        ).WithContext("Value", value)
    }
    if value > 16384 {
        return apperror.New(
            ErrContextBudgetTooLarge,
            "ContextTokenBudget must be <= 16384",
        ).WithContext("Value", value)
    }
    return nil
}
```

---

## Configuration Priority

Settings are resolved in the following priority order (highest first):

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        RAG CONFIGURATION PRIORITY                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   1. App-Level Override (Highest Priority)                                   │
│      └─▶ data/{appName}/settings/config.db → Settings table                │
│                                                                              │
│   2. Root DB Setting                                                         │
│      └─▶ data/aibridge.db → Settings WHERE Key = 'Rag.ChunkSize'           │
│                                                                              │
│   3. Seed Default (Lowest Priority)                                          │
│      └─▶ config.seed.json → Categories.Rag.Settings.ChunkSize.Default      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Use Case Guidelines

### Recommended Settings by Content Type

| Content Type | ChunkSize | ChunkOverlap | Rationale |
|--------------|-----------|--------------|-----------|
| **Code** | 4096-8192 | 200-300 | Larger context for function boundaries |
| **Documentation** | 1024-2048 | 100-150 | Balanced for paragraphs |
| **Dense Technical** | 256-512 | 50-100 | Smaller for precise retrieval |
| **Long-form Articles** | 2048-4096 | 150-200 | Maintain narrative flow |
| **API References** | 512-1024 | 50-100 | Endpoint-level granularity |

---

## Database Schema

### Settings Table (Root DB)

```sql
-- In aibridge.db
INSERT INTO Settings (Key, Value, ValueType, Source, Description) VALUES
('Rag.ChunkSize', '2048', 'int', 'seed', 'Tokens per RAG chunk'),
('Rag.ChunkSizeMin', '256', 'int', 'seed', 'Minimum allowed chunk size'),
('Rag.ChunkSizeMax', '8192', 'int', 'seed', 'Maximum allowed chunk size'),
('Rag.ChunkOverlap', '100', 'int', 'seed', 'Token overlap between chunks'),
('Rag.ChunkOverlapMin', '0', 'int', 'seed', 'Minimum overlap'),
('Rag.ChunkOverlapMax', '512', 'int', 'seed', 'Maximum overlap'),
('Rag.ContextTokenBudget', '4096', 'int', 'seed', 'Max tokens for context injection'),
('Rag.EmbeddingModel', 'NomicEmbedText', 'string', 'seed', 'Vector embedding model'),
('Rag.SimilarityThreshold', '0.7', 'float', 'seed', 'Minimum similarity score'),
('Rag.TopK', '10', 'int', 'seed', 'Max chunks per search');
```

### App-Level Override Table

```sql
-- In data/{appName}/settings/config.db
CREATE TABLE IF NOT EXISTS Setting (
    SettingsId INTEGER PRIMARY KEY AUTOINCREMENT,
    Key TEXT UNIQUE NOT NULL,
    Value TEXT NOT NULL,
    ValueType TEXT NOT NULL CHECK(ValueType IN ('string', 'int', 'float', 'bool', 'json')),
    Source TEXT NOT NULL DEFAULT 'user',
    Description TEXT,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Example override for code-heavy application
INSERT INTO Settings (Key, Value, ValueType, Source, Description) VALUES
('Rag.ChunkSize', '4096', 'int', 'user', 'Larger chunks for code context');
```

---

## API Endpoints

### Get Current Settings

```
GET /api/v1/settings/rag

Response:
{
  "ChunkSize": 2048,
  "ChunkSizeMin": 256,
  "ChunkSizeMax": 8192,
  "ChunkOverlap": 100,
  "ContextTokenBudget": 4096,
  "EmbeddingModel": "NomicEmbedText",
  "SimilarityThreshold": 0.7,
  "TopK": 10,
  "Source": "seed"
}
```

### Update Settings

```
PUT /api/v1/settings/rag
{
  "ChunkSize": 4096,
  "ChunkOverlap": 200
}

Response:
{
  "Updated": ["ChunkSize", "ChunkOverlap"],
  "Validated": true,
  "Source": "user"
}
```

### Get App-Level Override

```
GET /api/v1/apps/:appName/settings/rag

Response:
{
  "ChunkSize": 4096,
  "ChunkOverlap": 200,
  "Source": "app",
  "InheritedFrom": {
    "ContextTokenBudget": "root",
    "EmbeddingModel": "seed"
  }
}
```

---

## Migration from Previous Defaults

If upgrading from older chunk defaults (512/50), the system will:

1. **Preserve existing chunks** - Do not re-chunk on setting change
2. **Apply new settings to new ingestions only**
3. **Offer re-indexing** via `POST /api/v1/rag/reindex`

```
POST /api/v1/rag/reindex
{
  "AppName": "myapp",
  "ApplyNewChunkSettings": true
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Seedable Config Overview | `../00-overview.md` |
| AI Bridge Database | `../22-ai-bridge-cli/01-backend/12-database-architecture.md` |
| RAG Reindexing | `../22-ai-bridge-cli/01-backend/11-rag-reindexing.md` |
| Error Codes | `../03-error-manage/03-error-code-registry/00-overview.md` |
