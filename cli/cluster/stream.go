package cluster

// ProgressLog represents a single log or progress update from a client.
type ProgressLog struct {
	ClientID string
	Message  string
	Progress int
	IsDone   bool
}

// StreamHandler handles incoming progress logs from clients.
type StreamHandler interface {
	OnLogReceived(log ProgressLog)
}

// LogStreamServer represents a server capable of receiving streamed logs.
// This provides a skeleton for RPC, WebSockets, or SSE integrations.
type LogStreamServer struct {
	Handler StreamHandler
}

// NewLogStreamServer creates a new log stream server.
func NewLogStreamServer(handler StreamHandler) *LogStreamServer {
	return &LogStreamServer{
		Handler: handler,
	}
}

// ReceiveLog acts as the endpoint where clients push their logs.
func (s *LogStreamServer) ReceiveLog(log ProgressLog) {
	hasHandler := s.Handler != nil
	if hasHandler {
		s.Handler.OnLogReceived(log)
	}
}
