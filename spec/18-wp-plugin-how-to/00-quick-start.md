# Phase 00 — Quick Start Guide

> **Purpose:** Condensed onboarding for AI models. Read this first, then dive into individual phases as needed.  
> **Full index:** [readme.md](readme.md)  
> **Rule:** Every decision below links to the authoritative phase. When in doubt, the phase document wins.

---

## 0.1 What This Spec Covers

An 18-phase, self-contained guide for building WordPress plugins following the **Gold Standard** architecture. The spec covers everything from folder structure to deployment rollback, admin UI to REST APIs.

---

## 0.2 Before You Write Any Code

| Decision | Answer | Phase |
|----------|--------|-------|
| Minimum PHP | **8.1+** (backed enums, readonly, fibers) | [Phase 1, §1.1](01-foundation-and-architecture.md) |
| Minimum WordPress | **5.6+** | [Phase 1, §1.1](01-foundation-and-architecture.md) |
| Array syntax | `[]` only — never `array()` | [Phase 1, §1.1](01-foundation-and-architecture.md) |
| Exception catching | Always `Throwable`, never `Exception` | [Phase 1, §1.1](01-foundation-and-architecture.md) |
| Magic strings | **Forbidden** — use backed enums everywhere | [Phase 2](02-enums-and-coding-style/00-overview.md) |

---

## 0.3 Canonical Folder Structure

```
plugin-slug/
├── plugin-slug.php              ← Bootstrap (non-namespaced, 3 tasks only)
├── includes/
│   ├── Autoloader.php           ← PSR-4 autoloader (non-namespaced)
│   ├── Core/
│   │   └── Plugin.php           ← Singleton, composes traits, wires hooks
│   ├── Enums/                   ← All backed enums
│   ├── Helpers/                 ← Stateless utility classes
│   ├── Logging/
│   │   └── FileLogger.php       ← Structured file logger (singleton)
│   └── Traits/
│       ├── Auth/                ← Authentication traits
│       ├── Core/                ← ResponseTrait, StatusHandlerTrait
│       ├── Route/               ← RouteRegistrationTrait
│       └── [Feature]/           ← One subfolder per feature domain
├── templates/                   ← Admin page templates (200-line max)
├── assets/                      ← CSS/JS/images
├── data/                        ← JSON data files (colors, endpoints)
└── spec/                        ← Documentation
```

> **Phase 1, §1.2** — [Full folder structure rules](01-foundation-and-architecture.md)

---

## 0.4 The 5 Files You Always Create First

Build these in order. Phase 7 has copy-paste-ready versions of each.

### 1. Bootstrap (`plugin-slug.php`)

Three tasks only: plugin header, require autoloader, instantiate singleton.

```php
if (!defined('ABSPATH')) { exit; }
define('PLUGIN_NAME_DEBUG', false);
require_once __DIR__ . '/includes/Autoloader.php';
add_action('rest_api_init', function (): void {
    PluginName\Core\Plugin::getInstance();
});
```

> [Phase 7, §7.1](07-reference-implementations.md) — Full bootstrap

### 2. Autoloader (`includes/Autoloader.php`)

PSR-4 mapping from plugin namespace → `includes/` directory. Must log to `autoloader.log`.

> [Phase 7, §7.2](07-reference-implementations.md) — Full autoloader

### 3. Plugin Singleton (`includes/Core/Plugin.php`)

Composes all traits. No business logic here — only `use` statements and hook wiring.

> [Phase 7, §7.3](07-reference-implementations.md) — Full Plugin.php

### 4. ResponseTrait (`includes/Traits/Core/ResponseTrait.php`)

Provides `safeExecute()` error boundary and standard response builders.

> [Phase 3, §3.4](03-traits-and-composition.md) — safeExecute pattern  
> [Phase 7, §7.4](07-reference-implementations.md) — Full ResponseTrait

### 5. FileLogger (`includes/Logging/FileLogger.php`)

Structured file-based logging with rotation, deduplication, and stack traces.

> [Phase 4, §4.3](04-logging-and-error-handling.md) — FileLogger spec  
> [Phase 7, §7.5](07-reference-implementations.md) — Full FileLogger

---

## 0.5 Critical Patterns (Must-Know)

### Backed Enums — No Magic Strings

Every string constant, status code, option name, or configuration key is a backed enum.

```php
enum HttpMethodType: string {
    case GET  = 'GET';
    case POST = 'POST';
    // Metadata via match expressions — see Phase 2
}
```

> [Phase 2](02-enums-and-coding-style/00-overview.md) — Enum architecture (4 subcategories)

### safeExecute() — Universal Error Boundary

Every endpoint handler wraps its logic in `safeExecute()`. Never let exceptions escape.

```php
public function handleSomething(WP_REST_Request $request): WP_REST_Response {
    return $this->safeExecute('handleSomething', function () use ($request) {
        // ... business logic ...
        return $this->successResponse($data, 'Operation completed');
    });
}
```

> [Phase 3, §3.4](03-traits-and-composition.md) — safeExecute anatomy

### Two-Tier Logging

| Tier | When available | Use |
|------|---------------|-----|
| Tier 1 — `error_log()` / `ErrorLogHelper` | Always | Bootstrap failures, pre-autoloader |
| Tier 2 — `FileLogger` | After autoloader | Everything else |

> [Phase 4, §4.1](04-logging-and-error-handling.md) — Two-tier architecture

### Standard Response Envelope

All API responses use `EnvelopeBuilder` with PascalCase keys:

```json
{
    "Status": "success",
    "Message": "Operation completed",
    "Results": { ... },
    "Attributes": { "Count": 5, "Page": 1 }
}
```

> [Phase 5](05-helpers-responses-and-integration.md) — EnvelopeBuilder API

### Fail-Fast Input Validation

Guard clauses at the top of every handler. Use enum-based validation.

> [Phase 6](06-input-validation-patterns.md) — Validation patterns

---

## 0.6 Adding a New Feature — Checklist

When adding a feature (e.g., "manage widgets"):

| Step | Action | Phase |
|------|--------|-------|
| 1 | Create enum(s) in `Enums/` for any new constants | [Phase 2](02-enums-and-coding-style/00-overview.md) |
| 2 | Create trait(s) in `Traits/[Feature]/` | [Phase 3](03-traits-and-composition.md) |
| 3 | Add `use` statement in `Plugin.php` | [Phase 1, §1.5](01-foundation-and-architecture.md) |
| 4 | Register routes in `RouteRegistrationTrait` | [Phase 3, §3.6](03-traits-and-composition.md) |
| 5 | Add input validation guard clauses | [Phase 6](06-input-validation-patterns.md) |
| 6 | Add database migration if needed | [Phase 8, §8.5](08-wordpress-integration-patterns.md) |
| 7 | Add admin page template if needed (≤200 lines) | [Phase 11](11-frontend-and-template-patterns.md) |
| 8 | Add tests | [Phase 9](09-testing-patterns.md) |

> [Phase 5](05-helpers-responses-and-integration.md) — Full integration checklist

---

## 0.7 File Size Limits

| Category | Max lines | Ideal |
|----------|-----------|-------|
| Template files | 200 | 50–100 |
| Trait files | 200 | 50–150 |
| Helper classes | 200 | 50–100 |
| Enum files | No hard limit | Group by domain |

If a file exceeds 200 lines, decompose using the **Orchestrator Pattern**: a parent template that `include`s partial components.

> [Phase 11, §11.3](11-frontend-and-template-patterns.md) — Orchestrator pattern

---

## 0.8 Phase Reference Map

Read in this order for fastest onboarding:

| Priority | Phases | What you learn |
|----------|--------|----------------|
| **Read first** | 1, 2 | Structure, enums, coding style |
| **Core runtime** | 3, 4, 5 | Traits, logging, response envelope |
| **Input handling** | 6 | Validation |
| **See it all together** | 7 | Complete reference implementations |
| **WordPress integration** | 8, 9, 10 | Admin pages, testing, deployment |
| **Frontend** | 11, 12, 13 | Templates, design system, UI patterns |
| **APIs & config** | 14, 15 | REST conventions, settings architecture |
| **Advanced** | 16, 17, 18 | Error diagnostics, data files, JS patterns |

---

## 0.9 Common Mistakes to Avoid

| ❌ Don't | ✅ Do | Phase |
|----------|-------|-------|
| Call WordPress functions in constructors | Use lazy initialization / hook callbacks | [Phase 1](01-foundation-and-architecture.md) |
| Use `is_array()` directly | Use `TypeCheckerTrait::isArray()` | [Phase 3, §3.8](03-traits-and-composition.md) |
| Hardcode strings | Use backed enums | [Phase 2](02-enums-and-coding-style/00-overview.md) |
| Let exceptions escape handlers | Wrap in `safeExecute()` | [Phase 3, §3.4](03-traits-and-composition.md) |
| Use `array()` syntax | Use `[]` short syntax | [Phase 1, §1.1](01-foundation-and-architecture.md) |
| Put business logic in Plugin.php | Compose via traits | [Phase 3](03-traits-and-composition.md) |
| Skip input validation | Guard clauses at handler top | [Phase 6](06-input-validation-patterns.md) |
| Use regex in route patterns | Plain string endpoints | [Phase 14](14-rest-api-conventions.md) |
| Include source maps in production | Dev-only source maps | [Phase 11, §11.8](11-frontend-and-template-patterns.md) |
| Store option names as strings | `OptionNameType` enum | [Phase 15](15-settings-architecture.md) |

---

## 0.10 Quick Decision Matrix

| "I need to…" | Go to |
|---------------|-------|
| Create a new plugin from scratch | Phase 7 (copy reference implementations) |
| Add a REST endpoint | Phases 3, 6, 14 |
| Add an admin settings page | Phases 8, 12, 13, 15 |
| Add a database table | Phase 8, §8.5 |
| Seed data from JSON | Phase 8, §8.5.1 + Phase 17 |
| Build deployment/self-update | Phase 10 |
| Style admin pages | Phases 12, 13 |
| Add JS interactivity | Phases 11, 18 |
| Handle errors properly | Phases 4, 16 |
| Add settings with toggles | Phase 15 |

---

*Start with Phase 1 for the full specification. Return here whenever you need to orient.*
