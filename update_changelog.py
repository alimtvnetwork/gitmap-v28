import sys

with open("changelog.md", "r", encoding="utf-8") as f:
    content = f.read()

entry = "## [6.104.0] - 2026-08-25\n\n### Added\n- SSH login and join features plan\n\n"

# replace first occurrence only
content = content.replace("## [", entry + "## [", 1)

with open("changelog.md", "w", encoding="utf-8") as f:
    f.write(content)
