package cluster

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestResolveTargetNodes_MutualExclusivity(t *testing.T) {
	allNodes := []ClusterNode{
		{ID: "node1", DisplayId: 1, IP: "192.168.1.1", IsServer: false},
	}

	filter := NodeFilter{
		Except: []string{"192.168.1.2"},
		IPs:    []string{"192.168.1.1"},
	}

	_, err := ResolveTargetNodes(ServersClients, filter, allNodes)
	if err == nil || err.Error() != constants.ErrFilterExclusive {
		t.Errorf("Expected %v, got %v", constants.ErrFilterExclusive, err)
	}

	filter2 := NodeFilter{
		Except: []string{"192.168.1.2"},
		IDs:    []int{1},
	}

	_, err = ResolveTargetNodes(ServersClients, filter2, allNodes)
	if err == nil || err.Error() != constants.ErrFilterExclusive {
		t.Errorf("Expected %v, got %v", constants.ErrFilterExclusive, err)
	}
}

func TestResolveTargetNodes_Exclusions(t *testing.T) {
	allNodes := []ClusterNode{
		{ID: "node1", DisplayId: 1, IP: "192.168.1.10", IsServer: false},
		{ID: "node2", DisplayId: 2, IP: "192.168.1.11", IsServer: false},
		{ID: "node3", DisplayId: 3, IP: "192.168.1.12", IsServer: false},
		{ID: "node4", DisplayId: 4, IP: "192.168.1.13", IsServer: false},
		{ID: "node5", DisplayId: 5, IP: "192.168.2.14", IsServer: false},
	}

	tests := []struct {
		name     string
		except   []string
		expected []int // DisplayIds expected to remain
	}{
		{
			name:     "integer exact match",
			except:   []string{"2"},
			expected: []int{1, 3, 4, 5},
		},
		{
			name:     "IP exact match",
			except:   []string{"192.168.1.11"},
			expected: []int{1, 3, 4, 5},
		},
		{
			name:     "partial octet exact match",
			except:   []string{".12"},
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "ID range match",
			except:   []string{"2-4"},
			expected: []int{1, 5},
		},
		{
			name:     "IP octet range match",
			except:   []string{"10-12"},
			expected: []int{4, 5},
		},
		{
			name:     "mixed exclusions",
			except:   []string{"1", "192.168.1.13", ".14"},
			expected: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NodeFilter{Except: tt.except}
			result, err := ResolveTargetNodes(ServersClients, filter, allNodes)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d nodes, got %d", len(tt.expected), len(result))
			}

			expectedMap := make(map[int]bool)
			for _, id := range tt.expected {
				expectedMap[id] = true
			}

			for _, n := range result {
				if !expectedMap[n.DisplayId] {
					t.Errorf("did not expect to find DisplayId %d in results", n.DisplayId)
				}
			}
		})
	}
}

func TestResolveTargetNodes_Selector(t *testing.T) {
	allNodes := []ClusterNode{
		{ID: "server", DisplayId: 1, IP: "10.0.0.1", IsServer: true},
		{ID: "client", DisplayId: 2, IP: "10.0.0.2", IsServer: false},
	}

	res, err := ResolveTargetNodes(ClientsOnly, NodeFilter{}, allNodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 || res[0].DisplayId != 2 {
		t.Errorf("expected only client, got %v", res)
	}
}
