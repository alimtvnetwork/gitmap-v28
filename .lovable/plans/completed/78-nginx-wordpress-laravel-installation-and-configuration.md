# 78-nginx-wordpress-laravel-installation-and-configuration.md

**Title:** Nginx, WordPress, and Laravel Installation, Setup, and Configuration Engine  
**Status:** Completed  
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration  
**Prompt Version:** 2.1.0  
**Budget (N):** 250 steps  
**Target Codebase:** `gitmap/constants/`, `gitmap/cmd/`, `gitmap/templates/`

---

## 1. Task Overview & Objectives

Deliver production-ready, cross-platform installation, setup, and configuration commands for **Nginx**, **WordPress**, and **Laravel** in GitMap:
1. **First-Class Installer Targets (`gitmap install` / `gitmap in`):**
   - Package manager mapping across `choco`, `winget`, `apt`, and `brew`.
   - Tool canonicalization, aliases (`ngx`, `wp`, `artisan`), and pre-install detection.
   - Version extraction with `stderr` handling (`nginx -v`).
2. **Nginx Virtual Host Engine (`gitmap vhost`):**
   - VHost generator for WordPress, Laravel, generic PHP, and static sites.
   - Idempotent template markers (`# >>> gitmap:vhost/... >>> ... # <<< gitmap:vhost/... <<<`).
   - Dynamic PHP-FPM socket discovery (`/run/php/php*-fpm.sock` with TCP fallback `127.0.0.1:9000`).
   - Configuration validation (`nginx -t`) and reload (`nginx -s reload`).
3. **WordPress Setup Engine (`gitmap setup wordpress` / `gitmap setup wp`):**
   - Automated cryptographic salt generation (8 keys x 64 chars) via pure Go `crypto/rand`.
   - Database credentials binding, debug presets, and table prefix configuration.
   - Automatic Nginx vhost generation option (`--vhost`) and permission fixing (`--fix-perms`).
4. **Laravel Setup Engine (`gitmap setup laravel` / `gitmap setup art`):**
   - Base `.env` generation, ordered key-value merging, and cryptographic `APP_KEY` generation.
   - Public directory document root configuration, query string routing stubs, and `.env` access denial.
   - Automated `storage:link` creation.
5. **Permissions & Security Engine (`gitmap perms`):**
   - Cross-platform permission enforcement (`0755`/`0644` standard, `0775`/`0664` for uploads & storage, `0600` for credentials).
   - Windows NTFS ACL management (`icacls` & `attrib -r`).

---

## 2. Task-Specific Rule Set & Architectural Invariants

1. **Rule 1 (Zero Shell Concatenation):** All system process executions (`nginx -t`, `chmod`, `icacls`, `composer`, `php`) MUST execute discrete argument lists using `exec.Command(name, args...)`. Never pass unescaped shell strings to `cmd /c` or `sh -c`.
2. **Rule 2 (Cryptographic Entropy):** WordPress salts and Laravel `APP_KEY` generation MUST utilize Go standard library `crypto/rand` for cryptographically secure randomness.
3. **Rule 3 (Marker-Block Idempotency):** Any config file injection or template merge MUST use GitMap's sentinel comments (`# >>> gitmap:<tag> >>> ... # <<< gitmap:<tag> <<<`) preserving external user edits byte-for-byte.
4. **Rule 4 (Universal AppError Wrapping):** Every public and internal function returning an error MUST wrap failures in `apperror.New(...)` with specific error codes (`E9000:EXECUTION`, `E1001:VALIDATION`). No bare panics, no swallowed errors.
5. **Rule 5 (Function Sizing & Zero Nesting):** All Go functions must strictly stay within 8–15 lines with blank lines before return statements, and depth <= 1 (zero nested ifs).

---

## 3. Subtask Decomposition Ledger

| Subtask File | Scope & Deliverables | Primary Files |
| :--- | :--- | :--- |
| `01-installer-engine.md` | Installer package IDs, aliases, categories, descriptions, and verification | `gitmap/constants/constants_install.go`<br>`gitmap/cmd/install_packages.go`<br>`gitmap/cmd/install_packages_extra.go`<br>`gitmap/cmd/installverify.go` |
| `02-vhost-templates.md` | WordPress and Laravel Nginx vhost templates with security rules | `gitmap/cmd/vhost_templates.go`<br>`gitmap/cmd/vhost_templates_test.go` |
| `03-vhost-engine.md` | `gitmap vhost` CLI commands (list, create, enable, disable, test, reload) | `gitmap/cmd/vhost.go`<br>`gitmap/cmd/vhost_ops.go`<br>`gitmap/cmd/vhost_types.go` |
| `04-setup-wp-laravel.md` | `gitmap setup wordpress` and `gitmap setup laravel` configuration engines | `gitmap/cmd/setup_wp.go`<br>`gitmap/cmd/setup_wp_salts.go`<br>`gitmap/cmd/setup_laravel.go`<br>`gitmap/cmd/setup_laravel_env.go` |
| `05-perms-engine.md` | Cross-platform permission auditing and hardening engine | `gitmap/cmd/setup_perms.go`<br>`gitmap/cmd/setup_perms_unix.go`<br>`gitmap/cmd/setup_perms_windows.go` |
| `06-verification-and-ci.md` | Unit test suite, local CI runner validation, and binary synchronization | `gitmap/cmd/install_packages_test.go`<br>`gitmap/cmd/vhost_test.go`<br>`gitmap/cmd/setup_test.go` |
