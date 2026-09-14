package cmdssh

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func hasHelpToken(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if hasHelpToken(a) {
			return true
		}
	}
	return false
}

func hasJoinToken(args []string) bool {
	for _, a := range args {
		if a == "join" || a == "sj" {
			return true
		}
	}
	return false
}

func checkHelp(command string, args []string) {
	if hasHelpFlag(args) {
		helptext.Print(command)
		cliexit.Exit(0)
	}
}

func checkSSHHelp(args []string) {
	if hasJoinToken(args) && hasHelpFlag(args) {
		helptext.Print("ssh-join")
		cliexit.Exit(0)
	}
	checkHelp("ssh", args)
}

func homeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return h
}

func expandHome(p string) string {
	if p == "~" {
		return homeDir()
	}
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		return filepath.Join(homeDir(), p[2:])
	}
	return p
}

func isGitRepoCWD() bool {
	_, err := os.Stat(".git")
	return err == nil
}

func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}
	if err := db.Migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func reorderFlagsBeforeArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}
	return append(flags, positional...)
}
