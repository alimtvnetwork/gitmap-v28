# Kubernetes Cheat Sheet

Quick-reference commands for managing your cluster.

## Cluster Info

```bash
kubectl cluster-info                    # Cluster endpoint and services
kubectl get nodes                       # List all nodes
kubectl get nodes -o wide               # Nodes with IPs and OS info
kubectl version --short                 # Client + server versions
```

## Pods

```bash
kubectl get pods                        # List pods in default namespace
kubectl get pods -A                     # List pods in ALL namespaces
kubectl get pods -o wide                # Pods with node placement
kubectl describe pod <name>             # Detailed pod info
kubectl logs <pod-name>                 # Pod logs
kubectl logs <pod-name> -f              # Follow (tail) logs
kubectl exec -it <pod-name> -- bash     # Shell into a pod
kubectl delete pod <name>               # Delete a pod
```

## Deployments

```bash
kubectl get deployments                 # List deployments
kubectl create deployment nginx --image=nginx   # Quick deploy
kubectl scale deployment nginx --replicas=3     # Scale up/down
kubectl rollout status deployment nginx         # Rollout status
kubectl rollout undo deployment nginx           # Rollback
kubectl delete deployment nginx                 # Remove
```

## Services

```bash
kubectl get svc                         # List services
kubectl expose deployment nginx --port=80 --type=NodePort  # Expose
kubectl describe svc <name>             # Service details
kubectl delete svc <name>               # Remove service
```

## Namespaces

```bash
kubectl get namespaces                  # List namespaces
kubectl create namespace dev            # Create namespace
kubectl get pods -n kube-system         # Pods in specific namespace
kubectl config set-context --current --namespace=dev  # Switch default
```

## Config & Context

```bash
kubectl config view                     # Show kubeconfig
kubectl config current-context          # Current context
kubectl config get-contexts             # List all contexts
kubectl config use-context <name>       # Switch context
```

## Troubleshooting

```bash
kubectl describe node <name>            # Node events and conditions
kubectl get events --sort-by=.metadata.creationTimestamp  # Recent events
kubectl top nodes                       # Node resource usage (metrics-server required)
kubectl top pods                        # Pod resource usage
journalctl -xeu kubelet                 # Kubelet logs (on the node)
```

## kubeadm Commands

```bash
kubeadm init                            # Initialize control plane
kubeadm join <ip>:6443 --token <t> --discovery-token-ca-cert-hash <h>
kubeadm token create --print-join-command   # Regenerate join command
kubeadm token list                      # List active tokens
kubeadm reset --force                   # Tear down node (start fresh)
```

## Helm Commands

```bash
helm repo add <name> <url>              # Add chart repo
helm repo update                        # Refresh repos
helm search repo <keyword>              # Search charts
helm install <release> <chart>          # Install a chart
helm upgrade <release> <chart>          # Upgrade release
helm uninstall <release>                # Remove release
helm list                               # List installed releases
helm status <release>                   # Release info
```

## Useful Aliases

Add to `~/.bashrc` or `~/.zshrc`:

```bash
alias k="kubectl"
alias kgp="kubectl get pods"
alias kgn="kubectl get nodes"
alias kgs="kubectl get svc"
alias kga="kubectl get all"
alias kd="kubectl describe"
alias kl="kubectl logs"
alias kx="kubectl exec -it"
alias kns="kubectl config set-context --current --namespace"
```

## Network Plugins

```bash
# Weave Net (lightweight, easy setup)
kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml

# Calico (production-grade, network policies)
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml
```
