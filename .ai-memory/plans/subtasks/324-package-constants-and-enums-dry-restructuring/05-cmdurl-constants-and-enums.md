# Subtask 324-05: cmdurl Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdurl` (2 files: `cmd_url.go`, `cmd_url_test.go`)
- Leaf Packages to Create:
  - `cli/cmdurl/enums/`
    - `url_enums.go`
    - `enums_test.go`
  - `cli/cmdurl/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `UrlSubcommandType`: Subcommands (`read`, `get`, `ls`, `list`, `info`, `purge`, `rm`, `delete`, `help`).
2. `UrlOutputFormatType`: Output format (`raw`, `headers`, `json`, `table`).
3. `UrlCacheStatusType`: Cache telemetry status (`HIT`, `MISS`, `SAVED`, `PURGED`, `REMOVED`).
4. `UrlSchemeType`: Supported schemes (`http`, `https`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdUrl = "url"`, `CmdReadUrl = "read-url"`, `CmdReadUrlAlias = "readurl"`, `CmdGet = "get"`, `CmdGetUrl = "get-url"`.
- Subcommands: `SubcmdRead = "read"`, `SubcmdGet = "get"`, `SubcmdList = "ls"`, `SubcmdListAlias = "list"`, `SubcmdInfo = "info"`, `SubcmdPurge = "purge"`, `SubcmdRemove = "rm"`, `SubcmdRemoveAliasDelete = "delete"`, `SubcmdHelp = "help"`.
- Flags: `--fresh`, `--headers`, `--json`, `--quiet`, `-q`, `--all`, `--ttl`, `--limit`, `--help`, `-h`.
- Scheme prefixes: `PrefixHttp = "http://"`, `PrefixHttps = "https://"`.
- Defaults & limits: `DefaultLimit = 50`, `DefaultTTL = 24 * time.Hour`, `DefaultHttpTimeout = 30 * time.Second`, `DefaultUserAgent = "gitmap/1.0 (url-reader)"`, `DefaultDbFile = "url_cache.db"`.
- HTTP headers: `HeaderUserAgent`, `HeaderContentType`, `HeaderContentLength`, `HeaderETag`, `HeaderLastModified`.
- Column headers & formats: `ColHeaderHash`, `ColHeaderSize`, `ColHeaderType`, `ColHeaderHits`, `ColHeaderExpires`, `ColHeaderURL`, table templates, header templates.
- Telemetry & error codes: `MsgHitLogFmt`, `MsgMissLogFmt`, `MsgSavedLogFmt`, `MsgPurgedLogFmt`, `MsgRemovedLogFmt`, `ErrCodeE1301`–`E1305`, error messages, and operation tags.
- Exit codes: `ExitSuccess = 0`, `ExitGeneralError = 1`, `ExitUsageError = 2`.
- Helpers: `AllSimpleFlags() []string`, `AllSubcommandTokens() []string`.

## Consuming Callers to Refactor
- `cli/cmdurl/cmd_url.go`: Replace inline strings and numbers across flag parsing, HTTP dispatch, and cache telemetry with `constants.*` and `enums.*`. Maintain leaf error sovereignty.
- `cli/cmdurl/cmd_url_test.go`: Refactor to use `constants.*` and `enums.*`.
- Caller integration: `cli/cmd/rootutility.go`: Import `cli/cmdurl/constants` and dispatch `constants.CmdUrl`, `constants.CmdReadUrl`, `constants.CmdReadUrlAlias`, etc.

## Verification
- Unit tests in `cli/cmdurl/enums/enums_test.go` and `cli/cmdurl/constants/constants_test.go`.
- Unit tests in `cli/cmdurl/cmd_url_test.go`.
- `go vet ./cli/cmdurl/...`
- `go test -count=1 ./cli/cmdurl/...`
