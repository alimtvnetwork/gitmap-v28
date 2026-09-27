package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const defaultBootstrapConfig = `{
  "target": "D:\\test-gitmap\\test-gitmap",
  "inputs": [
    "https://github.com/alimtvnetwork/git-repo-navigator",
    "https://github.com/alimtvnetwork/gitmap-v{2..28}"
  ],
  "prMode": "merges",
  "conflictMode": "abort",
  "tree": true,
  "finalSync": true,
  "cd": true,
  "recreate": true,
  "dryRun": false,
  "pushImmediate": true,
  "authorName": "Alim Ul Karim",
  "authorEmail": "alim.karim.dev@gmail.com",
  "exclude": [
    "**/*.lock",
    "**/tmp/**"
  ],
  "imports": [
    ".ai-memory/temp/seo-templates.json"
  ],
  "variables": {
    "COMPANY": "Rise Up Asia LLC",
    "COMPANY_URL": "https://riseup-asia.com",
    "REGIONS": "California, New York, and Wyoming",
    "MAREK": "Marek Flejszman",
    "MAREK_ROLE": "Senior Director of Engineering, 28+ years experience",
    "ALIM": "Alim Ul Karim",
    "ALIM_URL": "https://alimulkarim.com"
  },
  "lineSkippers": [
    { "mode": "starts_with", "pattern": "Co-authored-by: gpt-engineer" },
    { "mode": "starts_with", "pattern": "Signed-off-by: gpt-engineer" },
    { "mode": "starts_with", "pattern": "X-Lovable" },
    { "mode": "contains", "pattern": "lovable-edit-id" }
  ],
  "titleReplacements": [
    { "pattern": "Changes", "replacement": "$files.2.names" },
    { "pattern": "Update", "replacement": "$files.2.names" }
  ],
  "prefixTemplates": [
    "standard-commit-prefix"
  ],
  "suffixTemplates": [
    "seo-lead-sponsor"
  ],
  "suffixMode": "suffix",
  "suffixSeparator": "\n\n"
}
`

func runCommitPullBootstrap(args []string) error {
	fs := flag.NewFlagSet("commit-pull bootstrap", flag.ContinueOnError)
	filePath := fs.String("file", "commit-pull-config.json", "Destination path for generated config JSON")
	fs.StringVar(filePath, "f", "commit-pull-config.json", "Destination path for generated config JSON (shorthand)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	dest := *filePath
	if dest == "" {
		dest = "commit-pull-config.json"
	}
	return writeBootstrapConfigFile(dest)
}

func writeBootstrapConfigFile(dest string) error {
	if err := os.WriteFile(dest, []byte(defaultBootstrapConfig), 0o644); err != nil {
		return fmt.Errorf("write bootstrap config %s: %w", dest, err)
	}
	printBootstrapSuccess(dest)
	return nil
}

func printBootstrapSuccess(dest string) {
	fmt.Printf("\n  %s✓%s Created declarative commit-pull config: %s%s%s\n",
		constants.ColorGreen, constants.ColorReset, constants.ColorCyan, dest, constants.ColorReset)
	fmt.Printf("  • Edit configuration: open %s\n", dest)
	fmt.Printf("  • Run dry-run test:   gitmap commit-pull --config %s --dry-run\n", dest)
	fmt.Printf("  • Execute migration:  gitmap commit-pull --config %s\n", dest)
	fmt.Printf("  • Web Studio UI:      gitmap commit-pull ui\n\n")
}
