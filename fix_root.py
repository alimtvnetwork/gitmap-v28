import re

with open("gitmap/cmd/root.go", "r", encoding="utf-8") as f:
    content = f.read()

# Replace: if command == "agy" {
content = content.replace('if command == "agy" {', 'if command == "agy" || command == "ag" || command == "antigravity" {')

with open("gitmap/cmd/root.go", "w", encoding="utf-8") as f:
    f.write(content)

print("done")
