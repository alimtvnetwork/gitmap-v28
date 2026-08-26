with open("gitmap/helptext/clone-sync.md", "r", encoding="utf-8") as f:
    content = f.read()

if "## Examples" not in content:
    content += """\n## Examples\n\n```bash\ngitmap clone-sync https://github.com/example/repo\n```\n"""
    with open("gitmap/helptext/clone-sync.md", "w", encoding="utf-8") as f:
        f.write(content)
