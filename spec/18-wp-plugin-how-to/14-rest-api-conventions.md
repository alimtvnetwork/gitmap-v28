# Phase 14 — REST API Conventions

> **Purpose:** Define the complete REST API design standard — route naming, namespace, HTTP method selection, pagination, filtering, category grouping, `endpoints.json` registry, and controller organisation — so any AI can design and implement plugin REST APIs consistently.
> **Audience:** AI code generators and human developers.
> **Prerequisite:** Phases 1–5 must be read first (Foundation, Enums, Traits, Logging, Envelope).

---

## 14.1 Namespace Convention

Every plugin registers its REST routes under a single namespace:

```
{plugin-slug}/v{major}
```

| Component | Source | Example |
|-----------|--------|---------|
| Plugin slug | `PluginConfigType::Slug->value` | `my-plugin` |
| API version | `PluginConfigType::ApiVersion->value` | `v1` |
| Full namespace | `PluginConfigType::apiFullNamespace()` | `my-plugin/v1` |

### Rules

| Rule | Detail |
|------|--------|
| One namespace per plugin | Never register routes under someone else's namespace |
| Version is major only | `v1`, `v2` — never `v1.2` |
| Bump major only for breaking changes | New endpoints do NOT require a version bump |
| Namespace is constructed via enum | Never hardcode the namespace string in route registration |

### apiFullNamespace() implementation

```php
// In PluginConfigType enum
public static function apiFullNamespace(): string
{
    return self::Slug->value . '/' . self::ApiVersion->value;
}
```

---

## 14.2 Route Naming Conventions

### Resource-based naming

Routes follow a **resource/action** pattern using kebab-case:

```
/{resource}                    → List or create
/{resource}/{action}           → Perform action on resource
/{resource}/(?P<id>\d+)        → Single resource by ID
/{resource}/{sub-resource}     → Nested resource
```

### Naming rules

| Rule | ✅ Correct | ❌ Wrong |
|------|-----------|---------|
| Plural nouns for collections | `plugins`, `agents`, `snapshots` | `plugin`, `agent` |
| Kebab-case for multi-word paths | `upload-active`, `sync-manifest` | `uploadActive`, `sync_manifest` |
| Action verbs as path segments | `plugins/enable`, `logs/clear` | `enablePlugin`, `clearLogs` |
| Nested resources use parent prefix | `snapshots/settings`, `agents/plugins` | `snapshot-settings` |
| No file extensions | `/status` | `/status.json` |
| No trailing slashes | `/plugins` | `/plugins/` |

### Standard endpoint categories

Group endpoints by domain. Each group shares a route prefix:

| Category | Prefix | Purpose | Examples |
|----------|--------|---------|----------|
| **System** | (root) | Plugin health, diagnostics | `status`, `openapi`, `opcache-reset` |
| **Plugins** | `plugins/` | Plugin lifecycle management | `plugins`, `plugins/enable`, `plugins/delete` |
| **Sync** | `plugins/sync` | File synchronization | `plugins/sync`, `plugins/sync-manifest` |
| **Content** | (root) | WordPress content management | `posts`, `categories`, `media` |
| **Logs** | `logs/` | Log management and retrieval | `logs`, `logs/status`, `logs/clear` |
| **Diagnostics** | (root) | Error logs, debugging | `error-logs`, `error-sessions` |
| **Agents** | `agents/` | Multi-site agent management | `agents`, `agents/add`, `agents/sync` |
| **Snapshots** | `snapshots/` | Backup and snapshot management | `snapshots/list`, `snapshots/restore` |
| **Cloud Storage** | `cloud-storage/` | Remote storage providers | `cloud-storage/accounts`, `cloud-storage/upload` |
| **Users** | `users/` | User management | `users`, `users/export` |
| **Settings** | (root) | Plugin configuration | `site-settings`, `site-health-summary` |
| **Debug** | `debug/` | Development-only routes | `debug/routes` |

---

## 14.3 HTTP Method Selection

### Method rules

| Method | When to use | Side effects | Idempotent |
|--------|-------------|-------------|------------|
| **GET** | Retrieve data, list resources, read status | None | Yes |
| **POST** | Create resource, execute action, trigger operation | Yes | No |
| **PUT** | Replace entire resource | Yes | Yes |
| **PATCH** | Partial update of a resource | Yes | Yes |
| **DELETE** | Remove a resource or initiate removal | Yes | Yes |

### Decision guide

```
Is this a read operation?
  → YES → GET

Does this create or modify something?
  → Creates a new resource → POST
  → Replaces the whole resource → PUT
  → Updates specific fields → PATCH (or POST for simpler plugins)

Does this remove something?
  → YES → DELETE
```

### Plugin-specific patterns

| Pattern | Method | Example | Why |
|---------|--------|---------|-----|
| List all items | GET | `GET /plugins` | Read-only retrieval |
| Get single item | GET or POST | `GET /plugins/(?P<id>\d+)` or `POST /plugins/info` | GET when ID is in URL; POST when body carries the identifier |
| Create item | POST | `POST /agents/add` | Side effect: creates resource |
| Delete item | POST or DELETE | `POST /agents/remove` or `DELETE /logs/clear` | POST is acceptable for actions; DELETE for explicit removal |
| Execute action | POST | `POST /plugins/enable` | Actions are not idempotent |
| Retrieve with complex filters | GET | `GET /logs?action=Upload&from=2026-01-01` | Filters as query params |
| Upload file | POST | `POST /upload` | Multipart or base64 body |
| Export/download | GET or POST | `GET /export-self` or `POST /snapshots/export` | GET for simple exports; POST when body specifies parameters |

### When POST is acceptable for reads

WordPress plugins commonly use POST for "read with body" operations where query parameters would be too complex or expose sensitive data. This is acceptable when:

1. The request body carries a slug or identifier that should not be in the URL
2. The request carries multiple complex filter objects
3. The response is generated dynamically based on body parameters

```php
// Acceptable: POST with slug in body for single-item retrieval
// Route: POST /plugins/info
// Body: { "slug": "some-plugin" }
```

---

## 14.4 EndpointType Enum — Route Registry

Every route path is defined as an `EndpointType` enum case. No string literals in route registration.

### Structure

```php
enum EndpointType: string
{
    // ── System ──────────────────────────────────────────────
    case Status       = 'status';
    case Openapi      = 'openapi';
    case OpcacheReset = 'opcache-reset';

    // ── Plugins ─────────────────────────────────────────────
    case Plugins       = 'plugins';
    case PluginInfo    = 'plugins/info';
    case PluginEnable  = 'plugins/enable';
    case PluginDisable = 'plugins/disable';
    case PluginDelete  = 'plugins/delete';

    // ── Logs ────────────────────────────────────────────────
    case Logs       = 'logs';
    case LogsStatus = 'logs/status';
    case LogsClear  = 'logs/clear';

    /** Prefixes value with '/' for register_rest_route(). */
    public function route(): string
    {
        return '/' . $this->value;
    }

    // ── Group helpers ───────────────────────────────────────
    public function isPlugin(): bool  { return str_starts_with($this->value, 'plugins/'); }
    public function isLog(): bool     { return str_starts_with($this->value, 'logs/'); }
    public function isAgent(): bool   { return str_starts_with($this->value, 'agents'); }

    // Standard comparison methods
    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Enum rules

| Rule | Detail |
|------|--------|
| Case name is PascalCase | `PluginEnable`, not `plugin_enable` |
| Value is the path fragment (kebab-case) | `'plugins/enable'` |
| Group by category with comment headers | `// ── Plugins ──` |
| `route()` always prepends `/` | Used by `register_rest_route()` |
| Group helpers use `str_starts_with()` | For category-level checks |
| Dynamic segments use WordPress regex | `'users/(?P<id>\d+)'` |

---

## 14.5 Route Registration Pattern

Routes are registered in `RouteRegistrationTrait` using a grouped, fault-tolerant pattern:

```php
trait RouteRegistrationTrait
{
    public function registerRoutes(): void
    {
        $namespace = PluginConfigType::apiFullNamespace();

        // Fault-tolerant closure — logs failures without stopping other routes
        $safeRegister = function (
            string $route,
            array $args,
            string $groupName = '',
        ) use ($namespace): bool {
            try {
                register_rest_route($namespace, $route, $args);

                return true;
            } catch (Throwable $e) {
                $this->fileLogger->logException($e, "Route:{$groupName}{$route}");

                return false;
            }
        };

        // ── Group: System ────────────────────────────────────
        $safeRegister(EndpointType::Status->route(), [
            'methods'             => HttpMethodType::Get->value,
            'callback'            => [$this, 'handleStatus'],
            'permission_callback' => [$this, 'checkStatusPermission'],
        ], 'system');

        // ── Group: Plugins ───────────────────────────────────
        $safeRegister(EndpointType::PluginEnable->route(), [
            'methods'             => HttpMethodType::Post->value,
            'callback'            => [$this, 'handlePluginEnable'],
            'permission_callback' => [$this, 'checkPluginPermission'],
        ], 'plugins');

        // ... additional groups follow the same pattern
    }
}
```

### Registration rules

| Rule | Detail |
|------|--------|
| Every route uses `EndpointType` for path | No string literals |
| Every route uses `HttpMethodType` for method | `HttpMethodType::Get->value`, not `'GET'` |
| Every route has a `permission_callback` | **Never `__return_true`** — full rules in §14.5.1 |
| Every mutating route verifies a nonce | See §14.5.2 — applies to POST/PUT/PATCH/DELETE |
| Routes are grouped by domain | Comment separators between groups |
| Each group is wrapped in `$safeRegister` | Failure in one group does not block others |
| Callback points to a public trait method | `[$this, 'handleMethodName']` |

---

## 14.5.1 Security — Permission Callbacks (🔴 CODE RED)

> **Severity:** 🔴 CODE RED — non-negotiable. Linter rule
> [`WP-PERM-001`](../02-coding-guidelines/06-cicd-integration/06-rules-mapping.md)
> blocks merge on any violation.

### Hard rules

| # | Rule | Rationale |
|---|------|-----------|
| 1 | **`permission_callback` is mandatory on every route.** Omitting it is a fatal error, not a warning. | WordPress allows registration without one but emits a `_doing_it_wrong` notice and still serves the route — silent public exposure. |
| 2 | **`__return_true` is forbidden as the value of `permission_callback`.** No exceptions, no "but it's only a read endpoint", no "ping is harmless". | A public-by-default callback turns every future change to that route into an unaudited public endpoint. Public routes must declare intent **explicitly** (rule 3). |
| 3 | **Public routes MUST use a named, documented callback** — not an anonymous closure, not `__return_true`, not a string literal. The callback's docblock MUST state *why* the route is public. | Forces an audit trail; greppable by reviewers. |
| 4 | **Authenticated routes MUST call `current_user_can( $cap )` with a specific capability** (e.g. `manage_options`, `edit_posts`). Never `is_user_logged_in()` alone for mutating routes. | A logged-in subscriber is not authorized to flush caches. |
| 5 | **Capability strings MUST come from a `CapabilityType` enum**, never inline strings. | Same rule as `EndpointType` / `HttpMethodType` — no magic strings. |
| 6 | **The callback MUST return `bool` or `WP_Error`**, never `int`, `string`, or `null`. | Truthy non-bool values (e.g. `1`, `'allow'`) are silently coerced and bypass intent checks in custom guards. |

### Canonical patterns

**❌ Forbidden:**

```php
'permission_callback' => '__return_true',                    // rule 2
'permission_callback' => fn() => true,                       // rule 3 (closure, no audit)
'permission_callback' => 'is_user_logged_in',                // rule 4 (no capability)
'permission_callback' => [$this, 'isLoggedIn'],              // rule 4 (no capability)
// missing 'permission_callback' entirely                    // rule 1
```

**✅ Required — public route:**

```php
/**
 * Public on purpose: the /ping route exposes only static metadata
 * (author, company, version) declared in PluginConfigType. Approved
 * by security review YYYY-MM-DD, ticket SEC-123.
 */
public function allowPublicPing(): bool
{
    return true;
}
```

**✅ Required — authenticated route:**

```php
public function checkSettingsPermission(): bool|WP_Error
{
    if (!current_user_can(CapabilityType::ManagePluginSettings->value)) {
        return new WP_Error(
            'rest_forbidden',
            __('You are not allowed to manage these settings.', $pluginSlug),
            ['status' => 403]
        );
    }
    return true;
}
```

### Lint enforcement

The CI linter pack ships rule `WP-PERM-001` (PHP) which fails the
build on any of the forbidden patterns above. See
[`linters-cicd/checks/wp-permission-callback/php.py`](../../linters-cicd/checks/).

---

## 14.5.2 Security — Nonce Verification for Mutating Routes (🔴 CODE RED)

> **Severity:** 🔴 CODE RED — non-negotiable. Linter rule
> [`WP-NONCE-001`](../02-coding-guidelines/06-cicd-integration/06-rules-mapping.md)
> blocks merge on any violation.

REST authentication via cookies is vulnerable to CSRF unless the
request also presents a valid nonce. WordPress core ships
`wp_rest` for the REST API and `check_admin_referer()` /
`check_ajax_referer()` for `admin-post.php` and `admin-ajax.php`.

### Hard rules

| # | Rule | Surface |
|---|------|---------|
| 1 | Every `POST`, `PUT`, `PATCH`, `DELETE` route MUST verify a nonce **before** any side-effect runs. | REST + admin-post + admin-ajax |
| 2 | REST mutating routes MUST require the `X-WP-Nonce` header (issued via `wp_create_nonce('wp_rest')`). The permission callback MUST call `wp_verify_nonce()` or rely on WordPress's built-in cookie+nonce auth and reject requests without a valid nonce. | REST |
| 3 | `admin-post.php` handlers MUST call `check_admin_referer( $action, $field )` as the **first executable line** of the handler. | Admin form POST |
| 4 | `admin-ajax.php` handlers MUST call `check_ajax_referer( $action, 'nonce' )` as the **first executable line**. | Admin AJAX |
| 5 | Nonce action strings MUST come from an `AjaxActionType::nonceAction()` / `EndpointType::nonceAction()` enum method — never inline literals. | All |
| 6 | Nonces are **in addition to**, never a replacement for, the permission callback / capability check from §14.5.1. | All |
| 7 | A failed nonce check MUST short-circuit with HTTP `403` and a structured `WP_Error('rest_cookie_invalid_nonce', …)` — never a redirect, never a silent pass. | All |

### Canonical pattern — REST mutating route

```php
public function checkSettingsPermission(WP_REST_Request $request): bool|WP_Error
{
    // Rule 4 — capability first.
    if (!current_user_can(CapabilityType::ManagePluginSettings->value)) {
        return new WP_Error('rest_forbidden', __('Forbidden.', $slug), ['status' => 403]);
    }
    // Rule 2 — nonce second.
    $nonce = $request->get_header('x_wp_nonce');
    if (!wp_verify_nonce($nonce, 'wp_rest')) {
        return new WP_Error(
            'rest_cookie_invalid_nonce',
            __('Cookie nonce is invalid.', $slug),
            ['status' => 403]
        );
    }
    return true;
}
```

### Canonical pattern — admin-post.php handler

```php
public function handleSaveSettings(): void
{
    // Rule 3 — first executable line.
    check_admin_referer(
        AdminActionType::SaveSettings->nonceAction(),
        AdminActionType::SaveSettings->nonceField()
    );
    // Rule 6 — capability check still required.
    if (!current_user_can(CapabilityType::ManagePluginSettings->value)) {
        wp_die(__('Forbidden.', $slug), 403);
    }
    // … sanitize + persist …
}
```

### Lint enforcement

The CI linter pack ships rule `WP-NONCE-001` (PHP) which scans every
`register_rest_route()` with a non-GET method, every
`add_action('admin_post_*')`, and every `add_action('wp_ajax_*')`,
and asserts a nonce verification call appears in the handler. See
[`linters-cicd/checks/wp-nonce-verification/php.py`](../../linters-cicd/checks/).

---

## 14.6 Pagination

### Request parameters

Pagination uses query parameters for GET endpoints and JSON body fields for POST endpoints:

| Parameter | Type | Default | Max | Source |
|-----------|------|---------|-----|--------|
| `limit` | int | 50 | 500 | `PaginationConfigType::DefaultLimit` |
| `offset` | int | 0 | — | Computed from page number or passed directly |
| `page` | int | 1 | — | Alternative to offset: `offset = (page - 1) * limit` |

### PaginationConfigType enum

```php
enum PaginationConfigType: int
{
    case DefaultLimit = 50;
    case MaxLimit     = 500;

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Pagination extraction pattern

```php
private function extractPagination(WP_REST_Request $request): array
{
    $rawLimit = $request->get_param('limit');
    $rawOffset = $request->get_param('offset');
    $rawPage = $request->get_param('page');

    $defaultLimit = PaginationConfigType::DefaultLimit->value;
    $maxLimit = PaginationConfigType::MaxLimit->value;

    $limit = ($rawLimit !== null) ? min(absint($rawLimit), $maxLimit) : $defaultLimit;
    $isLimitZero = ($limit === 0);

    if ($isLimitZero) {
        $limit = $defaultLimit;
    }

    // Support both offset and page-based pagination
    $hasPage = ($rawPage !== null);
    $offset = $hasPage
        ? (max(1, absint($rawPage)) - 1) * $limit
        : absint($rawOffset ?? 0);

    return [
        'limit'  => $limit,
        'offset' => $offset,
    ];
}
```

### Pagination in response

Include pagination metadata in the `Attributes` section of the envelope:

```json
{
  "Status": { "IsSuccess": true, "Code": 200, "Message": "OK" },
  "Attributes": {
    "RequestedAt": "/my-plugin/v1/logs",
    "TotalRecords": 1250,
    "Limit": 50,
    "Offset": 100,
    "Page": 3,
    "TotalPages": 25
  },
  "Results": [ ... ]
}
```

---

## 14.7 Filtering

### FilterKeyType enum

Standardise all filter parameter names via an enum:

```php
enum FilterKeyType: string
{
    case Status        = 'status';
    case Plugin        = 'plugin';
    case Action        = 'action';
    case User          = 'user';
    case TriggeredBy   = 'triggeredBy';
    case UploadSource  = 'uploadSource';
    case From          = 'from';
    case To            = 'to';
    case SourceMachine = 'sourceMachine';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Filter extraction pattern

```php
private function extractFilters(WP_REST_Request $request): array
{
    $filters = [];

    foreach (FilterKeyType::cases() as $filter) {
        $value = $request->get_param($filter->value);
        $hasValue = ($value !== null && $value !== '');

        if ($hasValue) {
            $filters[$filter->value] = sanitize_text_field($value);
        }
    }

    return $filters;
}
```

### Filter naming rules

| Rule | Detail |
|------|--------|
| camelCase for filter keys | `triggeredBy`, not `triggered_by` |
| Match FilterKeyType enum values | No ad-hoc filter parameter names |
| Date filters use ISO 8601 format | `from=2026-01-01`, `to=2026-12-31` |
| String filters are sanitised | Always `sanitize_text_field()` |
| Empty string means "no filter" | Never treat `""` as a valid filter value |

---

## 14.8 RequestFieldType Enum — Body Field Names

For POST endpoints that accept JSON bodies, all field names are defined in `RequestFieldType`:

```php
enum RequestFieldType: string
{
    case PluginZip     = 'plugin_zip';
    case Slug          = 'slug';
    case Activate      = 'activate';
    case UploadSource  = 'upload_source';
    case PluginVersion = 'plugin_version';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Body field naming rules

| Rule | Detail |
|------|--------|
| snake_case for body fields | `plugin_zip`, `upload_source` |
| Defined in `RequestFieldType` enum | No inline string literals |
| Usage in handler: `$body[RequestFieldType::Slug->value]` | Type-safe field access |
| Required fields validated with guard clauses | See Phase 6 — Input Validation |

---

## 14.9 Response Key Convention

All response JSON keys use PascalCase, defined in `ResponseKeyType`:

```php
// Response keys — PascalCase
enum ResponseKeyType: string
{
    // Envelope keys
    case Success  = 'Success';
    case Error    = 'Error';
    case Message  = 'Message';
    case Data     = 'Data';

    // Pagination keys
    case Total    = 'Total';
    case Limit    = 'Limit';
    case Offset   = 'Offset';

    // Domain keys — named per resource
    case Plugins  = 'Plugins';
    case Agents   = 'Agents';
    case Logs     = 'Logs';
    // ...
}
```

### Key naming rules

| Layer | Convention | Example |
|-------|-----------|---------|
| Request body fields | snake_case | `plugin_zip`, `upload_source` |
| Query parameters | camelCase | `triggeredBy`, `uploadSource` |
| Response keys | PascalCase | `PluginSlug`, `TotalRecords`, `IsSuccess` |
| URL path segments | kebab-case | `upload-active`, `sync-manifest` |

---

## 14.10 Controller Organisation

### One handler trait per endpoint (preferred)

For clarity and single-responsibility, each endpoint (or tightly related pair) gets its own trait:

```
Traits/
├── Core/
│   ├── StatusHandlerTrait.php         ← GET /status
│   └── PluginInventoryTrait.php       ← GET /plugins
├── Upload/
│   └── UploadHandlerTrait.php         ← POST /upload, POST /upload-active
├── Activate/
│   └── ActivateHandlerTrait.php       ← POST /plugins/enable
├── Log/
│   ├── LogRetrievalTrait.php          ← GET /logs
│   ├── LogStatusTrait.php             ← GET /logs/status
│   ├── LogClearingTrait.php           ← DELETE /logs/clear, POST /logs/clear/confirm
│   └── LogEmailTrait.php              ← POST /logs/email
```

### When to combine endpoints in one trait

Combine only when two endpoints share the same private business logic method:

| ✅ Combine | ❌ Separate |
|-----------|-----------|
| `POST /upload` and `POST /upload-active` (same upload logic, different activation flag) | `GET /logs` and `POST /logs/email` (completely different logic) |
| `GET /agents` and `POST /agents/add` if add logic is trivial | `POST /plugins/enable` and `POST /plugins/delete` (different side effects) |

### Trait anatomy (recap from Phase 3)

```php
trait SomeHandlerTrait
{
    // Public: route handler → wraps in safeExecute
    public function handleSomething(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executeSomething($request),
            'something',
        );
    }

    // Private: business logic
    private function executeSomething(WP_REST_Request $request): WP_REST_Response
    {
        // 1. Extract and validate input (Phase 6)
        // 2. Execute business logic
        // 3. Log result
        // 4. Return envelope response

        return EnvelopeBuilder::success('Operation complete')
            ->setRequestedAt($request->get_route())
            ->setSingleResult($data)
            ->toResponse();
    }
}
```

---

## 14.11 endpoints.json — Endpoint Registry Data File

Every plugin maintains an `endpoints.json` file in `data/` that documents all registered REST routes. This file is:

1. **Human-readable documentation** — developers can scan all endpoints in one place
2. **Machine-consumable** — admin UI can render endpoint tables from this data
3. **Synchronised with `EndpointType`** — every enum case should have a matching entry

### File location

```
plugin-slug/
├── data/
│   └── endpoints.json
```

### Format specification

```json
{
  "namespace": "{plugin-slug}/v{major}",
  "version": "2.0.0",
  "description": "Complete listing of all REST API endpoints.",
  "endpoints": [
    {
      "path": "status",
      "methods": ["GET"],
      "category": "system",
      "description": "Plugin status, version, and registered routes",
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
    },
    {
      "path": "logs",
      "methods": ["GET"],
      "category": "diagnostics",
      "description": "Query transaction logs",
      "auth": true
    }
  ]
}
```

### Field definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | ✅ | Route path fragment (matches `EndpointType` value) |
| `methods` | string[] | ✅ | HTTP methods: `["GET"]`, `["POST"]`, `["GET", "POST"]` |
| `category` | string | ✅ | Grouping key from §14.2 categories |
| `description` | string | ✅ | One-line human-readable description |
| `auth` | boolean | ✅ | Whether authentication is required |
| `body` | object | ❌ | For POST/PUT — field names and their types |

### Body field type values

| Type string | Meaning |
|-------------|---------|
| `"string"` | Text field |
| `"boolean"` | `true` / `false` |
| `"integer"` | Whole number |
| `"base64"` | Base64-encoded binary data |
| `"array"` | JSON array |
| `"object"` | JSON object |
| `"string (slug)"` | Plugin slug identifier |

### Maintenance rules

| Rule | Detail |
|------|--------|
| Every `EndpointType` case must have a matching entry | Except internal-only cases like `WpJson` |
| Update `endpoints.json` when adding/removing endpoints | Keep synchronised with code |
| `version` field matches the file's revision | Bump when endpoints change |
| Categories match the grouping in §14.2 | Consistent across code and documentation |

---

## 14.12 Standard Endpoint Patterns

### Status endpoint (every plugin must have one)

```php
// GET /status — Health check and version info
public function handleStatus(WP_REST_Request $request): WP_REST_Response
{
    return $this->safeExecute(
        fn() => $this->executeStatus($request),
        'status',
    );
}

private function executeStatus(WP_REST_Request $request): WP_REST_Response
{
    return EnvelopeBuilder::success('OK')
        ->setRequestedAt($request->get_route())
        ->setSingleResult([
            ResponseKeyType::Version->value => PluginConfigType::Version->value,
            ResponseKeyType::Slug->value    => PluginConfigType::Slug->value,
            ResponseKeyType::Status->value  => 'active',
        ])
        ->toResponse();
}
```

### List endpoint with pagination and filters

```php
// GET /logs — Paginated, filterable list
private function executeLogs(WP_REST_Request $request): WP_REST_Response
{
    $pagination = $this->extractPagination($request);
    $filters = $this->extractFilters($request);

    $results = $this->queryLogs($filters, $pagination['limit'], $pagination['offset']);
    $total = $this->countLogs($filters);

    $totalPages = ($pagination['limit'] > 0)
        ? (int) ceil($total / $pagination['limit'])
        : 1;

    return EnvelopeBuilder::success('OK')
        ->setRequestedAt($request->get_route())
        ->setListResult($results)
        ->setTotalRecords($total)
        ->toResponse();
}
```

### Action endpoint with body validation

```php
// POST /plugins/enable — Activate a plugin
private function executePluginEnable(WP_REST_Request $request): WP_REST_Response
{
    $body = $request->get_json_params();
    $hasBody = ($body !== null && $this->isArray($body));

    if (!$hasBody) {
        return $this->validationError('Request body must be a JSON object', $request);
    }

    $slug = $body[RequestFieldType::Slug->value] ?? null;
    $hasSlug = ($slug !== null && $this->isString($slug));

    if (!$hasSlug) {
        return $this->validationError('Missing required field: slug', $request);
    }

    // ... business logic ...

    $this->fileLogger->info('Plugin enabled', ['slug' => $slug]);

    return EnvelopeBuilder::success('Plugin activated')
        ->setRequestedAt($request->get_route())
        ->setSingleResult([
            ResponseKeyType::PluginSlug->value => $slug,
            ResponseKeyType::Activated->value  => true,
        ])
        ->toResponse();
}
```

### Two-phase confirmation pattern

For destructive operations, use a two-phase confirm flow:

```
Phase 1: DELETE /logs/clear → returns confirmation token
Phase 2: POST /logs/clear/confirm → consumes token, executes deletion
```

```php
// Phase 1: Request deletion — returns a token
private function executeLogsClear(WP_REST_Request $request): WP_REST_Response
{
    $token = bin2hex(random_bytes(16));
    set_transient('plugin_clear_logs_token', $token, 300); // 5-minute TTL

    return EnvelopeBuilder::success('Confirmation required')
        ->setRequestedAt($request->get_route())
        ->setSingleResult([
            ResponseKeyType::Confirm->value => $token,
            ResponseKeyType::Message->value => 'Send this token to /logs/clear/confirm to execute',
        ])
        ->toResponse();
}

// Phase 2: Confirm deletion — validates token
private function executeLogsClearConfirm(WP_REST_Request $request): WP_REST_Response
{
    $body = $request->get_json_params();
    $token = $body[ResponseKeyType::Confirm->value] ?? '';
    $storedToken = get_transient('plugin_clear_logs_token');

    $isValid = ($token !== '' && $token === $storedToken);

    if (!$isValid) {
        return $this->validationError('Invalid or expired confirmation token', $request);
    }

    delete_transient('plugin_clear_logs_token');

    // ... execute deletion ...

    return EnvelopeBuilder::success('Logs cleared')
        ->setRequestedAt($request->get_route())
        ->setSingleResult([ResponseKeyType::Deleted->value => true])
        ->toResponse();
}
```

---

## 14.13 Dynamic Route Segments

For endpoints that accept IDs or identifiers in the URL path:

```php
// Route with regex capture group
case UserId = 'users/(?P<id>\d+)';

// Registration
$safeRegister(EndpointType::UserId->route(), [
    'methods'             => HttpMethodType::Get->value,
    'callback'            => [$this, 'handleUserById'],
    'permission_callback' => [$this, 'checkPluginPermission'],
    'args'                => [
        'id' => [
            'required'          => true,
            'validate_callback' => fn($param) => is_numeric($param),
            'sanitize_callback' => 'absint',
        ],
    ],
], 'users');

// Handler — extract from URL params
private function executeUserById(WP_REST_Request $request): WP_REST_Response
{
    $userId = absint($request->get_param('id'));
    // ...
}
```

### Regex patterns for dynamic segments

| Pattern | Matches | Example |
|---------|---------|---------|
| `(?P<id>\d+)` | Numeric ID | `users/42` |
| `(?P<slug>[a-zA-Z0-9-]+)` | Slug string | `plugins/my-plugin` |
| `(?P<provider>[a-zA-Z]+)` | Alpha provider name | `cloud-storage/settings/google` |

---

## 14.14 OpenAPI Documentation (Optional)

For complex plugins, maintain an `openapi.json` file in `data/`:

```
plugin-slug/
├── data/
│   ├── endpoints.json     ← Lightweight endpoint registry (required)
│   └── openapi.json       ← Full OpenAPI 3.0 spec (optional)
```

Optionally expose it via a REST endpoint:

```php
case Openapi = 'openapi';

// GET /openapi — serve the OpenAPI spec
private function executeOpenapi(WP_REST_Request $request): WP_REST_Response
{
    $specPath = plugin_dir_path(dirname(__DIR__, 2)) . 'data/openapi.json';
    $hasSpec = file_exists($specPath);

    if (!$hasSpec) {
        return $this->validationError('OpenAPI specification not found', $request);
    }

    $spec = json_decode(file_get_contents($specPath), true);

    return EnvelopeBuilder::success('OK')
        ->setRequestedAt($request->get_route())
        ->setSingleResult($spec)
        ->toResponse();
}
```

---

## 14.15 Summary Table

| Aspect | Convention | Reference |
|--------|-----------|-----------|
| Namespace | `{slug}/v{major}` via `PluginConfigType::apiFullNamespace()` | §14.1 |
| Route paths | Resource-based, kebab-case, defined in `EndpointType` | §14.2, §14.4 |
| HTTP methods | Via `HttpMethodType` enum; GET=read, POST=action/create, DELETE=remove | §14.3 |
| Route registration | Grouped, fault-tolerant `$safeRegister` closure | §14.5 |
| Permission callbacks | Mandatory; `__return_true` forbidden; capability via `CapabilityType` enum | §14.5.1 |
| Nonce verification | Required on every POST/PUT/PATCH/DELETE; `wp_rest` for REST, `check_admin_referer`/`check_ajax_referer` for admin | §14.5.2 |
| Pagination | `limit`/`offset` or `page` params; `PaginationConfigType` defaults | §14.6 |
| Filters | `FilterKeyType` enum; camelCase query params | §14.7 |
| Request body fields | `RequestFieldType` enum; snake_case | §14.8 |
| Response keys | `ResponseKeyType` enum; PascalCase | §14.9 |
| Controller org | One handler trait per endpoint; group by domain subfolder | §14.10 |
| Endpoint registry | `data/endpoints.json` file; synchronised with `EndpointType` | §14.11 |
| Standard endpoints | Status (required), List+paginate, Action+validate, Two-phase confirm | §14.12 |
| Dynamic segments | WordPress regex capture groups in `EndpointType` values | §14.13 |

---

## Cross-References

- [Phase 3 — Traits and Composition](03-traits-and-composition.md) — handler trait anatomy, route registration trait
- [Phase 4 — Logging and Error Handling](04-logging-and-error-handling.md) — safeExecute, error responses
- [Phase 5 — Helpers, Response Envelope](05-helpers-responses-and-integration.md) — EnvelopeBuilder, response format
- [Phase 6 — Input Validation](06-input-validation-patterns.md) — guard clauses, body validation
- [Phase 2 — Enums and Coding Style](02-enums-and-coding-style/00-overview.md) — enum patterns for all enum types used here

---

*Phase 14 completes the REST API design standard: from namespace to route naming, HTTP method selection, pagination, filtering, data files, and controller organisation.*
