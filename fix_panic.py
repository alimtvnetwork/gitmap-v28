import re
import os

with open("gitmap/cmd/testpaths_test.go", "r", encoding="utf-8") as f:
    c = f.read()

c = c.replace('panic("runtime.Caller failed resolving cmd package dir")', 'panic(apperror.New("resolve cmd package dir", "ERR_TEST_PANIC", nil))')

if '"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"' not in c:
    c = c.replace('import (', 'import (\n\t"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"\n')

with open("gitmap/cmd/testpaths_test.go", "w", encoding="utf-8") as f:
    f.write(c)
print("Fixed panic")
