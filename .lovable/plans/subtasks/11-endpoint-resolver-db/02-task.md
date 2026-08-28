# Subtask 2: Core Resolver Implementation

1. Create `gitmap/cmd/resolver.go`.
2. Add `package cmd`.
3. Implement `func resolveEndpointString(raw string) string`.
   - If `raw` starts with `https://`, `http://`, `ssh://`, or `git@`, return it unchanged.
   - Run `filepath.Abs(raw)` and `os.Stat(absPath)`. If error is nil (meaning folder exists), return `absPath`.
   - Call `db, err := openDB()`. If err != nil, return `raw`.
   - Handle ID: `if id, err := strconv.ParseInt(raw, 10, 64); err == nil`. Call `db.FindByID(id)`. If found, return `AbsolutePath`.
   - Handle Alias: `db.ResolveAlias(raw)`. If found without error, return `AbsolutePath`.
   - Handle Slug: `db.FindBySlug(raw)`. If found, return `AbsolutePath`.
   - Fallback: return `raw`.
4. Ensure `gitmap/cmd` package compiles.
