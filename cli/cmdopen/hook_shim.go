package cmdopen

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// ResolveEndpointStringFn is wired by cmd/di_hooks.go.
var ResolveEndpointStringFn func(endpoint string) string

func resolveEndpointString(endpoint string) string {
	if ResolveEndpointStringFn != nil {
		return ResolveEndpointStringFn(endpoint)
	}
	return ""
}
