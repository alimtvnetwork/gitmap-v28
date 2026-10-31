package cmdcdp

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func initCDPFlagSet(opts *CDPOptions) *flag.FlagSet {
	fs := flag.NewFlagSet("cdp", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.IntVar(&opts.Port, "port", 0, "CDP port")
	fs.IntVar(&opts.Port, "p", 0, "CDP port")
	fs.StringVar(&opts.TargetID, "target", "", "Target ID")
	fs.StringVar(&opts.TargetID, "t", "", "Target ID")
	fs.BoolVar(&opts.IsJSON, "json", false, "Output JSON")
	fs.IntVar(&opts.TimeoutSec, "timeout", 5, "Timeout in seconds")

	return fs
}

func parseCDPFlags(args []string) (CDPOptions, []string, error) {
	opts := CDPOptions{TimeoutSec: 5}
	fs := initCDPFlagSet(&opts)
	if err := fs.Parse(args); err != nil {
		return opts, nil, apperror.WrapSimple(err, "cdp.flags")
	}

	return opts, fs.Args(), nil
}

func extractSubcommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", args
	}

	if strings.HasPrefix(args[0], "-") {
		return "", args
	}

	return strings.ToLower(args[0]), args[1:]
}

func resolveSubcommandAndOptions(args []string) (string, CDPOptions, []string, error) {
	subcmd, tail := extractSubcommand(args)
	opts, rest, err := parseCDPFlags(tail)
	if err != nil {
		return "", opts, nil, err
	}

	if subcmd == "" && len(rest) > 0 {
		subcmd = strings.ToLower(rest[0])
		rest = rest[1:]
	}

	return subcmd, opts, rest, nil
}

func resolvePort(port int) int {
	if port > 0 {
		return port
	}

	return DiscoverDefaultPort()
}

func findTargetWebSocketByID(port int, targetID string) (string, error) {
	targets, err := DiscoverTargets(port)
	if err != nil {
		return "", err
	}

	for i := range targets {
		isMatch := targets[i].ID == targetID || strings.Contains(targets[i].URL, targetID)
		if isMatch {
			return targets[i].WebSocketDebuggerURL, nil
		}
	}

	return "", fmt.Errorf("target %q not found on port %d", targetID, port)
}

func resolveWebSocketURL(port int, targetID string) (string, error) {
	if targetID != "" {
		return findTargetWebSocketByID(port, targetID)
	}

	target, err := FindActivePageTarget(port)
	if err != nil {
		return "", err
	}

	return target.WebSocketDebuggerURL, nil
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(v)
}

func printRawJSON(raw json.RawMessage) error {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		fmt.Println(string(raw))

		return nil
	}

	fmt.Println(pretty.String())

	return nil
}

func renderSingleTarget(t CDPTarget) {
	tag := termout.TokyoPill(t.Type, termout.SafeTokyoStyle(termout.TokyoOK))
	title := termout.TokyoAccentStyle.Render(t.Title)
	fmt.Printf("  %s %s\n", tag, title)
	fmt.Printf("    ID:  %s\n", t.ID)
	fmt.Printf("    URL: %s\n\n", t.URL)
}

func renderTargetsTable(port int, targets []CDPTarget) {
	hdr := termout.TokyoHeaderStyle.Render(fmt.Sprintf("CDP Targets on port %d (%d found)", port, len(targets)))
	fmt.Println("\n" + hdr + "\n")
	for _, t := range targets {
		renderSingleTarget(t)
	}
}

func runCDPTargets(opts CDPOptions) error {
	port := resolvePort(opts.Port)
	targets, err := DiscoverTargets(port)
	if err != nil {
		return err
	}

	if opts.IsJSON {
		return printJSON(targets)
	}

	renderTargetsTable(port, targets)

	return nil
}

func fetchProbeResult(opts CDPOptions) (*ProbeResult, int, error) {
	port := resolvePort(opts.Port)
	wsURL, err := resolveWebSocketURL(port, opts.TargetID)
	if err != nil {
		return nil, 0, err
	}

	res, err := RunDOMInputProbe(wsURL)

	return res, port, err
}

func renderProbeHeader(port int, res *ProbeResult) {
	hdr := termout.TokyoHeaderStyle.Render(fmt.Sprintf("DOM Probe Results (port %d)", port))
	fmt.Println("\n" + hdr)
	fmt.Printf("Active Element: %s\n", termout.TokyoAccentStyle.Render(res.Active))
	fmt.Printf("Selector Ladder Hits: %d | ContentEditables: %d\n\n", len(res.Found), len(res.AllContentEditables))
}

func renderHitDetails(hit ProbeHit) {
	if hit.Placeholder != "" {
		fmt.Printf("    Placeholder: %s\n", hit.Placeholder)
	}

	if hit.AriaLabel != "" {
		fmt.Printf("    Aria-Label:  %s\n", hit.AriaLabel)
	}

	if hit.Text != "" {
		fmt.Printf("    Text:        %s\n", hit.Text)
	}
}

func renderSingleProbeHit(hit ProbeHit) {
	pill := termout.TokyoPill(hit.Tag, termout.SafeTokyoStyle(termout.TokyoAccent))
	fmt.Printf("  %s %s\n", pill, termout.TokyoHeaderStyle.Render(hit.Selector))
	renderHitDetails(hit)
	fmt.Println()
}

func renderProbeResult(port int, res *ProbeResult) {
	renderProbeHeader(port, res)
	for _, h := range res.Found {
		renderSingleProbeHit(h)
	}
}

func runCDPProbe(opts CDPOptions) error {
	res, port, err := fetchProbeResult(opts)
	if err != nil {
		return err
	}

	if opts.IsJSON {
		return printJSON(res)
	}

	renderProbeResult(port, res)

	return nil
}

func extractEvalExpression(rest []string) (string, error) {
	if len(rest) == 0 {
		return "", fmt.Errorf("eval requires a JavaScript expression argument")
	}

	return strings.Join(rest, " "), nil
}

func renderEvalResult(raw json.RawMessage) {
	hdr := termout.TokyoHeaderStyle.Render("CDP Evaluation Result:")
	fmt.Println(hdr)
	fmt.Println(termout.TokyoBoxStyle(string(raw)))
}

func executeAndPrintEval(wsURL, expr string, opts CDPOptions) error {
	raw, err := EvaluateJS(wsURL, expr, time.Duration(opts.TimeoutSec)*time.Second)
	if err != nil {
		return err
	}

	if opts.IsJSON {
		return printRawJSON(raw)
	}

	renderEvalResult(raw)

	return nil
}

func runCDPEval(opts CDPOptions, rest []string) error {
	expr, err := extractEvalExpression(rest)
	if err != nil {
		return err
	}

	port := resolvePort(opts.Port)
	wsURL, err := resolveWebSocketURL(port, opts.TargetID)
	if err != nil {
		return err
	}

	return executeAndPrintEval(wsURL, expr, opts)
}

func renderCDPHelp() {
	hdr := termout.TokyoHeaderStyle.Render("gitmap cdp - Native Chrome DevTools Protocol Engine")
	fmt.Println("\n" + hdr + "\n")
	fmt.Println("Usage:")
	fmt.Println("  gitmap cdp targets [-p <port>] [--json]")
	fmt.Println("  gitmap cdp probe   [-p <port>] [-t <target-id>] [--json]")
	fmt.Println("  gitmap cdp eval <expr> [-p <port>] [-t <target-id>] [--json]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -p, --port     CDP debugging port (auto-discovered if omitted)")
	fmt.Println("  -t, --target   Target ID or URL substring")
	fmt.Println("      --json     Output structured JSON")
	fmt.Println("      --timeout  Execution timeout in seconds (default: 5)")
	fmt.Println()
}

func dispatchCDPSubcommand(subcmd string, opts CDPOptions, rest []string) error {
	switch subcmd {
	case "targets", "list", "ls":
		return runCDPTargets(opts)
	case "probe":
		return runCDPProbe(opts)
	case "eval":
		return runCDPEval(opts, rest)
	}

	renderCDPHelp()

	return nil
}

// RunCDP dispatches the cdp subcommand, executing probe, eval, or targets inspection.
func RunCDP(args []string) error {
	subcmd, opts, rest, err := resolveSubcommandAndOptions(args)
	if err != nil {
		return err
	}

	return dispatchCDPSubcommand(subcmd, opts, rest)
}
