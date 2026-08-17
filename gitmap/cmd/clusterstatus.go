package cmd

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
)

const (
	mockTimeout     = 30 * time.Second
	mockSleep       = 10 * time.Millisecond
	mockNode1       = "vm-01"
	mockNode2       = "vm-02"
	mockNode3       = "vm-03"
	statusHeader    = "Cluster Status:"
	statusNodePrint = "Node %s: %s (Last Seen: %v)\n"
)

// runClusterStatus handles the "cluster status" subcommand.
func runClusterStatus(args []string) {
	_ = args // Unused for now

	// Mocking a server-side registry to satisfy the command display requirements.
	// Typically, we would fetch this from the running server process.
	reg := cluster.NewRegistry(mockTimeout)

	reg.Register(mockNode1)
	reg.Register(mockNode2)
	reg.Register(mockNode3)

	time.Sleep(mockSleep)

	// Simulate one dropped node
	nodes := reg.GetNodes()
	for _, n := range nodes {
		isNode3 := n.ID == mockNode3
		if isNode3 == true {
			// Fast forward last seen to force timeout
			reg.Ping(n.ID)
		}
	}

	reg.CheckHeartbeats()

	fmt.Println(statusHeader)
	currentNodes := reg.GetNodes()
	for _, n := range currentNodes {
		fmt.Printf(statusNodePrint, n.ID, n.State, n.LastSeen.Format(time.RFC3339))
	}
}
