package pipelinedb

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// QueryRunById returns the pipeline run record by its GitHub run ID.
func (p *PipelineSplitDb) QueryRunById(runId uint64) (*PipelineRunRecord, error) {
	query := `SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
		EtaSeconds, DurationSeconds, RunUrl, IsSuccess, Notes, Comments, CreatedAt, UpdatedAt
		FROM PipelineRun WHERE RunId = ? LIMIT 1;`
	var r PipelineRunRecord
	var isSuccessInt int
	var notes, comments sql.NullString
	err := p.conn.QueryRow(query, runId).Scan(
		&r.RunId, &r.RepoSlug, &r.WorkflowName, &r.Status, &r.Conclusion,
		&r.Branch, &r.Sha, &r.EtaSeconds, &r.DurationSeconds, &r.RunUrl,
		&isSuccessInt, &notes, &comments, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query run by id")
	}
	r.IsSuccess = isSuccessInt != 0
	r.Notes = notes.String
	r.Comments = comments.String
	return &r, nil
}
