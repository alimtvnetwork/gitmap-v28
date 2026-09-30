package cmdssh

import (
	"fmt"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func printAuthKeyAddHelp() {
	fmt.Printf("\n%sInstall SSH Public Key into Local Authorized Keys%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Usage:")
	fmt.Println("  gitmap ssh auth-key-add [public-key-or-file]")
	fmt.Println("  gitmap ssh-key add [public-key-or-file]")
	fmt.Println("  gitmap ssh key add [public-key-or-file]")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap ssh auth-key-add \"ssh-ed25519 AAAAC3NzaC... user@laptop\"")
	fmt.Println("  gitmap ssh auth-key-add ~/.ssh/id_ed25519.pub")
	fmt.Println("  gitmap ssh auth-key-add   (prompts to paste key interactively)")
	fmt.Println("\n💡 Related SSH Key Commands:")
	fmt.Println("  gitmap ssh deploy-keys                - Sync all authorized keys across all cluster nodes")
	fmt.Println("  gitmap ssh copy-id <alias|user@ip>    - Push your local public key to a remote machine")
	fmt.Println("  gitmap ssh key ls                     - List local SSH keys and fingerprints")
}

func installAuthorizedKeyCrossPlatform(key, keyBlob string) (int, error) {
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		return installWindowsAuthKeys(key, keyBlob), nil
	}

	return installPosixAuthKeys(key, keyBlob)
}

func printAuthKeyAddSummary(count int) {
	hasAdded := count > 0
	if hasAdded {
		fmt.Printf("\n%s✓ Public key installed successfully (%d file(s) updated).%s\n", constants.ColorGreen, count, constants.ColorReset)
	} else {
		fmt.Println("\nKey is already authorized.")
	}
	fmt.Printf("\n  %s💡 Next Steps & Key Deployment Suggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    • Propagate all keys across cluster: gitmap ssh deploy-keys  (or: gitmap deploy-keys-all)")
	fmt.Println("    • Push your key to a remote node:    gitmap ssh copy-id <alias|user@ip>")
	fmt.Println("    • Inspect registered cluster nodes:  gitmap ssh ls")
	fmt.Println()
}

// RunSSHAuthKeyAddCLI adds a public key from another machine to local authorized_keys.
func RunSSHAuthKeyAddCLI(args []string) error {
	hasHelp := hasHelpFlag(args)
	if hasHelp {
		printAuthKeyAddHelp()

		return nil
	}

	raw, err := readPublicKeyInput(args)
	if err != nil {
		return err
	}

	validKey, keyBlob, valErr := validateSSHPublicKey(raw)
	if valErr != nil {
		return valErr
	}

	count, installErr := installAuthorizedKeyCrossPlatform(validKey, keyBlob)
	if installErr != nil {
		return installErr
	}

	printAuthKeyAddSummary(count)

	return nil
}
