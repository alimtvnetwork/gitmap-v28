import os
import re

for root, _, files in os.walk("gitmap/cmd"):
    for file in files:
        if file.endswith(".go"):
            path = os.path.join(root, file)
            with open(path, "r", encoding="utf-8") as f:
                c = f.read()
            
            new_c = re.sub(r'(\w+)\s*==\s*false\(', r'!\1(', c)
            
            if c != new_c:
                print("Fixing:", path)
                with open(path, "w", encoding="utf-8") as f:
                    f.write(new_c)
