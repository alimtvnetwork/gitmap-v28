package cmdssh

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// GetFleetNodeCompletions returns dynamic completions for cluster nodes.
func GetFleetNodeCompletions() []string {
	conns, _, err := loadFleetInventory()
	if err != nil || len(conns) == 0 {
		return []string{"all\tAll cluster nodes"}
	}

	results := make([]string, 0, len(conns)+1)
	results = append(results, "all\tDeploy or delegate to all cluster nodes")
	for i, c := range conns {
		desc := fmt.Sprintf("#%d %s (%s)", i+1, c.IPAddress, c.OS)
		if c.Alias != "" {
			results = append(results, fmt.Sprintf("%s\t%s", c.Alias, desc))
		}
		if c.IPAddress != "" && !strings.EqualFold(c.Alias, c.IPAddress) {
			results = append(results, fmt.Sprintf("%s\t#%d %s", c.IPAddress, i+1, c.Alias))
		}
	}
	return results
}

// DeployValidArgsFunction provides dynamic completions for deploy commands.
func DeployValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return GetFleetNodeCompletions(), cobra.ShellCompDirectiveNoFileComp
	}
	if len(args) == 1 {
		// Second argument is local source file or directory path
		return nil, cobra.ShellCompDirectiveDefault
	}
	if len(args) == 2 {
		return []string{
			"D:/work\tDefault Windows work folder",
			"~/work\tDefault Linux work folder",
		}, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// SSHValidArgsFunction provides dynamic node completions for ssh commands.
func SSHValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return GetFleetNodeCompletions(), cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}
