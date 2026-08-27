# find-read-json

Same as find-read, but outputs the file paths and contents as a JSON payload for machine piping or LLM usage.

## Why use this instead of shell commands?
Instead of manually aggregating content and wrapping it in JSON for your prompt payload, just use:
`gitmap find-read-json "*.md" -l 2`

## Examples
```bash
gitmap find-read-json "constants.go"
gitmap find-read-json "*.md" -l 2
```
