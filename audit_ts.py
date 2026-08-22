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

ts_files = find_files(['.ts', '.tsx'])

def scan_ts_file(filepath):
    results = []
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            lines = f.readlines()
            
            for i, line in enumerate(lines):
                ln = i + 1
                
                # swallowed errors
                if re.search(r'catch\s*\([^)]*\)\s*\{\s*\}', line) or re.search(r'catch\s*\{\s*\}', line):
                    results.append((filepath, ln, "swallowed error (catch {})", line.strip()))
                
                # inverted booleans
                if re.search(r'\b(isNot|hasNo|cannot|shouldNot)\w+', line):
                    results.append((filepath, ln, "inverted boolean", line.strip()))
                    
                # Nested ifs (quick hack: check indentation and if)
                if line.strip().startswith('if ') or line.strip().startswith('if('):
                    # Check next few lines for nested ifs (simplified)
                    # We'll skip complex logic for this python script and just do simple text search
                    pass
                    
    except Exception as e:
        pass
    return results

all_ts_results = []
for f in ts_files:
    all_ts_results.extend(scan_ts_file(f))

import collections
counter = collections.Counter()
for r in all_ts_results:
    counter[r[2]] += 1

print(f"Found {len(all_ts_results)} issues in TS/TSX files.")
for k, v in counter.items():
    print(f"{k}: {v}")

with open('audit_results_ts.csv', 'w', encoding='utf-8') as f:
    for r in all_ts_results:
        f.write(f"{r[0]}|{r[1]}|{r[2]}|{r[3]}\n")
