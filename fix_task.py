import re
with open('.lovable/plans/subtasks/01-coding-guideline-fixes/01-inverted-booleans.md', 'r') as f:
    lines = f.readlines()

new_lines = []
for l in lines:
    if 'jsonschema_test.go:172' in l or 'cluster\distribution.go:28' in l or 'cluster\exec_proj.go:62' in l or 'cluster\node_resolver.go:41' in l or 'cluster\ui.go:56' in l:
        continue
    new_lines.append(l)
    
with open('.lovable/plans/subtasks/01-coding-guideline-fixes/01-inverted-booleans.md', 'w') as f:
    f.writelines(new_lines)
