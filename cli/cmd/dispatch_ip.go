// Package cmd — dispatch_ip.go: thin ip dispatch adapter.
//
// Kept from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// runIPCmd still lives in package cmd; when ip_cmd.go moves to cmdip
// (later wave), this adapter is rewired.
package cmd

import "context"

func runIP(args []string) error {
	return runIPCmd(nil, args, context.Background())
}
