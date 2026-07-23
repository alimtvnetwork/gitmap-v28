package cmd

// Shared JSON status emitter for `gitmap startup-add` and
// `gitmap startup-remove`. Both commands surface the same shape
// so a single jq filter covers both:
//
//   {
//     "command":   "startup-add" | "startup-remove",
//     "action":    "created" | "overwritten" | "exists" |
//                  "refused" | "bad_name" | "deleted" | "noop",
//     "name":      "<user-provided --name>",
//     "target":    "<abs path>" | "<HKCU\\...>" | "",
//     "owner":     "gitmap" | "third-party" | "none" | "unknown",
//     "force_used": <bool>,
//     "dry_run":    <bool>
//   }
//
// Routes through stablejson so the field order is the same on every
// Go release — downstream jq scripts can rely on positional parsing.
// Empty target is emitted as the literal "" rather than omitted so
// the schema is rectangular (every record has every field).

import (
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// startupStatus is the in-memory representation built by the
// per-command translators. Kept as a struct (not a map) so the
// translators get compile-time field-name checks; the wire format
// is built field-by-field in writeStartupStatusJSON to preserve
// stable key order without reflection.
type startupStatus struct {
	command   string
	action    string
	name      string
	target    string
	owner     string
	forceUsed bool
	dryRun    bool
}

// startupStatusJSON wire-format key names. Centralized so a future
// rename is a single diff and the contract test in
// startupstatusjson_test.go pins both order and labels.
const (
	startupStatusKeyCommand   = "command"
	startupStatusKeyAction    = "action"
	startupStatusKeyName      = "name"
	startupStatusKeyTarget    = "target"
	startupStatusKeyOwner     = "owner"
	startupStatusKeyForceUsed = "force_used"
	startupStatusKeyDryRun    = "dry_run"
)

// Action labels. Mirror the AddStatus / RemoveStatus enums but with
// snake_case strings safe for shell pipelines (jq, awk on the
// `.action` value). Kept as constants so the translators can't
// disagree on spelling.
const (
	StartupActionCreated     = "created"
	StartupActionOverwritten = "overwritten"
	StartupActionExists      = "exists"
	StartupActionRefused     = "refused"
	StartupActionBadName     = "bad_name"
	StartupActionDeleted     = "deleted"
	StartupActionNoOp        = "noop"
)

// Owner labels. Tells the consumer who CURRENTLY holds the on-disk
// entry — independent of what we did to it. "none" appears only
// when there's no entry at all (NoOp / BadName).
const (
	StartupOwnerGitmap     = "gitmap"
	StartupOwnerThirdParty = "third-party"
	StartupOwnerNone       = "none"
	StartupOwnerUnknown    = "unknown"
)

// addResultToStatus translates a startup.AddResult into the wire
// shape. Mirrors printAddResult's switch so the terminal and JSON
// renderers stay in lockstep — adding a new AddStatus value will
// fail the exhaustiveness check in both functions.
func addResultToStatus(name string, force bool, res startup.AddResult) startupStatus {
	s := startupStatus{
		command: constants.CmdStartupAdd, name: name,
		target: res.Path, forceUsed: force,
	}
	switch res.Status {
	case startup.AddCreated:
		s.action, s.owner = StartupActionCreated, StartupOwnerGitmap
	case startup.AddOverwritten:
		s.action, s.owner = StartupActionOverwritten, StartupOwnerGitmap
	case startup.AddExists:
		s.action, s.owner = StartupActionExists, StartupOwnerGitmap
	case startup.AddRefused:
		s.action, s.owner = StartupActionRefused, StartupOwnerThirdParty
	case startup.AddBadName:
		s.action, s.owner = StartupActionBadName, StartupOwnerUnknown
		s.target = ""
	}

	return s
}

// removeResultToStatus translates a startup.RemoveResult into the
// wire shape. Owner inference is more nuanced than Add because a
// successful Delete proves the entry WAS gitmap-managed (only ours
// can be deleted), Refused proves it was third-party, and NoOp
// means nobody owns that name — so the consumer can tell the
// three "we did nothing" cases apart from `.owner` alone without
// re-parsing `.action`.
func removeResultToStatus(name string, res startup.RemoveResult) startupStatus {
	s := startupStatus{
		command: constants.CmdStartupRemove, name: name,
		target: res.Path, dryRun: res.DryRun,
	}
	switch res.Status {
	case startup.RemoveDeleted:
		s.action, s.owner = StartupActionDeleted, StartupOwnerGitmap
	case startup.RemoveNoOp:
		s.action, s.owner = StartupActionNoOp, StartupOwnerNone
		s.target = ""
	case startup.RemoveRefused:
		s.action, s.owner = StartupActionRefused, StartupOwnerThirdParty
	case startup.RemoveBadName:
		s.action, s.owner = StartupActionBadName, StartupOwnerUnknown
		s.target = ""
	}

	return s
}

// writeStartupStatusJSON emits one stablejson-encoded object per
// status. Uses WriteArrayIndent with a single-element array so
// callers get the exact same key-order guarantee as
// `startup-list --format=json`. Default indent is 2 spaces to
// match `startup-list`'s default; pass jsonIndent=0 for a minified
// single-line `{...}\n` shape suitable for line-oriented log
// pipelines.
func writeStartupStatusJSON(w io.Writer, s startupStatus, jsonIndent int) error {
	fields := []stablejson.Field{
		{Key: startupStatusKeyCommand, Value: s.command},
		{Key: startupStatusKeyAction, Value: s.action},
		{Key: startupStatusKeyName, Value: s.name},
		{Key: startupStatusKeyTarget, Value: s.target},
		{Key: startupStatusKeyOwner, Value: s.owner},
		{Key: startupStatusKeyForceUsed, Value: s.forceUsed},
		{Key: startupStatusKeyDryRun, Value: s.dryRun},
	}
	indent := indentSpaces(jsonIndent)

	return stablejson.WriteArrayIndent(w, [][]stablejson.Field{fields}, indent)
}

// emitStartupStatus is the dispatch boundary every CLI runner
// calls. Centralizes the "is the user asking for JSON?" decision so
// printAddResult / printRemoveResult only run when format is the
// default terminal mode. Wired against os.Stdout in the dispatchers;
// tests pass their own writer via writeStartupStatusJSON directly.
func emitStartupStatus(format string, jsonIndent int, s startupStatus) error {
	if format != constants.OutputJSON {
		return nil
	}

	return writeStartupStatusJSON(os.Stdout, s, jsonIndent)
}
