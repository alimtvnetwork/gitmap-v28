// Package cmd — sshjoin_distribute.go distributes public keys to multiple machines.
package cmd

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// DistributeKeysToHosts pushes the local public key to each machine in the targets list.

func DistributeKeysToHosts(ctx context.Context, hosts []string, user string, port int) error {
	if len(hosts) == 0 {
		return apperror.New("DistributeKeysToHosts", "E_INVALID_ARGS", map[string]any{"error": "no hosts provided"})
	}

	pubKey, errKey := getLocalPublicKey(ctx, "", false)
	if errKey != nil {
		return errKey
	}

	if user == "" {
		user = "root"
	}
	if port <= 0 {
		port = 22
	}

	for _, host := range hosts {
		target, errParse := ParseSSHTarget(host, user, port)
		if errParse != nil {
			fmt.Printf("⚠️ Failed to parse host %s: %v\n", host, errParse)
			continue
		}
		if errAppend := appendKeyRemote(ctx, pubKey, *target); errAppend != nil {
			fmt.Printf("✗ Failed key append to %s: %v\n", target.String(), errAppend)
		} else {
			fmt.Printf("✓ Key distributed to %s\n", target.String())
		}
	}

	return nil
}
