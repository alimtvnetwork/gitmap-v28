package cmdpushfix

import (
	"errors"
	"flag"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

// RunPushFix is the main entry point for the gitmap push-fix command.
func RunPushFix(args []string) error {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			helptext.Print("push_fix")
			return nil
		}
	}

	opts, err := parsePushFixFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	return ExecutePushFix(opts)
}

func parsePushFixFlags(args []string) (PushFixOptions, error) {
	fs := flag.NewFlagSet("push-fix", flag.ContinueOnError)
	remoteFlag := fs.String("remote", "origin", "Upstream remote name")
	fs.StringVar(remoteFlag, "r", "origin", "Upstream remote name (short)")
	branchFlag := fs.String("branch", "", "Target branch name")
	fs.StringVar(branchFlag, "b", "", "Target branch name (short)")
	dryRunFlag := fs.Bool("dry-run", false, "Simulate push without modifying remote")
	fs.BoolVar(dryRunFlag, "n", false, "Simulate push (short)")
	forceFlag := fs.Bool("force", false, "Push with --force-with-lease")
	fs.BoolVar(forceFlag, "f", false, "Push with force (short)")
	yesFlag := fs.Bool("yes", false, "Bypass confirmation prompts")
	fs.BoolVar(yesFlag, "y", false, "Bypass confirmation (short)")
	sshFlag := fs.Bool("ssh", false, "Force remote URL conversion to SSH")
	httpsFlag := fs.Bool("https", false, "Force remote URL conversion to HTTPS")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return PushFixOptions{}, err
		}
		return PushFixOptions{}, apperror.WrapSimple(err, "cmdpushfix.parseFlags")
	}

	opts := PushFixOptions{
		Remote:   *remoteFlag,
		Branch:   *branchFlag,
		IsDryRun: *dryRunFlag,
		IsForce:  *forceFlag,
		IsYes:    *yesFlag,
		IsSSH:    *sshFlag,
		IsHTTPS:  *httpsFlag,
	}
	return opts, nil
}
