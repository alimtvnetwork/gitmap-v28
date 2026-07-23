package movemerge

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Choice is the outcome of resolving one conflict.
type Choice int

const (
	// ChoiceLeft writes LEFT's version onto the destination side.
	ChoiceLeft Choice = iota
	// ChoiceRight writes RIGHT's version onto the destination side.
	ChoiceRight
	// ChoiceSkip leaves both sides untouched.
	ChoiceSkip
	// ChoiceQuit aborts the run; partial changes are kept.
	ChoiceQuit
)

// Resolver picks a Choice for each conflict. Stateful: All-Left/Right
// stickiness is held inside the resolver instance.
type Resolver struct {
	policy PreferPolicy
	sticky Choice
	hasStk bool
	in     io.Reader
	out    io.Writer
}

// NewResolver builds a Resolver for the run. When policy is non-None,
// the resolver short-circuits without reading from in.
func NewResolver(policy PreferPolicy, in io.Reader, out io.Writer) *Resolver {
	return &Resolver{policy: policy, in: in, out: out}
}

// Resolve returns the Choice for one conflict. ResolveAuto handles
// the bypass policies; otherwise the interactive prompt is used.
func (r *Resolver) Resolve(rel string, l, rgt FileMeta) (Choice, error) {
	if c, done := r.resolveSticky(); done {
		return c, nil
	}
	if c, done := r.resolveByPolicy(l, rgt); done {
		return c, nil
	}

	return r.resolveInteractive(rel, l, rgt)
}

// resolveSticky returns the sticky choice if All-Left/Right was set.
func (r *Resolver) resolveSticky() (Choice, bool) {
	if r.hasStk {
		return r.sticky, true
	}

	return 0, false
}

// resolveByPolicy applies non-interactive --prefer-* policies.
func (r *Resolver) resolveByPolicy(l, rgt FileMeta) (Choice, bool) {
	switch r.policy {
	case PreferNone:
		return 0, false
	case PreferLeft:
		return ChoiceLeft, true
	case PreferRight:
		return ChoiceRight, true
	case PreferSkip:
		return ChoiceSkip, true
	case PreferNewer:
		if l.Info.ModTime().After(rgt.Info.ModTime()) {
			return ChoiceLeft, true
		}

		return ChoiceRight, true
	}

	return 0, false
}

// resolveInteractive prints the prompt and reads one keystroke.
func (r *Resolver) resolveInteractive(rel string, l, rgt FileMeta) (Choice, error) {
	fmt.Fprintf(r.out, "  conflict: %s\n", rel)
	fmt.Fprintf(r.out, "    LEFT  : %d B  modified %s\n", l.Info.Size(), l.Info.ModTime().Format("2006-01-02 15:04"))
	fmt.Fprintf(r.out, "    RIGHT : %d B  modified %s\n", rgt.Info.Size(), rgt.Info.ModTime().Format("2006-01-02 15:04"))
	fmt.Fprintln(r.out, "  [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit")
	fmt.Fprint(r.out, "  > ")
	scanner := bufio.NewScanner(r.in)
	if !scanner.Scan() {
		return ChoiceQuit, fmt.Errorf("conflict prompt: stdin closed")
	}

	return r.parseKey(strings.TrimSpace(scanner.Text())), nil
}

// parseKey maps a single keystroke to a Choice; sets sticky when A/B.
func (r *Resolver) parseKey(key string) Choice {
	switch strings.ToUpper(key) {
	case "L":
		return ChoiceLeft
	case "R":
		return ChoiceRight
	case "A":
		r.sticky, r.hasStk = ChoiceLeft, true

		return ChoiceLeft
	case "B":
		r.sticky, r.hasStk = ChoiceRight, true

		return ChoiceRight
	case "Q":
		return ChoiceQuit
	}

	return ChoiceSkip
}
