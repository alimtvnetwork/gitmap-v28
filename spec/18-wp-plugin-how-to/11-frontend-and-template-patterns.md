# Phase 11 — Frontend & Template Patterns

> **Purpose:** Define how to write PHP templates, compose them from reusable partials, enforce file size limits, and optionally integrate React.js with Tailwind CSS for modern admin UIs. Covers asset enqueuing, source maps, and production builds.
> **Audience:** AI code generators and human developers.
> **Prerequisite:** Phases 1–10 must be read first.

---

## 11.1 File Size Limits

Every file in the plugin — PHP classes, traits, templates, partials, JS, CSS — must respect strict line limits.

| File type | Ideal | Maximum | Action when exceeded |
|-----------|-------|---------|---------------------|
| PHP template (page) | 50–100 | 200 | Extract sections into `partials/` |
| PHP partial | 30–50 | 100 | Split into smaller partials |
| PHP class / trait | 50–100 | 200 | Decompose into additional traits (Phase 3) |
| JavaScript file | 50–100 | 200 | Extract into modules |
| CSS file | 50–100 | 200 | Split by component or section |

**Rule:** If a file exceeds 200 lines, it **must** be refactored before merging. No exceptions.

### Why this matters

- Files under 100 lines are easier to review, test, and debug
- Smaller files reduce merge conflicts
- AI code generators produce better output when working with focused, single-purpose files
- Long templates are a sign of missing abstraction

---

## 11.2 Template Architecture

### Directory structure

```
plugin-slug/
├── templates/
│   ├── admin-settings.php           ← Page template (orchestrator)
│   ├── admin-logs.php               ← Page template
│   ├── admin-agents.php             ← Page template
│   └── partials/
│       ├── shared/                   ← Cross-page reusable partials
│       │   ├── page-header.php       ← Standard header with icon + title + version
│       │   ├── pagination.php        ← Pagination controls
│       │   ├── modal-wrapper.php     ← Modal shell (content injected)
│       │   ├── notice.php            ← Admin notice partial
│       │   └── empty-state.php       ← "No data" placeholder
│       ├── settings/                 ← Settings-page-specific partials
│       │   ├── section-general.php
│       │   ├── section-advanced.php
│       │   └── section-update.php
│       ├── logs/                     ← Logs-page-specific partials
│       │   ├── log-filters.php
│       │   └── log-table.php
│       └── agents/                   ← Agents-page-specific partials
│           ├── agent-form.php
│           └── agent-list.php
```

### Naming conventions

| Convention | Example | Rule |
|-----------|---------|------|
| Page templates | `admin-{page}.php` | One per admin page |
| Shared partials | `partials/shared/{name}.php` | Reused across 2+ pages |
| Page partials | `partials/{page}/{name}.php` | Specific to one page |
| Prefix pattern | `section-`, `form-`, `list-`, `table-`, `modal-` | Describes the UI element |

---

## 11.3 Page Templates — The Orchestrator Pattern

A page template is an **orchestrator**: it sets variables and includes partials. It should contain **minimal HTML** itself.

### ✅ Correct — Orchestrator pattern (under 100 lines)

```php
<?php
/**
 * Admin Settings Page Template
 *
 * Orchestrates settings sections via partials.
 * Each section is a self-contained partial under partials/settings/.
 *
 * @package PluginName
 * @since   1.0.0
 */

use PluginName\Enums\PluginConfigType;

if (!defined('ABSPATH')) {
    exit;
}

$pluginName = PluginConfigType::Name->value;
$pluginSlug = PluginConfigType::Slug->value;
?>
<div class="wrap plugin-admin">
    <?php
    // ── Page Header (shared partial) ────────────────────────
    $pageIcon = 'dashicons-admin-generic';
    $pageTitle = $pluginName . ' - ' . __('Settings', $pluginSlug);
    $pageDescription = __('Configure plugin behaviour.', $pluginSlug);
    include __DIR__ . '/partials/shared/page-header.php';
    ?>

    <!-- General Settings Section -->
    <?php include __DIR__ . '/partials/settings/section-general.php'; ?>

    <!-- Advanced Settings Section -->
    <?php include __DIR__ . '/partials/settings/section-advanced.php'; ?>

    <!-- Update Settings Section -->
    <?php include __DIR__ . '/partials/settings/section-update.php'; ?>
</div>
```

### ❌ Wrong — Monolithic template (400+ lines)

```php
<!-- DO NOT do this — everything inline in one file -->
<div class="wrap">
    <h1>Settings</h1>
    <!-- 400 lines of mixed HTML, PHP logic, and inline styles -->
</div>
```

---

## 11.4 Partial Templates — The Component Pattern

Each partial is a **self-contained UI component** that receives data via PHP variables set in the parent scope.

### Partial anatomy

```php
<?php
/**
 * Partial: Section Name
 *
 * Renders the [description of what this partial renders].
 *
 * Required variables (from parent scope):
 *   $pluginSlug — Plugin text domain for i18n
 *
 * Optional variables:
 *   $sectionTitle — Override default title
 *
 * @package PluginName
 * @since   1.0.0
 */

if (!defined('ABSPATH')) {
    exit;
}
?>
<div class="plugin-card">
    <h2>
        <span class="dashicons dashicons-admin-tools"></span>
        <?php echo esc_html($sectionTitle ?? __('Default Title', $pluginSlug)); ?>
    </h2>

    <table class="form-table">
        <tr>
            <th scope="row">
                <label for="setting_name">
                    <?php esc_html_e('Setting Label', $pluginSlug); ?>
                </label>
            </th>
            <td>
                <input type="text"
                       id="setting_name"
                       name="setting_name"
                       class="regular-text"
                       value="<?php echo esc_attr($settingValue ?? ''); ?>">
            </td>
        </tr>
    </table>
</div>
```

### Partial rules

| Rule | Detail |
|------|--------|
| **Document required variables** | Every partial must list which `$variables` it expects from parent scope |
| **ABSPATH guard** | Every partial starts with `if (!defined('ABSPATH')) { exit; }` |
| **Escape all output** | Use `esc_html()`, `esc_attr()`, `esc_url()`, `wp_kses()` — never raw `echo` |
| **No business logic** | Partials render data — they don't fetch, compute, or mutate |
| **No direct DB queries** | Data must be passed in via variables, not queried inside the partial |
| **i18n all strings** | Use `__()`, `_e()`, `esc_html_e()`, `esc_attr_e()` for all user-facing text |
| **Max 100 lines** | If a partial exceeds 100 lines, split it further |

### Shared partial reuse

When the same UI element appears on 2+ pages, extract it to `partials/shared/`:

```php
// In any page template — reuse the shared page header
$pageIcon = 'dashicons-database';
$pageTitle = $pluginName . ' - ' . __('Logs', $pluginSlug);
$pageDescription = __('View transaction logs.', $pluginSlug);
include __DIR__ . '/partials/shared/page-header.php';
```

### Complete Partial Example — Data Flow from Orchestrator to Partial

This example shows the full lifecycle: orchestrator prepares data, partial receives and renders it.

**Orchestrator:** `templates/admin-agents.php`

```php
<?php
/**
 * Admin Agents Page — Manages remote site agents.
 *
 * @package PluginName
 * @since   1.0.0
 */

use PluginName\Enums\PluginConfigType;

if (!defined('ABSPATH')) {
    exit;
}

$pluginName = PluginConfigType::Name->value;
$pluginSlug = PluginConfigType::Slug->value;

// ── Prepare data for partials ──────────────────────────────
$agents = $this->agentService->getAllAgents();
$agentCount = count($agents);
$hasAgents = ($agentCount > 0);
$statusSummary = $this->agentService->getStatusSummary();
?>
<div class="wrap pluginname-admin pluginname-agents">
    <?php
    // ── Shared header partial ───────────────────────────────
    $pageIcon = 'dashicons-networking';
    $pageTitle = $pluginName . ' - ' . __('Agents', $pluginSlug);
    $pageDescription = sprintf(
        __('%d agent(s) registered.', $pluginSlug),
        $agentCount
    );
    include __DIR__ . '/partials/shared/page-header.php';
    ?>

    <!-- Agent List (receives $agents, $hasAgents from parent scope) -->
    <?php include __DIR__ . '/partials/agents/agent-list.php'; ?>

    <!-- Add Agent Form -->
    <?php include __DIR__ . '/partials/agents/agent-form.php'; ?>
</div>
```

**Partial:** `templates/partials/agents/agent-list.php`

```php
<?php
/**
 * Partial: Agent List Table
 *
 * Renders a table of registered remote site agents.
 *
 * Required variables (from parent scope):
 *   $pluginSlug  — Plugin text domain for i18n
 *   $agents      — array<int, array{id: int, url: string, label: string, status: string}>
 *   $hasAgents   — bool — whether any agents exist
 *
 * @package PluginName
 * @since   1.0.0
 */

if (!defined('ABSPATH')) {
    exit;
}
?>
<div class="pluginname-card">
    <h2>
        <span class="dashicons dashicons-admin-site-alt3"></span>
        <?php esc_html_e('Registered Agents', $pluginSlug); ?>
    </h2>

    <?php if (!$hasAgents) : ?>
        <?php
        $emptyIcon = 'dashicons-networking';
        $emptyMessage = __('No agents registered yet. Add one below.', $pluginSlug);
        include __DIR__ . '/../shared/empty-state.php';
        ?>
    <?php else : ?>
        <table class="wp-list-table widefat fixed striped">
            <thead>
                <tr>
                    <th><?php esc_html_e('Label', $pluginSlug); ?></th>
                    <th><?php esc_html_e('URL', $pluginSlug); ?></th>
                    <th><?php esc_html_e('Status', $pluginSlug); ?></th>
                    <th><?php esc_html_e('Actions', $pluginSlug); ?></th>
                </tr>
            </thead>
            <tbody>
                <?php foreach ($agents as $agent) : ?>
                    <tr>
                        <td><?php echo esc_html($agent['label']); ?></td>
                        <td><code><?php echo esc_url($agent['url']); ?></code></td>
                        <td>
                            <span class="pluginname-status-badge pluginname-status-<?php echo esc_attr($agent['status']); ?>">
                                <?php echo esc_html(ucfirst($agent['status'])); ?>
                            </span>
                        </td>
                        <td>
                            <button class="button button-small btn-test-agent"
                                    data-agent-id="<?php echo esc_attr($agent['id']); ?>">
                                <?php esc_html_e('Test', $pluginSlug); ?>
                            </button>
                        </td>
                    </tr>
                <?php endforeach; ?>
            </tbody>
        </table>
    <?php endif; ?>
</div>
```

### Key takeaways from this example

| Pattern | Demonstrated |
|---------|-------------|
| **Data preparation in orchestrator** | `$agents`, `$hasAgents`, `$agentCount` computed before includes |
| **Variable documentation** | Partial header lists all required `$variables` with types |
| **Empty state delegation** | Empty state uses another shared partial (`empty-state.php`) |
| **Output escaping** | Every dynamic value uses `esc_html()`, `esc_attr()`, `esc_url()` |
| **No business logic in partial** | No DB queries, no service calls — pure rendering |
| **Nested partial reuse** | `agent-list.php` includes `shared/empty-state.php` when needed |

---

## 11.5 When to Extract a Partial

| Signal | Action |
|--------|--------|
| Template exceeds 100 lines | **Must** extract sections into partials |
| Same HTML block appears in 2+ templates | **Must** extract to `partials/shared/` |
| A `<div class="card">` or `<section>` block is self-contained | **Should** extract — it's a natural component boundary |
| A form has 5+ fields | **Should** extract the form to its own partial |
| A table with custom rendering logic | **Should** extract to `partials/{page}/table-{name}.php` |
| A modal or dialog | **Must** extract — modals are always reusable |

---

## 11.6 JavaScript and CSS — Traditional (Non-React) Pattern

### File structure

```
plugin-slug/
├── assets/
│   ├── css/
│   │   ├── admin-shared.css          ← Shared styles (loaded on all admin pages)
│   │   ├── admin-settings.css        ← Settings-page-specific styles
│   │   └── admin-logs.css            ← Logs-page-specific styles
│   └── js/
│       ├── admin-settings.js         ← Settings-page-specific behaviour
│       └── admin-logs.js             ← Logs-page-specific behaviour
```

### Enqueuing pattern (in AdminPageTrait)

```php
public function enqueueAdminAssets(string $hook): void
{
    $isPluginPage = $this->isPluginAdminPage($hook);

    if (!$isPluginPage) {
        return;
    }

    $pluginSlug = PluginConfigType::Slug->value;
    $version = PluginConfigType::Version->value;
    $baseUrl = plugin_dir_url(dirname(__DIR__));

    // Shared styles — all plugin admin pages
    wp_enqueue_style(
        $pluginSlug . '-admin-shared',
        $baseUrl . 'assets/css/admin-shared.css',
        [],
        $version,
    );

    // Page-specific assets
    $page = $this->getCurrentAdminPage($hook);
    $cssFile = 'assets/css/admin-' . $page . '.css';
    $jsFile = 'assets/js/admin-' . $page . '.js';

    $hasCss = file_exists(plugin_dir_path(dirname(__DIR__)) . $cssFile);

    if ($hasCss) {
        wp_enqueue_style(
            $pluginSlug . '-admin-' . $page,
            $baseUrl . $cssFile,
            [$pluginSlug . '-admin-shared'],
            $version,
        );
    }

    $hasJs = file_exists(plugin_dir_path(dirname(__DIR__)) . $jsFile);

    if ($hasJs) {
        wp_enqueue_script(
            $pluginSlug . '-admin-' . $page,
            $baseUrl . $jsFile,
            ['jquery'],
            $version,
            true,  // Load in footer
        );

        // Pass data to JS
        wp_localize_script($pluginSlug . '-admin-' . $page, 'PluginNameAdmin', [
            'ajaxUrl' => admin_url('admin-ajax.php'),
            'restUrl' => rest_url(PluginConfigType::apiFullNamespace() . '/'),
            'nonce'   => wp_create_nonce('wp_rest'),
            'slug'    => $pluginSlug,
        ]);
    }
}
```

---

## 11.7 React.js Integration (Optional — Developer Confirmation Required)

> **⚠️ Important:** Using React + Tailwind for a WordPress plugin admin UI is a valid architectural choice but must be **explicitly confirmed by the developer** before implementation. It adds build tooling complexity (Node.js, Vite/Webpack, npm scripts) that not all teams are prepared to maintain.

### When to use React

| Use React when | Stick with PHP templates when |
|----------------|------------------------------|
| Admin UI has complex interactive state (drag-and-drop, real-time updates) | Pages are mostly forms and tables |
| Multiple components share state (dashboard with linked widgets) | Each section is independent |
| Developer team has React experience | Team is PHP-only |
| Plugin will have a standalone SPA-like admin experience | Standard WordPress admin look-and-feel is preferred |

### Directory structure (React mode)

```
plugin-slug/
├── frontend/                        ← React source (NOT shipped in ZIP)
│   ├── src/
│   │   ├── main.tsx                  ← Entry point
│   │   ├── App.tsx                   ← Root component
│   │   ├── components/
│   │   │   ├── settings/
│   │   │   │   ├── SettingsPage.tsx
│   │   │   │   ├── GeneralSection.tsx
│   │   │   │   └── UpdateSection.tsx
│   │   │   └── shared/
│   │   │       ├── PageHeader.tsx
│   │   │       ├── Card.tsx
│   │   │       └── Modal.tsx
│   │   ├── hooks/
│   │   │   └── useApi.ts
│   │   ├── types/
│   │   │   └── index.ts
│   │   └── styles/
│   │       └── index.css             ← Tailwind entry
│   ├── index.html                    ← Dev server entry
│   ├── vite.config.ts                ← Build configuration
│   ├── tailwind.config.ts
│   ├── postcss.config.js
│   ├── tsconfig.json
│   └── package.json
├── assets/
│   └── dist/                         ← Built output (shipped in ZIP)
│       ├── admin.js                  ← Bundled React app
│       ├── admin.css                 ← Bundled Tailwind CSS
│       ├── admin.js.map              ← Source map (DEV ONLY — see §11.8)
│       └── admin.css.map             ← Source map (DEV ONLY — see §11.8)
├── includes/
├── templates/
│   └── admin-react-root.php          ← Minimal mount-point template
└── .distignore                       ← Must exclude frontend/ source
```

### Mount-point template

When using React, the PHP template is minimal — just a mount point:

```php
<?php
/**
 * Admin React App Mount Point
 *
 * Renders the #plugin-root div where the React application mounts.
 * All UI is handled by React — this template only provides the container.
 *
 * @package PluginName
 * @since   1.0.0
 */

if (!defined('ABSPATH')) {
    exit;
}
?>
<div class="wrap">
    <div id="plugin-name-root"></div>
</div>
```

### React component size rules

The same file size limits apply to React components:

| File type | Ideal | Maximum |
|-----------|-------|---------|
| Page component (e.g. `SettingsPage.tsx`) | 50–80 lines | 150 |
| UI component (e.g. `Card.tsx`, `Modal.tsx`) | 20–50 lines | 100 |
| Custom hook (e.g. `useApi.ts`) | 30–60 lines | 100 |
| Type definitions | 20–40 lines | 100 |

**Rule:** If a React component exceeds 100 lines, extract sub-components or custom hooks.

---

## 11.8 Source Maps — Dev vs. Production

Source maps **must** be included in development builds and **must NOT** be included in production/distribution builds.

### Why

| Environment | Source maps? | Reason |
|-------------|-------------|--------|
| **Development** | ✅ Yes | Enables debugging in browser DevTools with original source |
| **Production** | ❌ No | Prevents exposing source code, reduces file size, improves security |

### Vite configuration

```typescript
// frontend/vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig(({ mode }) => {
    const isDev = mode === 'development';

    return {
        plugins: [react()],
        build: {
            outDir: resolve(__dirname, '../assets/dist'),
            emptyOutDir: true,
            sourcemap: isDev,   // ← Source maps ONLY in dev
            rollupOptions: {
                input: resolve(__dirname, 'src/main.tsx'),
                output: {
                    entryFileNames: 'admin.js',
                    chunkFileNames: 'chunks/[name]-[hash].js',
                    assetFileNames: (assetInfo) => {
                        const isCss = assetInfo.name?.endsWith('.css');

                        return isCss ? 'admin.css' : 'assets/[name]-[hash][extname]';
                    },
                },
            },
        },
    };
});
```

### Build commands

```json
// frontend/package.json (scripts section)
{
  "scripts": {
    "dev": "vite",
    "build": "vite build --mode production",
    "build:dev": "vite build --mode development",
    "preview": "vite preview"
  }
}
```

### .distignore additions for React projects

```

# Add to .distignore when using React

frontend
node_modules
assets/dist/*.map
```

### Verification checklist

```
✅ `npm run build` produces assets/dist/admin.js + admin.css
✅ `npm run build` does NOT produce .map files
✅ `npm run build:dev` DOES produce .map files
✅ assets/dist/ is committed to version control (built artifacts ship with plugin)
✅ frontend/ source is NOT in the distribution ZIP
✅ .map files are NOT in the distribution ZIP
```

---

## 11.9 Enqueuing React Assets

```php
namespace PluginName\Traits\Admin;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\PluginConfigType;

trait AdminReactAssetsTrait
{
    /**
     * Enqueue the React application bundle.
     * Only loads on plugin admin pages.
     */
    public function enqueueReactAssets(string $hook): void
    {
        $isPluginPage = $this->isPluginAdminPage($hook);

        if (!$isPluginPage) {
            return;
        }

        $pluginSlug = PluginConfigType::Slug->value;
        $version = PluginConfigType::Version->value;
        $baseUrl = plugin_dir_url(dirname(__DIR__, 2));
        $distPath = plugin_dir_path(dirname(__DIR__, 2)) . 'assets/dist/';

        // CSS bundle
        $hasCss = file_exists($distPath . 'admin.css');

        if ($hasCss) {
            wp_enqueue_style(
                $pluginSlug . '-react',
                $baseUrl . 'assets/dist/admin.css',
                [],
                $version,
            );
        }

        // JS bundle — depends on wp-element (React provided by WordPress)
        $hasJs = file_exists($distPath . 'admin.js');

        if ($hasJs) {
            wp_enqueue_script(
                $pluginSlug . '-react',
                $baseUrl . 'assets/dist/admin.js',
                ['wp-element'],
                $version,
                true,
            );

            // Inject runtime config for the React app
            wp_localize_script($pluginSlug . '-react', 'PluginNameConfig', [
                'restUrl'  => rest_url(PluginConfigType::apiFullNamespace() . '/'),
                'nonce'    => wp_create_nonce('wp_rest'),
                'version'  => $version,
                'isDebug'  => PluginConfigType::isDebugMode(),
            ]);
        }
    }
}
```

---

## 11.10 Decision Matrix — PHP Templates vs. React

Before implementing the frontend, the developer must decide:

```
┌─────────────────────────────────────────────────────────┐
│  Does the admin UI need complex interactivity?          │
│  (drag-and-drop, real-time updates, shared state)       │
│                                                         │
│  YES → Confirm React with developer → §11.7–§11.9      │
│  NO  → Use PHP templates + partials → §11.2–§11.6      │
│                                                         │
│  HYBRID is also valid:                                  │
│  PHP templates for simple pages (settings, logs)        │
│  React for complex pages (dashboard, visual editors)    │
└─────────────────────────────────────────────────────────┘
```

**The AI must ask the developer before choosing React.** Never default to React without explicit confirmation.

---

## 11.11 Summary Table

| Aspect | Pattern | Reference |
|--------|---------|-----------|
| File size limits | 50–100 ideal, 200 max | §11.1 |
| Page templates | Orchestrator — set variables, include partials | §11.3 |
| Partials | Self-contained, documented variables, max 100 lines | §11.4 |
| Shared partials | `partials/shared/` for cross-page reuse | §11.2, §11.5 |
| Traditional JS/CSS | Per-page files in `assets/css/` and `assets/js/` | §11.6 |
| React (optional) | `frontend/` source → `assets/dist/` build output | §11.7 |
| Source maps | Dev: included; Production: excluded | §11.8 |
| React enqueuing | `AdminReactAssetsTrait` with `wp-element` dependency | §11.9 |
| Decision | Developer must confirm React before implementation | §11.10 |

---

*Phase 11 completes the frontend layer: from PHP template composition through optional React integration, with strict file size discipline throughout.*
