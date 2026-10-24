package cluster

import "github.com/alimtvnetwork/gitmap-v28/cli/cluster/enums"

// CommandType is the canonical broadcast command type, defined in enums.
type CommandType = enums.CommandType

const (
	// CommandPullAll is the command to tell clients to pull all local repositories.
	CommandPullAll = enums.CommandPullAll
)

// Event represents a broadcast event.
type Event struct {
	Type    CommandType
	Payload []byte
}

// ClientConnectioner abstracts the communication layer with a single client.
type (
	ClientConnectioner interface {
		SendEvent(event Event) error
	}

	ClientConnection = ClientConnectioner
)

// BroadcastEvent sends an event to all connected clients.
func BroadcastEvent(clients []ClientConnectioner, event Event) []error {
	var errors []error

	for _, client := range clients {
		err := client.SendEvent(event)
		hasError := err != nil
		if hasError {
			errors = append(errors, err)
		}
	}

	return errors
}
