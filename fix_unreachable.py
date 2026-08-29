import os
import subprocess

files = ['cmd/branch.go', 'cmd/cluster.go', 'cmd/dbreset.go', 'cmd/historyreset.go', 'cmd/probe.go', 'cmd/regoldens.go', 'cmd/replaceflags.go', 'cmd/rm.go', 'cmd/selfuninstallhandoff.go', 'cmd/tasksync.go', 'cmd/watch.go']

for f in files:
    path = os.path.join('gitmap', f)
    with open(path, 'r', encoding='utf-8') as file:
        lines = file.readlines()
    
    # Simple fix: if a line is just 'return nil' and the previous non-empty line is also 'return nil' or panic or os.Exit, we comment it out.
    new_lines = []
    for i, line in enumerate(lines):
        if line.strip() in ['return nil', 'return false', 'return true']:
            # check previous line
            j = i - 1
            while j >= 0 and lines[j].strip() == '':
                j -= 1
            if j >= 0:
                prev = lines[j].strip()
                if prev in ['return nil', 'return false', 'return true'] or prev.startswith('panic(') or prev.startswith('os.Exit(') or prev.startswith('return apperror.'):
                    new_lines.append('// ' + line)
                    continue
        new_lines.append(line)
        
    with open(path, 'w', encoding='utf-8') as file:
        file.writelines(new_lines)
