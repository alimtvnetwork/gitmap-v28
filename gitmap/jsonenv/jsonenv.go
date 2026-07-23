// Package jsonenv provides a uniform `--json` envelope for every
// gitmap command that supports machine-readable output (#12).
//
// Envelope shape:
//
//	{
//	  "schema":  "gitmap.v1",
//	  "version": "6.61.0",
//	  "command": "scan",
//	  "ok":      true,
//	  "data":    <command-specific payload>,
//	  "error":   "..."   // present only when ok=false
//	}
//
// Existing per-command JSON schemas under spec/08-json-schemas/ are
// embedded as the `data` payload — no breaking change to the inner
// shape, only an outer wrapper so tooling can detect schema version
// + dispatch by command without sniffing keys.
package jsonenv

import (
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Schema identifier surfaced in every envelope. Bumped on
// breaking changes to the wrapper (NOT to inner payloads).
const Schema = "gitmap.v1"

// Envelope is the wire shape emitted by Write.
type Envelope struct {
	Schema  string      `json:"schema"`
	Version string      `json:"version"`
	Command string      `json:"command"`
	OK      bool        `json:"ok"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// WriteOK emits a success envelope to w with the supplied payload.
func WriteOK(w io.Writer, command string, data interface{}) error {
	return write(w, Envelope{
		Schema: Schema, Version: constants.Version, Command: command, OK: true, Data: data,
	})
}

// WriteErr emits a failure envelope. payload may be nil.
func WriteErr(w io.Writer, command string, errMsg string, data interface{}) error {
	return write(w, Envelope{
		Schema: Schema, Version: constants.Version, Command: command, OK: false, Error: errMsg, Data: data,
	})
}

func write(w io.Writer, env Envelope) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", constants.JSONIndent)
	return enc.Encode(env)
}
