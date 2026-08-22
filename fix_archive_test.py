import re

with open("gitmap/archive/archive_test.go", "r", encoding="utf-8") as f:
    c = f.read()

c = c.replace('if matchAny("test.txt", nil, false) {', 'if matchAny("test.txt", nil, false).Data {')
c = c.replace('if matchAny("dir/test.txt", []string{"*.txt"}, false) == false {', 'if !matchAny("dir/test.txt", []string{"*.txt"}, false).Data {')
c = c.replace('if matchAny("dir/test.txt", []string{"dir/*"}, false) == false {', 'if !matchAny("dir/test.txt", []string{"dir/*"}, false).Data {')

with open("gitmap/archive/archive_test.go", "w", encoding="utf-8") as f:
    f.write(c)

print("Fixed archive_test.go boolean checks")
