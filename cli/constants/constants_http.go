package constants

// HTTP header and media-type constants shared by JSON API handlers,
// HTTP clients, and mail/schedule senders. Centralized so header
// names and media types cannot drift between handlers, clients, and
// embedded UI scripts.
const (
	// HeaderContentType is the HTTP Content-Type header name.
	HeaderContentType = "Content-Type"
	// MediaTypeJSON is the JSON media type (without parameters).
	MediaTypeJSON = "application/json"
	// AuthBearerPrefix is the Authorization bearer-scheme prefix,
	// including the trailing space.
	AuthBearerPrefix = "Bearer "
)
