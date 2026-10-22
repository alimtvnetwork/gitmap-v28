package store

import (
	"context"
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func appendUniqueHost(hosts []SSHHost, seen map[string]bool, h SSHHost) []SSHHost {
	_, isSeen := seen[h.ID]
	if isSeen {
		return hosts
	}

	seen[h.ID] = true

	return append(hosts, h)
}

func mergeUniqueHosts(resultList []SSHHost, seen map[string]bool, additions []SSHHost) []SSHHost {
	for _, host := range additions {
		resultList = appendUniqueHost(resultList, seen, host)
	}

	return resultList
}

func resolveTargetPart(ctx context.Context, part string, db *sql.DB) ([]SSHHost, error) {
	trimmed := strings.TrimSpace(part)
	if trimmed == "" {
		return nil, nil
	}

	return ListHostsByTarget(ctx, trimmed, db)
}

func resolveCommaTargets(ctx context.Context, target string, db *sql.DB) ([]SSHHost, error) {
	seen := make(map[string]bool)
	var out []SSHHost
	for _, part := range strings.Split(target, ",") {
		sub, err := resolveTargetPart(ctx, part, db)
		if err != nil {
			return nil, err
		}

		out = mergeUniqueHosts(out, seen, sub)
	}

	return ensureHostsSlice(out), nil
}

func resolveRoleTarget(ctx context.Context, target string, db *sql.DB) ([]SSHHost, bool, error) {
	if target == "control" || target == "master" {
		hosts, err := ListHostsByRole(ctx, "control", db)
		return hosts, true, err
	}

	if target == "workers" || target == "worker" || target == "nodes" {
		hosts, err := ListHostsByRole(ctx, "worker", db)
		return hosts, true, err
	}

	return nil, false, nil
}

func resolveGroupTarget(ctx context.Context, target string, db *sql.DB) ([]SSHHost, bool, error) {
	if target == "all" || target == "" {
		hosts, err := ListHosts(ctx, db)
		return hosts, true, err
	}

	return resolveRoleTarget(ctx, target, db)
}

func findHostByAliasOrIP(ctx context.Context, target string, db *sql.DB) ([]SSHHost, bool, error) {
	host, err := GetHostByAlias(ctx, target, db)
	if err == nil {
		return []SSHHost{host}, true, nil
	}

	host, err = GetHostByIP(ctx, target, db)
	if err == nil {
		return []SSHHost{host}, true, nil
	}

	return nil, false, nil
}

func wrapHostNotFoundError(target string) *apperror.AppError {
	appErr := apperror.Wrap(apperror.ErrNotFound, "ListHostsByTarget", map[string]any{"target": target})
	appErr.Code = "E_INTERNAL_ERROR"

	return appErr
}

func resolveCommaWrap(ctx context.Context, target string, db *sql.DB) SSHHostSliceWrap {
	hosts, err := resolveCommaTargets(ctx, target, db)
	if err != nil {
		return result.FailSlice[SSHHost](err)
	}

	return result.OkSlice(hosts)
}

func resolveGroupWrap(ctx context.Context, target string, db *sql.DB) (SSHHostSliceWrap, bool) {
	hosts, isGroup, err := resolveGroupTarget(ctx, target, db)
	if !isGroup {
		return SSHHostSliceWrap{}, false
	}

	if err != nil {
		return result.FailSlice[SSHHost](err), true
	}

	return result.OkSlice(hosts), true
}

func resolveSingleWrap(ctx context.Context, target string, db *sql.DB) (SSHHostSliceWrap, bool) {
	hosts, isSingle, err := findHostByAliasOrIP(ctx, target, db)
	if !isSingle {
		return SSHHostSliceWrap{}, false
	}

	if err != nil {
		return result.FailSlice[SSHHost](err), true
	}

	return result.OkSlice(hosts), true
}

// ListHostsByTargetWrap resolves SSH hosts by target specifier encapsulated in SSHHostSliceWrap.
func ListHostsByTargetWrap(ctx context.Context, target string, db *sql.DB) SSHHostSliceWrap {
	ctx = safeContext(ctx)
	if strings.Contains(target, ",") {
		return resolveCommaWrap(ctx, target, db)
	}

	groupRes, hasGroup := resolveGroupWrap(ctx, target, db)
	if hasGroup {
		return groupRes
	}

	singleRes, hasSingle := resolveSingleWrap(ctx, target, db)
	if hasSingle {
		return singleRes
	}

	return result.FailSlice[SSHHost](wrapHostNotFoundError(target))
}

// ListHostsByTarget resolves SSH hosts by target specifier (role, alias, IP, or comma list).
func ListHostsByTarget(ctx context.Context, target string, db *sql.DB) ([]SSHHost, error) {
	res := ListHostsByTargetWrap(ctx, target, db)
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

// GetSSHHostByAlias retrieves an SSHHost pointer by alias.
func GetSSHHostByAlias(ctx context.Context, alias string, db *sql.DB) (*SSHHost, error) {
	host, err := GetHostByAlias(ctx, alias, db)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

// GetSSHHostByIP retrieves an SSHHost pointer by IP address.
func GetSSHHostByIP(ctx context.Context, ip string, db *sql.DB) (*SSHHost, error) {
	host, err := GetHostByIP(ctx, ip, db)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

// ListSSHHostsWrap retrieves all SSH hosts encapsulated in SSHHostSliceWrap.
func ListSSHHostsWrap(ctx context.Context, db *sql.DB) SSHHostSliceWrap {
	return ListHostsWrap(ctx, db)
}

// ListSSHHosts retrieves all SSH hosts from the database.
func ListSSHHosts(ctx context.Context, db *sql.DB) ([]SSHHost, error) {
	res := ListSSHHostsWrap(ctx, db)
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}
