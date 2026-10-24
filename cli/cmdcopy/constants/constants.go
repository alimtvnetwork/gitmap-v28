package constants

import "os"

// Command names and aliases for copy and paste operations.
const (
	CmdCopy           = "copy"
	CmdCopyAliasMem   = "cp-mem"
	CmdCopyAliasToMem = "copy-to-memory"
	CmdPaste          = "paste"
	CmdPasteAliasMem  = "paste-mem"
)

// Flag literals for file and output parameters.
const (
	FlagFile       = "--file"
	FlagFileShort  = "-f"
	FlagFilepath   = "--filepath"
	FlagOut        = "--out"
	FlagOutShort   = "-o"
	FlagOutput     = "--output"
	FlagPrefixFile = "--file="
	FlagPrefixOut  = "--out="
)

// Directory paths, filenames, and default file permissions.
const (
	DirGitmap                 = ".gitmap"
	DirMemory                 = "memory"
	FileClipboard             = "clipboard.txt"
	DefaultFilePerm           = os.FileMode(0644)
	DefaultDirPerm            = os.FileMode(0755)
	DotDirSentinel            = "."
)

// Git CLI executable and subcommands.
const (
	GitBin             = "git"
	GitRevParse        = "rev-parse"
	GitArgShowTopLevel = "--show-toplevel"
)

// ANSI terminal color formatting codes.
const (
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

// Unicode status glyphs.
const (
	GlyphClipboard = "📋"
	GlyphSuccess   = "✅"
	GlyphWarning   = "⚠️ "
)

// Output and status format strings.
const (
	MsgCopiedSuccessFmt    = "  %s📋 Copied %d bytes to memory & clipboard ✅%s\n"
	MsgPastedSuccessFmt    = "  %s📋 Pasted %d bytes to %s ✅%s\n"
	MsgEmptyPasteNoticeFmt = "  %s⚠️  Clipboard and memory buffer are empty.%s\n"
)

// Help and usage documentation strings.
const (
	HelpCopyUsage  = "Usage: gitmap copy <text...> [--file <path>]"
	HelpCopyDesc   = "Copies text, items, or file contents to OS clipboard and persistent repository memory."
	HelpPasteUsage = "Usage: gitmap paste [--file <path>] [--out <path>]"
	HelpPasteDesc  = "Pastes content from OS clipboard or memory buffer to stdout or a file."
)

// Error codes, static messages, and diagnostic error format strings.
const (
	ErrCodeCopyEmpty            = "COPY_EMPTY"
	ErrCodeReadFileFailed       = "READ_FILE_FAILED"
	ErrCodeReadStdinFailed      = "READ_STDIN_FAILED"
	ErrCodeWriteMemoryFailed    = "WRITE_MEMORY_FAILED"
	ErrMsgNothingToCopy         = "nothing to copy: text or file required"
	ErrMsgEmptyTopLevel         = "empty top-level"
	ErrFmtReadFile              = "read file %q for copy"
	ErrFmtReadStdin             = "read stdin for copy"
	ErrFmtWriteMemoryClipboard  = "write memory clipboard buffer"
	ErrFmtCreateMemoryDir       = "create memory directory"
	ErrFmtWritePastedContent    = "write pasted content to %q"
)

// Process exit codes for command termination.
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitUsageError   = 2
)

// AllCopyFileFlags returns all recognized CLI flags designating a file argument.
func AllCopyFileFlags() []string {
	return []string{
		FlagFile,
		FlagFileShort,
		FlagFilepath,
		FlagOut,
		FlagOutShort,
		FlagOutput,
	}
}
