import re, os, subprocess

def fix_errors():
    out = subprocess.run(['go', 'build', './...'], cwd='gitmap', capture_output=True, text=True)
    if out.returncode == 0:
        print("Build passed!")
        return True

    lines = out.stderr.splitlines()
    for i, line in enumerate(lines):
        if '.go:' not in line: continue
        parts = line.split(':')
        file = 'gitmap/' + parts[0]
        line_num = int(parts[1]) - 1
        msg = parts[3].strip() if len(parts) > 3 else ''

        with open(file, 'r', encoding='utf8') as f:
            content = f.readlines()

        if 'too many return values' in line:
            # Change `return apperror.X()` to `apperror.X(); return`
            if 'return apperror' in content[line_num]:
                content[line_num] = content[line_num].replace('return apperror', 'apperror') + '\treturn\n'
        elif 'missing return' in line:
            # Add `return` or `return nil` before closing brace
            if content[line_num].strip() == '}':
                content.insert(line_num, '\treturn\n')
        elif 'cannot use apperror' in line and 'as' in line:
            # Change `return apperror.X()` to `fmt.Fprintln(os.Stderr, "error"); os.Exit(1)`
            match = re.search(r'return (apperror\.[^(]+\([^)]+\))', content[line_num])
            if match:
                content[line_num] = f'\tfmt.Fprintln(os.Stderr, {match.group(1)}.FullString())\n\tos.Exit(1)\n'

        with open(file, 'w', encoding='utf8') as f:
            f.writelines(content)
    
    return False

for _ in range(5):
    if fix_errors():
        break
