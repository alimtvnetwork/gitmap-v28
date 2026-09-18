# Subtask 03: Kubernetes Help Documentation & Catalog Integration

## Objective
Author comprehensive user guides and reference manuals for the Kubernetes Cluster Lifecycle suite in `cli/helptext/cluster.md` and register new topics in `cli/helptext/catalog.go`.

## Disjoint Files Assigned
- `cli/helptext/cluster.md`
- `cli/helptext/catalog.go`

## Implementation Details
1. `cli/helptext/cluster.md`:
   - Document complete lifecycle with ASCII topology diagram.
   - Document all `gitmap cluster k8s` subcommands (`prereq`, `install`, `init`, `cni`, `join`, `nfs`, `helm-install`, `helm-nfs`, `status`, `reset`).
   - Add concrete copy-pasteable examples for single-master multi-worker clusters.
   - Document default ports, kernel modules, CRI-O service management, and Weave Net / Calico overlay networks.

2. `cli/helptext/catalog.go`:
   - Register topics:
     * `cluster-k8s`: "Manage Kubernetes cluster lifecycle, runtime, and components"
     * `cluster-k8s-init`: "Initialize Kubernetes control plane"
     * `cluster-k8s-join`: "Join worker nodes to Kubernetes control plane"
     * `cluster-k8s-helm`: "Deploy Helm and storage class provisioners"

## Constraints
- Strict Unix LF line endings.
- GitHub flavored markdown formatting.
