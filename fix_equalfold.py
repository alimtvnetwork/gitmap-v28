import os, re

gitmap_dir = r'd:\work\gitmap\gitmap'

count = 0
for root, _, files in os.walk(gitmap_dir):
    for file in files:
        if file.endswith('.go'):
            path = os.path.join(root, file)
            with open(path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            def repl_eq(m):
                return f'strings.EqualFold({m.group(1)}, {m.group(2)})'
            def repl_neq(m):
                return f'!strings.EqualFold({m.group(1)}, {m.group(2)})'
            
            new_content = re.sub(r'strings\.ToLower\(([^)]+)\)\s*==\s*("[^"]+")', repl_eq, content)
            new_content = re.sub(r'strings\.ToLower\(([^)]+)\)\s*!=\s*("[^"]+")', repl_neq, new_content)
            
            new_content = re.sub(r'("[^"]+")\s*==\s*strings\.ToLower\(([^)]+)\)', lambda m: f'strings.EqualFold({m.group(2)}, {m.group(1)})', new_content)
            new_content = re.sub(r'("[^"]+")\s*!=\s*strings\.ToLower\(([^)]+)\)', lambda m: f'!strings.EqualFold({m.group(2)}, {m.group(1)})', new_content)
            
            if new_content != content:
                with open(path, 'w', encoding='utf-8', newline='') as f:
                    f.write(new_content)
                count += 1
                print(f'Fixed {path}')
print(f'Fixed {count} files')
