package cluster

import (
	"errors"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type TargetSelectorType int

const (
	ServersClients TargetSelectorType = 1
	ClientsOnly    TargetSelectorType = 2
)

// ClusterNode represents a cluster node for resolving purposes.
type ClusterNode struct {
	ID        string
	DisplayId int
	IP        string
	IsServer  bool
}

type NodeFilter struct {
	Except []string
	IPs    []string
	IDs    []int
}

// ResolveTargetNodes implements IP/ID inclusion and exclusion matching.
func ResolveTargetNodes(selector TargetSelectorType, filter NodeFilter, allNodes []ClusterNode) ([]ClusterNode, error) {
	hasNoNodes := len(allNodes) == constants.EmptySliceLength
	if hasNoNodes == true {
		return nil, errors.New(constants.ErrClusterNoNodes)
	}

	hasExcept := len(filter.Except) > constants.EmptySliceLength
	hasInclude := len(filter.IPs) > constants.EmptySliceLength || len(filter.IDs) > constants.EmptySliceLength

	isExclusiveError := hasExcept == true && hasInclude == true
	if isExclusiveError == true {
		return nil, errors.New(constants.ErrFilterExclusive)
	}

	var selectorFiltered []ClusterNode
	for _, n := range allNodes {
		isClientOnlyServer := selector == ClientsOnly && n.IsServer == true
		if isClientOnlyServer == true {
			continue
		}
		selectorFiltered = append(selectorFiltered, n)
	}

	if hasInclude == true {
		var included []ClusterNode
		for _, n := range selectorFiltered {
			isIncluded := isNodeIncluded(n, filter)
			if isIncluded == true {
				included = append(included, n)
			}
		}
		selectorFiltered = included
	}

	if hasExcept == true {
		var excluded []ClusterNode
		for _, n := range selectorFiltered {
			isExcluded := matchesExclusion(n, filter.Except)
			if isExcluded == false {
				excluded = append(excluded, n)
			}
		}
		selectorFiltered = excluded
	}

	return selectorFiltered, nil
}

func isNodeIncluded(n ClusterNode, filter NodeFilter) bool {
	for _, ip := range filter.IPs {
		isMatchIP := n.IP == ip
		if isMatchIP == true {
			return true
		}
	}
	for _, id := range filter.IDs {
		isMatchID := n.DisplayId == id
		if isMatchID == true {
			return true
		}
	}
	return false
}

func matchesExclusion(n ClusterNode, excepts []string) bool {
	ipParts := strings.Split(n.IP, constants.DotSeparator)
	trailingOctet := ""
	hasTrailing := len(ipParts) == constants.IPv4OctetCount
	if hasTrailing == true {
		trailingOctet = ipParts[constants.IPv4TrailingOctetIndex]
	}

	for _, ex := range excepts {
		isExactIPMatch := ex == n.IP
		if isExactIPMatch == true {
			return true
		}

		hasSeparator := strings.Contains(ex, constants.RangeSeparator)
		if hasSeparator == true {
			isRangeMatch := matchesRange(n, ex, trailingOctet, hasTrailing)
			if isRangeMatch == true {
				return true
			}
		} else {
			isExactMatch := matchesExact(n, ex, trailingOctet, hasTrailing)
			if isExactMatch == true {
				return true
			}
		}
	}
	return false
}

func matchesRange(n ClusterNode, ex string, trailingOctet string, hasTrailing bool) bool {
	parts := strings.Split(ex, constants.RangeSeparator)
	isTwoParts := len(parts) == constants.RangePartsCount
	if isTwoParts == false {
		return false
	}

	start, err1 := strconv.Atoi(parts[constants.RangeStartPartIndex])
	end, err2 := strconv.Atoi(parts[constants.RangeEndPartIndex])
	isValidRange := err1 == nil && err2 == nil
	if isValidRange == true {
		inDisplayIdRange := n.DisplayId >= start && n.DisplayId <= end
		if inDisplayIdRange == true {
			return true
		}

		if hasTrailing == true {
			tInt, err := strconv.Atoi(trailingOctet)
			isValidTrailing := err == nil
			if isValidTrailing == true {
				inTrailingRange := tInt >= start && tInt <= end
				if inTrailingRange == true {
					return true
				}
			}
		}
	}
	return false
}

func matchesExact(n ClusterNode, ex string, trailingOctet string, hasTrailing bool) bool {
	exInt, err := strconv.Atoi(ex)
	isValidInt := err == nil
	if isValidInt == true {
		isMatchDisplayId := n.DisplayId == exInt
		if isMatchDisplayId == true {
			return true
		}

		if hasTrailing == true {
			tInt, err2 := strconv.Atoi(trailingOctet)
			isValidTrailing := err2 == nil
			if isValidTrailing == true {
				isMatchTrailing := tInt == exInt
				if isMatchTrailing == true {
					return true
				}
			}
		}
	} else {
		isDotPrefix := strings.HasPrefix(ex, constants.DotSeparator)
		if isDotPrefix == true && hasTrailing == true {
			expected := constants.DotSeparator + trailingOctet
			isMatchDotSuffix := ex == expected
			if isMatchDotSuffix == true {
				return true
			}
		}
	}
	return false
}
