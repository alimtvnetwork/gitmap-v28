package constants

// Join CLI commands.
// gitmap:cmd top-level
const (
	CmdJoin      = "join"
	CmdJoinAlias = "jn"
)

// Join help text.
const HelpJoin = "  join (jn)           Join an existing orchestrator daemon"

// Join terminal messages.
const (
	MsgJoinStarting = "Joining cluster at %s..."
	MsgJoinSuccess  = "Successfully joined the cluster."
)

// Join error messages.
const (
	ErrJoinMissingAddress = "Missing server address. Usage: gitmap join <IP>:<PORT> --token <TOKEN>"
	ErrJoinMissingToken   = "Missing join token. Usage: gitmap join <IP>:<PORT> --token <TOKEN>"
	ErrJoinFailed         = "Failed to join cluster: %v\n"
)

// Join flag names.
const (
	FlagJoinToken = "token"
)

// Join flag descriptions.
const (
	FlagDescJoinToken = "Token required to join the cluster"
)
