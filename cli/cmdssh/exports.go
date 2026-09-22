package cmdssh

import (
	"context"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/spf13/cobra"
)

// JoinRunner and ProfileRunner are injected callbacks from root orchestrator.
var (
	JoinRunner    func([]string) error
	ProfileRunner func([]string) error
)

type SEOptions = seOptions
type SJOptions = sjOptions

// RunSSH executes ssh subcommands.
func RunSSH(args []string) error {
	return runSSH(args)
}

// RunSSHExec executes ssh exec subcommand.
func RunSSHExec(args []string) error {
	return runSSHExec(args)
}

// RunSSHCopyCLI executes ssh copy subcommand.
func RunSSHCopyCLI(args []string) error {
	return runSSHCopyCLI(args)
}

// RunSSHMvCLI executes ssh mv subcommand.
func RunSSHMvCLI(args []string) error {
	return runSSHMvCLI(args)
}

// RunSSHBind executes ssh bind subcommand.
func RunSSHBind(args []string) error {
	return runSSHBind(args)
}

// RunSJAddAuth executes ssh-join add-auth subcommand.
func RunSJAddAuth(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSJAddAuth(cmd, args, ctx)
}

// RunSJHistory executes ssh-join history subcommand.
func RunSJHistory(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSJHistory(cmd, args, ctx)
}

// RunSJLs executes ssh-join ls subcommand.
func RunSJLs(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSJLs(cmd, args, ctx)
}

// RunSJRm executes ssh-join rm subcommand.
func RunSJRm(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSJRm(cmd, args, ctx)
}

// RunSSHJoin executes ssh join command.
func RunSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSSHJoin(cmd, args, ctx)
}

// RunSSHJoinCLI executes unified ssh join CLI entrypoint.
func RunSSHJoinCLI(args []string) error {
	return runSSHJoinCLI(args)
}

// RunSSHErrorLogsCLI exports SSH error logs reader.
func RunSSHErrorLogsCLI(args []string) error {
	return runSSHErrorLogsCLI(args)
}

// ParseSEFlags parses ssh exec flags.
func ParseSEFlags(args []string) SEOptions {
	return parseSEFlags(args)
}

// ParseSJFlags parses ssh join flags.
func ParseSJFlags(args []string) SJOptions {
	return parseSJFlags(args)
}

// EnsureSSHDir ensures directory exists.
func EnsureSSHDir(dir string) error {
	return ensureSSHDir(dir)
}

// CopyPubKeyAndAnnounce copies public key to clipboard and prints advice.
func CopyPubKeyAndAnnounce(pub string) {
	copyPubKeyAndAnnounce(pub)
}

// ResolveGitEmail resolves git email.
func ResolveGitEmail() string {
	return resolveGitEmail()
}

// ValidateSSHKeygen validates keygen tool.
func ValidateSSHKeygen() error {
	return validateSSHKeygen()
}

// EncodeSSHListJSON encodes keys to JSON.
func EncodeSSHListJSON(w io.Writer, keys []model.SSHKey) error {
	return encodeSSHListJSON(w, keys)
}

// RunClusterAddCLI executes cluster add command.
func RunClusterAddCLI(args []string) error {
	if hasHelpFlag(args) {
		helptext.Print("cluster-add")
		return nil
	}
	if len(args) == 0 {
		helptext.Print("cluster-add")
		return apperror.NewValidationError("missing target host. Usage: gitmap cluster add <user@ip|ip> [alias] [flags]")
	}
	return RunSSHJoinCLI(append([]string{"add"}, args...))
}

// RunClusterJoinCLI executes cluster join command.
func RunClusterJoinCLI(args []string) error {
	if hasHelpFlag(args) {
		helptext.Print("cluster-join")
		return nil
	}
	if len(args) == 0 {
		helptext.Print("cluster-join")
		return apperror.NewValidationError("missing target host. Usage: gitmap cluster join <user@ip|ip> [alias] [flags]")
	}
	return RunSSHJoinCLI(args)
}

// IsInteractiveMacroAdd checks if arguments represent an interactive macro add with zero commands.
func IsInteractiveMacroAdd(args []string) bool {
	return isInteractiveMacroAdd(args)
}

// PrintInteractiveMacroAdvice prints guidance for interactive macro over SSH.
func PrintInteractiveMacroAdvice() {
	printInteractiveMacroAdvice()
}
