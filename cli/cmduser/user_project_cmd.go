package cmduser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

func runUserProject(args []string) error {
	if len(args) == 0 {
		return showProjectStatus()
	}

	switch args[0] {
	case "bind":
		return handleProjectBind(args[1:])
	case "unbind":
		return handleProjectUnbind()
	case "status":
		return showProjectStatus()
	default:
		fmt.Fprintf(os.Stderr, "Unknown project command: %s. Use bind, unbind, or status\n", args[0])
		return nil
	}
}

func handleProjectBind(args []string) error {
	alias := extractFirstPositional(args)
	if len(alias) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap user project bind <alias>\n")
		return nil
	}

	if err := usercontext.BindProject(".", alias); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding project: %v\n", err)
		return nil
	}

	usercontext.SyncProjectUser(".")
	fmt.Printf("✔ Bound current project to profile %q\n", alias)

	return nil
}

func handleProjectUnbind() error {
	if err := usercontext.UnbindProject("."); err != nil {
		fmt.Fprintf(os.Stderr, "Error unbinding project: %v\n", err)
		return nil
	}

	fmt.Println("✔ Unbound current project profile")

	return nil
}

func showProjectStatus() error {
	prof, isBound := usercontext.ResolveProjectUser(".")
	if !isBound || prof == nil {
		fmt.Println("Project is not bound to any profile (using global Git config)")
		return nil
	}

	fmt.Printf("Project bound to profile: %s (%s <%s>)\n", prof.ID, prof.Name, prof.Email)

	return nil
}

func runUserConfig(args []string) error {
	name, email, isGlobal, isInteractive := parseConfigArgs(args)
	if isInteractive {
		return runInteractiveUserConfig(isGlobal)
	}
	if len(name) > 0 || len(email) > 0 {
		return applyUserConfigUpdate(name, email, isGlobal)
	}

	return displayUserConfig(isGlobal)
}

func parseConfigArgs(args []string) (string, string, bool, bool) {
	isGlobal := hasUserFlag(args, "--global") || hasUserFlag(args, "-g") || hasPositionalWord(args, "global")
	isInteractive := hasUserFlag(args, "--interactive") || hasUserFlag(args, "-i")
	name := extractUserArgFlag(args, "--name")
	email := extractUserArgFlag(args, "--email")
	if len(name) > 0 || len(email) > 0 || isInteractive {
		return name, email, isGlobal, isInteractive
	}
	posArgs := filterPositionalArgs(args)
	posName, posEmail := resolvePositionalIdentity(posArgs)

	return posName, posEmail, isGlobal, isInteractive
}

func hasPositionalWord(args []string, word string) bool {
	for _, a := range args {
		if strings.EqualFold(a, word) {
			return true
		}
	}

	return false
}

func filterPositionalArgs(args []string) []string {
	var res []string
	for _, a := range args {
		if isConfigNoiseArg(a) {
			continue
		}
		res = append(res, a)
	}

	return res
}

func isConfigNoiseArg(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return true
	}
	lower := strings.ToLower(arg)

	return lower == "global" || lower == "local" || lower == "config" || lower == "set"
}

func resolvePositionalIdentity(args []string) (string, string) {
	if len(args) == 0 {
		return "", ""
	}
	if len(args) == 1 {
		return resolveSinglePositional(args[0])
	}
	if strings.Contains(args[0], "@") {
		return args[1], args[0]
	}

	return args[0], args[1]
}

func resolveSinglePositional(arg string) (string, string) {
	if strings.Contains(arg, "@") {
		return "", arg
	}

	return arg, ""
}

func runInteractiveUserConfig(isGlobal bool) error {
	reader := bufio.NewReader(os.Stdin)
	scope := resolveConfigScope(isGlobal)
	fmt.Printf("Configure %s Git Identity:\n", scope)
	fmt.Print("  Enter user.name : ")
	name, _ := reader.ReadString('\n')
	fmt.Print("  Enter user.email: ")
	email, _ := reader.ReadString('\n')

	return applyUserConfigUpdate(strings.TrimSpace(name), strings.TrimSpace(email), isGlobal)
}

func applyUserConfigUpdate(name, email string, isGlobal bool) error {
	if err := usercontext.SetGitIdentity(name, email, isGlobal); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating git config: %v\n", err)

		return nil
	}
	scope := resolveConfigScope(isGlobal)
	ident, _ := usercontext.GetGitIdentity(isGlobal)
	fmt.Printf("✔ Updated %s Git config: %s <%s>\n", scope, ident.Name, ident.Email)

	return nil
}

func displayUserConfig(isGlobal bool) error {
	ident, err := usercontext.GetGitIdentity(isGlobal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading git config: %v\n", err)

		return nil
	}
	scope := resolveConfigScope(isGlobal)
	nameDisplay := resolveDisplayVal(ident.Name)
	emailDisplay := resolveDisplayVal(ident.Email)
	fmt.Printf("%s Git Identity:\n  Name:  %s\n  Email: %s\n", scope, nameDisplay, emailDisplay)
	printConfigGuidance(scope, ident)

	return nil
}

func resolveDisplayVal(val string) string {
	if val == "" {
		return "(not configured)"
	}

	return val
}

func printConfigGuidance(scope string, ident usercontext.GitIdentity) {
	if ident.Name != "" && ident.Email != "" {
		return
	}
	scopeFlag := strings.ToLower(scope)
	fmt.Printf("\nTip: Set your Git identity from the terminal:\n")
	fmt.Printf("  gitmap user config %s --name \"<Your Name>\" --email \"<Your Email>\"\n", scopeFlag)
	fmt.Printf("  gitmap user config %s \"<Your Name>\" \"<Your Email>\"\n", scopeFlag)
	fmt.Printf("  gitmap user config %s -i   # interactive prompt\n", scopeFlag)
}

func resolveConfigScope(isGlobal bool) string {
	if isGlobal {
		return "Global"
	}

	return "Local"
}

func runUserSync(args []string) error {
	prof, isBound := usercontext.ResolveProjectUser(".")
	if !isBound || prof == nil {
		fmt.Println("No project profile binding found to sync.")
		return nil
	}

	if err := usercontext.SyncProjectUser("."); err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing project user: %v\n", err)
		return nil
	}

	fmt.Printf("✔ Synchronized local repository with profile %q (%s <%s>)\n", prof.ID, prof.Name, prof.Email)

	return nil
}
