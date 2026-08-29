import os

def fix_file(path, old, new):
    with open(path, 'r', encoding='utf-8') as f:
        content = f.read()
    if old in content:
        with open(path, 'w', encoding='utf-8') as f:
            f.write(content.replace(old, new))
        print('Fixed ' + path)

fix_file('gitmap/cmd/taskops.go', 'panic("error")\n\n\treturn model.TaskEntry{}', 'return model.TaskEntry{}')
fix_file('gitmap/cmd/update.go', 'panic("error")\n\treturn ""', 'return ""')
fix_file('gitmap/cmd/reinstall.go', 'panic("error")\n\treturn constants.ReinstallModeSelf, true', 'return constants.ReinstallModeSelf, false')
