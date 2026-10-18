package cmdnodes

import "strings"

type nodesFilterOptions struct {
	isJSON       bool
	isFast       bool
	openOnly     bool
	isNoCache    bool
	isRefresh    bool
	filterSSH    bool
	filterClust  bool
	filterSC     bool
	targetFilter string
}

func parseNodesFilterOptions(args []string) nodesFilterOptions {
	var opts nodesFilterOptions
	for _, arg := range args {
		if arg == "--json" || arg == "-j" {
			opts.isJSON = true
			continue
		}
		if arg == "--fast" || arg == "--no-probe" || arg == "-f" {
			opts.isFast = true
			continue
		}
		if arg == "--open-only" || arg == "--running-only" || arg == "-ro" || arg == "--online-only" {
			opts.openOnly = true
			continue
		}
		if arg == "--no-cache" {
			opts.isNoCache = true
			continue
		}
		if arg == "--refresh" || arg == "-r" {
			opts.isRefresh = true
			InvalidateNodesLivenessCache()
			continue
		}
		if arg == "--ssh" {
			opts.filterSSH = true
			continue
		}
		if arg == "--cluster" || arg == "--clst" {
			opts.filterClust = true
			continue
		}
		if arg == "--sc" || arg == "--server-clients" || arg == "--servers-clients" {
			opts.filterSC = true
			continue
		}
		hasDashPrefix := strings.HasPrefix(arg, "-")
		isTargetCandidate := !hasDashPrefix && opts.targetFilter == ""
		if isTargetCandidate {
			opts.targetFilter = arg
		}
	}
	return opts
}
