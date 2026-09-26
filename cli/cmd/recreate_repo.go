package cmd

import (
	"fmt"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runRecreate(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("missing repository name to recreate", "cmd.recreate")
	}

	name := args[0]
	isPrivate := false
	for _, arg := range args {
		if arg == "--private" {
			isPrivate = true
		}
	}

	cmdRm := exec.Command("git", "remote", "remove", "origin")
	_ = cmdRm.Run()

	var ghArgs []string
	ghArgs = append(ghArgs, "repo", "create", name)
	if isPrivate {
		ghArgs = append(ghArgs, "--private")
	} else {
		ghArgs = append(ghArgs, "--public")
	}
	ghArgs = append(ghArgs, "--source=.", "--remote=origin", "--push")

	cmdGh := exec.Command("gh", ghArgs...)
	out, err := cmdGh.CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("gh repo create failed: %s", string(out)))
	}

	fmt.Println(string(out))
	return nil
}
