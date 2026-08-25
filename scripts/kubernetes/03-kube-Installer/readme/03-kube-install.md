# Kube Install

## Set IP and make sure network is bridge mode

Do it in imperative mode. Written instruction are declarative instructions.

### Imperative mode bridge mode

```bash
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
sudo modprobe overlay
sudo modprobe br_netfilter
# sysctl params required by setup, params persist across reboots
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
```

### Apply sysctl params without reboot

`sudo sysctl --system`

#### Verify that the br_netfilter, overlay modules are loaded by running the following commands:

```bash
lsmod | grep br_netfilter
lsmod | grep overlay
```

#### Verify that the net.bridge.bridge-nf-call-iptables, net.bridge.bridge-nf-call-ip6tables, and net.ipv4.ip_forward system variables are set to 1 in your sysctl config by running the following command Must run (imperative mode bridge mode)[#Imperative mode bridge mode]

```bash
sysctl net.bridge.bridge-nf-call-iptables net.bridge.bridge-nf-call-ip6tables net.ipv4.ip_forward
```

#### Must run (imperative mode bridge mode)[#Imperative mode bridge mode]

```bash
sudo sysctl -w net.ipv4.ip_forward=1
```

### Set IP

Set IP using /02-ubuntu-install/00-set-ip.sh

```bash
sudo /02-ubuntu-install/00-set-ip.sh (your Ip) [your default gateway]
sudo /02-ubuntu-install/00-set-ip.sh 192.168.0.85
```

#### This is the network config written by 'subiquity'

```yaml
network:
  renderer: networkd
  ethernets:
    ens33:
      dhcp4: false
      addresses:
        - 192.168.0.66/24
      routes:
        - to: default
          via: 192.168.0.1
      nameservers:
          addresses: [8.8.8.8, 1.1.1.1]
  version: 2
```

Now validate and apply

```bash
sudo netplan try && sudo netplan apply 
```

## Add Repo, Kube and CRIO

### Kubernetes

```bash
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /" | tee /etc/apt/sources.list.d/kubernetes.list
```

### CRIO

```bash
curl -fsSL https://pkgs.k8s.io/addons:/cri-o:/prerelease:/main/deb/Release.key |
    gpg --dearmor -o /etc/apt/keyrings/cri-o-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/cri-o-apt-keyring.gpg] https://pkgs.k8s.io/addons:/cri-o:/prerelease:/main/deb/ /" |
    tee /etc/apt/sources.list.d/cri-o.list
```

## OS Updates

First get into `sudo user` by running

```bash
sudo su
```

```bash
apt update -y &&\
apt-get update -y &&\
apt upgrade -y 
```

### OS Update with Faster

```bash
sudo add-apt-repository ppa:apt-fast/stable &&
sudo apt-get update -y &&
sudo apt-get install apt-fast -y &&
sudo apt-get install -y aria2
```

## Basic Installations

First get into `sudo user` by running

```bash
sudo su
```

Then 


```bash
add-apt-repository ppa:apt-fast/stable &&
add-apt-repository ppa:git-core/ppa &&
apt update -y &&\
apt-get update -y &&\
apt upgrade -y &&\
apt-get install -y vim build-essential wget nano curl file git libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev git-core && \
apt install -y zsh git-lfs && \
git version && \
git lfs --version
```

## Oh-my-Zsh

```bash
wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh && \
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
```

### Fast downloader

```bash
sudo apt-get install apt-fast -y &&
sudo apt-get install -y aria2
```

## Swap Off

From SU user

`sudo nano /etc/fstab`

Comment out the swap section by `#`

Then Reboot 

`sudo reboot`

Then check the status of swap, it should be empty or nothing.

`sudo swapon --show`


## Kube Pre-installations

```bash
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
```

## Kube Installation with Kubeadm

```bash
apt update -y &&\
apt-get update -y &&\
apt upgrade -y
apt-get install -y apt-transport-https ca-certificates curl gpg && \
apt update -y &&\
apt-get update -y &&\
apt upgrade -y && \
apt-get install -y cri-o kubelet kubeadm kubectl socat
```

## Fix Network and Hostnamectl

```bash
 kubeadm reset
 hostnamectl set-hostname master # master node only
 hostnamectl set-hostname worker1 # worker node 1 node only
 hostnamectl set-hostname worker2 # worker node 2 node only
```

After install in the control plane machine run

```bash
systemctl enable crio.service && \
systemctl start crio.service && \
systemctl status crio.service && \
kubeadm config images pull
```

```bash
systemctl status crio.service
```

```bash
sudo sysctl -w net.ipv4.ip_forward=1
```


if not ran before.

### kubeadm init (only from master)

Make sure every machine has different host-names

```bash
kubeadm init
```
If any issues occurred, reset the network

```bash
kubeadm reset
```

### Don't run the apiserver part

```bash
kubeadm init --apiserver-advertise-address=192.168.0.66 --pod-network-cidr=192.168.0.0/16 --ignore-preflight-errors=all
```

## CRIO Services (Before Joining run this on master)

```bash
systemctl enable crio.service && \
systemctl start crio.service && \
systemctl status crio.service
```

## kubectl Services (Before Joining run this on master)

```bash
systemctl enable kubeapiserver && \
systemctl start crio.service && \
systemctl status crio.service
```

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

No need to run this

```bash
export KUBECONFIG=/etc/kubernetes/admin.conf
```

## After Running (init on master) - Do join on Worker - then Apply Yaml on Master node

```bash
kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml
kubectl -n kube-system get pods
```

### Deploy Calico network (no need)

```bash
kubectl --kubeconfig=/etc/kubernetes/admin.conf create -f https://docs.projectcalico.org/v3.14/manifests/calico.yaml
```

### Kube init

```bash
kubeadm init
```

### To get the join command

```bash
kubeadm token create --print-join-command
sudo kubeadm token create --print-join-command
```

### To reset network

```bash
kubeadm reset --force
```

## Install Quickly

```bash
chmod +x ./10-kube-init.sh && \
sudo ./10-kube-init.sh
```

```bash
mkdir -p $HOME/.kube && \
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
chown $(id -u):$(id -g) "$HOME/.kube/config"
```

```bash
kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml"
```

```bash
kubectl -n kube-system get pods
```