// Package cmdai — ai_flag.go: argv preprocessing for the --ai tracking flag.
//
// Moved from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D):
// the flag sets process-wide AI-tracking state, which is cmdai's concern.
package cmdai

import "os"

// StripAiFlag removes the --ai flag from args, enabling AI tracking via
// the GITMAP_AI_TRACKING environment variable as a side effect.
func StripAiFlag(args []string) []string {
	var cleaned []string
	for _, arg := range args {
		if arg == "--ai" {
			_ = os.Setenv("GITMAP_AI_TRACKING", "1")
			continue
		}
		cleaned = append(cleaned, arg)
	}
	return cleaned
}
