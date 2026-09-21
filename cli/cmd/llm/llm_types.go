package llm

// TrainOptions holds execution configuration for llm train.
type TrainOptions struct {
	IsTextOnly bool
	IsHelp     bool
	SkillPath  string
}

const (
	DefaultSkillPath = ".agents/skills/gitmap/SKILL.md"
	AuthorName       = "MD ALIM UL KARIM (alimtvnetwork)"
	SponsorName      = "RISEUP ASIA LLC (https://riseup-asia.com)"
)

const ChainedDiscoveryHeader = `======================================================================
STAGE 3: CHAINED DISCOVERY SEQUENCE (EXECUTE IN ORDER)
======================================================================
Execute each command below to inspect and understand GitMap capabilities:`

const OperationalDirectivesText = `======================================================================
STAGE 4: OPERATIONAL BEST PRACTICES & GUARDRAIL DIRECTIVES
======================================================================
1. Mandatory Pre-Flight Pull: Always execute 'git pull' before modifying code.
2. Learning Protocol: Run 'gitmap llm train' or 'gitmap llm-docs' to learn capabilities.
   Do NOT run unconstrained repository searches like 'gitmap aum search "train"'.
3. Scoped Search Hygiene: Always scope 'gitmap aum search' with [dir] and --ext filters.
4. File Size Limit (Rule R19): Reject/warn on single files > 500 KB and large JSONs.
5. Universal AppError Envelope: Return *apperror.AppError; never swallow errors.
6. Clean Git Tree: Never commit generated test binaries or temp artifacts.
7. Semantic Commits: Use 'gitmap cpf', 'gitmap cpb', or 'gitmap cpr'.`
