import os

p = 'gitmap/cmd/cdops.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = c.replace('return repos, nil\\n}', 'return repos, nil\n}')
c = c.replace('return p.Path, nil\\n}', 'return p.Path, nil\n}')
with open(p, 'w', encoding='utf8') as f: f.write(c)

p = 'gitmap/cmd/bookmarksave.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = c.replace('if err := if err := checkBookmarkNotExists(db, name); err != nil { return err }; err != nil { return err }', 'if err := checkBookmarkNotExists(db, name); err != nil { return err }')
with open(p, 'w', encoding='utf8') as f: f.write(c)
