# Specification 343: Improvements 631 to 660 (Interactive Bubble Tea Visual Mailbox Graph, ISP Feedback Loop Parser & Bounded Socket Latency Heatmap)

- **Sequence:** 343
- **Slug:** `343-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/03-improvements-631-to-660-interactive-bubbletea-and-inbox-analytics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (631 - 660)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 64** | Items 631–640 | High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator |
| **Cluster 65** | Items 641–650 | ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression |
| **Cluster 66** | Items 651–660 | Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI |

---

# Cluster 64: High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator (Items 631–640)

---

### Improvement 631: Visual Mailbox Topology Tree Navigator with Collapsible Node Hierarchy
- **ID:** `MAILBOX-GRAPH-631`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Render a real-time collapsible hierarchical tree in Bubble Tea representing multi-domain sender mailboxes, relay nodes, warming stages, and health statuses with keyboard navigation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Collapsible tree view (`+` expanded, `-` collapsed) displaying sender identity nodes, allocated IP relays, outbound queues, and warming status badges.
    - Lipgloss token palette: Active nodes (`#10b981`, emerald), Warming nodes (`#3b82f6`, blue), Stalled/Throttled (`#f59e0b`, amber), Quarantine/FBL Spike (`#ef4444`, red), Hierarchy connectors (`#64748b`, slate grey).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX TOPOLOGY TREE NAVIGATOR] ────────────────────────────────────────────────────────┐
      │ [-] domain: deliverability-core.io [Health: 98.4%] [Stage: Scaling]                        │
      │   ├─ [+] mbx: sender-01@deliverability-core.io [IP: 198.51.100.12] [Sent: 420/500] (READY) │
      │   ├─ [-] mbx: sender-02@deliverability-core.io [IP: 198.51.100.13] [Sent: 490/500] (WARM)  │
      │   │   ├─ Target: Google Workspace [Inbox: 99.1%] [Delay: 42ms]                            │
      │   │   ├─ Target: Microsoft 365    [Inbox: 97.4%] [Delay: 68ms]                            │
      │   │   └─ Target: Yahoo Mail       [Inbox: 96.0%] [Delay: 110ms]                           │
      │   └─ [+] mbx: sender-03@deliverability-core.io [IP: 198.51.100.14] [Sent: 150/200] (RAMP)  │
      │ [Enter: Expand/Collapse] [w: Focus Node] [p: Pause Node] [q: Return to Root]               │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline SVG directed tree canvas using SVG `<g>`, `<path>`, and `<text>` tags with SVG bezier curves (`M ... C ...`).
    - Embedded CSS transforms and drag-to-pan / wheel-zoom support with zero external JavaScript libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mailbox-tree`
  - Flags:
    - `--domain <string>`: Filter tree display by root domain identity (default: `all`).
    - `--depth <int>`: Hierarchy depth limit (`1`: domain, `2`: mailboxes, `3`: target ISP relays) (default: `3`).
    - `--show-throttled`: Only display nodes currently experiencing throttling or backpressure.
    - `--watch`: Persistent interactive Bubble Tea full-screen alternate buffer mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxTopologyNode struct {
        NodeID         string                `json:"nodeId"`
        ParentID       string                `json:"parentId"`
        DomainName     string                `json:"domainName"`
        MailboxAddress string                `json:"mailboxAddress"`
        AssignedIP     string                `json:"assignedIp"`
        HealthScore    float64               `json:"healthScore"`
        SentToday      int                   `json:"sentToday"`
        DailyQuota     int                   `json:"dailyQuota"`
        Children       []MailboxTopologyNode `json:"children"`
        IsExpanded     bool                  `json:"isExpanded"`
        IsHealthy      bool                  `json:"isHealthy"`
        HasActiveAlert bool                  `json:"hasActiveAlert"`
    }

    type MailboxTreeState struct {
        SelectedNodeID string                `json:"selectedNodeId"`
        RootNodes      []MailboxTopologyNode `json:"rootNodes"`
        TotalMailboxes int                   `json:"totalMailboxes"`
        IsLiveStream   bool                  `json:"isLiveStream"`
        HasFilter      bool                  `json:"hasFilter"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxTopologyNodes (
        nodeId TEXT PRIMARY KEY,
        parentId TEXT NOT NULL DEFAULT '',
        domainName TEXT NOT NULL,
        mailboxAddress TEXT NOT NULL,
        assignedIp TEXT NOT NULL,
        healthScore REAL NOT NULL DEFAULT 100.0,
        sentToday INTEGER NOT NULL DEFAULT 0,
        dailyQuota INTEGER NOT NULL DEFAULT 500,
        isExpanded INTEGER NOT NULL DEFAULT 1,
        isHealthy INTEGER NOT NULL DEFAULT 1,
        hasActiveAlert INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_mbx_node_parent ON MailboxTopologyNodes(parentId);
    CREATE INDEX IF NOT EXISTS idx_mbx_node_domain ON MailboxTopologyNodes(domainName);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/mailbox_tree.go` interfacing with `pkg/warmup/topology` and `pkg/mailboxtree/bubbletea_model.go`.

---

### Improvement 632: Interactive Node-Action Modal with Pause, Drain, and Reroute Capabilities
- **ID:** `MAILBOX-GRAPH-632`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Provide a contextual floating Bubble Tea modal dialog on top of the mailbox topology graph allowing operators to immediately Pause, Drain, or Reroute traffic from problematic mailboxes without leaving the terminal UI.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-contrast centered floating modal window with double-line border (`╔══╗`) overlaying the dimmed tree graph.
    - Lipgloss token palette: Modal border (`#8b5cf6`, purple), Accent selection (`#ec4899`, pink highlight), Safe action (`#10b981`, green), Destructive action (`#ef4444`, red).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX TOPOLOGY TREE NAVIGATOR] ────────────────────────────────────────────────────────┐
      │   ├─ [-] mbx: sender-02@deliverability-core.io [IP: 198.51.100.13]                         │
      │   │   ╔══════════════════ [ACTION: sender-02@deliverability-core.io] ══════════════════╗    │
      │   │   ║ Current Status: WARMING (490/500 quota) | Inactive Alert: None                 ║    │
      │   │   ║                                                                                ║    │
      │   │   ║  > [1] Pause Inbound & Outbound Traffic (Soft Hold)                            ║    │
      │   │   ║    [2] Drain Active Queue to Peer Mailbox (Zero Loss)                          ║    │
      │   │   ║    [3] Force Re-route Traffic to Pool: primary-relay-02                        ║    │
      │   │   ║    [4] Trigger Instant Health Re-Probe & Auth Audit (DKIM/SPF)                 ║    │
      │   │   ║                                                                                ║    │
      │   │   ║  [Press 1-4 or arrows + Enter to execute]  [Esc to Cancel]                     ║    │
      │   │   ╚════════════════════════════════════════════════════════════════════════════════╝    │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - HTML5 native `<dialog>` element with backdrop blur and CSS grid action menu buttons with SVG action icons.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup node-action`
  - Flags:
    - `--node-id <string>`: Target node or mailbox identifier (required).
    - `--action <string>`: Direct action execution (`pause`, `drain`, `reroute`, `probe`).
    - `--target-peer <string>`: Peer mailbox ID if draining or rerouting.
    - `--confirm`: Bypass confirmation prompt in non-interactive CI/CD runs.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type NodeActionRequest struct {
        NodeID       string `json:"nodeId"`
        ActionType   string `json:"actionType"`
        TargetPeerID string `json:"targetPeerId"`
        Reason       string `json:"reason"`
        IsImmediate  bool   `json:"isImmediate"`
        HasConfirmed bool   `json:"hasConfirmed"`
    }

    type NodeActionResult struct {
        ExecutionID  string    `json:"executionId"`
        NodeID       string    `json:"nodeId"`
        ActionType   string    `json:"actionType"`
        CompletedAt  time.Time `json:"completedAt"`
        IsSuccess    bool      `json:"isSuccess"`
        HasDrained   bool      `json:"hasDrained"`
        ErrorMessage string    `json:"errorMessage,omitempty"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxNodeActions (
        actionId TEXT PRIMARY KEY,
        nodeId TEXT NOT NULL,
        actionType TEXT NOT NULL,
        targetPeerId TEXT NOT NULL DEFAULT '',
        reason TEXT NOT NULL DEFAULT '',
        isImmediate INTEGER NOT NULL DEFAULT 1,
        isSuccess INTEGER NOT NULL DEFAULT 0,
        hasDrained INTEGER NOT NULL DEFAULT 0,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_node_actions_node ON MailboxNodeActions(nodeId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/mailbox_node_action.go` and connected to `pkg/warmup/drainer`.

---

### Improvement 633: Directed Acyclic Graph (DAG) Mailbox Relay Flow Visualizer
- **ID:** `MAILBOX-GRAPH-633`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Render a layered Directed Acyclic Graph (DAG) in terminal ANSI showing mail item movement from source warm-up generator through rate-limiter buckets to destination MX relays.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-stage horizontal pipeline DAG rendered using UTF-8 box drawing characters (`│`, `─`, `┌`, `└`, `►`).
    - Color tokens: Sources (`#38bdf8`, light blue), Token buckets (`#a855f7`, violet), Sockets (`#06b6d4`, cyan), Final MX (`#10b981`, green).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX DISPATCH DAG RELAY FLOW] ────────────────────────────────────────────────────────┐
      │ [WARMUP GEN] ──► [TOKEN BUCKET: 10/s] ──► [SOCKET POOL] ──► [MX: aspmx.l.google.com] (OK)  │
      │   (35 items)        (Cap: 100/100)          (8 conns)         [RTT: 38ms] [Sent: 342]      │
      │       │                                         │                                          │
      │       └────────► [TOKEN BUCKET: 5/s]  ──► [SOCKET POOL] ──► [MX: outlook.com] (SLOW)       │
      │                     (Cap: 25/50)            (3 conns)         [RTT: 184ms] [Sent: 148]     │
      │ [r: Refresh Rates] [s: Step Animation] [f: Freeze Stage] [q: Quit]                         │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline SVG DAG canvas with animated stroke dashes (`stroke-dashoffset`) visualizing simulated packet flow.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup relay-dag`
  - Flags:
    - `--filter-mx <string>`: Filter flow diagram by target MX destination host.
    - `--speed <string>`: Animation frame rate (`slow`, `normal`, `fast`) (default: `normal`).
    - `--json`: Output current DAG vertices and edges in JSON schema.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type RelayDAGVertex struct {
        VertexID     string   `json:"vertexId"`
        Label        string   `json:"label"`
        StageType    string   `json:"stageType"`
        QueueDepth   int      `json:"queueDepth"`
        CurrentRate  float64  `json:"currentRate"`
        IsSaturated  bool     `json:"isSaturated"`
        HasErrors    bool     `json:"hasErrors"`
    }

    type RelayDAGEdge struct {
        EdgeID       string  `json:"edgeId"`
        SourceID     string  `json:"sourceId"`
        TargetID     string  `json:"targetId"`
        Throughput   float64 `json:"throughput"`
        LatencyMs    float64 `json:"latencyMs"`
        IsCongested  bool    `json:"isCongested"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RelayDAGMetrics (
        metricId TEXT PRIMARY KEY,
        sourceId TEXT NOT NULL,
        targetId TEXT NOT NULL,
        throughput REAL NOT NULL DEFAULT 0.0,
        latencyMs REAL NOT NULL DEFAULT 0.0,
        isCongested INTEGER NOT NULL DEFAULT 0,
        recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_dag_edge ON RelayDAGMetrics(sourceId, targetId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/cmdemail/relay_dag.go` pulling from `pkg/warmup/dagengine`.

---

### Improvement 634: Mailbox Peer-to-Peer Warmup Reputation Closeness Matrix
- **ID:** `MAILBOX-GRAPH-634`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Display a compact terminal matrix detailing reciprocal interaction, message threading, and reputation closeness between warmup peer mailboxes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Cross-tabulated terminal matrix with row and column headers representing mailbox handles and color-coded score cells.
    - Color tokens: High closeness (`#10b981`, bold green `[99]`), Medium (`#3b82f6`, blue `[84]`), Low (`#eab308`, yellow `[42]`), Unpaired (`#475569`, grey `[--]`).
    - Model view snippet:
      ```text
      ┌─ [WARMUP PEER REPUTATION CLOSENESS MATRIX] ────────────────────────────────────────────────┐
      │ SENDER / RECEIVER │ MBX-01 │ MBX-02 │ MBX-03 │ MBX-04 │ MBX-05 │ RECIPROCITY │ THREAD DEPTH  │
      │ MBX-01 (alpha)    │  [--]  │  [98]  │  [94]  │  [82]  │  [91]  │    96.4%    │ 4.2 msgs/thd  │
      │ MBX-02 (beta)     │  [95]  │  [--]  │  [89]  │  [76]  │  [88]  │    92.1%    │ 3.8 msgs/thd  │
      │ MBX-03 (gamma)    │  [92]  │  [90]  │  [--]  │  [85]  │  [94]  │    94.5%    │ 4.0 msgs/thd  │
      │ [↑/↓/←/→: Select Pair] [Space: View Thread Details] [p: Export Matrix] [q: Quit]          │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure HTML table with inline CSS heatmap cell styling and dynamic tooltip popovers for interaction history.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup matrix`
  - Flags:
    - `--group <string>`: Mailbox pool or cluster name (default: `default`).
    - `--min-closeness <int>`: Filter connections below threshold (default: `0`).
    - `--format <string>`: Display format (`tui`, `csv`, `table`) (default: `tui`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PeerClosenessRecord struct {
        SourceMailboxID string  `json:"sourceMailboxId"`
        TargetMailboxID string  `json:"targetMailboxId"`
        ClosenessScore  int     `json:"closenessScore"`
        ExchangedCount  int     `json:"exchangedCount"`
        AverageThread   float64 `json:"averageThread"`
        IsReciprocal    bool    `json:"isReciprocal"`
        HasRecentBounce bool    `json:"hasRecentBounce"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupPeerMatrix (
        matrixEntryId TEXT PRIMARY KEY,
        sourceMailboxId TEXT NOT NULL,
        targetMailboxId TEXT NOT NULL,
        closenessScore INTEGER NOT NULL DEFAULT 0,
        exchangedCount INTEGER NOT NULL DEFAULT 0,
        averageThread REAL NOT NULL DEFAULT 1.0,
        isReciprocal INTEGER NOT NULL DEFAULT 1,
        hasRecentBounce INTEGER NOT NULL DEFAULT 0,
        lastInteraction DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_matrix_pair ON WarmupPeerMatrix(sourceMailboxId, targetMailboxId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured in `cli/cmdemail/peer_matrix.go` interfacing with `pkg/warmup/reputation`.

---

### Improvement 635: Fuzzy Interactive Search and Instant Mailbox Focus in Tree View
- **ID:** `MAILBOX-GRAPH-635`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Embed real-time fuzzy text search (`/` trigger) into the Bubble Tea mailbox tree navigator, auto-expanding parent nodes and highlighting matching substrings.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Bottom command line input strip featuring instant fuzzy matching across domain names, mailbox usernames, and outbound IP addresses.
    - Lipgloss token palette: Search query prompt (`#eab308`, bright yellow), Highlighted matches (`#ec4899`, bold magenta), Match counter badge (`#10b981`, green pill).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX TOPOLOGY TREE NAVIGATOR] ────────────────────────────────────────────────────────┐
      │ [-] domain: enterprise-outreach.com                                                        │
      │   ├─ [MATCH] mbx: sales-04@enterprise-outreach.com [IP: 192.0.2.45]                       │
      │   └─ [MATCH] mbx: sales-09@enterprise-outreach.com [IP: 192.0.2.50]                       │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ /sales- [Matches: 2/48 nodes] (Press Enter to cycle, Esc to cancel)                         │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline DOM text search filter with native client-side string ranking and automated node unfolding.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup tree-find`
  - Flags:
    - `--query <string>`: Initial search query string.
    - `--auto-expand`: Automatically expand all matched nodes on launch.
    - `--select-first`: Automatically highlight and focus first match.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FuzzySearchState struct {
        QueryPattern   string   `json:"queryPattern"`
        MatchedNodeIDs []string `json:"matchedNodeIds"`
        ActiveIndex    int      `json:"activeIndex"`
        MatchCount     int      `json:"matchCount"`
        IsSearching    bool     `json:"isSearching"`
        HasMatches     bool     `json:"hasMatches"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxSearchHistory (
        searchId TEXT PRIMARY KEY,
        queryPattern TEXT NOT NULL,
        matchCount INTEGER NOT NULL DEFAULT 0,
        hasMatches INTEGER NOT NULL DEFAULT 1,
        executedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/mailbox_search.go` integrated with `pkg/search/fuzzy`.

---

### Improvement 636: Hierarchical Mailbox Grouping by Warmup Stage and Health Score
- **ID:** `MAILBOX-GRAPH-636`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Enable dynamic regrouping of the visual mailbox tree by Warmup Stage (Seed, Ramp, Scale, Mature) or Health Score tier (Optimal, Warning, Degraded) with a single keystroke.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Tab bar toggling between `[1] By Domain`, `[2] By Warmup Stage`, `[3] By Health Tier`.
    - Stage icons: Seed (`🌱`), Ramp (`🚀`), Scale (`📈`), Mature (`🏆`).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX HIERARCHY] ── [1: DOMAIN] ─ [2: BY STAGE] ─ [3: BY HEALTH] ──────────────────────┐
      │ ▼ [STAGE: RAMP (100 - 500 msgs/day)] (14 Mailboxes)                                        │
      │   ├─ sender-02@domain-a.com [Day 8/30] [Quota: 220/250] [Health: 99.1%]                    │
      │   ├─ sender-05@domain-b.com [Day 11/30] [Quota: 310/350] [Health: 98.4%]                   │
      │ ▼ [STAGE: SEED (10 - 50 msgs/day)] (6 Mailboxes)                                           │
      │   └─ sender-01@fresh-domain.org [Day 2/30] [Quota: 25/30] [Health: 100.0%]                 │
      │ [Tab: Switch View Grouping] [Space: Toggle Group] [q: Quit]                                │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Accordion-style responsive card groups with pure CSS transitions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mailbox-groups`
  - Flags:
    - `--group-by <string>`: Grouping strategy (`domain`, `stage`, `health`) (default: `stage`).
    - `--sort <string>`: Sort direction within group (`asc`, `desc`, `health`) (default: `health`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxGroupContainer struct {
        GroupName      string               `json:"groupName"`
        GroupType      string               `json:"groupType"`
        TotalCount     int                  `json:"totalCount"`
        AverageHealth  float64              `json:"averageHealth"`
        MailboxNodes   []MailboxTopologyNode `json:"mailboxNodes"`
        IsExpanded     bool                 `json:"isExpanded"`
        HasWarningNode bool                 `json:"hasWarningNode"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxGroupingPreferences (
        userId TEXT PRIMARY KEY,
        defaultGroupBy TEXT NOT NULL DEFAULT 'stage',
        sortBy TEXT NOT NULL DEFAULT 'health',
        isStageExpanded INTEGER NOT NULL DEFAULT 1,
        hasFilterApplied INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/mailbox_grouping.go`.

---

### Improvement 637: Real-Time Node Telemetry Sidecar Panel with Mini Sparkline HUD
- **ID:** `MAILBOX-GRAPH-637`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Render a split-screen sidecar details panel adjacent to the tree navigator displaying real-time deliverability metrics, error ratios, and mini throughput sparklines for the active node.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - 60/40 horizontal split pane: Left pane holds tree hierarchy; Right pane renders node telemetry sidecar.
    - Lipgloss token palette: Sidecar border (`#3b82f6`, cobalt blue), Deliverability gauge (`#10b981`, green), Latency sparkline (`#06b6d4`, cyan ` ▂▃▅`).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX TREE] ────────────┬─ [NODE DETAILS: sender-02@deliverability-core.io] ───────────┐
      │ domain: deliverability-core │ Outbound IP: 198.51.100.13   [RDNS: mail02.deliverability.io]│
      │   ├─ sender-01              │ Current Quota: 490/500       [Rate: 1.2 msgs/sec]            │
      │ ► ├─ sender-02 (SELECTED)   │ Inbox Placement: 98.4%       [SPF: PASS | DKIM: PASS]        │
      │   └─ sender-03              │ Throughput Sparkline (5m):    ▂▃▅▇█▇▆▅▃▂ [Avg RTT: 48ms]     │
      │                             │ Active Sockets: 4 open       [Backpressure: NONE]            │
      └─────────────────────────────┴──────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - CSS Flexbox 2-column layout with instant WebSocket node detail streaming.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup sidecar-view`
  - Flags:
    - `--node-id <string>`: Initial focused node.
    - `--split-ratio <string>`: Ratio between tree and sidecar (default: `60-40`).
    - `--watch`: Persistent interactive mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxNodeSidecarTelemetry struct {
        NodeID             string    `json:"nodeId"`
        OutboundIP         string    `json:"outboundIp"`
        ReverseDNS         string    `json:"reverseDns"`
        CurrentRatePerSec  float64   `json:"currentRatePerSec"`
        InboxPlacementRate float64   `json:"inboxPlacementRate"`
        AverageRTTMs       float64   `json:"averageRttMs"`
        SparklineHistory   []int     `json:"sparklineHistory"`
        IsSPFPassing       bool      `json:"isSpfPassing"`
        IsDKIMPassing      bool      `json:"isDkimPassing"`
        HasBackpressure    bool      `json:"hasBackpressure"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxSidecarSnapshots (
        snapshotId TEXT PRIMARY KEY,
        nodeId TEXT NOT NULL,
        ratePerSec REAL NOT NULL DEFAULT 0.0,
        inboxPlacement REAL NOT NULL DEFAULT 100.0,
        rttMs REAL NOT NULL DEFAULT 0.0,
        isSpfPassing INTEGER NOT NULL DEFAULT 1,
        isDkimPassing INTEGER NOT NULL DEFAULT 1,
        hasBackpressure INTEGER NOT NULL DEFAULT 0,
        capturedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_sidecar_node ON MailboxSidecarSnapshots(nodeId, capturedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Added to `cli/cmdemail/mailbox_sidecar.go` interfacing with `pkg/telemetry/sidecar`.

---

### Improvement 638: Keyboard-Driven Multi-Select Batch Quarantine and Warming State Modifiers
- **ID:** `MAILBOX-GRAPH-638`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Permit bulk selection of multiple mailbox tree nodes using `Space` or `Shift+Down` to batch-modify warming stages or trigger instant quarantine isolation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Checkbox selection markers (`[✓]` selected, `[ ]` unselected) next to tree leaves with batch counter indicator on footer.
    - Lipgloss token palette: Selection counter (`#f59e0b`, bright amber), Selected items (`#38bdf8`, cyan background).
    - Model view snippet:
      ```text
      ┌─ [BATCH MAILBOX SELECTION] ────────────────────────────────────────────────────────────────┐
      │ [✓] sender-01@domain-a.com [IP: 198.51.100.12]                                             │
      │ [✓] sender-02@domain-a.com [IP: 198.51.100.13]                                             │
      │ [ ] sender-03@domain-a.com [IP: 198.51.100.14]                                             │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [Selected: 2 mailboxes] [Actions: 'q': Quarantine, 's': Step Up Stage, 'p': Pause, 'c': Clear] │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - HTML multi-checkbox list with sticky bottom floating batch command bar.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup batch-modify`
  - Flags:
    - `--node-ids <string>`: Comma-separated list of target node IDs.
    - `--action <string>`: Batch action (`quarantine`, `promote-stage`, `pause`, `resume`).
    - `--force`: Apply without confirmation prompt.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BatchModifyRequest struct {
        TargetNodeIDs []string `json:"targetNodeIds"`
        ActionType    string   `json:"actionType"`
        BatchReason   string   `json:"batchReason"`
        IsDryRun      bool     `json:"isDryRun"`
        HasConfirmed  bool     `json:"hasConfirmed"`
    }

    type BatchModifyResponse struct {
        SuccessCount  int      `json:"successCount"`
        FailureCount  int      `json:"failureCount"`
        ModifiedNodes []string `json:"modifiedNodes"`
        IsCompleted   bool     `json:"isCompleted"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxBatchOperations (
        batchOpId TEXT PRIMARY KEY,
        actionType TEXT NOT NULL,
        affectedCount INTEGER NOT NULL DEFAULT 0,
        isCompleted INTEGER NOT NULL DEFAULT 1,
        executedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured in `cli/cmdemail/mailbox_batch.go`.

---

### Improvement 639: Compact Topological ASCII Minimap for Fleet-Scale Viewports
- **ID:** `MAILBOX-GRAPH-639`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Render a high-density Braille or ASCII minimap (`⠿⠾⠽⠻`) at the top right of the terminal displaying cluster density and health hotspots when navigating mailboxes numbering in the thousands.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - 16x8 character block minimap displaying fleet distribution; current viewport represented by inverted bracket borders (`[ ]`).
    - Colors: Healthy pixels (green dots `⣿`), Degraded pixels (amber dots `⣦`), Quarantined (red dots `⣶`).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX TOPOLOGY TREE] ──────────────────────────┬─ [FLEET MINIMAP: 1,200 MBX] ──┐
      │ [-] domain-01.com [12 mbx]                         │ ⣿⣿⣿⣿⣿⣿⣾⣽⣻⢿ [VIEWPORT] │
      │   ├─ mbx-001 [OK]                                  │ ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿ │
      │   ├─ mbx-002 [OK]                                  │ ⣿⣿⣿⣿⣦⣦⣿⣿⣿⣿⣿⣿ │
      │   └─ mbx-003 [OK]                                  │ ⣿⣿⣿⣿⣿⣿⣿⣿⣶⣶⣿⣿ │
      └────────────────────────────────────────────────────┴───────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - HTML5 `<canvas>` element (120x60 px) rendered via native 2D canvas API.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup minimap`
  - Flags:
    - `--width <int>`: Minimap character width (default: `24`).
    - `--height <int>`: Minimap character height (default: `8`).
    - `--enable-minimap`: Toggle minimap visibility in tree mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FleetMinimapGrid struct {
        WidthCols     int      `json:"widthCols"`
        HeightRows    int      `json:"heightRows"`
        BrailleCells  []string `json:"brailleCells"`
        ViewportStart int      `json:"viewportStart"`
        ViewportEnd   int      `json:"viewportEnd"`
        IsVisible     bool     `json:"isVisible"`
        HasAlerts     bool     `json:"hasAlerts"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FleetMinimapCache (
        cacheKey TEXT PRIMARY KEY,
        compressedGrid TEXT NOT NULL,
        totalNodes INTEGER NOT NULL DEFAULT 0,
        isFresh INTEGER NOT NULL DEFAULT 1,
        generatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/mailbox_minimap.go`.

---

### Improvement 640: Live Event Log Tail Strip with Interactive Log Expansion
- **ID:** `MAILBOX-GRAPH-640`
- **Cluster:** 64 (High-Performance Bubble Tea Visual Mailbox Graph & Interactive Tree Navigator)
- **Goal:** Embed an interactive 3-line rolling live event tail strip at the bottom of the Bubble Tea mailbox graph, expandable into full-screen log stream via `Ctrl+L`.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Bottom 3-row footer displaying colorized real-time dispatch, bounce, and handshake logs.
    - Lipgloss token palette: Info logs (`#64748b`, slate), Warn (`#eab308`, amber), Err/Bounce (`#ef4444`, red).
    - Model view snippet:
      ```text
      ┌─ [LIVE DISPATCH & EVENT LOG TAIL] ─────────────────────────────────────────────────────────┐
      │ 21:42:01 [INFO]  [mbx-02] TLS 1.3 handshake ok with aspmx.l.google.com (142.250.102.27)     │
      │ 21:42:03 [WARN]  [mbx-05] Throttled by mx1.mail.icloud.com: 451 4.7.1 Try again later      │
      │ 21:42:04 [INFO]  [mbx-01] Delivered warm-up seed item #4928 to test-seed-99@yahoo.com     │
      │ [Ctrl+L: Fullscreen Log Inspector] [Ctrl+C: Quit] [Space: Pause Stream]                    │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Embedded terminal log console with virtualized scrolling (`pre` formatted text buffer).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup log-tail`
  - Flags:
    - `--lines <int>`: Number of tail lines to display (default: `3`).
    - `--filter-level <string>`: Log level filter (`debug`, `info`, `warn`, `error`) (default: `info`).
    - `--fullscreen`: Launch directly into full-screen log viewer.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxEventLogEntry struct {
        LogID       string    `json:"logId"`
        Timestamp   time.Time `json:"timestamp"`
        Level       string    `json:"level"`
        MailboxID   string    `json:"mailboxId"`
        RemoteHost  string    `json:"remoteHost"`
        Message     string    `json:"message"`
        IsError     bool      `json:"isError"`
        HasWarning  bool      `json:"hasWarning"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxEventLogs (
        logId TEXT PRIMARY KEY,
        mailboxId TEXT NOT NULL,
        remoteHost TEXT NOT NULL DEFAULT '',
        level TEXT NOT NULL,
        message TEXT NOT NULL,
        isError INTEGER NOT NULL DEFAULT 0,
        hasWarning INTEGER NOT NULL DEFAULT 0,
        loggedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_mbx_logs_time ON MailboxEventLogs(loggedAt DESC);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/cmdemail/mailbox_logs.go` interfacing with `pkg/logger/streambuffer`.

---

# Cluster 65: ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression (Items 641–650)

---

### Improvement 641: Autonomous ARF (Abuse Reporting Format) RFC 5965 Ingest Engine
- **ID:** `FBL-SPAM-641`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Ingest, decrypt, and parse RFC 5965 Abuse Reporting Format (ARF) MIME messages from major ISP feedback loops (Yahoo CFL, Microsoft SNDS/JMRP, Comcast, Fastmail) and extract originating envelope data.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live ARF complaint ingest table showing ISP source, reporting mailbox, feedback type (`abuse`, `fraud`, `not-spam`), original Message-ID, and auto-suppression action.
    - Lipgloss token palette: Abuse complaint (`#ef4444`, bright red badge `[ABUSE]`), ISP tag (`#8b5cf6`, violet `[YAHOO-FBL]`), Suppressed (`#10b981`, green check).
    - Model view snippet:
      ```text
      ┌─ [AUTONOMOUS ISP FEEDBACK LOOP (FBL) INGEST ENGINE] ───────────────────────────────────────┐
      │ TIME     │ ISP SOURCE    │ TYPE    │ RECIPIENT COMPLAINANT      │ ACTION     │ HEALTH IMPACT │
      │ 21:44:12 │ Yahoo CFL     │ [ABUSE] │ user-3921@yahoo.com        │ SUPPRESSED │ -0.05% (98.9) │
      │ 21:44:08 │ MSFT JMRP     │ [ABUSE] │ executive@hotmail.com      │ SUPPRESSED │ -0.04% (98.9) │
      │ 21:43:55 │ Fastmail FBL  │ [ABUSE] │ contact@fastmail.com       │ SUPPRESSED │ -0.02% (99.0) │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [Active Ingest: 3 feeds] [Total Complaints 24h: 3] [Current Global Abuse Ratio: 0.008%]    │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Real-time tabular event ledger with expandable raw RFC 5965 MIME headers and collapsible MIME body parts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fbl-ingest`
  - Flags:
    - `--listen-addr <string>`: Inbound SMTP or IMAP FBL mailbox address.
    - `--auto-suppress`: Automatically add complaining addresses to global suppression database (default: `true`).
    - `--dump-raw <path>`: Dump raw incoming RFC 5965 multipart payloads to disk for forensic review.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FBLComplaintRecord struct {
        ComplaintID       string    `json:"complaintId"`
        ISPSource         string    `json:"ispSource"`
        FeedbackType      string    `json:"feedbackType"`
        ReportedRecipient string    `json:"reportedRecipient"`
        OriginalMessageID string    `json:"originalMessageId"`
        OriginatingIP     string    `json:"originatingIp"`
        SenderMailbox     string    `json:"senderMailbox"`
        ReceivedAt        time.Time `json:"receivedAt"`
        IsSuppressed      bool      `json:"isSuppressed"`
        HasDKIMVerified   bool      `json:"hasDkimVerified"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FBLComplaints (
        complaintId TEXT PRIMARY KEY,
        ispSource TEXT NOT NULL,
        feedbackType TEXT NOT NULL,
        reportedRecipient TEXT NOT NULL,
        originalMessageId TEXT NOT NULL,
        originatingIp TEXT NOT NULL,
        senderMailbox TEXT NOT NULL,
        isSuppressed INTEGER NOT NULL DEFAULT 1,
        hasDkimVerified INTEGER NOT NULL DEFAULT 1,
        receivedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_fbl_recipient ON FBLComplaints(reportedRecipient);
    CREATE INDEX IF NOT EXISTS idx_fbl_sender ON FBLComplaints(senderMailbox);
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded in `cli/cmdemail/fbl_ingest.go` connecting to `pkg/fbl/arf_parser.go`.

---

### Improvement 642: Instant Multi-Tier Spam Complaint Auto-Suppression Engine
- **ID:** `FBL-SPAM-642`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Instantly blacklist complaining recipients across three tiers: (1) In-memory cache (<1ms), (2) SQLite Local DB (<5ms), (3) Distributed Redis/Cluster sync, canceling any inflight warm-up threads.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time suppression audit feed displaying tier convergence times and cancellation of queued mail items.
    - Lipgloss token palette: Tier 1 memory (`#06b6d4`, cyan `[L1: 0.2ms]`), Tier 2 SQLite (`#10b981`, green `[L2: 1.8ms]`), Tier 3 Fleet (`#8b5cf6`, purple `[L3: 12ms]`).
    - Model view snippet:
      ```text
      ┌─ [INSTANT COMPLAINT AUTO-SUPPRESSION ENGINE] ──────────────────────────────────────────────┐
      │ TARGET EMAIL           │ SUPPRESSION TIERS                 │ IN-FLIGHT QUEUE CANCELLED       │
      │ user-99@yahoo.com      │ [L1-MEM: OK] [L2-DB: OK] [L3: OK] │ 2 pending warm-up seeds ABORTED │
      │ alert-4@hotmail.com    │ [L1-MEM: OK] [L2-DB: OK] [L3: OK] │ 0 pending items (Clean)         │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Active Suppressed Cache: 4,120 entries | Average Cascade Latency: 4.8ms                    │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live status pills with micro-badges indicating suppression status and timestamp.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup suppress`
  - Flags:
    - `--email <string>`: Email address to manually add or verify in suppression tiers.
    - `--reason <string>`: Suppression reason (`fbl_complaint`, `hard_bounce`, `unsubscribe`, `manual`).
    - `--cascade-fleet`: Broadcast suppression event to all connected fleet nodes immediately.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SuppressionEntry struct {
        EmailAddress   string    `json:"emailAddress"`
        SuppressionTag string    `json:"suppressionTag"`
        SourceReason   string    `json:"sourceReason"`
        CascadeLatency float64   `json:"cascadeLatencyMs"`
        SuppressedAt   time.Time `json:"suppressedAt"`
        IsGlobal       bool      `json:"isGlobal"`
        HasPropagated  bool      `json:"hasPropagated"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SuppressionBlacklist (
        emailAddress TEXT PRIMARY KEY,
        suppressionTag TEXT NOT NULL DEFAULT 'fbl_complaint',
        sourceReason TEXT NOT NULL,
        isGlobal INTEGER NOT NULL DEFAULT 1,
        hasPropagated INTEGER NOT NULL DEFAULT 1,
        suppressedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_suppress_time ON SuppressionBlacklist(suppressedAt DESC);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/suppression_engine.go` interfacing with `pkg/suppression/tier_cache.go`.

---

### Improvement 643: ISP Complaint-Ratio Velocity Gauge and Safety Ceiling Breaker
- **ID:** `FBL-SPAM-643`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Continuously calculate rolling 24-hour complaint percentages per destination ISP (e.g., Yahoo/Gmail 0.10% threshold) and trip an automated circuit breaker if velocity exceeds 0.08%.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Radial or linear ANSI gauge comparing current complaint rates against the strict 0.10% industry danger line.
    - Color tokens: Normal (<0.03% `#10b981`, green), Elevated (0.05-0.08% `#f59e0b`, amber), Tripped (>0.08% `#ef4444`, flashing red `[CIRCUIT OPEN]`).
    - Model view snippet:
      ```text
      ┌─ [ISP COMPLAINT RATIO VELOCITY & SAFETY BREAKER] ──────────────────────────────────────────┐
      │ DESTINATION ISP │ 24H DELIVERED │ 24H COMPLAINTS │ RATIO   │ SAFETY STATUS │ BREAKER ACTION   │
      │ Google (Gmail)  │ 14,820 msgs   │ 2 complaints   │ 0.013%  │ [SAFE]        │ NORMAL DISPATCH  │
      │ Yahoo Mail      │ 8,400 msgs    │ 6 complaints   │ 0.071%  │ [ELEVATED]    │ THROTTLE -50%    │
      │ Microsoft 365   │ 6,100 msgs    │ 1 complaints   │ 0.016%  │ [SAFE]        │ NORMAL DISPATCH  │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Maximum Allowed Ratio: 0.100% | Protective Auto-Trip Level: 0.080%                         │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG progress meters with dynamic threshold markers and danger gradient fills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup complaint-breaker`
  - Flags:
    - `--isp <string>`: Target ISP domain (`yahoo`, `gmail`, `microsoft`, `all`) (default: `all`).
    - `--threshold-trip <float>`: Complaint ratio that trips circuit breaker (default: `0.08`).
    - `--reset`: Reset tripped circuit breaker manually after remediation.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ISPComplaintVelocity struct {
        ISPDomain        string    `json:"ispDomain"`
        DeliveredCount   int       `json:"deliveredCount"`
        ComplaintCount   int       `json:"complaintCount"`
        CurrentRatio     float64   `json:"currentRatio"`
        ThresholdRatio   float64   `json:"thresholdRatio"`
        IsTripped        bool      `json:"isTripped"`
        HasWarningStatus bool      `json:"hasWarningStatus"`
        LastEvaluatedAt  time.Time `json:"lastEvaluatedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ISPComplaintGauges (
        ispDomain TEXT PRIMARY KEY,
        deliveredCount INTEGER NOT NULL DEFAULT 0,
        complaintCount INTEGER NOT NULL DEFAULT 0,
        currentRatio REAL NOT NULL DEFAULT 0.0,
        isTripped INTEGER NOT NULL DEFAULT 0,
        hasWarningStatus INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/complaint_breaker.go` connected to `pkg/warmup/circuitbreaker.go`.

---

### Improvement 644: Google Postmaster Tools (GPT) API Telemetry Importer and Sync HUD
- **ID:** `FBL-SPAM-644`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Ingest Google Postmaster Tools API data (User-Reported Spam Rate, Domain Reputation, IP Reputation, DKIM/SPF Success Rates) and present synced timeline cards in Bubble Tea.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Tabular summary of domain reputation levels (`HIGH`, `MEDIUM`, `LOW`, `BAD`) with 7-day trend arrows (`↑`, `→`, `↓`).
    - Color tokens: High reputation (`#10b981`, bold green), Medium (`#3b82f6`, blue), Low (`#eab308`, amber), Bad (`#ef4444`, red).
    - Model view snippet:
      ```text
      ┌─ [GOOGLE POSTMASTER TOOLS (GPT) TELEMETRY HUD] ────────────────────────────────────────────┐
      │ DOMAIN                 │ DOMAIN REP │ IP REP  │ SPAM RATE │ SPF PASS │ DKIM PASS │ 7D TREND  │
      │ deliverability-core.io │ HIGH       │ HIGH    │ 0.01%     │ 100.0%   │ 100.0%    │ [STABLE ↑]│
      │ mail-outreach-v2.org   │ MEDIUM     │ HIGH    │ 0.04%     │ 99.8%    │ 100.0%    │ [SLIGHT ↓]│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Last API Sync: 12 minutes ago | Next Scheduled Polling: in 48 minutes                      │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Metric overview cards with reputation badge pills and SVG trend sparklines.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup gpt-sync`
  - Flags:
    - `--domain <string>`: Google Postmaster verified domain.
    - `--days <int>`: History sync window in days (default: `7`).
    - `--credentials <path>`: Path to Google Service Account JSON key.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PostmasterDomainRecord struct {
        DomainName       string    `json:"domainName"`
        DomainReputation string    `json:"domainReputation"`
        IPReputation     string    `json:"ipReputation"`
        UserSpamRate     float64   `json:"userSpamRate"`
        SPFPassRate      float64   `json:"spfPassRate"`
        DKIMPassRate     float64   `json:"dkimPassRate"`
        IsHighReputation bool      `json:"isHighReputation"`
        HasReputationDrop bool     `json:"hasReputationDrop"`
        SyncedAt         time.Time `json:"syncedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GooglePostmasterSnapshots (
        snapshotId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        domainReputation TEXT NOT NULL,
        ipReputation TEXT NOT NULL,
        userSpamRate REAL NOT NULL DEFAULT 0.0,
        spfPassRate REAL NOT NULL DEFAULT 100.0,
        dkimPassRate REAL NOT NULL DEFAULT 100.0,
        isHighReputation INTEGER NOT NULL DEFAULT 1,
        hasReputationDrop INTEGER NOT NULL DEFAULT 0,
        syncedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_gpt_domain ON GooglePostmasterSnapshots(domainName, syncedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/gpt_sync.go` interfacing with `pkg/postmaster/client.go`.

---

### Improvement 645: Microsoft SNDS / JMRP Automated Log Parser and IP Reputation Matrix
- **ID:** `FBL-SPAM-645`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Automate ingestion of Microsoft Smart Network Data Services (SNDS) automated data feeds, evaluating IP filter results (Green/Yellow/Red) and complaint ratios for Hotmail/Outlook delivery.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density matrix mapping outbound warmup IP blocks to Microsoft SNDS color status indicators.
    - Lipgloss tokens: Normal Green (`#10b981`), Warning Yellow (`#eab308`), Blocklist Red (`#ef4444`).
    - Model view snippet:
      ```text
      ┌─ [MICROSOFT SNDS / JMRP IP REPUTATION MATRIX] ─────────────────────────────────────────────┐
      │ IP ADDRESS     │ REVERSE DNS               │ SNDS STATUS │ COMPLAINT RATE │ TRAP HITS │ RCPT │
      │ 198.51.100.12  │ mail01.deliverability.io  │ [GREEN]     │ < 0.1%         │ 0 hits    │ OK   │
      │ 198.51.100.13  │ mail02.deliverability.io  │ [GREEN]     │ < 0.1%         │ 0 hits    │ OK   │
      │ 198.51.100.14  │ mail03.deliverability.io  │ [YELLOW]    │ 0.3% (ELEVATED)│ 0 hits    │ HOLD │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [Key: Green=Good, Yellow=Elevated, Red=Blocked] [Action: 'r' Re-check, 'h' Put on Hold]    │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Color-coded grid showing SNDS diagnostic output with one-click quarantine buttons.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup snds-inspect`
  - Flags:
    - `--ip-range <cidr>`: Inspect specific CIDR subnet (e.g. `198.51.100.0/24`).
    - `--api-key <string>`: Microsoft SNDS automated data export key.
    - `--quarantine-yellow`: Auto-quarantine mailboxes on IPs reporting yellow status.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SNDSIPStatusRecord struct {
        IPAddress        string    `json:"ipAddress"`
        ReverseDNS       string    `json:"reverseDns"`
        FilterResult     string    `json:"filterResult"`
        ComplaintRatio   float64   `json:"complaintRatio"`
        SpamTrapHits     int       `json:"spamTrapHits"`
        IsGreenStatus    bool      `json:"isGreenStatus"`
        HasTrapHit       bool      `json:"hasTrapHit"`
        RecordedAt       time.Time `json:"recordedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MicrosoftSNDSSnapshots (
        recordId TEXT PRIMARY KEY,
        ipAddress TEXT NOT NULL,
        reverseDns TEXT NOT NULL,
        filterResult TEXT NOT NULL,
        complaintRatio REAL NOT NULL DEFAULT 0.0,
        spamTrapHits INTEGER NOT NULL DEFAULT 0,
        isGreenStatus INTEGER NOT NULL DEFAULT 1,
        hasTrapHit INTEGER NOT NULL DEFAULT 0,
        recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_snds_ip ON MicrosoftSNDSSnapshots(ipAddress);
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured in `cli/cmdemail/snds_inspect.go` pulling from `pkg/postmaster/snds.go`.

---

### Improvement 646: List-Unsubscribe Header Compliance Verifier & One-Click Mailto HUD
- **ID:** `FBL-SPAM-646`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Audit RFC 8058 One-Click `List-Unsubscribe` and `List-Unsubscribe-Post` headers across all outgoing warm-up seed templates to ensure zero RFC compliance violations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Header inspection pane highlighting RFC 8058 conformance, HTTPS one-click POST targets, and mailto fallback endpoints.
    - Lipgloss tokens: Compliant (`#10b981`, bold green `[RFC 8058 PASS]`), Missing header (`#ef4444`, red `[NON-COMPLIANT]`).
    - Model view snippet:
      ```text
      ┌─ [RFC 8058 LIST-UNSUBSCRIBE AUDIT] ────────────────────────────────────────────────────────┐
      │ TEMPLATE / POOL        │ ONE-CLICK HTTPS ENDPOINT            │ MAILTO FALLBACK     │ STATUS  │
      │ Seed-Warmup-Template-A │ https://unsub.deliverability.io/u/1 │ mailto:unsub@deliv  │ [PASS]  │
      │ Seed-Warmup-Template-B │ https://unsub.deliverability.io/u/2 │ mailto:unsub@deliv  │ [PASS]  │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ POST Body Format: 'List-Unsubscribe=One-Click' | SSL Certificate: VALID (Exp: 180d)        │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Interactive header inspector with click-to-test mock unsubscribe POST simulation.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup verify-unsub`
  - Flags:
    - `--template-id <string>`: Check headers on specific warmup template.
    - `--test-post`: Send test RFC 8058 HTTP POST request to endpoint to verify 200 OK.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeComplianceResult struct {
        TemplateID       string `json:"templateId"`
        HasListUnsub     bool   `json:"hasListUnsub"`
        HasOneClickPost  bool   `json:"hasOneClickPost"`
        HasMailtoFallback bool   `json:"hasMailtoFallback"`
        HTTPSEndpoint    string `json:"httpsEndpoint"`
        IsRFC8058Valid   bool   `json:"isRfc8058Valid"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubscribeAuditLogs (
        auditId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        httpsEndpoint TEXT NOT NULL,
        isRfc8058Valid INTEGER NOT NULL DEFAULT 1,
        hasOneClickPost INTEGER NOT NULL DEFAULT 1,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/unsub_verify.go`.

---

### Improvement 647: Spam Trap Detection and Seed Pool Cross-Check Sanitizer
- **ID:** `FBL-SPAM-647`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Run heuristic and honey-pot fingerprinting algorithms across active peer seed mailboxes to detect and quarantine suspected pristine or recycled spam traps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Threat analysis log flagging abnormal dormant recipient addresses that suddenly generate hard bounce loops or silent blackhole FBL flags.
    - Lipgloss token palette: Warning flag (`#eab308`, amber `[SUSPECT TRAP]`), Quarantined (`#ef4444`, red `[PURGED]`).
    - Model view snippet:
      ```text
      ┌─ [SPAM TRAP DETECTION & SEED SANITIZER] ───────────────────────────────────────────────────┐
      │ ADDRESS SCANNED         │ HEURISTIC ANOMALY                 │ RISK SCORE │ REMEDIATION       │
      │ inactive-seed-40@msn.com│ Zero interaction > 90d + FBL ping │ 94 / 100   │ QUARANTINED       │
      │ deadbox-12@yahoo.com    │ Repeated MX 550 then silent accept│ 88 / 100   │ PURGED FROM POOL  │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Clean Seed Pool: 4,892 verified addresses | Quarantined Today: 2 addresses                 │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Visual risk assessment scatter plot with isolated seed list management.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup sanitize-seeds`
  - Flags:
    - `--pool <string>`: Name of recipient seed pool (default: `global-seeds`).
    - `--risk-threshold <int>`: Minimum score to trigger automatic purge (default: `85`).
    - `--dry-run`: Display flagged addresses without modifying pool database.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedSanitizeAudit struct {
        SeedAddress    string    `json:"seedAddress"`
        AnomalyType    string    `json:"anomalyType"`
        RiskScore      int       `json:"riskScore"`
        IsQuarantined  bool      `json:"isQuarantined"`
        HasPurged      bool      `json:"hasPurged"`
        ScannedAt      time.Time `json:"scannedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedSanitizerAudits (
        auditId TEXT PRIMARY KEY,
        seedAddress TEXT NOT NULL,
        anomalyType TEXT NOT NULL,
        riskScore INTEGER NOT NULL,
        isQuarantined INTEGER NOT NULL DEFAULT 1,
        hasPurged INTEGER NOT NULL DEFAULT 0,
        scannedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_seed_sanitize ON SeedSanitizerAudits(seedAddress);
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured in `cli/cmdemail/seed_sanitizer.go`.

---

### Improvement 648: FBL Feedback Loop Incident Post-Mortem Terminal Report Generator
- **ID:** `FBL-SPAM-648`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Autonomously generate an end-to-end incident root cause analysis (RCA) report in Markdown and terminal ANSI whenever an ISP spam complaint spike occurs.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Structured multi-section RCA terminal viewer featuring timeline chronology, contributing factors, affected mailbox nodes, and automated remediations.
    - Lipgloss token palette: Incident header (`#ef4444`, red background), Section headers (`#3b82f6`, bold blue), Timeline bullets (`#06b6d4`, cyan).
    - Model view snippet:
      ```text
      ┌─ [INCIDENT POST-MORTEM: FBL-INC-2026-10-27-01] ────────────────────────────────────────────┐
      │ ROOT CAUSE: Yahoo Mail complaint cluster traced to template subject line permutation #4    │
      │ ├─ Incident Trigger: 4 complaints received in 18 minutes (Yahoo CFL)                       │
      │ ├─ Affected Senders: sender-02@domain-a.com, sender-04@domain-a.com                        │
      │ ├─ Automated Action: Circuit breaker tripped; outbound Yahoo throttling reduced by 75%     │
      │ └─ Current State: Restored to Nominal; Affected subjects suppressed                        │
      │ [p: Save Markdown Report to 02-spec/...] [e: Export JSON] [q: Close Viewer]                │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Clean rendered Markdown report with collapsible timeline event nodes and printable layout.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fbl-rca`
  - Flags:
    - `--incident-id <string>`: Specific incident identifier.
    - `--output-file <path>`: Destination path for generated Markdown report.
    - `--print-terminal`: Print formatted ANSI post-mortem directly to stdout.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FBLIncidentPostMortem struct {
        IncidentID       string    `json:"incidentId"`
        TriggeredAt      time.Time `json:"triggeredAt"`
        ResolvedAt       time.Time `json:"resolvedAt"`
        RootCauseSummary string    `json:"rootCauseSummary"`
        ComplaintCount   int       `json:"complaintCount"`
        AffectedMailboxes []string `json:"affectedMailboxes"`
        RemediationSteps []string  `json:"remediationSteps"`
        IsResolved       bool      `json:"isResolved"`
        HasCircuitTripped bool     `json:"hasCircuitTripped"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FBLPostMortems (
        incidentId TEXT PRIMARY KEY,
        rootCauseSummary TEXT NOT NULL,
        complaintCount INTEGER NOT NULL,
        isResolved INTEGER NOT NULL DEFAULT 1,
        hasCircuitTripped INTEGER NOT NULL DEFAULT 1,
        triggeredAt DATETIME NOT NULL,
        resolvedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated in `cli/cmdemail/fbl_rca.go` leveraging `pkg/fbl/rca_engine.go`.

---

### Improvement 649: Automated DKIM Alignment and DMARC Feedback (RUF/RUA) Digest TUI
- **ID:** `FBL-SPAM-649`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Ingest and parse incoming DMARC aggregate XML (RUA) and forensic failure (RUF) reports, visualizing SPF/DKIM alignment ratios and spoofing attempts in Bubble Tea.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual split table: Top section displays DMARC alignment pass percentages by ISP; bottom section highlights unauthorized IP spoofing origins.
    - Lipgloss tokens: 100% Alignment (`#10b981`, green `[ALIGNED]`), Misaligned (`#ef4444`, red `[FAIL]`).
    - Model view snippet:
      ```text
      ┌─ [DMARC / DKIM ALIGNMENT & AGGREGATE DIGEST (RUA/RUF)] ────────────────────────────────────┐
      │ REPORTER (ISP) │ PASS COUNT │ FAIL COUNT │ SPF ALIGNED │ DKIM ALIGNED │ POLICY APPLIED       │
      │ google.com     │ 24,100     │ 0          │ 100.0%      │ 100.0%       │ none (PASS)          │
      │ yahoo.com      │ 12,450     │ 1          │ 99.9%       │ 100.0%       │ none (PASS)          │
      │ enterprise.ru  │ 0          │ 42         │ 0.0%        │ 0.0%         │ reject (SPOOF BLOCKED│
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - DMARC aggregate visualizer with geographical source IP map and alignment breakdown donuts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup dmarc-digest`
  - Flags:
    - `--rua-dir <path>`: Directory containing downloaded gzipped DMARC RUA XML files.
    - `--domain <string>`: Target domain to analyze.
    - `--show-spoofs`: Filter display to show only unauthorized source IPs failing alignment.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DMARCDigestRecord struct {
        ReportID       string    `json:"reportId"`
        DomainName     string    `json:"domainName"`
        ReporterOrg    string    `json:"reporterOrg"`
        TotalPass      int       `json:"totalPass"`
        TotalFail      int       `json:"totalFail"`
        SPFPassRate    float64   `json:"spfPassRate"`
        DKIMPassRate   float64   `json:"dkimPassRate"`
        IsFullAligned  bool      `json:"isFullAligned"`
        HasSpoofThreat bool      `json:"hasSpoofThreat"`
        ProcessedAt    time.Time `json:"processedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DMARCDigests (
        reportId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        reporterOrg TEXT NOT NULL,
        totalPass INTEGER NOT NULL DEFAULT 0,
        totalFail INTEGER NOT NULL DEFAULT 0,
        spfPassRate REAL NOT NULL DEFAULT 100.0,
        dkimPassRate REAL NOT NULL DEFAULT 100.0,
        isFullAligned INTEGER NOT NULL DEFAULT 1,
        hasSpoofThreat INTEGER NOT NULL DEFAULT 0,
        processedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_dmarc_domain ON DMARCDigests(domainName);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/dmarc_digest.go`.

---

### Improvement 650: Whitelist Verification and Feedback Loop Enrollment Health Check
- **ID:** `FBL-SPAM-650`
- **Cluster:** 65 (ISP Feedback Loop (FBL) Autonomous Parser & Spam Complaint Auto-Suppression)
- **Goal:** Autonomously verify enrollment status and reverse-DNS verification across major ISP feedback programs (Yahoo CFL, Microsoft SNDS, Spamhaus, CSA).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Status checklist verifying MX/IP registration, active postmaster email forwarders, and valid DKIM public keys across all participating domains.
    - Lipgloss tokens: Enrolled (`#10b981`, check `[ACTIVE]`), Pending Verification (`#eab308`, clock `[PENDING]`), Missing (`#ef4444`, cross `[UNENROLLED]`).
    - Model view snippet:
      ```text
      ┌─ [ISP FEEDBACK LOOP (FBL) ENROLLMENT HEALTH CHECK] ────────────────────────────────────────┐
      │ DOMAIN / IP        │ YAHOO CFL   │ MSFT SNDS/JMRP │ SPAMHAUS DBL │ POSTMASTER FORWARDER      │
      │ 198.51.100.12 (d1) │ [ACTIVE]    │ [ACTIVE]       │ [CLEAN]      │ postmaster@d1 (ACTIVE)    │
      │ 198.51.100.13 (d1) │ [ACTIVE]    │ [ACTIVE]       │ [CLEAN]      │ postmaster@d1 (ACTIVE)    │
      │ 198.51.100.14 (d2) │ [PENDING]   │ [ACTIVE]       │ [CLEAN]      │ postmaster@d2 (ACTIVE)    │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [All 3 Outbound IPs Registered with Abuse Desks] [Next Verification Check: 6 days]          │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Compliance checklist dashboard with direct link buttons to provider enrollment portals.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fbl-check`
  - Flags:
    - `--domain <string>`: Target domain to verify.
    - `--ip-list <string>`: Comma-separated outbound IP list.
    - `--fix-forwarders`: Auto-create local postmaster forwarder aliases if missing.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FBLEnrollmentStatus struct {
        Identifier         string    `json:"identifier"`
        YahooCFLStatus     string    `json:"yahooCflStatus"`
        MicrosoftSNDSStatus string   `json:"microsoftSndsStatus"`
        SpamhausStatus     string    `json:"spamhausStatus"`
        HasValidPostmaster bool      `json:"hasValidPostmaster"`
        IsFullyCovered     bool      `json:"isFullyCovered"`
        LastCheckedAt      time.Time `json:"lastCheckedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FBLEnrollments (
        identifier TEXT PRIMARY KEY,
        yahooCflStatus TEXT NOT NULL,
        microsoftSndsStatus TEXT NOT NULL,
        spamhausStatus TEXT NOT NULL,
        hasValidPostmaster INTEGER NOT NULL DEFAULT 1,
        isFullyCovered INTEGER NOT NULL DEFAULT 1,
        lastCheckedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/fbl_check.go`.

---

# Cluster 66: Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI (Items 651–660)

---

### Improvement 651: High-Density ANSI Latency Heatmap Matrix for Destination MX Relays
- **ID:** `SOCKET-JITTER-651`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Render a high-density 2D ANSI heatmap matrix in Bubble Tea showing round-trip latency (RTT in ms) and TCP connection handshake timings across destination MX relays.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual terminal grid dividing latency into 5-tier colored blocks: `<30ms` (`#10b981`, deep green `■`), `30-60ms` (`#06b6d4`, cyan `■`), `60-120ms` (`#eab308`, amber `■`), `120-250ms` (`#f97316`, orange `■`), `>250ms` (`#ef4444`, bright red `■`).
    - Model view snippet:
      ```text
      ┌─ [BOUNDED SOCKET LATENCY HEATMAP MATRIX] ──────────────────────────────────────────────────┐
      │ DESTINATION MX            │ IP ENDPOINT    │ TCP RTT │ TLS 1.3 │ DATA ACK │ 60S JITTER MATRIX│
      │ aspmx.l.google.com        │ 142.250.102.27 │ 24.2ms  │ 18.5ms  │ 12.1ms   │ ■■■■■■■■■■■■■■■■ │
      │ alt1.aspmx.l.google.com   │ 142.250.153.27 │ 28.6ms  │ 21.0ms  │ 14.2ms   │ ■■■■■■■■■■■■■■■■ │
      │ mta5.am0.yahoodns.net     │ 67.195.204.77  │ 84.1ms  │ 62.4ms  │ 41.8ms   │ ■■■■■■■■■■■■■■■■ │
      │ outlook-com.olc.prot.ms   │ 104.47.74.33   │ 112.5ms │ 89.2ms  │ 55.0ms   │ ■■■■■■■■■■■■■■■■ │
      │ mx1.mail.icloud.com       │ 17.57.155.1    │ 194.0ms │ 142.1ms │ 86.4ms   │ ■■■■■■■■■■■■■■■■ │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [Average Fleet RTT: 48.2ms] [P99 Latency: 198.4ms] [Legend: ■<30ms ■<60ms ■<120ms ■>250ms] │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Dynamic HTML5 canvas matrix rendering live 60-frame-per-second heatmap pixels with hovering tooltip overlays.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup socket-heatmap`
  - Flags:
    - `--mx-host <string>`: Filter by MX target hostname (default: `all`).
    - `--interval <string>`: Probe sampling interval (default: `2s`).
    - `--p99-alert <int>`: Trip visual alert if P99 RTT exceeds milliseconds threshold (default: `200`).
    - `--watch`: Persistent interactive Bubble Tea full-screen alternate buffer mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SocketLatencySample struct {
        MXHost       string    `json:"mxHost"`
        RemoteIP     string    `json:"remoteIp"`
        TCPHandshake float64   `json:"tcpHandshakeMs"`
        TLSHandshake float64   `json:"tlsHandshakeMs"`
        DataAckTime  float64   `json:"dataAckTimeMs"`
        TotalRTT     float64   `json:"totalRttMs"`
        IsUnderSLO   bool      `json:"isUnderSlo"`
        HasJitterSpike bool    `json:"hasJitterSpike"`
        SampledAt    time.Time `json:"sampledAt"`
    }

    type LatencyHeatmapState struct {
        ActiveSamples []SocketLatencySample `json:"activeSamples"`
        P99LatencyMs  float64              `json:"p99LatencyMs"`
        AverageRTTMs  float64              `json:"averageRttMs"`
        IsHealthy     bool                 `json:"isHealthy"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SocketLatencySamples (
        sampleId TEXT PRIMARY KEY,
        mxHost TEXT NOT NULL,
        remoteIp TEXT NOT NULL,
        tcpHandshakeMs REAL NOT NULL DEFAULT 0.0,
        tlsHandshakeMs REAL NOT NULL DEFAULT 0.0,
        dataAckTimeMs REAL NOT NULL DEFAULT 0.0,
        totalRttMs REAL NOT NULL DEFAULT 0.0,
        isUnderSlo INTEGER NOT NULL DEFAULT 1,
        hasJitterSpike INTEGER NOT NULL DEFAULT 0,
        sampledAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_latency_mx ON SocketLatencySamples(mxHost, sampledAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured in `cli/cmdemail/socket_heatmap.go` leveraging `pkg/netpoll/latency_probe.go`.

---

### Improvement 652: Dynamic Exponential Jitter Compensator & Pacing Regulator
- **ID:** `SOCKET-JITTER-652`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Autonomously dynamically inject microsecond pseudo-random Gaussian jitter and exponential backoff pacing between outgoing warm-up SMTP packets to mimic organic human keystroke rhythms and bypass ISP anti-bot rate heuristics.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live pacing regulator status displaying current Gaussian distribution curve (`σ` standard deviation, `μ` mean) and active artificial pause intervals.
    - Lipgloss token palette: Pacing active (`#3b82f6`, blue `[PACING ACTIVE]`), Natural dispersion (`#10b981`, emerald `[HUMAN LIKE]`).
    - Model view snippet:
      ```text
      ┌─ [DYNAMIC JITTER COMPENSATOR & PACING REGULATOR] ──────────────────────────────────────────┐
      │ DESTINATION MX            │ BASE DELAY │ JITTER SIGMA (σ) │ LAST APPLIED PAUSE │ STATUS      │
      │ aspmx.l.google.com        │ 850ms      │ ±120ms (Gauss)   │ 912ms              │ [NATURAL]   │
      │ mta5.am0.yahoodns.net     │ 1,200ms    │ ±350ms (Gauss)   │ 1,480ms            │ [THROTTLED] │
      │ outlook-com.olc.prot.ms   │ 600ms      │ ±95ms (Gauss)    │ 640ms              │ [NATURAL]   │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Active Pacing Algorithm: Gaussian Box-Muller Distribution | Organic Jitter: ACTIVE (100%)  │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Real-time SVG probability bell curve charting observed inter-packet delay intervals against standard normal distribution.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup jitter-tune`
  - Flags:
    - `--base-delay <duration>`: Base delay between consecutive email dispatches (default: `1s`).
    - `--sigma <duration>`: Standard deviation for Gaussian jitter (default: `200ms`).
    - `--adaptive-backoff`: Automatically increase base delay when receiving HTTP 421 or SMTP 451 temporary rate rejects.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type JitterPacingConfig struct {
        TargetMX          string  `json:"targetMx"`
        BaseDelayMs       int     `json:"baseDelayMs"`
        SigmaJitterMs     int     `json:"sigmaJitterMs"`
        MinDelayMs        int     `json:"minDelayMs"`
        MaxDelayMs        int     `json:"maxDelayMs"`
        IsAdaptive        bool    `json:"isAdaptive"`
        HasSpikeDetected  bool    `json:"hasSpikeDetected"`
    }

    type JitterExecutionEvent struct {
        EventID       string    `json:"eventId"`
        TargetMX      string    `json:"targetMx"`
        AppliedDelay  int       `json:"appliedDelayMs"`
        IsBackoffTrip bool      `json:"isBackoffTrip"`
        ExecutedAt    time.Time `json:"executedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS JitterPacingConfigs (
        targetMx TEXT PRIMARY KEY,
        baseDelayMs INTEGER NOT NULL DEFAULT 1000,
        sigmaJitterMs INTEGER NOT NULL DEFAULT 200,
        minDelayMs INTEGER NOT NULL DEFAULT 250,
        maxDelayMs INTEGER NOT NULL DEFAULT 10000,
        isAdaptive INTEGER NOT NULL DEFAULT 1,
        hasSpikeDetected INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/jitter_tune.go` integrated with `pkg/pacing/gaussian.go`.

---

### Improvement 653: Bounded Connection Pool Health & Idle Socket Eviction Inspector
- **ID:** `SOCKET-JITTER-653`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Display real-time active, idle, and half-open TCP socket counts across the outbound SMTP connection pool, enforcing bounded caps and proactive TCP keepalive/eviction sweeps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live socket gauge indicating pool capacity utilization (e.g. `24/64 sockets`), average idle lifetime, and closed zombie connections.
    - Lipgloss token palette: Utilization bar (`#10b981`, green `[████░░░░░░] 38%`), Eviction sweep (`#06b6d4`, cyan pulse `[SWEEP OK]`).
    - Model view snippet:
      ```text
      ┌─ [BOUNDED SMTP SOCKET POOL & EVICTION INSPECTOR] ──────────────────────────────────────────┐
      │ POOL IDENTIFIER │ ACTIVE │ IDLE │ MAX CAP │ UTILIZATION │ IDLE TIMEOUT │ EVICTED (1H)        │
      │ pool-google     │ 8      │ 4    │ 32      │ 37.5%       │ 45s          │ 14 sockets          │
      │ pool-yahoo      │ 12     │ 2    │ 16      │ 87.5%       │ 30s          │ 22 sockets          │
      │ pool-microsoft  │ 4      │ 6    │ 24      │ 41.7%       │ 60s          │ 6 sockets           │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [Active Pool Conns: 36 / 128 max] [Zero Leaks Detected] [Keepalive Frequency: 15s]          │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Circular donut chart showing active vs idle socket allocation with instantaneous eviction flush button.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup socket-pool`
  - Flags:
    - `--pool-id <string>`: Inspect specific connection pool (default: `all`).
    - `--max-conns <int>`: Enforce maximum bounded concurrent sockets (default: `64`).
    - `--evict-idle-now`: Manually purge all idle TCP sockets exceeding timeout.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SocketPoolMetrics struct {
        PoolID         string    `json:"poolId"`
        ActiveConns    int       `json:"activeConns"`
        IdleConns      int       `json:"idleConns"`
        MaxCapacity    int       `json:"maxCapacity"`
        EvictedCount   int       `json:"evictedCount"`
        UtilizationPct float64   `json:"utilizationPct"`
        IsNearLimit    bool      `json:"isNearLimit"`
        HasStaleSockets bool     `json:"hasStaleSockets"`
        LastEvictedAt  time.Time `json:"lastEvictedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SocketPoolSnapshots (
        snapshotId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        activeConns INTEGER NOT NULL DEFAULT 0,
        idleConns INTEGER NOT NULL DEFAULT 0,
        maxCapacity INTEGER NOT NULL DEFAULT 64,
        evictedCount INTEGER NOT NULL DEFAULT 0,
        isNearLimit INTEGER NOT NULL DEFAULT 0,
        hasStaleSockets INTEGER NOT NULL DEFAULT 0,
        capturedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_pool_time ON SocketPoolSnapshots(poolId, capturedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/socket_pool.go` connected to `pkg/netpoll/conn_pool.go`.

---

### Improvement 654: TLS 1.3 Cipher Suite & 0-RTT Handshake Duration Profiler
- **ID:** `SOCKET-JITTER-654`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Benchmark and profile TLS handshake durations, ALPN negotiation, and session ticket resumption (0-RTT / 1-RTT) across target mail exchangers in an interactive terminal audit.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Security and cryptographic profile table outlining negotiated cipher (`TLS_AES_256_GCM_SHA384`), curve (`X25519`), resumption status, and handshake duration in ms.
    - Lipgloss token palette: Fast handshake (<25ms `#10b981`, green), Resumption active (`#8b5cf6`, purple `[0-RTT TICKET]`).
    - Model view snippet:
      ```text
      ┌─ [TLS 1.3 CIPHER & HANDSHAKE PROFILER] ────────────────────────────────────────────────────┐
      │ TARGET MX            │ CIPHER SUITE               │ KEY EXCHANGE │ HANDSHAKE │ RESUMPTION   │
      │ aspmx.l.google.com   │ TLS_AES_256_GCM_SHA384     │ X25519       │ 18.2ms    │ [0-RTT OK]   │
      │ mta5.am0.yahoodns    │ TLS_CHACHA20_POLY1305_SHA  │ X25519       │ 24.5ms    │ [1-RTT NEW]  │
      │ outlook-com.olc.prot │ TLS_AES_128_GCM_SHA256     │ SecP256r1    │ 36.1ms    │ [0-RTT OK]   │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Supported TLS Versions: TLS 1.3 (Primary), TLS 1.2 (Fallback) | Min Cipher Grade: A+      │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Visual TLS handshake waterfall diagram dissecting ClientHello, ServerHello, and Finished states.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup tls-profile`
  - Flags:
    - `--host <string>`: Target hostname or MX endpoint to test.
    - `--force-tls13`: Test strict TLS 1.3 compatibility only.
    - `--enable-0rtt`: Test early data session ticket resumption.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TLSProfileResult struct {
        TargetHost      string  `json:"targetHost"`
        TLSVersion      string  `json:"tlsVersion"`
        CipherSuite     string  `json:"cipherSuite"`
        KeyExchange     string  `json:"keyExchange"`
        HandshakeTimeMs float64 `json:"handshakeTimeMs"`
        IsSessionResumed bool   `json:"isSessionResumed"`
        IsZeroRTTReady  bool    `json:"isZeroRttReady"`
        HasModernCipher bool    `json:"hasModernCipher"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TLSProfiles (
        profileId TEXT PRIMARY KEY,
        targetHost TEXT NOT NULL,
        tlsVersion TEXT NOT NULL,
        cipherSuite TEXT NOT NULL,
        handshakeTimeMs REAL NOT NULL DEFAULT 0.0,
        isSessionResumed INTEGER NOT NULL DEFAULT 0,
        isZeroRttReady INTEGER NOT NULL DEFAULT 0,
        hasModernCipher INTEGER NOT NULL DEFAULT 1,
        testedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/tls_profile.go`.

---

### Improvement 655: Kernel TCP Buffer Tuning & BBR Congestion Control Diagnostic
- **ID:** `SOCKET-JITTER-655`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Audit and recommend OS kernel network parameters (`tcp_congestion_control=bbr`, `tcp_nodelay=1`, `rmem`/`wmem` buffer scaling) for low-latency warm-up dispatch nodes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - System kernel diagnostic table comparing current sysctl/socket options with high-throughput warm-up recommendations.
    - Color tokens: Optimal setting (`#10b981`, green check `[PASS]`), Sub-optimal (`#eab308`, amber `[SUB-OPTIMAL]`).
    - Model view snippet:
      ```text
      ┌─ [KERNEL TCP & BBR CONGESTION CONTROL DIAGNOSTIC] ─────────────────────────────────────────┐
      │ PARAMETER                  │ CURRENT VALUE      │ RECOMMENDED VALUE   │ AUDIT RESULT         │
      │ TCP Congestion Algorithm   │ bbr                │ bbr                 │ [PASS] OPTIMAL       │
      │ TCP_NODELAY (Nagle Algo)   │ 1 (Disabled)       │ 1 (Disabled)        │ [PASS] ZERO DELAY    │
      │ SO_RCVBUF / SO_SNDBUF      │ 262,144 bytes      │ 262,144 bytes       │ [PASS] TUNED         │
      │ TCP Fast Open (TFO)        │ 3 (In/Out)         │ 3 (In/Out)          │ [PASS] ZERO-RTT      │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [All 4 Network Sockets Tuned for Sub-Millisecond Dispatch] [Action: 'a' Apply Recommended] │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - System tuning checklist with copy-pasteable Linux sysctl and Windows PowerShell netsh configuration snippets.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup kernel-net-tune`
  - Flags:
    - `--audit-only`: Inspect current kernel TCP socket parameters without applying changes.
    - `--apply`: Automatically apply optimized socket flags to current process socket listeners.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type KernelNetAudit struct {
        CongestionAlgorithm string `json:"congestionAlgorithm"`
        IsBBRActive         bool   `json:"isBbrActive"`
        IsTCPNoDelayEnabled bool   `json:"isTcpNoDelayEnabled"`
        HasOptimalBuffers   bool   `json:"hasOptimalBuffers"`
        HasFastOpenSupport  bool   `json:"hasFastOpenSupport"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS KernelNetAudits (
        auditId TEXT PRIMARY KEY,
        congestionAlgorithm TEXT NOT NULL,
        isBbrActive INTEGER NOT NULL DEFAULT 1,
        isTcpNoDelayEnabled INTEGER NOT NULL DEFAULT 1,
        hasOptimalBuffers INTEGER NOT NULL DEFAULT 1,
        hasFastOpenSupport INTEGER NOT NULL DEFAULT 1,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/kernel_tune.go`.

---

### Improvement 656: Multi-Path DNS Resolution & Unbound Local Resolver Cache HUD
- **ID:** `SOCKET-JITTER-656`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Benchmark and profile multi-path DNS resolution timings across Cloudflare, Google, and embedded local Unbound caching resolvers to eliminate DNS RTT overhead prior to TCP connection.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time DNS resolution benchmark table comparing lookup latency (MX, A, AAAA, TXT/SPF) across upstream resolvers.
    - Lipgloss tokens: Local cache (`<1ms` `#10b981`, bold green `[L-CACHE: 0.4ms]`), Remote resolver (`15-30ms` `#06b6d4`, cyan).
    - Model view snippet:
      ```text
      ┌─ [MULTI-PATH DNS RESOLVER & LOCAL CACHE HUD] ──────────────────────────────────────────────┐
      │ QUERY DOMAIN   │ RECORD │ UNBOUND (LOCAL) │ CLOUDFLARE (1.1.1.1) │ GOOGLE (8.8.8.8) │ STATUS │
      │ google.com     │ MX     │ 0.3ms (CACHE)   │ 14.2ms               │ 16.8ms           │ [FAST] │
      │ yahoo.com      │ MX     │ 0.4ms (CACHE)   │ 18.5ms               │ 21.0ms           │ [FAST] │
      │ outlook.com    │ MX     │ 0.3ms (CACHE)   │ 16.1ms               │ 19.4ms           │ [FAST] │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Local DNS Cache Hit Ratio: 99.4% | Total DNS Latency Saved: 14.8 seconds / hour            │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live DNS performance graph with cache hit/miss ratio gauge and one-click flush.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup dns-hud`
  - Flags:
    - `--test-domain <string>`: Benchmark lookup for specific domain.
    - `--flush-cache`: Flush internal DNS resolver memory cache.
    - `--show-misses`: Highlight DNS requests that missed local cache.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DNSResolutionMetrics struct {
        QueryHost       string    `json:"queryHost"`
        RecordType      string    `json:"recordType"`
        LookupDuration  float64   `json:"lookupDurationMs"`
        ResolverSource  string    `json:"resolverSource"`
        IsCacheHit      bool      `json:"isCacheHit"`
        HasDNSSecValid  bool      `json:"hasDnsSecValid"`
        ResolvedAt      time.Time `json:"resolvedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DNSResolutionLogs (
        logId TEXT PRIMARY KEY,
        queryHost TEXT NOT NULL,
        recordType TEXT NOT NULL,
        lookupDurationMs REAL NOT NULL DEFAULT 0.0,
        resolverSource TEXT NOT NULL,
        isCacheHit INTEGER NOT NULL DEFAULT 1,
        hasDnsSecValid INTEGER NOT NULL DEFAULT 1,
        resolvedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_dns_host ON DNSResolutionLogs(queryHost);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/cmdemail/dns_hud.go` leveraging `pkg/netpoll/dns_resolver.go`.

---

### Improvement 657: Socket Buffer Memory Leak & Goroutine Leak Sentinel
- **ID:** `SOCKET-JITTER-657`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Continuously monitor memory footprint (`pprof` heap allocs) and active Goroutine counts tied to open socket connections, catching hanging read/write deadlines before leaks develop.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time diagnostic tracker reporting active netpoll Goroutines, buffer pool byte allocations, and deadlocked socket candidates.
    - Lipgloss token palette: Normal operation (`#10b981`, green `[NOMINAL]`), Stalled Goroutine (`#ef4444`, red `[STALLED]`).
    - Model view snippet:
      ```text
      ┌─ [SOCKET MEMORY FOOTPRINT & GOROUTINE SENTINEL] ───────────────────────────────────────────┐
      │ RUNTIME COMPONENT     │ ACTIVE COUNT │ HEAP ALLOC   │ DEADLINE VIOLATIONS │ HEALTH STATE     │
      │ SMTP Dispatch Routines│ 42 routines  │ 12.4 MB      │ 0 violations        │ [NOMINAL]        │
      │ TLS Handshake Workers │ 8 routines   │ 4.2 MB       │ 0 violations        │ [NOMINAL]        │
      │ Event Tail Buffers    │ 2 routines   │ 1.1 MB       │ 0 violations        │ [NOMINAL]        │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Fleet Goroutines: 52 | Memory Buffer Cap: 64 MB | Zero Goroutine Leaks in 72h continuous  │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Memory alloc timeline with active Goroutine thread counter graph.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup socket-sentinel`
  - Flags:
    - `--max-goroutines <int>`: Alert threshold for concurrent netpoll workers (default: `200`).
    - `--dump-stack`: Output Goroutine stack traces to terminal if leak detected.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SocketSentinelReport struct {
        GoroutineCount   int       `json:"goroutineCount"`
        AllocatedHeapMB  float64   `json:"allocatedHeapMb"`
        StalledCount     int       `json:"stalledCount"`
        IsMemoryStable   bool      `json:"isMemoryStable"`
        HasLeakedThreads bool      `json:"hasLeakedThreads"`
        CapturedAt       time.Time `json:"capturedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SocketSentinelAudits (
        auditId TEXT PRIMARY KEY,
        goroutineCount INTEGER NOT NULL DEFAULT 0,
        allocatedHeapMb REAL NOT NULL DEFAULT 0.0,
        stalledCount INTEGER NOT NULL DEFAULT 0,
        isMemoryStable INTEGER NOT NULL DEFAULT 1,
        hasLeakedThreads INTEGER NOT NULL DEFAULT 0,
        capturedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured in `cli/cmdemail/socket_sentinel.go`.

---

### Improvement 658: Packet Loss & TCP Retransmission Rate Terminal Sparkline
- **ID:** `SOCKET-JITTER-658`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Parse TCP socket info (`tcp_info`) to extract real-time segment retransmission counts and packet loss ratios, visualizing degradation via terminal micro-sparklines.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live row view per ISP showing total segments sent, retransmitted packets, and 60-second loss sparklines (` ▂▃▅`).
    - Lipgloss tokens: 0% Retrans (`#10b981`, green `0.0%`), Retransmission spike (`#ef4444`, red `▲ 3.4%`).
    - Model view snippet:
      ```text
      ┌─ [TCP RETRANSMISSION & PACKET LOSS SPARKLINE HUD] ─────────────────────────────────────────┐
      │ DESTINATION MX            │ SEGMENTS SENT │ RETRANSMITS │ LOSS RATIO │ 60S LOSS SPARKLINE    │
      │ aspmx.l.google.com        │ 128,400       │ 12          │ 0.009%     │                [CLEAN]│
      │ mta5.am0.yahoodns.net     │ 84,200        │ 34          │ 0.040%     │   ▂            [CLEAN]│
      │ mx1.mail.icloud.com       │ 14,100        │ 88          │ 0.624%     │  ▂▃▅▃▂         [SPIKE]│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Maximum Allowable Retransmit Rate: 0.100% | Network Health: EXCELLENT                     │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - TCP packet delivery quality sparklines embedded in destination MX cards.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup tcp-loss`
  - Flags:
    - `--alert-threshold <float>`: Retransmission percentage to trigger notification (default: `0.5`).
    - `--watch`: Persistent interactive Bubble Tea full-screen alternate buffer mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TCPRetransmitSample struct {
        TargetHost      string    `json:"targetHost"`
        SegmentsSent    int64     `json:"segmentsSent"`
        Retransmits     int64     `json:"retransmits"`
        LossRatio       float64   `json:"lossRatio"`
        IsUnderLossCap  bool      `json:"isUnderLossCap"`
        HasSpikeAlert   bool      `json:"hasSpikeAlert"`
        SampledAt       time.Time `json:"sampledAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TCPLossSnapshots (
        snapshotId TEXT PRIMARY KEY,
        targetHost TEXT NOT NULL,
        segmentsSent INTEGER NOT NULL DEFAULT 0,
        retransmits INTEGER NOT NULL DEFAULT 0,
        lossRatio REAL NOT NULL DEFAULT 0.0,
        isUnderLossCap INTEGER NOT NULL DEFAULT 1,
        hasSpikeAlert INTEGER NOT NULL DEFAULT 0,
        sampledAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_tcp_host ON TCPLossSnapshots(targetHost);
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded in `cli/cmdemail/tcp_loss.go`.

---

### Improvement 659: Socket Connection Throttling & TCP Window Starvation Detector
- **ID:** `SOCKET-JITTER-659`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Detect TCP Zero-Window advertisements and receiver buffer exhaustion from destination mail exchangers, immediately signaling queue managers to throttle dispatch.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live alert card showing receive window (rcv_wnd) sizing and flags for zero-window freeze events.
    - Lipgloss token palette: Window healthy (`#10b981`, green `[WINDOW: 64KB]`), Zero window freeze (`#ef4444`, flashing red `[ZERO-WINDOW FREEZE]`).
    - Model view snippet:
      ```text
      ┌─ [TCP WINDOW STARVATION & ZERO-WINDOW DETECTOR] ───────────────────────────────────────────┐
      │ DESTINATION MX            │ CURRENT SND_CWND │ PEER RCV_WND │ WINDOW EVENTS │ ACTION TAKEN   │
      │ aspmx.l.google.com        │ 64 KB            │ 128 KB       │ 0 Zero-Wnd    │ FULL SPEED     │
      │ mta5.am0.yahoodns.net     │ 32 KB            │ 64 KB        │ 0 Zero-Wnd    │ FULL SPEED     │
      │ slow-mta.corporate.org    │ 4 KB             │ 0 KB         │ 2 Zero-Wnd    │ AUTO-THROTTLED │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Zero-Window Protection: ACTIVE | Starvation Incidents 24h: 1 | Recovery Time: 4.2s        │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Visual TCP window flow visualizer with automated backpressure triggers.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup tcp-window`
  - Flags:
    - `--host <string>`: Target MX host to inspect.
    - `--auto-throttle`: Automatically throttle dispatch concurrency when receiver window drops below 4KB (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TCPWindowAudit struct {
        TargetHost        string    `json:"targetHost"`
        SenderCWND        int       `json:"senderCwnd"`
        ReceiverWindow    int       `json:"receiverWindow"`
        ZeroWindowCount   int       `json:"zeroWindowCount"`
        IsWindowHealthy   bool      `json:"isWindowHealthy"`
        HasAutoThrottled  bool      `json:"hasAutoThrottled"`
        AuditedAt         time.Time `json:"auditedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TCPWindowAudits (
        auditId TEXT PRIMARY KEY,
        targetHost TEXT NOT NULL,
        senderCwnd INTEGER NOT NULL DEFAULT 65536,
        receiverWindow INTEGER NOT NULL DEFAULT 65536,
        zeroWindowCount INTEGER NOT NULL DEFAULT 0,
        isWindowHealthy INTEGER NOT NULL DEFAULT 1,
        hasAutoThrottled INTEGER NOT NULL DEFAULT 0,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/cmdemail/tcp_window.go`.

---

### Improvement 660: Integrated Bounded Socket & Jitter Master Control HUD
- **ID:** `SOCKET-JITTER-660`
- **Cluster:** 66 (Bounded Socket Pool Latency Heatmap & Jitter Compensator TUI)
- **Goal:** Consolidate latency heatmaps, Gaussian jitter regulators, pool metrics, and TCP diagnostics into a single unified full-screen Bubble Tea master dashboard.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Quad-split full terminal dashboard: (1) Latency Heatmap Top-Left, (2) Jitter Pacing Top-Right, (3) Socket Pool Utilization Bottom-Left, (4) TCP Quality Sentinel Bottom-Right.
    - Lipgloss token palette: Master border (`#6366f1`, indigo), Metric badges (`#10b981`, green), Active alert flags (`#ef4444`, red).
    - Model view snippet:
      ```text
      ┌─ [SOCKET POOL & JITTER MASTER CONTROL HUD] ───────────────────────────────────────────────┐
      │ [1: LATENCY HEATMAP]                │ [2: GAUSSIAN JITTER REGULATOR]                       │
      │ Google MX : 24ms ■■■■■■■■■■■■■■■■   │ Base Delay: 850ms | Jitter σ: ±120ms [HUMAN PACING]  │
      │ Yahoo MX  : 84ms ■■■■■■■■■■■■■■■■   │ Inter-Packet Delay Distribution: OPTIMAL (Gaussian)  │
      ├─────────────────────────────────────┼──────────────────────────────────────────────────────┤
      │ [3: SOCKET POOL UTILIZATION]        │ [4: TCP HEALTH & LOSS SENTINEL]                      │
      │ Active: 36 | Idle: 12 | Max: 128    │ Packet Loss: 0.012% | Zero-Window: NONE              │
      │ Utilization: [███░░░░░░░] 28.1%     │ TLS 1.3 0-RTT Tickets: 84% resumed                   │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [Keybindings: 'h': Heatmap View, 'j': Jitter Tuning, 'p': Pool Flush, 'q': Return to CLI]  │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Unified 4-tile responsive grid dashboard with live WebSocket telemetry synchronization.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup socket-master`
  - Flags:
    - `--refresh-rate <string>`: HUD refresh interval (default: `1s`).
    - `--export-json`: Export current quad-metric snapshot to JSON and exit.
    - `--watch`: Persistent interactive Bubble Tea full-screen alternate buffer mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SocketMasterHUDState struct {
        AverageFleetRTTMs float64             `json:"averageFleetRttMs"`
        ActiveSockets     int                 `json:"activeSockets"`
        MaxCapacity       int                 `json:"maxCapacity"`
        CurrentJitterMode string              `json:"currentJitterMode"`
        GlobalLossRatio   float64             `json:"globalLossRatio"`
        IsFleetHealthy    bool                `json:"isFleetHealthy"`
        HasActiveAlarms   bool                `json:"hasActiveAlarms"`
        LastUpdated       time.Time           `json:"lastUpdated"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SocketMasterSnapshots (
        snapshotId TEXT PRIMARY KEY,
        averageFleetRttMs REAL NOT NULL DEFAULT 0.0,
        activeSockets INTEGER NOT NULL DEFAULT 0,
        maxCapacity INTEGER NOT NULL DEFAULT 128,
        currentJitterMode TEXT NOT NULL DEFAULT 'gaussian',
        globalLossRatio REAL NOT NULL DEFAULT 0.0,
        isFleetHealthy INTEGER NOT NULL DEFAULT 1,
        hasActiveAlarms INTEGER NOT NULL DEFAULT 0,
        capturedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/socket_master.go` binding all socket, jitter, and netpoll components.

---
