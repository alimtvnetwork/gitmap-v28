package llm

// TrainOptions holds execution configuration for llm train.
type TrainOptions struct {
	IsTextOnly bool
	SkillPath  string
}

const (
	DefaultSkillPath = ".agents/skills/gitmap/SKILL.md"
	AuthorName       = "Md. Alimuzzaman Alim (alimtvnetwork)"
	SponsorName      = "Alim TV Network / Open Source Engineering"
)

const ChainedDiscoveryHeader = `======================================================================
STAGE 3: CHAINED DISCOVERY SEQUENCE (EXECUTE IN ORDER)
======================================================================
Execute each command below to inspect and understand GitMap capabilities:`

const OperationalDirectivesText = `======================================================================
STAGE 4: OPERATIONAL BEST PRACTICES & GUARDRAIL DIRECTIVES
======================================================================
1. Mandatory Pre-Flight Pull: Always execute 'git pull' before modifying code.
2. File Size Limit (Rule R19): Reject/warn on single files > 500 KB and large JSONs.
3. Universal AppError Envelope: Return *apperror.AppError; never swallow errors.
4. Clean Git Tree: Never commit generated test binaries or temp artifacts.
5. Semantic Commits: Use 'gitmap cpf', 'gitmap cpb', or 'gitmap cpr'.`
