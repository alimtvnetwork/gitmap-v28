# Subtask 2: Context Menu and Install Tools

1. gitmap/constants/constants_install.go: Add ToolAntigravity = "antigravity" and ToolAgCtx = "ag-ctx". Add to descriptions and categories.
2. gitmap/cmd/installctxentries.go: Add {KeyName: "80_antigravity", MUIVerb: "Open project with Antigravity", Args: []string{"ag"}, Mode: constants.CtxModeTerminal, Icon: constants.CtxIconGitmap} to ctxMenu().
3. gitmap/cmd/install.go: Map ToolAgCtx in specialToolHandler to run a new func 
unAgContextMenu (or just print instruction to use install ctx).
