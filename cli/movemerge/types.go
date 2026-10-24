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
import "github.com/alimtvnetwork/gitmap-v28/cli/movemerge/enums"

// EndpointKindType classifies a positional argument once at command start.
// EndpointKindType aliases the strongly-typed enum from the enums package.
type EndpointKindType = enums.EndpointKindType

const (
	EndpointFolder = enums.EndpointFolder
	EndpointURL    = enums.EndpointURL
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
// PreferPolicyType aliases the strongly-typed enum from the enums package.
type PreferPolicyType = enums.PreferPolicyType

const (
	PreferNone  = enums.PreferNone
	PreferLeft  = enums.PreferLeft
	PreferRight = enums.PreferRight
	PreferNewer = enums.PreferNewer
	PreferSkip  = enums.PreferSkip
)

// DirectionType selects which side(s) the operation writes to.
// DirectionType aliases the strongly-typed enum from the enums package.
type DirectionType = enums.DirectionType

const (
	DirBoth      = enums.DirBoth
	DirLeftOnly  = enums.DirLeftOnly
	DirRightOnly = enums.DirRightOnly
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
// DiffKindType aliases the strongly-typed enum from the enums package.
type DiffKindType = enums.DiffKindType

const (
	DiffMissingLeft  = enums.DiffMissingLeft
	DiffMissingRight = enums.DiffMissingRight
	DiffConflict     = enums.DiffConflict
	DiffIdentical    = enums.DiffIdentical
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
// ChoiceType aliases the strongly-typed enum from the enums package.
type ChoiceType = enums.ChoiceType

const (
	ChoiceLeft  = enums.ChoiceLeft
	ChoiceRight = enums.ChoiceRight
	ChoiceSkip  = enums.ChoiceSkip
	ChoiceQuit  = enums.ChoiceQuit
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
type DiffEntryResult = result.Wrap[DiffEntry]

// FileMetaMapResult wraps a map of relative paths to FileMeta.
type FileMetaMapResult = result.ResultMap[string, FileMeta]
