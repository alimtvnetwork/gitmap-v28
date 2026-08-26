# 02. Coding Guidelines (Redirect)

**Status:** Superseded  
**Canonical:** [`31-compiled-simple-coding-guidelines.md`](./31-compiled-simple-coding-guidelines.md)

This file previously duplicated the coding guidelines. To eliminate drift, the single source of truth is now:

- `spec/17-consolidated-guidelines/31-compiled-simple-coding-guidelines.md` (source)
- `.lovable/coding-guidelines/coding-guidelines.md` (mirror, auto-generated)
- `.cursorrules` (mirror, auto-generated)

Mirrors are produced by `scripts/sync-guidelines.mjs`. Do not hand-edit mirrors, and do not restore the old consolidated content here: any long-form coding rule additions must land in file 31 so the mirrors stay authoritative for agent search.

For language-specific standards previously duplicated in this file (Go, TypeScript, PHP, Rust, C#, PowerShell), read the source specs directly under `spec/02-coding-guidelines/` (subfolders `02-typescript/`, `03-golang/`, `04-php/`, `05-rust/`, `07-csharp/`, `09-powershell-integration/`).
