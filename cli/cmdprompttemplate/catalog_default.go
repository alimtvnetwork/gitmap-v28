package cmdprompttemplate

import "time"

// GetDefaultCategoryTemplates returns 20 distinct verification and quality prefix templates.
func GetDefaultCategoryTemplates() []PromptTemplate {
	now := time.Now().UTC()
	entries := []struct {
		slug string
		name string
		text string
	}{
		{
			slug: "fix-below-code",
			name: "Fix Below Code",
			text: "Fix the below code.",
		},
		{
			slug: "confident-completed-task",
			name: "Confident About Completed Task",
			text: "Have you fixed the below code? Are you confident about your completed task? I don't think that you have done a good job. Can you please recheck and verify the work has been done or not?",
		},
		{
			slug: "avoid-sloppy-work",
			name: "Avoid Sloppy Quick Work",
			text: "Are you stupid? Do you think that doing sloppy work will yield something? Please don't do quick work. Spend some time, make sure the task you are working on, that actually is done professionally. So make sure the below task is completed without issue, without a fault.",
		},
		{
			slug: "confirm-no-stupidity",
			name: "Confirm Confidence and Rigor",
			text: "Can you please confirm the below task has not missed out anyway, and there is no stupidity that you have done? Can you please confirm your confidence and do the work?",
		},
		{
			slug: "checked-quality-alert",
			name: "Alert Quality Review",
			text: "I have checked the quality of the work. It was sloppy. I do not appreciate that. WTF. You should make sure the tasks are done professionally. Please be alert and do the task professionally.",
		},
		{
			slug: "follow-guideline-rigor",
			name: "Follow Guideline Rigor",
			text: "Does anyone pay you to be sloppy and stupid? WTF. What you are thinking of? Are you stupid? Why don't you fix and complete the below task professionally? And also, while you do that, make sure you follow the coding guideline and nothing is missed.",
		},
		{
			slug: "engineering-not-single-prompt",
			name: "Engineering Depth Discipline",
			text: "Do you think engineering is just reading a single prompt and doing the work? What do you think of yourself? Engineering is reading through the coding guideline, making sure instructions are followed properly, there is nothing remaining. You haven't done it, so that's stupid. Please fix the below stuff properly.",
		},
		{
			slug: "patience-zero-shortcuts",
			name: "Patience and Zero Shortcuts",
			text: "Stop rushing through your code changes. Craftsmanship requires patience, deep verification, and zero shortcuts. Ensure every edge case is handled cleanly in the below task.",
		},
		{
			slug: "eliminate-nesting-swallowed",
			name: "Eliminate Nesting and Swallowed Errors",
			text: "Verify your implementation against strict repository standards. Check every error envelope, eliminate nesting, and ensure zero swallowed errors before marking this task complete.",
		},
		{
			slug: "resilient-architecture",
			name: "Resilient Architecture First",
			text: "Take a step back and examine the full architecture before editing. Sloppy patches degrade maintainability. Deliver clean, elegant, and resilient code for the task below.",
		},
		{
			slug: "inspect-commits-context",
			name: "Inspect Commits and Context",
			text: "Did you actually inspect recent git commits and memory files, or did you make blind assumptions? Review the context thoroughly and execute the following task with precision.",
		},
		{
			slug: "conventions-double-check",
			name: "Conventions Double Check",
			text: "Ensure that no unit tests break and all conventions are strictly respected. Re-check your logic twice and complete the below task with full confidence.",
		},
		{
			slug: "professional-anti-patterns",
			name: "Avoid Anti-Patterns",
			text: "Professional engineers do not write half-baked solutions. Follow the repository's explicit rules, avoid all anti-patterns, and finalize the below work properly.",
		},
		{
			slug: "meticulous-type-safety",
			name: "Meticulous Type Safety",
			text: "Careless mistakes cost time and compute. Pay meticulous attention to type safety, error wrapping, and naming standards while completing the task below.",
		},
		{
			slug: "zero-regressions-guarantee",
			name: "Zero Regressions Guarantee",
			text: "Can you guarantee that this change introduces zero regressions? Double-check all affected packages, verify error handling, and complete the instructions below.",
		},
		{
			slug: "ground-specifications",
			name: "Ground in Specifications",
			text: "Do not assume or guess. Ground every design decision in the specifications. Complete the below task with verified accuracy and thoroughness.",
		},
		{
			slug: "strictest-reviewer-mindset",
			name: "Strictest Reviewer Mindset",
			text: "Review your code changes as if you were the strictest code reviewer on earth. Eliminate any lazy shortcuts and execute the below task cleanly.",
		},
		{
			slug: "anticipate-failures",
			name: "Anticipate Failures Defensively",
			text: "A true senior engineer anticipates failures before they occur. Defensively validate every parameter, handle every error, and solve the task below thoroughly.",
		},
		{
			slug: "naming-and-silent-regressions",
			name: "Check Naming and Silent Regressions",
			text: "Inspect your changes for silent regressions, improper error swallowing, and naming rule violations. Refactor and complete the task below properly.",
		},
		{
			slug: "commit-to-excellence",
			name: "Commit to Excellence",
			text: "Commit to excellence. No sloppy logic, no skipped guidelines, no rushed conclusions. Implement the below specification completely and cleanly.",
		},
	}

	result := make([]PromptTemplate, 0, len(entries))
	for _, e := range entries {
		result = append(result, PromptTemplate{
			ID:          "def-" + e.slug,
			Name:        e.name,
			Category:    CategoryDefault,
			Slug:        e.slug,
			Description: e.name,
			Content:     e.text,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	return result
}
