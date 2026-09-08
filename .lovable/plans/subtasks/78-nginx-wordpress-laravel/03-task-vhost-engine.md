# 03-task-vhost-engine.md: Nginx Virtual Host CLI & Management Operations

**Parent Plan:** [.lovable/plans/pending/78-nginx-wordpress-laravel-installation-and-configuration.md](../../pending/78-nginx-wordpress-laravel-installation-and-configuration.md)  
**Status:** Completed  
**Objective:** Provide `gitmap vhost` commands to list, create, enable, disable, test (`nginx -t`), and reload Nginx sites with dynamic PHP-FPM socket discovery.

## Scope & Target Files
- `gitmap/cmd/vhost.go`: Command dispatcher and flag parsing.
- `gitmap/cmd/vhost_ops.go`: VHost operations (create, enable, disable, test, reload).
- `gitmap/cmd/vhost_types.go`: Struct definitions for VHostConfig, SiteType, and VHostOptions.
- `gitmap/cmd/vhost_phpfpm.go`: Dynamic PHP-FPM discovery logic (`/run/php/php*-fpm.sock`, systemd probe, TCP loopback 9000).
