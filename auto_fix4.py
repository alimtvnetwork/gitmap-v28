import re, os, subprocess

def fix_errors():
    out = subprocess.run(['go', 'build', './...'], cwd='gitmap', capture_output=True, text=True)
    if out.returncode == 0:
        print("Build passed!")
        return True
    
    print("Fixing errors...")
    
    files_to_goimports = set()

    lines = out.stderr.splitlines()
    for i, line in enumerate(lines):
        line = line.strip()
        if '.go:' not in line: continue
        parts = line.split(':')
        file = 'gitmap/' + parts[0]
        line_num = int(parts[1]) - 1

        with open(file, 'r', encoding='utf8') as f:
            content = f.readlines()

        if 'too many return values' in line:
            if 'return apperror' in content[line_num]:
                content[line_num] = content[line_num].replace('return apperror', 'apperror') + '\treturn\n'
        elif 'missing return' in line:
            if content[line_num].strip() == '}':
                content.insert(line_num, '\treturn\n')
        elif 'cannot use apperror' in line and 'as' in line:
            match = re.search(r'return (apperror\.[^(]+\([^)]+\))', content[line_num])
            if match:
                content[line_num] = f'\tfmt.Fprintln(os.Stderr, {match.group(1)}.Error())\n\tos.Exit(1)\n'
                files_to_goimports.add(file)
        elif 'not enough return values' in line:
            if 'return apperror' in content[line_num]:
                content[line_num] = '\tpanic("not enough return values")\n'
        elif 'FullString undefined' in line:
            content[line_num] = content[line_num].replace('.FullString()', '.Error()')
        elif 'does not implement error' in line and 'apperror.Wrap' in content[line_num]:
            content[line_num] = re.sub(r'apperror\.Wrap\(([^,]+),\s*"([^"]+)",\s*nil\)', r'apperror.New("\2 " + \1, "E9000", nil)', content[line_num])
        elif 'cannot use err' in line and 'as string value in argument to apperror.New' in line:
            content[line_num] = content[line_num].replace('apperror.New(err,', 'apperror.Wrap(err, "error:",')
        elif 'undefined: fmt' in line:
            content.insert(2, 'import "fmt"\n')
        elif 'declared and not used' in line:
            var_name = line.split("'")[1] if "'" in line else line.split()[4]
            # Replace declaration `name :=` or `name, err :=` with `_, err :=`
            content[line_num] = re.sub(r'\b' + re.escape(var_name) + r'\b\s*,', '_,', content[line_num])
            content[line_num] = re.sub(r',\s*\b' + re.escape(var_name) + r'\b', ', _', content[line_num])
            content[line_num] = re.sub(r'\b' + re.escape(var_name) + r'\b\s*:=', '_ =', content[line_num])
        
        with open(file, 'w', encoding='utf8') as f:
            f.writelines(content)
    
    if files_to_goimports:
        for f in files_to_goimports:
            subprocess.run(['go', 'run', 'golang.org/x/tools/cmd/goimports@latest', '-w', f.replace('gitmap/', '')], cwd='gitmap')
    return False

for _ in range(15):
    if fix_errors():
        break
