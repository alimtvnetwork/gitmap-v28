# gitmap nginx

Manage Nginx HTTP server, virtual hosts, SQLite site registry, syntax validation, service reloading, and installation.

## Alias

```bash
gitmap ngx
```

## Usage

```bash
gitmap nginx [subcommand] [flags]
gitmap ngx [subcommand]
```

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| status | st | Display Nginx installation status, version, and binary path |
| test | t, check | Test Nginx configuration file syntax (`nginx -t`) |
| reload | r, restart | Gracefully reload Nginx service (`nginx -s reload`) |
| list | ls | List all virtual hosts tracked in SQLite `sites.db` |
| add | create, new | Add domain or subdomain virtual host and persist to SQLite `sites.db` |
| rm | remove, del | Remove virtual host configuration and delete from SQLite `sites.db` |
| ini | showcase | Showcase PHP/Nginx INI directives, tuning, and marker blocks |
| vhost | vh | Manage virtual hosts (list, create, enable, disable) |
| enable | en | Enable a virtual host via `sites-enabled` symlink |
| disable | dis | Disable a virtual host by removing `sites-enabled` symlink |
| install | in | Install Nginx web server across apt, choco, winget, or brew |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --type | php | Site profile: wordpress, laravel, php, static |
| --port | 80 | HTTP listen port |
| --ssl | false | Enable SSL port 443 |
| --fastcgi | auto | FastCGI socket or host:port upstream |
| --dry-run, -n | false | Simulate actions without applying filesystem or database changes |
| -h, --help | false | Show command help |

## Examples

### Check Nginx Status

```bash
gitmap nginx status
```

### Add an Apex Domain

```bash
gitmap nginx add domain.com
```

### Add a Subdomain for WordPress

```bash
gitmap nginx add sub.domain --type=wordpress
```

### Add a Subdomain for Laravel with Custom Root

```bash
gitmap nginx add api.mysite.com /var/www/api --type=laravel
```

### Remove a Virtual Host and Purge from SQLite DB

```bash
gitmap nginx rm domain.com
```

### List Tracked Sites from SQLite Database

```bash
gitmap nginx list
```

### Showcase PHP/Nginx INI Configuration Updates

```bash
gitmap nginx ini
```

### Validate Configuration Syntax

```bash
gitmap nginx test
```

### Gracefully Reload Nginx

```bash
gitmap nginx reload
```

### Install Nginx

```bash
gitmap install nginx
gitmap nginx install
```
