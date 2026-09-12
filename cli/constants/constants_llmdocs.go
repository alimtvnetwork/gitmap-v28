package constants

// llm-docs messages.
const (
	MsgLLMDocsWritten = "  ✓ LLM.md written to %s\n"
	MsgLLMDocsGenning = "  ↻ Generating LLM.md from command registry...\n"
	ErrLLMDocsWrite   = "  ✗ Could not write LLM.md: %v\n"
	HelpLLMDocs       = "  llm-docs (ld)       Generate LLM.md reference for AI assistants"
)

// llm-docs flags.
const (
	FlagLLMDocsStdout       = "stdout"
	FlagDescLLMDocsStdout   = "Print to stdout instead of writing LLM.md file"
	FlagLLMDocsFormat       = "format"
	FlagDescLLMDocsFormat   = "Output format: markdown (default) or json"
	ErrLLMDocsFormat        = "  ✗ Unknown format %q — use markdown or json\n"
	FlagLLMDocsSections     = "sections"
	FlagDescLLMDocsSections = "Comma-separated sections to include (default: all)"
	ErrLLMDocsSections      = "  ✗ Unknown section %q — valid: commands,architecture,flags,conventions,structure,database,installation,patterns\n"
	LLMDocsValidSections    = "commands,architecture,flags,conventions,structure,database,installation,patterns"
)
