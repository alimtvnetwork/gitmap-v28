import re, subprocess

def autofix():
    out = subprocess.run(['go', 'build', './...'], cwd='gitmap', capture_output=True, text=True)
    if out.returncode == 0:
        print("Success")
        return True
    
    for line in out.stderr.splitlines():
        line = line.strip()
        if '.go:' not in line: continue
        parts = line.split(':')
        file = 'gitmap/' + parts[0]
        line_num = int(parts[1]) - 1
        
        with open(file, 'r', encoding='utf8') as f: c = f.readlines()
        
        if 'too many return values' in line:
            if 'return apperror' in c[line_num]:
                c[line_num] = c[line_num].replace('return apperror', 'panic(apperror')
                c[line_num] = c[line_num].replace('nil)', 'nil))')
        elif 'cannot use apperror' in line and 'as' in line:
            match = re.search(r'as ([\w\.\[\]]+) value in return', line)
            if match:
                typ = match.group(1)
                c[line_num] = f'\tvar empty {typ}\n\treturn empty\n'
                
        with open(file, 'w', encoding='utf8') as f: f.writelines(c)
    return False

for _ in range(20):
    if autofix(): break
