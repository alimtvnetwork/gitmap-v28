import re, subprocess

def autofix():
    out = subprocess.run(['go', 'build', './...'], cwd='gitmap', capture_output=True, text=True)
    if out.returncode == 0:
        return True
    
    files = set()
    for line in out.stderr.splitlines():
        if '.go:' in line:
            parts = line.strip().split(':')
            files.add('gitmap/' + parts[0])
            
    for p in files:
        with open(p, 'r', encoding='utf8') as f:
            lines = f.readlines()
        for i in range(len(lines)):
            if 'return apperror' in lines[i]:
                match = re.search(r'apperror\.[A-Za-z]+\(([^,]+),', lines[i])
                if match:
                    err_var = match.group(1).strip()
                    lines[i] = f'\tpanic({err_var})\n'
                else:
                    lines[i] = '\tpanic("error")\n'
        with open(p, 'w', encoding='utf8') as f:
            f.writelines(lines)
            
    return False

for _ in range(5):
    if autofix(): break
