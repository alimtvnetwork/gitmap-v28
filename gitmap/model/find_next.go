package model

// FindNextRow is the result of joining Repo with its latest available
// VersionProbe (Phase 2.4). Used by `gitmap find-next` to surface every
// repo with a new tag without re-running the probe.
type FindNextRow struct {
	Repo           ScanRecord `json:"repo"`
	NextVersionTag string     `json:"nextVersionTag"`
	NextVersionNum int64      `json:"nextVersionNum"`
	Method         string     `json:"method"`
	ProbedAt       string     `json:"probedAt"`
}
