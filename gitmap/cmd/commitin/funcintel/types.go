package funcintel

// Detector returns names of top-level declarations present in newSrc
// but not in prevSrc, sorted ascending and deduped. The detection is
// best-effort line-level regex matching per spec §6.4.
type Detector interface {
	Detect(prevSrc, newSrc string) []string
}

// FileChange is the per-file input the renderer takes after walking
// the commit's file list. NewlyAddedFile means the whole file was
// added in this commit (renderer prints the file even if no functions
// were detected).
type FileChange struct {
	RelativePath   string
	Language       string // "" => no detector for this extension
	PrevSource     string
	NewSource      string
	NewlyAddedFile bool
}
