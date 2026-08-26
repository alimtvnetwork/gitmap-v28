# Consolidated: Seedable Config Architecture — Complete Reference

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for version-controlled configuration seeding. An AI reading only this file must be able to implement the full pattern — including seed files, JSON schema validation, merge strategy, version bumping, changelog generation, validation data seeding, RAG config, theme system, and typed Go accessors — without consulting source specs.

**Source:** `spec/06-seedable-config-architecture/` (01-fundamentals + 02-features/01-05)

---

## CRITICAL: Naming Convention

**All field names use PascalCase. No underscores allowed.**

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `chunk_size` | `ChunkSize` |
| `created_at` | `CreatedAt` |
| `seed_version` | `SeedVersion` |

---

## 1. Pattern Summary

A reusable pattern for managing application configuration where:

1. **First-run seeding** populates SQLite DB from `config.seed.json`
2. **Every config change updates the version** (SemVer)
3. **Every version change logs to `changelog.md`**
4. **Subsequent runs respect version** to avoid duplicate seeds
5. **User customizations are never overwritten** by seed updates

### Version Flow

```
config.seed.json (source of truth)
        │
        ▼
  Version Change Detected?
        │
   ┌────┴────┐
   │         │
   NO        YES
   │         │
 Skip     Merge new settings (preserve existing)
           │
           ▼
      Update changelog.md
```

---

## 2. File Specifications

### config.seed.json

```json
{
  "$schema": "./config.schema.json",
  "Version": "1.2.0",
  "Changelog": "Added new cache settings for improved performance",
  "Categories": {
    "General": {
      "DisplayName": "General",
      "Description": "General application settings",
      "Settings": {
        "Theme": {
          "Type": "select",
          "Label": "Theme",
          "Description": "Application color theme",
          "Default": "system",
          "Options": ["light", "dark", "system", "high-contrast"]
        },
        "Language": {
          "Type": "select",
          "Label": "Language",
          "Default": "en",
          "Options": ["en", "es", "fr", "de", "zh", "ja"]
        },
        "AutoSave": {
          "Type": "boolean",
          "Label": "Auto Save",
          "Description": "Automatically save changes",
          "Default": true
        }
      }
    },
    "Cache": {
      "DisplayName": "Cache",
      "Description": "Caching configuration",
      "AddedIn": "1.2.0",
      "Settings": {
        "Enabled": {
          "Type": "boolean",
          "Label": "Enable Cache",
          "Default": true
        },
        "MaxSizeMb": {
          "Type": "number",
          "Label": "Max Cache Size (MB)",
          "Default": 100,
          "Min": 10,
          "Max": 1000
        },
        "TtlHours": {
          "Type": "number",
          "Label": "Cache TTL (hours)",
          "Default": 24,
          "Min": 1,
          "Max": 168
        }
      }
    }
  }
}
```

### config.schema.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Application Configuration",
  "type": "object",
  "required": ["Version", "Categories"],
  "properties": {
    "Version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+\\.\\d+$",
      "description": "Semantic version of configuration"
    },
    "Changelog": {
      "type": "string",
      "description": "Description of changes in this version"
    },
    "Categories": {
      "type": "object",
      "additionalProperties": { "$ref": "#/definitions/category" }
    }
  },
  "definitions": {
    "category": {
      "type": "object",
      "required": ["DisplayName", "Settings"],
      "properties": {
        "DisplayName": { "type": "string" },
        "Description": { "type": "string" },
        "Version": { "type": "string" },
        "AddedIn": { "type": "string" },
        "Settings": {
          "type": "object",
          "additionalProperties": { "$ref": "#/definitions/setting" }
        }
      }
    },
    "setting": {
      "type": "object",
      "required": ["Type", "Label", "Default"],
      "properties": {
        "Type": {
          "type": "string",
          "enum": ["string", "number", "boolean", "select", "array", "object"]
        },
        "Label": { "type": "string" },
        "Description": { "type": "string" },
        "Default": {},
        "Min": { "type": "number" },
        "Max": { "type": "number" },
        "Options": { "type": "array" },
        "AddedIn": { "type": "string" },
        "DeprecatedIn": { "type": "string" }
      }
    }
  }
}
```

### Setting Types

| Type | JSON Value | Example |
|------|-----------|---------|
| `string` | `"value"` | Free-text input |
| `number` | `42` | Numeric with optional Min/Max |
| `boolean` | `true` | Toggle switch |
| `select` | `"option"` | Dropdown with `Options` array |
| `array` | `["a", "b"]` | Multi-value (e.g., word lists) |
| `object` | `{...}` | Nested config |

---

## 3. Database Schema

### Table: ConfigMeta

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated seedable-config example; canonical version in 06-seedable-config-architecture/"
CREATE TABLE ConfigMeta (
    ConfigMetaId INTEGER PRIMARY KEY AUTOINCREMENT,
    SeedVersion TEXT NOT NULL,
    CurrentVersion TEXT NOT NULL,
    LastSeededAt DATETIME,
    ChangelogUpdatedAt DATETIME,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Table: Setting

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated seedable-config example; canonical version in 06-seedable-config-architecture/"
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,
    Key TEXT NOT NULL,
    Value TEXT NOT NULL,              -- JSON encoded
    Type TEXT NOT NULL,
    AddedInVersion TEXT,
    ModifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(Category, Key)
);

CREATE INDEX IdxSettingsCategory ON Setting(Category);
```

### Table: SettingHistory

```sql
CREATE TABLE SettingHistory (
    SettingsHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
    SettingId INTEGER NOT NULL,
    OldValue TEXT,
    NewValue TEXT NOT NULL,
    ChangedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ChangedBy TEXT,                   -- user, system, seed
    Version TEXT,
    FOREIGN KEY (SettingId) REFERENCES Setting(SettingId)
);

CREATE INDEX IdxHistorySetting ON SettingHistory(SettingId);
CREATE INDEX IdxHistoryChanged ON SettingHistory(ChangedAt);
```

### Table: ValidationData

```sql
-- linter-waive: MISSING-DESC-001 reason="Consolidated seedable-config example; canonical version in 06-seedable-config-architecture/"
CREATE TABLE ValidationData (
    ValidationDataId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,          -- 'Seo', 'Rag', 'Search'
    Key TEXT NOT NULL,               -- 'TransitionWords', 'StopWords'
    DataType TEXT NOT NULL,          -- 'array', 'map', 'number'
    Value TEXT NOT NULL,             -- JSON encoded
    Version TEXT NOT NULL,           -- Seed version
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(Category, Key)
);

CREATE INDEX IdxValidationDataCategory ON ValidationData(Category);
```

---

## 4. Strongly-Typed Value Container

**CRITICAL: No `interface{}` or `any` usage.** All values use a union struct:

```go
type SettingValue struct {
    StringVal  *string           `json:"StringVal,omitempty"`
    IntVal     *int              `json:"IntVal,omitempty"`
    FloatVal   *float64          `json:"FloatVal,omitempty"`
    BoolVal    *bool             `json:"BoolVal,omitempty"`
    StringsVal []string          `json:"StringsVal,omitempty"`
    MapVal     map[string]string `json:"MapVal,omitempty"`
}
```

---

## 5. Go Implementation

### ConfigService Struct & Types

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/Masterminds/semver/v3"
    "gorm.io/gorm"
)

type ConfigService struct {
    db            *gorm.DB
    seedPath      string
    changelogPath string
}

type SeedConfig struct {
    Version    string
    Changelog  string                    `json:",omitempty"`
    Categories map[string]CategoryConfig
}

type CategoryConfig struct {
    DisplayName string
    Description string                   `json:",omitempty"`
    Version     string                   `json:",omitempty"`
    AddedIn     string                   `json:",omitempty"`
    Settings    map[string]SettingConfig
}

type SettingConfig struct {
    Type        string
    Label       string
    Description string        `json:",omitempty"`
    Default     SettingValue  // Strongly typed — no interface{}
    Min         *float64      `json:",omitempty"`
    Max         *float64      `json:",omitempty"`
    Options     []string      `json:",omitempty"`
    AddedIn     string        `json:",omitempty"`
}

// ConfigMetaUpdate is the typed struct for GORM .Updates() — no map[string]interface{}
type ConfigMetaUpdate struct {
    SeedVersion    string    `gorm:"column:SeedVersion"`
    CurrentVersion string    `gorm:"column:CurrentVersion"`
    LastSeededAt   time.Time `gorm:"column:LastSeededAt"`
    UpdatedAt      time.Time `gorm:"column:UpdatedAt"`
}
```

### Seed With Version Check

```go
func (s *ConfigService) SeedWithVersionCheck() error {
    seed, err := s.loadSeedFile()
    if err != nil {
        return apperror.Wrap(err, ErrSeedLoadFailed, "load seed file")
    }

    meta, err := s.getMeta()
    if err != nil {
        return s.fullSeed(seed) // First time — full seed
    }

    currentVer, _ := semver.NewVersion(meta.SeedVersion)
    seedVer, _ := semver.NewVersion(seed.Version)

    if !seedVer.GreaterThan(currentVer) {
        return nil // No version change — skip
    }

    if err := s.mergeSeed(seed, meta.SeedVersion); err != nil {
        return err
    }

    return s.updateChangelog(seed)
}
```

### Merge Strategy

**New settings are inserted with defaults. Existing settings are never overwritten** — user customizations are preserved.

```go
func (s *ConfigService) mergeSeed(seed SeedConfig, previousVersion string) error {
    for catKey, cat := range seed.Categories {
        for settingKey, setting := range cat.Settings {
            var existing Setting
            err := s.db.Where("Category = ? AND Key = ?", catKey, settingKey).First(&existing).Error

            if err == gorm.ErrRecordNotFound {
                valueJson, _ := json.Marshal(setting.Default)
                s.db.Create(&Setting{
                    Category:       catKey,
                    Key:            settingKey,
                    Value:          string(valueJson),
                    Type:           setting.Type,
                    AddedInVersion: seed.Version,
                })
            }
            // Existing settings are preserved
        }
    }

    s.db.Model(&ConfigMeta{}).Where("ConfigMetaId = 1").Updates(ConfigMetaUpdate{
        SeedVersion:    seed.Version,
        CurrentVersion: seed.Version,
        LastSeededAt:   time.Now(),
        UpdatedAt:      time.Now(),
    })

    return nil
}
```

### Changelog Update

```go
func (s *ConfigService) updateChangelog(seed SeedConfig) error {
    if seed.Changelog == "" {
        return nil
    }

    entry := fmt.Sprintf("\n## [%s] - %s\n\n%s\n",
        seed.Version,
        time.Now().Format("2006-01-02"),
        seed.Changelog,
    )

    contentResult := pathutil.ReadFileIfExists(s.changelogPath)
    if contentResult.IsErr() {
        return contentResult.Err()
    }
    content := contentResult.Value()

    header := "# Changelog\n\nAll notable configuration changes are documented here.\n"
    if len(content) == 0 {
        content = []byte(header)
    }

    insertPos := len(header)
    newContent := string(content[:insertPos]) + entry + string(content[insertPos:])

    return pathutil.WriteFile(s.changelogPath, []byte(newContent))
}
```

---

## 6. Version Bumping Rules

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| New category added | Minor | 1.0.0 → 1.1.0 |
| New setting added | Minor | 1.1.0 → 1.2.0 |
| Default value changed | Patch | 1.2.0 → 1.2.1 |
| Setting deprecated | Patch | 1.2.1 → 1.2.2 |
| Breaking change (setting removed) | Major | 1.2.2 → 2.0.0 |

---

## 7. changelog.md Format

```markdown

# Changelog

All notable configuration changes are documented here.

## [1.2.0] - 2026-02-01

### Added

- Cache category with Enabled, MaxSizeMb, TtlHours settings

### Changed

- Theme options now include "high-contrast"

## [1.0.0] - 2026-01-01

### Initial Release

- General category with Theme, Language, AutoSave
```

---

## 8. Validation Data Seeding

### CRITICAL: No Hardcoded Arrays

**All validation arrays, lookup tables, and configurable data MUST use the CW Config → Root DB pattern.** Never hardcode validation data in Go source code.

```go
// ❌ WRONG: Hardcoded validation data
transitions := []string{"however", "therefore", "additionally"}

// ✅ CORRECT: Load from Root DB via typed accessor
transitions, err := v.validationData.GetSeoStringArray(SeoKeyTransitionWords)
```

### Typed Constants (Mandatory)

Never use magic strings for category/key lookups:

```go
type ValidationCategory string

const (
    CategorySeo    ValidationCategory = "Seo"
    CategoryRag    ValidationCategory = "Rag"
    CategoryFaq    ValidationCategory = "Faq"
    CategorySearch ValidationCategory = "Search"
)

type SeoKey string

const (
    SeoKeyTransitionWords            SeoKey = "TransitionWords"
    SeoKeyTransitionDensityThreshold SeoKey = "TransitionDensityThreshold"
    SeoKeyMaxSentenceWords           SeoKey = "MaxSentenceWords"
    SeoKeyMaxParagraphWords          SeoKey = "MaxParagraphWords"
    SeoKeyMinKeywordMentions         SeoKey = "MinKeywordMentions"
    SeoKeyForbiddenContainerTags     SeoKey = "ForbiddenContainerTags"
    SeoKeySlugMinWords               SeoKey = "SlugMinWords"
    SeoKeySlugMaxWords               SeoKey = "SlugMaxWords"
)

type RagKey string

const (
    RagKeyStopWords    RagKey = "StopWords"
    RagKeyMinChunkSize RagKey = "MinChunkSize"
    RagKeyMaxChunkSize RagKey = "MaxChunkSize"
    RagKeyChunkOverlap RagKey = "ChunkOverlap"
)
```

### ValidationDataService

```go
type ValidationDataService struct {
    db    *gorm.DB
    cache sync.Map // Thread-safe cache
}

func (s *ValidationDataService) GetStringArray(
    category ValidationCategory,
    key string,
) apperror.Result[[]string] {
    cacheKey := string(category) + ":" + key
    if cached, ok := s.cache.Load(cacheKey); ok {
        return cached.([]string), nil
    }

    var data ValidationData
    if err := s.db.Where("Category = ? AND Key = ?", string(category), key).First(&data).Error; err != nil {
        return nil, err
    }

    var result []string
    json.Unmarshal([]byte(data.Value), &result)
    s.cache.Store(cacheKey, result)

    return result, nil
}

func (s *ValidationDataService) GetNumber(
    category ValidationCategory,
    key string,
) apperror.Result[float64] {
    cacheKey := string(category) + ":" + key
    if cached, ok := s.cache.Load(cacheKey); ok {
        return cached.(float64), nil
    }

    var data ValidationData
    if err := s.db.Where("Category = ? AND Key = ?", string(category), key).First(&data).Error; err != nil {
        return 0, err
    }

    var result float64
    json.Unmarshal([]byte(data.Value), &result)
    s.cache.Store(cacheKey, result)

    return result, nil
}

// Domain-specific typed accessors
func (s *ValidationDataService) GetSeoStringArray(key SeoKey) apperror.Result[[]string] {
    return s.GetStringArray(CategorySeo, string(key))
}

func (s *ValidationDataService) GetSeoNumber(key SeoKey) apperror.Result[float64] {
    return s.GetNumber(CategorySeo, string(key))
}

func (s *ValidationDataService) GetRagStringArray(key RagKey) apperror.Result[[]string] {
    return s.GetStringArray(CategoryRag, string(key))
}

func (s *ValidationDataService) InvalidateCache() {
    s.cache = sync.Map{}
}
```

---

## 9. RAG Chunk Configuration

### Settings

| Setting Key | Type | Default | Min | Max | Description |
|-------------|------|---------|-----|-----|-------------|
| `Rag.ChunkSize` | int | 2048 | 256 | 8192 | Tokens per chunk (must be multiple of 256) |
| `Rag.ChunkOverlap` | int | 100 | 0 | 512 | Overlap between chunks (max 25% of ChunkSize) |
| `Rag.ContextTokenBudget` | int | 4096 | 512 | 16384 | Max tokens for RAG context |
| `Rag.EmbeddingModel` | string | `NomicEmbedText` | — | — | Embedding model |
| `Rag.SimilarityThreshold` | float | 0.7 | 0.0 | 1.0 | Min cosine similarity |
| `Rag.TopK` | int | 10 | 1 | 50 | Max chunks returned |

### Configuration Priority (Highest First)

1. **App-Level Override** — `data/{appName}/settings/config.db`
2. **Root DB Setting** — `data/root.db → Setting table`
3. **Seed Default** — `config.seed.json`

### Recommended Settings by Content Type

| Content Type | ChunkSize | ChunkOverlap | Rationale |
|--------------|-----------|--------------|-----------|
| Code | 4096–8192 | 200–300 | Larger context for function boundaries |
| Documentation | 1024–2048 | 100–150 | Balanced for paragraphs |
| Dense Technical | 256–512 | 50–100 | Smaller for precise retrieval |
| Long-form Articles | 2048–4096 | 150–200 | Maintain narrative flow |
| API References | 512–1024 | 50–100 | Endpoint-level granularity |

### API Endpoints

```
GET  /api/v1/settings/rag           → Current effective settings (merged priority)
PUT  /api/v1/settings/rag           → Update root-level settings
GET  /api/v1/apps/:appName/settings/rag  → App-level with inheritance info
POST /api/v1/rag/reindex            → Re-chunk with new settings
```

---

## 10. Theme System

The seed pattern supports a comprehensive theme system via the `Appearance` category:

| Setting | Type | Default | Options |
|---------|------|---------|---------|
| Theme | select | system | light, dark, system, high-contrast, dracula, monokai, solarized-light/dark, nord-light/dark, ocean-blue/dark, github-light/dark, etc. |
| AccentColor | select | blue | blue, indigo, violet, purple, pink, rose, red, orange, amber, green, teal, cyan, etc. |
| FontSize | select | medium | x-small, small, medium, large, x-large |
| FontFamily | select | system | system, inter, roboto, jetbrains-mono, fira-code, etc. |
| BorderRadius | select | medium | none, small, medium, large, full |
| AnimationSpeed | select | normal | none, reduced, normal, fast |
| CompactMode | boolean | false | — |

### Theme CSS Variable Mapping

```css
[data-theme="ocean-blue"] {
  --background: 200 30% 98%;
  --foreground: 200 50% 10%;
  --primary: 200 80% 50%;
  --primary-foreground: 200 10% 98%;
  --secondary: 180 40% 90%;
  --accent: 180 60% 45%;
}
```

### React Theme Hook

```typescript
export function useTheme() {
  const { settings, updateSetting } = useSettings();
  const theme = settings?.Appearance?.Theme ?? 'system';

  const setTheme = (newTheme: string) => {
    updateSetting('Appearance', 'Theme', newTheme);
    document.documentElement.setAttribute('data-theme', newTheme);
  };

  return { theme, setTheme };
}
```

---

## 11. UI Integration

### Version Badge Component

```typescript
export function VersionBadge() {
  const { meta } = useConfig();
  const isNew = meta.SeedVersion !== meta.CurrentVersion;

  return (
    <Badge variant={isNew ? "default" : "secondary"}>
      v{meta.CurrentVersion}
      {isNew && " (updated)"}
    </Badge>
  );
}
```

### New Settings Highlight

```typescript
function SettingItem({ setting, currentVersion }: Props) {
  const isNew = setting.AddedInVersion === currentVersion;

  return (
    <div className={cn("p-4 rounded-lg", isNew && "ring-2 ring-primary bg-primary/5")}>
      {isNew && <Badge className="mb-2">New in v{currentVersion}</Badge>}
      {/* setting content */}
    </div>
  );
}
```

---

## 12. Applicable Projects

| Project | Config Location |
|---------|-----------------|
| Spec Management | `go-backend/configs/` |
| GSearch CLI | `backend/configs/` |
| BRun CLI | `backend/configs/` |
| AI Bridge | `backend/configs/` |
| Nexus Flow | `backend/configs/` |
| WP SEO Publish CLI | `backend/configs/` |

---

## 13. Key Rules Summary

| # | Rule |
|---|------|
| 1 | **PascalCase fields** — zero tolerance for underscores in JSON or DB columns |
| 2 | **Typed SettingValue** — no `interface{}` or `any` — use union struct |
| 3 | **SemVer gating** — seed only runs when version increases |
| 4 | **Merge, don't overwrite** — user customizations are always preserved |
| 5 | **Changelog auto-updated** — every version bump appends to changelog.md |
| 6 | **JSON Schema validation** — all seed files validated against `config.schema.json` |
| 7 | **No hardcoded arrays** — all validation data seeded via CW Config → Root DB |
| 8 | **Typed constants** — never use magic strings for category/key lookups |
| 9 | **Typed GORM updates** — use `ConfigMetaUpdate` struct, not `map[string]interface{}` |
| 10 | **Cache with invalidation** — `sync.Map` cache in ValidationDataService, explicit invalidate |
| 11 | **3-tier priority** — App override > Root DB > Seed default |

---

## Dependencies

- Depends on [Split DB Architecture](./05-split-db-architecture.md) for the database layer
- Uses `github.com/Masterminds/semver/v3` for version comparison
- Uses GORM for database operations

---

*Consolidated seedable config — v3.2.0 — 2026-04-16*
