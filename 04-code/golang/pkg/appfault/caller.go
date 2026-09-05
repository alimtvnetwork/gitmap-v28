package appfault

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// CallerInfo encapsulates caller site metadata as a value object.
type CallerInfo struct {
	Function string `json:"Function,omitempty" yaml:"Function,omitempty"`
	File     string `json:"File,omitempty" yaml:"File,omitempty"`
	Line     int    `json:"Line,omitempty" yaml:"Line,omitempty"`
}

// String returns formatted "file:line (function)" or "file:line".
func (c CallerInfo) String() string {
	if c.Line == 0 && len(c.File) == 0 {
		return ""
	}

	if len(c.Function) > 0 {
		return fmt.Sprintf("%s:%d (%s)", c.File, c.Line, c.Function)
	}

	return fmt.Sprintf("%s:%d", c.File, c.Line)
}

// IsEmpty returns true if caller info is blank.
func (c CallerInfo) IsEmpty() bool {
	return c.Line == 0 && len(c.File) == 0
}

// extractFuncAndFile extracts function name and file details from pc.
func extractFuncAndFile(pc uintptr, file string, line int) CallerInfo {
	fnName := ""
	if fn := runtime.FuncForPC(pc); fn != nil {
		fnName = fn.Name()
	}

	return CallerInfo{Function: fnName, File: file, Line: line}
}

// CaptureCallerInfo captures caller metadata at a given stack skip depth.
func CaptureCallerInfo(skip int) CallerInfo {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return CallerInfo{}
	}

	return extractFuncAndFile(pc, file, line)
}

// MarshalJSON serializes CallerInfo. If empty, returns JSON null.
func (c CallerInfo) MarshalJSON() ([]byte, error) {
	if c.IsEmpty() {
		return []byte("null"), nil
	}

	type alias CallerInfo

	return json.Marshal(alias(c))
}

// UnmarshalJSON unmarshals from JSON object or formatted string.
func (c *CallerInfo) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" || len(trimmed) == 0 {
		*c = CallerInfo{}

		return nil
	}

	// Handle JSON string format: "file:line (function)" or "file:line"
	if strings.HasPrefix(trimmed, `"`) {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}

		*c = ParseCallerInfo(str)

		return nil
	}

	// Handle JSON object format: {"Function":"...","File":"...","Line":...}
	type alias CallerInfo
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*c = CallerInfo(a)

	return nil
}

// ToJson exports CallerInfo as indented JSON bytes.
func (c CallerInfo) ToJson() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// ToJsonString exports CallerInfo as a JSON string.
func (c CallerInfo) ToJsonString() string {
	b, err := c.ToJson()
	if err != nil {
		return "{}"
	}

	return string(b)
}

// ParseCallerInfo parses a formatted string into a CallerInfo value object.
func ParseCallerInfo(s string) CallerInfo {
	clean := strings.TrimSpace(s)
	if len(clean) == 0 {
		return CallerInfo{}
	}

	function := ""
	fileAndLine := clean

	// Pattern: "file:line (function)"
	if openParen := strings.Index(clean, "("); openParen > 0 {
		if closeParen := strings.LastIndex(clean, ")"); closeParen > openParen {
			function = strings.TrimSpace(clean[openParen+1 : closeParen])
			fileAndLine = strings.TrimSpace(clean[:openParen])
		}
	}

	// Pattern: "function at file:line"
	if atIdx := strings.Index(fileAndLine, " at "); atIdx > 0 {
		if len(function) == 0 {
			function = strings.TrimSpace(fileAndLine[:atIdx])
		}

		fileAndLine = strings.TrimSpace(fileAndLine[atIdx+4:])
	}

	file := fileAndLine
	line := 0
	if colonIdx := strings.LastIndex(fileAndLine, ":"); colonIdx > 0 {
		file = fileAndLine[:colonIdx]
		parsedLine, err := strconv.Atoi(fileAndLine[colonIdx+1:])
		if err == nil {
			line = parsedLine
		}
	}

	return CallerInfo{Function: function, File: file, Line: line}
}

// CallerInfoFromJson parses CallerInfo from JSON bytes.
func CallerInfoFromJson(data []byte) (CallerInfo, error) {
	var c CallerInfo
	if err := json.Unmarshal(data, &c); err != nil {
		return CallerInfo{}, err
	}

	return c, nil
}
