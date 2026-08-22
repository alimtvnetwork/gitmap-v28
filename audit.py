import os
import re

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

def scan_file(filepath):
    results = []
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            lines = f.readlines()
            
            for i, line in enumerate(lines):
                ln = i + 1
                # fmt.Errorf
                if 'fmt.Errorf(' in line:
                    results.append((filepath, ln, "fmt.Errorf", line.strip()))
                
                # panic
                if 'panic(' in line:
                    results.append((filepath, ln, "panic", line.strip()))
                    
                # swallowed error
                if '_ = err' in line or '_ , err :=' in line or '_, err =' in line:
                    results.append((filepath, ln, "swallowed error", line.strip()))

                # bare error return in function signature
                # func Foo() error {
                # func Foo() (int, error) {
                if line.strip().startswith('func ') and ' error' in line and '*apperror.AppError' not in line:
                    results.append((filepath, ln, "raw error return", line.strip()))
                    
                # boolean return in Go
                if line.strip().startswith('func ') and (' bool' in line or '(bool' in line) and 'Result' not in line:
                    results.append((filepath, ln, "bare bool return", line.strip()))
                    
    except Exception as e:
        pass
    return results

all_results = []
for f in go_files:
    all_results.extend(scan_file(f))

with open('audit_results_go.csv', 'w', encoding='utf-8') as f:
    for r in all_results:
        f.write(f"{r[0]}|{r[1]}|{r[2]}|{r[3]}\n")

print(f"Found {len(all_results)} issues in Go files.")
