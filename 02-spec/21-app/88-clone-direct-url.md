# Clone Direct URL

## Responsibility

Allow `gitmap-v28 clone` to accept a direct Git URL (HTTPS or SSH) instead of
requiring a pre-generated file (JSON/CSV/text). The repo is cloned into the
current directory (or a custom folder name), tracked in the database, and
the user is prompted to register it with GitHub Desktop.

## Usage

```
gitmap-v28 clone <url>                          # clone into repo-name folder
gitmap-v28 clone <url> <folder-name>            # clone into custom folder
gitmap-v28 clone <url> --github-desktop         # auto-register (no prompt)
```

## Behavior

1. **Detect direct URL** — if the `source` argument starts with `https://`,
   `http://`, or `git@`, treat it as a direct URL clone (not a file path).
2. **Resolve folder name**:
   - If a second positional argument is provided, use it as the target folder.
   - Otherwise, derive the folder name from the URL: strip `.git` suffix,
     take the last path segment (e.g., `wp-alim`).
3. **Build a ScanRecord** from the URL:
   - `HTTPSUrl` / `SSHUrl` — populated from the URL.
   - `RepoName` — derived folder/repo name.
   - `Branch` — empty (let `git clone` use the default branch).
   - `RelativePath` — the target folder name.
   - `AbsolutePath` — resolved after clone.
   - `Slug` — lowercase repo name.
4. **Clone** — run `git clone <url> <folder>`.
5. **Upsert to DB** — insert or update the repo record in the `Repos` table
   using the absolute path as the key.
6. **Enqueue pending task** — create a `Clone` pending task before cloning,
   mark complete/failed after.
7. **GitHub Desktop** — auto-register with GitHub Desktop (default behavior
   for direct URL clones; no prompt).
8. **VS Code** — if the `code` CLI is on PATH, open the cloned folder in
   VS Code using `--reuse-window`. Falls back to `--new-window` if reuse
   fails (handles admin-mode conflicts).
9. **Print summary** — single-repo success/failure message.

## URL Detection

```
isDirectURL(source) → true when:
  - strings.HasPrefix(source, "https://")
  - strings.HasPrefix(source, "http://")
  - strings.HasPrefix(source, "git@")
```

This check runs BEFORE `resolveCloneShorthand` so that URLs are never
matched against the JSON/CSV/text shorthand map.

## Flag Parsing Changes

`parseCloneFlags` returns an additional `folderName string`:
- `fs.Arg(0)` → source (URL or file path)
- `fs.Arg(1)` → optional folder name (only used for direct URL clones)

## Constants

| Constant | Value |
|----------|-------|
| `MsgCloneURLCloning` | `"Cloning %s into %s...\n"` |
| `MsgCloneURLDone` | `"Cloned %s successfully.\n"` |
| `ErrCloneURLFailed` | `"Error: clone failed for %s: %v (operation: git-clone)\n"` |
| `MsgCloneDesktopPrompt` | `"Add to GitHub Desktop? (y/n): "` |
| `ErrCloneUsage` | Updated to include URL syntax |

## Database Tracking

After a successful clone:
1. Build a `ScanRecord` with the resolved absolute path.
2. Call `db.UpsertRepos([]model.ScanRecord{record})` to persist.

## Error Handling

- If the URL is unreachable or auth fails, print the git error and exit 1.
- If the target folder already exists, print an error and exit 1.
- Pending task is marked as failed with the error reason.

## Examples

### Clone from HTTPS URL

```
gitmap-v28 clone https://github.com/alimtvnetwork/wp-alim.git
```

Output:
```
Cloning wp-alim into ./wp-alim...
Cloned wp-alim successfully.
Add to GitHub Desktop? (y/n): y
  + wp-alim registered with GitHub Desktop.
```

### Clone with custom folder name

```
gitmap-v28 clone https://github.com/alimtvnetwork/wp-alim.git "my-project"
```

Output:
```
Cloning wp-alim into ./my-project...
Cloned wp-alim successfully.
Add to GitHub Desktop? (y/n): n
```

### Clone SSH URL

```
gitmap-v28 clone git@github.com:alimtvnetwork/wp-alim.git
```

## Component Mapping

| Component | File | Change |
|-----------|------|--------|
| URL detection | `cmd/clone.go` | `isDirectURL()` function |
| Direct clone flow | `cmd/clone.go` | `executeDirectClone()` function |
| VS Code open | `cmd/clonevscode.go` | `openInVSCode()`, `isVSCodeAvailable()` |
| Flag parsing | `cmd/rootflags.go` | Return `folderName` from `parseCloneFlags` |
| Record builder | `cloner/cloner.go` | `CloneSingleURL()` public function |
| Constants | `constants/constants_messages.go` | Clone + VS Code message constants |
| Constants | `constants/constants_cli.go` | Updated `ErrCloneUsage` |
| Help text | `helptext/clone.md` | Add URL examples |
| DB upsert | `store/repo.go` | Existing `UpsertRepos` (no change) |

## See Also

- [CLI Interface](02-cli-interface.md) — Full clone command specification
- [Cloner](05-cloner.md) — File-based clone behavior
- [Clone Next Flatten](87-clone-next-flatten.md) — Version iteration cloning
- [Database](16-database.md) — Repo persistence model
