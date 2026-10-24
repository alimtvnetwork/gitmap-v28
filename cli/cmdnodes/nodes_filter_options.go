package cmdnodes

import (
	"slices"
	"strings"
)

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
		if slices.Contains([]string{"--json", "-j"}, arg) {
			opts.isJSON = true
			continue
		}
		if slices.Contains([]string{"--fast", "--no-probe", "-f"}, arg) {
			opts.isFast = true
			continue
		}
		if slices.Contains([]string{"--open-only", "--running-only", "-ro", "--online-only"}, arg) {
			opts.openOnly = true
			continue
		}
		if arg == "--no-cache" {
			opts.isNoCache = true
			continue
		}
		if slices.Contains([]string{"--refresh", "-r"}, arg) {
			opts.isRefresh = true
			InvalidateNodesLivenessCache()
			continue
		}
		if arg == "--ssh" {
			opts.filterSSH = true
			continue
		}
		if slices.Contains([]string{"--cluster", "--clst"}, arg) {
			opts.filterClust = true
			continue
		}
		if slices.Contains([]string{"--sc", "--server-clients", "--servers-clients"}, arg) {
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
