import os
import textwrap

CLI_DIR = "d:/work/gitmap/cli/cmdagy"
os.makedirs(CLI_DIR, exist_ok=True)

files = {
    "agy_lapp.go": """package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)
var lappSSH bool
var AgyLappCmd = &cobra.Command{
	Use:     "look-all-projects-prompts",
	Aliases: []string{"lapp"},
	Short:   "Look at all projects prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Looking at all projects prompts...")
		return nil
	},
}
func init() {
	AgyLappCmd.Flags().BoolVar(&lappSSH, "ssh", false, "Use SSH to aggregate")
	AgyCmd.AddCommand(AgyLappCmd)
}
""",
    "agy_wapp.go": """package cmdagy
import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
)
var wappSSH bool
var AgyWappCmd = &cobra.Command{
	Use:     "watch-all-projects-prompts",
	Aliases: []string{"wapp"},
	Short:   "Watch all projects prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Watching all projects prompts...")
		for {
			time.Sleep(30 * time.Second)
		}
	},
}
func init() {
	AgyWappCmd.Flags().BoolVar(&wappSSH, "ssh", false, "Use SSH to aggregate")
	AgyCmd.AddCommand(AgyWappCmd)
}
""",
    "agy_inject_txt.go": """package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)
var iptPromptName string
var iptWatch bool
var iptRerun int
var iptPrefix string
var iptSuffix string

var AgyInjectTxtCmd = &cobra.Command{
	Use:     "inject-prompts-txt <prompt-text>",
	Aliases: []string{"ipt"},
	Short:   "Inject raw prompt text",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Injecting text prompt...")
		return nil
	},
}
func init() {
	AgyInjectTxtCmd.Flags().StringVarP(&iptPromptName, "prompt-name", "p", "", "Prompt name")
	AgyInjectTxtCmd.Flags().BoolVar(&iptWatch, "watch", false, "Watch until prompt completes")
	AgyInjectTxtCmd.Flags().IntVar(&iptRerun, "rerun", 1, "Rerun count")
	AgyInjectTxtCmd.Flags().StringVar(&iptPrefix, "prefix", "", "Prefix template")
	AgyInjectTxtCmd.Flags().StringVar(&iptSuffix, "suffix", "", "Suffix template")
	AgyCmd.AddCommand(AgyInjectTxtCmd)
}
""",
    "agy_inject_rerun.go": """package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)
var iprPromptName string
var iprWatch bool
var iprPrefix string
var iprSuffix string
var iprNodes string
var iprN int

var AgyInjectRerunCmd = &cobra.Command{
	Use:     "inject-prompts-rerun <prompt-text>",
	Aliases: []string{"ipr"},
	Short:   "Inject prompt and rerun",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Injecting and rerunning...")
		return nil
	},
}

var AgyInjectRerunSSHCmd = &cobra.Command{
	Use:     "inject-prompts-rerun-ssh <prompt-text>",
	Aliases: []string{"iprs"},
	Short:   "Inject prompt and rerun via SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Injecting and rerunning over SSH...")
		return nil
	},
}

func init() {
	AgyInjectRerunCmd.Flags().StringVarP(&iprPromptName, "prompt-name", "p", "", "Prompt name")
	AgyInjectRerunCmd.Flags().BoolVar(&iprWatch, "watch", false, "Watch until prompt completes")
	AgyInjectRerunCmd.Flags().IntVarP(&iprN, "number", "n", 2, "Rerun count")
	AgyInjectRerunCmd.Flags().StringVar(&iprPrefix, "prefix", "", "Prefix template")
	AgyInjectRerunCmd.Flags().StringVar(&iprSuffix, "suffix", "", "Suffix template")
	AgyInjectRerunCmd.Flags().StringVar(&iprNodes, "ssh-only-node", "", "Specific SSH node")
	AgyCmd.AddCommand(AgyInjectRerunCmd)

	AgyInjectRerunSSHCmd.Flags().StringVarP(&iprPromptName, "prompt-name", "p", "", "Prompt name")
	AgyInjectRerunSSHCmd.Flags().BoolVar(&iprWatch, "watch", false, "Watch until prompt completes")
	AgyInjectRerunSSHCmd.Flags().IntVarP(&iprN, "number", "n", 2, "Rerun count")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprPrefix, "prefix", "", "Prefix template")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprSuffix, "suffix", "", "Suffix template")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprNodes, "nodes", "", "Specific SSH nodes")
	AgyCmd.AddCommand(AgyInjectRerunSSHCmd)
}
""",
    "agy_ssh_cmds.go": """package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)

var AgyLookProjectsSSHCmd = &cobra.Command{
	Use:   "look-projects-ssh",
	Short: "Table of all VM projects via SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching projects over SSH...")
		return nil
	},
}

var AgyRestEnableCmd = &cobra.Command{
	Use:   "rest-enable",
	Short: "Enable REST OS service",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Enabling REST service...")
		return nil
	},
}

func init() {
	AgyCmd.AddCommand(AgyLookProjectsSSHCmd)
	AgyCmd.AddCommand(AgyRestEnableCmd)
}
"""
}

for filename, content in files.items():
    path = os.path.join(CLI_DIR, filename)
    with open(path, "w") as f:
        f.write(content)

print("Scaffolded all missing CLI commands.")
