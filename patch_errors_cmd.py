import re

with open('cli/cmderrors/errors_cmd.go', 'r') as f:
    content = f.read()

content = content.replace('db.ClearErrors()', 'db.FederatedClearErrors()')
content = content.replace('db.GetError(id)', 'db.FederatedGetError(id)')
content = content.replace('db.ListErrors(limit, false)', 'db.FederatedListErrors(limit, false)')

with open('cli/cmderrors/errors_cmd.go', 'w') as f:
    f.write(content)
