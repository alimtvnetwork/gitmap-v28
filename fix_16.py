import re

def rewrite(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

rewrite('gitmap/cmd/file_search.go', 'apperror.New("Usage: gitmap file-search <file> <regex> [contextBefore] [contextAfter]", "E9000", nil)', 'apperror.New("Usage: gitmap file-search <file> <regex> [contextBefore] [contextAfter]\\n\\nAlternative LLM Examples:\\n  Get-ChildItem -Path . -Recurse -File | Select-String \\"type SearchResult struct\\"\\n  Get-ChildItem -Path cmd -Filter *.go | Select-String \\"func dispatch[A-Z]\\"\\n  cat cmd/root.go | Select-String \\"func finishCommandAudit\\" -Context 0,10", "E9000", nil)')

rewrite('gitmap/cmd/search_entry.go', 'fmt.Println("Usage: gitmap search <query> [--limit <n>]")', 'fmt.Println("Usage: gitmap search <query> [--limit <n>]\\n\\nAlternative LLM Examples:\\n  Get-ChildItem -Path . -Recurse -File | Select-String \\"type SearchResult struct\\"\\n  Get-ChildItem -Path cmd -Filter *.go | Select-String \\"func dispatch[A-Z]\\"\\n  cat cmd/root.go | Select-String \\"func finishCommandAudit\\" -Context 0,10")')
