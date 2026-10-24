package enums

import (
	"fmt"
	"strings"
)

// DocsCategoryType specifies documentation content categories.
type DocsCategoryType string

const (
	DocsCategoryCLI          DocsCategoryType = "cli"
	DocsCategoryArchitecture DocsCategoryType = "architecture"
	DocsCategorySpec         DocsCategoryType = "spec"
	DocsCategoryPrompt       DocsCategoryType = "prompt"
	DocsCategoryGuide        DocsCategoryType = "guide"
	DocsCategoryAPI          DocsCategoryType = "api"
	DocsCategoryAll          DocsCategoryType = "all"
)

// IsValid reports whether the category is a recognized variant.
func (c DocsCategoryType) IsValid() bool {
	switch c {
	case DocsCategoryCLI, DocsCategoryArchitecture, DocsCategorySpec,
		DocsCategoryPrompt, DocsCategoryGuide, DocsCategoryAPI, DocsCategoryAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the category.
func (c DocsCategoryType) String() string {
	return string(c)
}

// IsCLI reports whether the category is CLI documentation.
func (c DocsCategoryType) IsCLI() bool {
	return c == DocsCategoryCLI
}

// IsArchitecture reports whether the category is architecture documentation.
func (c DocsCategoryType) IsArchitecture() bool {
	return c == DocsCategoryArchitecture
}

// IsSpec reports whether the category is specification documentation.
func (c DocsCategoryType) IsSpec() bool {
	return c == DocsCategorySpec
}

// IsPrompt reports whether the category is prompt engineering documentation.
func (c DocsCategoryType) IsPrompt() bool {
	return c == DocsCategoryPrompt
}

// IsGuide reports whether the category is user guide documentation.
func (c DocsCategoryType) IsGuide() bool {
	return c == DocsCategoryGuide
}

// IsAPI reports whether the category is API documentation.
func (c DocsCategoryType) IsAPI() bool {
	return c == DocsCategoryAPI
}

// IsAll reports whether all categories are selected.
func (c DocsCategoryType) IsAll() bool {
	return c == DocsCategoryAll
}

// HasCategoryFilter reports whether a specific category constraint is applied.
func (c DocsCategoryType) HasCategoryFilter() bool {
	return c != DocsCategoryAll && c.IsValid()
}

// ParseDocsCategoryType parses a raw string into a DocsCategoryType.
func ParseDocsCategoryType(s string) (DocsCategoryType, error) {
	val := DocsCategoryType(strings.ToLower(strings.TrimSpace(s)))

	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid docs category: %q", s)
}

// ParseDocsCategory is an alias for ParseDocsCategoryType.
func ParseDocsCategory(s string) (DocsCategoryType, error) {
	return ParseDocsCategoryType(s)
}

// AllDocsCategoryTypes returns all supported documentation categories.
func AllDocsCategoryTypes() []DocsCategoryType {
	return []DocsCategoryType{
		DocsCategoryCLI,
		DocsCategoryArchitecture,
		DocsCategorySpec,
		DocsCategoryPrompt,
		DocsCategoryGuide,
		DocsCategoryAPI,
		DocsCategoryAll,
	}
}

// AllDocsCategories is an alias for AllDocsCategoryTypes.
func AllDocsCategories() []DocsCategoryType {
	return AllDocsCategoryTypes()
}

// DocsViewerModeType specifies the viewing mode or target for documentation.
type DocsViewerModeType string

const (
	DocsViewerModeBrowser  DocsViewerModeType = "browser"
	DocsViewerModeTerminal DocsViewerModeType = "terminal"
	DocsViewerModeServer   DocsViewerModeType = "server"
	DocsViewerModeWeb      DocsViewerModeType = "web"
)

// IsValid reports whether the viewer mode is a recognized variant.
func (m DocsViewerModeType) IsValid() bool {
	switch m {
	case DocsViewerModeBrowser, DocsViewerModeTerminal,
		DocsViewerModeServer, DocsViewerModeWeb:
		return true
	default:
		return false
	}
}

// String returns the string representation of the viewer mode.
func (m DocsViewerModeType) String() string {
	return string(m)
}

// IsBrowser reports whether the viewer target is a web browser.
func (m DocsViewerModeType) IsBrowser() bool {
	return m == DocsViewerModeBrowser
}

// IsTerminal reports whether the viewer target is terminal display.
func (m DocsViewerModeType) IsTerminal() bool {
	return m == DocsViewerModeTerminal
}

// IsServer reports whether the viewer target is a local HTTP server.
func (m DocsViewerModeType) IsServer() bool {
	return m == DocsViewerModeServer
}

// IsWeb reports whether the viewer target is hosted web documentation.
func (m DocsViewerModeType) IsWeb() bool {
	return m == DocsViewerModeWeb
}

// HasServerTarget reports whether the viewer mode utilizes an HTTP server.
func (m DocsViewerModeType) HasServerTarget() bool {
	return (m == DocsViewerModeServer || m == DocsViewerModeWeb) && m.IsValid()
}

// ParseDocsViewerModeType parses a raw string into a DocsViewerModeType.
func ParseDocsViewerModeType(s string) (DocsViewerModeType, error) {
	val := DocsViewerModeType(strings.ToLower(strings.TrimSpace(s)))

	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid docs viewer mode: %q", s)
}

// ParseDocsViewerMode is an alias for ParseDocsViewerModeType.
func ParseDocsViewerMode(s string) (DocsViewerModeType, error) {
	return ParseDocsViewerModeType(s)
}

// AllDocsViewerModeTypes returns all supported documentation viewer modes.
func AllDocsViewerModeTypes() []DocsViewerModeType {
	return []DocsViewerModeType{
		DocsViewerModeBrowser,
		DocsViewerModeTerminal,
		DocsViewerModeServer,
		DocsViewerModeWeb,
	}
}

// AllDocsViewerModes is an alias for AllDocsViewerModeTypes.
func AllDocsViewerModes() []DocsViewerModeType {
	return AllDocsViewerModeTypes()
}

// DocsFormatType specifies presentation formats for documentation.
type DocsFormatType string

const (
	DocsFormatMarkdown DocsFormatType = "markdown"
	DocsFormatHTML     DocsFormatType = "html"
	DocsFormatText     DocsFormatType = "text"
	DocsFormatJSON     DocsFormatType = "json"
)

// IsValid reports whether the format is a recognized variant.
func (f DocsFormatType) IsValid() bool {
	switch f {
	case DocsFormatMarkdown, DocsFormatHTML,
		DocsFormatText, DocsFormatJSON:
		return true
	default:
		return false
	}
}

// String returns the string representation of the format.
func (f DocsFormatType) String() string {
	return string(f)
}

// IsMarkdown reports whether the format is Markdown.
func (f DocsFormatType) IsMarkdown() bool {
	return f == DocsFormatMarkdown
}

// IsHTML reports whether the format is HTML.
func (f DocsFormatType) IsHTML() bool {
	return f == DocsFormatHTML
}

// IsText reports whether the format is plain text.
func (f DocsFormatType) IsText() bool {
	return f == DocsFormatText
}

// IsJSON reports whether the format is JSON.
func (f DocsFormatType) IsJSON() bool {
	return f == DocsFormatJSON
}

// HasFormatting reports whether the format has valid defined formatting.
func (f DocsFormatType) HasFormatting() bool {
	return f.IsValid()
}

// ParseDocsFormatType parses a raw string into a DocsFormatType.
func ParseDocsFormatType(s string) (DocsFormatType, error) {
	val := DocsFormatType(strings.ToLower(strings.TrimSpace(s)))

	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid docs format: %q", s)
}

// ParseDocsFormat is an alias for ParseDocsFormatType.
func ParseDocsFormat(s string) (DocsFormatType, error) {
	return ParseDocsFormatType(s)
}

// AllDocsFormatTypes returns all supported documentation formats.
func AllDocsFormatTypes() []DocsFormatType {
	return []DocsFormatType{
		DocsFormatMarkdown,
		DocsFormatHTML,
		DocsFormatText,
		DocsFormatJSON,
	}
}

// AllDocsFormats is an alias for AllDocsFormatTypes.
func AllDocsFormats() []DocsFormatType {
	return AllDocsFormatTypes()
}
