import re

with open('cli/cmdagy/agy_inject_prompts.go', 'r') as f:
    content = f.read()

# Remove ipt and ipr from AgyInjectPromptsCmd aliases to avoid conflicts
content = content.replace(
    'Aliases: []string{"ip", "ipt", "ipr"},',
    'Aliases: []string{"ip", "inject-prompts-ssh"},'
)

# Also, the user asked for:
# gitmap agy inject-prompts(ip) ... [--prefix/pfx ui-ux/--ui-ux/--uu] --suffix
content = content.replace(
    'AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "prefix", "", "Prefix template")',
    'AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "prefix", "", "Prefix template")\n\tAgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "pfx", "", "Prefix template (alias)")\n\tAgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsWatch, "ui-ux", false, "Use UI/UX template")\n\tAgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsWatch, "uu", false, "Use UI/UX template (alias)")\n\tAgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "suffix", "", "Suffix template")\n\tAgyInjectPromptsCmd.Flags().StringVar(&injectPromptsSSH, "nodes", "", "Nodes")'
)

with open('cli/cmdagy/agy_inject_prompts.go', 'w') as f:
    f.write(content)

