import re

with open('cli/cmd/repo_cmd_dispatch.go', 'r') as f:
    content = f.read()

content = content.replace(
    'case "create", "new", "c", "repo-create", "create-repo", "repoc", "crepo":',
    'case "create", "new", "c", "repo-create", "create-repo", "repoc", "crepo", "cr":'
)

with open('cli/cmd/repo_cmd_dispatch.go', 'w') as f:
    f.write(content)
