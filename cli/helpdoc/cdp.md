# gitmap cdp

Native pure Go Chrome DevTools Protocol (CDP) WebSocket inspection engine and JavaScript evaluator.

## Usage

```bash
gitmap cdp targets                         # List active debugging targets across detected DevTools ports
gitmap cdp targets -p 9222                 # List targets on a specific port
gitmap cdp targets --json                  # Output targets as structured JSON
gitmap cdp probe                           # Probe active page for DOM input and textarea elements
gitmap cdp probe -p 9222 -t <target-id>    # Probe a specific target page
gitmap cdp probe --json                    # Output DOM probe hierarchy as JSON
gitmap cdp eval "<javascript-expression>"  # Evaluate JavaScript expression on active page
gitmap cdp eval "document.title"           # Get current page title
gitmap cdp eval "window.location.href"     # Get active page URL
gitmap cdp eval "<expr>" --json            # Output evaluated value as JSON
```

## Options

| Flag | Shorthand | Description | Default |
|---|---|---|---|
| `--port` | `-p` | Chrome DevTools remote debugging port | Auto-discovered (9222, 49994, 9229, 8315) |
| `--target` | `-t` | Specific target ID or URL substring filter | Active "page" target |
| `--json` | | Output results in JSON format | `false` |
| `--timeout` | | Timeout for WebSocket handshake and RPC evaluation in seconds | `5` |

## Subcommands

- `targets`: Discovers open CDP endpoints via `http://127.0.0.1:<port>/json/list` and displays target title, type, ID, and WebSocket URL.
- `probe`: Connects via pure Go RFC 6455 WebSocket client, executes a non-intrusive selector ladder across `textarea`, `div[contenteditable]`, and `[role="textbox"]`, and returns active element metadata.
- `eval <expr>`: Executes a JavaScript expression in the context of the active target via `Runtime.evaluate` over WebSocket.

## Policy Note for Autonomous AI Agents

AI agents MUST NOT write ad-hoc Python raw socket WebSocket scripts (e.g., manual RFC 6455 masking, `struct.pack`, `send_frame`/`recv_frame`). Always use native `gitmap cdp probe` or `gitmap cdp eval`.
