import os

def find_files(exts):
    res = []
    for root, _, files in os.walk('.'):
        if 'node_modules' in root or '.git' in root or 'dist' in root or 'vendor' in root:
            continue
        for f in files:
            if any(f.endswith(ext) for ext in exts):
                res.append(os.path.join(root, f))
    return res

go_files = find_files(['.go'])
ts_files = find_files(['.ts', '.tsx'])
all_files = go_files + ts_files

results = []

for filepath in all_files:
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            lines = f.readlines()
            
            # Check return without blank line
            for i in range(1, len(lines)):
                curr = lines[i].strip()
                prev = lines[i-1].strip()
                if curr.startswith('return ') and prev != '' and not prev.startswith('{') and not prev.startswith('//'):
                    # To reduce noise, only log a few
                    results.append((filepath, i+1, "return without blank line", curr))
                    
            # Check monolithic functions (> 15 lines)
            in_func = False
            func_start = 0
            brace_count = 0
            for i, line in enumerate(lines):
                if line.strip().startswith('func ') or line.strip().startswith('function ') or '=>' in line:
                    if '{' in line and not in_func:
                        in_func = True
                        func_start = i
                        brace_count = line.count('{') - line.count('}')
                elif in_func:
                    brace_count += line.count('{') - line.count('}')
                    if brace_count <= 0:
                        in_func = False
                        func_len = i - func_start
                        if func_len > 15:
                            results.append((filepath, func_start+1, "monolithic function", f"length {func_len}"))
                            
    except Exception:
        pass

with open('audit_results_extra.csv', 'w', encoding='utf-8') as f:
    for r in results:
        f.write(f"{r[0]}|{r[1]}|{r[2]}|{r[3]}\n")

print(f"Found {len(results)} extra issues.")
