package jsonenvelope

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// EnvelopeAttributes holds metadata describing the payload, origin, and generation process.
type EnvelopeAttributes struct {
	Type      string `json:"type"`
	Source    string `json:"source,omitempty"`
	How       string `json:"how,omitempty"`
	Version   string `json:"version,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// Envelope is a generic envelope wrapping attributes and typed data.
type Envelope[T any] struct {
	Attributes EnvelopeAttributes `json:"attributes"`
	Data       T                  `json:"data"`
}

// RawEnvelope wraps attributes with unparsed json.RawMessage data.
type RawEnvelope struct {
	Attributes EnvelopeAttributes `json:"attributes"`
	Data       json.RawMessage    `json:"data"`
}

// NewEnvelope constructs a typed envelope with the given metadata and payload.
func NewEnvelope[T any](
	dataType string,
	source string,
	how string,
	version string,
	payload T,
) Envelope[T] {
	curTime := time.Now().UTC().Format(time.RFC3339)
	ver := version
	if ver == "" {
		ver = "1.0"
	}
	return Envelope[T]{
		Attributes: EnvelopeAttributes{
			Type:      dataType,
			Source:    source,
			How:       how,
			Version:   ver,
			Timestamp: curTime,
		},
		Data: payload,
	}
}

// IsEnvelope returns true if the raw JSON conforms to the {attributes, data} envelope schema.
func IsEnvelope(raw []byte) bool {
	var probe struct {
		Attributes *struct {
			Type string `json:"type"`
		} `json:"attributes"`
		Data json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(raw, &probe)
	if err != nil {
		return false
	}
	return probe.Attributes != nil && len(probe.Data) > 0
}

// ExtractPayload extracts the core payload bytes and attributes from either an envelope or legacy JSON.
func ExtractPayload(raw []byte) ([]byte, EnvelopeAttributes, error) {
	clean := strings.TrimSpace(string(raw))
	if len(clean) == 0 {
		return nil, EnvelopeAttributes{}, apperror.NewValidationError("ExtractPayload: empty input JSON")
	}

	if !IsEnvelope(raw) {
		attrs := EnvelopeAttributes{
			Version: "legacy",
		}
		return raw, attrs, nil
	}

	return extractEnvelopePayload(raw)
}

func extractEnvelopePayload(raw []byte) ([]byte, EnvelopeAttributes, error) {
	var env RawEnvelope
	err := json.Unmarshal(raw, &env)
	if err != nil {
		return nil, EnvelopeAttributes{}, apperror.WrapSimple(err, "ExtractPayload.UnmarshalEnvelope")
	}
	return env.Data, env.Attributes, nil
}
