package cmd

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ClusterFlags holds the parsed flags for cluster delegation commands.
type ClusterFlags struct {
	Selector       cluster.TargetSelectorType
	ExceptClause   string
	OnlyIPs        []string
	OnlyIDs        []int
	AutoConfirm    bool
	ForceLifecycle bool
	NoPreflight    bool
	Verbose        bool
	DryRun         bool
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}

type intSliceFlag []int

func (s *intSliceFlag) String() string {
	strs := make([]string, len(*s))
	for i, v := range *s {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, ",")
}

func (s *intSliceFlag) Set(value string) error {
	for _, p := range strings.Split(value, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return err
		}
		*s = append(*s, v)
	}
	return nil
}

// splitClusterFlagsAndArgs is a helper to properly separate flags from positional
// arguments before handing them to the flag parser.
func splitClusterFlagsAndArgs(args []string) (flags, positional []string) {
	expectVal := false
	for _, a := range args {
		switch {
		case expectVal:
			flags, expectVal = append(flags, a), false
		case isClusterFlag(a):
			flags, expectVal = append(flags, a), needsClusterValue(a)
		default:
			positional = append(positional, a)
		}
	}
	return flags, positional
}

func isClusterFlag(s string) bool {
	return len(s) >= 2 && s[0] == '-'
}

func needsClusterValue(token string) bool {
	if strings.Contains(token, "=") {
		return false
	}
	name := strings.TrimLeft(token, "-")
	return name == constants.ClusterFlagExcept || name == constants.ClusterFlagIP || name == constants.ClusterFlagID
}

func bindClusterFlags(fs *flag.FlagSet, opts *ClusterFlags, ips *stringSliceFlag, ids *intSliceFlag) (*bool, *bool) {
	fs.StringVar(&opts.ExceptClause, constants.ClusterFlagExcept, "", "")
	fs.Var(ips, constants.ClusterFlagIP, "")
	fs.Var(ids, constants.ClusterFlagID, "")
	yes := fs.Bool(constants.ClusterFlagYes, false, "")
	yesShort := fs.Bool(constants.ClusterFlagYesShort, false, "")
	fs.BoolVar(&opts.ForceLifecycle, constants.ClusterFlagForceLifecycle, false, "")
	fs.BoolVar(&opts.NoPreflight, constants.ClusterFlagNoPreflight, false, "")
	fs.BoolVar(&opts.Verbose, constants.ClusterFlagVerbose, false, "")
	fs.BoolVar(&opts.DryRun, constants.ClusterFlagDryRun, false, "")
	return yes, yesShort
}

func validateClusterFilter(opts ClusterFlags) error {
	isExceptProvided := opts.ExceptClause != ""
	isIncludeProvided := len(opts.OnlyIPs) > constants.EmptySliceLength || len(opts.OnlyIDs) > constants.EmptySliceLength
	if isExceptProvided && isIncludeProvided {
		return errors.New(constants.ErrFilterExclusive)
	}
	return nil
}

func parseClusterFlagSet(args []string, opts *ClusterFlags) ([]string, error) {
	var ips stringSliceFlag
	var ids intSliceFlag
	fs := flag.NewFlagSet("cluster", flag.ContinueOnError)
	yes, yesShort := bindClusterFlags(fs, opts, &ips, &ids)
	flags, pos := splitClusterFlagsAndArgs(args)
	if err := fs.Parse(flags); err != nil {
		return nil, errors.New("cluster: " + err.Error())
	}
	opts.AutoConfirm, opts.OnlyIPs, opts.OnlyIDs = *yes || *yesShort, ips, ids
	return pos, validateClusterFilter(*opts)
}

// ParseClusterFlags parses flags for cluster commands and returns remaining positional arguments.
func ParseClusterFlags(args []string) (ClusterFlags, []string, error) {
	var opts ClusterFlags
	pos, err := parseClusterFlagSet(args, &opts)
	if err != nil {
		return opts, nil, err
	}
	return opts, pos, nil
}
