# gitmap vhost

Manage Nginx virtual hosts for WordPress, Laravel, generic PHP, and static HTML.

## Alias

None

## Usage

    gitmap vhost [list]
    gitmap vhost create [flags] <domain> <root>
    gitmap vhost enable <domain>
    gitmap vhost disable <domain>
    gitmap vhost test
    gitmap vhost reload

## Subcommands

- `list`, `ls`: List configured virtual host sites.
- `create`, `add`, `new`: Generate production-grade Nginx configuration.
- `enable`, `en`: Enable virtual host via sites-enabled link.
- `disable`, `dis`: Disable virtual host by removing sites-enabled link.
- `test`, `check`, `t`: Validate Nginx configuration syntax (`nginx -t`).
- `reload`, `restart`, `r`: Gracefully reload Nginx (`nginx -s reload`).

## Flags (create)

| Flag | Default | Description |
|------|---------|-------------|
| --type \<type\> | php | Site profile: wordpress, laravel, php, static |
| --domain \<name\> | _(positional)_ | Primary domain name |
| --root \<path\> | _(positional)_ | Document root directory |
| --port \<n\> | 80 | HTTP listen port |
| --ssl | false | Enable SSL port 443 |
| --fastcgi \<pass\> | _(auto)_ | FastCGI socket or host:port upstream |
| --aliases \<list\> | _(empty)_ | Space-delimited server aliases |
| --enable | false | Link to sites-enabled immediately |
| --dry-run | false | Output configuration without saving |

## Examples

### Example 1: Create WordPress virtual host

```bash
gitmap vhost create --type=wordpress mysite.local /var/www/mysite
```

### Example 2: Create Laravel virtual host

```bash
gitmap vhost create --type=laravel api.local /var/www/api/public --enable
```

### Example 3: Test and reload Nginx

```bash
gitmap vhost test
gitmap vhost reload
```
