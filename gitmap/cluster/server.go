package cluster

import (
	"errors"
	"net"
	"net/rpc"
)

// Server represents the orchestrator daemon.
type Server struct {
	token string
}

// NewServer creates a new Server instance.
func NewServer(token string) *Server {
	return &Server{token: token}
}

// HandshakeArgs contains the arguments for the Handshake RPC.
type HandshakeArgs struct {
	Token string
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
	reply.Success = true
	return nil
}

// Serve starts accepting connections on the listener.
func (s *Server) Serve(listener net.Listener) {
	rpcServer := rpc.NewServer()
	rpcServer.Register(s)
	rpcServer.Accept(listener)
}
