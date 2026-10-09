# paste

Paste content from the system clipboard or repository memory buffer to stdout or a destination file.

## Aliases

`paste-mem`

## Usage

    gitmap paste
    gitmap paste --file <path>
    gitmap paste -o <path>

## Description

The `paste` command retrieves copied content first from the OS system clipboard. If the clipboard is empty or inaccessible, it automatically falls back to `.gitmap/memory/clipboard.txt`.

When `--file` or `-o` is provided, the content is saved directly to the specified destination path (automatically creating parent directories if needed). Without flags, it prints to stdout.

## Examples

```bash
# Print clipboard or memory content to terminal
gitmap paste

# Paste content into a new or existing file
gitmap paste --file configs/secrets.env

# Pipe pasted content to another command
gitmap paste | jq .
```
