package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const templatesUIDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><title>GitMap Templates & Variables Studio</title>
<style>
  body { background:#0d1117; color:#e6edf3; font-family:system-ui,sans-serif; margin:0; padding:24px; }
  h1,h2 { margin-top:0; color:#58a6ff; }
  .grid { display:grid; grid-template-columns:320px 1fr; gap:20px; }
  .card { background:#161b22; border:1px solid #30363d; border-radius:8px; padding:16px; margin-bottom:16px; }
  input,select,textarea,button { background:#0d1117; color:#e6edf3; border:1px solid #30363d; border-radius:6px; padding:8px; width:100%; box-sizing:border-box; margin-bottom:8px; }
  button { background:#238636; border-color:#2ea043; cursor:pointer; font-weight:600; width:auto; padding:8px 14px; }
  button.del { background:#da3633; border-color:#f85149; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th,td { border-bottom:1px solid #30363d; padding:8px; text-align:left; vertical-align:top; }
  code { color:#79c0ff; }
</style>
</head>
<body>
<h1>GitMap State Templates & Variables Studio</h1>
<div class="grid">
  <div>
    <div class="card">
      <h2>Add / Update Template</h2>
      <input id="t-id" placeholder="ID (optional, e.g. seo-01)">
      <input id="t-cat" placeholder="Category (seo, prompts, ui-ux, prefix)" value="seo">
      <input id="t-slug" placeholder="Slug (e.g. why-canonical)">
      <input id="t-title" placeholder="# Why question title?">
      <textarea id="t-text" rows="5" placeholder="Because ..."></textarea>
      <button onclick="saveTemplate()">Save Template</button>
    </div>
    <div class="card">
      <h2>Variables ($VAR)</h2>
      <input id="v-key" placeholder="Variable Key (e.g. BRAND)">
      <input id="v-val" placeholder="Variable Value">
      <button onclick="saveVar()">Set Variable</button>
      <div id="vars-list" style="margin-top:10px;"></div>
    </div>
  </div>
  <div>
    <div class="card">
      <h2>Templates State DB (gitmap-templates.db)</h2>
      <button onclick="exportJSON()">Export JSON</button>
      <table><thead><tr><th>ID / Slug</th><th>Category</th><th>Title & Reasoning</th><th>Action</th></tr></thead>
      <tbody id="tpl-rows"></tbody></table>
    </div>
  </div>
</div>
<script>
async function loadState() {
  const res = await fetch('/api/state');
  const data = await res.json();
  const tbody = document.getElementById('tpl-rows');
  tbody.innerHTML = (data.templates || []).map(t =>
    '<tr><td><code>'+t.id+'</code><br><small>'+t.slug+'</small></td><td>'+t.category+'</td><td><b>'+t.title+'</b><br><pre style="white-space:pre-wrap;margin:4px 0;">'+t.text+'</pre></td>' +
    '<td><button class="del" onclick="delTpl(\''+t.id+'\')">Delete</button></td></tr>'
  ).join('');
  const vdiv = document.getElementById('vars-list');
  const vars = data.variables || {};
  vdiv.innerHTML = Object.keys(vars).map(k => '<div><code>$'+k+'</code> = '+vars[k]+'</div>').join('');
}
async function saveTemplate() {
  await fetch('/api/templates', {method:'POST', body:JSON.stringify({
    id:document.getElementById('t-id').value, category:document.getElementById('t-cat').value,
    slug:document.getElementById('t-slug').value, title:document.getElementById('t-title').value,
    text:document.getElementById('t-text').value
  })});
  loadState();
}
async function delTpl(id) {
  await fetch('/api/templates?id='+encodeURIComponent(id), {method:'DELETE'});
  loadState();
}
async function saveVar() {
  await fetch('/api/variables', {method:'POST', body:JSON.stringify({
    key:document.getElementById('v-key').value, value:document.getElementById('v-val').value
  })});
  loadState();
}
function exportJSON() { window.location.href = '/api/export'; }
loadState();
</script>
</body></html>`

func runTemplatesStateUI(args []string) {
	fs := flag.NewFlagSet("templates-ui", flag.ExitOnError)
	port := fs.Int("port", 8787, "HTTP listen port")
	isNoOpen := fs.Bool("no-open", false, "Do not auto-open browser")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	ln, url, err := listenTemplatesUIServer(*port)
	if err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	serveTemplatesUI(ln, url, !*isNoOpen)
}

func listenTemplatesUIServer(port int) (net.Listener, string, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
	}
	if err != nil {
		return nil, "", apperror.WrapSimple(err, "cmd.templates.ui_listen")
	}

	url := "http://" + ln.Addr().String()

	return ln, url, nil
}

func serveTemplatesUI(ln net.Listener, url string, isAutoOpen bool) {
	fmt.Printf("GitMap Templates UI listening at %s\n", url)
	if isAutoOpen {
		openTemplatesUIBrowser(url)
	}

	mux := newTemplatesUIMux()
	_ = http.Serve(ln, mux)
}

func newTemplatesUIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleTemplatesUIIndex)
	mux.HandleFunc("/api/state", handleTemplatesUIState)
	mux.HandleFunc("/api/templates", handleTemplatesUIItems)
	mux.HandleFunc("/api/variables", handleTemplatesUIVars)
	mux.HandleFunc("/api/export", handleTemplatesUIExport)

	return mux
}

func handleTemplatesUIIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(templatesUIDashboardHTML))
}

func handleTemplatesUIState(w http.ResponseWriter, _ *http.Request) {
	db, err := store.OpenTemplatesSplitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer db.Close()

	cats, _ := db.ListCategories()
	items, _ := db.ListTemplateItems("")
	vars, _ := db.ListVariables("")
	writeUIJSON(w, map[string]any{"categories": cats, "templates": items, "variables": vars})
}

func handleTemplatesUIItems(w http.ResponseWriter, r *http.Request) {
	db, err := store.OpenTemplatesSplitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer db.Close()

	if r.Method == http.MethodDelete {
		_ = db.DeleteTemplateItem(r.URL.Query().Get("id"))
		writeUIJSON(w, map[string]bool{"ok": true})

		return
	}

	upsertUITemplateItem(w, r, db)
}

func upsertUITemplateItem(w http.ResponseWriter, r *http.Request, db *store.TemplatesSplitDB) {
	var item store.StateTemplateItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	saved, err := db.UpsertTemplateItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeUIJSON(w, saved)
}

func handleTemplatesUIVars(w http.ResponseWriter, r *http.Request) {
	db, err := store.OpenTemplatesSplitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer db.Close()

	var body store.StateTemplateVariable
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.Key != "" {
		_ = db.UpsertVariable(body.Key, body.Scope, body.Value, body.Description)
	}

	writeUIJSON(w, map[string]bool{"ok": true})
}

func handleTemplatesUIExport(w http.ResponseWriter, _ *http.Request) {
	db, err := store.OpenTemplatesSplitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer db.Close()

	cats, _ := db.ListCategories()
	items, _ := db.ListTemplateItems("")
	vars, _ := db.ListVariables("")
	payload := store.TemplateExportPayload{Version: "1.0", Categories: cats, Variables: vars, Templates: items}
	payload.ExportID = store.ComputeExportHashID(payload)
	writeUIJSON(w, payload)
}

func writeUIJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func openTemplatesUIBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	default:
		_ = exec.Command("xdg-open", url).Start()
	}
}
