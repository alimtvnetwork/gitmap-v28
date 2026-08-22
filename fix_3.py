import re

with open('gitmap/cmd/replacewalk_test.go', 'r', encoding='utf-8') as f:
    c = f.read()

# fix `isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release") == false)`
c = re.sub(r'isExcludedPrefix\(([^,]+), ([^=]+) == false\)', r'isExcludedPrefix(\1, \2) == false', c)

with open('gitmap/cmd/replacewalk_test.go', 'w', encoding='utf-8') as f:
    f.write(c)
    
print("Fixed replacewalk_test.go")
