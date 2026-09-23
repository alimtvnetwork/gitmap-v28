package cmdos

import (
	"os"

	"github.com/spf13/cobra"
)

var osInfoJSONFlag bool

// OSInfoCmd represents the 'gitmap os-info' command.
var OSInfoCmd = &cobra.Command{
	Use:     "os-info [flags]",
	Aliases: []string{"osinfo", "sysinfo", "system-info"},
	Short:   "Display operating system, distribution version, and hardware architecture",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunOSInfoCLI(args)
	},
}

func init() {
	OSInfoCmd.Flags().BoolVar(&osInfoJSONFlag, "json", false, "Output system information in JSON format")
	osCmd.AddCommand(OSInfoCmd)
}

// RunOSInfoCLI gathers and renders OS details.
func RunOSInfoCLI(args []string) error {
	isJSON := osInfoJSONFlag
	for _, a := range args {
		if a == "--json" {
			isJSON = true
		}
	}

	report := GetLocalOSInfo()
	return RenderOSInfo(os.Stdout, report, isJSON)
}
