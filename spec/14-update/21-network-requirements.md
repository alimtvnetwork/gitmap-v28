# Network Requirements

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

Define the HTTP client configuration, retry policies, proxy support, and error handling for all network operations in the self-update and updater binary flows. Any AI implementing the update system must configure networking according to these requirements.

---

## HTTP Client Configuration

### Timeouts

| Operation | Timeout | Rationale |
|-----------|---------|-----------|
| TCP connect | 30 seconds | Generous for slow networks |
| TLS handshake | 30 seconds | Included in connect timeout |
| Response headers | 30 seconds | Detect unresponsive servers |
| Download body | 5 minutes | Large binaries on slow connections |
| GitHub API request | 30 seconds | Small JSON payloads |

### Go Implementation

```go
import (
    "net"
    "net/http"
    "time"
)

func NewHTTPClient() *http.Client {
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout:   30 * time.Second,
        ResponseHeaderTimeout:  30 * time.Second,
        MaxIdleConns:           10,
        IdleConnTimeout:        90 * time.Second,
    }

    // Proxy: automatically uses HTTP_PROXY, HTTPS_PROXY, NO_PROXY
    // via http.ProxyFromEnvironment (Go default)

    return &http.Client{
        Transport: transport,
        Timeout:   5 * time.Minute, // overall request timeout
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("too many redirects (max 10)")
            }
            return nil
        },
    }
}
```

---

## Proxy Support

The HTTP client MUST respect standard proxy environment variables:

| Variable | Purpose |
|----------|---------|
| `HTTP_PROXY` | Proxy for HTTP requests |
| `HTTPS_PROXY` | Proxy for HTTPS requests |
| `NO_PROXY` | Comma-separated list of hosts to bypass |

In Go, `http.ProxyFromEnvironment` handles this automatically when using the default transport. **Do not override proxy settings** — let the user's environment control proxying.

### Bash (curl/wget)

```bash

# curl automatically respects HTTP_PROXY/HTTPS_PROXY

curl -fsSL "$url" -o "$output"

# wget also respects proxy env vars

wget -q "$url" -O "$output"
```

---

## Retry Policy

All network requests MUST implement retry with exponential backoff:

| Parameter | Value |
|-----------|-------|
| Max attempts | 3 |
| Initial delay | 1 second |
| Backoff multiplier | 2× |
| Max delay | 4 seconds |
| Retry on | Connection errors, 5xx responses, 429 (rate limited) |
| Do NOT retry on | 4xx (except 429), successful responses, checksum mismatches |

### Go Implementation

```go
func downloadWithRetry(client *http.Client, url, dest string) error {
    delays := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

    var lastErr error
    for attempt, delay := range delays {
        err := downloadFile(client, url, dest)
        if err == nil {
            return nil
        }
        lastErr = err

        if attempt < len(delays)-1 {
            fmt.Printf("  ⟳ Retry %d/%d in %v: %v\n", attempt+1, len(delays)-1, delay, err)
            time.Sleep(delay)
        }
    }
    return fmt.Errorf("download failed after %d attempts: %w", len(delays), lastErr)
}
```

---

## User-Agent

All HTTP requests MUST include a `User-Agent` header:

```
User-Agent: <binary>-updater/<version>
```

Example: `User-Agent: gitmap-updater/1.3.0`

**Why**: GitHub API rate limits are more generous for requests with a User-Agent. Requests without one may be deprioritized or blocked.

---

## GitHub API Rate Limiting

The GitHub API has rate limits:

| Auth Level | Rate Limit |
|------------|-----------|
| Unauthenticated | 60 requests/hour per IP |
| Authenticated (PAT) | 5,000 requests/hour |

The updater uses **unauthenticated** requests by default. For the `releases/latest` endpoint, a single API call per update check is sufficient.

### Handling Rate Limits

```go
if resp.StatusCode == 403 || resp.StatusCode == 429 {
    retryAfter := resp.Header.Get("Retry-After")
    return fmt.Errorf(
        "GitHub API rate limited (status %d). Try again after %s, or use --version to skip API lookup",
        resp.StatusCode, retryAfter,
    )
}
```

---

## TLS / Certificate Verification

- Use the **system certificate store** — do not bundle custom CA certificates
- Do **not** disable TLS verification (`InsecureSkipVerify: false`)
- On corporate networks with MITM proxies, users must add the proxy CA to their system trust store

---

## Download Progress

For large file downloads (binary archives), provide progress indication:

```go
type ProgressWriter struct {
    Total      int64
    Downloaded int64
    LastPrint  time.Time
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
    pw.Downloaded += int64(len(p))
    if time.Since(pw.LastPrint) > 500*time.Millisecond {
        pct := float64(pw.Downloaded) / float64(pw.Total) * 100
        fmt.Printf("\r  Downloading: %.1f%%", pct)
        pw.LastPrint = time.Now()
    }
    return len(p), nil
}
```

For non-interactive environments (piped output), skip progress display:

```go
if !term.IsTerminal(int(os.Stdout.Fd())) {
    // Skip progress bar in non-interactive mode
}
```

---

## GitHub Release URL Patterns

GitHub release asset URLs redirect. The download flow is:

```
1. Construct URL: https://github.com/<owner>/<repo>/releases/download/<tag>/<asset>
2. GitHub responds with 302 redirect to a CDN URL
3. Follow redirect and download from CDN
```

The HTTP client must follow redirects (up to 10). Go's default `http.Client` follows redirects automatically.

---

## Constraints

- Never disable TLS verification
- Always respect proxy environment variables
- Always include a User-Agent header
- Retry on transient errors (network, 5xx, 429) — never on 4xx (except 429)
- Progress display is mandatory for interactive terminals
- Download timeouts must be generous (5 min) for slow connections

---

## Cross-References

- [Updater Binary](./20-updater-binary.md) — Main consumer of network operations
- [Self-Update Overview](./02-self-update-overview.md) — Binary-based update flow
- [Checksums & Verification](./15-checksums-verification.md) — Post-download verification
- [Install Scripts](./19-install-scripts.md) — Bash/PowerShell download patterns

---

*Network requirements — v3.2.0 — 2026-04-13*
