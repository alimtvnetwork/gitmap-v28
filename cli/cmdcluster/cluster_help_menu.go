package cmdcluster

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderClusterHelp displays the styled two-column Cluster help menu.
func RenderClusterHelp() {
	termout.RenderMenu(buildClusterHelpMenu())
}

func buildClusterHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Cluster Fleet & K8s Orchestration (gitmap cluster)",
		UsageLines: []string{
			"gitmap cluster [command] [args]",
			"gitmap cluster exec <target> <command>",
			"gitmap cluster node <recipe> [args]",
			"gitmap cluster k8s <stage> [args]",
		},
		Sections: []termout.HelpSection{
			buildClusterNodeSection(),
			buildClusterOpsSection(),
			buildClusterK8sSection(),
			buildClusterTriadSection(),
		},
		FooterFlags: buildClusterFooterFlags(),
		Tips: []string{
			"Use 'gitmap cluster exec all <cmd>' to broadcast commands across the fleet.",
			"Run 'gitmap cluster compare' to inspect SSH vs Cluster vs SC differences.",
			"Run 'gitmap cluster k8s init control' to initialize a Kubernetes control plane.",
		},
	}
}

func buildClusterFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--sudo", Description: "Execute remote commands with elevated privileges"},
		{Command: "--parallel <N>", Description: "Concurrency worker count for parallel node execution"},
		{Command: "--json", Description: "Output node list or telemetry in structured JSON"},
		{Command: "--auth", Description: "Automatically deploy SSH key to remote authorized_keys"},
		{Command: "-h, --help", Description: "Show this cluster help menu"},
	}
}

func buildClusterNodeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Node & Topology Management",
		Entries: []termout.CommandEntry{
			{Command: "nodes (ls)", Description: "List all registered cluster nodes, roles, and status"},
			{Command: "add <u@ip> [alias]", Description: "Enroll and join remote node into cluster registry"},
			{Command: "node <sub>", Description: "Node provisioning recipes (set-ip, install-base, user)", HasSubcommands: true},
			{Command: "rm <alias|ip>", Description: "Remove node from cluster topology registry"},
			{Command: "import <file.json>", Description: "Import cluster topology from JSON configuration"},
			{Command: "export", Description: "Export cluster node registry to JSON or CSV"},
		},
	}
}

func buildClusterOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Cluster Operations & Execution",
		Entries: []termout.CommandEntry{
			{Command: "exec <target> <cmd>", Description: "Execute command across targeted nodes (control/workers/all)"},
			{Command: "run-script <target> <f>", Description: "Deploy and execute local script across targeted nodes"},
			{Command: "status (ping, health)", Description: "Check cluster connectivity, port reachability, and health"},
			{Command: "bootstrap <target>", Description: "Deploy RSA keys and configure passwordless sudo"},
			{Command: "install [target]", Description: "Install or update GitMap binary across cluster nodes"},
		},
	}
}

func buildClusterK8sSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Kubernetes (K8s) Lifecycle",
		Entries: []termout.CommandEntry{
			{Command: "k8s prereq <target>", Description: "Configure kernel modules, sysctl, and swap off"},
			{Command: "k8s install <target>", Description: "Install CRI-O container runtime, kubelet, and kubeadm"},
			{Command: "k8s init <target>", Description: "Initialize control plane with pod network CIDR"},
			{Command: "k8s cni <target>", Description: "Deploy Weave Net or Calico cluster network overlay"},
			{Command: "k8s join <target>", Description: "Generate token and join worker nodes to cluster"},
			{Command: "k8s nfs <target>", Description: "Setup NFS export directory and shared storage"},
			{Command: "k8s helm-install <t>", Description: "Install Helm v3 package manager binary on control plane"},
			{Command: "k8s reset <target>", Description: "Teardown Kubernetes setup and clean node iptables/IPVS"},
		},
	}
}

func buildClusterTriadSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Subsystems & Audit",
		Entries: []termout.CommandEntry{
			{Command: "compare (matrix)", Description: "Display architecture matrix: SSH vs Cluster vs SC"},
			{Command: "history [run-id]", Description: "Inspect audit trail of past cluster executions"},
		},
	}
}
