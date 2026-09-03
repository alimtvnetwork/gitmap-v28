# Phase 17 — Data File Patterns

> **Purpose:** Define conventions for static JSON data files that ship with the plugin, including schema design, enum-driven access, validation, and guidelines for adding new data files.
> **Audience:** AI code generators and human developers.
> **Prerequisite:** Phases 1–2 (enums) must be read first.

---

## 17.1 Data Directory Structure

```
plugin-slug/
├── data/
│   ├── .gitkeep              ← Ensures directory exists in version control
│   ├── colors.json           ← UI color tokens by group
│   ├── endpoints.json        ← REST API endpoint registry
│   └── openapi.json          ← OpenAPI 3.0 specification
```

### Rules

| Rule | Detail |
|------|--------|
| Location | Always `plugin-slug/data/` at plugin root |
| Format | JSON only — no YAML, no PHP arrays |
| Encoding | UTF-8, no BOM |
| Git | Include `.gitkeep` so the directory is tracked even if empty |
| Read-only | Data files are NEVER written to at runtime — they are static configuration |

---

## 17.2 colors.json — UI Color Tokens

### Schema

```json
{
    "logLevel": {
        "Error": "#dc3545",
        "Warn": "#fd7e14",
        "Info": "#0d6efd",
        "Debug": "#6c757d"
    },
    "status": {
        "success": "#46b450",
        "error": "#dc3232",
        "warning": "#dba617"
    },
    "wpAdmin": {
        "primary": "#2271b1",
        "primaryDark": "#135e96",
        "primaryBg": "#f0f6fc",
        "border": "#dcdcde",
        "textMuted": "#646970",
        "snapshotIncremental": "#7b1fa2",
        "headerAccent": "#667eea"
    }
}
```

### Structure Rules

| Rule | Detail |
|------|--------|
| Top-level keys | Each key is a **color group** — must have a matching `ColorGroupType` enum case |
| Group values | Object of `label → hex color` pairs |
| Color format | 6-digit hex with `#` prefix, lowercase hex digits |
| Naming | camelCase for group keys and color names; PascalCase for log level names (matches enum case names) |

### Formal JSON Schema

The `colors.json` file MUST conform to this schema. Use it for CI validation and AI code generation:

```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Plugin Color Tokens",
    "description": "UI color tokens grouped by purpose. Each top-level key maps to a ColorGroupType enum case.",
    "type": "object",
    "additionalProperties": {
        "type": "object",
        "description": "A color group — keys are token names, values are hex colors.",
        "additionalProperties": {
            "type": "string",
            "pattern": "^#[0-9a-fA-F]{6}$",
            "description": "6-digit hex color with # prefix."
        },
        "minProperties": 1
    },
    "minProperties": 1,
    "examples": [
        {
            "logLevel": { "Error": "#dc3545", "Info": "#0d6efd" },
            "status": { "success": "#46b450", "error": "#dc3232" }
        }
    ]
}
```

### Validation Rules

| Rule | Detail |
|------|--------|
| Every top-level key | MUST have a matching `ColorGroupType` enum case |
| Every value | MUST be a 6-digit hex string matching `^#[0-9a-fA-F]{6}$` |
| No empty groups | Each group must contain at least one color token |
| No duplicate keys | JSON parsing naturally prevents this, but CI should verify |

### Enum-Driven Access

```php
enum ColorGroupType: string {
    case LogLevel = 'logLevel';
    case Status   = 'status';
    case WpAdmin  = 'wpAdmin';
}
```

### ColorConfig — Static Cache Helper

Instead of reading `colors.json` inline, use a dedicated helper class with static caching:

```php
namespace PluginName\Helpers;

use PluginName\Enums\ColorGroupType;

class ColorConfig
{
    /** @var array<string, array<string, string>>|null */
    private static ?array $colors = null;

    /** Default fallback color (muted gray). */
    private const FALLBACK = '#6c757d';

    /** Load and cache the colors.json file. */
    private static function load(): array
    {
        $isLoaded = (self::$colors !== null);

        if ($isLoaded) {
            return self::$colors;
        }

        $path = PathHelper::getColorsJsonPath();
        $isFileMissing = PathHelper::isFileMissing($path);

        if ($isFileMissing) {
            self::$colors = [];
            return self::$colors;
        }

        $json = @file_get_contents($path);

        if ($json === false) {
            self::$colors = [];
            return self::$colors;
        }

        $decoded = json_decode($json, true);
        $isDecodeFailed = ($decoded === null);

        if ($isDecodeFailed) {
            self::$colors = [];
            return self::$colors;
        }

        self::$colors = $decoded;
        return self::$colors;
    }

    /** Get a color by group and key. */
    public static function get(ColorGroupType $group, string $key, string $fallback = self::FALLBACK): string
    {
        $colors = self::load();
        $g = $group->value;
        $hasGroup = isset($colors[$g]);

        if ($hasGroup) {
            $hasKey = isset($colors[$g][$key]);
            if ($hasKey) {
                return $colors[$g][$key];
            }
        }

        return $fallback;
    }

    /** Get an entire color group as an associative array. */
    public static function getGroup(ColorGroupType $group): array
    {
        $colors = self::load();
        $g = $group->value;

        return isset($colors[$g]) ? $colors[$g] : [];
    }

    /** Convenience: Get a log level color. */
    public static function logLevel(string $level): string
    {
        return self::get(ColorGroupType::LogLevel, $level);
    }

    /** Convenience: Get a status color. */
    public static function status(string $status): string
    {
        return self::get(ColorGroupType::Status, $status);
    }

    /** Reset the static cache (for testing). */
    public static function reset(): void
    {
        self::$colors = null;
    }
}
```

### Usage Examples

```php
// Direct access via ColorConfig
$errorColor = ColorConfig::logLevel('Error');          // '#dc3545'
$successColor = ColorConfig::status('success');        // '#46b450'
$primary = ColorConfig::get(ColorGroupType::WpAdmin, 'primary');  // '#2271b1'

// Get full group for iteration (e.g., building a legend)
$logColors = ColorConfig::getGroup(ColorGroupType::LogLevel);
foreach ($logColors as $level => $hex) {
    echo "<span style=\"color:{$hex}\">{$level}</span>";
}
```

### Adding a New Color Group

1. Add the group object to `colors.json`
2. Add a matching case to `ColorGroupType` enum
3. Optionally add a convenience method to `ColorConfig`
4. Reference via `ColorGroupType::NewGroup->value` — never hardcode the string key

---

## 17.3 endpoints.json — REST API Registry

### Schema

```json
{
    "namespace": "plugin-slug/v1",
    "version": "2.0.0",
    "description": "Human-readable description of the API",
    "endpoints": [
        {
            "path": "status",
            "methods": ["GET"],
            "category": "system",
            "description": "Plugin status and version",
            "auth": true
        },
        {
            "path": "upload",
            "methods": ["POST"],
            "category": "plugins",
            "description": "Upload plugin via base64-encoded ZIP",
            "auth": true,
            "body": {
                "plugin_zip": "base64",
                "slug": "string",
                "activate": "boolean"
            }
        }
    ]
}
```

### Endpoint Object Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | ✅ | Route path relative to namespace (no leading slash) |
| `methods` | string[] | ✅ | HTTP methods: `GET`, `POST`, `PUT`, `DELETE` |
| `category` | string | ✅ | Logical grouping (e.g., `system`, `plugins`, `agents`, `snapshots`, `content`, `diagnostics`, `sync`) |
| `description` | string | ✅ | Human-readable purpose |
| `auth` | boolean | ✅ | Whether authentication is required |
| `body` | object | ❌ | Request body schema for POST/PUT (field name → type hint) |

### Category Conventions

| Category | Purpose | Examples |
|----------|---------|---------|
| `system` | Plugin health, config, maintenance | `status`, `openapi`, `opcache-reset`, `export-self` |
| `plugins` | Plugin CRUD operations | `upload`, `plugins`, `plugins/enable` |
| `sync` | Delta file synchronization | `plugins/sync-manifest`, `plugins/sync` |
| `content` | WordPress content management | `posts`, `categories`, `media` |
| `agents` | Remote site management | `agents`, `agents/add`, `agents/test` |
| `snapshots` | Backup and restore operations | `snapshots/list`, `snapshots/restore` |
| `diagnostics` | Logging and error inspection | `logs`, `error-logs`, `error-sessions` |

### Synchronization Rule

**CRITICAL:** Every entry in `endpoints.json` MUST have a corresponding `EndpointType` enum case, and vice versa. They must stay synchronized. When adding a new endpoint:

1. Add the enum case to `EndpointType`
2. Add the entry to `endpoints.json`
3. Register the route in `RouteRegistrationTrait`

---

## 17.4 openapi.json — OpenAPI 3.0 Specification

### Purpose

Provides machine-readable API documentation for Swagger UI integration and external tooling.

### Structure

```json
{
    "openapi": "3.0.3",
    "info": {
        "title": "Plugin Name API",
        "description": "...",
        "version": "2.0.0",
        "contact": { "name": "...", "url": "..." },
        "license": { "name": "GPL v2 or later", "url": "..." }
    },
    "servers": [{
        "url": "{baseUrl}/wp-json/plugin-slug/v1",
        "variables": {
            "baseUrl": { "default": "https://example.com" }
        }
    }],
    "security": [{ "basicAuth": [] }],
    "paths": { /* ... */ }
}
```

### Rules

| Rule | Detail |
|------|--------|
| Version sync | `info.version` MUST match `endpoints.json` version |
| Server URL | Use template variable `{baseUrl}` for portability |
| Auth | Default security scheme is `basicAuth` |
| Paths | Must cover all endpoints listed in `endpoints.json` |

---

## 17.5 Adding a New Data File

### Checklist

1. **Create the JSON file** in `data/` with a clear, descriptive name (kebab-case)
2. **Create a corresponding enum** (if the file has keyed groups) in `includes/Enums/`
3. **Document the schema** — add field descriptions in this spec or a dedicated section
4. **Read via `file_get_contents` + `json_decode`** — never use `require` or `include` for JSON
5. **Cache if needed** — for frequently accessed data, cache the decoded array in a static property
6. **Never write at runtime** — data files are deployment artifacts, not runtime state

### Caching Pattern

```php
final class ColorRegistry
{
    private static ?array $cache = null;

    public static function get(ColorGroupType $group): array
    {
        if (self::$cache === null) {
            $path = plugin_dir_path(__FILE__) . '../../data/colors.json';
            self::$cache = json_decode(file_get_contents($path), true) ?: [];
        }
        return self::$cache[$group->value] ?? [];
    }
}
```

---

## 17.6 Validation Rules

### At Development Time

| Check | Tool | Trigger |
|-------|------|---------|
| Valid JSON syntax | `json_decode` returns non-null | Build/CI |
| Enum sync | All JSON keys have matching enum cases | Code review |
| No duplicate paths | `endpoints.json` paths are unique | CI lint |

### At Runtime

| Check | Action |
|-------|--------|
| File missing | Log error via `InitHelpers::errorLogWithPrefix()`, use empty defaults |
| Decode failure | Log error, return empty array — never throw |
| Missing key | Return `null` or default — never access without null-safe check |

---

## 17.7 Checklist

- [ ] `data/` directory with `.gitkeep`
- [ ] `colors.json` with `ColorGroupType` enum covering all top-level keys
- [ ] `endpoints.json` synchronized with `EndpointType` enum
- [ ] `openapi.json` version-matched with `endpoints.json`
- [ ] Enum-driven access for all keyed data (never hardcode string keys)
- [ ] Static cache pattern for frequently accessed data files
- [ ] Runtime fallbacks: missing file → log + empty default, never throw

---

*Last Updated: 2026-04-09*
