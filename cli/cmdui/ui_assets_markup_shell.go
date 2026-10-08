package cmdui

const uiAssetsMarkupShell = `<body>
  <div id="sidebar">
    <div class="logo">⚡ GitMap Fleet UI</div>
    <nav>
      <button class="nav-btn active" onclick="showTab('settings')">⚙️ Settings</button>
      <button class="nav-btn" onclick="showTab('commitin')">📦 Commitin</button>
      <button class="nav-btn" onclick="showTab('ssh')">🖥️ SSH Fleet</button>
      <button class="nav-btn" onclick="showTab('macro')">📜 Macros</button>
      <button class="nav-btn" onclick="showTab('installer')">🚀 Installers</button>
      <button class="nav-btn" onclick="showTab('prompts')">💬 Prompts & AI</button>
      <button class="nav-btn" onclick="showTab('import-export')">🔄 Import / Export</button>
      <button class="nav-btn" onclick="showTab('schedules')">⏱️ Schedules</button>
      <button class="nav-btn" onclick="showTab('editor')">📝 Remote Editor</button>
      <button class="nav-btn" onclick="showTab('help')">📖 CLI & AGY Help</button>
    </nav>
  </div>
  <div id="main">
    <header>
      <h2 id="header-title">Settings</h2>
      <div id="status-indicator" style="font-size: 0.85rem; color: var(--muted);">Connected: localhost</div>
    </header>

    <!-- SETTINGS TAB -->
`
