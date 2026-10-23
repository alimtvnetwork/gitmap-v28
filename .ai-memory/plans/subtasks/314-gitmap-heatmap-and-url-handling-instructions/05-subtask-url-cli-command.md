# Subtask 05: CLI URL Read Command & Root Utility Dispatch

> **Subtask ID:** `SUBTASK-314-05`  
> **Parent Task:** `314-gitmap-heatmap-and-url-handling-instructions`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/cmdurl/cmd_url.go` (new file)  
> - `cli/cmd/rootutility.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/02-component-spec.md`  

---

## 1. Objective

Implement the CLI URL interface in `cli/cmdurl/cmd_url.go` and wire command routing in `cli/cmd/rootutility.go`. The command supports `gitmap url read <url>` (with shorthand aliases `gitmap url <url>`, `gitmap read-url <url>`, and `gitmap get <url>`), plus subcommands `ls`, `info`, `purge`, and `rm`. The command guarantees byte-faithful streaming to stdout via `output.Raw()` while routing all diagnostic, telemetry, and cache-hit notifications to `output.UIErr()`.

---

## 2. Action Checklist

- [ ] 1. Author `cli/cmdurl/cmd_url.go`:
  - [ ] Package declaration `package cmdurl`.
  - [ ] Implement entry point `RunUrlCmd(args []string) error`.
  - [ ] Implement argument and flag parser `parseUrlArgs(args []string)`:
    - [ ] Detect implicit `read` subcommand when arg starts with `http://` or `https://`.
    - [ ] Parse flags: `--fresh` (`IsForceFresh`), `--ttl` (`TTL`), `--headers` (`HasHeaders`), `--json` (`IsJSON`), `--quiet` (`IsQuiet`).
  - [ ] Subcommand router (switch or dispatch map):
    - [ ] `read`: lookup cache in `repodb/url_cache.db`. If hit and not expired and `!IsForceFresh`, stream cached bytes to `output.Raw()`. If miss or forced fresh, invoke `downloader.Download()`, save to cache, and stream bytes to `output.Raw()`.
    - [ ] `ls`: query `ListUrlCache()` and print formatted table of cached URLs to `output.UI()` or JSON to `output.Raw()`.
    - [ ] `info`: query `GetUrlCache()`, print metadata, SHA-256 hash, ETag, headers, and hit stats to `output.UI()`.
    - [ ] `purge`: call `PurgeExpiredUrlCache()`, emit count of purged entries to `output.UIErr()`.
    - [ ] `rm`: call `DeleteUrlCache(urlText)`, confirm deletion to `output.UIErr()`.
  - [ ] Output channel separation:
    - [ ] **STRICT INVARIANT:** Unfiltered payload writes via `output.Raw().Write()` or `io.Copy(output.Raw(), reader)`.
    - [ ] Informational logs and diagnostics via `fmt.Fprintf(output.UIErr(), ...)` only.
  - [ ] Function constraints: all functions <= 15 lines, guard clauses, zero nested ifs.

- [ ] 2. Update `cli/cmd/rootutility.go`:
  - [ ] Import `github.com/alimtvnetwork/gitmap-v28/cli/cmdurl`.
  - [ ] In `utilityCoreEntries()` or `utilityDispatchEntries()`, register dispatch entries:
    - [ ] `{[]string{"url"}, func() error { return cmdurl.RunUrlCmd(argsTail()) }}`
    - [ ] `{[]string{"read-url", "readurl"}, func() error { return cmdurl.RunUrlCmd(append([]string{"read"}, argsTail()...)) }}`
    - [ ] `{[]string{"get"}, func() error { return cmdurl.RunUrlCmd(append([]string{"read"}, argsTail()...)) }}`
  - [ ] Keep functions <= 15 lines.

- [ ] 3. Quality & Coding Guideline Compliance:
  - [ ] Strictly relative paths in code and comments.
  - [ ] Strictly lowercase filenames (`cmd_url.go`, `rootutility.go`).
  - [ ] Positive booleans only: `IsForceFresh`, `HasHeaders`, `IsJSON`, `IsQuiet`, `HasFound`, `IsCacheHit`.
  - [ ] AppError error wrapping on all failure branches.
  - [ ] Zero git commands, zero builds, zero tests.

---

## 3. Concrete Code Specifications

### Specification: `cli/cmdurl/cmd_url.go`

```go
package cmdurl

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/downloader"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// UrlOptions holds parsed command line arguments and flags.
type UrlOptions struct {
	Subcommand   string
	URL          string
	IsForceFresh bool
	HasHeaders   bool
	IsJSON       bool
	IsQuiet      bool
	TTL          time.Duration
}

// RunUrlCmd routes and executes url subcommands.
func RunUrlCmd(args []string) error {
	opts, err := parseUrlArgs(args)
	if err != nil {
		return err
	}

	return dispatchSubcommand(opts)
}

func dispatchSubcommand(opts UrlOptions) error {
	switch opts.Subcommand {
	case "read":
		return runUrlRead(opts)
	case "ls":
		return runUrlList(opts)
	case "info":
		return runUrlInfo(opts)
	case "purge":
		return runUrlPurge(opts)
	case "rm":
		return runUrlRemove(opts)
	default:
		return apperror.NewSimple("unknown url subcommand: "+opts.Subcommand, "E1301")
	}
}

func runUrlRead(opts UrlOptions) error {
	cacheDb, err := store.OpenUrlCacheDB("")
	if err != nil {
		return err
	}
	defer cacheDb.Close()

	if !opts.IsForceFresh {
		if hasHit, err := tryServeFromCache(cacheDb, opts); err == nil && hasHit {
			return nil
		}
	}

	return fetchAndCacheUrl(cacheDb, opts)
}

func tryServeFromCache(db *store.UrlCacheDB, opts UrlOptions) (bool, error) {
	entry, hasHit, err := db.GetUrlCache(opts.URL)
	if err != nil || !hasHit {
		return false, err
	}

	isExpired := entry.ExpiresAt > 0 && entry.ExpiresAt < time.Now().Unix()
	if isExpired {
		return false, nil
	}

	logDiagnostic(opts.IsQuiet, fmt.Sprintf("[cache:hit] Serving %s from url_cache.db (%d bytes)", opts.URL, entry.ContentLength))
	_, writeErr := output.Raw().Write(entry.PayloadBlob)
	if writeErr != nil {
		return false, apperror.WrapSimple(writeErr, "cmdurl.tryServeFromCache.write")
	}

	return true, nil
}

func fetchAndCacheUrl(db *store.UrlCacheDB, opts UrlOptions) error {
	logDiagnostic(opts.IsQuiet, fmt.Sprintf("[cache:miss] Fetching %s...", opts.URL))
	dlOpts := downloader.DownloadOptions{
		URL:          opts.URL,
		IsForceFresh: opts.IsForceFresh,
		IsQuiet:      opts.IsQuiet,
		IsJSON:       opts.IsJSON,
	}

	result, err := downloader.Download(dlOpts)
	if err != nil {
		return err
	}

	return streamAndSaveResult(db, opts, result)
}

func streamAndSaveResult(db *store.UrlCacheDB, opts UrlOptions, res downloader.DownloadResult) error {
	payload, err := io.ReadAll(res.PayloadReader)
	if err != nil {
		return apperror.WrapSimple(err, "cmdurl.streamAndSaveResult.read")
	}

	ttl := opts.TTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	entry := store.UrlCacheEntry{
		UrlHash:       store.HashUrl(opts.URL),
		UrlText:       opts.URL,
		ContentLength: int64(len(payload)),
		PayloadBlob:   payload,
		ExpiresAt:     time.Now().Add(ttl).Unix(),
		IsActive:      true,
	}

	_ = db.SaveUrlCache(entry)
	_, writeErr := output.Raw().Write(payload)

	return writeErr
}

func logDiagnostic(isQuiet bool, msg string) {
	if !isQuiet {
		fmt.Fprintln(output.UIErr(), msg)
	}
}
```

### Specification: `cli/cmd/rootutility.go` Registration

```go
// Add in utilityCoreEntries() or utilityDispatchEntries():
{[]string{"url"}, func() error { return cmdurl.RunUrlCmd(argsTail()) }},
{[]string{"read-url", "readurl"}, func() error { return cmdurl.RunUrlCmd(append([]string{"read"}, argsTail()...)) }},
{[]string{"get"}, func() error { return cmdurl.RunUrlCmd(append([]string{"read"}, argsTail()...)) }},
```

---

## 4. Verification & Done Criteria

1. Invocations `gitmap url <url>`, `gitmap url read <url>`, `gitmap read-url <url>`, and `gitmap get <url>` route cleanly to `RunUrlCmd`.
2. Raw payloads are written exclusively to `output.Raw()`, preventing ANSI/binary stream corruption.
3. Diagnostic notifications are emitted exclusively to `output.UIErr()`.
4. Cache hits in `url_cache.db` respond in sub-millisecond (<0.5ms) latency without issuing network calls.
5. All Go functions conform to <= 15 lines with guard clauses and positive booleans.
