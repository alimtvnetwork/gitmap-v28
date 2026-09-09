// Package cmd — macro_add_history.go: interactive terminal history and arrow navigation for macro add.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type stdReadWriter struct {
	io.Reader
	io.Writer
}

type interactiveLineReader struct {
	terminal *term.Terminal
	scanner  *bufio.Scanner
	fd       int
	hasTerm  bool
}

func newInteractiveLineReader() *interactiveLineReader {
	fd := int(os.Stdin.Fd())
	hasTerm := term.IsTerminal(fd)
	if hasTerm {
		return createTerminalLineReader(fd)
	}

	return createScannerLineReader(fd)
}

func createTerminalLineReader(fd int) *interactiveLineReader {
	rw := &stdReadWriter{Reader: os.Stdin, Writer: os.Stdout}
	t := term.NewTerminal(rw, "")
	seedMacroHistory(t)

	return &interactiveLineReader{
		terminal: t,
		fd:       fd,
		hasTerm:  true,
	}
}

func createScannerLineReader(fd int) *interactiveLineReader {
	return &interactiveLineReader{
		scanner: bufio.NewScanner(os.Stdin),
		fd:      fd,
		hasTerm: false,
	}
}

func (r *interactiveLineReader) readLine(prompt string) (string, bool, error) {
	if r.hasTerm {
		return r.readRawTerminalLine(prompt)
	}

	return r.readCookedScannerLine(prompt)
}

func (r *interactiveLineReader) readRawTerminalLine(prompt string) (string, bool, error) {
	oldState, err := term.MakeRaw(r.fd)
	if err != nil {
		return r.readCookedScannerLine(prompt)
	}

	defer term.Restore(r.fd, oldState)
	r.terminal.SetPrompt(prompt)

	line, readErr := r.terminal.ReadLine()
	if readErr != nil {
		return handleTerminalReadError(readErr)
	}

	return strings.TrimSpace(line), false, nil
}

func handleTerminalReadError(err error) (string, bool, error) {
	if err == io.EOF {
		return "", true, nil
	}

	return "", false, err
}

func (r *interactiveLineReader) readCookedScannerLine(prompt string) (string, bool, error) {
	if r.scanner == nil {
		r.scanner = bufio.NewScanner(os.Stdin)
	}

	fmt.Print(prompt)
	if !r.scanner.Scan() {
		return "", true, r.scanner.Err()
	}

	clean := sanitizeRawEscapeCodes(r.scanner.Text())

	return strings.TrimSpace(clean), false, nil
}
