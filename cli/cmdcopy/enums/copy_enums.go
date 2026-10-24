package enums

import (
	"fmt"
	"strings"
)

// CopyActionType specifies the copy action operation.
type CopyActionType string

const (
	CopyActionCopy  CopyActionType = "copy"
	CopyActionPaste CopyActionType = "paste"
)

// IsValid reports whether the copy action is recognized.
func (a CopyActionType) IsValid() bool {
	switch a {
	case CopyActionCopy, CopyActionPaste:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a CopyActionType) String() string {
	return string(a)
}

// ParseCopyActionType parses a raw string into a CopyActionType.
func ParseCopyActionType(s string) (CopyActionType, error) {
	val := CopyActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid copy action: %q", s)
}

// AllCopyActionTypes returns all supported CopyActionType values.
func AllCopyActionTypes() []CopyActionType {
	return []CopyActionType{
		CopyActionCopy,
		CopyActionPaste,
	}
}

// IsCopy reports whether the action is copy.
func (a CopyActionType) IsCopy() bool {
	return a == CopyActionCopy
}

// IsPaste reports whether the action is paste.
func (a CopyActionType) IsPaste() bool {
	return a == CopyActionPaste
}

// HasAction reports whether the action is valid.
func (a CopyActionType) HasAction() bool {
	return a.IsValid()
}

// CopyTargetType specifies the destination target for copied content.
type CopyTargetType string

const (
	CopyTargetClipboard CopyTargetType = "clipboard"
	CopyTargetMemory    CopyTargetType = "memory"
	CopyTargetStdout    CopyTargetType = "stdout"
	CopyTargetFile      CopyTargetType = "file"
)

// IsValid reports whether the copy target is recognized.
func (t CopyTargetType) IsValid() bool {
	switch t {
	case CopyTargetClipboard, CopyTargetMemory, CopyTargetStdout, CopyTargetFile:
		return true
	default:
		return false
	}
}

// String returns the string representation of the target.
func (t CopyTargetType) String() string {
	return string(t)
}

// ParseCopyTargetType parses a raw string into a CopyTargetType.
func ParseCopyTargetType(s string) (CopyTargetType, error) {
	val := CopyTargetType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid copy target: %q", s)
}

// AllCopyTargetTypes returns all supported CopyTargetType values.
func AllCopyTargetTypes() []CopyTargetType {
	return []CopyTargetType{
		CopyTargetClipboard,
		CopyTargetMemory,
		CopyTargetStdout,
		CopyTargetFile,
	}
}

// IsClipboard reports whether the target is the clipboard.
func (t CopyTargetType) IsClipboard() bool {
	return t == CopyTargetClipboard
}

// IsMemory reports whether the target is memory.
func (t CopyTargetType) IsMemory() bool {
	return t == CopyTargetMemory
}

// IsStdout reports whether the target is standard output.
func (t CopyTargetType) IsStdout() bool {
	return t == CopyTargetStdout
}

// IsFile reports whether the target is a file.
func (t CopyTargetType) IsFile() bool {
	return t == CopyTargetFile
}

// HasTarget reports whether the target is valid.
func (t CopyTargetType) HasTarget() bool {
	return t.IsValid()
}

// CopyModeType specifies the execution mode for copy operations.
type CopyModeType string

const (
	CopyModeText CopyModeType = "text"
	CopyModeFile CopyModeType = "file"
	CopyModePipe CopyModeType = "pipe"
)

// IsValid reports whether the copy mode is recognized.
func (m CopyModeType) IsValid() bool {
	switch m {
	case CopyModeText, CopyModeFile, CopyModePipe:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode.
func (m CopyModeType) String() string {
	return string(m)
}

// ParseCopyModeType parses a raw string into a CopyModeType.
func ParseCopyModeType(s string) (CopyModeType, error) {
	val := CopyModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid copy mode: %q", s)
}

// AllCopyModeTypes returns all supported CopyModeType values.
func AllCopyModeTypes() []CopyModeType {
	return []CopyModeType{
		CopyModeText,
		CopyModeFile,
		CopyModePipe,
	}
}

// IsText reports whether the mode is text.
func (m CopyModeType) IsText() bool {
	return m == CopyModeText
}

// IsFile reports whether the mode is file.
func (m CopyModeType) IsFile() bool {
	return m == CopyModeFile
}

// IsPipe reports whether the mode is piped stdin.
func (m CopyModeType) IsPipe() bool {
	return m == CopyModePipe
}

// HasMode reports whether the mode is valid.
func (m CopyModeType) HasMode() bool {
	return m.IsValid()
}

// CopySourceType specifies the source of content to copy.
type CopySourceType string

const (
	CopySourceArgs      CopySourceType = "args"
	CopySourceFile      CopySourceType = "file"
	CopySourceStdin     CopySourceType = "stdin"
	CopySourceClipboard CopySourceType = "clipboard"
	CopySourceMemory    CopySourceType = "memory"
)

// IsValid reports whether the copy source is recognized.
func (s CopySourceType) IsValid() bool {
	switch s {
	case CopySourceArgs, CopySourceFile, CopySourceStdin, CopySourceClipboard, CopySourceMemory:
		return true
	default:
		return false
	}
}

// String returns the string representation of the source.
func (s CopySourceType) String() string {
	return string(s)
}

// ParseCopySourceType parses a raw string into a CopySourceType.
func ParseCopySourceType(s string) (CopySourceType, error) {
	val := CopySourceType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid copy source: %q", s)
}

// AllCopySourceTypes returns all supported CopySourceType values.
func AllCopySourceTypes() []CopySourceType {
	return []CopySourceType{
		CopySourceArgs,
		CopySourceFile,
		CopySourceStdin,
		CopySourceClipboard,
		CopySourceMemory,
	}
}

// IsArgs reports whether the source is command arguments.
func (s CopySourceType) IsArgs() bool {
	return s == CopySourceArgs
}

// IsFile reports whether the source is a file.
func (s CopySourceType) IsFile() bool {
	return s == CopySourceFile
}

// IsStdin reports whether the source is standard input.
func (s CopySourceType) IsStdin() bool {
	return s == CopySourceStdin
}

// IsClipboard reports whether the source is clipboard.
func (s CopySourceType) IsClipboard() bool {
	return s == CopySourceClipboard
}

// IsMemory reports whether the source is persistent memory.
func (s CopySourceType) IsMemory() bool {
	return s == CopySourceMemory
}

// HasSource reports whether the source is valid.
func (s CopySourceType) HasSource() bool {
	return s.IsValid()
}
