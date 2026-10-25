# Component and CLI Specification: Flexible Path Parsing and Safe Image Cat

- **Task Sequence:** 328
- **Task Slug:** `flexible-path-parsing-and-asset-image-reading`
- **Spec Path:** `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/02-component-and-cli-spec.md`
- **Created At:** 2026-10-25
- **Status:** `APPROVED`

---

## 1. Overview & Architectural Goals

This specification details the structural design, Go types, function signatures, CLI flags, and interactive terminal interfaces for Task 328.

The goals of this component suite are:
1. Provide a reusable parser in `cli/pathparse` capable of taking raw user or agent input in any of the 4 target formats and converting it into a canonical, normalized file path.
2. Upgrade `gitmap cat` in `cli/cmdmacro/` to integrate `cli/pathparse`, enabling direct execution on conversational strings, Markdown images, and bracketed paths.
3. Eliminate terminal buffer corruption when reading binary images by introducing an automated inspection card with an explicit `--raw` streaming flag.

---

## 2. Component Specification: `cli/pathparse` Package

### 2.1 Struct Definitions (`cli/pathparse/types.go`)

```go
package pathparse

// ParsedPathResult holds the parsed and resolved path properties.
type ParsedPathResult struct {
	RawInput            string `json:"rawInput"`
	CleanPath           string `json:"cleanPath"`
	ResolvedPath        string `json:"resolvedPath"`
	TrailingComment     string `json:"trailingComment"`
	HasMarkdownWrapper  bool   `json:"hasMarkdownWrapper"`
	HasBracketWrapper   bool   `json:"hasBracketWrapper"`
	HasListPrefix       bool   `json:"hasListPrefix"`
	HasTrailingComment  bool   `json:"hasTrailingComment"`
	IsAbsolute          bool   `json:"isAbsolute"`
	IsInsideRepo        bool   `json:"isInsideRepo"`
	IsParentRelative    bool   `json:"isParentRelative"`
	IsFileExists        bool   `json:"isFileExists"`
	IsDirectory         bool   `json:"isDirectory"`
	IsImage             bool   `json:"isImage"`
	ImageFormat         string `json:"imageFormat"`
}

// ImageMetadata contains decoded dimensions and properties of an image asset.
type ImageMetadata struct {
	Format      string `json:"format"`
	MimeType    string `json:"mimeType"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	SizeBytes   int64  `json:"sizeBytes"`
	Sha256Hex   string `json:"sha256Hex"`
	IsVector    bool   `json:"isVector"`
}
```

### 2.2 Core Parsing Functions (`cli/pathparse/parser.go`)

Each function strictly adheres to the 15-line limit, returns typed errors, and uses positive booleans.

```go
// ParseFlexiblePath converts raw conversational or markdown paths into resolved paths.
func ParseFlexiblePath(raw string, repoRoot string) (ParsedPathResult, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ParsedPathResult{}, apperror.NewSimple("path cannot be empty", "E_PATH_EMPTY")
	}

	step1, hasList := StripListPrefix(trimmed)
	step2, comment, hasComment := StripInlineComment(step1)
	step3, hasMd, hasBracket := StripWrappers(step2)
	decoded, err := DecodeURIPath(step3)
	if err != nil {
		return ParsedPathResult{}, err
	}

	return assembleParsedResult(raw, decoded, comment, hasList, hasComment, hasMd, hasBracket, repoRoot)
}
```

### 2.3 Stage Decompositions ($\le 15$ Lines Each)

#### StripListPrefix
Strips leading numbered lists (`1. `, `02) `) and bullet markers (`- `, `* `):
```go
func StripListPrefix(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	for i, r := range trimmed {
		if (r >= '0' && r <= '9') || r == '-' || r == '*' {
			continue
		}
		if (r == '.' || r == ')') && i+1 < len(trimmed) && trimmed[i+1] == ' ' {
			return strings.TrimSpace(trimmed[i+2:]), true
		}
		break
	}

	return trimmed, false
}
```

#### StripInlineComment
Extracts and removes trailing comments starting with `#` or `//`:
```go
func StripInlineComment(input string) (string, string, bool) {
	if idx := strings.Index(input, " #"); idx != -1 {
		return strings.TrimSpace(input[:idx]), strings.TrimSpace(input[idx+2:]), true
	}
	if idx := strings.Index(input, " //"); idx != -1 {
		return strings.TrimSpace(input[:idx]), strings.TrimSpace(input[idx+3:]), true
	}

	return input, "", false
}
```

#### StripWrappers
Removes Markdown image syntax `![alt](url)` and bracket syntax `@[path]`:
```go
func StripWrappers(input string) (string, bool, bool) {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "![") && strings.HasSuffix(trimmed, ")") {
		if openParen := strings.LastIndex(trimmed, "("); openParen != -1 {
			return strings.TrimSpace(trimmed[openParen+1 : len(trimmed)-1]), true, false
		}
	}
	if strings.HasPrefix(trimmed, "@[") && strings.HasSuffix(trimmed, "]") {
		return strings.TrimSpace(trimmed[2 : len(trimmed)-1]), false, true
	}
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1]), false, true
	}

	return trimmed, false, false
}
```

#### DecodeURIPath
Removes `file:///` and URL decodes escape sequences:
```go
func DecodeURIPath(input string) (string, error) {
	clean := strings.Trim(strings.TrimSpace(input), `"'`)
	if strings.HasPrefix(clean, "file:///") {
		clean = clean[8:]
	} else if strings.HasPrefix(clean, "file://") {
		clean = clean[7:]
	}

	unescaped, err := url.PathUnescape(clean)
	if err != nil {
		return "", apperror.WrapSimple(err, "invalid URL percent-encoding in path")
	}

	if len(unescaped) >= 3 && unescaped[0] == '/' && unescaped[2] == ':' {
		unescaped = unescaped[1:]
	}

	return filepath.ToSlash(unescaped), nil
}
```

---

## 3. Component Specification: `cli/cmdmacro/cat.go` Upgrade

### 3.1 Updated `RunCat` Dispatcher

```go
package cmdmacro

// RunCat streams or inspects files with automatic flexible path parsing and safe image protection.
func RunCat(args []string) error {
	opts, rawTarget, err := parseCatFlags(args)
	if err != nil {
		return err
	}

	parsed, err := pathparse.ParseFlexiblePath(rawTarget, ".")
	if err != nil {
		return err
	}

	if parsed.IsImage && !opts.HasRawFlag && isInteractiveTerminal() {
		return renderSafeImageInspection(parsed.ResolvedPath, parsed.CleanPath)
	}

	return streamRawFile(parsed.ResolvedPath)
}
```

### 3.2 Terminal Image Inspection Card (`cli/cmdmacro/cat_image.go`)

```go
// RenderImageCard outputs a styled terminal metadata box.
func RenderImageCard(meta pathparse.ImageMetadata, displayPath string) {
	termout.PrintNoticeCard("GitMap Asset Inspector: Safe Image Reader", []termout.Field{
		{Label: "File", Value: displayPath},
		{Label: "Format", Value: fmt.Sprintf("%s (%s)", meta.Format, meta.MimeType)},
		{Label: "Dimensions", Value: fmt.Sprintf("%d × %d px", meta.Width, meta.Height)},
		{Label: "File Size", Value: formatByteSize(meta.SizeBytes)},
		{Label: "SHA-256", Value: meta.Sha256Hex},
	}, "[SAFE IMAGE CAT: Terminal binary corruption prevented]\nTip: Use 'gitmap cat --raw <file>' to stream binary bytes.")
}
```

---

## 4. CLI Command Usage & Output Specifications

### 4.1 Format 1: Bracketed Windows Absolute Path
**Command:**
```powershell
gitmap cat "@[C:\Users\Administrator\.antigravity_tools\instances\v2-3127\home\.gemini\antigravity\brain\1ae0f958\media_1792894009084_45b4893f.png]"
```
**Parser Behavior:**
- Extracts `C:\Users\Administrator\...\media_1792894009084_45b4893f.png`
- Sets `HasBracketWrapper = true`, `IsAbsolute = true`
- Identifies `.png` extension and binary header
- Outputs safe image metadata card

### 4.2 Format 2: Markdown Image Syntax with URI Scheme
**Command:**
```powershell
gitmap cat "![Attached Image 1](file:///<drive>:/path/to/media_1792894009084_45b4893f.png)"
```
**Parser Behavior:**
- Extracts URL target from `(...)`
- Strips `file:///`, decodes unescaped path
- Sets `HasMarkdownWrapper = true`, `IsAbsolute = true`
- Outputs safe image metadata card

### 4.3 Format 3: Numbered List, Repo-Relative Path with Trailing Comment
**Command:**
```powershell
gitmap cat "1. @[assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo"
```
**Parser Behavior:**
- Strips `1. ` list prefix
- Strips ` # consider the file is relative from the repo` comment
- Extracts `assets/media_1792894009084_45b4893f.png`
- Resolves against repo root: `<repo-root>/assets/media_1792894009084_45b4893f.png`
- Sets `HasListPrefix = true`, `HasTrailingComment = true`, `IsInsideRepo = true`

### 4.4 Format 4: Numbered List, Parent-Relative Path with Trailing Comment
**Command:**
```powershell
gitmap cat "1. @[../repo-cache/assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo outside of the repo"
```
**Parser Behavior:**
- Strips `1. `, strips `@[...]`, strips `# ...` comment
- Detects `../` traversal
- Resolves relative to repo parent: `D:/work/repo-cache/assets/media_1792894009084_45b4893f.png`
- Sets `IsParentRelative = true`, `IsInsideRepo = false`

### 4.5 Binary Byte Streaming via `--raw` Flag
**Command:**
```powershell
gitmap cat --raw @[assets/logo.png] > output.png
```
**Behavior:**
- Detects `--raw` flag
- Bypasses terminal inspection card
- Streams raw binary bytes directly to stdout pipe
