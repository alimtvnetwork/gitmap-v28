import os
def fix_file(f):
    with open(f, 'r', encoding='utf-8') as file:
        c = file.read()
    c = c.replace('"constants.ErrBookmarkQuery+\\n"', 'constants.ErrBookmarkQuery')
    c = c.replace('"constants.ErrBookmarkQuery+\n"', 'constants.ErrBookmarkQuery')
    with open(f, 'w', encoding='utf-8') as file:
        file.write(c)

fix_file('gitmap/cmd/bookmarklist.go')
fix_file('gitmap/cmd/bookmarkrun.go')
