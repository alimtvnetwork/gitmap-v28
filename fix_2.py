import re

with open('gitmap/clonefrom/jsonschema_test.go', 'r', encoding='utf-8') as f:
    c = f.read()
c = c.replace('isNotNumber := !isNumber', 'isInvalidNumber := !isNumber')
c = c.replace('if isNotNumber == true', 'if isInvalidNumber')
with open('gitmap/clonefrom/jsonschema_test.go', 'w', encoding='utf-8') as f:
    f.write(c)

with open('gitmap/cluster/distribution.go', 'r', encoding='utf-8') as f:
    c = f.read()
c = c.replace('hasNoClients := len(clients) == EmptySize', 'isEmpty := len(clients) == EmptySize')
c = c.replace('if hasNoClients == true', 'if isEmpty')
with open('gitmap/cluster/distribution.go', 'w', encoding='utf-8') as f:
    f.write(c)

print("Applied 2 manual fixes.")
