package cmduser

import (
	"fmt"
	"os"

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
	isGlobal := hasUserFlag(args, "--global") || (len(args) > 0 && args[0] == "global")
	name := extractUserArgFlag(args, "--name")
	email := extractUserArgFlag(args, "--email")

	if len(name) > 0 || len(email) > 0 {
		return applyUserConfigUpdate(name, email, isGlobal)
	}

	return displayUserConfig(isGlobal)
}

func applyUserConfigUpdate(name, email string, isGlobal bool) error {
	if err := usercontext.SetGitIdentity(name, email, isGlobal); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating git config: %v\n", err)
		return nil
	}

	scope := resolveConfigScope(isGlobal)
	fmt.Printf("✔ Updated %s Git config: %s <%s>\n", scope, name, email)

	return nil
}

func displayUserConfig(isGlobal bool) error {
	ident, err := usercontext.GetGitIdentity(isGlobal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading git config: %v\n", err)
		return nil
	}

	scope := resolveConfigScope(isGlobal)
	fmt.Printf("%s Git Identity:\n  Name:  %s\n  Email: %s\n", scope, ident.Name, ident.Email)

	return nil
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
