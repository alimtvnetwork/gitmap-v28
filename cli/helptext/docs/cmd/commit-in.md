# gitmap commit-in & commit-write

Automated, intelligent commit engine featuring author rotation, SEO commit scheduling, function intelligence (funcintel), and profile synchronization.

## Syntax

```bash
gitmap commit-in [flags]
gitmap commit-write [flags]
```

## JSON Configuration Schema

```json
{
  "profile": "default",
  "author": {
    "name": "Jane Doe",
    "email": "jane@example.com"
  },
  "seo": {
    "url": "https://example.com",
    "keywords": ["fast", "cli", "git"]
  },
  "funcintel": {
    "enabled": true,
    "max_lines": 50
  },
  "dedupe": {
    "enabled": true
  }
}
```

## Flags & Options

- `--profile <name>`: Load specific profile configuration.
- `--seo-url <url>`: Website URL for template commit generation.
- `--dry-run`: Preview commits without modifying Git history.
- `--interval <range>`: Random commit delay intervals in seconds.
- `--files <glob>`: Match specific files to stage and commit.
