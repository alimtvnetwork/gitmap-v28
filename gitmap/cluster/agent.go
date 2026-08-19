package cluster

import (
	"context"
	"net"
	"net/rpc"
)

// Agent handles remote execution requests.
type Agent struct{}

func (a *Agent) ExecPS(args *AgentExecArgs, reply *AgentExecReply) error {
	stdout, stderr, exitCode, err := ExecPS(context.Background(), ClusterNode{}, args.Command)
	reply.Stdout = stdout
	reply.Stderr = stderr
	reply.ExitCode = exitCode
	return err
}

func (a *Agent) ExecCmd(args *AgentExecArgs, reply *AgentExecReply) error {
	stdout, stderr, exitCode, err := ExecCmd(context.Background(), ClusterNode{}, args.Command)
	reply.Stdout = stdout
	reply.Stderr = stderr
	reply.ExitCode = exitCode
	return err
}

// ServeAgent starts an RPC server for the agent on the given listener.
func ServeAgent(listener net.Listener) {
	server := rpc.NewServer()
	agent := &Agent{}
	server.Register(agent)
	server.Accept(listener)
}
