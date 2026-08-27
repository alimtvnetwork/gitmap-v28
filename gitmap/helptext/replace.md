# replace

Replaces text in repository files and tracks changes in the SplitDB for undo history.

## Why use this instead of shell commands?
Instead of running a dangerous multi-file sed command:
`find . -type f -exec sed -i 's/func run(/func execute(/g' {} +`

You should use GitMap's transactional replace:
`gitmap replace "func run(" "func execute("`

## Examples
```bash
gitmap replace "func run(" "func execute("
```
