package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func handleAutoLoginEnable(args []string) error {
	cfg := parseAutoLoginArgs(args)
	if len(cfg.Username) == 0 {
		return apperror.NewSimple("username is required (use -u <username> or interactive mode)", "E_AUTOLOGIN_NO_USER")
	}
	engine := GetAutoLoginEngine()
	if err := engine.Configure(cfg); err != nil {
		return err
	}
	fmt.Printf("✔ OS auto-login enabled for %s (domain: %s)\n", cfg.Username, cfg.Domain)
	return nil
}

func handleAutoLoginFallback(args []string) error {
	if strings.HasPrefix(args[0], "-") {
		return handleAutoLoginEnable(args)
	}
	printAutoLoginUsage()
	return apperror.NewSimple("unknown autologin command: "+args[0], "E_INVALID_AUTOLOGIN_CMD")
}

func parseAutoLoginArgs(args []string) AutoLoginConfig {
	cfg := AutoLoginConfig{Domain: ".", IsEnabled: true}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isMatchFlag(arg, "-u", "--user", "--username") && i+1 < len(args) {
			cfg.Username = args[i+1]
			i++
			continue
		}
		if isMatchFlag(arg, "-d", "--domain") && i+1 < len(args) {
			cfg.Domain = args[i+1]
			i++
			continue
		}
		if isMatchFlag(arg, "-p", "--pass", "--password") && i+1 < len(args) {
			cfg.Password = args[i+1]
			i++
			continue
		}
	}
	return cfg
}

func isMatchFlag(arg string, targets ...string) bool {
	for _, target := range targets {
		if arg == target {
			return true
		}
	}
	return false
}

func printAutoLoginUsage() {
	fmt.Println("Usage: gitmap os autologin [command] [flags]")
	fmt.Println("       gitmap os autologin                    (Interactive 3-parameter mode)")
	fmt.Println("       gitmap os autologin status             (Check current configuration)")
	fmt.Println("       gitmap os autologin enable -u <user> -d <domain> -p <pass>")
	fmt.Println("       gitmap os autologin disable            (Disable auto-login)")
}
