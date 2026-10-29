package linkexpander

// ExpandedLink holds the outcome of expanding one URL.
type ExpandedLink struct {
	Original   string `json:"original"`
	Expanded   string `json:"expanded"`
	IsExpanded bool   `json:"isExpanded"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
}

// ExpandResult represents the overall link expansion for a piece of text.
type ExpandResult struct {
	OriginalText string         `json:"originalText"`
	ExpandedText string         `json:"expandedText"`
	HasChanges   bool           `json:"hasChanges"`
	Count        int            `json:"count"`
	Links        []ExpandedLink `json:"links"`
}

// SendEmailJSONResult holds the structured JSON output for send-email.
type SendEmailJSONResult struct {
	IsSuccess     bool           `json:"isSuccess"`
	From          string         `json:"from,omitempty"`
	To            []string       `json:"to"`
	Subject       string         `json:"subject"`
	IsHTML        bool           `json:"isHTML"`
	IsExpanded    bool           `json:"isExpanded"`
	ExpandedCount int            `json:"expandedCount"`
	Links         []ExpandedLink `json:"links,omitempty"`
	Error         string         `json:"error,omitempty"`
}
