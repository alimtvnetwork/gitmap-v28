package cluster

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
