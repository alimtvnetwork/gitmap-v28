# Nginx, WordPress, and Laravel Setup, VHost & Perms Engine Architecture

## 1. Context & Objectives
Delivered production-grade installation, configuration, and management for **Nginx**, **WordPress**, and **Laravel** across Windows, Linux, and macOS:
- `gitmap install nginx|wordpress|laravel`: Universal package installation across `choco`, `winget`, `apt`, and `brew`.
- `gitmap vhost`: Virtual host configuration, site enabling/disabling, syntax testing (`nginx -t`), and reloading (`nginx -s reload`).
- `gitmap setup wordpress`: Automated `wp-config.php` generation with Go `crypto/rand` 64-character salts and database parameter binding.
- `gitmap setup laravel`: Automated `.env` key-value synthesis, 32-byte base64 `APP_KEY` generation, `public/` web root, and `storage:link`.
- `gitmap perms`: Cross-platform permission auditing and fixing (`0755`/`0644`/`0775`/`0600` on Unix; `icacls` & `attrib -r` on Windows).

## 2. Key Technical Findings & Learned Patterns

### 1. Positional Argument Flag Halting in Go
- **Problem**: Standard library `flag.FlagSet.Parse(args)` halts flag processing upon encountering the first non-flag positional argument (e.g. `gitmap setup wordpress /var/www/site --db-name=wp`). Any flags specified after the positional argument are ignored and remain in `fs.Args()`.
- **Solution**: Always sanitize input flags via `reorderFlagsBeforeArgs(args)` prior to `fs.Parse(reorderFlagsBeforeArgs(args))`. This bubbles all flags (`--*`, `-*`) ahead of positional arguments.

### 2. Typed Nil Interface Pointer Leak in Error Returns
- **Problem**: In Go, returning a concrete typed nil pointer (e.g., `(*apperror.AppError)(nil)`) from a function whose signature returns the built-in `error` interface causes `err != nil` to evaluate to `true` because the interface contains type metadata.
- **Solution**: Explicitly inspect concrete errors before returning them to `error` interface returns:
  ```go
  appErr := auditItem(...)
  if appErr != nil {
      return appErr
  }
  return nil
  ```

### 3. Nginx Version Output Exclusively to Stderr
- **Problem**: `nginx --version` is rejected by Nginx, and `nginx -v` writes its output exclusively to `stderr` rather than `stdout`. Using `cmd.Output()` resulted in empty version strings.
- **Solution**: In `installverify.go`, use flag `-v` for Nginx and invoke `cmd.CombinedOutput()` to capture both streams.

### 4. Positive Boolean Naming Guidelines
- **Problem**: Variables such as `hasNoMatches` violate the coding guidelines requiring strictly affirmative variable prefixes (`is*`, `has*`).
- **Solution**: Always use positive boolean logic, e.g. `hasMatches := len(matches) > 0; if !hasMatches { return "" }`.

### 5. Single-Level If Statements and Function Sizing
- **Problem**: Nested conditions (`if outer { if inner { ... } }`) violate repository linters.
- **Solution**: Extract inner blocks into dedicated helper functions or use early returns so that every conditional has maximum nesting depth 1.
