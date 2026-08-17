package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runJoin connects to an existing orchestrator daemon.
func runJoin(args []string) {
	checkHelp("join", args)
	
	fs := flag.NewFlagSet(constants.CmdJoin, flag.ExitOnError)
	token := fs.String(constants.FlagJoinToken, "", constants.FlagDescJoinToken)
	fs.Parse(args)

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		fmt.Fprintln(os.Stderr, constants.ErrJoinMissingAddress)
		os.Exit(1)
	}
	if *token == "" {
		fmt.Fprintln(os.Stderr, constants.ErrJoinMissingToken)
		os.Exit(1)
	}

	address := positionalArgs[0]
	fmt.Printf(constants.MsgJoinStarting+"\n", address)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-node"
	}
	client := cluster.NewNodeClient(hostname, address, *token)
	if err := client.Handshake(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrJoinFailed, err)
		os.Exit(1)
	}

	fmt.Println(constants.MsgJoinSuccess)
}
