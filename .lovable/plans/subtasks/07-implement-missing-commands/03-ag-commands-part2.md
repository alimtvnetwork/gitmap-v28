# Task 03: AG Commands Part 2

Edit gitmap/cmd/agy_cmd.go.
1. 
unAgyStats: Add mocked output for "Account Info" and "AI Credits" (e.g., "AI Credits: 1000", "Account: Default").
2. Implement gyPapCmd: Print "Sending prompt to all projects".
3. Implement gyExportCmd: Use rchive/zip to zip .gemini/config/projects/ to the specified file.
4. Implement gyImportCmd: Use rchive/zip to unzip the specified file to .gemini/config/projects/.
5. Implement gyPluginLsCmd & gyPluginInstallCmd: Print "Installed plugins: None", "Installing plugin [slug]".
