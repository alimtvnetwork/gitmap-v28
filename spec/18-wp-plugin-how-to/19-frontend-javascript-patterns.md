# Phase 18 — Frontend JavaScript Patterns

> **Purpose:** Define the conventions for writing admin-facing JavaScript in WordPress plugins. Covers the localized-object pattern, AJAX request helpers, DOM manipulation, modal interactions, status messaging, and file organization.
> **Audience:** AI code generators and human developers.
> **Prerequisite:** Phases 11 (Templates), 12 (Design System), 13 (Admin UI) must be read first.

---

## 18.1 File Organization

### Directory Structure

```
plugin-slug/
├── assets/
│   └── js/
│       ├── admin-settings.js     ← Settings page scripts
│       ├── admin-logs.js         ← Logs page scripts
│       ├── admin-agents.js       ← Agents page scripts
│       ├── admin-errors.js       ← Error management scripts
│       ├── admin-snapshots.js    ← Snapshots page scripts
│       ├── admin-license.js      ← License page scripts
│       └── admin-feedback.js     ← Feedback widget scripts
```

### Rules

| Rule | Detail |
|------|--------|
| Naming | `admin-{page-slug}.js` — one file per admin page |
| Max lines | 200 lines per file (see Phase 11) — extract helpers if exceeded |
| jQuery | Use WordPress-bundled jQuery via `jQuery(document).ready()` wrapper |
| No global pollution | All code inside the `jQuery(document).ready()` closure |
| Dependencies | Declare `jquery` as a dependency in `wp_enqueue_script()` |

---

## 18.2 Localized Object Pattern

Every JS file receives its PHP-dependent values via `wp_localize_script()`. This is the **only** way to pass data from PHP to JavaScript.

### PHP Side (Enqueue)

```php
wp_enqueue_script(
    'my-plugin-admin-agents',
    plugins_url('assets/js/admin-agents.js', __FILE__),
    ['jquery'],
    PluginConfigType::Version->value,
    true  // Load in footer
);

wp_localize_script('my-plugin-admin-agents', 'RiseupAgents', [
    'apiBase'      => rest_url(PluginConfigType::apiFullNamespace()),
    'nonce'        => wp_create_nonce('wp_rest'),
    'endpoints'    => [
        'agents' => EndpointType::Agents->value,
        'test'   => EndpointType::AgentsTest->value,
    ],
    'responseKeys' => [
        'agents'  => ResponseKeyType::Agents->value,
        'status'  => ResponseKeyType::Status->value,
    ],
    'i18n'         => [
        'noAgentsYet'  => __('No agent sites registered yet.', $pluginSlug),
        'confirmDelete'=> __('Are you sure you want to remove this agent?', $pluginSlug),
    ],
    'agentStatus'  => [
        'active'   => AgentStatusType::Active->value,
        'inactive' => AgentStatusType::Inactive->value,
    ],
]);
```

### JavaScript Side (Consumption)

```javascript
jQuery(document).ready(function($) {
    var C = window.RiseupAgents;   // Single reference to localized object
    var apiBase = C.apiBase;
    var nonce = C.nonce;

    var ENDPOINTS = C.endpoints;
    var LABELS = C.i18n;
    var RESPONSE_KEYS = C.responseKeys;
});
```

### Localized Object Shape Convention

| Key | Type | Purpose |
|-----|------|---------|
| `apiBase` | string | Full REST API base URL |
| `nonce` | string | WP REST nonce for `X-WP-Nonce` header |
| `endpoints` | object | Endpoint path values from `EndpointType` enum |
| `responseKeys` | object | Response field names from `ResponseKeyType` enum |
| `i18n` | object | All user-facing strings (translated via `__()`) |
| `actions` | object | AJAX action names for `wp_ajax_` handlers |
| `status` / `agentStatus` | object | Status enum values for conditional rendering |

### Critical Rules

| Rule | Detail |
|------|--------|
| No hardcoded strings | ALL user-facing text comes from `i18n` — never inline English strings in JS |
| No hardcoded URLs | API base and endpoints come from localized object — never construct URLs manually |
| No hardcoded keys | Response field names come from `responseKeys` — never hardcode `'Status'` or `'Agents'` |
| Enum-driven | Every localized value that maps to a PHP enum MUST use `EnumCase->value` in PHP |

---

## 18.3 AJAX Request Patterns

### Pattern A — REST API Requests

For REST endpoints registered via `register_rest_route()`:

```javascript
function apiRequest(method, endpoint, data) {
    return $.ajax({
        url: apiBase + '/' + endpoint,
        method: method,
        contentType: 'application/json',
        data: data ? JSON.stringify(data) : null,
        beforeSend: function(xhr) {
            xhr.setRequestHeader('X-WP-Nonce', nonce);
        }
    });
}

// Usage
apiRequest('GET', ENDPOINTS.agents).done(function(response) {
    var agents = response[RESPONSE_KEYS.agents];
    // ... render
}).fail(function(xhr) {
    showStatus('#status', 'Request failed', true);
});
```

### Pattern B — WordPress AJAX Requests

For `wp_ajax_` action handlers:

```javascript
$.post(ajaxurl, {
    action: C.actions.testConnection,
    nonce: ajaxNonce
}, function(response) {
    if (response.success) {
        showStatus('✓ ' + response.data.message, false);
    } else {
        showStatus('✗ ' + (response.data.message || 'Failed'), true);
    }
}).fail(function() {
    showStatus('✗ Request failed', true);
}).always(function() {
    // Re-enable button
});
```

### When to Use Each

| Pattern | Use When |
|---------|----------|
| REST API (Pattern A) | Endpoint is registered via `register_rest_route()` and returns `WP_REST_Response` |
| WordPress AJAX (Pattern B) | Handler uses `wp_ajax_{action}` hook and returns `wp_send_json_success/error` |

---

## 18.4 Status Message Helper

Every page needs a status message function:

```javascript
function showStatus(selector, message, isError) {
    $(selector).text(message)
        .removeClass('success error')
        .addClass(isError ? 'error' : 'success')
        .show();
    setTimeout(function() { $(selector).fadeOut(); }, 5000);
}
```

### Variant — Inline Status (no selector)

```javascript
function showStatus(message, isError) {
    $status.html(message)
        .css('color', isError ? '#dc3232' : '#46b450');
    setTimeout(function() { $status.fadeOut(); }, 5000);
    $status.show();
}
```

### Rules

| Rule | Detail |
|------|--------|
| Auto-dismiss | Always `setTimeout` with 5000ms fadeOut |
| Color | Success: `#46b450`, Error: `#dc3232` (from `colors.json` status group) |
| Prefix | Success: `✓`, Error: `✗` — applied by caller, not the helper |

---

## 18.5 Button State Management

Buttons MUST be disabled during AJAX operations to prevent double-submission:

```javascript
$('#btn_test_connection').on('click', function() {
    var $btn = $(this);

    // 1. Disable + show spinner
    $btn.prop('disabled', true)
        .find('.dashicons')
        .removeClass('dashicons-yes-alt')
        .addClass('dashicons-update spin');

    // 2. Make request
    $.post(ajaxurl, { /* ... */ })
        .done(function(response) { /* ... */ })
        .fail(function() { /* ... */ })
        .always(function() {
            // 3. Re-enable + restore icon
            $btn.prop('disabled', false)
                .find('.dashicons')
                .removeClass('dashicons-update spin')
                .addClass('dashicons-yes-alt');
        });
});
```

### Rules

| Rule | Detail |
|------|--------|
| Disable on click | `$btn.prop('disabled', true)` — ALWAYS in click handler, before AJAX |
| Spinner icon | Swap dashicon to `dashicons-update` + add `spin` class |
| Re-enable in `.always()` | Never in `.done()` or `.fail()` alone — `.always()` covers both |
| No double-click | Button stays disabled for entire request lifecycle |

---

## 18.6 Modal Interaction Patterns

### Opening a Modal

```javascript
// From button click
$('.toggle-details').on('click', function(e) {
    e.stopPropagation();
    var details = $(this).data('details');
    var formatted = JSON.stringify(details, null, 2);
    $('#riseup-details-content').text(formatted);
    $('#riseup-details-modal').show();
});

// From clickable row (data-details attribute)
$('.riseup-log-row.has-details').on('click', function(e) {
    if ($(e.target).is('button, a') || $(e.target).closest('button, a').length) {
        return;  // Don't open modal when clicking action buttons
    }
    var details = $(this).data('details');
    if (details) {
        $('#riseup-details-content').text(JSON.stringify(details, null, 2));
        $('#riseup-details-modal').show();
    }
});
```

### Closing a Modal

Three close triggers are ALWAYS required:

```javascript
// 1. Close button (×)
$('.riseup-modal-close').on('click', function() {
    $('#riseup-details-modal').hide();
});

// 2. Overlay click (click outside content)
$('.riseup-modal').on('click', function(e) {
    if (e.target === this) {
        $('#riseup-details-modal').hide();
    }
});

// 3. Escape key
$(document).on('keydown', function(e) {
    if (e.key === 'Escape') {
        $('#riseup-details-modal').hide();
    }
});
```

### Rules

| Rule | Detail |
|------|--------|
| Three close methods | ✕ button, overlay click, Escape key — ALL three required |
| Click delegation | Row clicks must exclude button/link targets via `e.target` check |
| `e.stopPropagation()` | Detail buttons must stop propagation to prevent row click handler |
| JSON formatting | Use `JSON.stringify(data, null, 2)` for readable display |

---

## 18.7 Table Rendering from API Data

### Dynamic Table Population

```javascript
function loadAgents() {
    $('#agents-loading').show();
    $('#agents-table').hide();

    apiRequest('GET', ENDPOINTS.agents).done(function(response) {
        var $tbody = $('#agents-tbody').empty();
        var agents = response[RESPONSE_KEYS.agents];

        if (!agents || agents.length === 0) {
            $tbody.append(
                '<tr class="no-agents"><td colspan="5">' +
                LABELS.noAgentsYet +
                '</td></tr>'
            );
        } else {
            agents.forEach(function(agent) {
                $tbody.append(buildAgentRow(agent));
            });
        }
    }).always(function() {
        $('#agents-loading').hide();
        $('#agents-table').show();
    });
}
```

### Rules

| Rule | Detail |
|------|--------|
| Loading state | Show loading indicator, hide table during fetch |
| Empty state | Display localized "no data" message in a full-colspan row |
| `.empty()` before append | Always clear tbody before re-rendering |
| `.always()` for cleanup | Hide loader in `.always()`, not `.done()` |
| HTML escaping | Use `escapeHtml()` helper for all user-supplied content |

### HTML Escape Helper

```javascript
function escapeHtml(str) {
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
}
```

---

## 18.8 URL Builder Helpers

### Per-Entity Endpoint URLs

```javascript
function buildAgentUrl(id, suffix) {
    return ENDPOINTS.agents + '/' + id + (suffix ? '/' + suffix : '');
}

function endpointSuffix(endpointValue) {
    var parts = endpointValue.split('/');
    return parts[parts.length - 1];
}
```

### Rules

- Build URLs from localized endpoint values — never hardcode paths
- Use helper functions for URL construction — never inline string concatenation
- Entity-specific endpoints follow `{resource}/{id}/{action}` pattern

---

## 18.9 Confirm-Before-Destructive-Action

```javascript
$('#riseup-clear-errors').on('click', function() {
    if (!confirm(C.i18n.confirmClearAll)) {
        return;
    }
    // ... proceed with AJAX
});
```

### Rules

| Rule | Detail |
|------|--------|
| Always confirm | Delete, clear, remove, reset — any destructive action |
| Localized message | Confirmation text from `i18n` — never hardcode |
| Return early | If cancelled, `return` immediately — don't nest the AJAX call |

---

## 18.10 Checklist

- [ ] One JS file per admin page, named `admin-{page}.js`
- [ ] All PHP values passed via `wp_localize_script()` — no inline `<script>` blocks
- [ ] Localized object uses enum `->value` for endpoints, response keys, status values
- [ ] All user-facing strings in `i18n` sub-object, translated with `__()`
- [ ] REST requests use `X-WP-Nonce` header via `beforeSend`
- [ ] WordPress AJAX uses `ajaxurl` + `action` + `nonce` pattern
- [ ] Buttons disabled during AJAX with spinner swap
- [ ] Modals have three close methods (button, overlay, Escape)
- [ ] Dynamic tables handle loading, empty, and populated states
- [ ] Destructive actions require `confirm()` with localized message
- [ ] `escapeHtml()` used for all user-supplied content in HTML

---

*Last Updated: 2026-04-09*
