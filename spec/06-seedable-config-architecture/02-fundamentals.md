# Seedable Config Architecture — Fundamentals

**Version:** 3.2.0
**Updated:** 2026-04-16
**Parent:** [01-index.md](./01-index.md)

---


```
┌─────────────────────────────────────────────────────────────────────────┐
│                        CW CONFIG VERSION FLOW                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  config.seed.json                                                        │
│  ┌─────────────────────────────────────────┐                             │
│  │ {                                       │                             │
│  │   "Version": "1.2.0",                   │ ← Source of truth           │
│  │   "Categories": { ... }                 │                             │
│  │ }                                       │                             │
│  └──────────────────┬──────────────────────┘                             │
│                     │                                                    │
│                     ▼                                                    │
│  ┌─────────────────────────────────────────┐                             │
│  │        Version Change Detected?          │                             │
│  └──────────────────┬──────────────────────┘                             │
│                     │                                                    │
│         ┌───────────┴───────────┐                                        │
│         │                       │                                        │
│         ▼                       ▼                                        │
│   ┌───────────┐          ┌───────────────┐                               │
│   │    NO     │          │     YES       │                               │
│   │ Skip Seed │          │ Merge + Seed  │                               │
│   └───────────┘          └───────┬───────┘                               │
│                                  │                                        │
│                                  ▼                                        │
│                          ┌───────────────┐                               │
│                          │ Update        │                               │
│                          │ changelog.md  │                               │
│                          └───────────────┘                               │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## File Specifications

### config.seed.json

The seed file contains default values and metadata:

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
      "Version": "1.2.0",
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

JSON Schema for validation:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Application Configuration",
  "Type": "object",
  "required": ["Version", "Categories"],
  "properties": {
    "Version": {
      "Type": "string",
      "pattern": "^\\\\d+\\\\.\\\\d+\\\\.\\\\d+$",
      "Description": "Semantic version of configuration"
    },
    "Changelog": {
      "Type": "string",
      "Description": "Description of changes in this version"
    },
    "Categories": {
      "Type": "object",
      "additionalProperties": {
        "$ref": "#/definitions/category"
      }
    }
  },
  "definitions": {
    "category": {
      "Type": "object",
      "required": ["DisplayName", "Settings"],
      "properties": {
        "DisplayName": { "Type": "string" },
        "Description": { "Type": "string" },
        "Version": { "Type": "string" },
        "AddedIn": { "Type": "string" },
        "Settings": {
          "Type": "object",
          "additionalProperties": {
            "$ref": "#/definitions/setting"
          }
        }
      }
    },
    "setting": {
      "Type": "object",
      "required": ["Type", "Label", "Default"],
      "properties": {
        "Type": {
          "Type": "string",
          "enum": ["string", "number", "boolean", "select", "array", "object"]
        },
        "Label": { "Type": "string" },
        "Description": { "Type": "string" },
        "Default": {},
        "Min": { "Type": "number" },
        "Max": { "Type": "number" },
        "Options": { "Type": "array" },
        "AddedIn": { "Type": "string" },
        "DeprecatedIn": { "Type": "string" }
      }
    }
  }
}
```

### changelog.md Format

```markdown

# Changelog

All notable configuration changes are documented here.

## [1.2.0] - 2026-02-01

### Added

- Cache category with Enabled, MaxSizeMb, TtlHours settings

### Changed

- Theme options now include "high-contrast"

## [1.1.0] - 2026-01-15

### Added

- Network category with port and timeout settings

## [1.0.0] - 2026-01-01

### Initial Release

- General category with Theme, Language, AutoSave
```

---

## Database Schema

### Table: ConfigMeta

```sql
-- linter-waive: MISSING-DESC-001 reason="Seedable-config example; focus on config schema mechanics"
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
-- linter-waive: MISSING-DESC-001 reason="Seedable-config example; focus on config schema mechanics"
CREATE TABLE Setting (
    SettingId INTEGER PRIMARY KEY AUTOINCREMENT,
    Category TEXT NOT NULL,
    Key TEXT NOT NULL,
    Value TEXT NOT NULL,           -- JSON encoded
    Type TEXT NOT NULL,
    AddedInVersion TEXT,         -- Version when setting was added
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
    ChangedBy TEXT,               -- user, system, seed
    Version TEXT,                  -- Version at time of change
    FOREIGN KEY (SettingId) REFERENCES Setting(SettingId)
);

CREATE INDEX IdxHistorySetting ON SettingsHistory(SettingId);
CREATE INDEX IdxHistoryChanged ON SettingsHistory(ChangedAt);
```

---

## Strongly-Typed Value Container

**CRITICAL: No `interface{}` or `any` usage in ConfigService code. All values use the `SettingValue` union struct.**

```go
// SettingValue is the strongly-typed union container for all config values.
// See 02-spec/22-ai-bridge-cli/01-backend/57-settings-service.md for the canonical definition.
type SettingValue struct {
    StringVal  *string            `json:"StringVal,omitempty"`
    IntVal     *int               `json:"IntVal,omitempty"`
    FloatVal   *float64           `json:"FloatVal,omitempty"`
    BoolVal    *bool              `json:"BoolVal,omitempty"`
    StringsVal []string           `json:"StringsVal,omitempty"`
    MapVal     map[string]string  `json:"MapVal,omitempty"`
}
```

---

## Go Implementation

### ConfigService

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
    db           *gorm.DB
    seedPath     string
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

// SeedWithVersionCheck seeds config if version changed
func (s *ConfigService) SeedWithVersionCheck() error {
    seed, err := s.loadSeedFile()
    if err != nil {
        return apperror.Wrap(
            err,
            ErrSeedLoadFailed,
            "load seed file",
        )
    }

    meta, err := s.getMeta()
    if err != nil {
        // First time - full seed
        return s.fullSeed(seed)
    }

    // Compare versions
    currentVer, _ := semver.NewVersion(meta.SeedVersion)
    seedVer, _ := semver.NewVersion(seed.Version)

    if !seedVer.GreaterThan(currentVer) {
        // No version change, skip seed
        return nil
    }

    // Version increased - merge new settings
    if err := s.mergeSeed(seed, meta.SeedVersion); err != nil {
        return err
    }

    // Update changelog
    return s.updateChangelog(seed)
}

// ConfigMetaUpdate is the typed struct for GORM .Updates() calls — no map[string]interface{}
type ConfigMetaUpdate struct {
    SeedVersion    string    `gorm:"column:SeedVersion"`
    CurrentVersion string    `gorm:"column:CurrentVersion"`
    LastSeededAt   time.Time `gorm:"column:LastSeededAt"`
    UpdatedAt      time.Time `gorm:"column:UpdatedAt"`
}

// mergeSeed adds new settings without overwriting existing
func (s *ConfigService) mergeSeed(seed SeedConfig, previousVersion string) error {
    for catKey, cat := range seed.Categories {
        for settingKey, setting := range cat.Settings {
            // Check if setting exists
            var existing Setting
            err := s.db.Where("category = ? AND key = ?", catKey, settingKey).First(&existing).Error

            if err == gorm.ErrRecordNotFound {
                // New setting - insert with default
                valueJson, _ := json.Marshal(setting.Default)
                newSetting := Setting{
                    ID:            generateId(),
                    Category:      catKey,
                    Key:           settingKey,
                    Value:         string(valueJson),
                    Type:          setting.Type,
                    AddedInVersion: seed.Version,
                }
                s.db.Create(&newSetting)
            }
            // Existing settings are preserved
        }
    }

    // Update meta using typed struct
    s.db.Model(&ConfigMeta{}).Where("ConfigMetaId = 1").Updates(ConfigMetaUpdate{
        SeedVersion:    seed.Version,
        CurrentVersion: seed.Version,
        LastSeededAt:   time.Now(),
        UpdatedAt:      time.Now(),
    })

    return nil
}

// updateChangelog appends version entry to changelog.md
func (s *ConfigService) updateChangelog(seed SeedConfig) error {
    if seed.Changelog == "" {
        return nil
    }

    entry := fmt.Sprintf("\n## [%s] - %s\n\n%s\n",
        seed.Version,
        time.Now().Format("2006-01-02"),
        seed.Changelog,
    )

    // Read existing changelog
    contentResult := pathutil.ReadFileIfExists(s.changelogPath)
    if contentResult.IsErr() {
        return contentResult.Err()
    }
    content := contentResult.Value()

    // Insert after header
    header := "# Changelog\n\nAll notable configuration changes are documented here.\n"

    if len(content) == 0 {
        content = []byte(header)
    }

    // Find insert position (after header)
    insertPos := len(header)
    if len(content) >= len(header) {
        insertPos = len(header)
    }

    newContent := string(content[:insertPos]) + entry + string(content[insertPos:])

    return pathutil.WriteFile(s.changelogPath, []byte(newContent))
}
```

---

## Version Bumping Rules

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| New category added | Minor | 1.0.0 → 1.1.0 |
| New setting added | Minor | 1.1.0 → 1.2.0 |
| Default value changed | Patch | 1.2.0 → 1.2.1 |
| Setting deprecated | Patch | 1.2.1 → 1.2.2 |
| Breaking change (setting removed) | Major | 1.2.2 → 2.0.0 |

---

## UI Integration

### Version Badge Component

```typescript
// components/VersionBadge.tsx
import { Badge } from '@/components/ui/badge';
import { useConfig } from '@/hooks/useConfig';

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
// Highlight settings added in current version
function SettingItem({ setting, currentVersion }: Props) {
  const isNew = setting.AddedInVersion === currentVersion;

  return (
    <div className={cn(
      "p-4 rounded-lg",
      isNew && "ring-2 ring-primary bg-primary/5"
    )}>
      {isNew && <Badge className="mb-2">New in v{currentVersion}</Badge>}
      {/* ... setting content */}
    </div>
  );
}
```

---

## Theme Support

### Comprehensive Theme System

The CW Config pattern supports a rich theme system with multiple customization options:

```json
{
  "Categories": {
    "Appearance": {
      "DisplayName": "Appearance",
      "Description": "Visual customization options",
      "Settings": {
        "Theme": {
          "Type": "select",
          "Label": "Theme",
          "Description": "Base color scheme",
          "Default": "system",
          "Options": [
            "light",
            "dark",
            "system",
            "high-contrast",
            "high-contrast-dark",
            "colorful-light",
            "colorful-dark",
            "ocean-blue",
            "ocean-dark",
            "forest-green",
            "forest-dark",
            "sunset-orange",
            "sunset-dark",
            "midnight-purple",
            "rose-pink",
            "slate-gray",
            "nord-light",
            "nord-dark",
            "solarized-light",
            "solarized-dark",
            "dracula",
            "monokai",
            "github-light",
            "github-dark"
          ]
        },
        "AccentColor": {
          "Type": "select",
          "Label": "Accent Color",
          "Description": "Primary action color",
          "Default": "blue",
          "Options": [
            "blue",
            "indigo",
            "violet",
            "purple",
            "fuchsia",
            "pink",
            "rose",
            "red",
            "orange",
            "amber",
            "yellow",
            "lime",
            "green",
            "emerald",
            "teal",
            "cyan",
            "sky"
          ]
        },
        "FontSize": {
          "Type": "select",
          "Label": "Font Size",
          "Description": "Base text size",
          "Default": "medium",
          "Options": ["x-small", "small", "medium", "large", "x-large"]
        },
        "FontFamily": {
          "Type": "select",
          "Label": "Font Family",
          "Description": "Text font style",
          "Default": "system",
          "Options": [
            "system",
            "inter",
            "roboto",
            "open-sans",
            "lato",
            "poppins",
            "source-sans",
            "jetbrains-mono",
            "fira-code"
          ]
        },
        "BorderRadius": {
          "Type": "select",
          "Label": "Border Radius",
          "Description": "Corner rounding style",
          "Default": "medium",
          "Options": ["none", "small", "medium", "large", "full"]
        },
        "AnimationSpeed": {
          "Type": "select",
          "Label": "Animation Speed",
          "Description": "UI transition speed",
          "Default": "normal",
          "Options": ["none", "reduced", "normal", "fast"]
        },
        "CompactMode": {
          "Type": "boolean",
          "Label": "Compact Mode",
          "Description": "Reduce padding and spacing",
          "Default": false
        },
        "ShowIcons": {
          "Type": "boolean",
          "Label": "Show Icons",
          "Description": "Display icons in navigation",
          "Default": true
        }
      }
    }
  }
}
```

### Theme CSS Variables

Each theme maps to CSS custom properties:

```css
/* Example: ocean-blue theme */
[data-theme="ocean-blue"] {
  --background: 200 30% 98%;
  --foreground: 200 50% 10%;
  --primary: 200 80% 50%;
  --primary-foreground: 200 10% 98%;
  --secondary: 180 40% 90%;
  --muted: 200 20% 95%;
  --accent: 180 60% 45%;
  --destructive: 0 70% 50%;
  --border: 200 20% 85%;
  --ring: 200 80% 50%;
  --radius: 0.5rem;
}

/* Example: dracula theme */
[data-theme="dracula"] {
  --background: 231 15% 18%;
  --foreground: 60 30% 96%;
  --primary: 265 89% 78%;
  --primary-foreground: 231 15% 18%;
  --secondary: 225 27% 26%;
  --muted: 232 14% 31%;
  --accent: 135 94% 65%;
  --destructive: 0 100% 67%;
  --border: 232 14% 31%;
  --ring: 265 89% 78%;
}
```

### React Theme Provider

```typescript
// hooks/useTheme.ts
import { useSettings } from './useSettings';

export function useTheme() {
  const { settings, updateSetting } = useSettings();

  const theme = settings?.Appearance?.Theme ?? 'system';
  const accentColor = settings?.Appearance?.AccentColor ?? 'blue';

  const setTheme = (newTheme: string) => {
    updateSetting('Appearance', 'Theme', newTheme);
    document.documentElement.setAttribute('data-theme', newTheme);
  };

  const setAccentColor = (color: string) => {
    updateSetting('Appearance', 'AccentColor', color);
    document.documentElement.setAttribute('data-accent', color);
  };

  return { theme, accentColor, setTheme, setAccentColor };
}
```

---

## Applicable Projects

This pattern is used by:

| Project | Config Location |
|---------|-----------------|
| Spec Management | `go-backend/configs/` |
| GSearch CLI | `backend/configs/` |
| BRun CLI | `backend/configs/` |
| AI Bridge | `backend/configs/` |
| Nexus Flow | `backend/configs/` |
| **WP SEO Publish CLI** | `backend/configs/` |

---

## CRITICAL: No Hardcoded Arrays

**All validation arrays, lookup tables, and configurable data MUST use CW Config → Root DB pattern.**

See: `06-validation-data-seeding.md` for complete implementation guide.

### Examples of Data That Must Be Seeded:

- Transition words for SEO validation
- Stop words for RAG indexing
- Allowed file types for search
- Forbidden HTML tags
- Thresholds (sentence length, paragraph length, etc.)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Split DB Architecture | [05-split-db-architecture/01-index.md](../05-split-db-architecture/01-index.md) |

---

*This pattern ensures all configuration changes are versioned and documented.*
