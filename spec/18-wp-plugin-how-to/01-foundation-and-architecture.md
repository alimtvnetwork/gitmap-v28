# Phase 1 — Foundation and Architecture

> **Purpose:** Establish the non-negotiable structural rules every WordPress plugin must follow.
> **Audience:** AI code generators and human developers.
> **Cross-reference:** [`../02-coding-guidelines/01-cross-language/04-code-style/00-overview.md`](../02-coding-guidelines/01-cross-language/04-code-style/00-overview.md) for formatting details.

---

## 1.1 Minimum Requirements

| Requirement | Value |
|-------------|-------|
| PHP version | 8.1+ (enables backed enums, readonly properties, fibers, intersection types) |
| WordPress version | 5.6+ |
| Array syntax | Short `[]` only — never `array()` |
| Exception catching | Always `Throwable`, never `Exception` |

---

## 1.2 Folder Structure

Every plugin follows a single canonical directory layout. No deviations.

```
plugin-slug/
├── plugin-slug.php              ← Bootstrap file (non-namespaced)
├── includes/
│   ├── Autoloader.php           ← PSR-4 autoloader (non-namespaced)
│   ├── Core/
│   │   └── Plugin.php           ← Singleton entry point
│   ├── Enums/                   ← All backed enums
│   ├── Helpers/                 ← Stateless utility classes
│   ├── Logging/
│   │   └── FileLogger.php       ← File-based logger (singleton)
│   └── Traits/
│       ├── Auth/                ← Authentication traits
│       ├── Core/                ← Response, status, inventory traits
│       ├── Route/               ← Route registration trait
│       ├── [Feature]/           ← One subfolder per feature domain
│       └── ...
├── spec/                        ← Documentation
└── .ai-instructions             ← AI-specific constraints
```

### Rules

| Rule | Detail |
|------|--------|
| One class per file | File name matches class name exactly (`PascalCase.php`) |
| Namespaced files | All files under `includes/` are namespaced |
| Non-namespaced files | Only `plugin-slug.php` (bootstrap) and `Autoloader.php` |
| Traits subfolder | Each feature domain gets its own subfolder under `Traits/` |
| No business logic in `Core/Plugin.php` | Plugin.php only composes traits and wires hooks |

---

## 1.3 Bootstrap File (`plugin-slug.php`)

The bootstrap file sits at the plugin root. It performs exactly three tasks:

1. Declare the WordPress plugin header comment
2. Require the autoloader
3. Instantiate the Plugin singleton on the appropriate hook

The bootstrap file is **non-namespaced** because WordPress loads it directly. It must not contain any business logic, helper functions, or class definitions.

---

## 1.4 Autoloader (`includes/Autoloader.php`)

The autoloader is a `final` class with a single responsibility: map the plugin's root namespace to the `includes/` directory using PSR-4 conventions.

### Required behaviours

| Behaviour | Detail |
|-----------|--------|
| Namespace prefix check | Only load classes within the plugin's own namespace |
| File existence check | Verify the resolved file exists before `require_once` |
| Diagnostic logging | Write autoloader events to a dedicated log file in `wp-content/uploads/{slug}/logs/autoloader.log` |
| Error re-throw | If `require_once` fails, log the error with stack trace and re-throw |
| Self-register | Call `spl_autoload_register()` at the bottom of the file |

### Why diagnostic logging matters

The autoloader runs before any other plugin code. If it fails, the FileLogger is not yet available. Writing directly to a known file path ensures failures are always captured.

---

## 1.5 Core Plugin Class (`Core/Plugin.php`)

This is the plugin's composition root. It is a **singleton** that:

1. Uses `use` statements to compose all traits
2. Holds a `FileLogger` instance as a private property
3. Registers WordPress hooks in the constructor (and nothing else)
4. Optionally supports a verbose boot mode via a `wp-config.php` constant

### Singleton pattern

```
private static ?self $instance = null;

public static function getInstance(): self { ... }
private function __construct() { ... }
```

### Constructor responsibilities (exhaustive list)

| Step | Action |
|------|--------|
| 1 | Obtain the FileLogger singleton |
| 2 | Record start time with `microtime(true)` |
| 3 | Register `rest_api_init` hook (and any other hooks) |
| 4 | Log initialization summary with elapsed time |

The constructor must **never** call WordPress functions that are not yet available. It must **never** perform database queries, HTTP requests, or file I/O beyond logging.

---

## 1.6 ABSPATH Guard

Every PHP file under `includes/` must include the ABSPATH guard immediately after the namespace declaration:

```
namespace PluginName\Some\Namespace;

if (!defined('ABSPATH')) {
    exit;
}
```

**Placement rule:** The guard goes **after** the namespace line, never before it.

---

## 1.7 Global Type Imports

All global PHP types and WordPress classes used in a namespaced file must be imported via `use` statements at the top, without a leading backslash.

### Common imports

| Type | Import |
|------|--------|
| Throwable | `use Throwable;` |
| WP_REST_Request | `use WP_REST_Request;` |
| WP_REST_Response | `use WP_REST_Response;` |
| WP_Error | `use WP_Error;` |
| WP_User | `use WP_User;` |
| DateTimeZone | `use \DateTimeZone;` (only when PHP built-in conflicts with namespace) |

For non-namespaced files (bootstrap, autoloader), global types are used without `use` statements since they are already in global scope.

---

## 1.8 Version Tracking

The plugin version lives in exactly one place: the `PluginConfigType` enum (see Phase 2). Every version reference in headers, log lines, and API responses reads from this single source.

Any code change bumps at least the minor version.
