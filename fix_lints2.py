import os
import re

with open('gitmap/cmd/changelog.go', 'r', encoding='utf-8') as f:
    t = f.read()

t = t.replace('\treturn dispatchChangelogOutput(version, latest, limit, source, pretty)\n\treturn nil', '\treturn dispatchChangelogOutput(version, latest, limit, source, pretty)')
t = re.sub(r'(\tif !latest && len\(version\) == 0 \{\n\t\treturn nil\n\t\}\n\treturn nil)+', r'\1', t)

with open('gitmap/cmd/changelog.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/cmd/cluster.go', 'r', encoding='utf-8') as f:
    t = f.read()

t = t.replace('\treturn nil\n\treturn nil', '\treturn nil')

with open('gitmap/cmd/cluster.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/cmd/workflow_open_pr.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('"GET"', 'http.MethodGet')
with open('gitmap/cmd/workflow_open_pr.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/cmd/clone.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('out, ok := "", false', 'ok := false')
with open('gitmap/cmd/clone.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/cmd/installer_export_git.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('\ttarget := ""\n', '')
with open('gitmap/cmd/installer_export_git.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/cmd/safety_snapshot.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('\tsrc := ""\n', '')
with open('gitmap/cmd/safety_snapshot.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/store/scheduler.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('func parseScheduleRows(rows *sql.Rows) ([]SchedulerTask, error)', 'func parseScheduleRows(rows *sql.Rows) []SchedulerTask')
t = t.replace('return tasks, nil', 'return tasks')
with open('gitmap/store/scheduler.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/store/scheduler_fetch.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('tasks, err := parseScheduleRows(rows)\n\tif err != nil {\n\t\treturn nil, err\n\t}', 'tasks := parseScheduleRows(rows)')
with open('gitmap/store/scheduler_fetch.go', 'w', encoding='utf-8') as f:
    f.write(t)

with open('gitmap/installer/execute_order.go', 'r', encoding='utf-8') as f:
    t = f.read()
t = t.replace('s *model.InstallerScript, osTarget, order string', 'osTarget, order string')
t = t.replace('err := m.dispatchOrder(ctx, s, osTarget, order)', 'err := m.dispatchOrder(ctx, osTarget, order)')
with open('gitmap/installer/execute_order.go', 'w', encoding='utf-8') as f:
    f.write(t)
