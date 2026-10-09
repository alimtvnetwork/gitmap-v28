// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

type agyPromptOptions struct {
	nodeAlias  string
	project    string
	promptText string
	title      string
	isEnqueue  bool
	timeout    time.Duration
}

func parseAgyPromptArgs(args []string) (agyPromptOptions, error) {
	opts := agyPromptOptions{
		timeout: 30 * time.Second,
	}

	var positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--enqueue" || a == "-q" {
			opts.isEnqueue = true
			continue
		}
		if (a == "--title" || a == "-t") && i+1 < len(args) {
			opts.title = args[i+1]
			i++
			continue
		}
		if a == "--timeout" && i+1 < len(args) {
			opts.timeout = parseDurationQuiet(args[i+1], opts.timeout)
			i++
			continue
		}
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
		}
	}

	if len(positional) < 3 {
		return opts, apperror.NewSimple("usage: gitmap nodes agy prompt <node> <project> \"<prompt>\" [--title <title>] [--enqueue]", "E9057")
	}

	opts.nodeAlias = positional[0]
	opts.project = positional[1]
	opts.promptText = strings.Join(positional[2:], " ")

	if opts.title == "" {
		opts.title = fmt.Sprintf("Prompt for %s", opts.project)
	}

	return opts, nil
}

func parseDurationQuiet(val string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(val); err == nil {
		return d
	}
	return fallback
}

// RunNodesAgyPrompt dispatches a prompt directly to a specific project on a local or remote node.
func RunNodesAgyPrompt(args []string) error {
	opts, err := parseAgyPromptArgs(args)
	if err != nil {
		return err
	}

	if isLocalTargetNode(opts.nodeAlias) {
		return dispatchLocalPrompt(opts)
	}

	return dispatchRemotePrompt(opts)
}

func isLocalTargetNode(alias string) bool {
	if strings.EqualFold(alias, "local") || strings.EqualFold(alias, "localhost") || alias == "." {
		return true
	}
	host, _ := os.Hostname()
	return strings.EqualFold(alias, host)
}

func dispatchLocalPrompt(opts agyPromptOptions) error {
	start := time.Now()
	if opts.isEnqueue {
		return enqueueLocalPrompt(opts, start)
	}

	err := cmdagy.ExecuteSendPrompt(opts.project, opts.promptText, opts.title)
	if err != nil {
		return err
	}

	printPromptDispatchConfirmation("local", opts.project, opts.title, "dispatched", "-", time.Since(start))
	return nil
}

func enqueueLocalPrompt(opts agyPromptOptions, start time.Time) error {
	payload := cmdagy.PromptPayload{
		ProjectTarget: opts.project,
		Title:         opts.title,
		PromptText:    opts.promptText,
		IsEnqueue:     true,
	}
	rec, err := cmdagy.EnqueuePromptPayload(payload)
	if err != nil {
		return err
	}
	printPromptDispatchConfirmation("local", opts.project, opts.title, "queued", rec.ID, time.Since(start))
	return nil
}

func dispatchRemotePrompt(opts agyPromptOptions) error {
	conn, err := resolveSSHConnectionByAlias(opts.nodeAlias)
	if err != nil {
		return err
	}

	start := time.Now()
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return apperror.WrapSimple(errConnect, fmt.Sprintf("dial node %q", conn.Alias))
	}
	defer client.Close()

	destPromptFile := resolveRemotePromptTempPath(conn)
	promptBytes := []byte(opts.promptText)
	if errStream := cmdssh.StreamFileToRemote(client, destPromptFile, promptBytes, conn.OS); errStream != nil {
		return apperror.WrapSimple(errStream, fmt.Sprintf("stage prompt to node %q", conn.Alias))
	}

	remoteCmd := buildRemoteSendPromptCmd(conn, opts.project, destPromptFile, opts.title, opts.isEnqueue)
	shell := "sh"
	if isWindowsNode(conn) {
		shell = "ps"
	}

	out, errExec := secrets.RunCommand(client, remoteCmd, shell)
	if errExec != nil {
		return apperror.WrapSimple(errExec, fmt.Sprintf("execute remote prompt on %q: %s", conn.Alias, out))
	}

	status := "dispatched"
	if opts.isEnqueue {
		status = "queued"
	}

	printPromptDispatchConfirmation(conn.Alias, opts.project, opts.title, status, "-", time.Since(start))
	return nil
}

func resolveRemotePromptTempPath(conn db.SSHConnection) string {
	ts := time.Now().UnixNano()
	if isWindowsNode(conn) {
		return fmt.Sprintf(`C:\Windows\Temp\agy_prompt_%d.txt`, ts)
	}
	return fmt.Sprintf("/tmp/agy_prompt_%d.txt", ts)
}

func buildRemoteSendPromptCmd(conn db.SSHConnection, project, stagedFile, title string, isEnqueue bool) string {
	enqueueFlag := ""
	if isEnqueue {
		enqueueFlag = " --enqueue"
	}
	titleFlag := ""
	if title != "" {
		titleFlag = fmt.Sprintf(" --title %q", title)
	}

	if isWindowsNode(conn) {
		return fmt.Sprintf("gitmap agy send-prompt --project %q --prompt %q%s%s", project, stagedFile, titleFlag, enqueueFlag)
	}
	return fmt.Sprintf("gitmap agy send-prompt --project '%s' --prompt '%s'%s%s", project, stagedFile, titleFlag, enqueueFlag)
}

func printPromptDispatchConfirmation(node, project, title, status, queueID string, dur time.Duration) {
	fmt.Printf("\n  %s Successfully %s prompt to [%s] %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, status, node, project)
	fmt.Printf("    • Target Node:  %s\n", node)
	fmt.Printf("    • Target Proj:  %s\n", project)
	fmt.Printf("    • Prompt Title: %s\n", title)
	if queueID != "-" && queueID != "" {
		fmt.Printf("    • Queue ID:     %s\n", queueID)
	}
	fmt.Printf("    • Status:       %s\n", strings.ToUpper(status))
	fmt.Printf("    • Duration:     %v\n\n", dur.Round(time.Millisecond))
}
