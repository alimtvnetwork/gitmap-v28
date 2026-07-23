// Package movemerge implements the file-level move and merge family:
// `gitmap mv`, `merge-both`, `merge-left`, `merge-right`. Each
// endpoint (LEFT or RIGHT) is either a local folder or a remote git
// URL with optional :branch suffix. The package resolves both into
// working folders, performs the requested file operation, and (for
// URL endpoints) commits + pushes the result.
//
// Spec: spec/01-app/97-move-and-merge.md
package movemerge

// EndpointKind classifies a positional argument once at command start.
type EndpointKind int

const (
	// EndpointFolder is a plain on-disk path (relative or absolute).
	EndpointFolder EndpointKind = iota
	// EndpointURL is an https/http/ssh/git@ remote, optionally :branch.
	EndpointURL
)

// Endpoint is a fully resolved LEFT or RIGHT argument.
type Endpoint struct {
	Raw         string       // original CLI token
	DisplayName string       // trimmed, used in commit messages and logs
	Kind        EndpointKind // folder or URL
	URL         string       // canonical URL when Kind == EndpointURL
	Branch      string       // optional :branch suffix; "" when omitted
	WorkingDir  string       // absolute resolved working folder
	IsGitRepo   bool         // true when WorkingDir contains .git/
	Existed     bool         // true when WorkingDir already existed pre-resolve
}

// PreferPolicy is how -y / --prefer-* resolve conflicts non-interactively.
type PreferPolicy int

const (
	// PreferNone means use the interactive prompt.
	PreferNone PreferPolicy = iota
	// PreferLeft makes LEFT always win.
	PreferLeft
	// PreferRight makes RIGHT always win.
	PreferRight
	// PreferNewer compares mtime; newer side wins.
	PreferNewer
	// PreferSkip skips every conflict (only missing files copied).
	PreferSkip
)

// Direction selects which side(s) the operation writes to.
type Direction int

const (
	// DirBoth writes into both sides (merge-both).
	DirBoth Direction = iota
	// DirLeftOnly writes only into LEFT (merge-left).
	DirLeftOnly
	// DirRightOnly writes only into RIGHT (merge-right).
	DirRightOnly
)

// Options bundles every CLI flag for the move/merge family.
type Options struct {
	Yes             bool
	Prefer          PreferPolicy
	NoPush          bool
	NoCommit        bool
	ForceFolder     bool
	PullFolder      bool
	InitNewRight    bool
	DryRun          bool
	IncludeVCS      bool
	IncludeNodeMods bool
	CommandName     string // "mv" | "merge-both" | "merge-left" | "merge-right"
	LogPrefix       string // "[mv]" etc.
	CommitMsgFmt    string // template; "%s" filled from other side's display
}
