package cmdpipeline

import (
	"bufio"
	"strings"
)

// IsLineStrippedByProfile checks if a line should be dropped based on the profile.
func IsLineStrippedByProfile(line string, p *PEFormatProfile) bool {
	if p == nil {
		return false
	}
	trimmed := strings.TrimSpace(line)
	if p.IsSkipEmpty && trimmed == "" {
		return true
	}
	for _, pfx := range p.StripPrefixes {
		if strings.HasPrefix(trimmed, pfx) || strings.HasPrefix(line, pfx) {
			return true
		}
	}
	for _, sub := range p.StripContains {
		if strings.Contains(line, sub) {
			return true
		}
	}
	return false
}

// IsWarningByProfile checks if a line begins a warning block.
func IsWarningByProfile(line string, p *PEFormatProfile) bool {
	markers := defaultWarningMarkers()
	if p != nil && len(p.WarningMarkers) > 0 {
		markers = p.WarningMarkers
	}
	trimmed := strings.TrimSpace(line)
	lower := strings.ToLower(trimmed)
	for _, m := range markers {
		if strings.HasPrefix(lower, strings.ToLower(m)) || strings.Contains(lower, strings.ToLower(m)) {
			return true
		}
	}
	return false
}

// IsErrorByProfile checks if a line matches any error markers.
func IsErrorByProfile(line string, p *PEFormatProfile) bool {
	markers := defaultErrorMarkers()
	if p != nil && len(p.ErrorMarkers) > 0 {
		markers = p.ErrorMarkers
	}
	for _, m := range markers {
		if strings.Contains(line, m) || strings.Contains(strings.ToLower(line), strings.ToLower(m)) {
			return true
		}
	}
	return false
}

// IsUntilMarkerByProfile checks if a line ends the active warning/context capture block.
func IsUntilMarkerByProfile(line string, p *PEFormatProfile) bool {
	markers := defaultCaptureUntilMarkers()
	if p != nil && len(p.CaptureUntilMarkers) > 0 {
		markers = p.CaptureUntilMarkers
	}
	trimmed := strings.TrimSpace(line)
	for _, m := range markers {
		if strings.HasPrefix(trimmed, m) || strings.Contains(line, m) {
			return true
		}
	}
	return false
}

// FilterLogWithProfile applies a declarative profile directly to raw log text.
func FilterLogWithProfile(rawLogs string, p *PEFormatProfile) string {
	if p == nil {
		defaultP := DefaultPEFormatProfile()
		p = &defaultP
	}
	scanner := bufio.NewScanner(strings.NewReader(rawLogs))
	var outLines []string
	isCapturingWarn := false
	warnCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		processFilterLine(line, p, &outLines, &isCapturingWarn, &warnCount)
	}
	return strings.Join(outLines, "\n")
}

func processFilterLine(
	line string,
	p *PEFormatProfile,
	out *[]string,
	isCapturing *bool,
	count *int,
) {
	if IsWarningByProfile(line, p) {
		*isCapturing = true
		*count = 1
		*out = append(*out, line)
		return
	}
	if *isCapturing {
		handleCapturingLine(line, p, out, isCapturing, count)
		return
	}
	if IsLineStrippedByProfile(line, p) {
		return
	}
	if IsErrorByProfile(line, p) {
		*out = append(*out, line)
	}
}

func handleCapturingLine(
	line string,
	p *PEFormatProfile,
	out *[]string,
	isCapturing *bool,
	count *int,
) {
	isTerminalLine := IsUntilMarkerByProfile(line, p) || *count >= p.MaxContextLines
	if isTerminalLine {
		terminateCaptureAndAppendIfError(out, line, p, isCapturing)

		return
	}
	*count++
	*out = append(*out, line)
}

func terminateCaptureAndAppendIfError(out *[]string, line string, p *PEFormatProfile, isCapturing *bool) {
	*isCapturing = false
	isMatchingError := !IsLineStrippedByProfile(line, p) && IsErrorByProfile(line, p)
	if isMatchingError {
		*out = append(*out, line)
	}
}
