import re

with open('cli/cmd/rootcore.go', 'r') as f:
    content = f.read()

content = content.replace(
    'constants.CmdCreateRepo, constants.CmdCreateRepoAlias,',
    'constants.CmdCreateRepo, constants.CmdCreateRepoAlias, "cr",'
)

with open('cli/cmd/rootcore.go', 'w') as f:
    f.write(content)
