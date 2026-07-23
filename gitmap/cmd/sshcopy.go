package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runSSHCopy prints the public key (like `ssh cat`) AND pushes it to the
// OS clipboard. Falls back to a "copy it manually" warning when no
// clipboard tool is available — never fails the user.
func runSSHCopy(args []string) {
	fs := flag.NewFlagSet("ssh-copy", flag.ExitOnError)
	nameFlag := fs.String("name", constants.DefaultSSHKeyName, "Key name")
	fs.StringVar(nameFlag, "n", constants.DefaultSSHKeyName, "Key name (short)")
	fs.Parse(args)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(*nameFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, *nameFlag)
		printAvailableKeys(db)
		os.Exit(1)
	}

	pub := strings.TrimSpace(key.PublicKey)
	fmt.Println(pub)

	copyPubKeyAndAnnounce(pub)
}

// copyPubKeyAndAnnounce pushes the public key to the OS clipboard and
// prints a friendly emoji status line to stderr. Soft-fails when no
// clipboard tool is on PATH so the caller never blocks.
func copyPubKeyAndAnnounce(pub string) {
	tool, err := writeClipboard(pub)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHClipboard, tool, err)

		return
	}
	if tool == "" {
		fmt.Fprint(os.Stderr, constants.MsgSSHCopyFallback)

		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgSSHCopied, len(pub))
}

// writeClipboard pipes text to the platform-native clipboard tool.
// Returns ("", nil) when no tool is available (soft fallback).
func writeClipboard(text string) (string, error) {
	tool, args := resolveClipboardTool()
	if tool == "" {
		return "", nil
	}
	cmd := exec.Command(tool, args...)
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return tool, err
	}

	return tool, nil
}

// resolveClipboardTool picks the right clipboard binary for the host OS.
func resolveClipboardTool() (string, []string) {
	switch runtime.GOOS {
	case "windows":
		return "clip", nil
	case "darwin":
		return "pbcopy", nil
	}
	if path, err := exec.LookPath("wl-copy"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("xclip"); err == nil {
		return path, []string{"-selection", "clipboard"}
	}
	if path, err := exec.LookPath("xsel"); err == nil {
		return path, []string{"--clipboard", "--input"}
	}

	return "", nil
}
