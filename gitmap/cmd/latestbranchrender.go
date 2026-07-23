package cmd

// JSON encoder for `gitmap latest-branch --json`.
//
// Migrated off json.NewEncoder(latestBranchJSON) onto gitmap/stablejson
// so key order becomes a compile-time decision rather than a
// reflection accident. The nested `top` array is pre-rendered in
// compact mode and embedded as json.RawMessage.
//
// Schema: spec/08-json-schemas/latest-branch.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// latest-branch top-level wire keys. Names + order are the contract.
const (
	lbKeyBranch     = "branch"
	lbKeyRemote     = "remote"
	lbKeySha        = "sha"
	lbKeyCommitDate = "commitDate"
	lbKeySubject    = "subject"
	lbKeyRef        = "ref"
	lbKeyTop        = "top"
)

// latest-branch top-item wire keys.
const (
	lbTopKeyBranch     = "branch"
	lbTopKeySha        = "sha"
	lbTopKeyCommitDate = "commitDate"
	lbTopKeySubject    = "subject"
)

// encodeLatestBranchJSON writes the latest-branch result as stable
// JSON with 2-space indentation. When top == 0 the `top` key is
// omitted; when top > 0 it is included as a nested array.
func encodeLatestBranchJSON(
	w io.Writer, result latestBranchResult,
	items []gitutil.RemoteBranchInfo, top int,
) error {
	fields := []stablejson.Field{
		{Key: lbKeyBranch, Value: result.branchNames},
		{Key: lbKeyRemote, Value: result.selectedRemote},
		{Key: lbKeySha, Value: result.shortSha},
		{Key: lbKeyCommitDate, Value: result.commitDate},
		{Key: lbKeySubject, Value: result.latest.Subject},
		{Key: lbKeyRef, Value: result.latest.RemoteRef},
	}

	if top > 0 {
		topRaw, err := renderLatestBranchTopRaw(items, top)
		if err != nil {
			return err
		}
		fields = append(fields, stablejson.Field{Key: lbKeyTop, Value: topRaw})
	}

	return stablejson.WriteObject(w, fields)
}

// renderLatestBranchTopRaw pre-renders the top-N array in compact
// mode for embedding as a top-level object value.
func renderLatestBranchTopRaw(items []gitutil.RemoteBranchInfo, top int) (json.RawMessage, error) {
	count := top
	if count > len(items) {
		count = len(items)
	}

	topItems := make([][]stablejson.Field, 0, count)
	for _, item := range items[:count] {
		topItems = append(topItems, []stablejson.Field{
			{Key: lbTopKeyBranch, Value: gitutil.StripRemotePrefix(item.RemoteRef)},
			{Key: lbTopKeySha, Value: gitutil.TruncSha(item.Sha)},
			{Key: lbTopKeyCommitDate, Value: gitutil.FormatDisplayDate(item.CommitDate)},
			{Key: lbTopKeySubject, Value: item.Subject},
		})
	}

	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, topItems, ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}
