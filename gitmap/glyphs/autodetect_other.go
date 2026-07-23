//go:build !windows

package glyphs

// On non-Windows platforms the system terminal + font combination
// renders emoji reliably, so auto-detect always picks rich mode.
