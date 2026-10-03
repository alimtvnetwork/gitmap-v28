# Phase 10 — Deployment Patterns

> **Purpose:** Define how to package, version, distribute, and auto-update WordPress plugins built with this architecture. Covers ZIP packaging, semantic versioning, self-hosted update servers, self-update with rollback, and CI/CD automation.
> **Audience:** AI code generators and human developers.
> **Prerequisite:** Phases 1–9 must be read first.

---

## 10.1 Versioning Strategy

### Semantic Versioning

All plugins follow [SemVer 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH
  │      │     └── Bug fixes, typo corrections, no API changes
  │      └──────── New features, backward-compatible
  └─────────────── Breaking changes to REST API, database schema, or hook contracts
```

### Version Source of Truth

The version is declared in **exactly two places** and must always match:

| Location | Example | Read by |
|----------|---------|---------|
| Main plugin file header | `* Version: 2.31.0` | WordPress core |
| `PluginConfigType::Version` enum case | `case Version = '2.31.0'` | All plugin code |

**Rule:** Every code change requires a minor version bump (per project versioning policy). The `.release` folder is exempt.

### Version Bump Checklist

```
1. Update main plugin file header → Version: X.Y.Z
2. Update PluginConfigType::Version → case Version = 'X.Y.Z'
3. Update changelog.md → add entry under ## [X.Y.Z] - YYYY-MM-DD
4. If composer.json has a version field → update it too
5. Commit with message: "Bump version to X.Y.Z"
```

---

## 10.2 Plugin File Structure for Distribution

### What ships vs. what stays behind

```
plugin-slug/                      ← ZIP root
├── plugin-slug.php               ← Main plugin file (required)
├── uninstall.php                 ← Clean removal hook
├── readme.md                     ← Plugin readme
├── changelog.md                  ← Version history
├── composer.json                 ← Dependency manifest
├── settings.json                 ← Default configuration
├── assets/                       ← Admin CSS/JS/images
├── data/
│   └── seeds/                    ← Seed JSON files + manifest
├── includes/                     ← All PHP source (PSR-4)
├── templates/                    ← PHP view templates
└── vendor/                       ← Composer autoloader (production only)
```

### Excluded from distribution ZIP

```
.git/
.github/
.ai-instructions
node_modules/
tests/
phpunit.xml
phpstan.neon
phpstan-bootstrap.php
composer.lock
spec/
*.log
.DS_Store
Thumbs.db
```

### .distignore file

Create a `.distignore` at the plugin root listing files excluded from the ZIP:

```
.git
.github
.ai-instructions
tests
phpunit.xml
phpstan.neon
phpstan-bootstrap.php
composer.lock
spec
*.log
.DS_Store
Thumbs.db
```

---

## 10.3 ZIP Packaging

### Manual packaging (development)

```bash

# From the directory containing the plugin folder

cd wp-plugins/

# Install production dependencies only

cd plugin-slug && composer install --no-dev --optimize-autoloader && cd ..

# Create ZIP excluding dev files

zip -r plugin-slug-v2.31.0.zip plugin-slug/ \
  -x "plugin-slug/.git/*" \
  -x "plugin-slug/.github/*" \
  -x "plugin-slug/.ai-instructions" \
  -x "plugin-slug/tests/*" \
  -x "plugin-slug/phpunit.xml" \
  -x "plugin-slug/phpstan.neon" \
  -x "plugin-slug/phpstan-bootstrap.php" \
  -x "plugin-slug/composer.lock" \
  -x "plugin-slug/spec/*" \
  -x "plugin-slug/*.log"
```

### Automated packaging script

Create `scripts/package.sh` in the plugin root:

```bash
#!/bin/bash
set -euo pipefail

PLUGIN_SLUG="plugin-slug"
VERSION=$(grep -oP "Version:\s*\K[0-9.]+" "${PLUGIN_SLUG}.php")
OUTPUT="${PLUGIN_SLUG}-v${VERSION}.zip"

echo "📦 Packaging ${PLUGIN_SLUG} v${VERSION}..."

# Production dependencies

composer install --no-dev --optimize-autoloader --quiet

# Build ZIP respecting .distignore

if command -v rsync &> /dev/null; then
    TMPDIR=$(mktemp -d)
    rsync -a --exclude-from=.distignore . "${TMPDIR}/${PLUGIN_SLUG}/"
    cd "${TMPDIR}"
    zip -r "${OLDPWD}/${OUTPUT}" "${PLUGIN_SLUG}/"
    rm -rf "${TMPDIR}"
else
    # Fallback: manual exclusions
    cd ..
    zip -r "${PLUGIN_SLUG}/${OUTPUT}" "${PLUGIN_SLUG}/" \
        -x "${PLUGIN_SLUG}/.git/*" \
        -x "${PLUGIN_SLUG}/tests/*" \
        -x "${PLUGIN_SLUG}/spec/*" \
        -x "${PLUGIN_SLUG}/phpunit.xml" \
        -x "${PLUGIN_SLUG}/phpstan.neon" \
        -x "${PLUGIN_SLUG}/composer.lock"
fi

echo "✅ Created ${OUTPUT} ($(du -h "${OUTPUT}" | cut -f1))"
```

### ZIP integrity requirements

| Check | Rule |
|-------|------|
| Main plugin file exists | `plugin-slug/plugin-slug.php` must be at ZIP root level |
| No nested folders | ZIP must not contain `plugin-slug/plugin-slug/` (double-nesting) |
| vendor/ present | Autoloader must be included — plugin won't boot without it |
| No dev dependencies | `vendor/phpunit/` must NOT appear in the ZIP |
| File permissions | PHP files: 644, directories: 755 |

---

## 10.4 Self-Hosted Update Server

### Why self-hosted?

Plugins not listed on wordpress.org need a custom update mechanism. The pattern hooks into WordPress's native update system via two filters.

### Update JSON endpoint

The update server (or a static JSON file) must serve an update-info response:

```json
{
  "name": "Plugin Name",
  "slug": "plugin-slug",
  "version": "2.31.0",
  "download_url": "https://updates.example.com/plugin-slug/plugin-slug-v2.31.0.zip",
  "tested": "6.7",
  "requires": "5.6",
  "requires_php": "8.2",
  "last_updated": "2026-04-08",
  "sections": {
    "description": "Plugin description here.",
    "changelog": "<h4>2.31.0</h4><ul><li>New feature X</li></ul>"
  }
}
```

### UpdateResolver — Hook Registration

The plugin hooks into WordPress's update system using two filters:

```php
namespace PluginName\Update;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\HookType;
use PluginName\Enums\OptionNameType;
use PluginName\Enums\UpdateConfigType;
use PluginName\Logging\FileLogger;

class UpdateResolver
{
    use Traits\UpdateResolverUrlTrait;
    use Traits\UpdateResolverFetchTrait;
    use Traits\UpdateResolverWpHooksTrait;
    use Traits\UpdateResolverIntegrityTrait;
    use Traits\UpdateResolverBackupTrait;

    private FileLogger $fileLogger;
    private static ?self $instance = null;

    public static function getInstance(): static
    {
        if (self::$instance === null) {
            self::$instance = new self();
        }

        return self::$instance;
    }

    private function __construct()
    {
        $this->fileLogger = FileLogger::getInstance();

        $settings = $this->getSettings();
        $isEnabled = !empty($settings['enabled']);

        if ($isEnabled) {
            add_filter(
                HookType::PreSetSiteTransientUpdatePlugins->value,
                [$this, 'checkForPluginUpdate'],
            );
            add_filter(
                HookType::PluginsApi->value,
                [$this, 'pluginInfo'],
                10,
                3,
            );
        }
    }
}
```

### Required HookType enum cases

```php
// In HookType enum — add these cases for update hooks
case PreSetSiteTransientUpdatePlugins = 'pre_set_site_transient_update_plugins';
case PluginsApi                       = 'plugins_api';
```

### UpdateResolverWpHooksTrait — Core Hooks

```php
namespace PluginName\Update\Traits;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\PluginConfigType;
use stdClass;

trait UpdateResolverWpHooksTrait
{
    /**
     * Hook into WordPress update check.
     * Called via pre_set_site_transient_update_plugins filter.
     */
    public function checkForPluginUpdate(mixed $transient): mixed
    {
        $hasChecked = is_object($transient) && property_exists($transient, 'checked');

        if (!$hasChecked) {
            return $transient;
        }

        $settings = $this->getSettings();
        $remoteVersion = $settings['new_version'] ?? '';
        $currentVersion = PluginConfigType::Version->value;
        $pluginBasename = PluginConfigType::Basename->value;
        $hasUpdate = (
            !empty($remoteVersion)
            && version_compare($remoteVersion, $currentVersion, '>')
        );

        if (!$hasUpdate) {
            return $transient;
        }

        $updateObj = new stdClass();
        $updateObj->slug        = PluginConfigType::Slug->value;
        $updateObj->plugin      = $pluginBasename;
        $updateObj->new_version = $remoteVersion;
        $updateObj->url         = PluginConfigType::PluginUri->value;
        $updateObj->package     = $settings['package_url'] ?? '';
        $updateObj->tested      = $settings['update_info']['tested'] ?? '';
        $updateObj->requires    = $settings['update_info']['requires'] ?? '';

        $transient->response[$pluginBasename] = $updateObj;

        return $transient;
    }

    /**
     * Provide plugin information for the WordPress "View Details" modal.
     * Called via plugins_api filter.
     */
    public function pluginInfo(mixed $result, string $action, object $args): mixed
    {
        $isQueryingThisPlugin = (
            $action === 'plugin_information'
            && isset($args->slug)
            && $args->slug === PluginConfigType::Slug->value
        );

        if (!$isQueryingThisPlugin) {
            return $result;
        }

        $settings = $this->getSettings();
        $info = $settings['update_info'] ?? [];
        $hasInfo = !empty($info);

        if (!$hasInfo) {
            return $result;
        }

        $pluginInfo = new stdClass();
        $pluginInfo->name          = $info['name'] ?? PluginConfigType::Name->value;
        $pluginInfo->slug          = PluginConfigType::Slug->value;
        $pluginInfo->version       = $info['version'] ?? '';
        $pluginInfo->author        = $info['author'] ?? '';
        $pluginInfo->download_link = $settings['package_url'] ?? '';
        $pluginInfo->tested        = $info['tested'] ?? '';
        $pluginInfo->requires      = $info['requires'] ?? '';
        $pluginInfo->requires_php  = $info['requires_php'] ?? '';
        $pluginInfo->last_updated  = $info['last_updated'] ?? '';
        $pluginInfo->sections      = $info['sections'] ?? [];

        return $pluginInfo;
    }
}
```

### Update settings storage

```php
// OptionNameType enum — add case for update settings
case UpdateSettings = 'plugin_slug_update_settings';
```

Settings structure stored in `wp_options`:

```php
$defaults = [
    'enabled'      => false,          // Master toggle
    'master_url'   => '',             // User-configured update URL
    'resolved_url' => '',             // After 301 redirect resolution
    'resolved_at'  => '',             // When URL was last resolved
    'cache_days'   => 7,              // Days to cache resolved URL
    'last_check'   => '',             // Last update check timestamp
    'last_error'   => '',             // Last error message
    'package_url'  => '',             // Direct ZIP download URL
    'new_version'  => '',             // Latest available version
    'update_info'  => [],             // Full update metadata
];
```

---

## 10.5 Self-Update with Rollback

### Update lifecycle

```
1. Admin triggers update (via WP dashboard or AJAX endpoint)
2. Plugin creates backup of current version → wp-content/uploads/plugin-slug/backups/
3. Plugin downloads ZIP from package_url
4. Plugin extracts ZIP to plugin directory
5. Pre-activation validation runs (critical files, syntax check)
6. WordPress activates the new version
7. Post-activation health check runs (classes loaded, hooks registered)
8. On ANY failure at steps 5-7 → automatic rollback from backup
9. On success → backup retained for manual rollback window
```

### SelfUpdateStatusType Enum

> **Full definition:** [Phase 2 — 03-self-update-status-enum.md](./02-enums-and-coding-style/04-self-update-status-enum.md)

This enum tracks every possible outcome of the self-update lifecycle. It uses `match`-based metadata methods (see [02-enum-metadata-pattern.md](./02-enums-and-coding-style/03-enum-metadata-pattern.md)) with per-case `is*()` helpers.

Key cases: `Success`, `RolledBack`, `RollbackFailed`, `BackupCreationFailed`, `ExtractionFailed`, `ValidationFailed`, `ActivationException`, `ActivationWpError`, `HealthCheckFailed`, `PluginFileNotFound`, `CriticalFileMissing`, `SyntaxError`, `FileUnreadable`, `DirectoryMissing`, `BootErrorDetected`, `CriticalClassMissing`, `RestHookMissing`.

Domain helpers: `isRollbackReason()`, `isSuccess()`, `label()`, `info()`.

### Pre-activation validation (SelfUpdateValidator)

```php
namespace PluginName\Update;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\SelfUpdateStatusType;

final class SelfUpdateValidator
{
    /**
     * Validate the extracted plugin before activation.
     *
     * @param string $pluginDir Path to the extracted plugin directory
     * @return array{valid: bool, errors: list<array{status: SelfUpdateStatusType, file: string}>}
     */
    public function validate(string $pluginDir): array
    {
        $errors = [];
        $criticalFiles = $this->getCriticalFiles();

        foreach ($criticalFiles as $relativePath) {
            $fullPath = $pluginDir . '/' . $relativePath;
            $fileExists = file_exists($fullPath);

            if (!$fileExists) {
                $errors[] = [
                    'status' => SelfUpdateStatusType::CriticalFileMissing,
                    'file'   => $relativePath,
                ];

                continue;
            }

            $isReadable = is_readable($fullPath);

            if (!$isReadable) {
                $errors[] = [
                    'status' => SelfUpdateStatusType::FileUnreadable,
                    'file'   => $relativePath,
                ];

                continue;
            }

            // Syntax check PHP files
            $isPhpFile = (pathinfo($fullPath, PATHINFO_EXTENSION) === 'php');

            if ($isPhpFile) {
                $output = [];
                $exitCode = 0;
                exec("php -l " . escapeshellarg($fullPath) . " 2>&1", $output, $exitCode);
                $hasSyntaxError = ($exitCode !== 0);

                if ($hasSyntaxError) {
                    $errors[] = [
                        'status' => SelfUpdateStatusType::SyntaxError,
                        'file'   => $relativePath,
                    ];
                }
            }
        }

        $hasErrors = (count($errors) > 0);

        return [
            'valid'  => !$hasErrors,
            'errors' => $errors,
        ];
    }

    /**
     * Files that MUST exist for the plugin to function.
     *
     * @return list<string>
     */
    private function getCriticalFiles(): array
    {
        return [
            'plugin-slug.php',
            'includes/Core/Plugin.php',
            'includes/Enums/PluginConfigType.php',
            'vendor/autoload.php',
        ];
    }
}
```

### Post-activation health check (SelfUpdateHealthCheck)

```php
namespace PluginName\Update;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\SelfUpdateStatusType;
use PluginName\ErrorHandling\BootErrorCollector;

final class SelfUpdateHealthCheck
{
    /**
     * Run health checks after the new version is activated.
     *
     * @return array{healthy: bool, errors: list<array{status: SelfUpdateStatusType, detail: string}>}
     */
    public function check(): array
    {
        $errors = [];

        // Check 1: Boot errors
        $bootErrors = BootErrorCollector::getErrors();
        $hasBootErrors = (count($bootErrors) > 0);

        if ($hasBootErrors) {
            $errors[] = [
                'status' => SelfUpdateStatusType::BootErrorDetected,
                'detail' => implode('; ', $bootErrors),
            ];
        }

        // Check 2: Critical classes loaded
        $criticalClasses = $this->getCriticalClasses();

        foreach ($criticalClasses as $className) {
            $isLoaded = class_exists($className, false);

            if (!$isLoaded) {
                $errors[] = [
                    'status' => SelfUpdateStatusType::CriticalClassMissing,
                    'detail' => $className,
                ];
            }
        }

        // Check 3: REST hooks registered
        $hasRestHooks = has_action('rest_api_init');

        if (!$hasRestHooks) {
            $errors[] = [
                'status' => SelfUpdateStatusType::RestHookMissing,
                'detail' => 'rest_api_init hook not registered',
            ];
        }

        $hasErrors = (count($errors) > 0);

        return [
            'healthy' => !$hasErrors,
            'errors'  => $errors,
        ];
    }

    /**
     * Classes that must be loaded after activation.
     *
     * @return list<string>
     */
    private function getCriticalClasses(): array
    {
        return [
            'PluginName\\Core\\Plugin',
            'PluginName\\Enums\\PluginConfigType',
            'PluginName\\Logging\\FileLogger',
        ];
    }
}
```

### Backup and rollback (UpdateResolverBackupTrait)

```php
namespace PluginName\Update\Traits;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\SelfUpdateStatusType;
use PluginName\Helpers\PathHelper;

trait UpdateResolverBackupTrait
{
    /**
     * Create a backup of the current plugin version.
     *
     * @param string $pluginDir  Current plugin directory
     * @param string $backupDir  Backup destination
     * @return array{success: bool, path: string, status?: SelfUpdateStatusType}
     */
    public function createBackup(string $pluginDir, string $backupDir): array
    {
        $timestamp = gmdate('Ymd_His');
        $backupPath = $backupDir . '/backup_' . $timestamp;

        if (PathHelper::isDirMissing($backupDir)) {
            PathHelper::makeDirectory($backupDir);
        }

        $copied = $this->recursiveCopy($pluginDir, $backupPath);

        if (!$copied) {
            return [
                'success' => false,
                'path'    => '',
                'status'  => SelfUpdateStatusType::BackupCreationFailed,
            ];
        }

        return [
            'success' => true,
            'path'    => $backupPath,
        ];
    }

    /**
     * Restore from backup after a failed update.
     *
     * @param string $backupPath Path to the backup
     * @param string $pluginDir  Plugin directory to restore to
     * @return bool Whether the rollback succeeded
     */
    public function rollback(string $backupPath, string $pluginDir): bool
    {
        $hasBackup = is_dir($backupPath);

        if (!$hasBackup) {
            return false;
        }

        // Remove failed new version
        $this->recursiveDelete($pluginDir);

        // Restore from backup
        return $this->recursiveCopy($backupPath, $pluginDir);
    }

    private function recursiveCopy(string $src, string $dst): bool
    {
        // Implementation: recursive directory copy
        // Uses opendir/readdir pattern with error handling
        // Returns false on any failure
    }

    private function recursiveDelete(string $dir): bool
    {
        // Implementation: recursive directory deletion
        // Uses RecursiveDirectoryIterator
        // Returns false on any failure
    }
}
```

---

## 10.6 URL Resolution with Redirect Handling

### Why resolve URLs?

Update URLs may use URL shorteners or redirectors that return 301/302 responses. The plugin must resolve through redirect chains to find the final download URL.

### UpdateResolverUrlTrait

```php
namespace PluginName\Update\Traits;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\UpdateConfigType;

trait UpdateResolverUrlTrait
{
    /**
     * Resolve a URL through redirect chains to find the final destination.
     *
     * @param string $url          Starting URL
     * @param int    $maxRedirects Maximum redirects to follow
     * @return array{resolved: bool, url: string, redirects: int}
     */
    public function resolveUrl(string $url, int $maxRedirects = 5): array
    {
        $currentUrl = $url;
        $redirectCount = 0;
        $limit = min($maxRedirects, UpdateConfigType::MaxRedirects->value);

        for ($i = 0; $i < $limit; $i++) {
            $response = wp_remote_head($currentUrl, [
                'redirection' => 0,  // Don't auto-follow
                'timeout'     => 10,
            ]);

            $isError = is_wp_error($response);

            if ($isError) {
                break;
            }

            $statusCode = wp_remote_retrieve_response_code($response);
            $isRedirect = in_array($statusCode, [301, 302, 307, 308], true);

            if (!$isRedirect) {
                return [
                    'resolved'  => true,
                    'url'       => $currentUrl,
                    'redirects' => $redirectCount,
                ];
            }

            $location = wp_remote_retrieve_header($response, 'location');
            $hasLocation = !empty($location);

            if (!$hasLocation) {
                break;
            }

            $currentUrl = $location;
            $redirectCount++;
        }

        return [
            'resolved'  => false,
            'url'       => $currentUrl,
            'redirects' => $redirectCount,
        ];
    }

    /**
     * Check if the cached resolved URL is still fresh.
     */
    public function isUrlCacheFresh(string $resolvedAt, int $cacheDays): bool
    {
        $hasResolvedAt = !empty($resolvedAt);

        if (!$hasResolvedAt) {
            return false;
        }

        $resolvedTime = strtotime($resolvedAt);
        $expiresAt = $resolvedTime + ($cacheDays * DAY_IN_SECONDS);
        $now = time();

        return ($now < $expiresAt);
    }
}
```

---

## 10.7 UpdateConfigType Enum

> **Pattern:** Config enums (int-backed, no metadata) — do NOT use info-object pattern. See [01-enum-architecture.md §2.1](./02-enums-and-coding-style/02-enum-architecture.md).

| Case | Value | Purpose |
|------|-------|---------|
| `CacheDaysDefault` | `7` | Days to cache resolved update URL |
| `MaxRedirects` | `5` | Maximum 301/302 redirects to follow |

---

## 10.8 changelog.md Format

Every version bump must include a CHANGELOG entry:

```markdown

# Changelog

All notable changes to this project will be documented in this file.

## [2.31.0] - 2026-04-08

### Added

- New feature description

### Changed

- Updated behavior description

### Fixed

- Bug fix description

### Removed

- Deprecated feature removed

## [2.30.0] - 2026-04-01

...
```

### Rules

| Rule | Detail |
|------|--------|
| Newest first | Latest version at the top |
| Date format | ISO 8601: `YYYY-MM-DD` |
| Categories | `Added`, `Changed`, `Fixed`, `Removed`, `Security`, `Deprecated` |
| One line per change | Each bullet is a single, clear sentence |
| Reference issues | Link to GitHub issues or spec files where applicable |

---

## 10.9 Uninstall Cleanup

### uninstall.php

Every plugin must include `uninstall.php` for clean removal:

```php
<?php
/**
 * Uninstall handler — runs when the plugin is deleted via WP admin.
 *
 * @package PluginName
 */

if (!defined('WP_UNINSTALL_PLUGIN')) {
    exit;
}

// Remove plugin options
$optionKeys = [
    'plugin_slug_update_settings',
    'plugin_slug_settings',
    'plugin_slug_version',
];

foreach ($optionKeys as $key) {
    delete_option($key);
}

// Remove transients
delete_transient('plugin_slug_cache');

// Remove uploaded files (SQLite databases, logs, backups)
$uploadDir = wp_upload_dir();
$pluginDataDir = $uploadDir['basedir'] . '/plugin-slug';
$hasDataDir = is_dir($pluginDataDir);

if ($hasDataDir) {
    // Recursive delete of plugin data directory
    $iterator = new RecursiveIteratorIterator(
        new RecursiveDirectoryIterator($pluginDataDir, RecursiveDirectoryIterator::SKIP_DOTS),
        RecursiveIteratorIterator::CHILD_FIRST,
    );

    foreach ($iterator as $file) {
        $isDir = $file->isDir();

        if ($isDir) {
            rmdir($file->getPathname());
        } else {
            unlink($file->getPathname());
        }
    }

    rmdir($pluginDataDir);
}
```

### What to clean

| Data type | Location | Clean on uninstall? |
|-----------|----------|-------------------|
| Plugin options | `wp_options` table | ✅ Always |
| Transients | `wp_options` table | ✅ Always |
| SQLite databases | `wp-content/uploads/plugin-slug/` | ✅ Always |
| Log files | `wp-content/uploads/plugin-slug/logs/` | ✅ Always |
| Backups | `wp-content/uploads/plugin-slug/backups/` | ✅ Always |
| Cron events | WordPress cron | ✅ Always |
| User meta | `wp_usermeta` table | ⚠️ Only if plugin added it |
| Custom tables (MySQL) | WordPress database | ⚠️ Only if plugin created them |

---

## 10.10 Trait Decomposition for UpdateResolver

Following the Gold Standard trait composition pattern (Phase 3), the UpdateResolver is decomposed into focused traits:

```
Update/
├── UpdateResolver.php                    ← Shell class (singleton)
├── SelfUpdateValidator.php               ← Pre-activation checks
├── SelfUpdateHealthCheck.php             ← Post-activation checks
├── SelfUpdateBackupHelper.php            ← Backup creation utilities
└── Traits/
    ├── UpdateResolverUrlTrait.php         ← URL resolution + redirect handling
    ├── UpdateResolverFetchTrait.php       ← HTTP fetching of update info
    ├── UpdateResolverWpHooksTrait.php     ← WordPress filter callbacks
    ├── UpdateResolverIntegrityTrait.php   ← ZIP integrity verification
    └── UpdateResolverBackupTrait.php      ← Backup creation + rollback
```

| Trait | Responsibility |
|-------|---------------|
| `UpdateResolverUrlTrait` | Resolve URLs through 301/302 redirects, cache resolved URLs |
| `UpdateResolverFetchTrait` | Fetch update JSON from remote server, parse response |
| `UpdateResolverWpHooksTrait` | `checkForPluginUpdate()` and `pluginInfo()` filter callbacks |
| `UpdateResolverIntegrityTrait` | Verify ZIP checksums, validate extracted contents |
| `UpdateResolverBackupTrait` | Create pre-update backup, restore on failure |

---

## 10.11 CI/CD Automation (GitHub Actions)

### Recommended workflow

```yaml

# .github/workflows/release.yml

name: Release Plugin

on:
  push:
    tags:
      - 'v*'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: shivammathur/setup-php@v2
        with:
          php-version: '8.2'
      - run: composer install
      - run: vendor/bin/phpunit
      - run: vendor/bin/phpstan analyse

  package:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: shivammathur/setup-php@v2
        with:
          php-version: '8.2'
      - run: composer install --no-dev --optimize-autoloader
      - name: Package ZIP
        run: |
          PLUGIN_SLUG="${{ github.event.repository.name }}"
          VERSION="${GITHUB_REF#refs/tags/v}"
          mkdir -p dist
          rsync -a --exclude-from=.distignore . "dist/${PLUGIN_SLUG}/"
          cd dist && zip -r "../${PLUGIN_SLUG}-v${VERSION}.zip" "${PLUGIN_SLUG}/"
      - name: Upload release asset
        uses: softprops/action-gh-release@v1
        with:
          files: '*.zip'
```

### Release checklist

```
1. ✅ All tests pass (phpunit + phpstan)
2. ✅ Version bumped in plugin header + PluginConfigType
3. ✅ changelog.md updated
4. ✅ composer install --no-dev succeeds
5. ✅ ZIP packages correctly (no double-nesting, vendor/ present)
6. ✅ ZIP installs cleanly on a fresh WordPress site
7. ✅ Auto-update detects new version from update server
8. ✅ Self-update completes with health check passing
9. ✅ Rollback works when validation fails
```

---

## 10.12 Summary Table

| Aspect | Pattern | Reference |
|--------|---------|-----------|
| Version source | Plugin header + `PluginConfigType::Version` | §10.1 |
| ZIP packaging | `.distignore` + `scripts/package.sh` | §10.3 |
| Update detection | `pre_set_site_transient_update_plugins` filter | §10.4 |
| Plugin info modal | `plugins_api` filter | §10.4 |
| Self-update | Download → Validate → Activate → Health Check → Rollback | §10.5 |
| URL resolution | Follow 301/302 chains with configurable max | §10.6 |
| Rollback | Backup before update, restore on failure | §10.5 |
| Uninstall cleanup | `uninstall.php` removes options, data, logs | §10.9 |
| CI/CD | GitHub Actions: test → package → release | §10.11 |

---

*Phase 10 completes the plugin lifecycle: from development (Phases 1–9) through packaging, distribution, updating, and clean removal.*
