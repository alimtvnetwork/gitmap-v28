import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { KeyRound, Terminal, Shield, FolderGit2, Settings, Network, Cpu, Download, RefreshCw, Layers } from "lucide-react";

const MOCK_KEYS = [
  { name: "default", path: "~/.ssh/id_rsa", fingerprint: "SHA256:abc123...", created: "2026-03-22" },
  { name: "work", path: "~/.ssh/id_rsa_work", fingerprint: "SHA256:def456...", created: "2026-03-22" },
];

const ListPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap ssh list</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold text-xs mb-1">
        {"  "}SSH Keys (2):
      </div>
      <div className="text-muted-foreground text-xs mb-1">
        {"  "}
        <span className="inline-block w-[100px]">Name</span>
        <span className="inline-block w-[200px]">Path</span>
        <span className="inline-block w-[180px]">Fingerprint</span>
        <span>Created</span>
      </div>
      {MOCK_KEYS.map((k) => (
        <div key={k.name} className="text-terminal-foreground text-xs">
          {"  "}
          <span className="inline-block w-[100px] text-foreground font-semibold">{k.name}</span>
          <span className="inline-block w-[200px] text-muted-foreground">{k.path}</span>
          <span className="inline-block w-[180px] text-primary">{k.fingerprint}</span>
          <span className="text-muted-foreground">{k.created}</span>
        </div>
      ))}
    </div>
  </div>
);

const GenPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap ssh --name work</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-green-400">{"  "}✓ SSH key "work" generated</div>
      <div className="text-terminal-foreground">{"    "}Path:        ~/.ssh/id_rsa_work</div>
      <div className="text-terminal-foreground">{"    "}Fingerprint: SHA256:def456...</div>
      <div className="text-terminal-foreground mt-1">{"    "}Public key:</div>
      <div className="text-primary mt-1">{"  "}ssh-rsa AAAA... user@example.com</div>
      <div className="text-blue-400 mt-2">{"  "}ℹ  Copy the public key above and add it to your Git provider.</div>
    </div>
  </div>
);

const ScanPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap ssh scan</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto">
      <div className="text-cyan-400 font-bold mb-2">SSH Fleet Liveness & Reachability Scan:</div>
      <div className="text-muted-foreground mb-1">
        <span className="inline-block w-[100px]">ALIAS</span>
        <span className="inline-block w-[140px]">IP</span>
        <span className="inline-block w-[80px]">USER</span>
        <span className="inline-block w-[60px]">PORT</span>
        <span className="inline-block w-[90px]">STATUS</span>
        <span className="inline-block w-[80px]">LATENCY</span>
        <span>DETAILS</span>
      </div>
      <div className="text-green-400">
        <span className="inline-block w-[100px] text-foreground font-semibold">devbox</span>
        <span className="inline-block w-[140px] text-muted-foreground">192.168.1.14</span>
        <span className="inline-block w-[80px] text-muted-foreground">alim</span>
        <span className="inline-block w-[60px] text-muted-foreground">22</span>
        <span className="inline-block w-[90px] text-green-400 font-bold">ONLINE</span>
        <span className="inline-block w-[80px]">12ms</span>
        <span className="text-muted-foreground">tcp reachable</span>
      </div>
      <div className="text-red-400">
        <span className="inline-block w-[100px] text-foreground font-semibold">worker-1</span>
        <span className="inline-block w-[140px] text-muted-foreground">192.168.1.20</span>
        <span className="inline-block w-[80px] text-muted-foreground">ubuntu</span>
        <span className="inline-block w-[60px] text-muted-foreground">22</span>
        <span className="inline-block w-[90px] text-red-400 font-bold">OFFLINE</span>
        <span className="inline-block w-[80px]">-</span>
        <span className="text-muted-foreground">connection timeout</span>
      </div>
      <div className="text-terminal-foreground mt-3 pt-2 border-t border-border/40">
        Summary: 1/2 nodes online (in-memory TTL cache: 45s)
      </div>
    </div>
  </div>
);

const ExecPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap ssh exec "gitmap --version"</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto space-y-1">
      <div>
        <span className="text-cyan-400 font-bold">[devbox|192.168.1.14]</span>{" "}
        <span className="text-foreground">gitmap version v6.260.0 linux/amd64</span>
      </div>
      <div>
        <span className="text-cyan-400 font-bold">[worker-1|192.168.1.20]</span>{" "}
        <span className="text-yellow-400">OFFLINE (skipped: connection timeout)</span>
      </div>
      <div className="text-muted-foreground pt-2">SSH Execution Done.</div>
    </div>
  </div>
);

const ComparePreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap ssh compare</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto">
      <div className="text-cyan-400 font-bold mb-2">GitMap Remote Subsystems Architecture Comparison:</div>
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="text-muted-foreground border-b border-border/50">
            <th className="pb-1 pr-3">SUBSYSTEM</th>
            <th className="pb-1 pr-3">PRIMARY FOCUS</th>
            <th className="pb-1 pr-3">JOIN COMMAND</th>
            <th className="pb-1 pr-3">EXEC COMMAND</th>
            <th className="pb-1 pr-3">MONITORING</th>
            <th className="pb-1">BEST USED WHEN</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border/30">
          <tr className="text-green-400">
            <td className="py-2 pr-3 font-semibold">gitmap ssh</td>
            <td className="py-2 pr-3">Direct node-level management</td>
            <td className="py-2 pr-3 text-muted-foreground">gitmap ssh join &lt;u@ip&gt;</td>
            <td className="py-2 pr-3 text-foreground">gitmap ssh exec &lt;cmd&gt;</td>
            <td className="py-2 pr-3">gitmap ssh scan</td>
            <td className="py-2 text-muted-foreground">Ad-hoc commands, install/update, AGY/code open</td>
          </tr>
          <tr className="text-cyan-400">
            <td className="py-2 pr-3 font-semibold">gitmap cluster</td>
            <td className="py-2 pr-3">Multi-node cluster orchestration</td>
            <td className="py-2 pr-3 text-muted-foreground">gitmap cluster node add &lt;ip&gt;</td>
            <td className="py-2 pr-3 text-foreground">gitmap cluster exec &lt;tgt&gt; &lt;cmd&gt;</td>
            <td className="py-2 pr-3">gitmap cluster node ls</td>
            <td className="py-2 text-muted-foreground">K8s bootstrap, cluster recipes, distributed scripts</td>
          </tr>
          <tr className="text-yellow-400">
            <td className="py-2 pr-3 font-semibold">gitmap sc</td>
            <td className="py-2 pr-3">Servers-Clients fleet daemon</td>
            <td className="py-2 pr-3 text-muted-foreground">gitmap sc join &lt;url&gt;</td>
            <td className="py-2 pr-3 text-foreground">gitmap sc exec &lt;cmd&gt;</td>
            <td className="py-2 pr-3">gitmap sc status</td>
            <td className="py-2 text-muted-foreground">Master-worker topology, continuous sync, live telemetry</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
);

const SSHPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3 mb-2">
          <KeyRound className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">SSH Management & Remote Execution</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Generate SSH keys, enroll remote fleet machines, probe node liveness, execute commands, install GitMap, and delegate AGY / VS Code workspaces.
        </p>
      </div>

      {/* Overview */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Shield className="h-5 w-5 text-primary" /> Overview
        </h2>
        <p className="text-muted-foreground mb-4">
          The <code className="text-primary">ssh</code> suite combines local SSH key lifecycle management with direct remote node execution.
          It provides automatic liveness caching (45s TTL), skips offline nodes gracefully, auto-discovers default SSH keys (<code className="text-xs">~/.ssh/id_ed25519</code>, <code className="text-xs">id_rsa</code>),
          and enables one-command remote installation, AGY delegation, and VS Code remote sessions.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: KeyRound, title: "Named Keys & Config", desc: "Automatic ~/.ssh/config generation and clone integration" },
            { icon: Network, title: "Liveness & Scan", desc: "In-memory reachability cache (45s) skipping offline nodes" },
            { icon: Cpu, title: "Remote Command Exec", desc: "Live terminal streaming with target filtering and default auth" },
            { icon: Download, title: "Remote GitMap Install", desc: "Auto-detects missing binary; installs or updates to latest" },
            { icon: RefreshCw, title: "AGY & VS Code Remote", desc: "Open remote folders directly in Google Antigravity or VS Code" },
            { icon: Layers, title: "Triad Architecture", desc: "Subsystem comparison matrix: ssh vs cluster vs sc" },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Subcommands */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Terminal className="h-5 w-5 text-primary" /> Complete Subcommands Reference
        </h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Subcommand</th>
                <th className="text-left px-4 py-2 font-medium">Alias</th>
                <th className="text-left px-4 py-2 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                { cmd: "(default)", alias: "create", desc: "Generate a new SSH key pair" },
                { cmd: "cat", alias: "view, v", desc: "Display the public key" },
                { cmd: "copy", alias: "cp", desc: "Copy public key to OS clipboard" },
                { cmd: "list", alias: "ls", desc: "List all stored SSH keys" },
                { cmd: "delete", alias: "rm", desc: "Delete a key record (optionally files)" },
                { cmd: "config", alias: "—", desc: "Regenerate ~/.ssh/config entries" },
                { cmd: "join", alias: "sj", desc: "Enroll machine into SSH registry with encrypted password/key" },
                { cmd: "exec", alias: "se", desc: "Execute command across online nodes with liveness check" },
                { cmd: "scan", alias: "—", desc: "Probe reachability and latency across SSH fleet nodes" },
                { cmd: "install", alias: "i", desc: "Install or update GitMap on remote node(s)" },
                { cmd: "update", alias: "u", desc: "Update GitMap binary across remote fleet" },
                { cmd: "agy", alias: "—", desc: "Run Antigravity CLI or open remote folder with AppError" },
                { cmd: "code", alias: "—", desc: "Open remote folder in VS Code via SSH Remote" },
                { cmd: "compare", alias: "matrix", desc: "Display Triad comparison matrix: SSH vs Cluster vs SC" },
              ].map((s) => (
                <tr key={s.cmd} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{s.cmd}</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">{s.alias}</td>
                  <td className="px-4 py-2 text-muted-foreground">{s.desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Liveness & Scan Section */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Fleet Liveness Scan</h2>
        <p className="text-muted-foreground mb-4">
          Probes all registered SSH nodes in parallel. Reachability results are stored in an in-memory TTL cache (45 seconds),
          preventing consecutive command hangs on dead or suspended virtual machines.
        </p>
        <ScanPreview />
        <CodeBlock code="gitmap ssh scan" />
      </section>

      {/* Remote Execution Section */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Remote Command Execution with Target Selection</h2>
        <p className="text-muted-foreground mb-4">
          Commands run concurrently across running nodes. Offline machines are detected and skipped cleanly with status notices.
          Target individual machines by alias (<code className="text-xs">--target devbox</code>), IP (<code className="text-xs">--ip 192.168.1.14</code>), or positional target (<code className="text-xs">gitmap ssh exec devbox "uptime"</code>).
        </p>
        <ExecPreview />
        <CodeBlock code={`# Execute on all running machines
gitmap ssh exec "gitmap --version"

# Target specific machine by alias or IP
gitmap ssh exec --target devbox "uname -a"
gitmap ssh exec 192.168.1.14 "gitmap status"

# Run IP inspection across fleet
gitmap ssh exec ip`} />
      </section>

      {/* Install & Update Section */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Remote GitMap Install & Update</h2>
        <p className="text-muted-foreground mb-4">
          The <code className="text-primary">ssh install</code> command checks if GitMap exists on the remote node via <code className="text-xs">gitmap --version</code>.
          If missing, it runs the official cross-platform installer (curl bash on Linux/macOS, PowerShell on Windows). If already installed, it triggers an update.
        </p>
        <CodeBlock code={`# Install on a single machine
gitmap ssh install gitmap devbox

# Install across all fleet machines
gitmap ssh install gitmap all

# Update GitMap binary across fleet
gitmap ssh update gitmap all`} />
      </section>

      {/* AGY & VS Code Remote Section */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Remote AGY & VS Code Delegation</h2>
        <p className="text-muted-foreground mb-4">
          Open remote workspaces or run Antigravity CLI tools directly over SSH. If the remote connection fails or the machine is offline,
          structured <code className="text-primary">AppError</code> diagnostics with full stack traces are displayed for immediate troubleshooting.
        </p>
        <CodeBlock code={`# Open remote folder in Google Antigravity
gitmap ssh agy open /var/www/my-project --target devbox

# Run AGY CLI tool remotely
gitmap ssh agy "agy --version" --target devbox

# Open remote workspace in VS Code via SSH Remote
gitmap ssh code open /opt/app --target devbox`} />
      </section>

      {/* Triad Architecture Comparison */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Layers className="h-5 w-5 text-primary" /> Subsystems Architecture Comparison (Triad)
        </h2>
        <p className="text-muted-foreground mb-4">
          GitMap provides three distinct remote management subsystems tailored for different operational workflows.
          Run <code className="text-primary">gitmap ssh compare</code> in your terminal anytime to display this comparison matrix.
        </p>
        <ComparePreview />
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
          <div className="rounded-lg border border-border p-4 bg-card">
            <h3 className="font-semibold text-sm text-green-400 mb-1">gitmap ssh</h3>
            <p className="text-xs text-muted-foreground">
              Direct, lightweight node management over standard SSH. Best for workstations, ad-hoc maintenance, liveness scans, GitMap installation, and AGY / VS Code remote launches.
            </p>
          </div>
          <div className="rounded-lg border border-border p-4 bg-card">
            <h3 className="font-semibold text-sm text-cyan-400 mb-1">gitmap cluster</h3>
            <p className="text-xs text-muted-foreground">
              Role-based infrastructure orchestration (<code className="text-xs">control</code> vs <code className="text-xs">workers</code>), Netplan IP configuration, automated user provisioning, and end-to-end Kubernetes lifecycle.
            </p>
          </div>
          <div className="rounded-lg border border-border p-4 bg-card">
            <h3 className="font-semibold text-sm text-yellow-400 mb-1">gitmap sc (servers-clients)</h3>
            <p className="text-xs text-muted-foreground">
              High-speed broadcast fan-out across entire fleets with bounded concurrency pools, multi-shell support (bash/ps/cmd), and continuous daemon synchronization.
            </p>
          </div>
        </div>
      </section>

      {/* Key Management Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">SSH Key Generation & Stored Keys</h2>
        <GenPreview />
        <CodeBlock code="gitmap ssh --name work --path ~/.ssh/id_rsa_work" />
        <ListPreview />
      </section>
    </div>
  </DocsLayout>
);

export default SSHPage;
