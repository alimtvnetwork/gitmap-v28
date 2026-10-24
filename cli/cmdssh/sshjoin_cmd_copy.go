package cmdssh

import (
	"github.com/spf13/cobra"
)

var SJCopyIdCmd = &cobra.Command{
	Use:   "copy-id [target] [flags]",
	Short: "Copy SSH public key to target machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHAuthKeyDeployCLI(args)
	},
}

var SJFixAuthCmd = &cobra.Command{
	Use:   "fix-auth [target] [flags]",
	Short: "Deploy SSH public key to target machines as authorized key",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHAuthKeyDeployCLI(args)
	},
}

func init() {
	SSHJoinCmd.AddCommand(SJAddCmd)
	SSHJoinCmd.AddCommand(SJAddWithPassCmd)
	SSHJoinCmd.AddCommand(SJRmCmd)
	SSHJoinCmd.AddCommand(SJAddAuthCmd)
	SSHJoinCmd.AddCommand(SJLsCmd)
	SSHJoinCmd.AddCommand(SJNodesCmd)
	SSHJoinCmd.AddCommand(SJHistCmd)
	SSHJoinCmd.AddCommand(SJClusterImportCmd)
	SSHJoinCmd.AddCommand(ClusterBootstrapCmd)
	SSHJoinCmd.AddCommand(ClusterInstallCmd)
	SSHJoinCmd.AddCommand(SJAuthKeyCmd)
	SSHJoinCmd.AddCommand(SJCopyIdCmd)
	SSHJoinCmd.AddCommand(SJFixAuthCmd)
}
