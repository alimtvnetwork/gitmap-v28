import re

with open('cli/cmdagy/agy_misc.go', 'r') as f:
    content = f.read()

content = content.replace(
    'Aliases: []string{"ip"},',
    'Aliases: []string{"imp", "import"},'
)

with open('cli/cmdagy/agy_misc.go', 'w') as f:
    f.write(content)
