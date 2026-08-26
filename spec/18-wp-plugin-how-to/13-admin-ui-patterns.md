# Phase 13 — Admin UI Patterns

> **Created:** 2026-04-09
> **Status:** ✅ Active
> **Depends on:** [Phase 12 — Design System](./12-design-system.md)

---

## 1. Overview

This specification defines the standard UI patterns, layouts, and component anatomy used across all WordPress plugin admin pages. Every admin page MUST follow these patterns to maintain visual and behavioral consistency. The design system tokens from Phase 12 are assumed.

---

## 2. Page Layout Architecture

### 2.1 Root Wrapper

Every admin page MUST be wrapped in a standard structure:

```php
<div class="wrap riseup-admin">
    <?php
    $pageIcon = 'dashicons-database';
    $pageTitle = __('Page Title', $pluginSlug);
    $pageDescription = __('Optional description text.', $pluginSlug);
    include __DIR__ . '/partials/shared/page-header.php';
    ?>

    <!-- Page content as .riseup-card sections -->
</div>
```

- `wrap` — WordPress admin wrapper class (required for proper margins)
- `riseup-admin` — Plugin namespace for CSS scoping
- Page-specific class (e.g., `riseup-agents`, `riseup-snapshots`) added when page-level CSS overrides are needed

### 2.2 Page Header Partial

The `page-header.php` partial renders a consistent header across all pages:

```
┌──────────────────────────────────────────────────────────┐
│ 🔧 Plugin Name — Page Title  v2.10.0  [optional badge]  │
│ Optional description paragraph                           │
└──────────────────────────────────────────────────────────┘
```

**Required variables:**
| Variable | Type | Description |
|----------|------|-------------|
| `$pageIcon` | string | Dashicons class (e.g., `dashicons-database`) |
| `$pageTitle` | string | Translated page title |
| `$pluginSlug` | string | Plugin text domain |

**Optional variables:**
| Variable | Type | Description |
|----------|------|-------------|
| `$pageDescription` | string | Translated description paragraph |
| `$headerExtra` | string | Raw HTML after version badge (e.g., error count badge) |

**Implementation:**

```php
<h1>
    <span class="dashicons <?php echo esc_attr($pageIcon); ?>"></span>
    <?php echo esc_html($pageTitle); ?>
    <span class="riseup-version-badge">v<?php echo esc_html(PluginConfigType::Version->value); ?></span>
    <?php if (!empty($headerExtra)) { echo $headerExtra; } ?>
</h1>
<?php if (!empty($pageDescription)): ?>
<p class="description"><?php echo esc_html($pageDescription); ?></p>
<?php endif; ?>
```

### 2.3 Content Sections (Cards)

All page content is organized into `.riseup-card` sections. Each card is a logical unit:

```
┌─ .riseup-card ──────────────────────────────────┐
│  <h2>                                            │
│    <dashicon> Section Title                      │
│    [optional: action buttons, badges]            │
│  </h2>                                           │
│                                                  │
│  [content: form, table, stats, etc.]             │
└──────────────────────────────────────────────────┘
```

**Rules:**
1. Each card has exactly one `<h2>` heading with a dashicon
2. Cards are stacked vertically with consistent spacing
3. Cards MAY contain nested `.form-table`, `.wp-list-table`, or custom content
4. Action buttons that affect card content go inside the `<h2>` or immediately below it

### 2.4 Grid Layouts

For side-by-side content (e.g., chart + calendar), use a grid row:

```php
<div class="riseup-analytics-row">
    <div class="riseup-card"><!-- Primary content (flexible) --></div>
    <div class="riseup-card"><!-- Secondary content (fixed width) --></div>
</div>
```

```css
.riseup-analytics-row {
    display: grid;
    grid-template-columns: 1fr 340px;
    gap: 20px;
}
@media (max-width: 1100px) {
    .riseup-analytics-row { grid-template-columns: 1fr; }
}
```

---

## 3. Actions Bar

The actions bar provides primary page-level operations:

```
┌─ .riseup-card ──────────────────────────────────────────────┐
│  ┌─ .riseup-actions-row ──────────────────────────────────┐ │
│  │ [🔵 Primary] [⬜ Secondary] [⬜ Secondary] [⬜ Refresh]│ │
│  │ <inline-status>                                        │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌─ Expandable Options (hidden by default) ──────────────┐  │
│  │ .form-table with contextual settings                   │  │
│  │ [✅ Confirm] [Cancel]                                  │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

**Structure:**

```php
<div class="riseup-card">
    <div class="riseup-actions-row">
        <button class="button button-primary">
            <span class="dashicons dashicons-camera"></span>
            <?php esc_html_e('Primary Action', $pluginSlug); ?>
        </button>
        <button class="button button-secondary">
            <span class="dashicons dashicons-update"></span>
            <?php esc_html_e('Refresh', $pluginSlug); ?>
        </button>
        <span id="action-status" class="riseup-inline-status"></span>
    </div>

    <!-- Expandable options panel (toggle via JS) -->
    <div id="options_panel" style="display: none;">
        <!-- Form content -->
        <p>
            <button class="button button-primary">Confirm</button>
            <button class="button button-secondary">Cancel</button>
        </p>
    </div>
</div>
```

**Rules:**
1. Primary action is always `button-primary`, secondary actions are `button-secondary`
2. Every button with an icon uses a dashicon `<span>` inside the button
3. Inline status text (`riseup-inline-status`) appears after the last button
4. Expandable panels are separated by a border-top and use `riseup-snapshot-options` styling
5. Actions row uses `display: flex; gap: 10px; align-items: center; flex-wrap: wrap`

---

## 4. Filter Bar

### 4.1 Anatomy

```
┌─ .riseup-filters ────────────────────────────────────────┐
│  ┌─ .filter-row (primary) ─────────────────────────────┐ │
│  │ ACTION ▾  STATUS ▾  TRIGGER ▾  USER [___]  PLUGIN   │ │
│  └─────────────────────────────────────────────────────┘ │
│  ┌─ .filter-row-secondary ─────────────────────────────┐ │
│  │ FROM [date]  TO [date]  [🔵 Filter] [Reset]         │ │
│  └─────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────┘
```

### 4.2 Structure

```php
<div class="riseup-filters">
    <form method="get">
        <input type="hidden" name="page" value="<?php echo esc_attr(AdminPageType::Logs->value); ?>">
        
        <div class="filter-row">
            <label>
                <span><?php esc_html_e('Action:', $pluginSlug); ?></span>
                <select name="filter_action">
                    <option value=""><?php esc_html_e('All Actions', $pluginSlug); ?></option>
                    <?php foreach ($actionLabels as $key => $label): ?>
                        <option value="<?php echo esc_attr($key); ?>" <?php selected($filters['action'], $key); ?>>
                            <?php echo esc_html($label); ?>
                        </option>
                    <?php endforeach; ?>
                </select>
            </label>
            <!-- More filters... -->
        </div>
        
        <div class="filter-row filter-row-secondary">
            <label>
                <span><?php esc_html_e('From:', $pluginSlug); ?></span>
                <input type="date" name="filter_from" value="<?php echo esc_attr($filters['from']); ?>">
            </label>
            <button type="submit" class="button button-primary"><?php esc_html_e('Filter', $pluginSlug); ?></button>
            <a href="<?php echo esc_url($resetUrl); ?>" class="button"><?php esc_html_e('Reset', $pluginSlug); ?></a>
        </div>
    </form>
</div>
```

### 4.3 Rules

1. Each filter label has an uppercase `<span>` label above the input
2. `<select>` options always start with an "All" option (empty value)
3. Enum-driven filters iterate over label arrays built from enum values
4. Date filters use native `<input type="date">`
5. Primary row holds category filters, secondary row holds date range + submit
6. Filter form submits via GET with `page` parameter preserved
7. Reset link navigates to the page URL without query params
8. Focus-within state highlights the entire filter container

---

## 5. Table Patterns

### 5.1 Standard Table

All data tables use WordPress `wp-list-table widefat fixed striped`:

```php
<table class="wp-list-table widefat fixed striped">
    <thead>
        <tr>
            <th class="column-name"><?php esc_html_e('Name', $pluginSlug); ?></th>
            <!-- Column widths defined in page-specific CSS -->
        </tr>
    </thead>
    <tbody>
        <?php if (empty($items)): ?>
            <tr><td colspan="5" class="no-items"><?php esc_html_e('No items found.', $pluginSlug); ?></td></tr>
        <?php else: ?>
            <?php foreach ($items as $item): ?>
                <tr><!-- row content --></tr>
            <?php endforeach; ?>
        <?php endif; ?>
    </tbody>
</table>
```

### 5.2 Column Width Definitions

Column widths MUST be defined in page-specific CSS using `.column-{name}` classes:

```css
.riseup-admin.riseup-agents .column-name    { width: 20%; }
.riseup-admin.riseup-agents .column-url     { width: 30%; }
.riseup-admin.riseup-agents .column-status  { width: 10%; }
.riseup-admin.riseup-agents .column-actions { width: 25%; }
```

### 5.3 Date Group Headers

Tables with chronological data MUST insert date group separator rows when the date changes:

```php
<?php
$currentDateGroup = '';
foreach ($logs as $log):
    $logDate = DateHelper::formatDateOnly($logTimestamp);
    
    if ($logDate !== $currentDateGroup):
        $currentDateGroup = $logDate;
        $relativeDayKey = DateHelper::relativeDayKey($logTimestamp);
        // Build label: "Today — March 15, 2026" or just "March 14, 2026"
        if ($relativeDayKey === 'today') {
            $dateLabel = __('Today', $pluginSlug) . ' — ' . $logDateDisplay;
        } elseif ($relativeDayKey === 'yesterday') {
            $dateLabel = __('Yesterday', $pluginSlug) . ' — ' . $logDateDisplay;
        } else {
            $dateLabel = $logDateDisplay;
        }
?>
    <tr class="date-group-header">
        <td colspan="11">
            <span class="date-group-label"><?php echo esc_html($dateLabel); ?></span>
        </td>
    </tr>
<?php endif; ?>
    <tr><!-- normal data row --></tr>
<?php endforeach; ?>
```

**Visual:**
```
┌──────────────────────────────────────────────────────┐
│ 📅 Today — April 9, 2026                            │  ← gradient bg, accent border-top
├──────────────────────────────────────────────────────┤
│ Row data...                                          │
│ Row data...                                          │
├──────────────────────────────────────────────────────┤
│ 📅 Yesterday — April 8, 2026                        │
├──────────────────────────────────────────────────────┤
│ Row data...                                          │
└──────────────────────────────────────────────────────┘
```

### 5.4 Clickable Rows

Rows with expandable details use a `has-details` class and store data in `data-details`:

```php
<tr class="riseup-log-row <?php echo $hasDetails ? 'has-details' : ''; ?>"
    <?php if ($hasDetails): ?>
        data-details="<?php echo esc_attr(json_encode($details)); ?>"
    <?php endif; ?>>
```

Hover effect: `inset 3px 0 0 #667eea` left border accent + background tint.

### 5.5 Endpoint Group Headers

Tables displaying REST API endpoints group by category:

```php
<tr class="endpoint-group-header">
    <td colspan="4"><?php echo esc_html($groupName); ?></td>
</tr>
```

### 5.6 Nested Rows

Child rows (e.g., incremental snapshots under full backups) use:

```php
<tr class="riseup-nested-row">
    <td><!-- 3px purple left border via CSS --></td>
</tr>
```

---

## 6. Badge System Usage

### 6.1 Enum-Driven Badge Rendering

Badges MUST be driven by enum values mapped to CSS classes and label arrays:

```php
// Build mappings from enums
$triggerClasses = [
    TriggerSourceType::Api->value       => 'trigger-api',
    TriggerSourceType::Dashboard->value => 'trigger-dashboard',
    TriggerSourceType::Agent->value     => 'trigger-agent',
];
$triggerLabels = [
    TriggerSourceType::Api->value       => __('API', $pluginSlug),
    TriggerSourceType::Dashboard->value => __('Dashboard', $pluginSlug),
];

// Render
<span class="trigger-badge <?php echo esc_attr($triggerClass); ?>">
    <?php echo esc_html($triggerLabel); ?>
</span>
```

### 6.2 Badge Categories in Use

| Category | Class Pattern | Shape | Use Case |
|----------|--------------|-------|----------|
| **Status** | `.status-{Value}` | pill (20px radius) | Connection status, operation result |
| **Action** | `.action-{Value}` | tag (4px radius) | Log action type |
| **Trigger** | `.trigger-{source}` | tag (4px radius) | How an operation was initiated |
| **Upload Source** | `.source-{method}` | tag (4px radius) | Upload method badge |
| **HTTP Method** | `.method-{verb}` | tag (6px radius) | REST endpoint method |
| **Version** | `.version-{state}` | tag (4px radius) | Current vs old version |
| **Level** | `.level-badge` | pill (20px radius) | Error severity |
| **Count** | `.error-count-badge` | pill (12px radius) | Header error count |
| **Tab** | `.tab-badge` | pill (10px radius) | Tab notification count |

### 6.3 Missing Value Handling

When a value is empty or null, render a muted placeholder instead of a badge:

```php
<?php if ($hasValue): ?>
    <span class="badge"><?php echo esc_html($label); ?></span>
<?php else: ?>
    <span class="na">—</span>
<?php endif; ?>
```

---

## 7. Modal Anatomy

### 7.1 Modal Wrapper Partial

All modals MUST use the `modal-wrapper.php` partial for consistency:

```php
<?php
$modalId    = 'confirm-delete-modal';
$modalTitle = __('Confirm Deletion', $pluginSlug);
$modalIcon  = 'dashicons-warning';
$modalIconColor = '#d63638';
$modalBody  = '<p>' . __('Are you sure?', $pluginSlug) . '</p>';
$modalFooter = '<button class="button button-primary">Confirm</button>';
include __DIR__ . '/partials/shared/modal-wrapper.php';
?>
```

**Available variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `$modalId` | `'riseup-modal'` | HTML id attribute |
| `$modalTitle` | `''` | Modal heading text |
| `$modalIcon` | `''` | Dashicons class for header icon |
| `$modalIconColor` | `''` | Icon color override |
| `$modalMaxWidth` | `'600px'` | Content max-width |
| `$modalCloseButton` | `true` | Show × close button |
| `$modalHeaderExtra` | `''` | Raw HTML after title |
| `$modalBody` | `''` | Modal body content |
| `$modalFooter` | `''` | Footer/actions content |

### 7.2 Variable Cleanup

The partial MUST `unset()` all modal variables after rendering to prevent bleed into subsequent includes:

```php
unset($modalId, $modalTitle, $modalIcon, $modalIconColor, 
      $modalMaxWidth, $modalCloseButton, $modalHeaderExtra, 
      $modalBody, $modalFooter);
```

### 7.3 Modal Sizes

| Size | Max Width | Use Case |
|------|-----------|----------|
| Standard | `600px` | Confirmations, simple forms |
| Wide | `800px` | Agent plugins, complex forms |
| Fullscreen | `1000px` | Error detail with tabs |

### 7.4 Modal Rendering (Inline vs Partial)

Simple modals that don't need the partial (e.g., legacy templates) may use inline HTML but MUST follow the same class structure:

```php
<div id="my-modal" class="riseup-modal" style="display: none;">
    <div class="riseup-modal-overlay"></div>
    <div class="riseup-modal-content">
        <div class="riseup-modal-header">
            <div class="modal-header-left">
                <span class="dashicons dashicons-admin-plugins"></span>
                <h3>Title</h3>
            </div>
            <button type="button" class="riseup-modal-close">&times;</button>
        </div>
        <div class="riseup-modal-body"><!-- content --></div>
    </div>
</div>
```

---

## 8. Notice Patterns

### 8.1 WordPress Admin Notices

Standard WordPress notices for page-level feedback:

```php
<div class="notice notice-success is-dismissible">
    <p><?php esc_html_e('Settings saved successfully.', $pluginSlug); ?></p>
</div>
```

### 8.2 Flash Banner (Contextual Alert)

For important contextual alerts (e.g., "new errors detected since last visit"):

```php
<div class="riseup-flash-banner">
    <span class="flash-icon">⚠️</span>
    <div class="flash-content">
        <strong><?php esc_html_e('Alert message', $pluginSlug); ?></strong>
        <span class="flash-time"><?php echo esc_html($timeAgo); ?></span>
    </div>
    <button class="button flash-dismiss"><?php esc_html_e('Dismiss', $pluginSlug); ?></button>
</div>
```

### 8.3 Warning Card

For persistent warnings that require user action:

```php
<div class="notice notice-warning riseup-feedback-warning-card">
    <div class="riseup-feedback-warning-inner">
        <span class="dashicons dashicons-warning"></span>
        <div class="riseup-feedback-warning-content">
            <strong><?php esc_html_e('Warning title', $pluginSlug); ?></strong>
            <p><?php esc_html_e('Warning explanation.', $pluginSlug); ?></p>
            <div class="riseup-feedback-warning-actions">
                <button class="button button-primary">Action</button>
            </div>
        </div>
    </div>
</div>
```

### 8.4 Inline Status

For inline operation feedback next to buttons:

```php
<span id="status-id" class="riseup-inline-status"></span>
```

Set dynamically via JavaScript:
```javascript
statusEl.className = 'riseup-inline-status success';
statusEl.textContent = '✓ Saved';
// OR
statusEl.className = 'riseup-inline-status error';
statusEl.textContent = '✕ Failed';
```

---

## 9. Empty States

### 9.1 Table Empty State

When a table has no data, show a single-row message:

```php
<?php if (empty($items)): ?>
    <tr>
        <td colspan="<?php echo $columnCount; ?>" class="no-items">
            <?php esc_html_e('No items found.', $pluginSlug); ?>
        </td>
    </tr>
<?php endif; ?>
```

### 9.2 Section Empty State

When an entire card/section has no data, show a descriptive empty state:

```php
<div id="snapshots_empty" style="display: none;">
    <p><em><?php esc_html_e('No snapshots found. Click "Snapshot Now" to create your first backup.', $pluginSlug); ?></em></p>
</div>
```

### 9.3 File Empty State (Success)

When a file/log is empty (positive state — no errors):

```php
<div class="file-empty">
    <span class="dashicons dashicons-yes-alt"></span>
    <p><?php esc_html_e('No errors found — looking good!', $pluginSlug); ?></p>
</div>
```

The dashicon uses green color `#22c55e` with a pulse animation.

### 9.4 Empty State Rules

1. Empty states MUST provide context about what would appear and how to create it
2. Table empty states use `class="no-items"` inside a full-colspan `<td>`
3. Section empty states use `<em>` for visual distinction
4. Positive empty states (no errors) celebrate with a green icon
5. Empty states are controlled via `style="display: none;"` and toggled by JavaScript

---

## 10. Loading States

### 10.1 Section Loading

Use the WordPress spinner for section-level loading:

```php
<div id="section_loading" style="display: none;">
    <span class="spinner is-active" style="float: none;"></span>
    <?php esc_html_e('Loading...', $pluginSlug); ?>
</div>
```

### 10.2 Button Loading

Buttons show a spinning dashicon during async operations:

```php
<button id="my-btn" class="button button-primary">
    <span class="dashicons dashicons-update"></span>
    <?php esc_html_e('Save', $pluginSlug); ?>
</button>
```

JavaScript toggles the spin class:
```javascript
btn.querySelector('.dashicons').classList.add('spin');
btn.disabled = true;
// After completion:
btn.querySelector('.dashicons').classList.remove('spin');
btn.disabled = false;
```

### 10.3 File Loading

For async-loaded content panels:

```php
<div class="file-loading">
    <?php esc_html_e('Loading file contents...', $pluginSlug); ?>
</div>
```

Uses `riseupPulse` animation for subtle breathing effect.

### 10.4 Loading State Rules

1. Loading indicators are hidden by default (`display: none`)
2. JavaScript shows the loader and hides the content container
3. After data loads, hide the loader and show the content (or empty state)
4. Three-state pattern: `loading → content OR empty`
5. The WordPress spinner class `spinner is-active` with `float: none` is the standard

### 10.5 Three-State Container Pattern

```php
<!-- State 1: Loading -->
<div id="section_loading">
    <span class="spinner is-active" style="float: none;"></span>
    Loading...
</div>

<!-- State 2: Content (hidden until loaded) -->
<div id="section_content" style="display: none;">
    <!-- actual content -->
</div>

<!-- State 3: Empty (hidden unless no data) -->
<div id="section_empty" style="display: none;">
    <p><em>No data available.</em></p>
</div>
```

JavaScript flow:
```javascript
// Show loading
show('section_loading'); hide('section_content'); hide('section_empty');

// After fetch:
hide('section_loading');
if (data.length > 0) {
    show('section_content');
} else {
    show('section_empty');
}
```

---

## 11. Stats Bar

### 11.1 Inline Stats

Simple stats shown between filter bar and table:

```php
<div class="riseup-stats">
    <span class="stat-item">
        <strong><?php echo esc_html($total); ?></strong>
        <?php esc_html_e('total records', $pluginSlug); ?>
    </span>
    <span class="stat-item">
        <?php esc_html_e('Page', $pluginSlug); ?> <?php echo esc_html($page); ?> 
        <?php esc_html_e('of', $pluginSlug); ?> <?php echo esc_html($totalPages); ?>
    </span>
</div>
```

### 11.2 Stat Cards

Rich stat displays with hero numbers:

```php
<div class="riseup-analytics-summary">
    <div class="riseup-stat-card">
        <span class="riseup-stat-value" id="stat_total">—</span>
        <span class="riseup-stat-label"><?php esc_html_e('Total Size', $pluginSlug); ?></span>
    </div>
    <!-- More stat cards -->
</div>
```

**Rules:**
1. Stat cards use `display: flex; gap: 14px; flex-wrap: wrap`
2. Each card has `flex: 1; min-width: 90px`
3. Value is monospace, 20px, bold, primary color
4. Label is uppercase, 11px, 0.5px letter-spacing, muted color
5. Default value is `—` (em dash) before data loads
6. Cards elevate on hover with `translateY(-1px)` and value scales `1.08`

---

## 12. Progress Panel

For long-running operations, show a progress panel:

```
┌─ .riseup-card ────────────────────────────────────────┐
│  ⚡ Snapshot In Progress  [0%]                        │
│  ┌─ progress bar ────────────────────────────────┐    │
│  │ ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │    │
│  └───────────────────────────────────────────────┘    │
│  Processing table wp_posts... (4 of 12)               │
│                                                        │
│  Table Progress:                                       │
│  ✓ wp_options  ✓ wp_users  ⟳ wp_posts  ○ wp_comments  │
└────────────────────────────────────────────────────────┘
```

**Structure:**

```php
<div id="progress_panel" class="riseup-card" style="display: none;">
    <h2>
        <span class="dashicons dashicons-performance"></span>
        <?php esc_html_e('Operation In Progress', $pluginSlug); ?>
        <span id="progress_percent_badge" class="riseup-badge">0%</span>
    </h2>
    <div class="riseup-progress-bar-wrap">
        <div id="progress_bar" class="riseup-progress-bar" style="width: 0%;"></div>
    </div>
    <div id="progress_meta" class="riseup-progress-meta"></div>
    <div id="progress_tables" class="riseup-progress-tables" style="display: none;">
        <h4>Table Progress</h4>
        <div id="progress_tables_list"></div>
    </div>
</div>
```

---

## 13. Form Patterns

### 13.1 Settings Form (form-table)

WordPress settings forms use the `.form-table` pattern:

```php
<form id="settings-form">
    <table class="form-table">
        <tr>
            <th scope="row">
                <label for="field_id">
                    <?php esc_html_e('Field Label', $pluginSlug); ?>
                    <span class="required">*</span>
                </label>
            </th>
            <td>
                <input type="text" id="field_id" class="regular-text" required>
                <p class="description"><?php esc_html_e('Help text.', $pluginSlug); ?></p>
            </td>
        </tr>
    </table>
    <p class="submit">
        <button type="submit" class="button button-primary">
            <span class="dashicons dashicons-yes"></span>
            <?php esc_html_e('Save', $pluginSlug); ?>
        </button>
        <span id="form-status" class="riseup-inline-status"></span>
    </p>
</form>
```

### 13.2 Selection Cards (Radio Alternative)

For mutually exclusive options, use visual selection cards instead of radio buttons:

```php
<div class="riseup-storage-cards">
    <label class="riseup-storage-card" data-mode="option_a">
        <input type="radio" name="setting" value="option_a">
        <div class="riseup-storage-card-inner">
            <span class="dashicons dashicons-media-archive"></span>
            <strong>Option A</strong>
            <span class="description">Short description</span>
        </div>
    </label>
    <label class="riseup-storage-card active" data-mode="option_b">
        <input type="radio" name="setting" value="option_b" checked>
        <div class="riseup-storage-card-inner">
            <span class="dashicons dashicons-grid-view"></span>
            <strong>Option B</strong>
            <span class="description">Short description</span>
        </div>
    </label>
</div>
```

**Rules:**
1. Radio inputs are hidden (`display: none`)
2. The `.active` class is toggled via JavaScript on selection
3. Active card shows primary border + primary background tint
4. Cards use `flex: 1; min-width: 160px; max-width: 220px`

### 13.3 Slider Input

For numeric range inputs:

```php
<div class="riseup-slider-row">
    <input type="range" id="setting_value" min="1" max="10" value="5" class="riseup-range-slider">
    <span id="value_display" class="riseup-slider-value">5</span>
</div>
```

### 13.4 Conditional Fields

Fields that depend on other field values use `display: none` toggled by JavaScript:

```php
<tr id="conditional_row" style="display: none;">
    <th scope="row"><label>Conditional Field</label></th>
    <td><!-- content --></td>
</tr>
```

### 13.5 Form Section Dividers

When a form has multiple logical sections, use an `<h3>` with border-top:

```php
<h3 style="margin-top: 25px; padding-top: 15px; border-top: 1px solid #eee;">
    <span class="dashicons dashicons-performance"></span>
    <?php esc_html_e('Section Title', $pluginSlug); ?>
</h3>
```

---

### 13.6 Security — Admin POST Endpoints (🔴 CODE RED)

> **Severity:** 🔴 CODE RED — non-negotiable. Mirrors and extends
> [Phase 14 §14.5.2](./14-rest-api-conventions.md#1452-security--nonce-verification-for-mutating-routes--code-red)
> for non-REST admin surfaces (`admin-post.php`, `admin-ajax.php`,
> Settings API, custom `<form>` POSTs).

Every admin form that POSTs to WordPress MUST satisfy **all** of the
following before any side effect runs (writing options, deleting rows,
enqueuing background jobs, calling external APIs, etc.). Each
requirement is enforced by a CI linter rule.

| # | Requirement | Linter rule |
|---|-------------|-------------|
| 1 | The form template MUST emit a nonce field via `wp_nonce_field( $action, $field )` — never a hand-rolled hidden input. The `$action` MUST come from an `AdminActionType::nonceAction()` enum method. | `WP-NONCE-002` |
| 2 | The handler MUST call `check_admin_referer( $action, $field )` (or `check_ajax_referer()` for AJAX) as its **first executable line**, before any DB write, option update, or external call. | `WP-NONCE-001` |
| 3 | The handler MUST verify capability with `current_user_can( CapabilityType::… )` after the nonce check. `is_user_logged_in()` alone is insufficient. | `WP-PERM-001` |
| 4 | All `$_POST` / `$_GET` reads MUST be wrapped in `wp_unslash()` then a typed sanitizer (`sanitize_text_field`, `absint`, `sanitize_email`, `wp_kses_post`, …). Raw `$_POST['x']` is forbidden. | `WP-INPUT-001` |
| 5 | The form MUST POST to `admin-post.php` (or `admin-ajax.php`) with a hidden `action` field whose value is an `AdminActionType` enum case — never an inline string. | `WP-ACTION-001` |
| 6 | A failed nonce MUST short-circuit with `wp_die( __( 'Security check failed.', $slug ), 403 )`. Never redirect on a failed nonce — that masks tampering. | `WP-NONCE-001` |
| 7 | After a successful save, the handler MUST `wp_safe_redirect()` to a whitelisted admin URL with a `_wpnonce`-protected query arg, then `exit;`. Bare `header( "Location: …" )` is forbidden. | `WP-REDIRECT-001` |

#### Canonical form template

```php
<form method="post" action="<?php echo esc_url( admin_url('admin-post.php') ); ?>">
    <input type="hidden" name="action"
           value="<?php echo esc_attr( AdminActionType::SaveSettings->value ); ?>">
    <?php
    // Rule 1 — enum-driven nonce field.
    wp_nonce_field(
        AdminActionType::SaveSettings->nonceAction(),
        AdminActionType::SaveSettings->nonceField()
    );
    ?>
    <table class="form-table">
        <!-- … fields … -->
    </table>
    <?php submit_button( __( 'Save', $pluginSlug ) ); ?>
</form>
```

#### Canonical handler

```php
add_action(
    'admin_post_' . AdminActionType::SaveSettings->value,
    [ $this, 'handleSaveSettings' ]
);

public function handleSaveSettings(): void
{
    // Rule 2 — nonce FIRST.
    check_admin_referer(
        AdminActionType::SaveSettings->nonceAction(),
        AdminActionType::SaveSettings->nonceField()
    );

    // Rule 3 — capability SECOND.
    if ( ! current_user_can( CapabilityType::ManagePluginSettings->value ) ) {
        wp_die( __( 'Forbidden.', $pluginSlug ), 403 );
    }

    // Rule 4 — sanitize every input.
    $email = sanitize_email( wp_unslash( $_POST['admin_email'] ?? '' ) );
    $limit = absint( wp_unslash( $_POST['daily_limit'] ?? 0 ) );

    // … persist via Options API or micro-ORM …

    // Rule 7 — safe redirect with success flag.
    wp_safe_redirect( add_query_arg(
        [ 'page' => $pluginSlug, 'updated' => 'true' ],
        admin_url( 'admin.php' )
    ) );
    exit;
}
```

#### Forbidden patterns (will fail CI)

```php
// ❌ Rule 1 — hand-rolled hidden input instead of wp_nonce_field()
<input type="hidden" name="my_nonce" value="<?php echo wp_create_nonce('save'); ?>">

// ❌ Rule 2 — DB write before nonce check
update_option( 'my_setting', $_POST['value'] );
check_admin_referer( 'save_settings' );

// ❌ Rule 3 — logged-in is not authorized
if ( ! is_user_logged_in() ) { wp_die(); }

// ❌ Rule 4 — raw superglobal
$email = $_POST['admin_email'];

// ❌ Rule 5 — inline action string
<input type="hidden" name="action" value="save_my_settings">

// ❌ Rule 6 — silent redirect on failed nonce
if ( ! wp_verify_nonce( $_POST['_wpnonce'], 'save' ) ) {
    wp_redirect( admin_url() ); exit;
}

// ❌ Rule 7 — bare Location header bypasses safe-redirect allow-list
header( 'Location: ' . $_POST['return_to'] ); exit;
```

#### Cross-references

- Phase 14 §14.5.1 — REST permission callbacks (forbids `__return_true`).
- Phase 14 §14.5.2 — REST nonce verification (`wp_rest`).
- Phase 8 — WordPress integration patterns (`wp_nonce_field` examples).
- Phase 16 — Error handling extraction (AJAX nonce action enum pattern).

---

## 14. Tab System

### 14.1 Page-Level Tabs

```php
<div class="riseup-tabs">
    <a href="#tab-errors" class="nav-tab nav-tab-active" data-tab="errors">
        <span class="dashicons dashicons-warning"></span>
        <?php esc_html_e('Errors', $pluginSlug); ?>
        <span class="tab-badge"><?php echo esc_html($errorCount); ?></span>
    </a>
    <a href="#tab-file" class="nav-tab" data-tab="file">
        <span class="dashicons dashicons-media-text"></span>
        <?php esc_html_e('File Viewer', $pluginSlug); ?>
    </a>
</div>

<div id="tab-errors" class="riseup-tab-content"><!-- content --></div>
<div id="tab-file" class="riseup-tab-content" style="display: none;"><!-- content --></div>
```

### 14.2 Modal Tabs

```php
<div class="modal-tabs">
    <button class="modal-tab active" data-tab="context">
        <span class="dashicons dashicons-info"></span>
        Context
    </button>
    <button class="modal-tab" data-tab="stack">
        <span class="dashicons dashicons-editor-code"></span>
        Stack Trace
    </button>
</div>
<div class="modal-tab-pane" id="pane-context"><!-- content --></div>
<div class="modal-tab-pane" id="pane-stack" style="display: none;"><!-- content --></div>
```

---

## 15. Pagination

### 15.1 Pagination Partial

Use the shared pagination partial for all paginated content:

```php
<?php
$totalPages = $totalPages;
$page = $currentPage;
include __DIR__ . '/partials/shared/pagination.php';
?>
```

The partial uses WordPress `paginate_links()` and renders nothing if `$totalPages <= 1`.

### 15.2 JavaScript Pagination

For AJAX-loaded tables, pagination links are rendered into a container:

```php
<div id="section_pagination" class="tablenav bottom" style="display: none;">
    <div class="tablenav-pages">
        <span class="displaying-num" id="section_count"></span>
        <span class="pagination-links" id="section_pages"></span>
    </div>
</div>
```

---

## 16. Chart & Calendar Components

### 16.1 Bar Chart

```
┌─ Chart Container ──────────────────────────────┐
│  Y-axis │ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██    │
│  labels │ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██    │
│         └──────────────────────────────────     │
│           x-axis labels (rotated -45°)          │
│                                                  │
│  ● Full  ● Incremental                          │
└──────────────────────────────────────────────────┘
```

### 16.2 Calendar View

```
┌─ Calendar Card ───────────────────────────────┐
│  [◀]     April 2026     [▶]                   │
│  Sun Mon Tue Wed Thu Fri Sat                   │
│  ... ... ... 1   2   3   4                     │
│  5   6   7   8•  9   10  11                    │
│  ...                                           │
│  ● Full  ● Incremental  ● Scheduled            │
└───────────────────────────────────────────────┘
```

### 16.3 Legend

```php
<div class="riseup-chart-legend">
    <span class="riseup-legend-item">
        <span class="riseup-legend-dot" style="background:#2271b1;"></span>
        <?php esc_html_e('Full', $pluginSlug); ?>
    </span>
</div>
```

---

## 17. Auto-Refresh Toggle

For real-time monitoring pages:

```php
<label class="auto-refresh-toggle">
    <input type="checkbox" id="auto_refresh_toggle">
    <span class="auto-refresh-label">
        <span class="live-dot" id="live_indicator"></span>
        <?php esc_html_e('Auto-refresh', $pluginSlug); ?>
    </span>
</label>
```

The live dot pulses green when active, stays gray when inactive.

---

## 18. Template Composition Rules

### 18.1 Partial Directory Structure

```
templates/
├── admin-{page}.php              # Page orchestrator
├── partials/
│   ├── shared/                   # Cross-page partials
│   │   ├── page-header.php
│   │   ├── modal-wrapper.php
│   │   └── pagination.php
│   ├── {page}/                   # Page-specific partials
│   │   └── modals.php
│   └── settings/                 # Settings section partials
│       ├── section-log-retrieval.php
│       └── section-snapshot-settings.php
└── unused/                       # Archived templates
```

### 18.2 When to Extract a Partial

Extract a partial when:
1. The same HTML block appears on 2+ pages
2. A single template file exceeds ~300 lines
3. A logical section (modals, settings group) is self-contained
4. A component needs consistent rendering (page header, pagination)

### 18.3 Partial Variable Contract

Every partial MUST document its required and optional variables in a docblock:

```php
/**
 * Required variables (from parent scope):
 *   $pageIcon   — Dashicons class
 *   $pageTitle  — Translated title
 *   $pluginSlug — Text domain
 *
 * Optional variables:
 *   $pageDescription — Translated description
 *   $headerExtra     — Raw HTML after badge
 */
```

### 18.4 Variable Cleanup

Partials that define temporary variables MUST `unset()` them after rendering to prevent scope bleed into subsequent includes.

---

## 19. Accessibility Rules

1. All `<th>` cells MUST have `scope="row"` or `scope="col"`
2. All `<label>` elements MUST have a `for` attribute matching an input `id`
3. All `<input>` elements MUST have a `type` attribute
4. Required fields MUST show `<span class="required">*</span>` in the label
5. All dashicons used as icons MUST be inside a `<span>` (not standalone)
6. Interactive elements MUST be `<button>` or `<a>`, never `<div>` or `<span>`
7. Modal close buttons use `&times;` character, not an icon font

---

## 20. Anti-Patterns (NEVER DO)

1. ❌ Use `<table>` for non-tabular layout
2. ❌ Place modals inside scrollable containers
3. ❌ Use inline `onclick` handlers — always use event delegation
4. ❌ Hardcode column counts in empty state `colspan` — use a variable
5. ❌ Skip the three-state pattern (loading/content/empty) for async sections
6. ❌ Use `display: block` to show elements — use `display: ''` (restores original)
7. ❌ Render user input without `esc_html()` or `esc_attr()`
8. ❌ Duplicate modal HTML across pages — use the shared partial
9. ❌ Create templates > 400 lines without extracting partials
10. ❌ Use `echo` for translatable strings — always use `esc_html_e()` or `esc_html(__(...))`
