# gitmap nginx

Manage Nginx HTTP server, virtual hosts, syntax validation, service reloading, and installation.

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
| vhost | vh | Manage virtual hosts (list, create, enable, disable) |
| list | ls | List all configured virtual host sites |
| create | add, new | Generate a new Nginx virtual host configuration |
| enable | en | Enable a virtual host via `sites-enabled` symlink |
| disable | dis | Disable a virtual host by removing `sites-enabled` symlink |
| install | in | Install Nginx web server across apt, choco, winget, or brew |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run, -n | false | Simulate actions without applying filesystem or service changes |
| -h, --help | false | Show command help |

## Examples

### Check Nginx Status

```bash
gitmap nginx status
```

### Validate Configuration Syntax

```bash
gitmap nginx test
```

Output:

```text
▶ Validating Nginx configuration...
  ✓ Configuration test successful (nginx -t)
```

### Gracefully Reload Nginx

```bash
gitmap nginx reload
```

Output:

```text
▶ Reloading Nginx...
  ✓ Configuration test passed (nginx -t)
  ✓ Nginx reloaded successfully (nginx -s reload)
```

### Create a Virtual Host for WordPress

```bash
gitmap nginx vhost create --type=wordpress mysite.local /var/www/mysite --enable
```

### Install Nginx

```bash
gitmap nginx install
```
