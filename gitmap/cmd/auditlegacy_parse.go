// Package cmd: flag parsing for `gitmap audit-legacy`.
package cmd

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// parseAuditLegacyArgs parses flags into an options struct.
func parseAuditLegacyArgs(args []string) (auditLegacyOpts, error) {
	fs := flag.NewFlagSet(constants.CmdAuditLegacy, flag.ContinueOnError)
	pats := fs.String(constants.FlagAuditLegacyPatterns,
		constants.DefaultAuditLegacyPatterns, constants.FlagDescAuditLegacyPatterns)
	root := fs.String(constants.FlagAuditLegacyPath, ".", constants.FlagDescAuditLegacyPath)
	asJSON := fs.Bool(constants.FlagAuditLegacyJSON, false, constants.FlagDescAuditLegacyJSON)
	report := fs.String(constants.FlagAuditLegacyReport, "",
		constants.FlagDescAuditLegacyReport)
	diffs := fs.Bool(constants.FlagAuditLegacyDiffs, false,
		constants.FlagDescAuditLegacyDiffs)
	if err := fs.Parse(args); err != nil {
		return auditLegacyOpts{}, err
	}
	reportSet := isAuditFlagSet(fs, constants.FlagAuditLegacyReport)
	compiled, raw, err := compileAuditPatterns(*pats)
	if err != nil {
		return auditLegacyOpts{}, err
	}

	return auditLegacyOpts{
		Patterns: compiled, Raw: raw, Root: *root,
		AsJSON: *asJSON, ReportPath: resolveReportPath(reportSet, *report),
		WriteDiffs: *diffs,
	}, nil
}

// isAuditFlagSet returns true when the user explicitly passed `name`.
func isAuditFlagSet(fs *flag.FlagSet, name string) bool {
	seen := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			seen = true
		}
	})

	return seen
}

// resolveReportPath returns "" when --report wasn't passed, the
// default file when passed without a value, or the user's value.
func resolveReportPath(set bool, value string) string {
	if !set {
		return ""
	}
	if value == "" {
		return constants.DefaultAuditLegacyReport
	}

	return value
}

// compileAuditPatterns compiles a comma-separated pattern list.
func compileAuditPatterns(csv string) ([]*regexp.Regexp, []string, error) {
	parts := strings.Split(csv, ",")
	out := make([]*regexp.Regexp, 0, len(parts))
	raw := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		re, err := compileOnePattern(p)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, re)
		raw = append(raw, p)
	}

	return out, raw, nil
}

// compileOnePattern wraps regexp.Compile with a standardized error.
func compileOnePattern(p string) (*regexp.Regexp, error) {
	re, err := regexp.Compile(p)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrAuditLegacyRegex, p, err)
	}

	return re, nil
}
