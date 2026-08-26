# Phase 15 — Settings Architecture

> **Created:** 2026-04-09
> **Status:** ✅ Active
> **Depends on:** [Phase 12 — Design System](./12-design-system.md), [Phase 13 — Admin UI Patterns](./13-admin-ui-patterns.md)

---

## 1. Overview

Plugin settings define how a WordPress plugin is configured at runtime. This specification covers the data model, enum-driven key management, settings groups, defaults, validation, sanitization, UI layout, toggle switches, conditional display, save feedback, and action buttons. All settings MUST be type-safe, enum-driven, and follow the WordPress Options API.

---

## 2. Data Model

### 2.1 Option Name Registry (`OptionNameType`)

Every WordPress option key used by the plugin MUST be registered as a case in the `OptionNameType` enum. This prevents typos, enables IDE autocomplete, and makes option usage searchable.

```php
enum OptionNameType: string
{
    case SnapshotSettings   = 'RiseupSnapshotSettings';
    case LogRetrieval       = 'RiseupLogRetrievalSettings';
    case UpdateSettings     = 'RiseupUpdateSettings';
    case PluginSettings     = 'RiseupAsiaSettings';
    case ErrorNotification  = 'RiseupErrorNotificationSettings';
    case SupportSettings    = 'RiseupSupportSettings';
    case LastPluginVersion  = 'riseup_asia_last_version';

    /** WordPress core — value must remain snake_case. */
    case ActivePlugins      = 'active_plugins';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

**Rules:**
1. Plugin option names use PascalCase prefix + PascalCase suffix (e.g., `RiseupSnapshotSettings`)
2. WordPress core option names retain their original format (e.g., `active_plugins`)
3. Every `get_option()` and `update_option()` call MUST use `OptionNameType::Case->value`
4. Helper methods (`isEqual`, `isAnyOf`) follow the standard enum metadata pattern (Phase 02)

### 2.2 Settings Key Registry (`SettingsKeyType`)

Individual keys within a settings array MUST be registered in `SettingsKeyType`. This enum also handles migration from legacy snake_case to PascalCase:

```php
enum SettingsKeyType: string
{
    case PreferredProvider     = 'PreferredProvider';
    case ScheduleEnabled       = 'ScheduleEnabled';
    case ScheduleFrequency     = 'ScheduleFrequency';
    case WorkerPoolSize        = 'WorkerPoolSize';
    case StorageMode           = 'StorageMode';
    // ...

    public static function legacyMap(): array {
        return [
            'preferred_provider' => self::PreferredProvider,
            'schedule_enabled'   => self::ScheduleEnabled,
            // ...
        ];
    }

    public static function migrateArray(array $data): array {
        $map = self::legacyMap();
        $migrated = [];
        foreach ($data as $key => $value) {
            $enumCase = $map[$key] ?? self::tryFrom($key);
            if ($enumCase instanceof self) {
                $migrated[$enumCase->value] = $value;
            } else {
                $migrated[$key] = $value;
            }
        }
        return $migrated;
    }

    public static function isLegacyKey(string $key): bool {
        return isset(self::legacyMap()[$key]);
    }
}
```

**Rules:**
1. All settings keys use PascalCase as the enum value
2. The `legacyMap()` provides backward compatibility with snake_case keys
3. `migrateArray()` transforms an entire array in one call
4. Settings accessed via `$settings[SettingsKeyType::WorkerPoolSize->value]`

### 2.3 Storage Pattern

Settings are stored as serialized arrays in WordPress options:

```php
// Read
$settings = get_option(OptionNameType::SnapshotSettings->value, []);

// Write
update_option(OptionNameType::SnapshotSettings->value, $settings);

// Delete
delete_option(OptionNameType::SnapshotSettings->value);
```

Each `OptionNameType` case maps to one WordPress option row containing a serialized associative array. Individual keys within the array are defined by `SettingsKeyType` or inline strings for simpler option groups.

---

## 3. Settings Groups

### 3.1 Group Architecture

Settings are organized into logical groups, each rendered as a `.riseup-card` section:

| Group | Option Name | Description | UI Section |
|-------|------------|-------------|------------|
| Plugin Settings | `RiseupAsiaSettings` | Endpoints, auth, log retrieval | Main settings card |
| Auto-Update | `RiseupUpdateSettings` | Master URL, cache, resolved URL | Auto-update card |
| Snapshot Settings | `RiseupSnapshotSettings` | Provider, schedule, retention, worker pool | Snapshot card (partial) |
| Log Retrieval | via `RiseupAsiaSettings[log_retrieval]` | Which logs to include in API | Log retrieval card (partial) |
| Error Notification | `RiseupErrorNotificationSettings` | Error alert configuration | Error notification card |
| Support | `RiseupSupportSettings` | Support email, fallback URL | Support card |

### 3.2 Nested Settings

Some option groups use nested arrays for logical sub-grouping:

```php
// Nested structure in RiseupAsiaSettings
$settings = [
    'endpoints' => [
        'upload' => ['enabled' => 1, 'auth_required' => 1],
        'status' => ['enabled' => 1, 'auth_required' => 0],
    ],
    'log_retrieval' => [
        'include_error_log' => 1,
        'include_full_log'  => 0,
        'max_lines'         => 500,
    ],
];
```

HTML `name` attributes encode the nesting:
```html
<input name="RiseupAsiaSettings[endpoints][upload][enabled]" value="1">
<input name="RiseupAsiaSettings[log_retrieval][max_lines]" value="500">
```

---

## 4. Default Values

### 4.1 Default Value Sources

Defaults come from two sources:

1. **Config enums** — Numeric/string defaults defined in a config enum:
   ```php
   enum SnapshotConfigType: int|string {
       case RetentionDaysDefault  = 30;
       case RetentionCountDefault = 10;
       case BatchSize             = 1000;
       case WorkerPoolDefault     = 5;
       case WorkerPoolMin         = 1;
       case MaxSizeMb             = 500;
   }
   ```

2. **Value enums** — Defaults selected from value-type enums:
   ```php
   $preferredProvider = $settings[SettingsKeyType::PreferredProvider->value]
       ?? SnapshotProviderType::Auto->value;
   $storageMode = $settings[SettingsKeyType::StorageMode->value]
       ?? StorageModeType::PerTable->value;
   ```

### 4.2 Default Extraction Pattern

At the top of each settings partial, extract all values with defaults:

```php
$preferredProvider = $snapshotSettings[SettingsKeyType::PreferredProvider->value]
    ?? SnapshotProviderType::Auto->value;
$scheduleEnabled   = $snapshotSettings[SettingsKeyType::ScheduleEnabled->value]
    ?? false;
$scheduleFrequency = $snapshotSettings[SettingsKeyType::ScheduleFrequency->value]
    ?? SnapshotFrequencyType::Daily->value;
$retentionDays     = $snapshotSettings[SettingsKeyType::RetentionDays->value]
    ?? SnapshotConfigType::RetentionDaysDefault->value;
```

**Rules:**
1. ALWAYS use null coalescing (`??`) with an explicit default
2. Enum defaults reference the enum's `->value`, not the case itself
3. Boolean defaults use `false` (opt-in) or `true` (opt-out) explicitly
4. Extract all values at the top of the template/partial, before HTML output

---

## 5. Validation & Sanitization

### 5.1 WordPress Registration

Settings MUST be registered with `register_setting()` including a sanitize callback:

```php
register_setting(
    PluginConfigType::SettingsGroup->value,
    OptionNameType::PluginSettings->value,
    [
        'type'              => 'array',
        'sanitize_callback' => [$this, 'sanitizePluginSettings'],
    ]
);
```

### 5.2 Sanitize Callback Pattern

```php
public function sanitizePluginSettings(array $input): array {
    $sanitized = [];

    // Boolean toggle → checkbox sends '1' or nothing
    $sanitized['endpoints'] = [];
    foreach ($input['endpoints'] ?? [] as $endpoint => $config) {
        $sanitized['endpoints'][$endpoint] = [
            'enabled'       => isset($config['enabled']) ? 1 : 0,
            'auth_required' => isset($config['auth_required']) ? 1 : 0,
        ];
    }

    // Numeric with range
    $maxLines = intval($input['log_retrieval']['max_lines'] ?? 500);
    $sanitized['log_retrieval']['max_lines'] = max(50, min(5000, $maxLines));

    // URL
    $sanitized['master_url'] = esc_url_raw($input['master_url'] ?? '');

    // Email
    $sanitized['support_email'] = sanitize_email($input['support_email'] ?? '');

    return $sanitized;
}
```

### 5.3 Validation Rules by Field Type

| Field Type | Sanitization | Validation |
|-----------|-------------|------------|
| Text | `sanitize_text_field()` | Max length check |
| URL | `esc_url_raw()` | Must start with `https://` |
| Email | `sanitize_email()` | WordPress email validation |
| Number | `intval()` / `floatval()` | Min/max range clamping |
| Checkbox | `isset() ? 1 : 0` | Boolean coercion |
| Select | `sanitize_text_field()` | Must match enum values |
| Textarea | `sanitize_textarea_field()` | Max length check |

### 5.4 Enum-Constrained Selects

Select fields whose values come from enums MUST validate against the enum:

```php
$frequency = $input['schedule_frequency'] ?? '';
$validFrequency = SnapshotFrequencyType::tryFrom($frequency);
$sanitized['schedule_frequency'] = $validFrequency?->value
    ?? SnapshotFrequencyType::Daily->value;
```

---

## 6. Settings Page Layout

### 6.1 Page Structure

```
┌─ .wrap .riseup-admin ──────────────────────────────────────┐
│  Page Header (via partial)                                  │
│                                                              │
│  [Success Notice — if settings-updated]                     │
│                                                              │
│  <form method="post" action="options.php">                  │
│    <?php settings_fields('group'); ?>                        │
│                                                              │
│    ┌─ Card: Plugin Information ────────────────────────────┐ │
│    │  Version, API Namespace, REST Base (read-only)        │ │
│    └───────────────────────────────────────────────────────┘ │
│                                                              │
│    ┌─ Card: REST API Endpoints ───────────────────────────┐ │
│    │  Endpoint table with toggle switches                  │ │
│    └───────────────────────────────────────────────────────┘ │
│                                                              │
│    ┌─ Card: Auto-Update Settings ─────────────────────────┐ │
│    │  Toggle, URL, cache, diagnostics, action buttons      │ │
│    └───────────────────────────────────────────────────────┘ │
│                                                              │
│    <?php include 'partials/settings/section-snapshots.php' ?>│
│    <?php include 'partials/settings/section-logs.php' ?>     │
│                                                              │
│    ┌─ Card: Support & Feedback ───────────────────────────┐ │
│    │  Email, fallback URL                                  │ │
│    └───────────────────────────────────────────────────────┘ │
│                                                              │
│    <?php submit_button('Save Settings'); ?>                  │
│  </form>                                                     │
└──────────────────────────────────────────────────────────────┘
```

### 6.2 Form Submission

Settings pages use the WordPress Options API flow:

```php
<form method="post" action="options.php">
    <?php settings_fields(PluginConfigType::SettingsGroup->value); ?>
    <!-- Settings cards -->
    <?php submit_button(__('Save Settings', $pluginSlug)); ?>
</form>
```

This uses WordPress's built-in `options.php` handler which:
1. Verifies the nonce (from `settings_fields()`)
2. Calls the registered sanitize callback
3. Saves to `wp_options`
4. Redirects back with `?settings-updated=true`

### 6.3 Save Feedback

Success notices appear after redirect:

```php
<?php if (isset($_GET['settings-updated']) && $_GET['settings-updated']): ?>
    <div class="notice notice-success is-dismissible">
        <p><?php esc_html_e('Settings saved successfully.', $pluginSlug); ?></p>
    </div>
<?php endif; ?>
```

### 6.4 Settings Partials

Large settings pages MUST delegate sections to partials:

```php
// In admin-settings.php (orchestrator)
<?php include __DIR__ . '/partials/settings/section-snapshot-settings.php'; ?>
<?php include __DIR__ . '/partials/settings/section-log-retrieval.php'; ?>
```

Each partial documents its required variables:

```php
/**
 * Variables expected: $pluginSlug, $snapshotSettings, $snapshotProviders.
 */
```

---

## 7. Field Types & UI Patterns

### 7.1 Text Input

```php
<input type="text" id="field_id"
       name="<?php echo esc_attr(OptionNameType::Settings->value); ?>[key]"
       value="<?php echo esc_attr($value); ?>"
       class="regular-text"
       placeholder="Example text">
<p class="description"><?php esc_html_e('Help text.', $pluginSlug); ?></p>
```

### 7.2 URL Input

```php
<input type="url" id="field_id"
       name="<?php echo esc_attr(OptionNameType::Settings->value); ?>[key]"
       value="<?php echo esc_attr($value); ?>"
       class="regular-text"
       placeholder="https://example.com">
```

### 7.3 Email Input

```php
<input type="email" id="field_id"
       name="<?php echo esc_attr(OptionNameType::Settings->value); ?>[key]"
       value="<?php echo esc_attr($value ?? ''); ?>"
       class="regular-text"
       placeholder="support@example.com">
```

### 7.4 Number Input

```php
<input type="number" id="field_id"
       name="<?php echo esc_attr(OptionNameType::Settings->value); ?>[key]"
       value="<?php echo esc_attr($value); ?>"
       min="50" max="5000" step="50"
       class="small-text">
```

### 7.5 Select (Enum-Driven)

```php
<select id="field_id"
        name="<?php echo esc_attr(OptionNameType::Settings->value); ?>[key]">
    <?php foreach (MyEnumType::cases() as $case): ?>
        <option value="<?php echo esc_attr($case->value); ?>"
                <?php selected($currentValue, $case->value); ?>>
            <?php echo esc_html($case->label()); ?>
        </option>
    <?php endforeach; ?>
</select>
```

### 7.6 Time Input

```php
<input type="time" id="field_id" value="<?php echo esc_attr($scheduleTime); ?>">
```

### 7.7 Read-Only Display

```php
<tr>
    <th><?php esc_html_e('Version', $pluginSlug); ?></th>
    <td><code><?php echo esc_html(PluginConfigType::Version->value); ?></code></td>
</tr>
```

### 7.8 Diagnostic Display (Conditional)

```php
<?php $hasValue = BooleanHelpers::hasValue($settings['resolved_url'] ?? null); ?>
<?php if ($hasValue): ?>
    <code><?php echo esc_html($settings['resolved_url']); ?></code>
    <br><small class="text-muted">
        <?php printf(esc_html__('Cached on: %s', $pluginSlug), esc_html($settings['resolved_at'])); ?>
    </small>
<?php else: ?>
    <em><?php esc_html_e('Not resolved yet', $pluginSlug); ?></em>
<?php endif; ?>
```

---

## 8. Toggle Switch

### 8.1 HTML Structure

```php
<label class="toggle-switch">
    <input type="checkbox"
           id="field_id"
           name="<?php echo esc_attr(OptionNameType::Settings->value); ?>[key]"
           value="1"
           <?php checked($isEnabled); ?>>
    <span class="toggle-slider"></span>
</label>
```

### 8.2 CSS

```css
.toggle-switch {
    position: relative;
    display: inline-block;
    width: 44px;
    height: 24px;
    vertical-align: middle;
}
.toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
}
.toggle-slider {
    position: absolute;
    cursor: pointer;
    inset: 0;
    background: #cbd5e1;
    border-radius: 24px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.toggle-slider::before {
    content: '';
    position: absolute;
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background: white;
    border-radius: 50%;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);
}
.toggle-switch input:checked + .toggle-slider {
    background: var(--riseup-primary, #1d4ed8);
}
.toggle-switch input:checked + .toggle-slider::before {
    transform: translateX(20px);
}
.toggle-switch input:focus + .toggle-slider {
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}
```

### 8.3 Boolean Extraction

Use `BooleanHelpers::hasValue()` for consistent truthy detection:

```php
$isEnabled = BooleanHelpers::hasValue($settings['enabled'] ?? null);
checked($isEnabled);
```

### 8.4 Toggle in Tables

Toggles in endpoint configuration tables use nested `name` attributes:

```php
<label class="toggle-switch">
    <input type="checkbox"
           name="<?php echo esc_attr(OptionNameType::PluginSettings->value); ?>[endpoints][<?php echo esc_attr($endpoint); ?>][enabled]"
           value="1"
           <?php checked($isEnabled); ?>>
    <span class="toggle-slider"></span>
</label>
```

---

## 9. Conditional Display

### 9.1 Server-Side Conditional Rendering

Fields that depend on other values use inline `style` with PHP logic:

```php
<tr id="retention_days_row"
    style="<?php echo $retentionType !== RetentionType::Days->value ? 'display:none;' : ''; ?>">
```

### 9.2 Enum-Driven Conditional Logic

Complex conditionals use enum helper methods:

```php
<tr id="day_row"
    style="<?php
        $freq = SnapshotFrequencyType::tryFrom($scheduleFrequency);
        echo ($freq !== null && $freq->isAnyOf(
            SnapshotFrequencyType::Hourly,
            SnapshotFrequencyType::Daily,
            SnapshotFrequencyType::Manual
        )) ? 'display:none;' : '';
    ?>">
```

### 9.3 Client-Side Toggle (JavaScript)

```javascript
document.getElementById('retention_type').addEventListener('change', function() {
    document.getElementById('retention_days_row').style.display =
        this.value === 'days' ? '' : 'none';
    document.getElementById('retention_count_row').style.display =
        this.value === 'count' ? '' : 'none';
});
```

### 9.4 Rules

1. Initial visibility MUST be set server-side in PHP (no FOUC)
2. JavaScript handles dynamic toggling after page load
3. Hidden fields' values are preserved even when hidden
4. Conditional rows use `display: none` (not `visibility: hidden`)
5. Use `style.display = ''` to restore, not `style.display = 'table-row'`

---

## 10. Selection Cards (Radio Alternative)

For mutually exclusive settings with visual emphasis:

```php
<div class="riseup-storage-mode-cards" style="display: flex; gap: 12px;">
    <label class="riseup-mode-card" id="mode_card_single"
           style="flex: 1; cursor: pointer; padding: 12px;
                  border: 2px solid <?php echo $storageMode === 'single' ? '#2271b1' : '#dcdcde'; ?>;
                  border-radius: 8px;
                  background: <?php echo $storageMode === 'single' ? '#f0f6fc' : '#fff'; ?>;">
        <input type="radio" name="storage_mode" value="single"
               <?php checked($storageMode, 'single'); ?> style="display: none;">
        <span class="dashicons dashicons-database"></span>
        <strong>Single File</strong>
        <span class="description">All tables in one database.</span>
    </label>
    <!-- More cards... -->
</div>
```

**JavaScript toggles active state:**
```javascript
document.querySelectorAll('.riseup-mode-card input[type="radio"]').forEach(radio => {
    radio.addEventListener('change', () => {
        document.querySelectorAll('.riseup-mode-card').forEach(card => {
            const isActive = card.querySelector('input').checked;
            card.style.borderColor = isActive ? '#2271b1' : '#dcdcde';
            card.style.background = isActive ? '#f0f6fc' : '#fff';
        });
    });
});
```

---

## 11. Slider Input

For numeric range settings with visual feedback:

```php
<div style="display: flex; align-items: center; gap: 12px; max-width: 340px;">
    <input type="range" id="worker_pool_size"
           min="<?php echo esc_attr(SnapshotConfigType::WorkerPoolMin->value); ?>"
           max="<?php echo esc_attr(SnapshotConfigType::workerPoolMax()); ?>"
           value="<?php echo esc_attr($workerPoolSize); ?>"
           style="flex: 1; accent-color: #2271b1;">
    <span id="worker_pool_value"
          style="font-family: monospace; font-size: 14px; min-width: 24px;
                 text-align: center; font-weight: 600; color: #2271b1;">
        <?php echo esc_html($workerPoolSize); ?>
    </span>
</div>
```

JavaScript updates the display:
```javascript
document.getElementById('worker_pool_size').addEventListener('input', function() {
    document.getElementById('worker_pool_value').textContent = this.value;
});
```

---

## 12. Action Buttons

### 12.1 Settings Action Row

Settings sections may include action buttons that perform AJAX operations without saving:

```php
<tr>
    <th scope="row"><?php esc_html_e('Actions', $pluginSlug); ?></th>
    <td>
        <button type="button" id="btn_test_connection" class="button button-secondary">
            <span class="dashicons dashicons-yes-alt"></span>
            <?php esc_html_e('Test Connection', $pluginSlug); ?>
        </button>
        <button type="button" id="btn_clear_cache" class="button button-secondary">
            <span class="dashicons dashicons-trash"></span>
            <?php esc_html_e('Clear Cache', $pluginSlug); ?>
        </button>
        <button type="button" id="btn_check_updates" class="button button-secondary">
            <span class="dashicons dashicons-update"></span>
            <?php esc_html_e('Check Now', $pluginSlug); ?>
        </button>
        <span id="action_status" class="riseup-inline-status"></span>
    </td>
</tr>
```

### 12.2 Rules

1. Action buttons MUST use `type="button"` (not `submit`) to prevent form submission
2. Each button has a dashicon for visual identification
3. Inline status text shows operation result
4. Buttons in action rows use `button-secondary` variant
5. The primary save button uses `submit_button()` or `button-primary`

### 12.3 AJAX Action Flow

```javascript
document.getElementById('btn_test_connection').addEventListener('click', async function() {
    const btn = this;
    const status = document.getElementById('action_status');
    const icon = btn.querySelector('.dashicons');

    // Loading state
    icon.classList.add('spin');
    btn.disabled = true;
    status.textContent = '';

    try {
        const response = await fetch(ajaxUrl, { method: 'POST', body: formData });
        const data = await response.json();

        status.className = 'riseup-inline-status ' + (data.success ? 'success' : 'error');
        status.textContent = data.success ? '✓ Connected' : '✕ ' + data.message;
    } catch (e) {
        status.className = 'riseup-inline-status error';
        status.textContent = '✕ Request failed';
    } finally {
        icon.classList.remove('spin');
        btn.disabled = false;
    }
});
```

---

## 13. Warning Messages

### 13.1 Inline Warnings

Critical settings sections include warnings:

```php
<p class="riseup-warning">
    <span class="dashicons dashicons-warning"></span>
    <?php esc_html_e('Warning: Disabling authentication can expose your site.', $pluginSlug); ?>
</p>
```

### 13.2 Conditional Error Display

Error states from previous operations are shown conditionally:

```php
<?php $hasLastError = BooleanHelpers::hasValue($settings['last_error'] ?? null); ?>
<?php if ($hasLastError): ?>
<tr>
    <th scope="row"><?php esc_html_e('Last Error', $pluginSlug); ?></th>
    <td>
        <span class="riseup-error-text"><?php echo esc_html($settings['last_error']); ?></span>
    </td>
</tr>
<?php endif; ?>
```

### 13.3 Version Comparison

```php
<?php if (version_compare($settings['new_version'], PluginConfigType::Version->value, '>')): ?>
    <span class="dashicons dashicons-arrow-up-alt" style="color: #46b450;"></span>
    <span style="color: #46b450;"><?php esc_html_e('Update available!', $pluginSlug); ?></span>
<?php endif; ?>
```

---

## 14. Endpoint Configuration Table

A specialized settings pattern for per-endpoint toggles:

```php
<table class="wp-list-table widefat fixed striped riseup-endpoints-table">
    <thead>
        <tr>
            <th><?php esc_html_e('Endpoint', $pluginSlug); ?></th>
            <th><?php esc_html_e('Description', $pluginSlug); ?></th>
            <th><?php esc_html_e('Enabled', $pluginSlug); ?></th>
            <th><?php esc_html_e('Auth Required', $pluginSlug); ?></th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($endpointGroups as $groupKey => $group): ?>
            <tr class="endpoint-group-header">
                <td colspan="4">
                    <span class="dashicons <?php echo esc_attr($group['icon']); ?>"></span>
                    <?php echo esc_html($group['label']); ?>
                </td>
            </tr>
            <?php foreach ($group['endpoints'] as $endpoint => $meta): ?>
                <tr>
                    <td>
                        <strong><?php echo esc_html($meta['label']); ?></strong><br>
                        <code>/<?php echo esc_html($endpoint); ?></code>
                    </td>
                    <td><?php echo esc_html($meta['desc']); ?></td>
                    <td>
                        <label class="toggle-switch">
                            <input type="checkbox"
                                   name="Settings[endpoints][<?php echo esc_attr($endpoint); ?>][enabled]"
                                   value="1" <?php checked($isEnabled); ?>>
                            <span class="toggle-slider"></span>
                        </label>
                    </td>
                    <td>
                        <label class="toggle-switch">
                            <input type="checkbox"
                                   name="Settings[endpoints][<?php echo esc_attr($endpoint); ?>][auth_required]"
                                   value="1" <?php checked($isAuthRequired); ?>>
                            <span class="toggle-slider"></span>
                        </label>
                    </td>
                </tr>
            <?php endforeach; ?>
        <?php endforeach; ?>
    </tbody>
</table>
```

**Rules:**
1. Endpoints are grouped by category with group headers
2. Each endpoint has independent enabled/auth toggles
3. Group data comes from `data/endpoints.json` processed by PHP
4. Toggle names use nested array notation for WordPress serialization

---

## 15. Settings Section Partials

### 15.1 Extraction Criteria

Extract a settings section into a partial when:
1. The section has ≥5 form fields
2. The section has its own save/action button (AJAX-based)
3. The section is logically independent (snapshot settings, log settings)
4. The section is reused across multiple pages

### 15.2 Partial Structure

```php
<?php
/**
 * Settings Partial — Database Snapshot Settings card.
 *
 * Variables expected: $pluginSlug, $snapshotSettings, $snapshotProviders.
 *
 * @package RiseupAsiaUploader
 * @since   1.64.0
 */

if (!defined('ABSPATH')) {
    exit;
}

use RiseupAsia\Enums\SettingsKeyType;
// ... more use statements

// Extract values with defaults
$preferredProvider = $snapshotSettings[SettingsKeyType::PreferredProvider->value]
    ?? SnapshotProviderType::Auto->value;
// ... more extractions
?>
<!-- HTML card -->
<div class="riseup-card">
    <h2><span class="dashicons dashicons-database"></span> Title</h2>
    <!-- Sub-sections with <h3> dividers -->
    <h3>Provider</h3>
    <table class="form-table"><!-- fields --></table>

    <h3>Scheduling</h3>
    <table class="form-table"><!-- fields --></table>

    <!-- Actions -->
    <table class="form-table">
        <tr>
            <th>Actions</th>
            <td>
                <button type="button" class="button button-primary">Save</button>
                <button type="button" class="button button-secondary">Run Cleanup</button>
                <span class="riseup-inline-status"></span>
            </td>
        </tr>
    </table>
</div>
```

### 15.3 Sub-Section Dividers

Within a card, logical groups are separated by `<h3>` headings:

```php
<h3>
    <span class="dashicons dashicons-performance"></span>
    <?php esc_html_e('Worker Pool & Storage', $pluginSlug); ?>
</h3>
```

---

## 16. Dual Save Patterns

### 16.1 WordPress Form Save (Global)

The main settings form uses WordPress Options API with `submit_button()`:
- Saves via `POST` to `options.php`
- Redirects back with `?settings-updated=true`
- Success notice rendered from query param

### 16.2 AJAX Save (Section-Level)

Individual sections may have their own save buttons that use AJAX:
- Uses `type="button"` to prevent form submission
- Calls REST API or AJAX endpoint
- Shows inline status feedback
- Does NOT redirect

### 16.3 When to Use Each

| Pattern | Use Case |
|---------|----------|
| WordPress form save | Main settings page with `register_setting()` |
| AJAX save | Section-level settings managed by custom code (e.g., snapshot settings stored in SQLite) |

---

## 17. Anti-Patterns (NEVER DO)

1. ❌ Hardcode option names as strings — always use `OptionNameType::Case->value`
2. ❌ Store settings in individual options — group related settings in arrays
3. ❌ Skip sanitize callbacks in `register_setting()`
4. ❌ Use `$_POST` directly — always use WordPress Options API or `wp_verify_nonce()`
5. ❌ Default to `null` — always provide explicit defaults via enums or constants
6. ❌ Mix snake_case and PascalCase keys — migrate legacy keys via `SettingsKeyType::migrateArray()`
7. ❌ Use `<input type="submit">` for AJAX actions — use `type="button"`
8. ❌ Show/hide conditional fields with JavaScript only — set initial state server-side
9. ❌ Place action buttons outside a `<table class="form-table">` in settings context
10. ❌ Skip the `BooleanHelpers::hasValue()` check for checkbox values
