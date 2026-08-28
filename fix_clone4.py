import os, re

p = 'gitmap/cmd/clone.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()

c = re.sub(r'(fmt\.Fprintf\(os\.Stderr, "  \\033\[1;92m.+?Task %s spawned and running\.\\n", taskID\)\n\t\})', r'\1\n\treturn nil', c)
c = re.sub(r'(fmt\.Printf\(constants\.MsgCloneComplete, summary\.Succeeded, summary\.Failed\)\n)', r'\1\treturn nil\n', c)

c = c.replace('func resolveShorthand(source string) string {', 'func resolveShorthand(source string) (string, *apperror.AppError) {')
c = c.replace('return source\n}', 'return source, nil\n}')
c = c.replace('source = resolveShorthand(source)', 'source, err = resolveShorthand(source)\n\tif err != nil { return err }')

with open(p, 'w', encoding='utf8') as f: f.write(c)
