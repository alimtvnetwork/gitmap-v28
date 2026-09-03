package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runProfiles(args []string) error {
	checkHelp("profiles", args)
	if len(args) == 0 {
		return runProfilesList(args)
	}
	sub := args[0]
	rest := args[1:]
	return routeProfilesSub(sub, rest)
}

func routeProfilesSub(sub string, args []string) error {
	switch sub {
	case "ls", "list":
		return runProfilesList(args)
	case "set-default", "default":
		return runProfilesSetDefault(args)
	case "switch", "use":
		return runProfilesSwitch(args)
	case "add", "create":
		return runProfilesAdd(args)
	case "rm", "remove", "delete":
		return runProfilesRemove(args)
	case "status":
		return runProfilesStatus(args)
	default:
		return runProfilesList(append([]string{sub}, args...))
	}
}

func runProfilesList(args []string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return apperror.WrapSimple(err, "load profiles:")
	}
	if hasArgFlag(args, "--json") {
		return outputProfilesJSON(cfg)
	}
	printProfilesTable(cfg)
	return nil
}

func outputProfilesJSON(cfg model.GitProfileConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal json:")
	}
	fmt.Println(string(data))
	return nil
}

func printProfilesTable(cfg model.GitProfileConfig) {
	fmt.Printf("\n  %s● Git Accounts & Profiles (Total: %d)%s\n",
		constants.ColorCyan, len(cfg.Profiles), constants.ColorReset)
	fmt.Println("  --------------------------------------------------------------------------------")
	fmt.Println("  #    Name                 Provider   Type           Default  Usage  Last Used")
	fmt.Println("  --------------------------------------------------------------------------------")
	for i, p := range cfg.Profiles {
		printProfileRow(i+1, p, cfg.Default)
	}
	fmt.Println("  --------------------------------------------------------------------------------")
	fmt.Println("  Tip: Switch default with 'gitmap profiles set-default <1-N|name>'")
	fmt.Println()
}

func printProfileRow(seq int, p model.GitProfile, defaultName string) {
	defTag := "-"
	if p.Name == defaultName || p.IsDefault {
		defTag = "* (default)"
	}
	lastUsed := "-"
	if !p.LastUsedAt.IsZero() {
		lastUsed = p.LastUsedAt.Format("2006-01-02")
	}
	fmt.Printf("  [%d]  %-20s %-10s %-14s %-8s %-6d %s\n",
		seq, p.Name, p.Provider, p.Type, defTag, p.UsageCount, lastUsed)
}
