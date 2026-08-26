// Package gitutil — conflict_evaluator.go checks merge collision probability.
package gitutil

func HasPotentialConflictRisk(d DirtyDiagnosis) bool {
	return d.ModifiedCount > 0 || d.DeletedCount > 0
}
