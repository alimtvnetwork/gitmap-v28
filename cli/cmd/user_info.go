package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

// runUserInfo executes the user info inspection command.
func runUserInfo(args []string) error {
	return executeUserInfo(os.Stdout, args)
}

func executeUserInfo(w io.Writer, args []string) error {
	summary := usercontext.CollectUserSummary()
	if hasUserInfoJSONFlag(args) {
		return emitUserInfoJSON(w, summary)
	}

	return usercontext.RenderUserInfoCard(w, summary)
}

func hasUserInfoJSONFlag(args []string) bool {
	return hasUserFlag(args, "--json") || hasUserFlag(args, "-j")
}

func emitUserInfoJSON(w io.Writer, summary usercontext.UserSummary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(data))

	return nil
}
