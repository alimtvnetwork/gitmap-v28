# gitmap perms

Audit and fix web application file and directory permissions cross-platform.

## Alias

    gitmap permissions

## Usage

    gitmap perms [flags] [<path>]
    gitmap setup perms [flags] [<path>]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dir \<path\> | . | Target application root directory |
| --app \<type\> | generic | Application type: generic, wordpress, laravel |
| --owner \<user:group\> | www-data:www-data | Unix user:group ownership |
| --fix | false | Apply permission fixes (read-only audit if omitted) |
| --dry-run | false | Simulate fixes without modifying filesystem |

## Permissions Specifications

### Unix / Linux / macOS
- **Directories**: 0755 (standard), 0775 (writable paths such as `storage/`, `bootstrap/cache/`, `wp-content/uploads/`)
- **Files**: 0644 (standard), 0664 (writable files)
- **Sensitive Credentials**: 0600 (`wp-config.php`, `.env`, `*.key`, `*.pem`)

### Windows
- Grants full read/write inheritance via `icacls` to current user
- Clears read-only attributes on writable paths via `attrib -r`

## Examples

### Example 1: Audit WordPress permissions (read-only)

```bash
gitmap perms --app=wordpress /var/www/html
```

### Example 2: Fix Laravel permissions

```bash
gitmap perms --app=laravel --fix /var/www/my-laravel-app
```
