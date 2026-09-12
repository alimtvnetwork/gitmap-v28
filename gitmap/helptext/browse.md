# browse

Open web URLs in the operating system's default browser or Google Chrome.

## Aliases

`open-url`, `browse-url`, `open-browser`

## Usage

    gitmap browse <url> [--chrome]
    gitmap open-url <url> [--chrome]

## Description

Opens any target URL across platforms:
- **Default browser**: Uses the system default browser handler (`rundll32 url.dll,FileProtocolHandler` on Windows, `open` on macOS, `xdg-open` on Linux).
- **Google Chrome**: When `--chrome` or `-c` is passed, GitMap resolves Google Chrome across Windows, macOS, and Linux, and launches the URL directly in Chrome.
- Automatically prepends `https://` if no protocol is specified.

## Examples

```bash
# Open repository documentation in the default browser
gitmap browse https://github.com/alimtvnetwork/gitmap-v28

# Open local development server in Google Chrome
gitmap browse http://localhost:3000 --chrome

# Short URL with automatic https:// resolution
gitmap open-url google.com
```
