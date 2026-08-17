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

	if hasInclude {
		selectorFiltered = filterIncluded(selectorFiltered, filter)
	}

	if hasExcept {
		selectorFiltered = filterExcluded(selectorFiltered, filter)
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
		isRangeMatch := hasSeparator && matchesRange(n, ex, trailingOctet, hasTrailing)
		if isRangeMatch {
			return true
		}
		isExactMatch := !hasSeparator && matchesExact(n, ex, trailingOctet, hasTrailing)
		if isExactMatch {
			return true
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
	if !isValidRange {
		return false
	}

	inDisplayIdRange := n.DisplayId >= start && n.DisplayId <= end
	if inDisplayIdRange {
		return true
	}

	if !hasTrailing {
		return false
	}

	tInt, err := strconv.Atoi(trailingOctet)
	isValidTrailing := err == nil
	if !isValidTrailing {
		return false
	}

	inTrailingRange := tInt >= start && tInt <= end
	if inTrailingRange {
		return true
	}

	return false
}

func matchesExact(n ClusterNode, ex string, trailingOctet string, hasTrailing bool) bool {
	exInt, err := strconv.Atoi(ex)
	isValidInt := err == nil

	isDotPrefix := strings.HasPrefix(ex, constants.DotSeparator)
	expected := constants.DotSeparator + trailingOctet
	isMatchDotSuffix := isDotPrefix && hasTrailing && ex == expected

	if !isValidInt && isMatchDotSuffix {
		return true
	}
	if !isValidInt {
		return false
	}

	isMatchDisplayId := n.DisplayId == exInt
	if isMatchDisplayId {
		return true
	}

	if !hasTrailing {
		return false
	}

	tInt, err2 := strconv.Atoi(trailingOctet)
	isValidTrailing := err2 == nil
	if !isValidTrailing {
		return false
	}

	isMatchTrailing := tInt == exInt
	if isMatchTrailing {
		return true
	}

	return false
}

func filterIncluded(nodes []ClusterNode, filter NodeFilter) []ClusterNode {
	var included []ClusterNode
	for _, n := range nodes {
		if isNodeIncluded(n, filter) {
			included = append(included, n)
		}
	}
	return included
}

func filterExcluded(nodes []ClusterNode, filter NodeFilter) []ClusterNode {
	var excluded []ClusterNode
	for _, n := range nodes {
		if !matchesExclusion(n, filter.Except) {
			excluded = append(excluded, n)
		}
	}
	return excluded
}
