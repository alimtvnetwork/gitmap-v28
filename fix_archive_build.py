import re
import os

files_to_check = [
    "gitmap/archive/create.go",
    "gitmap/archive/extract.go",
    "gitmap/archive/list.go",
    "gitmap/archive/source.go"
]

for p in files_to_check:
    if os.path.exists(p):
        with open(p, "r", encoding="utf-8") as f:
            c = f.read()
        
        # very naive remove of fmt import if it's unused
        if 'fmt.' not in c:
            c = re.sub(r'\t"fmt"\n', '', c)
            c = re.sub(r'"fmt"\n', '', c)
        
        with open(p, "w", encoding="utf-8") as f:
            f.write(c)

with open("gitmap/archive/archive_test.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("matchAny(\"test.txt\", nil, true) == false", "!matchAny(\"test.txt\", nil, true).Data")
c = c.replace("if matchAny(\"dir/test.txt\", []string{\"*\"}, false) {", "if matchAny(\"dir/test.txt\", []string{\"*\"}, false).Data {")
c = c.replace("matchAny(\"dir/test.txt\", []string{\".\"}, false) == false", "!matchAny(\"dir/test.txt\", []string{\".\"}, false).Data")
with open("gitmap/archive/archive_test.go", "w", encoding="utf-8") as f:
    f.write(c)

print("Fixed imports and tests")
