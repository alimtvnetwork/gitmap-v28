import re

with open('cli/cmd/repo_create_params.go', 'r') as f:
    content = f.read()

content = content.replace(
    'IsCG:         hasArgFlag(args, "--cg"),',
    'IsCG:         hasArgFlag(args, "--cg") || hasArgFlag(args, "--cd") || hasArgFlag(args, "--coding-guideline"),'
)

with open('cli/cmd/repo_create_params.go', 'w') as f:
    f.write(content)
