package cluster

// CommandType represents the type of command broadcasted to clients.
type CommandType string

const (
	// CommandPullAll is the command to tell clients to pull all local repositories.
	CommandPullAll CommandType = "PullAll"
)

// Event represents a broadcast event.
type Event struct {
	Type    CommandType
	Payload []byte
}

// ClientConnection abstracts the communication layer with a single client.
type ClientConnection interface {
	SendEvent(event Event) error
}

// BroadcastEvent sends an event to all connected clients.
func BroadcastEvent(clients []ClientConnection, event Event) []error {
	var errors []error

	for _, client := range clients {
		err := client.SendEvent(event)
		hasError := err != nil
		if hasError == true {
			errors = append(errors, err)
		}
	}

	return errors
}
