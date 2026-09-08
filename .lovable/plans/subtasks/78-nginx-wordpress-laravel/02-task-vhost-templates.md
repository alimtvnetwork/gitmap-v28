# 02-task-vhost-templates.md: Nginx VHost Templates & Configuration Generator

**Parent Plan:** [.lovable/plans/pending/78-nginx-wordpress-laravel-installation-and-configuration.md](../../pending/78-nginx-wordpress-laravel-installation-and-configuration.md)  
**Status:** Completed  
**Objective:** Provide production-grade, secure Nginx virtual host templates for WordPress, Laravel, PHP, and static websites with marker-block idempotency.

## Scope & Target Files
- `gitmap/cmd/vhost_templates.go`: Define templates for WordPress (`try_files $uri $uri/ /index.php?$args;`, upload directory PHP block, xmlrpc block, fastcgi tuning), Laravel (`root /public`, `try_files $uri $uri/ /index.php?$query_string;`, `.env` deny block), generic PHP, and static.
- `gitmap/cmd/vhost_templates_test.go`: Unit tests ensuring valid template rendering, marker generation, and parameter injection.
