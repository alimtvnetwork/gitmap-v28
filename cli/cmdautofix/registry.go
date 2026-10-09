package cmdautofix

var categoryRegistry = []Category{
	{
		Name:  "encoding",
		Desc:  "UTF-8 no-BOM + LF normalization",
		Check: encodingCheck,
		Fix:   encodingFix,
	},
	{
		Name:  "newlines",
		Desc:  "CRLF→LF, trim trailing whitespace, one final newline",
		Check: newlinesCheck,
		Fix:   newlinesFix,
	},
	{
		Name:  "naming",
		Desc:  "boolean-comparison style audit (report-only)",
		Exts:  codeExts,
		Check: namingCheck,
		// Fix is nil: report-only by design (D11).
	},
	{
		Name:  "paths",
		Desc:  "Windows absolute path sanitizer (docs)",
		Check: pathsCheck,
		Fix:   pathsFix,
	},
	{
		Name:  "gofmt",
		Desc:  "gofmt -w over .go files",
		Exts:  []string{".go"},
		Check: gofmtCheck,
		Fix:   gofmtFix,
		Exec:  true,
	},
	{
		Name:  "misspell",
		Desc:  "British→American spelling, case-preserving",
		Exts:  []string{".md", ".go", ".ts", ".tsx", ".js", ".py", ".json", ".sh", ".ps1", ".txt"},
		Check: misspellCheck,
		Fix:   misspellFix,
	},
	{
		Name:  "markdown",
		Desc:  "collapse 3+ blank lines in .md",
		Exts:  []string{".md"},
		Check: markdownCheck,
		Fix:   markdownFix,
	},
	{
		Name:  "guidelines",
		Desc:  "composite: newlines + naming",
		Check: guidelinesCheck,
		Fix:   guidelinesFix,
	},
}
