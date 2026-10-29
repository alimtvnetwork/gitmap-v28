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

// InstructionLink defines a single URL target and its notification configuration.
type InstructionLink struct {
	URL          string   `json:"url"`
	Recipients   []string `json:"recipients,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	BodyTemplate string   `json:"bodyTemplate,omitempty"`
}

// InstructionSafety defines sandbox constraints and resource limits.
type InstructionSafety struct {
	AllowedProtocols []string `json:"allowedProtocols"`
	MaxBodyBytes     int64    `json:"maxBodyBytes"`
	TimeoutSeconds   int      `json:"timeoutSeconds"`
	IsNoJavaScript   bool     `json:"noJavaScript"`
}

// InstructionConfig defines the root configuration model parsed from JSON input.
type InstructionConfig struct {
	Links      []InstructionLink `json:"links"`
	MaxLevels  int               `json:"maxLevels"`
	IsInfinite bool              `json:"infinite"`
	Safety     InstructionSafety `json:"safety"`
	IsAutoSend bool              `json:"autoSend"`
}

// InstructionLinkReport captures the evaluation outcome for an individual link instruction.
type InstructionLinkReport struct {
	OriginalURL         string   `json:"originalUrl"`
	ExpandedURL         string   `json:"expandedUrl"`
	IsExpanded          bool     `json:"isExpanded"`
	StatusCode          int      `json:"statusCode"`
	RedirectHops        int      `json:"redirectHops"`
	HopTrace            []string `json:"hopTrace,omitempty"`
	IsSent              bool     `json:"isSent"`
	IsEmailSent         bool     `json:"isEmailSent,omitempty"`
	SentRecipients      []string `json:"sentRecipients,omitempty"`
	Error               string   `json:"error,omitempty"`
	InterpolatedSubject string   `json:"interpolatedSubject,omitempty"`
	InterpolatedBody    string   `json:"interpolatedBody,omitempty"`
}

// InstructionResult represents the complete top-level output envelope.
type InstructionResult struct {
	TotalLinks    int                     `json:"totalLinks"`
	ExpandedCount int                     `json:"expandedCount"`
	SentCount     int                     `json:"sentCount"`
	HasErrors     bool                    `json:"hasErrors"`
	Results       []InstructionLinkReport `json:"results"`
}

// DefaultInstructionSafety returns hardened default security parameters.
func DefaultInstructionSafety() InstructionSafety {
	return InstructionSafety{
		AllowedProtocols: []string{"http", "https"},
		MaxBodyBytes:     4096,
		TimeoutSeconds:   10,
		IsNoJavaScript:   true,
	}
}

// DefaultInstructionConfig constructs an InstructionConfig populated with safe defaults.
func DefaultInstructionConfig() InstructionConfig {
	return InstructionConfig{
		MaxLevels:  10,
		IsInfinite: false,
		Safety:     DefaultInstructionSafety(),
		IsAutoSend: false,
	}
}
