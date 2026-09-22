package cmdclone

// MultiCloneItem represents a parsed and normalized repository candidate.
type MultiCloneItem struct {
	URL           string
	RepoName      string
	TargetDir     string
	TargetAbsPath string
}

// MultiCloneOptions encapsulates flags and inputs for batch multiclone.
type MultiCloneOptions struct {
	TargetDir    string
	RawInput     string
	FilePath     string
	IsDryRun     bool
	IsConcurrent bool
	Concurrency  int
	NoReplace    bool
	NoVSCodeSync bool
	GHDesktop    bool
}

// MultiCloneSummary aggregates the batch cloning results.
type MultiCloneSummary struct {
	TotalParsed    int
	UniqueCount    int
	DuplicateCount int
	Succeeded      int
	Failed         int
}
