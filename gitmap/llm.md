

## AI File Search Patterns

When searching codebases, LLMs can use native gitmap commands OR standard terminal tools.
Here are equivalent alternative command samples for LLM search operations:

- **Find a specific struct definition**:
  - gitmap file-search . "type SearchResult struct"
  - Get-ChildItem -Path gitmap -Recurse -File | Select-String "type SearchResult struct"

- **Search functions with Regex context**:
  - gitmap file-search cmd/ "func dispatch[A-Z]" 0 10
  - Get-ChildItem -Path gitmap/cmd -Filter *.go | Select-String "func dispatch[A-Z]"

- **Find specific function contexts**:
  - gitmap file-search cmd/root.go "func finishCommandAudit" 0 10
  - cat gitmap/cmd/root.go | Select-String "func finishCommandAudit" -Context 0,10
