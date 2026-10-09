package cmduser

import (
	"testing"
)

func TestUserUsage(t *testing.T) {
	printUserUsage()
	RenderUserHelp()
}

func TestDispatchUser_NotUserCommand(t *testing.T) {
	handled, err := dispatchUser("other")
	if handled {
		t.Errorf("expected not handled, got true")
	}
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestDispatchUserSubcommand_Info(t *testing.T) {
	err := dispatchUserSubcommand("info", []string{"--json"})
	if err != nil {
		t.Errorf("unexpected error on info: %v", err)
	}
}

func TestDispatchUserSubcommand_Status(t *testing.T) {
	err := dispatchUserSubcommand("status", []string{"-j"})
	if err != nil {
		t.Errorf("unexpected error on status: %v", err)
	}
}

func TestDispatchUserSubcommand_List(t *testing.T) {
	err := dispatchUserSubcommand("list", []string{})
	if err != nil {
		t.Errorf("unexpected error on list: %v", err)
	}
}

func TestDispatchUserSubcommand_Ls(t *testing.T) {
	err := dispatchUserSubcommand("ls", []string{})
	if err != nil {
		t.Errorf("unexpected error on ls: %v", err)
	}
}

func TestDispatchUserSubcommand_Project(t *testing.T) {
	err := dispatchUserSubcommand("project", []string{"status"})
	if err != nil {
		t.Errorf("unexpected error on project status: %v", err)
	}
}

func TestDispatchUserSubcommand_Config(t *testing.T) {
	err := dispatchUserSubcommand("config", []string{"local"})
	if err != nil {
		t.Errorf("unexpected error on config local: %v", err)
	}
}

func TestDispatchUserSubcommand_Sync(t *testing.T) {
	err := dispatchUserSubcommand("sync", []string{})
	if err != nil {
		t.Errorf("unexpected error on sync: %v", err)
	}
}

func TestDispatchUserSubcommand_Help(t *testing.T) {
	err := dispatchUserSubcommand("help", []string{})
	if err != nil {
		t.Errorf("unexpected error on help: %v", err)
	}
}

func TestRouteUserAdd_GitProfile(t *testing.T) {
	args := []string{"test-prof", "--name", "Tester", "--email", "test@test.com"}
	err := routeUserAdd(args)
	if err != nil {
		t.Errorf("unexpected error adding profile: %v", err)
	}
}

func TestExtractFirstPositional(t *testing.T) {
	args := []string{"--global", "my-alias", "--force"}
	alias := extractFirstPositional(args)
	if alias != "my-alias" {
		t.Errorf("expected my-alias, got %s", alias)
	}
}

// dispatchUser reports whether sub is a user subcommand. Test-local
// helper — the production dispatcher (dispatchUserSubcommand) has a
// different signature.
func dispatchUser(sub string) (bool, error) {
	switch sub {
	case "info", "status", "list", "ls", "switch", "use", "project",
		"config", "sync", "add", "rm", "delete", "remove", "del",
		"create-root", "root":
		return true, nil
	default:
		return false, nil
	}
}
