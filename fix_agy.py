import re

with open("gitmap/cmd/agy_cmd.go", "r", encoding="utf-8") as f:
    content = f.read()

# I will append the new commands to the init() function.
init_addition = """
	AgyCmd.AddCommand(agyClearCmd)
	AgyCmd.AddCommand(agyOpenCmd)
	AgyCmd.AddCommand(agyPromptCmd)
	AgyCmd.AddCommand(agyRwCmd)
	AgyCmd.AddCommand(agySyncCmd)
	AgyCmd.AddCommand(agyPapCmd)
	AgyCmd.AddCommand(agyExportCmd)
	AgyCmd.AddCommand(agyImportCmd)
	AgyCmd.AddCommand(agyPluginsCmd)
"""
content = content.replace("	AgyCmd.AddCommand(agyUpdateCmd)\n}", "	AgyCmd.AddCommand(agyUpdateCmd)" + init_addition + "}\n")

# Now I'll append the new Cobra commands at the end of the file.
new_commands = """
var agyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear Antigravity cache and projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [clear] is not yet implemented")
		return nil
	},
}

var agyOpenCmd = &cobra.Command{
	Use:   "open [slug or path]",
	Short: "Open Antigravity or a specific project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [open] is not yet implemented")
		return nil
	},
}

var agyPromptCmd = &cobra.Command{
	Use:   "prompt [project slug/id/path] [prompt or text file]",
	Short: "Send a prompt to Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [prompt] is not yet implemented")
		return nil
	},
}

var agyRwCmd = &cobra.Command{
	Use:   "rw [path or project slug]",
	Short: "Enable rewrite both for project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [rw] is not yet implemented")
		return nil
	},
}

var agySyncCmd = &cobra.Command{
	Use:   "sync [path or dir]",
	Short: "Load and sync all projects to Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [sync] is not yet implemented")
		return nil
	},
}

var agyPapCmd = &cobra.Command{
	Use:     "prompt-all-project [title] [prompt or text file]",
	Aliases: []string{"pap"},
	Short:   "Send prompt to all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [prompt-all-project] is not yet implemented")
		return nil
	},
}

var agyExportCmd = &cobra.Command{
	Use:     "export-projects [file or path]",
	Aliases: []string{"ep"},
	Short:   "Create a zip backup of Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [export-projects] is not yet implemented")
		return nil
	},
}

var agyImportCmd = &cobra.Command{
	Use:     "import-projects [file or path]",
	Aliases: []string{"ip"},
	Short:   "Import a zip backup of Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [import-projects] is not yet implemented")
		return nil
	},
}

var agyPluginsCmd = &cobra.Command{
	Use:   "plugin",
	Aliases: []string{"plugins"},
	Short: "Manage Antigravity plugins",
}

var agyPluginLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed and installable plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [plugin ls] is not yet implemented")
		return nil
	},
}

var agyPluginInstallCmd = &cobra.Command{
	Use:   "install [slug]",
	Short: "Install an Antigravity plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [plugin install] is not yet implemented")
		return nil
	},
}

func initPlugins() {
	agyPluginsCmd.AddCommand(agyPluginLsCmd)
	agyPluginsCmd.AddCommand(agyPluginInstallCmd)
}
"""

content += new_commands

# We also need to make sure initPlugins() is called in the main init()
content = content.replace("	AgyCmd.AddCommand(agyPluginsCmd)\n}", "	AgyCmd.AddCommand(agyPluginsCmd)\n\tinitPlugins()\n}\n")

# And add aliases to agyAddCmd ("add-project"), agyRmCmd ("remove")
content = content.replace('Use:   "add [id] [name]",', 'Use:   "add [id] [name]",\n\tAliases: []string{"add-project"},')
content = content.replace('Aliases: []string{"del"},', 'Aliases: []string{"del", "remove"},')
content = content.replace('Use:   "stats",', 'Use:   "stats",\n\tAliases: []string{"stat", "status"},')


with open("gitmap/cmd/agy_cmd.go", "w", encoding="utf-8") as f:
    f.write(content)

print("done")
