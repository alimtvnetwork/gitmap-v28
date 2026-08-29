# Subtask 21.02: LLM Specification & Public URL Support

## Target Files

- `gitmap/cmd/llm/llm.go`
- `gitmap/llm.md`
- `llm.md`

## Actions

- [ ] Enrich `gitmap/cmd/llm/llm.go` with complete architectural documentation, AI command alternatives, and file search context.
- [ ] Prominently display the public GitHub MD link (`https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md`).
- [ ] Ensure `gitmap llm --url` prints the raw GitHub URL directly for scripting.

## Acceptance Criteria

- [ ] `gitmap llm` outputs full specification + URL.
- [ ] `gitmap llm --url` prints only the URL.
