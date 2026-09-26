import re

with open('cli/cmdagy/agy_look_prompts.go', 'r') as f:
    content = f.read()

content = content.replace(
    'AgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "compact", 0, "Compact words")',
    'AgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "compact", 0, "Compact words")\n\tAgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "words", 0, "Words (alias for compact)")'
)

with open('cli/cmdagy/agy_look_prompts.go', 'w') as f:
    f.write(content)

with open('cli/cmdagy/agy_watch_prompts.go', 'r') as f:
    content = f.read()

content = content.replace(
    'AgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "compact", 0, "Compact words")',
    'AgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "compact", 0, "Compact words")\n\tAgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "words", 0, "Words (alias for compact)")'
)

with open('cli/cmdagy/agy_watch_prompts.go', 'w') as f:
    f.write(content)
