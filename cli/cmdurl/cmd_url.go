package cmdurl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// UrlOptions holds parsed command line arguments and flags for url commands.
type UrlOptions struct {
	Subcommand   string
	URL          string
	IsForceFresh bool
	HasHeaders   bool
	IsJSON       bool
	IsQuiet      bool
	IsAll        bool
	Limit        int
	TTL          time.Duration
}

// UrlJsonEnvelope defines a JSON output structure when --json is passed.
type UrlJsonEnvelope struct {
	URL           string            `json:"url"`
	UrlHash       string            `json:"urlHash"`
	StatusCode    int               `json:"statusCode"`
	ContentType   string            `json:"contentType"`
	ContentLength int64             `json:"contentLength"`
	ETag          string            `json:"etag,omitempty"`
	IsCacheHit    bool              `json:"isCacheHit"`
	DurationMs    int64             `json:"durationMs"`
	Headers       map[string]string `json:"headers,omitempty"`
	Payload       string            `json:"payload,omitempty"`
}

// RunUrlCommand routes and executes url subcommands.
func RunUrlCommand(args []string) error {
	opts, err := ParseUrlArgs(args)
	if err != nil {
		return err
	}

	if opts.Subcommand == "help" {
		return printUrlHelp()
	}

	return dispatchSubcommand(opts)
}

// RunUrlCmd is an alias for RunUrlCommand.
func RunUrlCmd(args []string) error {
	return RunUrlCommand(args)
}

// ParseUrlArgs parses command-line arguments into UrlOptions.
func ParseUrlArgs(args []string) (UrlOptions, error) {
	opts := defaultUrlOptions()
	if len(args) == 0 {
		opts.Subcommand = "help"

		return opts, nil
	}

	return parseSubcommandAndFlags(args, opts)
}

func defaultUrlOptions() UrlOptions {
	return UrlOptions{
		Subcommand: "read",
		Limit:      50,
		TTL:        24 * time.Hour,
	}
}

func parseSubcommandAndFlags(args []string, opts UrlOptions) (UrlOptions, error) {
	first := args[0]
	if first == "-h" || first == "--help" || first == "help" {
		opts.Subcommand = "help"

		return opts, nil
	}

	return parseUrlOrSubcommand(first, args[1:], opts)
}

func parseUrlOrSubcommand(first string, rest []string, opts UrlOptions) (UrlOptions, error) {
	if strings.HasPrefix(first, "http://") || strings.HasPrefix(first, "https://") {
		opts.Subcommand = "read"
		opts.URL = first

		return parseFlagList(rest, opts)
	}

	opts.Subcommand = first

	return parseFlagList(rest, opts)
}

func parseFlagList(tokens []string, opts UrlOptions) (UrlOptions, error) {
	for i := 0; i < len(tokens); i++ {
		advance, err := processFlagToken(tokens, i, &opts)
		if err != nil {
			return opts, err
		}

		i += advance
	}

	return opts, nil
}

func processFlagToken(tokens []string, i int, opts *UrlOptions) (int, error) {
	token := tokens[i]
	if isSimpleFlag(token) {
		*opts = applySimpleFlag(token, *opts)

		return 0, nil
	}

	return processComplexFlag(tokens, i, token, opts)
}

func processComplexFlag(tokens []string, i int, token string, opts *UrlOptions) (int, error) {
	advance, err := parseValueFlag(tokens, i, opts)
	if err != nil || advance > 0 {
		return advance, err
	}

	*opts = applyPositionalArg(token, *opts)

	return 0, nil
}

func isSimpleFlag(token string) bool {
	return token == "--fresh" || token == "--headers" || token == "--json" ||
		token == "-q" || token == "--quiet" || token == "--all"
}

func applySimpleFlag(token string, opts UrlOptions) UrlOptions {
	opts = applyOutputFlags(token, opts)

	return applyControlFlags(token, opts)
}

func applyOutputFlags(token string, opts UrlOptions) UrlOptions {
	if token == "--headers" {
		opts.HasHeaders = true
	}

	if token == "--json" {
		opts.IsJSON = true
	}

	if token == "-q" || token == "--quiet" {
		opts.IsQuiet = true
	}

	return opts
}

func applyControlFlags(token string, opts UrlOptions) UrlOptions {
	if token == "--fresh" {
		opts.IsForceFresh = true
	}

	if token == "--all" {
		opts.IsAll = true
	}

	return opts
}

func applyPositionalArg(token string, opts UrlOptions) UrlOptions {
	if opts.URL == "" && !strings.HasPrefix(token, "-") {
		opts.URL = token
	}

	return opts
}

func parseValueFlag(tokens []string, i int, opts *UrlOptions) (int, error) {
	if tokens[i] == "--ttl" {
		return parseTtlFlag(tokens, i, opts)
	}

	if tokens[i] == "--limit" {
		return parseLimitFlag(tokens, i, opts)
	}

	return 0, nil
}

func parseTtlFlag(tokens []string, i int, opts *UrlOptions) (int, error) {
	if i+1 >= len(tokens) {
		return 0, nil
	}

	dur, err := time.ParseDuration(tokens[i+1])
	if err != nil {
		return 0, apperror.WrapSimple(err, "cmdurl.parseTTL")
	}

	opts.TTL = dur

	return 1, nil
}

func parseLimitFlag(tokens []string, i int, opts *UrlOptions) (int, error) {
	if i+1 >= len(tokens) {
		return 0, nil
	}

	lim, err := strconv.Atoi(tokens[i+1])
	if err != nil {
		return 0, apperror.WrapSimple(err, "cmdurl.parseLimit")
	}

	opts.Limit = lim

	return 1, nil
}

func dispatchSubcommand(opts UrlOptions) error {
	handler, hasHandler := urlDispatchMap()[opts.Subcommand]
	if !hasHandler {
		return apperror.NewSimple("unknown url subcommand: "+opts.Subcommand, "E1301")
	}

	return handler(opts)
}

func urlDispatchMap() map[string]func(UrlOptions) error {
	return map[string]func(UrlOptions) error{
		"read": runUrlRead, "get": runUrlRead,
		"ls": runUrlList, "list": runUrlList,
		"info": runUrlInfo, "purge": runUrlPurge,
		"rm": runUrlRemove, "delete": runUrlRemove,
	}
}

func runUrlRead(opts UrlOptions) error {
	if opts.URL == "" {
		return apperror.NewSimple("url read requires a URL argument", "E1302")
	}

	db, err := store.OpenUrlCacheDB("")
	if err != nil {
		return err
	}

	defer db.Close()

	return serveOrFetchUrl(db, opts)
}

func serveOrFetchUrl(db *store.UrlCacheSplitDB, opts UrlOptions) error {
	if opts.IsForceFresh {
		return fetchAndCacheUrl(db, opts)
	}

	hasHit, err := tryServeFromCache(db, opts)
	if err == nil && hasHit {
		return nil
	}

	return fetchAndCacheUrl(db, opts)
}

func tryServeFromCache(db *store.UrlCacheSplitDB, opts UrlOptions) (bool, error) {
	start := time.Now()
	entry, hasHit, err := findFreshCacheEntry(db, opts.URL)
	if err != nil || !hasHit {
		return false, err
	}

	dur := time.Since(start)
	logHitTelemetry(db, entry, opts, dur)

	return true, emitReadResult(opts, *entry, true, dur.Milliseconds())
}

func findFreshCacheEntry(db *store.UrlCacheSplitDB, url string) (*store.UrlCacheEntry, bool, error) {
	entry, hasHit, err := db.GetUrlCache(url)
	if err != nil || !hasHit {
		return nil, false, err
	}

	if isCacheExpired(entry.ExpiresAt) {
		return nil, false, nil
	}

	return entry, true, nil
}

func logHitTelemetry(db *store.UrlCacheSplitDB, entry *store.UrlCacheEntry, opts UrlOptions, dur time.Duration) {
	logAccessAudit(db, entry, true, dur)
	logDiagnostic(opts.IsQuiet, formatHitLog(dur, opts.URL, entry.ContentLength))
}

func isCacheExpired(expiresAt int64) bool {
	return expiresAt > 0 && expiresAt < time.Now().Unix()
}

func formatHitLog(dur time.Duration, url string, size int64) string {
	ms := float64(dur.Microseconds()) / 1000.0

	return fmt.Sprintf("[CACHE HIT %.2fms] %s (%d bytes)", ms, url, size)
}

func fetchAndCacheUrl(db *store.UrlCacheSplitDB, opts UrlOptions) error {
	logDiagnostic(opts.IsQuiet, fmt.Sprintf("[CACHE MISS] Fetching %s...", opts.URL))
	start := time.Now()
	resp, payload, err := executeHttpGet(opts.URL)
	if err != nil {
		return err
	}

	dur := time.Since(start)
	entry := buildCacheEntry(opts, resp, payload, dur.Milliseconds())
	_ = db.SaveUrlCache(entry)
	logSaveTelemetry(db, &entry, opts, dur, len(payload))

	return emitReadResult(opts, entry, false, dur.Milliseconds())
}

func logSaveTelemetry(db *store.UrlCacheSplitDB, entry *store.UrlCacheEntry, opts UrlOptions, dur time.Duration, byteCount int) {
	logAccessAudit(db, entry, false, dur)
	logDiagnostic(opts.IsQuiet, fmt.Sprintf("[CACHE SAVED] %s (%d bytes, %dms)", opts.URL, byteCount, dur.Milliseconds()))
}

func executeHttpGet(rawUrl string) (*http.Response, []byte, error) {
	resp, err := sendHttpGetRequest(rawUrl)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	return readHttpResponseBody(resp)
}

func sendHttpGetRequest(rawUrl string) (*http.Response, error) {
	req, err := createHttpGetRequest(rawUrl)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdurl.executeHttpGet.do")
	}

	return resp, nil
}

func createHttpGetRequest(rawUrl string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdurl.createHttpGetRequest")
	}

	req.Header.Set("User-Agent", "gitmap/1.0 (url-reader)")

	return req, nil
}

func readHttpResponseBody(resp *http.Response) (*http.Response, []byte, error) {
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "cmdurl.readHttpResponseBody")
	}

	return resp, payload, nil
}

func buildCacheEntry(opts UrlOptions, resp *http.Response, payload []byte, durMs int64) store.UrlCacheEntry {
	ttl := resolveTTL(opts.TTL)
	now := time.Now().Unix()
	entry := initBaseCacheEntry(opts.URL, payload, durMs, now, ttl)
	headerBytes, _ := json.Marshal(resp.Header)
	entry.StatusCode = resp.StatusCode
	entry.ETag = resp.Header.Get("ETag")
	entry.LastModified = resp.Header.Get("Last-Modified")
	entry.ContentType = resp.Header.Get(constants.HeaderContentType)
	entry.HeadersJson = string(headerBytes)

	return entry
}

func resolveTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return 24 * time.Hour
	}

	return ttl
}

func initBaseCacheEntry(url string, payload []byte, durMs, now int64, ttl time.Duration) store.UrlCacheEntry {
	return store.UrlCacheEntry{
		UrlHash:       store.HashUrl(url),
		RawUrl:        url,
		UrlText:       url,
		HttpMethod:    http.MethodGet,
		ContentLength: int64(len(payload)),
		PayloadBlob:   payload,
		DurationMs:    durMs,
		ExpiresAt:     now + int64(ttl.Seconds()),
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func emitReadResult(opts UrlOptions, entry store.UrlCacheEntry, isCacheHit bool, durMs int64) error {
	if opts.HasHeaders {
		return emitHeadersOnly(entry)
	}

	if opts.IsJSON {
		return emitJsonEnvelope(opts.URL, entry, isCacheHit, durMs)
	}

	return writePayloadBytes(entry.PayloadBlob)
}

func writePayloadBytes(blob []byte) error {
	_, err := output.Raw().Write(blob)
	if err != nil {
		return apperror.WrapSimple(err, "cmdurl.emitReadResult")
	}

	return nil
}

func emitHeadersOnly(entry store.UrlCacheEntry) error {
	printFixedHeaders(entry)
	printCustomHeaders(entry.HeadersJson)

	return nil
}

func printFixedHeaders(entry store.UrlCacheEntry) {
	fmt.Fprintf(output.Raw(), "HTTP/1.1 %d\n", entry.StatusCode)
	fmt.Fprintf(output.Raw(), constants.HeaderContentType+": %s\n", entry.ContentType)
	fmt.Fprintf(output.Raw(), "Content-Length: %d\n", entry.ContentLength)
	if entry.ETag != "" {
		fmt.Fprintf(output.Raw(), "ETag: %s\n", entry.ETag)
	}
}

func printCustomHeaders(headersJson string) {
	var headers map[string][]string
	_ = json.Unmarshal([]byte(headersJson), &headers)
	for k, vals := range headers {
		for _, v := range vals {
			fmt.Fprintf(output.Raw(), "%s: %s\n", k, v)
		}
	}
}

func emitJsonEnvelope(url string, entry store.UrlCacheEntry, isCacheHit bool, durMs int64) error {
	env := makeJsonEnvelope(url, entry, isCacheHit, durMs)
	bytes, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdurl.emitJsonEnvelope")
	}

	_, writeErr := output.Raw().Write(append(bytes, '\n'))

	return writeErr
}

func makeJsonEnvelope(url string, entry store.UrlCacheEntry, isCacheHit bool, durMs int64) UrlJsonEnvelope {
	return UrlJsonEnvelope{
		URL:           url,
		UrlHash:       entry.UrlHash,
		StatusCode:    entry.StatusCode,
		ContentType:   entry.ContentType,
		ContentLength: entry.ContentLength,
		ETag:          entry.ETag,
		IsCacheHit:    isCacheHit,
		DurationMs:    durMs,
		Payload:       string(entry.PayloadBlob),
	}
}

func logAccessAudit(db *store.UrlCacheSplitDB, entry *store.UrlCacheEntry, isCacheHit bool, dur time.Duration) {
	log := makeAuditLogEntry(entry, isCacheHit, dur)
	_ = db.LogUrlAccess(log)
	if isCacheHit && entry.UrlCacheEntryId > 0 {
		_ = db.TouchAccess(entry.UrlCacheEntryId)
	}
}

func makeAuditLogEntry(entry *store.UrlCacheEntry, isCacheHit bool, dur time.Duration) store.UrlAccessLog {
	return store.UrlAccessLog{
		UrlHash:      entry.UrlHash,
		RawUrl:       entry.EffectiveURL(),
		IsCacheHit:   isCacheHit,
		StatusCode:   entry.StatusCode,
		BytesServed:  entry.ContentLength,
		DurationMs:   dur.Milliseconds(),
		CallerSource: "cli",
	}
}

func runUrlList(opts UrlOptions) error {
	db, err := store.OpenUrlCacheDB("")
	if err != nil {
		return err
	}

	defer db.Close()

	return renderUrlList(db, opts)
}

func renderUrlList(db *store.UrlCacheSplitDB, opts UrlOptions) error {
	entries, err := db.ListUrlCache(opts.Limit, 0)
	if err != nil {
		return err
	}

	if opts.IsJSON {
		return printJsonList(entries)
	}

	return printTableList(entries)
}

func printJsonList(entries []store.UrlCacheEntry) error {
	bytes, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdurl.printJsonList")
	}

	_, writeErr := output.Raw().Write(append(bytes, '\n'))

	return writeErr
}

func printTableList(entries []store.UrlCacheEntry) error {
	now := time.Now().Unix()
	fmt.Fprintf(output.UI(), "%-10s  %-8s  %-12s  %-6s  %-8s  %s\n", "HASH", "SIZE", "TYPE", "HITS", "EXPIRES", "URL")
	for _, e := range entries {
		printTableRow(e, now)
	}

	return nil
}

func printTableRow(e store.UrlCacheEntry, now int64) {
	shortHash := truncateHash(e.UrlHash)
	expiresText := formatExpiresText(e.ExpiresAt, now)
	fmt.Fprintf(output.UI(), "%-10s  %-8d  %-12s  %-6d  %-8s  %s\n", shortHash, e.ContentLength, truncateType(e.ContentType), e.AccessCount, expiresText, e.EffectiveURL())
}

func truncateHash(h string) string {
	if len(h) > 8 {
		return h[:8]
	}

	return h
}

func formatExpiresText(expiresAt, now int64) string {
	if expiresAt <= 0 {
		return "never"
	}

	diff := expiresAt - now
	if diff <= 0 {
		return "expired"
	}

	return fmt.Sprintf("%dm", diff/60)
}

func truncateType(ct string) string {
	if len(ct) > 12 {
		return ct[:12]
	}

	return ct
}

func runUrlInfo(opts UrlOptions) error {
	if opts.URL == "" {
		return apperror.NewSimple("url info requires a URL argument", "E1303")
	}

	db, err := store.OpenUrlCacheDB("")
	if err != nil {
		return err
	}

	defer db.Close()

	return fetchAndEmitInfo(db, opts)
}

func fetchAndEmitInfo(db *store.UrlCacheSplitDB, opts UrlOptions) error {
	entry, hasHit, err := db.GetUrlCache(opts.URL)
	if err != nil {
		return err
	}

	if !hasHit {
		return apperror.NewSimple("url not found in cache: "+opts.URL, "E1304")
	}

	return emitInfoResult(opts, *entry)
}

func emitInfoResult(opts UrlOptions, entry store.UrlCacheEntry) error {
	if opts.IsJSON {
		return printJsonInfo(entry)
	}

	return printFormattedInfo(entry)
}

func printJsonInfo(entry store.UrlCacheEntry) error {
	bytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdurl.printJsonInfo")
	}

	_, writeErr := output.Raw().Write(append(bytes, '\n'))

	return writeErr
}

func printFormattedInfo(entry store.UrlCacheEntry) error {
	fmt.Fprintf(output.UI(), "URL:          %s\n", entry.EffectiveURL())
	fmt.Fprintf(output.UI(), "Hash:         %s\n", entry.UrlHash)
	fmt.Fprintf(output.UI(), "Status:       %d\n", entry.StatusCode)
	fmt.Fprintf(output.UI(), constants.HeaderContentType+": %s\n", entry.ContentType)
	fmt.Fprintf(output.UI(), "Size:         %d bytes\n", entry.ContentLength)
	fmt.Fprintf(output.UI(), "Access Count: %d\n", entry.AccessCount)
	fmt.Fprintf(output.UI(), "Expires At:   %s\n", time.Unix(entry.ExpiresAt, 0).Format(time.RFC3339))

	return nil
}

func runUrlPurge(opts UrlOptions) error {
	db, err := store.OpenUrlCacheDB("")
	if err != nil {
		return err
	}

	defer db.Close()

	return executeUrlPurge(db, opts.IsAll)
}

func executeUrlPurge(db *store.UrlCacheSplitDB, isAll bool) error {
	targetTime := resolvePurgeTargetTime(isAll)
	purgedCount, err := db.PurgeExpired(targetTime)
	if err != nil {
		return err
	}

	fmt.Fprintf(output.UIErr(), "[cache:purge] Purged %d entries from url_cache.db\n", purgedCount)

	return nil
}

func resolvePurgeTargetTime(isAll bool) int64 {
	if isAll {
		return time.Now().Add(100 * 365 * 24 * time.Hour).Unix()
	}

	return time.Now().Unix()
}

func runUrlRemove(opts UrlOptions) error {
	if opts.URL == "" {
		return apperror.NewSimple("url rm requires a URL argument", "E1305")
	}

	db, err := store.OpenUrlCacheDB("")
	if err != nil {
		return err
	}

	defer db.Close()

	return executeUrlRemove(db, opts.URL)
}

func executeUrlRemove(db *store.UrlCacheSplitDB, rawUrl string) error {
	if err := db.DeleteUrlCache(rawUrl); err != nil {
		return err
	}

	fmt.Fprintf(output.UIErr(), "[cache:rm] Removed %s from url_cache.db\n", rawUrl)

	return nil
}

func logDiagnostic(isQuiet bool, msg string) {
	if !isQuiet {
		fmt.Fprintln(output.UIErr(), msg)
	}
}

const urlHelpText = `gitmap url — URL streaming and caching utility

Usage:
  gitmap url read <url> [flags]    Fetch and stream URL content
  gitmap url <url> [flags]         Shorthand for url read
  gitmap get <url> [flags]         Shorthand for url read
  gitmap read-url <url> [flags]    Shorthand for url read
  gitmap url ls [flags]            List cached URLs
  gitmap url info <url> [flags]    Inspect cache entry details
  gitmap url purge [flags]         Purge expired cache entries
  gitmap url rm <url>              Remove URL from cache

Flags:
  --fresh          Force re-fetch bypassing local cache
  --ttl <dur>      Cache expiration duration (default: 24h)
  --headers        Output HTTP response headers only
  --json           Format output as JSON envelope
  --quiet, -q      Suppress diagnostic log messages
  --all            Purge all entries (used with url purge)
  --limit <n>      Limit number of results (used with url ls)
`

func printUrlHelp() error {
	fmt.Fprint(output.UI(), urlHelpText)

	return nil
}
