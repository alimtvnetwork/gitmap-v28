# replace

Replaces exact text in repository files and tracks changes in the SplitDB for undo history.

## Why use this instead of shell scripts?
Instead of writing complex PowerShell replacement scripts like:
```powershell
$content = Get-Content gitmap/searcher/db_search.go
$content = $content -replace 'func SearchRepoDB\(ctx context.Context, db \*sql.DB, query string, useCache bool\)', 'func SearchRepoDB(ctx context.Context, db *sql.DB, query string, limit int, useCache bool)'
Set-Content gitmap/searcher/db_search.go -Value $content
```
or running a dangerous multi-file sed command:
`find . -type f -exec sed -i 's/func run(/func execute(/g' {} +`

You should use GitMap's transactional replace:
`gitmap replace "func SearchRepoDB(ctx context.Context, db *sql.DB, query string, useCache bool)" "func SearchRepoDB(ctx context.Context, db *sql.DB, query string, limit int, useCache bool)"`

## Examples
```bash
gitmap replace "func run(" "func execute("
gitmap replace "oldString" "newString"
```
