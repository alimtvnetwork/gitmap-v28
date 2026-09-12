# copy

Copy text, items, or file contents to the OS system clipboard and the persistent repository memory buffer (`.gitmap/memory/clipboard.txt`).

## Aliases

`cp-mem`, `copy-to-memory`

## Usage

    gitmap copy <text...>
    gitmap copy --file <path>
    gitmap copy <filepath>

## Description

The `copy` command provides unified cross-platform copying:
1. Writes content directly to the operating system clipboard (supported across Windows, macOS, and Linux).
2. Persists a copy in `.gitmap/memory/clipboard.txt` so content remains retrievable even in headless environments, CI runners, or across terminal sessions.

## Examples

```bash
# Copy a string of text
gitmap copy "npm run build && npm run test"

# Copy a configuration file into memory & clipboard
gitmap copy --file .env.example

# Copy by direct file path
gitmap copy package.json
```
