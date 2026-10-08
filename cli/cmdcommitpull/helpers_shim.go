package cmdcommitpull

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitin"
)

func runCommitPullUI(args []string) error {
	fs := flag.NewFlagSet("commit-pull ui", flag.ContinueOnError)
	port := fs.Int("port", 8925, "Port for web UI")
	noOpen := fs.Bool("no-open", false, "Do not open browser automatically")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ln, err := cmdcommitin.ListenCommitPullUIPort(*port)
	if err != nil {
		return fmt.Errorf("listen ui: %w", err)
	}
	return cmdcommitin.StartCommitPullUIServer(ln, *noOpen)
}
