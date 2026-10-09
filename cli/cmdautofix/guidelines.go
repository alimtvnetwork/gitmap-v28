package cmdautofix

// guidelinesComposite ports 05-guideline-autofixer.py: it runs the newlines
// fix, then the naming check (D12). Unlike script 05's inverted
// default-fix behavior, this composite obeys the unified dry-run/--apply
// rule: dry-run only checks, --apply fixes newlines but naming stays
// report-only (its Fix is nil).
func guidelinesCheck(relPath string, src []byte, opts Options) []Violation {
	violations := []Violation{}
	violations = append(violations, newlinesCheck(relPath, src, opts)...)
	// Naming applies to code files only, same as its own category scope.
	if isCodeFile(relPath) {
		violations = append(violations, namingCheck(relPath, src, opts)...)
	}
	return violations
}

func guidelinesFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	fixed, fixVios := newlinesFix(relPath, src, opts)
	// Naming hits were already recorded by guidelinesCheck; Fix reports
	// only the newlines rewrite outcome.
	return fixed, fixVios
}
