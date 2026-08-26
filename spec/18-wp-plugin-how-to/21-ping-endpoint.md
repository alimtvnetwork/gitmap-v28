# Phase 21 — Ping Endpoint

> **Version:** 1.0.0  
> **Updated:** 2026-04-16  
> **AI Confidence:** Production-Ready  
> **Ambiguity:** None  
> **Purpose:** Define the standard ping endpoint that every WordPress plugin must implement. The ping endpoint returns author, company, and version information. It supports both authorized and non-authorized modes, configurable per plugin.

---

## 21.1 Overview

Every WordPress plugin built under these conventions **must** expose a `/ping` endpoint. This endpoint serves as a lightweight health-check and identity beacon — returning who built the plugin, which company owns it, and what version is running.

### Why Ping?

| Reason | Detail |
|--------|--------|
| **Health check** | External monitors can verify the plugin is active and responding |
| **Identity** | Multi-vendor environments need to know who owns each endpoint |
| **Version audit** | Deployment pipelines verify the expected version is live |
| **Debugging** | Support teams confirm which version a site is running |

---

## 21.2 Authorization Modes

The ping endpoint supports **two authorization modes**. The plugin developer chooses which mode to use based on their security requirements. Both modes are documented here; the developer selects one via `PluginConfigType`.

> **🔴 CODE RED — security update (v4.24+):** Per
> [Phase 14 §14.5.1](./14-rest-api-conventions.md#1451-security--permission-callbacks--code-red),
> `__return_true` is **forbidden** as the value of `permission_callback`,
> including for "public" routes. Mode 1 below now resolves to a named
> `allowPublicPing()` method whose docblock declares why the route is
> public. Plugins authored against earlier versions of this spec MUST
> migrate; the `WP-PERM-001` linter blocks merge on any remaining
> `__return_true` callback.

### Mode 1: Non-Authorized (Public)

```php
'permission_callback' => [$this, 'allowPublicPing'],
```

| Aspect | Detail |
|--------|--------|
| **When to use** | Public-facing plugins, SaaS integrations, uptime monitors |
| **Pros** | Zero friction — any HTTP client can ping |
| **Cons** | Exposes author/version to unauthenticated users |
| **Security note** | Version exposure may aid targeted attacks; acceptable when the plugin is open-source or version info is already public |
| **Required helper** | `allowPublicPing(): bool` — returns `true`; docblock MUST state why the route is public (audit trail per §14.5.1 rule 3) |

### Mode 2: Authorized (Authenticated)

```php
'permission_callback' => [$this, 'checkPingPermission'],
```

| Aspect | Detail |
|--------|--------|
| **When to use** | Private/enterprise plugins, security-sensitive environments |
| **Pros** | Only authenticated users see author/version info |
| **Cons** | External uptime monitors need credentials |
| **Security note** | Recommended for closed-source or enterprise plugins |
| **Required helper** | `checkPingPermission(): bool\|WP_Error` — MUST call `current_user_can()` with a `CapabilityType` enum value, never `is_user_logged_in()` alone (§14.5.1 rule 4) |

### Configuration via PluginConfigType

```php
enum PluginConfigType: string
{
    case Author  = 'Md. Alim Ul Karim';
    case Company = 'Riseup Asia LLC';
    case Version = '2.5.0';

    /** Whether the ping endpoint requires authentication. */
    case IsPingAuthorized = 'true';  // 'true' or 'false'

    /**
     * Resolves the permission callback for the ping endpoint.
     *
     * Always returns a [object, method] pair — never the forbidden
     * '__return_true' string (see Phase 14 §14.5.1).
     */
    public static function pingPermissionCallback(object $handler): callable
    {
        $isAuthorized = self::IsPingAuthorized->value === 'true';

        if (!$isAuthorized) {
            return [$handler, 'allowPublicPing'];
        }

        return [$handler, 'checkPingPermission'];
    }
}
```

---

## 21.3 EndpointType Registration

Add the `Ping` case to the `EndpointType` enum:

```php
enum EndpointType: string
{
    // ── System ──────────────────────────────────────────────
    case Status       = 'status';
    case Ping         = 'ping';          // ← NEW
    case Openapi      = 'openapi';
    case OpcacheReset = 'opcache-reset';

    // ── Business ────────────────────────────────────────────
    case UserProfile  = 'user-profile';
    case Settings     = 'settings';
}
```

### Route registration

```php
// In RouteRegistrationTrait::registerRoutes()

// ── Ping ────────────────────────────────────────────────
$safeRegister(EndpointType::Ping->route(), [
    'methods'             => HttpMethodType::Get->value,
    'callback'            => [$this, 'handlePing'],
    'permission_callback' => PluginConfigType::pingPermissionCallback($this),
], 'system');
```

---

## 21.4 Response Format

The ping endpoint returns **exactly three fields** inside the standard response envelope:

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `Author` | string | `PluginConfigType::Author->value` | Full name of the plugin author |
| `Company` | string | `PluginConfigType::Company->value` | Company or organization name |
| `Version` | string | `PluginConfigType::Version->value` | SemVer version of the plugin |

### Response example

```json
{
  "Status": {
    "IsSuccess": true,
    "Code": 200,
    "Message": "pong"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin/v1/ping"
  },
  "Results": {
    "Author": "Md. Alim Ul Karim",
    "Company": "Riseup Asia LLC",
    "Version": "2.5.0"
  }
}
```

### ResponseKeyType additions

```php
enum ResponseKeyType: string
{
    case Status  = 'Status';
    case Message = 'Message';

    // ── Ping ────────────────────────────────────────────────
    case Author  = 'Author';
    case Company = 'Company';
    case Version = 'Version';
}
```

---

## 21.5 Handler Implementation

### PingHandlerTrait

```php
<?php

namespace PluginSlug\Traits\Core;

use WP_REST_Request;
use WP_REST_Response;
use PluginSlug\Enums\PluginConfigType;
use PluginSlug\Enums\ResponseKeyType;

/**
 * Handles the GET /ping endpoint.
 *
 * Returns author name, company, and plugin version.
 * Authorization mode is controlled by PluginConfigType::IsPingAuthorized.
 */
trait PingHandlerTrait
{
    /**
     * Route handler — wraps business logic in safeExecute.
     */
    public function handlePing(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executePing($request),
            'ping',
        );
    }

    /**
     * Business logic — returns author, company, version.
     */
    private function executePing(WP_REST_Request $request): WP_REST_Response
    {
        return EnvelopeBuilder::success('pong')
            ->setRequestedAt($request->get_route())
            ->setSingleResult([
                ResponseKeyType::Author->value  => PluginConfigType::Author->value,
                ResponseKeyType::Company->value => PluginConfigType::Company->value,
                ResponseKeyType::Version->value => PluginConfigType::Version->value,
            ])
            ->toResponse();
    }

    /**
     * Permission check for authorized mode.
     * Uses the same permission logic as the status endpoint.
     */
    public function checkPingPermission(WP_REST_Request $request): bool
    {
        return current_user_can('manage_options');
    }
}
```

### File location

```
includes/
└── Traits/
    └── Core/
        └── PingHandlerTrait.php
```

### Plugin.php integration

```php
class Plugin
{
    use PingHandlerTrait;          // ← Add this
    use StatusHandlerTrait;
    use RouteRegistrationTrait;
    use LoggerTrait;
}
```

---

## 21.6 PluginConfigType Additions

Add the following cases to `PluginConfigType` if not already present:

```php
enum PluginConfigType: string
{
    /** Plugin author — full legal name. */
    case Author  = 'Md. Alim Ul Karim';

    /** Company or organization name. */
    case Company = 'Riseup Asia LLC';

    /** Plugin version. */
    case Version = '2.5.0';

    /** Whether the ping endpoint requires authentication. */
    case IsPingAuthorized = 'true';
}
```

---

## 21.7 endpoints.json Entry

Add to `data/endpoints.json`:

```json
{
  "path": "ping",
  "methods": ["GET"],
  "category": "system",
  "description": "Health check — returns author, company, and plugin version",
  "auth": true
}
```

> **Note:** Set `"auth": false` if `IsPingAuthorized` is `'false'`.

---

## 21.8 Testing

### Manual test (non-authorized mode)

```bash
curl -s https://example.com/wp-json/my-plugin/v1/ping | jq .
```

### Manual test (authorized mode)

```bash
curl -s -H "Authorization: Bearer <token>" \
  https://example.com/wp-json/my-plugin/v1/ping | jq .
```

### Expected assertions

| Assertion | Check |
|-----------|-------|
| Status code | `200` |
| `Status.IsSuccess` | `true` |
| `Status.Message` | `"pong"` |
| `Results.Author` | Non-empty string |
| `Results.Company` | Non-empty string |
| `Results.Version` | Valid SemVer string |

---

## 21.9 Checklist

- [ ] `EndpointType::Ping` case added to enum
- [ ] Route registered in `RouteRegistrationTrait` with correct permission callback
- [ ] `PingHandlerTrait` created in `Traits/Core/`
- [ ] `PingHandlerTrait` composed into `Plugin.php`
- [ ] `PluginConfigType::Author`, `Company`, `IsPingAuthorized` cases present
- [ ] `ResponseKeyType::Author`, `Company`, `Version` cases present
- [ ] `endpoints.json` updated with ping entry
- [ ] Authorization mode matches plugin security requirements
- [ ] Manual test passes with expected response format

---

## Cross-References

- [Phase 14 — REST API Conventions](14-rest-api-conventions.md) — Route registration, EndpointType, response format
- [Phase 5 — Helpers, Response Envelope](05-helpers-responses-and-integration.md) — EnvelopeBuilder pattern
- [Phase 3 — Traits and Composition](03-traits-and-composition.md) — Handler trait anatomy
- [Phase 4 — Logging and Error Handling](04-logging-and-error-handling.md) — safeExecute wrapper

---

*Phase 21 — Ping Endpoint — v1.0.0 — 2026-04-16*
