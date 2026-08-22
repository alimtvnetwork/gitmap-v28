import re

with open("gitmap/cmd/auditlegacy.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isMatch == false(f, exclude)", "!isMatch(f, exclude)")
with open("gitmap/cmd/auditlegacy.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/cdops.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isEmptyGroup :=", "hasNoGroup :=")
with open("gitmap/cmd/cdops.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/chrome_bookmarks.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("hasNoRootsAfterRootName :=", "isEmptyRootsAfterRootName :=")
c = c.replace("hasNoRootsAfterFolderPath :=", "isEmptyRootsAfterFolderPath :=")
c = c.replace("hasNoRootsAfterMatch :=", "isEmptyRootsAfterMatch :=")
with open("gitmap/cmd/chrome_bookmarks.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cmd/chromeprofile_merge.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isMatch == false(f, exclude)", "!isMatch(f, exclude)")
with open("gitmap/cmd/chromeprofile_merge.go", "w", encoding="utf-8") as f:
    f.write(c)

print("Fixed more boolean errors")
