package cluster

import (
	"errors"
	"net"
	"net/rpc"
	"time"
)

// Server represents the orchestrator daemon.
type Server struct {
	token    string
	registry *Registry
}

// NewServer creates a new Server instance.
func NewServer(token string, heartbeatTimeout time.Duration) *Server {
	return &Server{
		token:    token,
		registry: NewRegistry(heartbeatTimeout),
	}
}

// HandshakeArgs contains the arguments for the Handshake RPC.
type HandshakeArgs struct {
	Token string
	ID    string
}

// HandshakeReply contains the reply for the Handshake RPC.
type HandshakeReply struct {
	Success bool
}

// Handshake verifies the join token.
func (s *Server) Handshake(args *HandshakeArgs, reply *HandshakeReply) error {
	if args.Token != s.token {
		reply.Success = false
		return errors.New("invalid join token")
	}
	s.registry.Register(args.ID)
	reply.Success = true
	return nil
}

// PingArgs contains the arguments for the Ping RPC.
type PingArgs struct {
	ID string
}

// PingReply contains the reply for the Ping RPC.
type PingReply struct {
	Success bool
}

// Ping updates the last seen time for a node.
func (s *Server) Ping(args *PingArgs, reply *PingReply) error {
	s.registry.Ping(args.ID)
	reply.Success = true
	return nil
}

// DisconnectArgs contains the arguments for the Disconnect RPC.
type DisconnectArgs struct {
	ID string
}

// DisconnectReply contains the reply for the Disconnect RPC.
type DisconnectReply struct {
	Success bool
}

// Disconnect gracefully marks a node as disconnected in the registry.
func (s *Server) Disconnect(args *DisconnectArgs, reply *DisconnectReply) error {
	s.registry.Disconnect(args.ID)
	reply.Success = true
	return nil
}

// Serve starts accepting connections on the listener.
func (s *Server) Serve(listener net.Listener) {
	rpcServer := rpc.NewServer()
	rpcServer.Register(s)
	rpcServer.Accept(listener)
}
