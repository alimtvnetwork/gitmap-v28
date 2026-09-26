import re

with open('cli/cmd/commons.go', 'r') as f:
    content = f.read()

content = content.replace('parseSyncFlags', 'parseCommonFlags')
content = content.replace('runSyncLines', 'runCommonLines')
content = content.replace('runSyncLFSInstall', 'runCommonLFSInstall')
content = content.replace('runSyncPrettierRC', 'runCommonPrettierRC')

with open('cli/cmd/commons.go', 'w') as f:
    f.write(content)

with open('cli/cmd/root.go', 'r') as f:
    content = f.read()

content = re.sub(r'found, err = dispatchSync\(command\).*?return\n\t\}', 'found, err = dispatchCommon(command)\n\tif handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {\n\t\treturn\n\t}', content, flags=re.DOTALL)

with open('cli/cmd/root.go', 'w') as f:
    f.write(content)
