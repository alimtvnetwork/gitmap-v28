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
// arguments before handing them to the flag parser. We treat anything starting
// with `-` as a flag. If it matches a known flag that takes a value, we keep the
// next token as well.
func splitClusterFlagsAndArgs(args []string) (flags, positional []string) {
	expectValue := false
	for _, a := range args {
		if expectValue {
			flags = append(flags, a)
			expectValue = false
			continue
		}
		if isClusterFlag(a) {
			flags = append(flags, a)
			expectValue = needsClusterValue(a)
			continue
		}
		positional = append(positional, a)
	}
	return
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

// ParseClusterFlags parses flags for cluster commands and returns remaining positional arguments.
func ParseClusterFlags(args []string) (ClusterFlags, []string, error) {
	var opts ClusterFlags

	fs := flag.NewFlagSet("cluster", flag.ContinueOnError)

	var ips stringSliceFlag
	var ids intSliceFlag

	fs.StringVar(&opts.ExceptClause, constants.ClusterFlagExcept, "", "")
	fs.Var(&ips, constants.ClusterFlagIP, "")
	fs.Var(&ids, constants.ClusterFlagID, "")

	yes := fs.Bool(constants.ClusterFlagYes, false, "")
	yesShort := fs.Bool(constants.ClusterFlagYesShort, false, "")

	fs.BoolVar(&opts.ForceLifecycle, constants.ClusterFlagForceLifecycle, false, "")
	fs.BoolVar(&opts.NoPreflight, constants.ClusterFlagNoPreflight, false, "")

	flags, positional := splitClusterFlagsAndArgs(args)
	if err := fs.Parse(flags); err != nil {
		return opts, nil, errors.New("cluster: " + err.Error())
	}

	opts.AutoConfirm = *yes || *yesShort
	opts.OnlyIPs = ips
	opts.OnlyIDs = ids

	isExceptProvided := opts.ExceptClause != ""
	isIncludeProvided := len(opts.OnlyIPs) > constants.EmptySliceLength || len(opts.OnlyIDs) > constants.EmptySliceLength

	if isExceptProvided == true && isIncludeProvided == true {
		return opts, nil, errors.New(constants.ErrFilterExclusive)
	}

	return opts, positional, nil
}
