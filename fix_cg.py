import os

p = 'gitmap/cmd/cg.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = c.replace('func parseCgFlags(args []string) cgOptions {', 'func parseCgFlags(args []string) (cgOptions, *apperror.AppError) {')
c = c.replace('return apperror.New("fatal error", "E9000", nil)', 'return cgOptions{}, apperror.New("fatal error", "E9000", nil)')
c = c.replace('return cgOptions{\n\t\tdryRun: *dryRun,\n\t\tpush:   *push,\n\t\tall:    *all,\n\t}\n}', 'return cgOptions{\n\t\tdryRun: *dryRun,\n\t\tpush:   *push,\n\t\tall:    *all,\n\t}, nil\n}')
c = c.replace('opts := parseCgFlags(args)', 'opts, err := parseCgFlags(args)\n\tif err != nil { return err }')
with open(p, 'w', encoding='utf8') as f: f.write(c)
