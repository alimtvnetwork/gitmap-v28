package cmddb

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const dbRestHelp = `gitmap db-rest — Launch REST API server or bind database alias.

Usage:
  gitmap db-rest "<db-path>" as <alias>
  gitmap db-rest <targets> [ports]

Arguments:
  <targets>   Single alias, comma-separated aliases, or direct file paths
              Examples: 'master', 'master,pipeline', 'myapp,./data/store.db'
  [ports]     Optional comma-separated ports to probe (default: 8080,8081,8082,3000,5000)

Examples:
  gitmap db-rest "./data/app.db" as appdb
  gitmap db-rest appdb 8080,8081
  gitmap db-rest master,pipeline 8080,8081,3000
  gitmap db-rest appdb,./external/store.db 9000
`

// RunDbRest executes the REST API server daemon or alias binding command.
func RunDbRest(args []string) *appfault.AppError {
	if len(args) == 0 || (len(args) == 1 && isDbHelpToken(args[0])) {
		fmt.Print(dbRestHelp)
		return nil
	}

	if isAliasBindingArgs(args) {
		return handleAliasBinding(args)
	}

	return handleRestServerStartup(args)
}

func isAliasBindingArgs(args []string) bool {
	for i, a := range args {
		if strings.EqualFold(a, "as") && i > 0 && i < len(args)-1 {
			return true
		}
	}
	return false
}

func handleAliasBinding(args []string) *appfault.AppError {
	asIdx := 0
	for i, a := range args {
		if strings.EqualFold(a, "as") {
			asIdx = i
			break
		}
	}

	rawPath := strings.Trim(strings.TrimSpace(strings.Join(args[:asIdx], " ")), `"'`)
	alias := strings.Trim(strings.TrimSpace(strings.Join(args[asIdx+1:], " ")), `"'`)

	if rawPath == "" || alias == "" {
		return appfault.NewValidationError("usage: gitmap db-rest \"<db-path>\" as <alias>")
	}

	cleanPath := filepath.Clean(rawPath)
	if _, err := os.Stat(cleanPath); err != nil {
		return appfault.NewNotFoundError(fmt.Sprintf("database file not found: %s", cleanPath))
	}

	if err := RegisterAlias(alias, cleanPath, "Registered via db-rest"); err != nil {
		return appfault.NewValidationError(fmt.Sprintf("register alias '%s': %v", alias, err))
	}

	size := getFileSize(cleanPath)
	tables := countTablesInDb(cleanPath)
	fmt.Printf("[SUCCESS] Registered database alias '%s' -> %s (%s, %d tables)\n", alias, cleanPath, formatBytes(size), tables)
	fmt.Printf("You can now start the REST server with: gitmap db-rest %s 8080\n", alias)
	return nil
}

func handleRestServerStartup(args []string) *appfault.AppError {
	targetArg := args[0]
	portArg := ""
	if len(args) > 1 {
		portArg = args[1]
	}

	candidates, err := ParseCandidatePorts(portArg)
	if err != nil {
		return appfault.NewValidationError(err.Error())
	}

	boundPort, err := PickAvailablePort(candidates)
	if err != nil {
		return appfault.NewValidationError(err.Error())
	}

	targets, err := ResolveMountTargets(targetArg)
	if err != nil {
		return appfault.NewNotFoundError(err.Error())
	}

	printRestBanner(boundPort, targets)

	return startRestListenerAndServe(targets, boundPort)
}

func startRestListenerAndServe(targets []DbMountTarget, port int) *appfault.AppError {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return appfault.NewValidationError(fmt.Sprintf("listen on port %d: %v", port, err))
	}
	defer listener.Close()

	dbs := make(map[string]string, len(targets))
	for _, t := range targets {
		dbs[t.Alias] = t.DatabasePath
	}

	server, srvErr := StartMultiDbRestServer(dbs, listener)
	if srvErr != nil {
		return appfault.NewValidationError(fmt.Sprintf("REST server failure: %v", srvErr))
	}

	waitForServerShutdown(server)
	return nil
}

func waitForServerShutdown(server *http.Server) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func printRestBanner(port int, targets []DbMountTarget) {
	fmt.Println("================================================================================")
	fmt.Println("                    GITMAP DATABASE REST SERVER (v28)                           ")
	fmt.Println("================================================================================")
	fmt.Printf(" Server Listening: http://127.0.0.1:%d\n", port)
	fmt.Printf(" Mounted Databases (%d):\n", len(targets))
	for _, t := range targets {
		fmt.Printf("   * %-12s -> %s (%s, %d tables)\n", t.Alias, t.DatabasePath, formatBytes(t.FileSizeBytes), t.TableCount)
	}
	fmt.Println("")
	fmt.Println(" Endpoints:")
	fmt.Printf("   * Database Catalog:    http://127.0.0.1:%d/api/v1/databases\n", port)
	fmt.Printf("   * Interactive Swagger: http://127.0.0.1:%d/swagger\n", port)
	fmt.Printf("   * Database Visualizer: http://127.0.0.1:%d/ui\n", port)
	fmt.Println("================================================================================")
	fmt.Println(" Press Ctrl+C to stop the server.")
}
