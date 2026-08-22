import re
import os

for root, _, files in os.walk('gitmap'):
    for file in files:
        if file.endswith('.go'):
            path = os.path.join(root, file)
            with open(path, 'r', encoding='utf-8') as f:
                content = f.read()
            # replace `isX == false(y)` with `isX(y) == false`
            new_content = re.sub(r'(\w+) == false\((.*?)\)', r'\1(\2) == false', content)
            if new_content != content:
                with open(path, 'w', encoding='utf-8') as f:
                    f.write(new_content)

print("Fixed syntax errors!")
