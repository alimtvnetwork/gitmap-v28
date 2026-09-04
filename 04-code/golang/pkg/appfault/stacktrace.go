package appfault

import (
	"fmt"
	"runtime"
	"strings"
)

// StackFrame holds metadata for a single caller frame.
type StackFrame struct {
	Function string `json:"Function,omitempty" yaml:"Function,omitempty"`
	File     string `json:"File,omitempty" yaml:"File,omitempty"`
	Line     int    `json:"Line,omitempty" yaml:"Line,omitempty"`
}

// NewStackFrame constructs a StackFrame with line-by-line assignment.
func NewStackFrame(function string, file string, line int) StackFrame {
	frame := StackFrame{}
	frame.Function = function
	frame.File = file
	frame.Line = line

	return frame
}

// StackTrace is a collection of structured call frames.
type StackTrace []StackFrame

// NewStackTrace creates a StackTrace from a slice of StackFrames.
func NewStackTrace(frames ...StackFrame) StackTrace {
	return StackTrace(frames)
}

// appendFrameIfApp appends frame if not in runtime.
func appendFrameIfApp(trace StackTrace, f runtime.Frame) StackTrace {
	if isAppFrame(f.File) {
		frame := NewStackFrame(f.Function, f.File, f.Line)

		return append(trace, frame)
	}

	return trace
}

// parseFrames extracts non-runtime frames from runtime.Frames.
func parseFrames(frames *runtime.Frames) StackTrace {
	var trace StackTrace
	for {
		f, more := frames.Next()
		trace = appendFrameIfApp(trace, f)
		if !more {
			break
		}
	}

	return trace
}

// CaptureStackTrace captures caller frames starting at skip offset.
func CaptureStackTrace(skip int) StackTrace {
	pc := make([]uintptr, 32)
	n := runtime.Callers(skip+1, pc)
	if n == 0 {
		return StackTrace{}
	}

	return parseFrames(runtime.CallersFrames(pc[:n]))
}

// CaptureCaller captures the top caller line string, reusing CaptureStackTrace.
func CaptureCaller(skip int) string {
	trace := CaptureStackTrace(skip + 1)

	return trace.CallerLine()
}

// isAppFrame filters out runtime frames.
func isAppFrame(file string) bool {
	return !strings.Contains(file, "runtime/")
}

// CallerLine returns a compact "file:line" string of the top frame.
func (st StackTrace) CallerLine() string {
	if len(st) == 0 {
		return "unknown:0"
	}

	return fmt.Sprintf("%s:%d", st[0].File, st[0].Line)
}

// String formats the multi-line stack trace.
func (st StackTrace) String() string {
	var builder strings.Builder
	for idx, frame := range st {
		builder.WriteString(fmt.Sprintf("#%d %s\n   %s:%d\n", idx, frame.Function, frame.File, frame.Line))
	}

	return builder.String()
}
