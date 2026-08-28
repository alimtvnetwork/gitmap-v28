import re

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/clustercommand.go', 'return apperror.Wrap(err, "Error generating run ref:", nil)', 'return "RUN-ERROR"')
r('gitmap/cmd/clustercommand.go', 'func performPreflight(flags ClusterFlags, selector cluster.TargetSelectorType, effective []cluster.ClusterNode, cmdStr string, runRef string) {', 'func performPreflight(flags ClusterFlags, selector cluster.TargetSelectorType, effective []cluster.ClusterNode, cmdStr string, runRef string) *apperror.AppError {')
r('gitmap/cmd/clustercommand.go', 'if flags.NoPreflight == true {\n\t\treturn\n\t}', 'if flags.NoPreflight == true {\n\t\treturn nil\n\t}')
r('gitmap/cmd/clustercommand.go', 'return apperror.New("Operation aborted.", "E9000", nil)\n\t}', 'return apperror.New("Operation aborted.", "E9000", nil)\n\t}\n\treturn nil')
r('gitmap/cmd/clustercommand.go', 'return apperror.Wrap(err, "Error inserting ClusterRun:", nil)', 'return 0')
r('gitmap/cmd/clustercommand.go', 'performPreflight(flags, selector, effective, cmdStr, runRef)', 'if err := performPreflight(flags, selector, effective, cmdStr, runRef); err != nil { return err }')

r('gitmap/cmd/code.go', 'func upsertCodeEntry(rootPath, name string) {', 'func upsertCodeEntry(rootPath, name string) *apperror.AppError {')
r('gitmap/cmd/code.go', 'func appendCodePathsToDB(rootPath string, extras []string) {', 'func appendCodePathsToDB(rootPath string, extras []string) *apperror.AppError {')
r('gitmap/cmd/code.go', 'return apperror.New("fatal error", "E9000", nil)\n\t}', 'return apperror.New("fatal error", "E9000", nil)\n\t}\n\treturn nil')
r('gitmap/cmd/code.go', 'if len(merged) == len(row.Paths) {\n\t\treturn\n\t}', 'if len(merged) == len(row.Paths) {\n\t\treturn nil\n\t}')
r('gitmap/cmd/code.go', 'return apperror.New("fatal error", "E9000", nil)\n\t\t}\n\t\tout = append(out, abs)', 'return nil\n\t\t}\n\t\tout = append(out, abs)')
r('gitmap/cmd/code.go', 'upsertCodeEntry(rootPath, repoName)', 'if err := upsertCodeEntry(rootPath, repoName); err != nil { return err }')
r('gitmap/cmd/code.go', 'appendCodePathsToDB(rootPath, cf.Add)', 'if err := appendCodePathsToDB(rootPath, cf.Add); err != nil { return err }')

