// Package installer — prompt_output_stream.go captures output bytes.
package installer

import "bytes"

type PromptOutputBuffer struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}
