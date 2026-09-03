import os

svg_template = """<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 900 480" width="900" height="480">
  <defs>
    <linearGradient id="winGrad_{id}" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#1e1e2e"/>
      <stop offset="100%" stop-color="#11111b"/>
    </linearGradient>
    <linearGradient id="barGrad_{id}" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#313244"/>
      <stop offset="100%" stop-color="#181825"/>
    </linearGradient>
    <filter id="shadow_{id}" x="-5%" y="-5%" width="115%" height="120%">
      <feDropShadow dx="0" dy="10" stdDeviation="16" flood-color="#000000" flood-opacity="0.6"/>
    </filter>
  </defs>

  <rect width="900" height="480" fill="#0b0b14"/>

  <!-- Grid lines -->
  <g opacity="0.04">
    <line x1="0" y1="60" x2="900" y2="60" stroke="#cdd6f4" stroke-width="1"/>
    <line x1="0" y1="120" x2="900" y2="120" stroke="#cdd6f4" stroke-width="1"/>
    <line x1="0" y1="180" x2="900" y2="180" stroke="#cdd6f4" stroke-width="1"/>
    <line x1="0" y1="240" x2="900" y2="240" stroke="#cdd6f4" stroke-width="1"/>
    <line x1="0" y1="300" x2="900" y2="300" stroke="#cdd6f4" stroke-width="1"/>
    <line x1="0" y1="360" x2="900" y2="360" stroke="#cdd6f4" stroke-width="1"/>
    <line x1="0" y1="420" x2="900" y2="420" stroke="#cdd6f4" stroke-width="1"/>
  </g>

  <!-- Window Frame -->
  <g filter="url(#shadow_{id})">
    <rect x="25" y="20" width="850" height="440" rx="10" ry="10" fill="url(#winGrad_{id})" stroke="#45475a" stroke-width="1"/>
    <rect x="25" y="20" width="850" height="36" rx="10" ry="10" fill="url(#barGrad_{id})"/>
    <rect x="25" y="44" width="850" height="12" fill="url(#barGrad_{id})"/>
    <circle cx="50" cy="38" r="6" fill="#f38ba8"/>
    <circle cx="70" cy="38" r="6" fill="#f9e2af"/>
    <circle cx="90" cy="38" r="6" fill="#a6e3a1"/>
    <text x="450" y="42" text-anchor="middle" fill="#a6adc8" font-size="12" font-weight="600" font-family="Menlo, Monaco, Consolas, monospace">gitmap {title} — terminal session</text>
  </g>

  <!-- Output Lines -->
  <g font-family="Menlo, Monaco, Consolas, monospace" font-size="12">
    <!-- Prompt Line 1 -->
    <text x="48" y="86" fill="#a6e3a1">❯</text>
    <text x="64" y="86" fill="#89dceb">{dir}</text>
    <text x="{cmd1_x}" y="86" fill="#f5c2e7">$</text>
    <text x="{cmd1_text_x}" y="86" fill="#cdd6f4" font-weight="bold">{cmd1}</text>

{out1}

    <!-- Prompt Line 2 -->
    <text x="48" y="{p2_y}" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.45;0.48;1" dur="10s" repeatCount="indefinite"/>❯</text>
    <text x="64" y="{p2_y}" fill="#89dceb"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.45;0.48;1" dur="10s" repeatCount="indefinite"/>{dir}</text>
    <text x="{cmd1_x}" y="{p2_y}" fill="#f5c2e7"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.45;0.48;1" dur="10s" repeatCount="indefinite"/>$</text>
    <text x="{cmd1_text_x}" y="{p2_y}" fill="#cdd6f4" font-weight="bold"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.45;0.48;1" dur="10s" repeatCount="indefinite"/>{cmd2}</text>

{out2}
  </g>
</svg>
"""

items = {
    "agy": {
        "title": "agy",
        "dir": "~/.gemini",
        "cmd1_x": "135", "cmd1_text_x": "150",
        "cmd1": "gitmap agy ls show-projects-with-empty-conversations",
        "p2_y": "216",
        "cmd2": "gitmap agy remove-projects-with-empty-conversations --except \"prompts-alpha\" --dry-run",
        "out1": """    <text x="48" y="112" fill="#bac2de"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>PROJECT ID                             NAME             CONVERSATIONS  STATUS</text>
    <line x1="48" y1="120" x2="830" y2="120" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="140" fill="#f38ba8"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.15;1" dur="10s" repeatCount="indefinite"/>0349c4d0-5a91-4f3e-800f-81fd53fc724f   prompts-alpha    0 active       ⚠ EMPTY (0 turns)</text>
    <text x="48" y="160" fill="#f38ba8"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>1b75371b-19ab-4a0e-a60c-1786508aeb6c   test-runner      1 aborted      ⚠ EMPTY (init error)</text>
    <text x="48" y="184" fill="#f9e2af"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.18;0.22;1" dur="10s" repeatCount="indefinite"/>Found 2 empty project(s). Prune with: gitmap agy remove-projects-with-empty-conversations</text>""",
        "out2": """    <text x="48" y="242" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>[DRY-RUN] Evaluating 2 empty conversation projects against whitelist...</text>
    <text x="48" y="264" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>  ✓ Preserved (whitelist match): prompts-alpha (0349c4d0-5a91-4f3e-800f-81fd53fc724f)</text>
    <text x="48" y="286" fill="#fab387"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>  ● Would remove: test-runner (1b75371b-19ab-4a0e-a60c-1786508aeb6c)</text>
    <text x="48" y="310" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.62;0.66;1" dur="10s" repeatCount="indefinite"/>✓ Dry-run completed: 1 candidate identified for deletion, 1 protected.</text>"""
    },
    "db": {
        "title": "db",
        "dir": "~/gitmap",
        "cmd1_x": "130", "cmd1_text_x": "145",
        "cmd1": "gitmap db ls",
        "p2_y": "230",
        "cmd2": "gitmap db repo-db list",
        "out1": """    <text x="48" y="112" fill="#89b4fa"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>● Primary Master SQLite Database: <tspan fill="#a6e3a1">bin/data/gitmap.db</tspan> (1.4 MB) [WAL Enabled]</text>
    <text x="48" y="132" fill="#bac2de"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.12;1" dur="10s" repeatCount="indefinite"/>  Role: Global repository index, groups, bookmarks, profiles, and scan roots.</text>
    <text x="48" y="156" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.14;0.17;1" dur="10s" repeatCount="indefinite"/>● Split Search Databases: <tspan fill="#a6e3a1">bin/data/repo_search/*.db</tspan> (48 split files, 12.8 MB total)</text>
    <text x="48" y="176" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.17;0.20;1" dur="10s" repeatCount="indefinite"/>  Rationale: Lock elimination across parallel multi-goroutine directory scans.</text>
    <text x="48" y="196" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.20;0.24;1" dur="10s" repeatCount="indefinite"/>✓ Database integrity verified. Zero write-lock collisions detected.</text>""",
        "out2": """    <text x="48" y="256" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>SLUG                    DB FILE                     SIZE      FILES   CACHES  STATUS</text>
    <line x1="48" y1="264" x2="830" y2="264" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="284" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>gitmap-v28              repo_search/gitmap-v28.db   245 KB    340     12      <tspan fill="#a6e3a1">TRACKED</tspan></text>
    <text x="48" y="304" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>api-gateway             repo_search/api-gateway.db  180 KB    192     8       <tspan fill="#a6e3a1">TRACKED</tspan></text>
    <text x="48" y="328" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.62;0.66;1" dur="10s" repeatCount="indefinite"/>✓ 48 split repository databases registered and synced with master index.</text>"""
    },
    "find-duplicates": {
        "title": "find-duplicates",
        "dir": "~/work",
        "cmd1_x": "115", "cmd1_text_x": "130",
        "cmd1": "gitmap find-duplicates",
        "p2_y": "220",
        "cmd2": "gitmap agy optimize-projects",
        "out1": """    <text x="48" y="112" fill="#f9e2af"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>── Antigravity (AGY) Duplicate Projects ──</text>
    <text x="48" y="132" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.15;1" dur="10s" repeatCount="indefinite"/>Found 1 duplicate group (2 repeat JSON files pointing to D:\\wp-work\\riseup-asia\\gitmap)</text>
    <text x="48" y="156" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>● Remediation Command: <tspan fill="#a6e3a1">gitmap agy optimize-projects</tspan></text>
    <text x="48" y="180" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.18;0.22;1" dur="10s" repeatCount="indefinite"/>✓ Audited 4 platforms: AGY (1 dup), VS Code (0 dups), Chrome (0 dups), Git (0 dups)</text>""",
        "out2": """    <text x="48" y="246" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Deduplicating Antigravity workspaces across 1 duplicate path group...</text>
    <text x="48" y="268" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>  ✓ Kept newest: gitmap-v28 (0349c4d0-5a91-4f3e-800f-81fd53fc724f)</text>
    <text x="48" y="290" fill="#fab387"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>  ✓ Removed duplicate: 2c666070-443d-489a-8cb7-72e67d3e2859</text>
    <text x="48" y="314" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.62;0.66;1" dur="10s" repeatCount="indefinite"/>✓ Optimization complete: 1 duplicate removed, workspace registry clean.</text>"""
    },
    "vscode": {
        "title": "vscode",
        "dir": "~/work",
        "cmd1_x": "115", "cmd1_text_x": "130",
        "cmd1": "gitmap vscode sync",
        "p2_y": "210",
        "cmd2": "gitmap vscode optimize-projects",
        "out1": """    <text x="48" y="112" fill="#89b4fa"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>Scanning tracked repositories for VS Code Project Manager...</text>
    <text x="48" y="134" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.16;1" dur="10s" repeatCount="indefinite"/>✓ Synchronized 48 repositories into ~/AppData/Roaming/Code/User/projects.json</text>
    <text x="48" y="156" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.16;0.20;1" dur="10s" repeatCount="indefinite"/>  Added 3 new projects, updated 45 existing mappings.</text>""",
        "out2": """    <text x="48" y="236" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Checking projects.json for duplicate repository paths...</text>
    <text x="48" y="258" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.60;1" dur="10s" repeatCount="indefinite"/>✓ No redundant project entries found. Workspace list is optimized.</text>"""
    },
    "chrome": {
        "title": "chrome",
        "dir": "~/chrome",
        "cmd1_x": "125", "cmd1_text_x": "140",
        "cmd1": "gitmap chrome list",
        "p2_y": "220",
        "cmd2": "gitmap chrome copy Default BackupProfile --offline",
        "out1": """    <text x="48" y="112" fill="#bac2de"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>NAME              DIRECTORY              EXTENSIONS  BOOKMARKS  STATUS</text>
    <line x1="48" y1="120" x2="830" y2="120" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="140" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.15;1" dur="10s" repeatCount="indefinite"/>Default           User Data/Default      14          320        <tspan fill="#a6e3a1">ACTIVE</tspan></text>
    <text x="48" y="160" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>WorkDev           User Data/Profile 1    6           85         <tspan fill="#bac2de">IDLE</tspan></text>""",
        "out2": """    <text x="48" y="246" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Creating offline clone: 'Default' -> 'BackupProfile'...</text>
    <text x="48" y="268" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>  ✓ Copied preferences, feature flags, and bookmark trees.</text>
    <text x="48" y="290" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>✓ Offline profile created successfully. Ready for launch or merge.</text>"""
    },
    "cloning": {
        "title": "clone",
        "dir": "~/repos",
        "cmd1_x": "120", "cmd1_text_x": "135",
        "cmd1": "gitmap clone ./scan-output.json --safe-pull",
        "p2_y": "210",
        "cmd2": "gitmap clone-next --delete --ssh-key work",
        "out1": """    <text x="48" y="112" fill="#89b4fa"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>Cloning repositories from manifest with safe-pull retry protection...</text>
    <text x="48" y="134" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.15;1" dur="10s" repeatCount="indefinite"/>✓ [1/4] gitmap-v28: Pull updated (up-to-date)</text>
    <text x="48" y="154" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>✓ [2/4] cloud-worker: Cloned to ./cloud-worker</text>""",
        "out2": """    <text x="48" y="236" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Detected current version: v28. Cloning v29 from remote...</text>
    <text x="48" y="258" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>✓ Cloned next iteration: gitmap-v29 (authenticated via key: work)</text>
    <text x="48" y="280" fill="#fab387"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>✓ Purged prior directory: gitmap-v28 (--delete enabled)</text>"""
    },
    "git-ops": {
        "title": "git-ops",
        "dir": "~/repos",
        "cmd1_x": "120", "cmd1_text_x": "135",
        "cmd1": "gitmap pull --all",
        "p2_y": "220",
        "cmd2": "gitmap status",
        "out1": """    <text x="48" y="112" fill="#bac2de"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>REPO           BRANCH   LATEST   STATUS   SHA       DIAGNOSTICS</text>
    <line x1="48" y1="120" x2="830" y2="120" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="140" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.15;1" dur="10s" repeatCount="indefinite"/>gitmap-v28     main     main     ✓ OK     1b9bf00   Clean pull completed</text>
    <text x="48" y="160" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>api-core       main     main     ✓ OK     3e841fa   Clean pull completed</text>""",
        "out2": """    <text x="48" y="246" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>REPO           STATE    AHEAD/BEHIND  UNCOMMITTED  STASHES</text>
    <line x1="48" y1="254" x2="830" y2="254" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="274" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>gitmap-v28     ✓ Clean  0 / 0         0 files      0</text>
    <text x="48" y="294" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>api-core       ✓ Clean  0 / 0         0 files      0</text>"""
    },
    "release": {
        "title": "release",
        "dir": "~/gitmap",
        "cmd1_x": "130", "cmd1_text_x": "145",
        "cmd1": "gitmap release --bump minor --draft --compress",
        "p2_y": "220",
        "cmd2": "gitmap list-versions --limit 2",
        "out1": """    <text x="48" y="112" fill="#89b4fa"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>Executing GitMap Release Protocol [minor bump: v6.165.0 -> v6.166.0]...</text>
    <text x="48" y="134" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.15;1" dur="10s" repeatCount="indefinite"/>  ✓ Generated changelog between tags</text>
    <text x="48" y="154" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>  ✓ Cross-compiled binaries & compressed zip archives</text>
    <text x="48" y="174" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.18;0.22;1" dur="10s" repeatCount="indefinite"/>✓ Draft release v6.166.0 published to GitHub successfully.</text>""",
        "out2": """    <text x="48" y="246" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>TAG         COMMIT    DATE        NOTES</text>
    <line x1="48" y1="254" x2="830" y2="254" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="274" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>v6.166.0    3c365d7   2026-09-03  Antigravity empty conversation pruner</text>
    <text x="48" y="294" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.58;0.62;1" dur="10s" repeatCount="indefinite"/>v6.165.0    8ab12e4   2026-09-02  Database suite & cross-platform deduplication</text>"""
    },
    "navigation": {
        "title": "navigation",
        "dir": "~",
        "cmd1_x": "75", "cmd1_text_x": "90",
        "cmd1": "gitmap cd gitmap-v28",
        "p2_y": "200",
        "cmd2": "gitmap group create backend api auth worker",
        "out1": """    <text x="48" y="112" fill="#a6e3a1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>✓ Matched slug 'gitmap-v28' -> D:\\wp-work\\riseup-asia\\gitmap</text>
    <text x="48" y="132" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.16;1" dur="10s" repeatCount="indefinite"/>Switched directory successfully.</text>""",
        "out2": """    <text x="48" y="226" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Creating repository group 'backend' with 3 repositories...</text>
    <text x="48" y="248" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>✓ Group 'backend' created. Execute batch commands with: gitmap group backend pull</text>"""
    },
    "schedule": {
        "title": "schedule",
        "dir": "~",
        "cmd1_x": "75", "cmd1_text_x": "90",
        "cmd1": "gitmap schedule add \"nightly-pull\" \"gitmap pull --all\" \"0 2 * * *\"",
        "p2_y": "200",
        "cmd2": "gitmap schedule list",
        "out1": """    <text x="48" y="112" fill="#a6e3a1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>✓ Enqueued scheduled job 'nightly-pull' (ID: task-8821) with cron: 0 2 * * *</text>""",
        "out2": """    <text x="48" y="226" fill="#bac2de"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>ID         NAME          SCHEDULE     STATUS    NEXT RUN</text>
    <line x1="48" y1="234" x2="830" y2="234" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="254" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>task-8821  nightly-pull  0 2 * * *    ENABLED   Tomorrow 02:00:00</text>"""
    },
    "cluster": {
        "title": "cluster",
        "dir": "~/cluster",
        "cmd1_x": "130", "cmd1_text_x": "145",
        "cmd1": "gitmap cluster nodes",
        "p2_y": "200",
        "cmd2": "gitmap serve --port 9999",
        "out1": """    <text x="48" y="112" fill="#bac2de"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>NODE ID        HOST          OS       ROLE       STATUS</text>
    <line x1="48" y1="120" x2="830" y2="120" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="140" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.16;1" dur="10s" repeatCount="indefinite"/>node-main      192.168.1.10  windows  server     CONNECTED</text>
    <text x="48" y="160" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>node-worker1   192.168.1.25  linux    client     CONNECTED</text>""",
        "out2": """    <text x="48" y="226" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Starting cluster orchestrator daemon on :9999...</text>
    <text x="48" y="248" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>✓ Orchestrator listening. Join token: gm_tok_9b2e88a104</text>"""
    },
    "templates": {
        "title": "templates",
        "dir": "~/new-project",
        "cmd1_x": "155", "cmd1_text_x": "170",
        "cmd1": "gitmap templates init go python --lfs",
        "p2_y": "200",
        "cmd2": "gitmap commons",
        "out1": """    <text x="48" y="112" fill="#89b4fa"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>Scaffolding curated .gitignore and .gitattributes for: go, python...</text>
    <text x="48" y="134" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.16;1" dur="10s" repeatCount="indefinite"/>✓ Created .gitignore with preserved marker-blocks</text>
    <text x="48" y="154" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>✓ Configured Git LFS binary file patterns in .gitattributes</text>""",
        "out2": """    <text x="48" y="226" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>✓ Union-merged all standard configs: ignore, attributes, prettier, lfs</text>"""
    },
    "data": {
        "title": "data",
        "dir": "~/gitmap",
        "cmd1_x": "130", "cmd1_text_x": "145",
        "cmd1": "gitmap profile create work-env",
        "p2_y": "200",
        "cmd2": "gitmap bookmark save pull-ci \"gitmap pull --all\"",
        "out1": """    <text x="48" y="112" fill="#a6e3a1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>✓ Profile 'work-env' created and switched. Isolated database ready.</text>""",
        "out2": """    <text x="48" y="226" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>✓ Saved bookmark 'pull-ci' -> 'gitmap pull --all'. Replay with: gitmap bookmark run pull-ci</text>"""
    },
    "automation": {
        "title": "automation",
        "dir": "~/gitmap",
        "cmd1_x": "130", "cmd1_text_x": "145",
        "cmd1": "gitmap installer ls",
        "p2_y": "210",
        "cmd2": "gitmap install neovim",
        "out1": """    <text x="48" y="112" fill="#bac2de"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>TOOL SLUG    NAME          CATEGORY    PLATFORM</text>
    <line x1="48" y1="120" x2="830" y2="120" stroke="#45475a" stroke-width="1"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/></line>
    <text x="48" y="140" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.16;1" dur="10s" repeatCount="indefinite"/>neovim       Neovim IDE    editor      windows, darwin, linux</text>
    <text x="48" y="160" fill="#cdd6f4"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>ripgrep      Ripgrep Tool  search      windows, darwin, linux</text>""",
        "out2": """    <text x="48" y="236" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Running verified installer for 'neovim'...</text>
    <text x="48" y="258" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>✓ Neovim installed and added to PATH successfully.</text>"""
    },
    "utilities": {
        "title": "utilities",
        "dir": "~/gitmap",
        "cmd1_x": "130", "cmd1_text_x": "145",
        "cmd1": "gitmap doctor --fix-path",
        "p2_y": "210",
        "cmd2": "gitmap fix-repo -2 --dry-run",
        "out1": """    <text x="48" y="112" fill="#89b4fa"><animate attributeName="opacity" values="0;1;1" keyTimes="0;0.1;1" dur="10s" repeatCount="indefinite"/>Diagnosing environment and GitMap installation health...</text>
    <text x="48" y="134" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.12;0.16;1" dur="10s" repeatCount="indefinite"/>✓ Git binary: C:\\Program Files\\Git\\cmd\\git.exe</text>
    <text x="48" y="154" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.15;0.18;1" dur="10s" repeatCount="indefinite"/>✓ Database permissions: Read/Write OK</text>
    <text x="48" y="174" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.18;0.22;1" dur="10s" repeatCount="indefinite"/>✓ Fixed PATH entry: added C:\\Users\\Alim\\bin</text>""",
        "out2": """    <text x="48" y="236" fill="#89b4fa"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.52;0.55;1" dur="10s" repeatCount="indefinite"/>Previewing token rewrites for prior 2 versions (v26..v27 -> v28)...</text>
    <text x="48" y="258" fill="#a6e3a1"><animate attributeName="opacity" values="0;0;1;1" keyTimes="0;0.55;0.58;1" dur="10s" repeatCount="indefinite"/>✓ Dry-run: 0 stale version tokens require rewrite.</text>"""
    }
}

out_dir = "docs/assets"
os.makedirs(out_dir, exist_ok=True)

for key, data in items.items():
    svg_content = svg_template.format(
        id=key,
        title=data["title"],
        dir=data["dir"],
        cmd1=data["cmd1"],
        cmd1_x=data["cmd1_x"],
        cmd1_text_x=data["cmd1_text_x"],
        out1=data["out1"],
        p2_y=data["p2_y"],
        cmd2=data["cmd2"],
        out2=data["out2"]
    )
    path = os.path.join(out_dir, f"{key}.svg")
    with open(path, "w", encoding="utf-8") as f:
        f.write(svg_content.strip() + "\n")
    print(f"Generated {path}")

print("All 15 animated SVGs generated successfully!")
