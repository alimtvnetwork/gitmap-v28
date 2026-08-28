import re

with open('cmd/roottooling.go', 'r', encoding='utf-8') as f:
    c = f.read()

c = c.replace('func() { checkHelp("desktop-sync", argsTail()); runDesktopSync() }', 'func() error { checkHelp("desktop-sync", argsTail()); return runDesktopSync() }')
c = c.replace('func() { checkHelp("rescan", argsTail()); runRescan() }', 'func() error { checkHelp("rescan", argsTail()); return runRescan() }')
c = c.replace('func() { checkHelp("doctor", argsTail()); runDoctor(argsTail()) }', 'func() error { checkHelp("doctor", argsTail()); return runDoctor(argsTail()) }')
c = c.replace('func() { RunInstallerCLI(argsTail()) }', 'func() error { RunInstallerCLI(argsTail()); return nil }')
c = c.replace('func() { RunOSCLI(argsTail()) }', 'func() error { RunOSCLI(argsTail()); return nil }')

with open('cmd/roottooling.go', 'w', encoding='utf-8') as f:
    f.write(c)

with open('cmd/rootcore.go', 'r', encoding='utf-8') as f:
    c = f.read()

c = re.sub(r'func\(\) \{ run([a-zA-Z0-9_]+)\((.*?)\) \}', r'func() error { return run\1(\2) }', c)
c = re.sub(r'func\(\) \{ dispatch([a-zA-Z0-9_]+)\((.*?)\) \}', r'func() error { dispatch\1(\2); return nil }', c)

with open('cmd/rootcore.go', 'w', encoding='utf-8') as f:
    f.write(c)
