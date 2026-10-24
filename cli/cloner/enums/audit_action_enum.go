package enums

// AuditActionType labels the planned action for a single record.
type AuditActionType byte

const (
	AuditActionInvalid AuditActionType = iota
	AuditActionClone
	AuditActionPull
	AuditActionCached
	AuditActionConflict
)

var actionStrings = [...]string{
	AuditActionInvalid:  "invalid",
	AuditActionClone:    "clone",
	AuditActionPull:     "pull",
	AuditActionCached:   "cached",
	AuditActionConflict: "conflict",
}

var actionMarkers = [...]string{
	AuditActionInvalid:  "!",
	AuditActionClone:    "+",
	AuditActionPull:     "~",
	AuditActionCached:   "=",
	AuditActionConflict: "?",
}

// IsValid reports whether the action is a valid executable or planned audit action.
func (a AuditActionType) IsValid() bool {
	return a >= AuditActionClone && a <= AuditActionConflict
}

// IsInvalid reports whether the action represents an invalid record.
func (a AuditActionType) IsInvalid() bool {
	return a == AuditActionInvalid
}

// String returns the canonical text label for the audit action.
func (a AuditActionType) String() string {
	if int(a) < len(actionStrings) {
		return actionStrings[a]
	}

	return "unknown"
}

// Marker returns the diff-style single-character marker for the audit action.
func (a AuditActionType) Marker() string {
	if int(a) < len(actionMarkers) {
		return actionMarkers[a]
	}

	return " "
}
