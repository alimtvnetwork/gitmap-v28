import re

def rewrite(p):
    with open(p, 'r', encoding='utf8') as f:
        lines = f.readlines()
    for i in range(len(lines)):
        if 'return apperror' in lines[i]:
            match = re.search(r'apperror\.[A-Za-z]+\(([^,]+),', lines[i])
            if match:
                err_var = match.group(1).strip()
                if err_var == '"fatal error"' or err_var.startswith('constants'):
                    lines[i] = f'\tpanic({err_var})\n'
                else:
                    lines[i] = f'\tpanic({err_var})\n'
            else:
                lines[i] = '\tpanic("error")\n'
    with open(p, 'w', encoding='utf8') as f:
        f.writelines(lines)

for p in ['gitmap/cmd/reinstall.go', 'gitmap/cmd/releasealias.go', 'gitmap/cmd/releasepull.go', 'gitmap/cmd/releaserebase.go', 'gitmap/cmd/releaseself.go', 'gitmap/cmd/replaceflags.go']:
    rewrite(p)
