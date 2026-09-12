package heavy_test

import (
	"crypto/tls"
	"net"
	"net/rpc"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
)

func TestAgentIntegration(t *testing.T) {
	// 1. Generate TLS config
	tlsConf, err := cluster.GenerateTLSConfig()
	if err != nil {
		t.Fatalf("Failed to generate TLS config: %v", err)
	}

	// 2. Start Agent listener
	listener, err := tls.Listen("tcp", "127.0.0.1:0", tlsConf)
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	go cluster.ServeAgent(listener)

	// wait for server to start
	time.Sleep(100 * time.Millisecond)

	// 3. Dispatch a command to it
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to split host port: %v", err)
	}

	conf := &tls.Config{InsecureSkipVerify: true}
	client, err := tls.Dial("tcp", "127.0.0.1:"+port, conf)
	if err != nil {
		t.Fatalf("Failed to dial agent: %v", err)
	}
	defer client.Close()

	rpcClient := rpc.NewClient(client)
	args := &cluster.AgentExecArgs{Command: "echo hello"}
	var reply cluster.AgentExecReply

	err = rpcClient.Call("Agent.ExecCmd", args, &reply)
	if err != nil {
		t.Fatalf("Agent.ExecCmd failed: %v", err)
	}
}
