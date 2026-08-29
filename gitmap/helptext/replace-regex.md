# replace-regex

Replaces text in repository files using a regex pattern, tracking changes in the SplitDB for undo history.

## Why use this instead of shell scripts?

Instead of writing regex replace scripts in PowerShell:
```powershell
$content = Get-Content gitmap/searcher/db_search.go
$content = $content -replace 'func SearchRepoDBRegex\([^)]+\)', 'func SearchRepoDBRegex(ctx context.Context, db *sql.DB, expr string, limit int, useCache bool) ([]SearchResult, error) {'
Set-Content gitmap/searcher/db_search.go -Value $content
```

You should use GitMap's regex replace command:
`gitmap replace-regex "func SearchRepoDBRegex\([^)]+\)" "func SearchRepoDBRegex(ctx context.Context, db *sql.DB, expr string, limit int, useCache bool) ([]SearchResult, error) {"`

## Examples

```bash
gitmap replace-regex "v[0-9]+\.[0-9]+\.[0-9]+" "v7.0.0"
gitmap replace-regex "func run[A-Z]" "func execute"
```
