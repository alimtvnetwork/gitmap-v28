package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runJoin connects to an existing orchestrator daemon.
func runJoin(args []string) *apperror.AppError {
	CheckHelpOrEmpty(constants.CmdJoin, args)

	address, token, appErr := parseJoinArgs(args)
	if appErr != nil {
		return appErr
	}

	return executeJoinHandshake(address, token)
}

func parseJoinArgs(args []string) (string, string, *apperror.AppError) {
	flags, positional := splitJoinArgs(args)
	fs := flag.NewFlagSet(constants.CmdJoin, flag.ContinueOnError)
	token := fs.String(constants.FlagJoinToken, "", constants.FlagDescJoinToken)
	if err := fs.Parse(flags); err != nil {
		return "", "", apperror.NewValidationError(err.Error())
	}

	positional = append(positional, fs.Args()...)
	if appErr := validateJoinInputs(positional, *token); appErr != nil {
		return "", "", appErr
	}

	return positional[0], *token, nil
}

func splitJoinArgs(args []string) ([]string, []string) {
	var flags, positional []string
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if isJoinFlagToken(arg) {
			flags = append(flags, arg)
			skipNext = appendJoinFlagValue(args, i, &flags)
			continue
		}

		positional = append(positional, arg)
	}

	return flags, positional
}

func isJoinFlagToken(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func appendJoinFlagValue(args []string, i int, flags *[]string) bool {
	if (args[i] == "--token" || args[i] == "-token") && i+1 < len(args) {
		*flags = append(*flags, args[i+1])

		return true
	}

	return false
}

func validateJoinInputs(positional []string, token string) *apperror.AppError {
	if len(positional) < 1 {
		return apperror.NewValidationError(constants.ErrJoinMissingAddress)
	}

	if token == "" {
		return apperror.NewValidationError(constants.ErrJoinMissingToken)
	}

	return nil
}

func executeJoinHandshake(address, token string) *apperror.AppError {
	fmt.Printf(constants.MsgJoinStarting+"\n", address)
	hostname := resolveJoinHostname()
	client := cluster.NewNodeClient(hostname, address, token)
	if handshakeErr := client.Handshake(); handshakeErr != nil {
		return wrapJoinHandshakeError(handshakeErr, address, hostname)
	}

	fmt.Println(constants.MsgJoinSuccess)

	return nil
}

func resolveJoinHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-node"
	}

	return hostname
}

func wrapJoinHandshakeError(err error, address, hostname string) *apperror.AppError {
	msg := fmt.Sprintf("Failed to join cluster at %s: %v", address, err)
	ctx := map[string]any{"address": address, "hostname": hostname}

	return apperror.WrapWithDetails(
		err,
		"join.Handshake",
		"E8005",
		msg,
		"cmd.runJoin",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		ctx,
	)
}
