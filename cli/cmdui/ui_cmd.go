package cmdui

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

// uiListenHost and uiNoBrowser carry the parsed `gitmap ui` flags for the
// server bootstrap in ui_server.go. They default to the secure loopback-only
// configuration and reset per process invocation.
var uiListenHost = "127.0.0.1"
var uiNoBrowser = false

func init() {
	cmdssh.RunSSHUIFn = RunUI
}

// RunUICmd dispatches gitmap ui and gitmap <module> ui.
func RunUICmd(args []string) error {
	checkHelp("ui", args)

	fs := flag.NewFlagSet("ui", flag.ContinueOnError)
	port := fs.Int("port", 8080, "HTTP listen port for the Web UI server")
	fs.IntVar(port, "p", 8080, "shorthand for --port")
	host := fs.String("host", "127.0.0.1", "bind IP address for the Web UI server (loopback only by default)")
	noBrowser := fs.Bool("no-browser", false, "do not automatically launch the system default web browser")

	// Flags may trail the page argument (gitmap ui settings --port 8080), so
	// hoist them ahead of the positional page name before parsing. This fixes
	// the old bug where --port itself became the page name.
	flagArgs, positional := splitUIFlags(args)
	if err := fs.Parse(flagArgs); err != nil {
		return err
	}

	page := "settings"
	if len(positional) > 0 {
		page = positional[0]
	}

	uiListenHost = *host
	uiNoBrowser = *noBrowser

	if !isLoopbackHost(uiListenHost) {
		fmt.Println("WARNING: binding a non-loopback address exposes the embedded terminal to the network.")
	}

	return RunUI(page, *port)
}

// valueFlagsUI are the ui flags that consume the following argument.
var valueFlagsUI = map[string]bool{"-p": true, "--port": true, "--host": true}

// splitUIFlags separates flag tokens (plus their values) from the positional
// page argument, regardless of order on the command line.
func splitUIFlags(args []string) (flagArgs, positional []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		name := a
		if idx := strings.Index(a, "="); idx >= 0 {
			name = a[:idx]
		}
		if len(a) > 1 && strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
			if valueFlagsUI[name] && !strings.Contains(a, "=") && i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
			continue
		}
		positional = append(positional, a)
	}
	return flagArgs, positional
}

func isLoopbackHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	return h == "127.0.0.1" || h == "::1" || h == "localhost"
}
