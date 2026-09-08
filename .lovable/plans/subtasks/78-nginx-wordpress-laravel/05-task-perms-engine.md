# 05-task-perms-engine.md: Web Application Permissions & Hardening Engine

**Parent Plan:** [.lovable/plans/pending/78-nginx-wordpress-laravel-installation-and-configuration.md](../../pending/78-nginx-wordpress-laravel-installation-and-configuration.md)  
**Status:** Completed  
**Objective:** Provide `gitmap perms` / `gitmap permissions` to audit and fix directory and file permissions across Linux (`www-data:www-data`, `0755`/`0644`, `0775`/`0664`, `0600`) and Windows (`icacls`, `attrib -r`).

## Scope & Target Files
- `gitmap/cmd/setup_perms.go`: CLI dispatcher and platform router.
- `gitmap/cmd/setup_perms_unix.go`: Linux/Unix permission audit and fixer.
- `gitmap/cmd/setup_perms_windows.go`: Windows NTFS ACL audit and fixer.
