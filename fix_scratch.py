import re
import os

files = [
    ".lovable/scratch/find_ignored_errs.go",
    ".lovable/scratch/find_monolithic.go",
    ".lovable/scratch/find_nested_ifs.go",
    ".lovable/scratch/find_single_chars.go"
]

for p in files:
    if os.path.exists(p):
        with open(p, "r", encoding="utf-8") as f:
            c = f.read()
        c = c.replace('panic(err)', 'panic(apperror.New("scratch failure", "ERR_SCRATCH", map[string]any{"err": err}))')
        if '"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"' not in c:
            c = c.replace('import (', 'import (\n\t"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"\n', 1)
        with open(p, "w", encoding="utf-8") as f:
            f.write(c)

print("Fixed scratch panics")
