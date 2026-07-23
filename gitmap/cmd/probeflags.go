package cmd

// probeflags.go — argument parsing for `gitmap probe`.
//
// Split out of probe.go to honor the 200-line per-file budget. The
// dispatcher (runProbe) lives in probe.go and consumes probeOptions
// produced here. parseProbeArgs is order-agnostic and supports both
// `--flag N` and `--flag=N` forms for every value flag.
//
// Flag map (v3.135.0+):
//   --json                       boolean
//   --probe-workers N            int, capped to [1, ProbeMaxWorkers]
//   --workers N                  deprecated alias for --probe-workers
//   --probe-depth N              int, passed to shallow-clone fallback
//   <positional>                 forwarded to opts.rest

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// probeOptions captures the parsed CLI flags for `gitmap probe`.
// workers is already clamped into [1, ProbeMaxWorkers] by parseProbeArgs;
// depth defaults to constants.ProbeDefaultDepth (1) and is coerced to >=1
// inside the shallow-clone fallback regardless.
type probeOptions struct {
	jsonOut bool
	// termOut, when true, prints the standardized RepoTermBlock
	// per repo so the per-probe summary matches the shape used by
	// scan/clone-from/clone-next. Set via `--output terminal`.
	termOut bool
	workers int
	depth   int
	rest    []string
}

// parseProbeArgs walks the arg list, peeling off recognized flags and
// returning everything else as positional args. Order-agnostic.
func parseProbeArgs(args []string) (probeOptions, error) {
	opts := probeOptions{
		workers: constants.ProbeDefaultWorkers,
		depth:   constants.ProbeDefaultDepth,
		rest:    make([]string, 0, len(args)),
	}
	for i := 0; i < len(args); i++ {
		next, consumed, err := applyProbeFlag(&opts, args, i)
		if err != nil {
			return opts, err
		}
		if !consumed {
			opts.rest = append(opts.rest, args[i])
		}
		i = next
	}

	return opts, nil
}

// applyProbeFlag dispatches a single token to the right handler.
// Returns the new loop index and a `consumed` bool — false means the
// caller should treat args[i] as positional.
func applyProbeFlag(opts *probeOptions, args []string, i int) (int, bool, error) {
	a := args[i]
	if a == constants.ProbeFlagJSON {
		opts.jsonOut = true
		return i, true, nil
	}
	if next, ok, err := applyOutputFlag(opts, args, i); ok || err != nil {
		return next, true, err
	}
	if next, ok, err := applyWorkersFlag(opts, args, i); ok || err != nil {
		return next, true, err
	}
	if next, ok, err := applyDepthFlag(opts, args, i); ok || err != nil {
		return next, true, err
	}

	return i, false, nil
}

// applyOutputFlag handles `--output terminal` (and the inline
// `--output=terminal` form). Any other value is rejected with a
// pointed error so users don't silently get the wrong format. The
// flag flips opts.termOut; emission lives in probereport.go.
func applyOutputFlag(opts *probeOptions, args []string, i int) (int, bool, error) {
	if !matchesFlag(args[i], constants.ProbeFlagOutput) {
		return i, false, nil
	}
	val, next, err := readStringFlag(args, i)
	if err != nil {
		return i, true, err
	}
	if val != constants.OutputTerminal {
		return i, true, fmt.Errorf(
			"version probe: --output only supports %q, got %q",
			constants.OutputTerminal, val)
	}
	opts.termOut = true

	return next, true, nil
}

// readStringFlag mirrors readIntFlag for string-valued flags.
// Supports both `--flag value` and `--flag=value`.
func readStringFlag(args []string, i int) (string, int, error) {
	if eq := strings.IndexByte(args[i], '='); eq >= 0 {
		return args[i][eq+1:], i, nil
	}
	if i+1 >= len(args) {
		return "", i, fmt.Errorf("version probe: %s requires a value", args[i])
	}

	return args[i+1], i + 1, nil
}

// applyWorkersFlag handles --probe-workers / --workers (with a
// deprecation notice on the latter), in both `--flag N` and
// `--flag=N` forms.
func applyWorkersFlag(opts *probeOptions, args []string, i int) (int, bool, error) {
	a := args[i]
	if matchesFlag(a, constants.ProbeFlagWorkers) {
		fmt.Fprint(os.Stderr, constants.MsgProbeWorkersAlias)
	} else if !matchesFlag(a, constants.ProbeFlagProbeWorkers) {
		return i, false, nil
	}
	n, next, err := readIntFlag(args, i, constants.ErrProbeWorkersMissing, constants.ErrProbeWorkersValue)
	if err != nil {
		return i, true, err
	}
	opts.workers = clampProbeWorkers(n)

	return next, true, nil
}

// applyDepthFlag handles --probe-depth in both forms. depth<1 is
// rejected with the standard error message; coercion to >=1 happens
// inside tryShallowClone as a defensive safety net.
func applyDepthFlag(opts *probeOptions, args []string, i int) (int, bool, error) {
	if !matchesFlag(args[i], constants.ProbeFlagDepth) {
		return i, false, nil
	}
	n, next, err := readIntFlag(args, i, constants.ErrProbeDepthMissing, constants.ErrProbeDepthValue)
	if err != nil {
		return i, true, err
	}
	opts.depth = n

	return next, true, nil
}

// matchesFlag returns true when token is exactly `flag` or starts
// with `flag=` (the inline-value form).
func matchesFlag(token, flag string) bool {
	return token == flag || strings.HasPrefix(token, flag+"=")
}

// readIntFlag reads the int value for the flag at args[i] in either
// `--flag N` or `--flag=N` form. Returns (value, newIndex, error).
// newIndex is i for inline form and i+1 for two-arg form.
func readIntFlag(args []string, i int, missingFmt, valueFmt string) (int, int, error) {
	if eq := strings.IndexByte(args[i], '='); eq >= 0 {
		n, err := parsePositiveInt(args[i][eq+1:], valueFmt)
		return n, i, err
	}
	if i+1 >= len(args) {
		return 0, i, errors.New(missingFmt)
	}
	n, err := parsePositiveInt(args[i+1], valueFmt)

	return n, i + 1, err
}

// parsePositiveInt parses a strictly-positive int, formatting the
// caller-supplied error template on failure.
func parsePositiveInt(s, errFmt string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, fmt.Errorf(errFmt, s)
	}

	return n, nil
}

// clampProbeWorkers enforces the [1, ProbeMaxWorkers] cap, printing a
// notice to stderr when the user asked for more than we'll grant.
func clampProbeWorkers(n int) int {
	if n > constants.ProbeMaxWorkers {
		fmt.Fprintf(os.Stderr, constants.MsgProbeWorkersClamped, n, constants.ProbeMaxWorkers)
		return constants.ProbeMaxWorkers
	}

	return n
}
