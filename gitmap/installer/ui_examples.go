// Package installer — ui_examples.go formats sample interactive composer configurations.
package installer

import "strings"

// PrintComposerExample returns an example JSON and CLI walkthrough for setting up Composer.
func PrintComposerExample() string {
	var b strings.Builder
	b.WriteString("Composer Installer Configuration Example\n")
	b.WriteString("=======================================\n\n")
	b.WriteString("JSON Definition:\n")
	b.WriteString("{\n")
	b.WriteString("  \"name\": \"PHP Composer\",\n")
	b.WriteString("  \"slug\": \"composer\",\n")
	b.WriteString("  \"target_os\": \"all\",\n")
	b.WriteString("  \"version\": \"v2.7.0\",\n")
	b.WriteString("  \"instructions\": \"php -r \\\"copy('https://getcomposer.org/installer', 'composer-setup.php');\\\"\"\n")
	b.WriteString("}\n")

	return b.String()
}
