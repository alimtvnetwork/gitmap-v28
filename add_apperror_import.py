import re

with open('gitmap/archive/create.go', 'r', encoding='utf-8') as f:
    c = f.read()

if '"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"' not in c:
    c = re.sub(r'import \(', 'import (\n\t"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"\n', c, count=1)

with open('gitmap/archive/create.go', 'w', encoding='utf-8') as f:
    f.write(c)

print("Added apperror import")
