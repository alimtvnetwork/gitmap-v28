package orchestrator

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/runlog"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/workspace"
)

// inputRepoID returns the InputRepo PK for a staged input, inserting
// the row on first call and caching it thereafter so a single input
// produces exactly one InputRepo row regardless of commit count.
func (c *runContext) inputRepoID(staged workspace.StagedInput) (int64, error) {
	if c.inputRepoIds == nil {
		c.inputRepoIds = map[int]int64{}
	}
	idx := staged.Input.OrderIndex
	if id, ok := c.inputRepoIds[idx]; ok {
		return id, nil
	}
	id, err := runlog.InsertInputRepo(c.DB.Conn(), c.RunID, idx, staged.Input.Original, staged.WorkPath, staged.Input.Kind)
	if err != nil {
		return 0, err
	}
	c.inputRepoIds[idx] = id
	return id, nil
}

// firstLine returns the first line of msg (handy for the success log).
func firstLine(msg string) string {
	if i := strings.IndexByte(msg, '\n'); i >= 0 {
		return msg[:i]
	}
	return msg
}
