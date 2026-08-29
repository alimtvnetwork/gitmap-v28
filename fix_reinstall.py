import os

with open('gitmap/cmd/reinstall.go', 'r', encoding='utf-8') as f:
    r = f.read()

r = r.replace('return "", false', 'return apperror.NewSimple("reinstall aborted or failed", "E9025")')

with open('gitmap/cmd/reinstall.go', 'w', encoding='utf-8') as f:
    f.write(r)
