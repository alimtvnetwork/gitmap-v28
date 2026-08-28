package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runJoin connects to an existing orchestrator daemon.
func runJoin(args []string) error {
	checkHelp("join", args)

	fs := flag.NewFlagSet(constants.CmdJoin, flag.ExitOnError)
	token := fs.String(constants.FlagJoinToken, "", constants.FlagDescJoinToken)
	fs.Parse(args)

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		return apperror.NewSimple(constants.ErrJoinMissingAddress, "E9000")
	}
	if *token == "" {
		return apperror.NewSimple(constants.ErrJoinMissingToken, "E9000")
	}

	address := positionalArgs[0]
	fmt.Printf(constants.MsgJoinStarting+"\n", address)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-node"
	}
	client := cluster.NewNodeClient(hostname, address, *token)
	if err := client.Handshake(); err != nil {
		return apperror.WrapSimple(err, constants.ErrJoinFailed)
	}

	fmt.Println(constants.MsgJoinSuccess)
	return nil
}
