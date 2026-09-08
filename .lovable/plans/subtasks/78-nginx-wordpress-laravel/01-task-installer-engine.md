# 01-task-installer-engine.md: First-Class Nginx, WordPress, and Laravel Installer Targets

**Parent Plan:** [.lovable/plans/pending/78-nginx-wordpress-laravel-installation-and-configuration.md](../../pending/78-nginx-wordpress-laravel-installation-and-configuration.md)  
**Status:** Completed  
**Objective:** Register Nginx, WordPress, and Laravel as first-class install targets across Windows (`choco` / `winget`), Linux (`apt`), and macOS (`brew`).

## Scope & Target Files
- `gitmap/constants/constants_install.go`: Add `ToolNginx`, `ToolWordPress`, `ToolLaravel`, package IDs, descriptions, and categories.
- `gitmap/cmd/install_packages.go`: Add mappings for `chocoPackageMap`, `wingetPackageMap`, and aliases (`ngx`, `wp`, `wp-cli`, `artisan`).
- `gitmap/cmd/install_packages_extra.go`: Add mappings for `aptPackageMap` and `brewPackageMap`.
- `gitmap/cmd/installverify.go`: Add binary mapping (`nginx`, `wp`, `laravel`), update `expectedExePath`, and enhance `getInstalledVersion` with `CombinedOutput()` and `-v` for Nginx.
- `gitmap/cmd/install_packages_test.go`: Add unit tests for package resolution and alias mapping.

## Invariants & Rules
- All functions <= 15 lines.
- Blank line before return statements.
- Positive booleans only (`is` / `has`).
