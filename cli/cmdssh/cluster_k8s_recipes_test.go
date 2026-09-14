package cmdssh

import (
	"testing"
)

func TestGenerateK8sPrereqScript_Modules(t *testing.T) {
	script := GenerateK8sPrereqScript()
	assertContains(t, script, "/etc/modules-load.d/k8s.conf")
	assertContains(t, script, "overlay")
	assertContains(t, script, "br_netfilter")
	assertContains(t, script, "modprobe overlay")
	assertContains(t, script, "modprobe br_netfilter")
}

func TestGenerateK8sPrereqScript_Sysctl(t *testing.T) {
	script := GenerateK8sPrereqScript()
	assertContains(t, script, "/etc/sysctl.d/k8s.conf")
	assertContains(t, script, "net.bridge.bridge-nf-call-iptables = 1")
	assertContains(t, script, "net.bridge.bridge-nf-call-ip6tables = 1")
	assertContains(t, script, "net.ipv4.ip_forward = 1")
	assertContains(t, script, "sysctl --system")
}

func TestGenerateK8sInstallScript_Defaults(t *testing.T) {
	script := GenerateK8sInstallScript("", "")
	assertContains(t, script, "core:/stable:/v1.31/deb/")
	assertContains(t, script, "addons:/cri-o:/prerelease:/main/deb/")
	assertContains(t, script, "cri-o kubelet kubeadm kubectl socat")
	assertContains(t, script, "systemctl start crio.service")
	assertContains(t, script, "kubeadm config images pull")
	assertNotContains(t, script, "hostnamectl")
}

func TestGenerateK8sInstallScript_Custom(t *testing.T) {
	script := GenerateK8sInstallScript("1.32", "worker-1")
	assertContains(t, script, "hostnamectl set-hostname \"worker-1\"")
	assertContains(t, script, "core:/stable:/v1.32/deb/")
}

func TestGenerateK8sInstallScript_Prefix(t *testing.T) {
	script := GenerateK8sInstallScript("v1.30", "")
	assertContains(t, script, "core:/stable:/v1.30/deb/")
	assertNotContains(t, script, "core:/stable:/vv1.30")
}

func TestGenerateK8sInitScript_Default(t *testing.T) {
	script := GenerateK8sInitScript("")
	assertContains(t, script, "systemctl start crio.service")
	assertContains(t, script, "kubeadm init\n")
	assertContains(t, script, "$HOME/.kube/config")
	assertContains(t, script, "chown $(id -u):$(id -g) $HOME/.kube/config")
	assertNotContains(t, script, "--pod-network-cidr")
}

func TestGenerateK8sInitScript_WithCIDR(t *testing.T) {
	script := GenerateK8sInitScript("10.244.0.0/16")
	assertContains(t, script, "kubeadm init --pod-network-cidr=10.244.0.0/16")
	assertContains(t, script, "$HOME/.kube/config")
}

func TestGenerateCNIScript_Weave(t *testing.T) {
	scriptEmpty := GenerateCNIScript("")
	assertContains(t, scriptEmpty, "https://reweave.azurewebsites.net/k8s/v1.31/net.yaml")
	scriptWeave := GenerateCNIScript("weave")
	assertContains(t, scriptWeave, "https://reweave.azurewebsites.net/k8s/v1.31/net.yaml")
}

func TestGenerateCNIScript_Calico(t *testing.T) {
	script := GenerateCNIScript("calico")
	assertContains(t, script, "https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml")
	scriptUpper := GenerateCNIScript("CALICO")
	assertContains(t, scriptUpper, "calico.yaml")
}

func TestGenerateK8sJoinScript_WithoutSudo(t *testing.T) {
	script := GenerateK8sJoinScript("kubeadm join 10.0.0.1:6443 --token abc")
	assertContains(t, script, "sudo kubeadm join 10.0.0.1:6443 --token abc")
}

func TestGenerateK8sJoinScript_WithSudo(t *testing.T) {
	script := GenerateK8sJoinScript("sudo kubeadm join 10.0.0.1:6443 --token abc")
	assertContains(t, script, "sudo kubeadm join 10.0.0.1:6443 --token abc")
	assertNotContains(t, script, "sudo sudo")
}

func TestGenerateNFSServerScript_Default(t *testing.T) {
	script := GenerateNFSServerScript("")
	assertContains(t, script, "apt-get install -y nfs-server")
	assertContains(t, script, "mkdir -p /nfsexport")
	assertContains(t, script, "/nfsexport *(rw,no_root_squash)")
	assertContains(t, script, "systemctl restart nfs-server")
}

func TestGenerateNFSServerScript_Custom(t *testing.T) {
	script := GenerateNFSServerScript("/mnt/storage")
	assertContains(t, script, "mkdir -p /mnt/storage")
	assertContains(t, script, "/mnt/storage *(rw,no_root_squash)")
}

func TestGenerateHelmInstallScript_Default(t *testing.T) {
	script := GenerateHelmInstallScript("")
	assertContains(t, script, "helm-v3.16.2-linux-amd64.tar.gz")
	assertContains(t, script, "/usr/local/bin/helm")
	assertContains(t, script, "rm -rf")
}

func TestGenerateHelmInstallScript_Custom(t *testing.T) {
	script := GenerateHelmInstallScript("3.17.0")
	assertContains(t, script, "helm-v3.17.0-linux-amd64.tar.gz")
	scriptPrefixed := GenerateHelmInstallScript("v3.17.0")
	assertContains(t, scriptPrefixed, "helm-v3.17.0-linux-amd64.tar.gz")
	assertNotContains(t, scriptPrefixed, "helm-vv3.17.0")
}

func TestGenerateHelmNFSScript_Default(t *testing.T) {
	script := GenerateHelmNFSScript("192.168.0.20", "")
	assertContains(t, script, "helm repo add nfs-subdir-external-provisioner")
	assertContains(t, script, "helm repo update")
	assertContains(t, script, "helm uninstall nfs-provisioner")
	assertContains(t, script, "--set nfs.server=\"192.168.0.20\" --set nfs.path=\"/nfsexport\"")
}

func TestGenerateHelmNFSScript_Custom(t *testing.T) {
	script := GenerateHelmNFSScript("10.0.0.50", "/data/nfs")
	assertContains(t, script, "--set nfs.server=\"10.0.0.50\" --set nfs.path=\"/data/nfs\"")
}

func TestGenerateK8sResetScript(t *testing.T) {
	script := GenerateK8sResetScript()
	assertContains(t, script, "kubeadm reset --force")
	assertContains(t, script, "rm -rf $HOME/.kube /etc/cni/net.d")
}
