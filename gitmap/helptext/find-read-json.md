# find-read-json

Same as find-read, but outputs the file paths and contents as a JSON payload for machine piping or LLM usage.

## Examples
```bash
gitmap find-read-json "constants.go"
gitmap find-read-json "*.md" -l 2
```
