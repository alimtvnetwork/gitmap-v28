// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"encoding/json"
	"fmt"
	"html"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// settingEmailTaskTriggers is the Setting-table key holding the trigger
// list as a JSON blob. Defined here (not in constants) per subtask 06's
// contract — the sibling worker owns the shared constants files.
const settingEmailTaskTriggers = "email_task_triggers"

const (
	// triggerTimeFormat renders outcome timestamps in email bodies.
	triggerTimeFormat = "2006-01-02 15:04:05"
	// triggerSuccessSummary is the outcome summary for clean completions.
	triggerSuccessSummary = "completed successfully"
	// unknownMachineValue is the machine-info fallback — never an error.
	unknownMachineValue = "unknown"
	// outboundProbeAddr is dialed over UDP to discover the outbound IP;
	// connect() on UDP sends no traffic.
	outboundProbeAddr = "8.8.8.8:80"
	// statusOKText and statusFailText label the email status banner.
	statusOKText   = "SUCCESS"
	statusFailText = "FAILED"
	// bannerOKColor and bannerFailColor color the email status banner.
	bannerOKColor   = "#1a7f37"
	bannerFailColor = "#d1242f"
)

// EmailTrigger is one registered email notification: when task completes,
// an HTML status email is sent to To.
type EmailTrigger struct {
	Name string        `json:"name"`
	To   string        `json:"to"`
	Task EmailTaskType `json:"task"`
}

// TaskOutcome describes a finished gitmap operation for trigger emails.
type TaskOutcome struct {
	Repo       string `json:"repo"`
	Branch     string `json:"branch"`
	Operation  string `json:"operation"`
	IsSuccess  bool   `json:"is_success"`
	Summary    string `json:"summary"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

// MachineInfo carries host identity rows for the trigger email body.
type MachineInfo struct {
	Host string
	User string
	IP   string
}

// openTriggerStore opens the settings database for trigger persistence.
// A variable (not a direct call) so tests can point it at a temp DB.
var openTriggerStore = store.OpenDefault

// sendTriggerEmail sends one rendered trigger email; stubbed in tests.
var sendTriggerEmail = defaultSendTriggerEmail

// LoadEmailTriggers reads all registered triggers (empty when none).
// A corrupt blob surfaces as an error — it is trigger data, not
// unrelated state, so "behave as absent" would hide user config loss.
func LoadEmailTriggers(db *store.DB) ([]EmailTrigger, error) {
	raw := strings.TrimSpace(db.GetSetting(settingEmailTaskTriggers))
	if raw == "" {
		return []EmailTrigger{}, nil
	}

	var triggers []EmailTrigger
	if err := json.Unmarshal([]byte(raw), &triggers); err != nil {
		return nil, apperror.WrapSimple(err, "parse email task triggers:")
	}

	return triggers, nil
}

// SaveEmailTriggers persists the trigger list as a JSON blob.
func SaveEmailTriggers(db *store.DB, triggers []EmailTrigger) error {
	raw, err := json.Marshal(triggers)
	if err != nil {
		return apperror.WrapSimple(err, "encode email task triggers:")
	}

	if err := db.SetSetting(settingEmailTaskTriggers, string(raw)); err != nil {
		return apperror.WrapSimple(err, "save email task triggers:")
	}

	return nil
}

// RegisterEmailTrigger adds a trigger, replacing any existing trigger
// that shares its name.
func RegisterEmailTrigger(db *store.DB, trigger EmailTrigger) error {
	triggers, err := LoadEmailTriggers(db)
	if err != nil {
		return err
	}

	return SaveEmailTriggers(db, replaceOrAppendTrigger(triggers, trigger))
}

// replaceOrAppendTrigger swaps the same-named trigger or appends.
func replaceOrAppendTrigger(triggers []EmailTrigger, trigger EmailTrigger) []EmailTrigger {
	for i, existing := range triggers {
		if existing.Name == trigger.Name {
			triggers[i] = trigger

			return triggers
		}
	}

	return append(triggers, trigger)
}

// RemoveEmailTrigger deletes the trigger with the given name.
func RemoveEmailTrigger(db *store.DB, name string) error {
	triggers, err := LoadEmailTriggers(db)
	if err != nil {
		return err
	}

	kept := filterOutTrigger(triggers, name)
	if len(kept) == len(triggers) {
		return fmt.Errorf("no email trigger named %q", name)
	}

	return SaveEmailTriggers(db, kept)
}

// filterOutTrigger drops the named trigger from the list.
func filterOutTrigger(triggers []EmailTrigger, name string) []EmailTrigger {
	kept := make([]EmailTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		if trigger.Name != name {
			kept = append(kept, trigger)
		}
	}

	return kept
}

// TriggersForTask returns the triggers registered for one task type.
func TriggersForTask(triggers []EmailTrigger, task EmailTaskType) []EmailTrigger {
	var matched []EmailTrigger
	for _, trigger := range triggers {
		if trigger.Task == task {
			matched = append(matched, trigger)
		}
	}

	return matched
}

// InferTaskOutcome builds a TaskOutcome from an operation result,
// detecting repo and branch from the current git work tree. It never
// errors: every detection has a fallback.
func InferTaskOutcome(task EmailTaskType, err error) TaskOutcome {
	summary := triggerSuccessSummary
	isSuccess := true
	if err != nil {
		summary = err.Error()
		isSuccess = false
	}

	now := time.Now().UTC().Format(triggerTimeFormat)

	return TaskOutcome{
		Repo:       detectRepoName(),
		Branch:     detectBranchName(),
		Operation:  string(task),
		IsSuccess:  isSuccess,
		Summary:    summary,
		StartedAt:  now,
		FinishedAt: now,
	}
}

// detectRepoName returns the enclosing git repo's directory name, falling
// back to the cwd basename, then "unknown" — never an error.
func detectRepoName() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		return filepath.Base(strings.TrimSpace(string(out)))
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		return filepath.Base(cwd)
	}

	return unknownMachineValue
}

// detectBranchName returns the current branch, or "" outside a repo.
func detectBranchName() string {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// FireEmailTaskTriggers emails every trigger registered for the task.
// It never panics and never returns an error: a trigger must never break
// the operation that fired it — send failures go to stderr only.
func FireEmailTaskTriggers(task EmailTaskType, outcome TaskOutcome) {
	defer func() {
		_ = recover()
	}()

	fireEmailTaskTriggers(task, outcome)
}

// fireEmailTaskTriggers loads matching triggers and fires each one.
func fireEmailTaskTriggers(task EmailTaskType, outcome TaskOutcome) {
	db, err := openTriggerStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ email triggers skipped: cannot open settings db: %v\n", err)

		return
	}
	defer db.Close()

	triggers, loadErr := LoadEmailTriggers(db)
	if loadErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ email triggers skipped: %v\n", loadErr)

		return
	}

	info := CollectMachineInfo()
	for _, trigger := range TriggersForTask(triggers, task) {
		fireSingleTrigger(trigger, task, outcome, info)
	}
}

// fireSingleTrigger renders and sends one trigger email, reporting
// failures to stderr without disturbing the firing operation.
func fireSingleTrigger(trigger EmailTrigger, task EmailTaskType, outcome TaskOutcome, info MachineInfo) {
	subject := buildTriggerSubject(task, outcome)
	body := BuildTriggerEmailHTML(task, outcome, info)
	if err := sendTriggerEmail(trigger.To, subject, body); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ email trigger %q skipped: %v\n", trigger.Name, err)

		return
	}

	fmt.Printf("  ✉ email trigger %q notified %s\n", trigger.Name, trigger.To)
}

// buildTriggerSubject labels the notification with task + status.
func buildTriggerSubject(task EmailTaskType, outcome TaskOutcome) string {
	status := statusOKText
	if !outcome.IsSuccess {
		status = statusFailText
	}

	return fmt.Sprintf("[gitmap] %s %s: %s", task, status, outcome.Repo)
}

// CollectMachineInfo gathers host identity rows for trigger emails.
// It never errors: every field falls back to "unknown".
func CollectMachineInfo() MachineInfo {
	return MachineInfo{
		Host: machineHostname(),
		User: machineUsername(),
		IP:   machineOutboundIP(),
	}
}

// machineHostname returns os.Hostname() or "unknown".
func machineHostname() string {
	host, err := os.Hostname()
	if err == nil && len(strings.TrimSpace(host)) > 0 {
		return host
	}

	return unknownMachineValue
}

// machineUsername returns the OS login name or "unknown".
func machineUsername() string {
	current, err := user.Current()
	if err == nil && len(strings.TrimSpace(current.Username)) > 0 {
		return current.Username
	}

	envUser := strings.TrimSpace(os.Getenv("USER"))
	if len(envUser) > 0 {
		return envUser
	}

	return unknownMachineValue
}

// machineOutboundIP discovers the outbound IP via a UDP dial (no traffic
// sent), then the first non-loopback interface address, then "unknown".
func machineOutboundIP() string {
	if ip := outboundIPViaDial(); ip != "" {
		return ip
	}

	if ip := outboundIPViaInterfaces(); ip != "" {
		return ip
	}

	return unknownMachineValue
}

// outboundIPViaDial returns the local address of a UDP socket dialed at
// the probe address. Dialing UDP only resolves the route — no packet
// leaves the machine.
func outboundIPViaDial() string {
	conn, err := net.Dial("udp", outboundProbeAddr)
	if err != nil {
		return ""
	}
	defer conn.Close()

	udpAddr, isUDP := conn.LocalAddr().(*net.UDPAddr)
	if isUDP {
		return udpAddr.IP.String()
	}

	return ""
}

// outboundIPViaInterfaces returns the first non-loopback IPv4 address.
func outboundIPViaInterfaces() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if ip := firstUsableInterfaceIP(iface); ip != "" {
			return ip
		}
	}

	return ""
}

// firstUsableInterfaceIP returns the first non-loopback IPv4 on an
// interface, or "".
func firstUsableInterfaceIP(iface net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ip := usableIPFromAddr(addr); ip != "" {
			return ip
		}
	}

	return ""
}

// usableIPFromAddr extracts a non-loopback IPv4 from a net.Addr.
func usableIPFromAddr(addr net.Addr) string {
	ip := ipFromAddr(addr)
	if ip == nil || ip.IsLoopback() {
		return ""
	}

	flat := ip.To4()
	if flat == nil {
		return ""
	}

	return flat.String()
}

// ipFromAddr unwraps the IP from the two common net.Addr shapes.
func ipFromAddr(addr net.Addr) net.IP {
	switch typed := addr.(type) {
	case *net.IPNet:
		return typed.IP
	case *net.IPAddr:
		return typed.IP
	default:
		return nil
	}
}

// BuildTriggerEmailHTML renders the notification email: a green/red
// status banner plus a table of machine-info and outcome rows. Every
// interpolated value is HTML-escaped.
func BuildTriggerEmailHTML(task EmailTaskType, outcome TaskOutcome, info MachineInfo) string {
	color, status := triggerBannerStyle(outcome.IsSuccess)

	var sb strings.Builder
	sb.WriteString("<html><body style=\"font-family:sans-serif;max-width:640px;margin:0 auto;\">")
	sb.WriteString("<h2>gitmap " + html.EscapeString(string(task)) + " — " + html.EscapeString(status) + "</h2>")
	sb.WriteString("<div style=\"background:" + color + ";color:#ffffff;padding:12px 16px;border-radius:6px;\">")
	sb.WriteString("<strong>" + html.EscapeString(status) + "</strong> — gitmap " + html.EscapeString(outcome.Operation))
	sb.WriteString("</div>")
	sb.WriteString("<table style=\"margin-top:12px;border-collapse:collapse;width:100%;\">")
	writeTriggerRow(&sb, "Host", info.Host)
	writeTriggerRow(&sb, "User", info.User)
	writeTriggerRow(&sb, "Outbound IP", info.IP)
	writeTriggerRow(&sb, "Repo", outcome.Repo)
	writeTriggerRow(&sb, "Branch", outcome.Branch)
	writeTriggerRow(&sb, "Operation", outcome.Operation)
	writeTriggerRow(&sb, "Outcome", outcome.Summary)
	writeTriggerRow(&sb, "Started", outcome.StartedAt)
	writeTriggerRow(&sb, "Finished", outcome.FinishedAt)
	sb.WriteString("</table></body></html>")

	return sb.String()
}

// triggerBannerStyle picks the banner color and label per outcome.
func triggerBannerStyle(isSuccess bool) (color, label string) {
	if isSuccess {
		return bannerOKColor, statusOKText
	}

	return bannerFailColor, statusFailText
}

// writeTriggerRow appends one escaped label/value row to the email table.
func writeTriggerRow(sb *strings.Builder, label, value string) {
	sb.WriteString("<tr><td style=\"padding:6px 10px;border:1px solid #dddddd;font-weight:bold;\">")
	sb.WriteString(html.EscapeString(label))
	sb.WriteString("</td><td style=\"padding:6px 10px;border:1px solid #dddddd;\">")
	sb.WriteString(html.EscapeString(value))
	sb.WriteString("</td></tr>")
}
