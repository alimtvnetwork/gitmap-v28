import os

with open('gitmap/clonepick/clonepick.go', 'r', encoding='utf-8') as f:
    t = f.read()

t = t.replace('StatusCancelled = "canceled"', 'StatusCancelled = "canceled"')

with open('gitmap/clonepick/clonepick.go', 'w', encoding='utf-8') as f:
    f.write(t)
