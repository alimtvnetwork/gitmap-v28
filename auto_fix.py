import re
import os
import glob

def process_file(filepath):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. isNotX := !isX  => isNonX := isX == false
    def repl_isNot(m):
        name = m.group(1)
        val = m.group(2)
        return f"isNon{name} := {val} == false"
    
    content = re.sub(r'isNot(\w+)\s*:=\s*!(is\w+)', repl_isNot, content)
    content = re.sub(r'isNot(\w+)\s*:=\s*!(\w+\(.*\))', repl_isNot, content)
    
    # 2. if isNotX == true => if isNonX
    def repl_if(m):
        name = m.group(1)
        return f"if isNon{name}"
    content = re.sub(r'if isNot(\w+) == true', repl_if, content)
    content = re.sub(r'if isNot(\w+)\s*\{', r'if isNon\1 {', content)
    
    # 3. hasNoX := len(X) == 0 => isEmptyX := len(X) == 0
    def repl_hasNo(m):
        name = m.group(1)
        val = m.group(2)
        return f"isEmpty{name} := {val}"
    
    content = re.sub(r'hasNo(\w+)\s*:=\s*(len.*)', repl_hasNo, content)
    content = re.sub(r'if hasNo(\w+) == true', r'if isEmpty\1', content)
    content = re.sub(r'if hasNo(\w+)\s*\{', r'if isEmpty\1 {', content)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

for root, _, files in os.walk('gitmap'):
    for file in files:
        if file.endswith('.go'):
            process_file(os.path.join(root, file))

print("Automated replacements complete!")
