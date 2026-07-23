package render

// adapters.go — convenience builders that turn the gitmap data
// types most-likely-to-be-rendered (model.ScanRecord, etc.) into
// RepoTermBlock instances. Lives in the render package so callers
// don't have to know the field-name mapping; producers in
// cmd/clone-from/clone-next/probe import this helper to stay DRY.

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// FromScanRecord builds a RepoTermBlock from a model.ScanRecord.
// idx is the 1-based row number to surface to the user.
//
// CloneCommand defaults to the record's CloneInstruction, which
// the mapper has already shaped to "git clone [-b BRANCH] URL PATH".
// When the instruction is empty we synthesize a minimal command
// from the picked URL so the block is always populated.
func FromScanRecord(idx int, r model.ScanRecord) RepoTermBlock {
	original := pickURLForTransport(r.Transport, r.HTTPSUrl, r.SSHUrl)
	target := original
	cmd := strings.TrimSpace(r.CloneInstruction)
	if len(cmd) == 0 && len(target) > 0 {
		cmd = fmt.Sprintf("git clone %s", target)
	}

	return RepoTermBlock{
		Index:        idx,
		Name:         r.RepoName,
		Branch:       r.Branch,
		BranchSource: r.BranchSource,
		Transport:    r.Transport,
		HTTPSUrl:     r.HTTPSUrl,
		SSHUrl:       r.SSHUrl,
		OriginalURL:  original,
		TargetURL:    target,
		CloneCommand: cmd,
	}
}

// FromScanRecords maps a slice with 1-based indexing.
func FromScanRecords(records []model.ScanRecord) []RepoTermBlock {
	out := make([]RepoTermBlock, 0, len(records))
	for i, r := range records {
		out = append(out, FromScanRecord(i+1, r))
	}

	return out
}

// pickURLForTransport returns the URL whose transport matches the
// repo's identified `origin` transport. SSH-origin repos surface as
// SSH in the per-repo "from / to / command" block so the displayed
// command matches what gitmap will actually invoke (fixes the
// browser-auth-prompt class of bugs where an SSH repo was reported
// as being cloned over HTTPS). Empty preferred URL falls through to
// the other transport so the block is never blank when ANY URL is
// known.
func pickURLForTransport(transport, https, ssh string) string {
	httpsTrim := strings.TrimSpace(https)
	sshTrim := strings.TrimSpace(ssh)
	if transport == constants.ScanTransportSSH && len(sshTrim) > 0 {
		return ssh
	}
	if len(httpsTrim) > 0 {
		return https
	}

	return ssh
}
