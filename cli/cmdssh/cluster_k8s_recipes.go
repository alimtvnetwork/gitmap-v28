package cmdssh

import (
	"fmt"
	"strings"
)

const defaultK8sVersion = "1.31"
const defaultHelmVersion = "3.16.2"
const defaultNFSExportDir = "/nfsexport"
const defaultWeaveNetManifest = "https://reweave.azurewebsites.net/k8s/v1.31/net.yaml"
const defaultCalicoManifest = "https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml"

const k8sPrereqScriptTemplate = `#!/bin/bash
set -e
mkdir -p /etc/modules-load.d /etc/sysctl.d
cat <<'EOF' > /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
modprobe overlay
modprobe br_netfilter
cat <<'EOF' > /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward = 1
EOF
sysctl --system
`

const k8sInstallScriptTemplate = `#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive
%smkdir -p /etc/apt/keyrings /etc/apt/sources.list.d
curl -fsSL https://pkgs.k8s.io/core:/stable:/v%s/deb/Release.key | gpg --dearmor --yes -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v%s/deb/ /" > /etc/apt/sources.list.d/kubernetes.list
curl -fsSL https://pkgs.k8s.io/addons:/cri-o:/prerelease:/main/deb/Release.key | gpg --dearmor --yes -o /etc/apt/keyrings/cri-o-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/cri-o-apt-keyring.gpg] https://pkgs.k8s.io/addons:/cri-o:/prerelease:/main/deb/ /" > /etc/apt/sources.list.d/cri-o.list
apt-get update -y
apt-get install -y cri-o kubelet kubeadm kubectl socat apt-transport-https ca-certificates curl gpg sshpass
systemctl enable crio.service
systemctl start crio.service
kubeadm config images pull
`

const k8sInitScriptTemplate = `#!/bin/bash
set -e
systemctl enable crio.service
systemctl start crio.service
%smkdir -p $HOME/.kube
cp -f /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
`

const cniScriptTemplate = `#!/bin/bash
set -e
kubectl apply -f %s
`

const k8sJoinScriptTemplate = `#!/bin/bash
set -e
%s
`

const nfsServerScriptTemplate = `#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y nfs-server
mkdir -p %s
echo '%s *(rw,no_root_squash)' >> /etc/exports
systemctl restart nfs-server
`

const helmInstallScriptTemplate = `#!/bin/bash
set -e
TMP_DIR=$(mktemp -d)
curl -fsSL "https://get.helm.sh/helm-v%s-linux-amd64.tar.gz" -o "${TMP_DIR}/helm.tar.gz"
tar -xzf "${TMP_DIR}/helm.tar.gz" -C "${TMP_DIR}"
mv "${TMP_DIR}/linux-amd64/helm" /usr/local/bin/helm
chmod +x /usr/local/bin/helm
rm -rf "${TMP_DIR}"
`

const helmNFSScriptTemplate = `#!/bin/bash
set -e
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
helm repo update
helm uninstall nfs-provisioner 2>/dev/null || true
helm install nfs-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner --set nfs.server="%s" --set nfs.path="%s"
`

const k8sResetScriptTemplate = `#!/bin/bash
set -e
kubeadm reset --force && rm -rf $HOME/.kube /etc/cni/net.d
`

// GenerateK8sPrereqScript generates bash script for kernel modules and sysctl prerequisites.
func GenerateK8sPrereqScript() string {
	return k8sPrereqScriptTemplate
}

func buildHostnameSnippet(newHostname string) string {
	hasHostname := newHostname != ""
	if hasHostname {
		return fmt.Sprintf("hostnamectl set-hostname \"%s\"\n", newHostname)
	}

	return ""
}

func resolveK8sVersion(version string) string {
	cleanVer := strings.TrimPrefix(strings.TrimSpace(version), "v")
	hasVer := cleanVer != ""
	if hasVer {
		return cleanVer
	}

	return defaultK8sVersion
}

// GenerateK8sInstallScript generates bash script for installing Kubernetes tools and CRI-O.
func GenerateK8sInstallScript(version string, newHostname string) string {
	resolvedVer := resolveK8sVersion(version)
	hostSnippet := buildHostnameSnippet(newHostname)

	return fmt.Sprintf(k8sInstallScriptTemplate, hostSnippet, resolvedVer, resolvedVer)
}

func buildKubeadmInitCmd(podCIDR string) string {
	cleanCIDR := strings.TrimSpace(podCIDR)
	hasCIDR := cleanCIDR != ""
	if hasCIDR {
		return fmt.Sprintf("kubeadm init --pod-network-cidr=%s\n", cleanCIDR)
	}

	return "kubeadm init\n"
}

// GenerateK8sInitScript generates bash script for initializing Kubernetes master node.
func GenerateK8sInitScript(podCIDR string) string {
	initCmd := buildKubeadmInitCmd(podCIDR)

	return fmt.Sprintf(k8sInitScriptTemplate, initCmd)
}

func resolveCNIManifest(plugin string) string {
	isCalico := strings.EqualFold(strings.TrimSpace(plugin), "calico")
	if isCalico {
		return defaultCalicoManifest
	}

	return defaultWeaveNetManifest
}

// GenerateCNIScript generates bash script for applying CNI network overlay.
func GenerateCNIScript(plugin string) string {
	manifestURL := resolveCNIManifest(plugin)

	return fmt.Sprintf(cniScriptTemplate, manifestURL)
}

func resolveJoinCmd(joinCmd string) string {
	trimmedCmd := strings.TrimSpace(joinCmd)
	hasSudo := strings.HasPrefix(trimmedCmd, "sudo ")
	if hasSudo {
		return trimmedCmd
	}

	return fmt.Sprintf("sudo %s", trimmedCmd)
}

// GenerateK8sJoinScript generates bash script for joining worker node to cluster.
func GenerateK8sJoinScript(joinCmd string) string {
	cmdWithSudo := resolveJoinCmd(joinCmd)

	return fmt.Sprintf(k8sJoinScriptTemplate, cmdWithSudo)
}

func resolveNFSExportDir(exportDir string) string {
	cleanDir := strings.TrimSpace(exportDir)
	hasDir := cleanDir != ""
	if hasDir {
		return cleanDir
	}

	return defaultNFSExportDir
}

// GenerateNFSServerScript generates bash script for installing and configuring NFS server.
func GenerateNFSServerScript(exportDir string) string {
	resolvedDir := resolveNFSExportDir(exportDir)

	return fmt.Sprintf(nfsServerScriptTemplate, resolvedDir, resolvedDir)
}

func resolveHelmVersion(version string) string {
	cleanVer := strings.TrimPrefix(strings.TrimSpace(version), "v")
	hasVer := cleanVer != ""
	if hasVer {
		return cleanVer
	}

	return defaultHelmVersion
}

// GenerateHelmInstallScript generates bash script for downloading and installing Helm CLI.
func GenerateHelmInstallScript(version string) string {
	resolvedVer := resolveHelmVersion(version)

	return fmt.Sprintf(helmInstallScriptTemplate, resolvedVer)
}

// GenerateHelmNFSScript generates bash script for configuring NFS CSI provisioner via Helm.
func GenerateHelmNFSScript(nfsServerIP string, exportDir string) string {
	resolvedDir := resolveNFSExportDir(exportDir)
	cleanIP := strings.TrimSpace(nfsServerIP)

	return fmt.Sprintf(helmNFSScriptTemplate, cleanIP, resolvedDir)
}

// GenerateK8sResetScript generates bash script for resetting Kubernetes cluster node state.
func GenerateK8sResetScript() string {
	return k8sResetScriptTemplate
}
