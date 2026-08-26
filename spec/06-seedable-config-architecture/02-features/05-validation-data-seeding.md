# Validation Data Seeding Pattern

**Version:** 3.2.0  
**Created:** 2026-03-09  
**Status:** Active  
**Purpose:** Define pattern for loading validation arrays and lookup data from CW Config → Root DB

---

## Overview

All **validation arrays**, **lookup tables**, and **configurable data** used across CLI applications MUST follow the CW Config → Root DB pattern. This ensures:

1. **No hardcoded arrays** in Go source code
2. **Runtime configurability** via settings database
3. **Version-controlled changes** through seed versioning
4. **User customization** without code changes

---

## Anti-Pattern: Hardcoded Arrays ❌

```go
// ❌ WRONG: Hardcoded validation data
func validateTransitionDensity(content string, _ *SeoConfig) ValidationResult {
    transitions := []string{
        "however", "therefore", "additionally", "moreover", "furthermore",
        "consequently", "meanwhile", "nevertheless", "accordingly", "hence",
    }
    // ...
}
```

---

## Correct Pattern: CW Config → Root DB ✅

### Step 1: Define in config.seed.json

```json
{
  "$schema": "./config.schema.json",
  "Version": "1.3.0",
  "Changelog": "Added SEO validation data arrays",
  "Categories": {
    "Seo": {
      "DisplayName": "SEO Settings",
      "Description": "SEO content generation configuration",
      "Settings": {
        "TransitionWords": {
          "Type": "array",
          "Label": "Transition Words",
          "Description": "Words counted for transition density validation",
          "Default": [
            "however", "therefore", "additionally", "moreover", "furthermore",
            "consequently", "meanwhile", "nevertheless", "accordingly", "hence",
            "thus", "indeed", "specifically", "particularly", "notably",
            "significantly", "ultimately", "essentially", "primarily", "initially",
            "subsequently", "similarly", "likewise", "conversely", "alternatively",
            "otherwise", "regardless", "nonetheless", "certainly", "undoubtedly"
          ]
        },
        "TransitionDensityThreshold": {
          "Type": "number",
          "Label": "Transition Density Threshold",
          "Description": "Minimum percentage of transition words required",
          "Default": 40,
          "Min": 10,
          "Max": 80
        },
        "MaxSentenceWords": {
          "Type": "number",
          "Label": "Max Sentence Words",
          "Description": "Maximum words allowed per sentence",
          "Default": 18,
          "Min": 10,
          "Max": 50
        },
        "MaxParagraphWords": {
          "Type": "number",
          "Label": "Max Paragraph Words",
          "Description": "Maximum words allowed per paragraph",
          "Default": 180,
          "Min": 50,
          "Max": 500
        },
        "MinKeywordMentions": {
          "Type": "number",
          "Label": "Min Keyword Mentions",
          "Description": "Minimum keyword occurrences required",
          "Default": 8,
          "Min": 3,
          "Max": 20
        },
        "MinAreaMentions": {
          "Type": "number",
          "Label": "Min Area Mentions",
          "Description": "Minimum area/location mentions per section",
          "Default": 3,
          "Min": 1,
          "Max": 10
        },
        "MaxAreaMentions": {
          "Type": "number",
          "Label": "Max Area Mentions",
          "Description": "Maximum area/location mentions per section",
          "Default": 4,
          "Min": 2,
          "Max": 15
        },
        "LinksPerSentenceMin": {
          "Type": "number",
          "Label": "Min Links Per Sentence",
          "Description": "Minimum internal/external links per sentence",
          "Default": 2,
          "Min": 0,
          "Max": 5
        },
        "LinksPerSentenceMax": {
          "Type": "number",
          "Label": "Max Links Per Sentence",
          "Description": "Maximum internal/external links per sentence",
          "Default": 3,
          "Min": 1,
          "Max": 10
        },
        "StatisticalRangeMin": {
          "Type": "number",
          "Label": "Statistical Range Min",
          "Description": "Minimum value for credibility percentages",
          "Default": 2.51,
          "Min": 0.01,
          "Max": 10.0
        },
        "StatisticalRangeMax": {
          "Type": "number",
          "Label": "Statistical Range Max",
          "Description": "Maximum value for credibility percentages",
          "Default": 2.97,
          "Min": 0.01,
          "Max": 10.0
        },
        "TrustMetricsMax": {
          "Type": "number",
          "Label": "Trust Metrics Max",
          "Description": "Maximum percentage for company glorification metrics",
          "Default": 5,
          "Min": 1,
          "Max": 10
        },
        "ForbiddenContainerTags": {
          "Type": "array",
          "Label": "Forbidden Container Tags",
          "Description": "HTML tags not allowed inside seo-container-para contrast",
          "Default": ["p", "div"]
        },
        "ExperienceYearsMin": {
          "Type": "number",
          "Label": "Experience Years Min",
          "Description": "Minimum years for experience narratives",
          "Default": 5,
          "Min": 1,
          "Max": 20
        },
        "ExperienceYearsMax": {
          "Type": "number",
          "Label": "Experience Years Max",
          "Description": "Maximum years for experience narratives",
          "Default": 15,
          "Min": 5,
          "Max": 50
        },
        "SlugMaxWords": {
          "Type": "number",
          "Label": "Slug Max Words",
          "Description": "Maximum words in generated URL slugs",
          "Default": 4,
          "Min": 2,
          "Max": 8
        },
        "SlugMinWords": {
          "Type": "number",
          "Label": "Slug Min Words",
          "Description": "Minimum words in generated URL slugs",
          "Default": 3,
          "Min": 1,
          "Max": 5
        },
        "ExternalLinkNofollowCount": {
          "Type": "number",
          "Label": "External Nofollow Count",
          "Description": "Number of external links to mark as nofollow",
          "Default": 2,
          "Min": 0,
          "Max": 5
        }
      }
    }
  }
}
```

### Step 2: Database Storage

```sql
-- Root DB (settings.db) stores validation data
-- linter-waive: MISSING-DESC-001 reason="Seedable-config example; focus on config schema mechanics"
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

### Step 3: Go Enums/Constants (MANDATORY)

> **CRITICAL:** Never use magic strings. Always use typed constants for categories and keys.

```go
package validation

// ============================================
// Validation Category Constants
// ============================================
type ValidationCategory string

const (
    CategorySeo    ValidationCategory = "Seo"
    CategoryRag    ValidationCategory = "Rag"
    CategoryFaq    ValidationCategory = "Faq"
    CategorySearch ValidationCategory = "Search"
)

// ============================================
// SEO Validation Key Constants
// ============================================
type SeoKey string

const (
    SeoKeyTransitionWords           SeoKey = "TransitionWords"
    SeoKeyTransitionDensityThreshold SeoKey = "TransitionDensityThreshold"
    SeoKeyMaxSentenceWords          SeoKey = "MaxSentenceWords"
    SeoKeyMaxParagraphWords         SeoKey = "MaxParagraphWords"
    SeoKeyMinKeywordMentions        SeoKey = "MinKeywordMentions"
    SeoKeyMinAreaMentions           SeoKey = "MinAreaMentions"
    SeoKeyMaxAreaMentions           SeoKey = "MaxAreaMentions"
    SeoKeyLinksPerSentenceMin       SeoKey = "LinksPerSentenceMin"
    SeoKeyLinksPerSentenceMax       SeoKey = "LinksPerSentenceMax"
    SeoKeyStatisticalRangeMin       SeoKey = "StatisticalRangeMin"
    SeoKeyStatisticalRangeMax       SeoKey = "StatisticalRangeMax"
    SeoKeyTrustMetricsMax           SeoKey = "TrustMetricsMax"
    SeoKeyForbiddenContainerTags    SeoKey = "ForbiddenContainerTags"
    SeoKeyExperienceYearsMin        SeoKey = "ExperienceYearsMin"
    SeoKeyExperienceYearsMax        SeoKey = "ExperienceYearsMax"
    SeoKeySlugMinWords              SeoKey = "SlugMinWords"
    SeoKeySlugMaxWords              SeoKey = "SlugMaxWords"
    SeoKeyExternalLinkNofollowCount SeoKey = "ExternalLinkNofollowCount"
)

// ============================================
// RAG Validation Key Constants
// ============================================
type RagKey string

const (
    RagKeyStopWords    RagKey = "StopWords"
    RagKeyMinChunkSize RagKey = "MinChunkSize"
    RagKeyMaxChunkSize RagKey = "MaxChunkSize"
    RagKeyChunkOverlap RagKey = "ChunkOverlap"
)

// ============================================
// FAQ Validation Key Constants
// ============================================
type FaqKey string

const (
    FaqKeyDefaultOutputFormat       FaqKey = "DefaultOutputFormat"
    FaqKeyDefaultIncludeSchema      FaqKey = "DefaultIncludeSchema"
    FaqKeyDefaultSchemaVariation    FaqKey = "DefaultSchemaVariation"
    FaqKeyDefaultEncodeHtmlInJson   FaqKey = "DefaultEncodeHtmlInJson"
    FaqKeyDefaultWordLimit          FaqKey = "DefaultWordLimit"
    FaqKeyDefaultTransitionDensity  FaqKey = "DefaultTransitionDensity"
    FaqKeyDefaultKeywordMentions    FaqKey = "DefaultKeywordMentions"
    FaqKeyDefaultAreaMentions       FaqKey = "DefaultAreaMentions"
    FaqKeyDefaultMaxSentenceWords   FaqKey = "DefaultMaxSentenceWords"
    FaqKeyDefaultMaxParagraphWords  FaqKey = "DefaultMaxParagraphWords"
    FaqKeyDefaultEnableGSearch      FaqKey = "DefaultEnableGSearch"
    FaqKeyDefaultEnableSitemapLinking FaqKey = "DefaultEnableSitemapLinking"
    FaqKeyDefaultEnableYouTubeEmbed FaqKey = "DefaultEnableYouTubeEmbed"
    FaqKeySchemaParagraphs          FaqKey = "SchemaParagraphs"
    FaqKeySchemaTemplates           FaqKey = "SchemaTemplates"
    FaqKeyHtmlTemplates             FaqKey = "HtmlTemplates"
    FaqKeyTransitionWords           FaqKey = "TransitionWords"
    FaqKeyTrustPercentageMin        FaqKey = "TrustPercentageMin"
    FaqKeyTrustPercentageMax        FaqKey = "TrustPercentageMax"
    FaqKeyMonthlyImprovementMin     FaqKey = "MonthlyImprovementMin"
    FaqKeyMonthlyImprovementMax     FaqKey = "MonthlyImprovementMax"
    FaqKeyEffectivenessMin          FaqKey = "EffectivenessMin"
    FaqKeyEffectivenessMax          FaqKey = "EffectivenessMax"
    FaqKeyQuestionPatterns          FaqKey = "QuestionPatterns"
)

// ============================================
// Search Validation Key Constants
// ============================================
type SearchKey string

const (
    SearchKeyAllowedFileTypes     SearchKey = "AllowedFileTypes"
    SearchKeyExcludedDirectories  SearchKey = "ExcludedDirectories"
    SearchKeyMaxFileSize          SearchKey = "MaxFileSize"
)
```

### Step 4: ValidationDataService with Typed Methods

```go
package validation

import (
    "encoding/json"
    "sync"
)

// ValidationDataService loads validation data from Root DB
type ValidationDataService struct {
    db    *gorm.DB
    cache sync.Map  // Thread-safe cache
}

// GetStringArray retrieves a string array using typed constants
func (s *ValidationDataService) GetStringArray(category ValidationCategory, key string) apperror.Result[[]string] {
    cacheKey := string(category) + ":" + key
    // EXEMPTED: typed accessor internal — cache stores known []string values (§7.2)
    if cached, ok := s.cache.Load(cacheKey); ok {
        return cached.([]string), nil
    }
    
    var data ValidationData
    if err := s.db.Where("Category = ? AND Key = ?", string(category), key).First(&data).Error; err != nil {
        return nil, err
    }
    
    var result []string
    if err := json.Unmarshal([]byte(data.Value), &result); err != nil {
        return nil, err
    }
    
    s.cache.Store(cacheKey, result)
    return result, nil
}

// GetNumber retrieves a numeric value using typed constants
func (s *ValidationDataService) GetNumber(category ValidationCategory, key string) apperror.Result[float64] {
    cacheKey := string(category) + ":" + key
    // EXEMPTED: typed accessor internal — cache stores known float64 values (§7.2)
    if cached, ok := s.cache.Load(cacheKey); ok {
        return cached.(float64), nil
    }
    
    var data ValidationData
    if err := s.db.Where("Category = ? AND Key = ?", string(category), key).First(&data).Error; err != nil {
        return 0, err
    }
    
    var result float64
    if err := json.Unmarshal([]byte(data.Value), &result); err != nil {
        return 0, err
    }
    
    s.cache.Store(cacheKey, result)
    return result, nil
}

// SEO-specific typed accessors
func (s *ValidationDataService) GetSeoStringArray(key SeoKey) apperror.Result[[]string] {
    return s.GetStringArray(CategorySeo, string(key))
}

func (s *ValidationDataService) GetSeoNumber(key SeoKey) apperror.Result[float64] {
    return s.GetNumber(CategorySeo, string(key))
}

// RAG-specific typed accessors
func (s *ValidationDataService) GetRagStringArray(key RagKey) apperror.Result[[]string] {
    return s.GetStringArray(CategoryRag, string(key))
}

func (s *ValidationDataService) GetRagNumber(key RagKey) apperror.Result[float64] {
    return s.GetNumber(CategoryRag, string(key))
}

// FAQ-specific typed accessors
func (s *ValidationDataService) GetFaqStringArray(key FaqKey) apperror.Result[[]string] {
    return s.GetStringArray(CategoryFaq, string(key))
}

func (s *ValidationDataService) GetFaqNumber(key FaqKey) apperror.Result[float64] {
    return s.GetNumber(CategoryFaq, string(key))
}

func (s *ValidationDataService) GetFaqBool(key FaqKey) apperror.Result[bool] {
    return s.GetBool(CategoryFaq, string(key))
}

func (s *ValidationDataService) GetFaqString(key FaqKey) apperror.Result[string] {
    return s.GetString(CategoryFaq, string(key))
}

// Search-specific typed accessors
func (s *ValidationDataService) GetSearchStringArray(key SearchKey) apperror.Result[[]string] {
    return s.GetStringArray(CategorySearch, string(key))
}

func (s *ValidationDataService) GetSearchNumber(key SearchKey) apperror.Result[float64] {
    return s.GetNumber(CategorySearch, string(key))
}

// InvalidateCache clears cached validation data
func (s *ValidationDataService) InvalidateCache() {
    s.cache = sync.Map{}
}
```

### Step 5: Correct Validator Implementation (Using Typed Constants)

```go
// ✅ CORRECT: Using typed constants - no magic strings
func (v *GuidelineValidator) validateTransitionDensity(content string, _ *SeoConfig) ValidationResult {
    // Load transition words using typed accessor
    transitions, err := v.validationData.GetSeoStringArray(SeoKeyTransitionWords)
    if err != nil {
        return ValidationResult{
            Passed:  false,
            Message: "Failed to load transition words from config",
        }
    }
    
    // Load threshold using typed accessor
    threshold, err := v.validationData.GetSeoNumber(SeoKeyTransitionDensityThreshold)
    if err != nil {
        threshold = 40.0  // Fallback if not configured
    }
    
    words := strings.Fields(strings.ToLower(content))
    transitionCount := 0
    
    transitionSet := make(map[string]bool)
    for _, t := range transitions {
        transitionSet[strings.ToLower(t)] = true
    }
    
    for _, word := range words {
        cleanWord := strings.Trim(word, ".,!?;:")
        if transitionSet[cleanWord] {
            transitionCount++
        }
    }
    
    density := float64(transitionCount) / float64(len(words)) * 100
    passed := density >= threshold
    
    return ValidationResult{
        Passed:  passed,
        Message: fmt.Sprintf("Transition density: %.1f%% (required: %.0f%%+)", density, threshold),
    }
}

// Example: Using RAG constants
func (s *RagService) GetStopWords() apperror.Result[[]string] {
    return s.validationData.GetRagStringArray(RagKeyStopWords)
}

// Example: Using Search constants
func (s *SearchService) GetAllowedFileTypes() apperror.Result[[]string] {
    return s.validationData.GetSearchStringArray(SearchKeyAllowedFileTypes)
}
```

---

## Anti-Pattern vs Correct Pattern

### ❌ WRONG: Magic Strings

```go
// Never do this
transitions, _ := v.validationData.GetStringArray("Seo", "TransitionWords")
threshold, _ := v.validationData.GetNumber("Seo", "TransitionDensityThreshold")
```

### ✅ CORRECT: Typed Constants

```go
// Always use typed constants
transitions, _ := v.validationData.GetSeoStringArray(SeoKeyTransitionWords)
threshold, _ := v.validationData.GetSeoNumber(SeoKeyTransitionDensityThreshold)
```

---

## Common Validation Data Categories

### SEO Validation

| Key | Type | Description |
|-----|------|-------------|
| `TransitionWords` | array | Words for transition density |
| `TransitionDensityThreshold` | number | Min % required |
| `MaxSentenceWords` | number | Max words per sentence |
| `MaxParagraphWords` | number | Max words per paragraph |
| `MinKeywordMentions` | number | Min keyword occurrences |
| `ForbiddenContainerTags` | array | Tags not allowed in containers |

### RAG Validation

| Key | Type | Description |
|-----|------|-------------|
| `StopWords` | array | Words to exclude from indexing |
| `MinChunkSize` | number | Minimum chunk size |
| `MaxChunkSize` | number | Maximum chunk size |
| `ChunkOverlap` | number | Overlap between chunks |

### Search Validation

| Key | Type | Description |
|-----|------|-------------|
| `AllowedFileTypes` | array | File extensions to index |
| `ExcludedDirectories` | array | Directories to skip |
| `MaxFileSize` | number | Max file size to process |

---

## Version Seeding Behavior

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    VALIDATION DATA SEEDING FLOW                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  config.seed.json (v1.3.0)                                               │
│  └── Categories.Seo.Settings.TransitionWords: [...]                     │
│                                                                          │
│                          ↓                                               │
│                                                                          │
│  ConfigService.SeedWithVersionCheck()                                    │
│  └── Check: SeedVersion (1.3.0) > DbVersion (1.2.0)?                  │
│                                                                          │
│                          ↓ YES                                           │
│                                                                          │
│  INSERT INTO ValidationData                                              │
│  └── Category: 'Seo'                                                    │
│  └── Key: 'TransitionWords'                                             │
│  └── Value: '["however","therefore",...]'                               │
│  └── Version: '1.3.0'                                                   │
│                                                                          │
│                          ↓                                               │
│                                                                          │
│  Update ConfigMeta.SeedVersion = '1.3.0'                               │
│  Append to changelog.md                                                  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## API for Runtime Updates

```
PUT /api/v1/config/validation/:category/:key
  Body: { "Value": [...], "Reason": "Added new transition words" }
  Effect: Updates ValidationData, invalidates cache, logs change

GET /api/v1/config/validation/:category
  Returns: All validation data for category

GET /api/v1/config/validation/:category/:key
  Returns: Specific validation data entry
```

---

## Checklist for New Validation Data

- [ ] Define in `config.seed.json` under appropriate category
- [ ] Bump config version (minor for new setting)
- [ ] Add changelog entry
- [ ] Create Go accessor using `ValidationDataService`
- [ ] Never hardcode the array in source code
- [ ] Add API endpoint for runtime updates if needed
- [ ] Add tests for default values

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CW Config Overview | `../00-overview.md` |
| RAG Chunk Settings | `./02-rag-chunk-settings.md` |
| RAG Validation Helpers | `./03-rag-validation-helpers.md` |
| AI SEO Guidelines | `../22-ai-bridge-cli/01-backend/17-ai-seo-core-guidelines.md` |
