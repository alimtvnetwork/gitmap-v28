import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { KeyRound, Terminal, Shield, FolderGit2, Settings } from "lucide-react";

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

const SSHPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3 mb-2">
          <KeyRound className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">SSH Key Management</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Generate, store, and manage SSH keys for Git authentication directly from the CLI.
        </p>
      </div>

      {/* Overview */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Shield className="h-5 w-5 text-primary" /> Overview
        </h2>
        <p className="text-muted-foreground mb-4">
          The <code className="text-primary">ssh</code> command provides one-command SSH key generation
          with automatic database storage, public key display, and <code>~/.ssh/config</code> management.
          Keys are identified by name and integrate with <code>gitmap clone --ssh-key</code>.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: KeyRound, title: "Named Keys", desc: "Label keys (work, personal) for easy reference" },
            { icon: Settings, title: "Auto Config", desc: "~/.ssh/config managed automatically" },
            { icon: FolderGit2, title: "Clone Integration", desc: "Use --ssh-key flag with clone command" },
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
          <Terminal className="h-5 w-5 text-primary" /> Subcommands
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
                { cmd: "(default)", alias: "—", desc: "Generate a new SSH key pair" },
                { cmd: "cat", alias: "—", desc: "Display the public key" },
                { cmd: "list", alias: "ls", desc: "List all stored SSH keys" },
                { cmd: "delete", alias: "rm", desc: "Delete a key record (optionally files)" },
                { cmd: "config", alias: "—", desc: "Regenerate ~/.ssh/config entries" },
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

      {/* Generate Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Generate a Named Key</h2>
        <GenPreview />
        <CodeBlock code="gitmap ssh --name work --path ~/.ssh/id_rsa_work" />
      </section>

      {/* List Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">List Stored Keys</h2>
        <ListPreview />
      </section>

      {/* Flags */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Flags</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Flag</th>
                <th className="text-left px-4 py-2 font-medium">Short</th>
                <th className="text-left px-4 py-2 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                { flag: "--name", short: "-n", desc: "Key label in database (default: 'default')" },
                { flag: "--path", short: "-p", desc: "Private key file path" },
                { flag: "--email", short: "-e", desc: "Email comment for the key" },
                { flag: "--force", short: "-f", desc: "Skip regeneration prompt" },
                { flag: "--host", short: "—", desc: "Git provider hostname (default: github.com)" },
                { flag: "--confirm", short: "—", desc: "Skip interactive confirmation prompt" },
                { flag: "--files", short: "—", desc: "Delete key files from disk (delete only)" },
                { flag: "--json", short: "—", desc: "Output key list as JSON (list only)" },
                { flag: "--ssh-key", short: "-K", desc: "SSH key name for clone command" },
              ].map((f) => (
                <tr key={f.flag} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{f.flag}</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">{f.short}</td>
                  <td className="px-4 py-2 text-muted-foreground">{f.desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Examples */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <div className="space-y-3">
          <CodeBlock code={`# Generate default SSH key
gitmap ssh

# Generate named key for work
gitmap ssh --name work --path ~/.ssh/id_rsa_work

# Generate key for GitLab with custom host
gitmap ssh --name gitlab --host gitlab.com --email user@company.com

# Non-interactive generation (CI/scripting)
gitmap ssh --name deploy --confirm

# Display public key for copying
gitmap ssh cat --name work

# List keys as JSON for scripting
gitmap ssh list --json

# Clone using a specific SSH key
gitmap clone repos.json --ssh-key work

# Delete key and remove files from disk
gitmap ssh delete --name work --files

# Regenerate SSH config entries
gitmap ssh config`} />
        </div>
      </section>

      {/* SSH Config */}
      <section>
        <h2 className="text-xl font-semibold mb-3">SSH Config Management</h2>
        <p className="text-muted-foreground mb-4">
          When multiple keys exist, gitmap auto-manages a marked block in <code>~/.ssh/config</code>.
          User entries outside the markers are preserved.
        </p>
        <CodeBlock code={`# --- gitmap managed (do not edit) ---
Host github.com-default
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_rsa
    IdentitiesOnly yes

Host github.com-work
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_rsa_work
    IdentitiesOnly yes
# --- end gitmap managed ---`} />
      </section>
    </div>
  </DocsLayout>
);

export default SSHPage;
