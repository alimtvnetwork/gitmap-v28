package cmdpipeline

// PEFormatProfile defines a declarative filtering rule set for pipeline logs.
type PEFormatProfile struct {
	Name                string   `json:"name"`
	Alias               string   `json:"alias,omitempty"`
	Description         string   `json:"description,omitempty"`
	StripPrefixes       []string `json:"strip_prefixes,omitempty"`
	StripContains       []string `json:"strip_contains,omitempty"`
	IsSkipEmpty         bool     `json:"skip_empty,omitempty"`
	WarningMarkers      []string `json:"warning_markers,omitempty"`
	ErrorMarkers        []string `json:"error_markers,omitempty"`
	CaptureUntilMarkers []string `json:"capture_until_markers,omitempty"`
	MaxContextLines     int      `json:"max_context_lines,omitempty"`
}

// DefaultPEFormatProfile returns the standard polyglot extraction profile.
func DefaultPEFormatProfile() PEFormatProfile {
	return PEFormatProfile{
		Name:                "default",
		Alias:               "default",
		Description:         "Universal polyglot build and CI/CD error/warning filter",
		StripPrefixes:       defaultStripPrefixes(),
		StripContains:       defaultStripContains(),
		IsSkipEmpty:         true,
		WarningMarkers:      defaultWarningMarkers(),
		ErrorMarkers:        defaultErrorMarkers(),
		CaptureUntilMarkers: defaultCaptureUntilMarkers(),
		MaxContextLines:     25,
	}
}

// BuiltinTauriProfile returns specialized rules for Tauri + Rust desktop pipelines.
func BuiltinTauriProfile() PEFormatProfile {
	return PEFormatProfile{
		Name:                "tauri",
		Alias:               "tauri",
		Description:         "Tauri + Rust + Vite desktop build error and multi-line warning filter",
		StripPrefixes:       tauriStripPrefixes(),
		StripContains:       tauriStripContains(),
		IsSkipEmpty:         true,
		WarningMarkers:      tauriWarningMarkers(),
		ErrorMarkers:        tauriErrorMarkers(),
		CaptureUntilMarkers: tauriCaptureUntilMarkers(),
		MaxContextLines:     35,
	}
}

func defaultStripPrefixes() []string {
	return []string{
		"Compiling ",
		"   Compiling ",
		"Updating ",
		"    Updating ",
		"transforming...",
		"vite v",
		"Browserslist:",
		"dist/",
		"computing gzip size...",
	}
}

func defaultStripContains() []string {
	return []string{
		"Finished `release` profile",
		"modules transformed",
	}
}

func defaultWarningMarkers() []string {
	return []string{
		"warning:",
		"warning[",
		"warning ",
		": warning:",
		"##[warning]",
	}
}

func defaultErrorMarkers() []string {
	return []string{
		"##[error]",
		"--- FAIL:",
		"FAIL\t",
		"FAIL:",
		"FAILED",
		"failed to bundle",
		"Error failed to bundle",
		"failed to copy",
		"failed to build",
		"does not exist",
		"cannot find",
		"command not found",
		"No such file or directory",
		"fatal error:",
		"syntax error:",
		"AssertionError",
		"panic:",
		"Error:",
		"error:",
		"Process completed with exit code",
	}
}

func defaultCaptureUntilMarkers() []string {
	return []string{
		"Compiling ",
		"   Compiling ",
		"Finished `release`",
		"ok  \t",
		"PASS",
	}
}

func tauriStripPrefixes() []string {
	return append(defaultStripPrefixes(), "    Finished `release`")
}

func tauriStripContains() []string {
	return []string{
		"modules transformed",
		"Some chunks are larger than 500 kB",
	}
}

func tauriWarningMarkers() []string {
	return []string{
		"warning:",
		"warning[",
		"##[warning]",
	}
}

func tauriErrorMarkers() []string {
	return []string{
		"failed to bundle",
		"Error failed to bundle",
		"does not exist",
		"Failed to copy binary",
		"##[error]",
		"Error:",
		"error:",
		"Process completed with exit code",
	}
}

func tauriCaptureUntilMarkers() []string {
	return []string{
		"Compiling ",
		"   Compiling ",
		"Finished `release`",
		"ok  \t",
	}
}
