package cluster

import (
	"net/rpc"
)

// NodeClient represents a node joining the cluster.
type NodeClient struct {
	address string
	token   string
}

// NewNodeClient creates a new NodeClient instance.
func NewNodeClient(address, token string) *NodeClient {
	return &NodeClient{address: address, token: token}
}

// Handshake connects to the server and verifies the join token.
func (c *NodeClient) Handshake() error {
	client, err := rpc.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	defer client.Close()

	args := &HandshakeArgs{Token: c.token}
	var reply HandshakeReply
	if err := client.Call("Server.Handshake", args, &reply); err != nil {
		return err
	}

	return nil
}
