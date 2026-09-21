package cmdssh

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func printKeysTableHeader(count int) {
	fmt.Printf("\n  %s🔑 Registered SSH Keys (%d)%s\n\n", constants.ColorCyan, count, constants.ColorReset)
	fmt.Printf("    %-20s %-30s %-25s %s\n", "NAME", "FINGERPRINT", "EMAIL", "PRIVATE_KEY_PATH")
	fmt.Printf("    %s\n", constants.TermTableRule)
}

func printKeyRow(k model.SSHKey) {
	fmt.Printf("    %-20s %-30s %-25s %s\n",
		k.Name, k.Fingerprint, k.Email, k.PrivatePath)
}

func printKeysTableFooter() {
	fmt.Println()
	fmt.Printf("  %sTip: Use 'gitmap ssh keys remove <name>' to remove a key with undo support.%s\n\n",
		constants.ColorDim, constants.ColorReset)
}

func renderSSHKeysTable(keys []model.SSHKey) {
	printKeysTableHeader(len(keys))
	if len(keys) == 0 {
		fmt.Printf("    %s(no SSH keys registered yet)%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, k := range keys {
		printKeyRow(k)
	}
	printKeysTableFooter()
}

func runSSHKeysManage(args []string) error {
	_ = args
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "runSSHKeysManage.openDB")
	}
	defer dbConn.Close()

	keys, errList := dbConn.ListSSHKeys()
	if errList != nil {
		return apperror.WrapSimple(errList, "runSSHKeysManage.ListSSHKeys")
	}

	renderSSHKeysTable(keys)
	ctx := context.Background()
	_, _ = EnqueueSSHTask(ctx, ActionManageKeys, "all", fmt.Sprintf("%d keys", len(keys)), "")

	return nil
}
