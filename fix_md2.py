with open("gitmap/helptext/clone-sync.md", "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace("    gitmap clone-sync https://github.com/example/repo", "```bash\ngitmap clone-sync https://github.com/example/repo\n```")
content = content.replace("    gitmap cs https://github.com/example/repo1 https://github.com/example/repo2", "```bash\ngitmap cs https://github.com/example/repo1 https://github.com/example/repo2\n```")

with open("gitmap/helptext/clone-sync.md", "w", encoding="utf-8") as f:
    f.write(content)

print("done")
