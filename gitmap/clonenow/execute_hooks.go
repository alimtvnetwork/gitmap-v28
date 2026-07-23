package clonenow

// execute_hooks.go — adds a per-row BeforeRow callback to the
// executor so clone-related CLI commands can stream a per-repo
// terminal block immediately BEFORE each clone shells out, instead
// of dumping all blocks upfront. Keeps Execute (no hook) intact so
// existing call sites and golden tests stay byte-identical.
//
// Design contract (mirrors clonefrom.ExecuteWithHooks):
//
//   - BeforeRow fires synchronously, before the row's git clone
//     starts. The CLI uses it to print the standardized RepoTermBlock.
//   - It receives (index, total, row, resolvedURL, resolvedDest).
//     The CLI already needs all five to render the block; passing
//     them in avoids re-deriving URL/dest in the renderer.
//   - A nil hook is allowed and means "no-op" — that's how the
//     legacy Execute wrapper stays a one-line forward.
//   - The hook MUST NOT mutate the row; it's pass-by-value to make
//     that obvious at the call site.

import (
	"io"
	"os"
)

// BeforeRowHook is invoked once per row, just before the row's git
// clone runs. See file header for the parameter contract.
type BeforeRowHook func(index, total int, row Row, url, dest string)

// ExecuteWithHooks is Execute + a per-row BeforeRow callback. The
// body is a copy of Execute's loop with the hook insertion point;
// kept as a separate function (rather than refactoring Execute to
// accept an optional hook) so the existing Execute signature — and
// every test that calls it — stays untouched.
func ExecuteWithHooks(plan Plan, cwd string, progress io.Writer,
	beforeRow BeforeRowHook) []Result {
	if len(cwd) == 0 {
		if wd, err := os.Getwd(); err == nil {
			cwd = wd
		}
	}
	out := make([]Result, 0, len(plan.Rows))
	total := len(plan.Rows)
	for i, r := range plan.Rows {
		if beforeRow != nil {
			url := r.PickURL(plan.Mode)
			beforeRow(i+1, total, r, url, r.RelativePath)
		}
		res := executeRow(r, plan, cwd)
		out = append(out, res)
		writeProgress(progress, i+1, total, res)
	}

	return out
}
