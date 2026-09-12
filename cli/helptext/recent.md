# recent

Show the last 10 repos visited via the navigation helper. Pipe `--print` to fzf for fuzzy-jump.

## Examples

```bash
gitmap recent
gitmap rct
cd "$(gitmap recent --print | fzf)"
```
