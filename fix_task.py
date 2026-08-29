import os
import re

# taskops.go
with open('gitmap/cmd/taskops.go', 'r', encoding='utf-8') as f:
    t = f.read()

t = t.replace('func findTaskByName(tasks model.TaskFile, name string) model.TaskEntry {', 'func findTaskByName(tasks model.TaskFile, name string) (model.TaskEntry, error) {')
t = t.replace('return t\n\t\t}', 'return t, nil\n\t\t}')
t = t.replace('\tpanic("error")\n\n\treturn model.TaskEntry{}', '\treturn model.TaskEntry{}, apperror.NewSimple("task not found: " + name, "E9023")')
t = t.replace('entry := findTaskByName(tasks, name)', 'entry, err := findTaskByName(tasks, name)\n\tif err != nil {\n\t\treturn err\n\t}')

with open('gitmap/cmd/taskops.go', 'w', encoding='utf-8') as f:
    f.write(t)

# tasksync.go
with open('gitmap/cmd/tasksync.go', 'r', encoding='utf-8') as f:
    ts = f.read()
ts = ts.replace('entry := findTaskByName(tasks, name)', 'entry, err := findTaskByName(tasks, name)\n\tif err != nil {\n\t\treturn err\n\t}')
with open('gitmap/cmd/tasksync.go', 'w', encoding='utf-8') as f:
    f.write(ts)
