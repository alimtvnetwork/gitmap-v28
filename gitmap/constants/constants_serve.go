package constants

// Serve CLI commands.
// gitmap:cmd top-level
const (
	CmdServe      = "serve"
	CmdServeAlias = "sv"
)

// Serve help text.
const HelpServe = "  serve (sv)          Start the orchestrator daemon and generate a join token"

// Serve flag descriptions.
const (
	FlagDescServePort = "Port to bind the server to (default: 9999)"
)

// Serve terminal messages.
const (
	MsgServeStarting = "Starting orchestrator daemon..."
	MsgServeAddress  = "Server listening on %s:%d"
	MsgServeToken    = "Join Token: %s"
	MsgServeShutdown = "\nShutting down orchestrator daemon..."
)

// Serve error messages.
const (
	ErrServeBind          = "Failed to bind server to network interface: %v"
	ErrServeTokenGenerate = "Failed to generate join token: %v\n"
)

// Serve configuration defaults.
const (
	ServeDefaultPort = 9999
	ServeBindAddress = "0.0.0.0"
	ServeProtocol    = "tcp"
	ServeTokenBytes  = 32
)

// Serve flag names.
const (
	FlagServePort = "port"
)
