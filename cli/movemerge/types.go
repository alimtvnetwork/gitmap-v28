// Package movemerge implements the file-level move and merge family:
// `gitmap mv`, `merge-both`, `merge-left`, `merge-right`. Each
// endpoint (LEFT or RIGHT) is either a local folder or a remote git
// URL with optional :branch suffix. The package resolves both into
// working folders, performs the requested file operation, and (for
// URL endpoints) commits + pushes the result.
//
// Spec: 02-spec/01-app/97-move-and-merge.md
package movemerge

import (
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// EndpointKindType classifies a positional argument once at command start.
type EndpointKindType int

const (
	// EndpointFolder is a plain on-disk path (relative or absolute).
	EndpointFolder EndpointKindType = iota
	// EndpointURL is an https/http/ssh/git@ remote, optionally :branch.
	EndpointURL
)

// Endpoint is a fully resolved LEFT or RIGHT argument.
type Endpoint struct {
	Raw         string           // original CLI token
	DisplayName string           // trimmed, used in commit messages and logs
	Kind        EndpointKindType // folder or URL
	URL         string           // canonical URL when Kind == EndpointURL
	Branch      string           // optional :branch suffix; "" when omitted
	WorkingDir  string           // absolute resolved working folder
	IsGitRepo   bool             // true when WorkingDir contains .git/
	IsExisted   bool             // true when WorkingDir already existed pre-resolve
}

// PreferPolicyType is how -y / --prefer-* resolve conflicts non-interactively.
type PreferPolicyType int

const (
	// PreferNone means use the interactive prompt.
	PreferNone PreferPolicyType = iota
	// PreferLeft makes LEFT always win.
	PreferLeft
	// PreferRight makes RIGHT always win.
	PreferRight
	// PreferNewer compares mtime; newer side wins.
	PreferNewer
	// PreferSkip skips every conflict (only missing files copied).
	PreferSkip
)

// DirectionType selects which side(s) the operation writes to.
type DirectionType int

const (
	// DirBoth writes into both sides (merge-both).
	DirBoth DirectionType = iota
	// DirLeftOnly writes only into LEFT (merge-left).
	DirLeftOnly
	// DirRightOnly writes only into RIGHT (merge-right).
	DirRightOnly
)

// Options bundles every CLI flag for the move/merge family.
type Options struct {
	IsYes             bool
	Prefer            PreferPolicyType
	IsSkipPush        bool
	IsSkipCommit      bool
	IsForceFolder     bool
	IsPullFolder      bool
	IsInitNewRight    bool
	IsDryRun          bool
	IsIncludeVCS      bool
	IsIncludeNodeMods bool
	CommandName       string // "mv" | "merge-both" | "merge-left" | "merge-right"
	LogPrefix         string // "[mv]" etc.
	CommitMsgFmt      string // template; "%s" filled from other side's display
}

// DiffKindType classifies a path across LEFT and RIGHT.
type DiffKindType int

const (
	// DiffMissingLeft = present on RIGHT only.
	DiffMissingLeft DiffKindType = iota
	// DiffMissingRight = present on LEFT only.
	DiffMissingRight
	// DiffConflict = present on both with different content.
	DiffConflict
	// DiffIdentical = present on both with byte-equal content.
	DiffIdentical
)

// DiffEntry is one classified path with both sides' metadata.
type DiffEntry struct {
	RelPath string
	Kind    DiffKindType
	Left    FileMeta
	Right   FileMeta
}

// FileMeta is a single file's identity used by the diff stage.
type FileMeta struct {
	RelPath string
	Info    os.FileInfo
	SHA     string
}

// ChoiceType is the outcome of resolving one conflict.
type ChoiceType int

const (
	// ChoiceLeft writes LEFT's version onto the destination side.
	ChoiceLeft ChoiceType = iota
	// ChoiceRight writes RIGHT's version onto the destination side.
	ChoiceRight
	// ChoiceSkip leaves both sides untouched.
	ChoiceSkip
	// ChoiceQuit aborts the run; partial changes are kept.
	ChoiceQuit
)

// Resolver picks a ChoiceType for each conflict. Stateful: All-Left/Right
// stickiness is held inside the resolver instance.
type Resolver struct {
	policy PreferPolicyType
	sticky ChoiceType
	hasStk bool
	in     io.Reader
	out    io.Writer
}

// DiffEntryResult wraps a single DiffEntry in a Result envelope.
type DiffEntryResult = result.Result[DiffEntry]

// FileMetaMapResult wraps a map of relative paths to FileMeta.
type FileMetaMapResult = result.ResultMap[string, FileMeta]
