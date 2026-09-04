package appfault

import (
	"fmt"
	"runtime"
)

// CallerInfo encapsulates caller site metadata.
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
