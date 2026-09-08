# 04-task-setup-wp-laravel.md: WordPress & Laravel Setup and Configuration Engines

**Parent Plan:** [.lovable/plans/pending/78-nginx-wordpress-laravel-installation-and-configuration.md](../../pending/78-nginx-wordpress-laravel-installation-and-configuration.md)  
**Status:** Completed  
**Objective:** Implement `gitmap setup wordpress` and `gitmap setup laravel` commands for automated configuration generation.

## Scope & Target Files
- `gitmap/cmd/setup_wp.go`: WordPress setup runner, database parameter binding, table prefix, debug options.
- `gitmap/cmd/setup_wp_salts.go`: Cryptographic salt generator (8 keys × 64 chars) via `crypto/rand`.
- `gitmap/cmd/setup_laravel.go`: Laravel setup runner, `.env` synthesis, database driver configuration, `php artisan storage:link`.
- `gitmap/cmd/setup_laravel_env.go`: Laravel `APP_KEY` generation (32 bytes Base64) and ordered key-value `.env` merging.
