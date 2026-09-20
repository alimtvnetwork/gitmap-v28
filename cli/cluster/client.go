package cluster

import (
	"crypto/tls"
	"net"
	"net/rpc"
	"time"
)

// NodeClient represents a node joining the cluster.
type NodeClient struct {
	id      string
	address string
	token   string
}

// NewNodeClient creates a new NodeClient instance.
func NewNodeClient(id, address, token string) *NodeClient {
	return &NodeClient{id: id, address: address, token: token}
}

func (c *NodeClient) dialTLS() (*rpc.Client, error) {
	conf := &tls.Config{InsecureSkipVerify: true}
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := c.dialWithRetry(conf, dialer)
	if err != nil {
		return nil, err
	}

	return rpc.NewClient(conn), nil
}

func (c *NodeClient) dialWithRetry(conf *tls.Config, dialer *net.Dialer) (net.Conn, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		conn, err := tls.DialWithDialer(dialer, "tcp", c.address, conf)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}

	return nil, lastErr
}

// Handshake connects to the server and verifies the join token.
func (c *NodeClient) Handshake() error {
	client, err := c.dialTLS()
	if err != nil {
		return err
	}

	defer client.Close()

	args := &HandshakeArgs{Token: c.token, ID: c.id}
	var reply HandshakeReply
	if err := client.Call("Server.Handshake", args, &reply); err != nil {
		return err
	}

	return nil
}

// Ping sends a heartbeat to the server.
func (c *NodeClient) Ping() error {
	client, err := c.dialTLS()
	if err != nil {
		return err
	}

	defer client.Close()

	args := &PingArgs{ID: c.id}
	var reply PingReply

	return client.Call("Server.Ping", args, &reply)
}
