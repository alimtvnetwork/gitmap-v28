package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Config struct {
	Server    string
	Token     string
	Alias     string
	AutoStart bool
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install-startup":
			installStartup()
			return
		case "uninstall-startup":
			uninstallStartup()
			return
		case "status":
			printStatus()
			return
		}
	}

	serverFlag := flag.String("server", "127.0.0.1:9999", "Server address (host:port)")
	tokenFlag := flag.String("token", "", "Cluster join authorization token")
	aliasFlag := flag.String("alias", "", "Node alias / name")
	autoStartFlag := flag.Bool("auto-start", false, "Register in Windows Startup to run at boot")
	daemonFlag := flag.Bool("daemon", false, "Run in background daemon mode")
	flag.Parse()

	if *autoStartFlag {
		installStartup()
	}

	hostname, _ := os.Hostname()
	alias := *aliasFlag
	if alias == "" {
		alias = hostname
	}

	fmt.Printf("⚡ GitMap Node Agent starting (alias: %s, server: %s, daemon: %v)...\n",
		alias, *serverFlag, *daemonFlag)

	startHeartbeatLoop(*serverFlag, *tokenFlag, alias)
}

func startHeartbeatLoop(server, token, alias string) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s/api/v1/cluster/heartbeat?alias=%s&token=%s&os=%s",
		server, alias, token, runtime.GOOS)

	consecutiveFailures := 0
	for {
		resp, err := client.Get(url)
		if err != nil {
			consecutiveFailures = handleHeartbeatFailure(server, consecutiveFailures)
		} else {
			consecutiveFailures = handleHeartbeatSuccess(server, resp, consecutiveFailures)
		}
		time.Sleep(15 * time.Second)
	}
}

func handleHeartbeatFailure(server string, failures int) int {
	failures++
	if failures <= 3 || failures%10 == 0 {
		fmt.Printf("▲ [%s] Waiting for server %s (attempt %d)...\n",
			time.Now().Format("15:04:05"), server, failures)
	}
	return failures
}

func handleHeartbeatSuccess(server string, resp *http.Response, failures int) int {
	resp.Body.Close()
	if failures > 0 {
		fmt.Printf("✔ [%s] Connected to cluster orchestrator at %s\n",
			time.Now().Format("15:04:05"), server)
	}
	return 0
}

func installStartup() {
	if runtime.GOOS != "windows" {
		fmt.Println("Auto-startup registration is currently supported on Windows.")
		return
	}
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving executable: %v\n", err)
		return
	}
	cmd := exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`,
		"/v", "GitMapNodeJoin", "/t", "REG_SZ", "/d", fmt.Sprintf(`"%s" --daemon`, exe), "/f")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register startup: %v (%s)\n", err, string(out))
		return
	}
	fmt.Println("✔ Registered GitMap Node Join in Windows Startup.")
}

func uninstallStartup() {
	if runtime.GOOS != "windows" {
		return
	}
	cmd := exec.Command("reg", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`,
		"/v", "GitMapNodeJoin", "/f")
	_ = cmd.Run()
	fmt.Println("✔ Removed GitMap Node Join from Windows Startup.")
}

func printStatus() {
	fmt.Println("GitMap Node Agent Status:")
	fmt.Printf("  OS:       %s/%s\n", runtime.GOOS, runtime.GOARCH)
	host, _ := os.Hostname()
	fmt.Printf("  Hostname: %s\n", host)
}
