package constants

// AsyncActionType defines action verbs for the async command.
type AsyncActionType string

// Backwards-compatible type alias.
type AsyncAction = AsyncActionType

const (
	AsyncActionInvalid AsyncActionType = ""
	AsyncActionLs      AsyncActionType = "ls"
	AsyncActionList    AsyncActionType = "list"
	AsyncActionStop    AsyncActionType = "stop"
	AsyncActionKill    AsyncActionType = "kill"
)

// String returns the canonical serialized string representation.
func (a AsyncActionType) String() string {
	return string(a)
}

// Label returns the human-readable display label.
func (a AsyncActionType) Label() string {
	return a.String()
}

// IsValid reports whether the async action is valid.
func (a AsyncActionType) IsValid() bool {
	hasAction := a.IsList() || a.IsStop()

	return hasAction
}

// IsList reports whether the action indicates listing background jobs.
func (a AsyncActionType) IsList() bool {
	return a == AsyncActionLs || a == AsyncActionList
}

// IsStop reports whether the action indicates stopping or killing a job.
func (a AsyncActionType) IsStop() bool {
	return a == AsyncActionStop || a == AsyncActionKill
}
