import os
import re

file_path = r"d:\work\gitmap\src\data\commands.ts"

with open(file_path, "r", encoding="utf-8") as f:
    content = f.read()

# Add category if not exists
if "key: \"cluster\"" not in content:
    cat_idx = content.find("];\n", content.find("export const Categories: CommandCategory[] = ["))
    if cat_idx != -1:
        new_cat = '  { key: "cluster", label: "Cluster & Delegation", description: "Multi-machine clustering and command delegation", icon: "🕸️" },\n'
        content = content[:cat_idx] + new_cat + content[cat_idx:]


new_commands = """
  {
    name: "servers-clients",
    alias: "sc",
    description: "Broadcast shell commands, git operations, or installations across all server and client nodes in the cluster.",
    usage: "gitmap servers-clients <subcommand> [args] [flags]",
    category: "cluster",
    flags: [
      { flag: "--except <list>", description: "Exclude nodes by ID, IP, or trailing IP octet" },
      { flag: "--yes, -Y", description: "Bypass preflight confirmation prompt" }
    ],
    examples: [
      { command: 'gitmap servers-clients ps "Get-Service | Where Status -eq Running"', description: "Execute PowerShell command across cluster" },
      { command: 'gitmap servers-clients cmd "ipconfig /all" --except 24,151', description: "Execute cmd with exclusions" },
      { command: 'gitmap servers-clients install "git,nodejs,dotnet" --except 2', description: "Install packages simultaneously" },
      { command: 'gitmap servers-clients pull --all', description: "Run git pull --all on all nodes" },
      { command: 'gitmap servers-clients proj "api-backend" run --except 2', description: "Run project automation" }
    ]
  },
  {
    name: "clients",
    description: "Broadcast shell commands, git operations, or lifecycle actions across all client nodes in the cluster (excludes the server/orchestrator).",
    usage: "gitmap clients <subcommand> [args] [flags]",
    category: "cluster",
    flags: [
      { flag: "--except <list>", description: "Exclude nodes by ID, IP, or trailing IP octet" },
      { flag: "--ip <list>", description: "Run only on the named IPs" },
      { flag: "--id <list>", description: "Run only on the named Display IDs" },
      { flag: "--force-lifecycle", description: "Require for restart/shutdown/logoff" }
    ],
    examples: [
      { command: 'gitmap clients ps "Get-DiskUsage C:\\\\" --except 24,151', description: "Run PS command on clients with exclusions" },
      { command: 'gitmap clients cmd "whoami", ps "Get-Date" --ip 192.168.1.10,192.168.1.11', description: "Chain commands sequentially" },
      { command: 'gitmap clients restart --except 1', description: "Restart machines (requires password)" }
    ]
  },
  {
    name: "cluster history",
    description: "View the audit trail of past cluster executions.",
    usage: "gitmap cluster history [RunRef]",
    category: "cluster",
    examples: [
      { command: 'gitmap cluster history', description: "List all past cluster runs" },
      { command: 'gitmap cluster history RUN-20260817-001', description: "Expand a specific run ID for per-node outcomes" }
    ]
  },
  {
    name: "cluster nodes",
    description: "List registered nodes in the cluster.",
    usage: "gitmap cluster nodes [--json]",
    category: "cluster",
    flags: [
      { flag: "--json", description: "Emit listing as a JSON array" }
    ],
    examples: [
      { command: 'gitmap cluster nodes', description: "Display node list in a table" },
      { command: 'gitmap cluster nodes --json', description: "Export node list as JSON" }
    ]
  },
  {
    name: "cluster export",
    description: "Export the cluster node registry to a JSON or CSV file.",
    usage: "gitmap cluster export [--format json|csv] [--output <file>]",
    category: "cluster",
    flags: [
      { flag: "--format <type>", description: "json (default) or csv" },
      { flag: "--output <file>", description: "File path to write" }
    ],
    examples: [
      { command: 'gitmap cluster export --format json --output nodes.json', description: "Export as JSON" },
      { command: 'gitmap cluster export --format csv > nodes.csv', description: "Export as CSV to stdout" }
    ]
  },
  {
    name: "cluster import",
    description: "Import a cluster node registry file (JSON or CSV).",
    usage: "gitmap cluster import <file>",
    category: "cluster",
    examples: [
      { command: 'gitmap cluster import nodes.json', description: "Merge JSON records into cluster database" }
    ]
  },
  {
    name: "cluster set-password",
    description: "Store a client node password for privileged lifecycle commands (restart, shutdown, logoff).",
    usage: "gitmap cluster set-password --id <id>",
    category: "cluster",
    examples: [
      { command: 'gitmap cluster set-password --id 5', description: "Securely prompt and save bcrypt-hashed password for Node 5" }
    ]
  }
"""

cmd_idx = content.rfind("];")
if cmd_idx != -1:
    content = content[:cmd_idx] + new_commands + content[cmd_idx:]
    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)
    print("Added to commands.ts")
else:
    print("Could not find the end of commands array")
