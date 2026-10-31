package cmdcdp

import "encoding/json"

// CDPTarget describes a Chrome DevTools Protocol debug target.
type CDPTarget struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	Type                string `json:"type"`
	URL                 string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	DevtoolsFrontendURL string `json:"devtoolsFrontendUrl"`
}

// CDPRequest represents a JSON-RPC command sent to a CDP target.
type CDPRequest struct {
	ID     int            `json:"id"`
	Method string         `json:"method"`
	Params map[string]any `json:"params,omitempty"`
}

// CDPResponse represents a response received from a CDP target.
type CDPResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *CDPError       `json:"error,omitempty"`
}

// CDPError captures an error returned by the Chrome DevTools Protocol.
type CDPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ProbeHit represents an individual DOM element identified by the selector ladder.
type ProbeHit struct {
	Selector    string `json:"selector"`
	Tag         string `json:"tag"`
	AriaLabel   string `json:"ariaLabel"`
	Placeholder string `json:"placeholder"`
	Text        string `json:"text"`
}

// ContentEditableItem represents a contenteditable or rich-text DOM container.
type ContentEditableItem struct {
	Tag       string `json:"tag"`
	Aria      string `json:"aria"`
	ClassName string `json:"className"`
}

// ProbeResult stores the structured output of the DOM input probe execution.
type ProbeResult struct {
	Found               []ProbeHit            `json:"found"`
	Active              string                `json:"active"`
	AllContentEditables []ContentEditableItem `json:"allContentEditables"`
}

// CDPOptions holds parsed command-line flags for CDP operations.
type CDPOptions struct {
	Port       int
	TargetID   string
	IsJSON     bool
	TimeoutSec int
}

// runtimeEvalResult unpacks the outer envelope of a Runtime.evaluate response.
type runtimeEvalResult struct {
	Result struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"result"`
	ExceptionDetails any `json:"exceptionDetails,omitempty"`
}
