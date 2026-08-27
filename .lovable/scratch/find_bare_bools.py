import os
import re

root_dir = r"d:\wp-work\riseup-asia\gitmap"
pattern = re.compile(r"^\s*return\s+(true|false)\s*$")

with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\scratch\bare_bools.txt", "w", encoding="utf-8") as out:
    for dirpath, dirnames, filenames in os.walk(root_dir):
        if "vendor" in dirpath or "node_modules" in dirpath:
            continue
        for file in filenames:
            if file.endswith(".go"):
                filepath = os.path.join(dirpath, file)
                try:
                    with open(filepath, "r", encoding="utf-8") as f:
                        for i, line in enumerate(f):
                            if pattern.search(line):
                                out.write(f"{filepath}:{i+1}: {line.strip()}\n")
                except:
                    pass
