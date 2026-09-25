package pipelinedb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func (p *PipelineSplitDb) setupSchema() error {
	if err := p.InitSchema(); err != nil {
		_ = p.conn.Close()

		return err
	}

	return nil
}

func pipelineSchemaQueries() []string {
	return []string{
		sqlCreatePipelineRun,
		sqlCreatePipelineErrorLog,
		sqlCreatePipelineDetailErrorLog,
		sqlCreatePipelineCompactErrorLog,
		sqlCreatePipelineJob,
		sqlCreatePipelineSegment,
	}
}

func (p *PipelineSplitDb) executeSchemaQueries(queries []string) error {
	for _, q := range queries {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "init pipeline db schema")
		}
	}

	return nil
}

// InitSchema ensures all pipeline tables exist.
func (p *PipelineSplitDb) InitSchema() error {
	return p.executeSchemaQueries(pipelineSchemaQueries())
}
