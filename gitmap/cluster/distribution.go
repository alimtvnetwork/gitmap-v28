package cluster

import "time"

const (
	// EmptySize represents a length or size of zero.
	EmptySize = 0
)

// Client defines a node in the cluster.
type Client struct {
	ID string
}

// Repo defines a repository to be managed.
type Repo struct {
	URL string
}

// Workload defines the set of repositories assigned to a specific client.
type Workload struct {
	ClientID string
	Repos    []Repo
}

// DistributeWorkload splits a list of repositories as evenly as possible among available clients.
func DistributeWorkload(clients []Client, repos []Repo) []Workload {
	hasNoClients := len(clients) == EmptySize
	if hasNoClients == true {
		return nil
	}

	workloads := make([]Workload, len(clients))
	for i := range clients {
		workloads[i] = Workload{
			ClientID: clients[i].ID,
			Repos:    []Repo{},
		}
	}

	for i := range repos {
		clientIndex := i % len(clients)
		workloads[clientIndex].Repos = append(workloads[clientIndex].Repos, repos[i])
	}

	return workloads
}

// DistributionLoop monitors node health and redistributes workloads if a node drops.
// It calls onUpdate whenever a new workload distribution is generated.
func DistributionLoop(registry *Registry, repos []Repo, interval time.Duration, onUpdate func([]Workload)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	activeClients := getActiveClients(registry)
	currentWorkloads := DistributeWorkload(activeClients, repos)
	onUpdate(currentWorkloads)

	for range ticker.C {
		registry.CheckHeartbeats()

		newActiveClients := getActiveClients(registry)

		hasNodeChanged := clientsChanged(activeClients, newActiveClients)

		if hasNodeChanged == true {
			activeClients = newActiveClients
			currentWorkloads = DistributeWorkload(activeClients, repos)
			onUpdate(currentWorkloads)
		}
	}
}

// getActiveClients returns a list of clients that are currently connected.
func getActiveClients(registry *Registry) []Client {
	nodes := registry.GetNodes()
	active := []Client{}
	for _, node := range nodes {
		isConnected := node.State == StateConnected
		if isConnected == true {
			active = append(active, Client{ID: node.ID})
		}
	}
	return active
}

func clientsChanged(old, new []Client) bool {
	if len(old) != len(new) {
		return true
	}
	for i := range old {
		if old[i].ID != new[i].ID {
			return true
		}
	}
	return false
}
