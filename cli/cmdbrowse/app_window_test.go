package cmdbrowse

import (
	"testing"
)

func TestFindEdgeBinaryPath(t *testing.T) {
	// Should not panic on any platform
	_ = findEdgeBinaryPath()
}

func TestFindAppBrowserBinary(t *testing.T) {
	// Should return a string without panicking
	_ = findAppBrowserBinary()
}
