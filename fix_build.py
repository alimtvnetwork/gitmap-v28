import re

with open("gitmap/cmd/auditlegacy.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isMatch == false(f, exclude)", "!isMatch(f, exclude)")
with open("gitmap/cmd/auditlegacy.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/cdops.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("if hasNoGroup {", "if isEmptyGroup {")
with open("gitmap/cmd/cdops.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/chrome_bookmarks.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("if isEmptyRootsAfterRootName {", "if hasNoRootsAfterRootName {")
c = c.replace("if isEmptyRootsAfterFolderPath {", "if hasNoRootsAfterFolderPath {")
c = c.replace("if isEmptyRootsAfterMatch {", "if hasNoRootsAfterMatch {")
with open("gitmap/cmd/chrome_bookmarks.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/chromeprofile_merge.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isMatch == false(f, exclude)", "!isMatch(f, exclude)")
with open("gitmap/cmd/chromeprofile_merge.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/archive/archive_test.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("matchAny(\"test.txt\", nil, true) == false", "!matchAny(\"test.txt\", nil, true).Data")
c = c.replace("if matchAny(\"dir/test.txt\", []string{\"*\"}, false) {", "if matchAny(\"dir/test.txt\", []string{\"*\"}, false).Data {")
c = c.replace("matchAny(\"dir/test.txt\", []string{\".\"}, false) == false", "!matchAny(\"dir/test.txt\", []string{\".\"}, false).Data")
with open("gitmap/archive/archive_test.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/commitin/e2e/edge_cases_test.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isMatch == false(out)", "!isMatch(out)")
with open("gitmap/cmd/commitin/e2e/edge_cases_test.go", "w", encoding="utf-8") as f:
    f.write(c)

print("Fixed lingering build errors")
