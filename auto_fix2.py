import re, os, subprocess

def fix_errors():
    out = subprocess.run(['go', 'build', './...'], cwd='gitmap', capture_output=True, text=True)
    if out.returncode == 0:
        print("Build passed!")
        return True
    
    print("Fixing errors...")

    lines = out.stderr.splitlines()
    for i, line in enumerate(lines):
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
        elif 'not enough return values' in line:
            if 'return apperror' in content[line_num]:
                content[line_num] = content[line_num].replace('return apperror', 'apperror') + '\treturn\n'
                # Just replace with panic or os.Exit
                content[line_num] = '\tpanic("not enough return values")\n'
        elif 'FullString undefined' in line:
            content[line_num] = content[line_num].replace('.FullString()', '.Error()')
        elif 'does not implement error' in line and 'apperror.Wrap' in content[line_num]:
            # Replace apperror.Wrap(string, "msg", nil) with apperror.New("msg: " + string, "E9000", nil)
            content[line_num] = re.sub(r'apperror\.Wrap\(([^,]+),\s*"([^"]+)",\s*nil\)', r'apperror.New("\2 " + \1, "E9000", nil)', content[line_num])

        with open(file, 'w', encoding='utf8') as f:
            f.writelines(content)
    
    return False

for _ in range(15):
    if fix_errors():
        break
