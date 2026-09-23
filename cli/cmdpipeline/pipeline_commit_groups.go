package cmdpipeline

import (
	"sort"
	"strings"
)

// CommitWorkflowItem captures an individual workflow execution belonging to a commit.
type CommitWorkflowItem struct {
	DatabaseId uint64 `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Duration   int    `json:"duration"`
	Url        string `json:"url"`
}

// CommitPipelineGroup aggregates all pipeline workflow runs triggered by a single commit SHA.
type CommitPipelineGroup struct {
	HeadSha         string               `json:"headSha"`
	HeadBranch      string               `json:"headBranch"`
	Release         string               `json:"release,omitempty"`
	Status          string               `json:"status"`
	Conclusion      string               `json:"conclusion"`
	CreatedAt       string               `json:"createdAt"`
	UpdatedAt       string               `json:"updatedAt"`
	TotalDuration   int                  `json:"totalDuration"`
	Workflows       []CommitWorkflowItem `json:"workflows"`
	FailedWorkflows int                  `json:"failedWorkflows"`
	PassedWorkflows int                  `json:"passedWorkflows"`
	ActiveWorkflows int                  `json:"activeWorkflows"`
}

type groupBuilder struct {
	orderedShas []string
	groups      map[string]*CommitPipelineGroup
}

// GroupRunsByCommit groups raw workflow runs by commit SHA in chronological order.
func GroupRunsByCommit(runs []ghRunItem) []CommitPipelineGroup {
	if len(runs) == 0 {
		return nil
	}

	builder := newGroupBuilder()
	for _, run := range runs {
		builder.addRun(run)
	}

	return builder.build()
}

func newGroupBuilder() *groupBuilder {
	return &groupBuilder{
		groups: make(map[string]*CommitPipelineGroup),
	}
}

func (b *groupBuilder) addRun(r ghRunItem) {
	group := b.ensureGroup(r)
	wf := buildWorkflowItem(r)
	group.Workflows = append(group.Workflows, wf)
	updateGroupMetrics(group, wf)
	b.updateGroupTimestamps(group, r)
	b.updateGroupRelease(group, r)
}

func (b *groupBuilder) updateGroupRelease(group *CommitPipelineGroup, r ghRunItem) {
	if len(group.Release) > 0 {
		return
	}

	rel := resolveReleaseFromRun(r)
	if len(rel) > 0 {
		group.Release = rel
	}
}

func (b *groupBuilder) updateGroupTimestamps(group *CommitPipelineGroup, r ghRunItem) {
	if r.CreatedAt > group.CreatedAt {
		group.CreatedAt = r.CreatedAt
	}
	if r.UpdatedAt > group.UpdatedAt {
		group.UpdatedAt = r.UpdatedAt
	}
}

func (b *groupBuilder) ensureGroup(r ghRunItem) *CommitPipelineGroup {
	group, isFound := b.groups[r.HeadSha]
	if isFound {
		return group
	}

	return b.createNewGroup(r)
}

func (b *groupBuilder) createNewGroup(r ghRunItem) *CommitPipelineGroup {
	newGroup := initCommitGroup(r)
	b.orderedShas = append(b.orderedShas, r.HeadSha)
	b.groups[r.HeadSha] = newGroup

	return newGroup
}

func initCommitGroup(r ghRunItem) *CommitPipelineGroup {
	return &CommitPipelineGroup{
		HeadSha:    r.HeadSha,
		HeadBranch: r.HeadBranch,
		Release:    resolveReleaseFromRun(r),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func resolveReleaseFromRun(r ghRunItem) string {
	relFromBranch := extractReleaseFromBranch(r.HeadBranch)
	if len(relFromBranch) > 0 {
		return relFromBranch
	}

	relFromTitle := extractReleaseFromTitle(r.DisplayTitle)
	if len(relFromTitle) > 0 {
		return relFromTitle
	}

	return ""
}

func buildWorkflowItem(r ghRunItem) CommitWorkflowItem {
	return CommitWorkflowItem{
		DatabaseId: r.DatabaseId,
		Name:       r.Name,
		Status:     r.Status,
		Conclusion: r.Conclusion,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		Duration:   calculateRunDuration(r.CreatedAt, r.UpdatedAt),
		Url:        r.Url,
	}
}

func updateGroupMetrics(group *CommitPipelineGroup, wf CommitWorkflowItem) {
	group.TotalDuration += wf.Duration
	updateWorkflowCounts(group, wf)
	group.Status = computeAggregateStatus(group.Workflows)
	group.Conclusion = computeAggregateConclusion(group.Workflows)
}

func updateWorkflowCounts(group *CommitPipelineGroup, wf CommitWorkflowItem) {
	if isWorkflowFailure(wf) {
		group.FailedWorkflows++
		return
	}
	if isWorkflowActive(wf) {
		group.ActiveWorkflows++
		return
	}
	group.PassedWorkflows++
}

func isWorkflowFailure(wf CommitWorkflowItem) bool {
	return isFailingConclusion(wf.Conclusion)
}

func isWorkflowActive(wf CommitWorkflowItem) bool {
	return wf.Status == "in_progress" || wf.Status == "queued" || wf.Status == "waiting" || wf.Status == "pending"
}

func computeAggregateStatus(workflows []CommitWorkflowItem) string {
	for _, wf := range workflows {
		if isWorkflowActive(wf) {
			return "in_progress"
		}
	}

	return "completed"
}

func computeAggregateConclusion(workflows []CommitWorkflowItem) string {
	if hasFailingWorkflow(workflows) {
		return "failure"
	}
	if hasActiveWorkflow(workflows) {
		return "in_progress"
	}

	return "success"
}

func hasFailingWorkflow(workflows []CommitWorkflowItem) bool {
	for _, wf := range workflows {
		if isWorkflowFailure(wf) {
			return true
		}
	}

	return false
}

func hasActiveWorkflow(workflows []CommitWorkflowItem) bool {
	for _, wf := range workflows {
		if isWorkflowActive(wf) {
			return true
		}
	}

	return false
}

func (b *groupBuilder) build() []CommitPipelineGroup {
	result := make([]CommitPipelineGroup, 0, len(b.orderedShas))
	for _, sha := range b.orderedShas {
		result = append(result, *b.groups[sha])
	}
	sortCommitGroupsDesc(result)

	return result
}

func sortCommitGroupsDesc(groups []CommitPipelineGroup) {
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].CreatedAt > groups[j].CreatedAt
	})
}

// ResolveCommitGroupByOffset resolves a commit group by negative offset or 0 for latest.
func ResolveCommitGroupByOffset(groups []CommitPipelineGroup, offset int) (*CommitPipelineGroup, bool) {
	if len(groups) == 0 {
		return nil, false
	}

	idx := normalizeOffsetToIndex(offset)
	if idx >= 0 && idx < len(groups) {
		return &groups[idx], true
	}

	return nil, false
}

func normalizeOffsetToIndex(offset int) int {
	if offset < 0 {
		return -offset
	}

	return offset
}

// ResolveCommitGroupByTarget resolves a commit group by offset, keyword, or commit SHA.
func ResolveCommitGroupByTarget(groups []CommitPipelineGroup, target string) (*CommitPipelineGroup, bool) {
	if len(groups) == 0 {
		return nil, false
	}
	if isLatestTarget(target) {
		return &groups[0], true
	}
	if offset, isOffset := tryParseTargetOffset(target); isOffset {
		return ResolveCommitGroupByOffset(groups, offset)
	}

	return findGroupByShaPrefix(groups, target)
}

func isLatestTarget(target string) bool {
	lower := strings.ToLower(strings.TrimSpace(target))

	return len(lower) == 0 || lower == "latest" || lower == "head"
}

func tryParseTargetOffset(target string) (int, bool) {
	trimmed := strings.TrimSpace(target)
	if isZeroOffsetTarget(trimmed) {
		return 0, true
	}
	if offset, isNeg := ParseNegativeIndex(trimmed); isNeg {
		return offset, true
	}

	return 0, false
}

func isZeroOffsetTarget(trimmed string) bool {
	lower := strings.ToLower(trimmed)

	return lower == "0" || lower == "0n" || lower == "head~0" || lower == "~0"
}

func findGroupByShaPrefix(groups []CommitPipelineGroup, target string) (*CommitPipelineGroup, bool) {
	cleanTarget := strings.ToLower(strings.TrimSpace(target))
	for i := range groups {
		if matchSha(groups[i].HeadSha, cleanTarget) {
			return &groups[i], true
		}
	}

	return nil, false
}

func matchSha(sha, target string) bool {
	lower := strings.ToLower(sha)

	return lower == target || strings.HasPrefix(lower, target)
}

// CollapseConsecutiveEmptyLines collapses multiple consecutive empty lines into at most one empty line.
func CollapseConsecutiveEmptyLines(input string) string {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	filtered := filterConsecutiveEmpty(lines)

	return strings.Join(filtered, "\n")
}

func filterConsecutiveEmpty(lines []string) []string {
	var result []string
	hasPrevEmpty := false
	for _, line := range lines {
		result, hasPrevEmpty = appendNonDuplicateEmpty(result, line, hasPrevEmpty)
	}

	return result
}

func appendNonDuplicateEmpty(result []string, line string, hasPrevEmpty bool) ([]string, bool) {
	isEmpty := len(strings.TrimSpace(line)) == 0
	if isEmpty && hasPrevEmpty {
		return result, true
	}
	if isEmpty {
		return append(result, ""), true
	}

	return append(result, line), false
}
