// Package render — pretty-mode policy.
//
// PrettyMode lets callers override the default TTY+env auto-detection that
// gates ANSI rendering across the CLI. The shared decision function
// (Decide) centralizes precedence so every command surface that prints
// markdown answers the question "should I emit ANSI?" the same way:
//
//	user explicit  >  env opt-out  >  TTY auto-detect  >  content gate
//
// Adding a new pretty-render surface means: parse a PrettyMode from args
// (see cmd.parsePrettyFlag) and call Decide — never re-implement the
// gate inline.
package render

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/theme"
)

// PrettyMode is the tri-state requested by the caller.
type PrettyMode int

const (
	// PrettyAuto is the default: render when stdout is a real TTY and the
	// shared GITMAP_NO_PRETTY env opt-out is unset.
	PrettyAuto PrettyMode = iota
	// PrettyOn forces ANSI rendering regardless of TTY / env. Useful when
	// piping into a renderer that understands ANSI (e.g. `less -R`).
	PrettyOn
	// PrettyOff suppresses ANSI rendering regardless of TTY / env. Use
	// when redirecting to a file or feeding the output into tools that
	// choke on escape codes (diff, sha256sum, grep -F).
	PrettyOff
)

// EnvNoPretty is the shared opt-out env var name. Mirrors the convention
// of NO_COLOR but scoped to gitmap's pretty markdown pipeline so users
// can disable ANSI rendering across `help`, `templates show`, and
// `changelog` with a single export.
const EnvNoPretty = "GITMAP_NO_PRETTY"

// Decide returns true when the caller should emit ANSI-rendered markdown.
//
// Inputs:
//
//   - mode:       caller-supplied tri-state (parsed from --pretty / --no-pretty);
//   - isTTY:      whether stdout is connected to a real terminal;
//   - isMarkdown: whether the content is markdown at all (callers that
//     never produce markdown can pass true to skip this gate).
//
// Precedence (first match wins):
//
//  1. PrettyOff      → false (explicit user opt-out trumps everything).
//  2. !isMarkdown    → false (no point rendering plain text as markdown).
//  3. PrettyOn       → true  (explicit user opt-in trumps env / TTY).
//  4. EnvNoPretty=…  → false (shared environment opt-out).
//  5. isTTY          → true  (auto-detect: render only on a real terminal).
//  6. otherwise      → false (pipes / redirects stay byte-faithful).
func Decide(mode PrettyMode, isTTY, isMarkdown bool) bool {
	if mode == PrettyOff {
		return false
	}
	if !isMarkdown {
		return false
	}
	if mode == PrettyOn {
		return true
	}
	if os.Getenv(EnvNoPretty) != "" {
		return false
	}

	return isTTY
}

// StdoutIsTerminal reports whether os.Stdout is connected to a real TTY
// using a dependency-free os.Stat / ModeCharDevice check. Exported so
// callers in cmd/, helptext/, and any future markdown-printing surface
// can share one TTY probe instead of reimplementing it locally.
//
// When the theme package has wrapped os.Stdout with a pipe (standard /
// monochrome modes) the live Stat would report "pipe" and falsely
// disable pretty rendering — we defer to the pre-wrap TTY state
// theme.Install captured before swapping the handle.
func StdoutIsTerminal() bool {
	if theme.Active() != theme.ModeBright {
		return theme.IsStdoutTTY()
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (info.Mode() & os.ModeCharDevice) != 0
}

