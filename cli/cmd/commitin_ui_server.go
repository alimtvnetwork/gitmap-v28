package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const commitInUIDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><title>GitMap Commit-Pull Migration Studio</title>
<style>
  body { background:#0d1117; color:#e6edf3; font-family:system-ui,sans-serif; margin:0; padding:24px; }
  h1,h2 { margin-top:0; color:#58a6ff; }
  .grid { display:grid; grid-template-columns:1fr 1fr; gap:20px; }
  .card { background:#161b22; border:1px solid #30363d; border-radius:8px; padding:16px; margin-bottom:16px; }
  input,select,textarea,button { background:#0d1117; color:#e6edf3; border:1px solid #30363d; border-radius:6px; padding:8px; width:100%; box-sizing:border-box; margin-bottom:8px; }
  button { background:#238636; border-color:#2ea043; cursor:pointer; font-weight:600; width:auto; padding:8px 16px; margin-right:8px; }
  button.sec { background:#21262d; border-color:#363b42; }
  .cmd-box { background:#030712; border:1px solid #1f2937; border-radius:6px; padding:12px; font-family:monospace; color:#38bdf8; word-break:break-all; }
  label { font-weight:bold; font-size:13px; color:#8b949e; display:block; margin-bottom:4px; }
  .checkbox-group { display:flex; gap:16px; flex-wrap:wrap; margin-bottom:12px; }
  .checkbox-group label { display:flex; align-items:center; gap:6px; cursor:pointer; color:#e6edf3; }
  .checkbox-group input { width:auto; margin:0; }
</style>
</head>
<body>
<h1>GitMap Commit-Pull Migration Studio</h1>
<div class="grid">
  <div>
    <div class="card">
      <h2>Migration Configuration</h2>
      <label>Target Repository Path</label>
      <input id="cfg-target" value="D:\test-gitmap\test-gitmap" oninput="updateCmd()">
      <label>Source Inputs (comma-separated or brace-range)</label>
      <textarea id="cfg-inputs" rows="3" oninput="updateCmd()">https://github.com/alimtvnetwork/git-repo-navigator, https://github.com/alimtvnetwork/gitmap-v{2..28}</textarea>
      <label>PR Replay Mode</label>
      <select id="cfg-prmode" onchange="updateCmd()">
        <option value="merges">merges (Feature branches & release PRs)</option>
        <option value="feature-per-commit">feature-per-commit</option>
        <option value="direct">direct</option>
      </select>
      <div class="checkbox-group">
        <label><input type="checkbox" id="chk-tree" checked onchange="updateCmd()"> --tree</label>
        <label><input type="checkbox" id="chk-sync" checked onchange="updateCmd()"> --final-sync</label>
        <label><input type="checkbox" id="chk-cd" checked onchange="updateCmd()"> --cd</label>
        <label><input type="checkbox" id="chk-recreate" checked onchange="updateCmd()"> --recreate</label>
        <label><input type="checkbox" id="chk-dryrun" onchange="updateCmd()"> --dry-run</label>
      </div>
      <button onclick="saveConfig()">Save commit-pull-config.json</button>
      <button class="sec" onclick="loadSample()">Reset to Default</button>
    </div>
    <div class="card">
      <h2>Generated CLI Command</h2>
      <div class="cmd-box" id="cmd-preview">gitmap commit-pull --config commit-pull-config.json</div>
    </div>
  </div>
  <div>
    <div class="card">
      <h2>Line Skippers & Optimization</h2>
      <label>Fast Prefix / Contains / Regex Stripping</label>
      <textarea id="cfg-skippers" rows="5" readonly>[
  { "mode": "starts_with", "pattern": "Co-authored-by: gpt-engineer" },
  { "mode": "starts_with", "pattern": "Signed-off-by: gpt-engineer" },
  { "mode": "starts_with", "pattern": "X-Lovable" },
  { "mode": "contains", "pattern": "lovable-edit-id" }
]</textarea>
    </div>
    <div class="card">
      <h2>Self-Contained SEO Templates & Variables</h2>
      <label>Imported from .ai-memory/temp/seo-templates.json</label>
      <div id="vars-preview" style="font-size:13px; color:#7ee787; font-family:monospace; line-height:1.6;">
        $COMPANY = Rise Up Asia LLC<br>
        $COMPANY_URL = https://riseup-asia.com<br>
        $REGIONS = California, New York, and Wyoming<br>
        $MAREK = Marek Flejszman (28+ years experience)<br>
        $ALIM = Alim Ul Karim (https://alimulkarim.com)<br>
      </div>
    </div>
  </div>
</div>
<script>
function updateCmd() {
  const isDry = document.getElementById('chk-dryrun').checked;
  const flag = isDry ? ' --dry-run' : '';
  document.getElementById('cmd-preview').innerText = 'gitmap commit-pull --config commit-pull-config.json' + flag;
}
function saveConfig() {
  const payload = {
    target: document.getElementById('cfg-target').value,
    inputs: document.getElementById('cfg-inputs').value.split(',').map(s => s.trim()).filter(Boolean),
    prMode: document.getElementById('cfg-prmode').value,
    tree: document.getElementById('chk-tree').checked,
    finalSync: document.getElementById('chk-sync').checked,
    cd: document.getElementById('chk-cd').checked,
    recreate: document.getElementById('chk-recreate').checked,
    dryRun: document.getElementById('chk-dryrun').checked
  };
  fetch('/api/save', { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(payload) })
    .then(r => r.json()).then(d => alert('Saved commit-pull-config.json successfully!'));
}
function loadSample() {
  document.getElementById('cfg-target').value = 'D:\\test-gitmap\\test-gitmap';
  document.getElementById('cfg-inputs').value = 'https://github.com/alimtvnetwork/git-repo-navigator, https://github.com/alimtvnetwork/gitmap-v{2..28}';
  document.getElementById('chk-dryrun').checked = false;
  updateCmd();
}
</script>
</body>
</html>`

func runCommitPullUI(args []string) error {
	fs := flag.NewFlagSet("commit-pull ui", flag.ContinueOnError)
	port := fs.Int("port", 8925, "Port for web UI")
	noOpen := fs.Bool("no-open", false, "Do not open browser automatically")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ln, err := listenCommitPullUIPort(*port)
	if err != nil {
		return fmt.Errorf("listen ui: %w", err)
	}
	return startCommitPullUIServer(ln, *noOpen)
}

func listenCommitPullUIPort(port int) (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err == nil {
		return ln, nil
	}
	return net.Listen("tcp", "127.0.0.1:0")
}

func startCommitPullUIServer(ln net.Listener, noOpen bool) error {
	addr := fmt.Sprintf("http://127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port)
	fmt.Printf("\n  🚀 GitMap Commit-Pull Migration Studio running at %s\n", addr)
	fmt.Println("  Press Ctrl+C to stop.")
	if !noOpen {
		openCommitPullUIBrowser(addr)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleCommitPullUIHome)
	mux.HandleFunc("/api/save", handleCommitPullUISave)
	mux.HandleFunc("/api/templates", handleCommitPullUITemplates)
	return http.Serve(ln, mux)
}

func handleCommitPullUIHome(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(commitInUIDashboardHTML))
}

func handleCommitPullUISave(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, _ := json.MarshalIndent(body, "", "  ")
	_ = os.WriteFile("commit-pull-config.json", data, 0o644)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func handleCommitPullUITemplates(w http.ResponseWriter, _ *http.Request) {
	db, err := store.OpenTemplatesSplitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	items, _ := db.ListTemplateItems("seo")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func openCommitPullUIBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	default:
		_ = exec.Command("xdg-open", url).Start()
	}
}
