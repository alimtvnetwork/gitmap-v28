import re

with open('cli/cmdagy/agy_most_conv.go', 'r') as f:
    content = f.read()

content = content.replace(
    'Use:   "most-conv",',
    'Use:   "most-conv",\n\tAliases: []string{"most-conversation", "most-conversations"},'
)

with open('cli/cmdagy/agy_most_conv.go', 'w') as f:
    f.write(content)
