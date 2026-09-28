package cmdagy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const sugUIDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GitMap — Shutdown Until Green (SUG) Studio</title>
<style>
  :root {
    --bg: #0d1117;
    --surface: #161b22;
    --border: #30363d;
    --text: #e6edf3;
    --dim: #8b949e;
    --cyan: #58a6ff;
    --green: #2ea043;
    --green-bg: rgba(46, 160, 67, 0.15);
    --yellow: #d29922;
    --yellow-bg: rgba(210, 153, 34, 0.15);
    --red: #f85149;
    --red-bg: rgba(248, 81, 73, 0.15);
  }
  body { background: var(--bg); color: var(--text); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, monospace; margin: 0; padding: 24px; }
  .container { max-width: 960px; margin: 0 auto; }
  header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; border-bottom: 1px solid var(--border); padding-bottom: 16px; }
  h1 { font-size: 20px; font-weight: 700; color: var(--cyan); margin: 0; display: flex; align-items: center; gap: 8px; }
  .badge { font-size: 12px; font-weight: 600; padding: 4px 10px; border-radius: 12px; }
  .badge.running { background: var(--green-bg); color: var(--green); border: 1px solid var(--green); }
  .badge.idle { background: var(--yellow-bg); color: var(--yellow); border: 1px solid var(--yellow); }
  .card { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; padding: 16px; margin-bottom: 20px; }
  .btn-group { display: flex; gap: 10px; flex-wrap: wrap; margin-bottom: 16px; }
  button { background: var(--surface); color: var(--text); border: 1px solid var(--border); border-radius: 6px; padding: 8px 14px; font-size: 13px; font-weight: 600; cursor: pointer; transition: 0.15s; }
  button:hover { border-color: var(--cyan); }
  button.primary { background: #238636; border-color: #2ea043; }
  button.primary:hover { background: #2ea043; }
  button.danger { background: #da3633; border-color: #f85149; }
  .add-form { display: flex; gap: 8px; margin-top: 12px; }
  input[type="text"] { flex: 1; background: var(--bg); border: 1px solid var(--border); border-radius: 6px; padding: 8px 12px; color: var(--text); font-family: monospace; font-size: 13px; }
  input[type="text"]:focus { outline: none; border-color: var(--cyan); }
  table { width: 100%; border-collapse: collapse; margin-top: 12px; }
  th, td { padding: 10px 12px; text-align: left; border-bottom: 1px solid var(--border); font-size: 13px; }
  th { color: var(--dim); font-size: 11px; text-transform: uppercase; letter-spacing: 0.5px; }
  .status-pill { display: inline-flex; align-items: center; gap: 6px; font-weight: 600; font-size: 12px; padding: 2px 8px; border-radius: 10px; }
  .status-pill.green { background: var(--green-bg); color: var(--green); }
  .status-pill.pending { background: var(--yellow-bg); color: var(--yellow); }
  .action-col { text-align: right; }
  .btn-del { padding: 4px 8px; font-size: 11px; }
  .timer-banner { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: var(--dim); margin-top: 8px; }
</style>
</head>
<body>
<div class="container">
  <header>
    <div>
      <h1>⚡ GitMap Shutdown Until Green (SUG)</h1>
      <div style="font-size: 12px; color: var(--dim); margin-top: 4px;">Automated multi-project pipeline completion & system power governor</div>
    </div>
    <div id="runtime-status"><span class="badge idle">STATUS: EVALUATING...</span></div>
  </header>

  <div class="card">
    <div class="btn-group">
      <button class="primary" onclick="loadStatus()">↻ Refresh Pipelines</button>
      <button onclick="autoPopulateRunning()">+ Auto-Add Running AGY Projects</button>
      <button onclick="triggerOnceCheck()">Evaluate Pass Once</button>
    </div>
    <div class="timer-banner">
      <span id="target-count">Monitored Projects: 0</span>
      <span id="countdown">Auto-refresh in: 5m 00s</span>
    </div>
  </div>

  <div class="card">
    <h3 style="margin-top:0; font-size:14px;">Registered Watch List</h3>
    <table>
      <thead>
        <tr>
          <th>Target Project</th>
          <th>Status</th>
          <th class="action-col">Action</th>
        </tr>
      </thead>
      <tbody id="targets-body">
        <tr><td colspan="3" style="text-align:center; color:var(--dim); padding:24px;">Loading monitored projects...</td></tr>
      </tbody>
    </table>

    <div class="add-form">
      <input type="text" id="new-target" placeholder="Add local path, repository name, or Git URL (e.g. ./my-repo, antigravity-manager)...">
      <button class="primary" onclick="addNewTarget()">Add Target</button>
    </div>
  </div>
</div>

<script>
let remainingSeconds = 300;

async function loadStatus() {
  remainingSeconds = 300;
  try {
    const res = await fetch('/api/status');
    const data = await res.json();
    renderStatusBadge(data);
    renderTargetsTable(data);
  } catch (e) {
    console.error('Failed to load status:', e);
  }
}

function renderStatusBadge(data) {
  const el = document.getElementById('runtime-status');
  if (data.isLoopRunning) {
    el.innerHTML = '<span class="badge running">● WATCH LOOP ACTIVE (PID ' + data.pid + ')</span>';
  } else {
    el.innerHTML = '<span class="badge idle">○ WATCH LOOP IDLE</span>';
  }
  document.getElementById('target-count').innerText = 'Monitored Projects: ' + (data.targets ? data.targets.length : 0);
}

function renderTargetsTable(data) {
  const tbody = document.getElementById('targets-body');
  if (!data.targets || data.targets.length === 0) {
    tbody.innerHTML = '<tr><td colspan="3" style="text-align:center; color:var(--dim); padding:24px;">Watch list is empty. Add a target path or click "Auto-Add Running AGY Projects".</td></tr>';
    return;
  }
  let html = '';
  data.targets.forEach(t => {
    const pillClass = t.isGreen ? 'green' : 'pending';
    const pillText = t.isGreen ? '✔ GREEN' : '⏳ RUNNING/FAILING';
    html += '<tr>' +
      '<td><span style="font-family:monospace; font-weight:600;">' + escapeHtml(t.name) + '</span></td>' +
      '<td><span class="status-pill ' + pillClass + '">' + pillText + '</span></td>' +
      '<td class="action-col"><button class="danger btn-del" onclick="removeTarget(\'' + encodeURIComponent(t.name) + '\')">Remove</button></td>' +
      '</tr>';
  });
  tbody.innerHTML = html;
}

async function addNewTarget() {
  const input = document.getElementById('new-target');
  const target = input.value.trim();
  if (!target) return;
  const res = await fetch('/api/add?target=' + encodeURIComponent(target), { method: 'POST' });
  const data = await res.json();
  if (data.error) {
    alert(data.error);
  } else {
    input.value = '';
    loadStatus();
  }
}

async function removeTarget(target) {
  await fetch('/api/remove?target=' + target, { method: 'POST' });
  loadStatus();
}

async function autoPopulateRunning() {
  await fetch('/api/running-projects', { method: 'POST' });
  loadStatus();
}

async function triggerOnceCheck() {
  await fetch('/api/evaluate-once', { method: 'POST' });
  loadStatus();
}

function escapeHtml(s) {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

setInterval(() => {
  remainingSeconds--;
  if (remainingSeconds <= 0) {
    loadStatus();
  } else {
    const m = Math.floor(remainingSeconds / 60);
    const s = remainingSeconds % 60;
    document.getElementById('countdown').innerText = 'Auto-refresh in: ' + m + 'm ' + (s < 10 ? '0' : '') + s + 's';
  }
}, 1000);

loadStatus();
</script>
</body>
</html>`

// SUGTargetEvaluation summarizes an item's pipeline status for UI presentation.
type SUGTargetEvaluation struct {
	Name    string `json:"name"`
	IsGreen bool   `json:"isGreen"`
}

// SUGUIStatusPayload conveys state to the browser.
type SUGUIStatusPayload struct {
	IsLoopRunning   bool                  `json:"isLoopRunning"`
	PID             int                   `json:"pid"`
	IntervalSeconds int                   `json:"intervalSeconds"`
	Targets         []SUGTargetEvaluation `json:"targets"`
}

// RunSUGUI starts the local HTTP dark-mode studio dashboard.
func RunSUGUI(args []string) error {
	addr := "127.0.0.1:45199"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
	}
	if err != nil {
		return apperror.WrapSimple(err, "sug_ui_listen")
	}
	defer ln.Close()

	url := "http://" + ln.Addr().String()
	fmt.Printf("\n  %s✔ GitMap Shutdown Until Green UI started at %s%s\n", constants.ColorGreen, url, constants.ColorReset)
	fmt.Printf("  %sPress Ctrl+C to stop dashboard server%s\n\n", constants.ColorDim, constants.ColorReset)

	isNoOpen := false
	for _, a := range args {
		if a == "--no-open" {
			isNoOpen = true
		}
	}
	if !isNoOpen {
		openSUGUIBrowser(url)
	}

	mux := newSUGUIMux()
	return http.Serve(ln, mux)
}

func newSUGUIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleSUGUIIndex)
	mux.HandleFunc("/api/status", handleSUGUIStatus)
	mux.HandleFunc("/api/add", handleSUGUIAdd)
	mux.HandleFunc("/api/remove", handleSUGUIRemove)
	mux.HandleFunc("/api/running-projects", handleSUGUIRunningProjects)
	mux.HandleFunc("/api/evaluate-once", handleSUGUIEvaluateOnce)
	return mux
}

func handleSUGUIIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(sugUIDashboardHTML))
}

func handleSUGUIStatus(w http.ResponseWriter, _ *http.Request) {
	cfg := loadSUGConfig()
	st, isRunning := LoadSUGRuntimeStatus()

	evals := make([]SUGTargetEvaluation, 0, len(cfg.ProjectTargets))
	for _, p := range cfg.ProjectTargets {
		isGreen := checkSingleProjectPipelineGreen(p)
		evals = append(evals, SUGTargetEvaluation{Name: p, IsGreen: isGreen})
	}

	payload := SUGUIStatusPayload{
		IsLoopRunning:   isRunning,
		PID:             st.PID,
		IntervalSeconds: cfg.IntervalSec,
		Targets:         evals,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func handleSUGUIAdd(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimSpace(r.URL.Query().Get("target"))
	if target == "" {
		http.Error(w, `{"error":"target is required"}`, http.StatusBadRequest)
		return
	}
	_, isValid := ValidateSUGTarget(target)
	if !isValid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Target %q is not a valid directory, repository name, or Git URL", target),
		})
		return
	}
	cfg := loadSUGConfig()
	appendNewSUGTargets(&cfg, []string{target})
	_ = saveSUGConfig(cfg)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleSUGUIRemove(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimSpace(r.URL.Query().Get("target"))
	if target != "" {
		cfg := loadSUGConfig()
		cfg.ProjectTargets = filterKeptSUGProjects(cfg.ProjectTargets, []string{target})
		_ = saveSUGConfig(cfg)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleSUGUIRunningProjects(w http.ResponseWriter, _ *http.Request) {
	_ = setSUGRunningProjects()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleSUGUIEvaluateOnce(w http.ResponseWriter, _ *http.Request) {
	cfg := loadSUGConfig()
	for _, p := range cfg.ProjectTargets {
		_ = checkSingleProjectPipelineGreen(p)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func openSUGUIBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	default:
		_ = exec.Command("xdg-open", url).Start()
	}
}
