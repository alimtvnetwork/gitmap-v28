package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runServe starts the orchestrator daemon and generates a join token.
func runServe(args []string) error {
	checkHelp("serve", args)
	port := parseServeFlags(args)

	fmt.Println(constants.MsgServeStarting)

	// Generate a secure Join Token
	token, err := generateJoinToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrServeTokenGenerate, err)
		cliexit.HandleError(nil, 1)
	}

	// Bind to all network interfaces on the specified port
	address := fmt.Sprintf("%s:%d", constants.ServeBindAddress, port)
	listener, err := net.Listen(constants.ServeProtocol, address)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrServeBind+"\n", err)
		cliexit.HandleError(nil, 1)
	}
	defer listener.Close()

	// Display IP/Port and Token
	fmt.Printf(constants.MsgServeAddress+"\n", getLocalIP(), port)
	fmt.Printf(constants.MsgServeToken+"\n", token)
	fmt.Printf(constants.MsgServeJoinCommand, getLocalIP(), port, token)

	server := cluster.NewServer(token, 30*time.Second)
	go func() {
		server.Serve(listener)
	}()

	// Block until interrupted
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println(constants.MsgServeShutdown)
	return nil
}

// parseServeFlags parses serve-specific CLI flags.
func parseServeFlags(args []string) int {
	fs := flag.NewFlagSet(constants.CmdServe, flag.ExitOnError)
	port := fs.Int(constants.FlagServePort, constants.ServeDefaultPort, constants.FlagDescServePort)
	fs.Parse(args)
	return *port
}

// generateJoinToken generates a 32-byte secure hex token.
func generateJoinToken() (string, error) {
	bytes := make([]byte, constants.ServeTokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// getLocalIP returns the first non-loopback local IPv4 address, or "0.0.0.0" if none found.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return constants.ServeBindAddress
	}

	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)
		if ok == false {
			continue
		}
		if ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return constants.ServeBindAddress
}
