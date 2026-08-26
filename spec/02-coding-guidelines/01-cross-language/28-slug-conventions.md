# Slug Conventions

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Slugs are URL-safe, human-readable identifiers used in URLs, API endpoints, database keys, and file paths. This document defines the **universal slug format** for all languages and frameworks.

---

## Golden Rule

> **All slugs MUST be lowercase kebab-case (hyphen-separated).** No exceptions.

---

## Format

```
lowercase-words-separated-by-hyphens
```

| Rule | Convention |
|------|-----------|
| Case | **Always lowercase** |
| Separator | **Hyphens** (`-`) only — never underscores, spaces, or dots |
| Characters | `a-z`, `0-9`, `-` only |
| Leading/trailing hyphens | ❌ Forbidden |
| Consecutive hyphens | ❌ Forbidden (`my--slug`) |
| Max length | 80 characters recommended |

---

## Examples

### ✅ Correct Slugs

```
user-profile
admin-settings
blog-post-123
my-awesome-plugin
api-v2-endpoints
search-results-page
2026-04-02-url-error-casing-fix
```

### ❌ Incorrect Slugs

```
UserProfile          ← PascalCase
user_profile         ← underscores
user profile         ← spaces
User-Profile         ← mixed case
ADMIN-SETTINGS       ← uppercase
my--slug             ← consecutive hyphens
-leading-hyphen      ← leading hyphen
trailing-hyphen-     ← trailing hyphen
```

---

## REST API Slug Examples

### Resource URLs

```
GET    /api/v1/blog-posts
GET    /api/v1/blog-posts/my-first-post
GET    /api/v1/user-profiles/john-doe
POST   /api/v1/site-settings
PUT    /api/v1/plugin-configs/my-awesome-plugin
DELETE /api/v1/error-logs/2026-04-02-url-error
```

### Full REST API Design Sample

```yaml

# OpenAPI-style paths — all slugs are lowercase kebab-case

paths:
  /api/v1/blog-posts:
    get:
      summary: List all blog posts
      parameters:
        - name: category-slug        # ← query param slug
          in: query
          example: "tech-tutorials"

  /api/v1/blog-posts/{post-slug}:    # ← path param slug
    get:
      summary: Get a single blog post
      parameters:
        - name: post-slug
          in: path
          example: "getting-started-with-go"

  /api/v1/user-profiles/{user-slug}:
    get:
      summary: Get user profile
      parameters:
        - name: user-slug
          in: path
          example: "john-doe"

  /api/v1/error-codes/{error-slug}:
    get:
      summary: Get error details
      parameters:
        - name: error-slug
          in: path
          example: "auth-token-expired"

  /api/v1/site-configs/{config-slug}:
    put:
      summary: Update site configuration
      parameters:
        - name: config-slug
          in: path
          example: "smtp-email-settings"
```

### Database Slug Column

```sql
-- Slugs stored in database columns
CREATE TABLE blog_posts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,    -- lowercase kebab-case
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Example data
INSERT INTO blog_posts (title, slug) VALUES
    ('Getting Started with Go', 'getting-started-with-go'),
    ('URL Error Casing Fix',   '2026-04-02-url-error-casing-fix'),
    ('My First Blog Post',     'my-first-blog-post');

-- Query by slug
SELECT * FROM blog_posts WHERE slug = 'getting-started-with-go';
```

---

## Slug Generation — Code Examples

### Go

```go
func ToSlug(input string) string {
    slug := strings.ToLower(input)
    slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "-")
    slug = regexp.MustCompile(`-{2,}`).ReplaceAllString(slug, "-")
    slug = strings.Trim(slug, "-")

    return slug
}

// ToSlug("My Awesome Plugin!") → "my-awesome-plugin"
// ToSlug("URL Error — Casing") → "url-error-casing"
```

### TypeScript

```typescript
function toSlug(input: string): string {
    return input
        .toLowerCase()
        .replace(/[^a-z0-9-]+/g, "-")
        .replace(/-{2,}/g, "-")
        .replace(/^-|-$/g, "");
}

// toSlug("My Awesome Plugin!") → "my-awesome-plugin"
// toSlug("URL Error — Casing") → "url-error-casing"
```

### PHP

```php
function to_slug(string $input): string {
    $slug = strtolower($input);
    $slug = preg_replace('/[^a-z0-9-]+/', '-', $slug);
    $slug = preg_replace('/-{2,}/', '-', $slug);

    return trim($slug, '-');
}

// to_slug("My Awesome Plugin!") → "my-awesome-plugin"
```

---

## Where Slugs Are Used

| Context | Example |
|---------|---------|
| URL paths | `/blog/my-first-post` |
| REST API resources | `/api/v1/user-profiles/john-doe` |
| Database identifiers | `slug = 'getting-started-with-go'` |
| Plugin/theme slugs | `my-awesome-plugin` |
| Error documentation | `2026-04-02-url-error-casing-fix` |
| Config keys (URL-safe) | `smtp-email-settings` |
| File names (spec docs) | `08-cross-references.md` |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Overview | [./00-overview.md](./00-overview.md) |
| Variable Naming | [./22-variable-naming-conventions.md](./22-variable-naming-conventions.md) |
| Key Naming PascalCase | [./11-key-naming-pascalcase.md](./11-key-naming-pascalcase.md) |
