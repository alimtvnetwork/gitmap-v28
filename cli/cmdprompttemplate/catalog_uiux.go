package cmdprompttemplate

import "time"

// GetUIUXCategoryTemplates returns 20 distinct UI/UX design and interaction prefix templates.
func GetUIUXCategoryTemplates() []PromptTemplate {
	now := time.Now().UTC()
	entries := []struct {
		slug string
		name string
		text string
	}{
		{
			slug: "gstac-design-principles",
			name: "Learn GSTAC UX Principles",
			text: "What a dumb design that you have done. You have not followed the user experience principles. You should be learning the user experience from GSTAC, the coding guideline design system folder, and also the references that has been given below. A lot of places the sloppy work has been done.",
		},
		{
			slug: "engineering-ux-depth",
			name: "Engineering UX Depth",
			text: "Make sure that you have completed the design aspects and UX is done in terms of proper engineering. It's not just doing the work. Okay? Can you please confirm that you have followed through things properly? Can you please confirm that?",
		},
		{
			slug: "avoid-stupidity-design-rating",
			name: "Raise Design Rating and Avoid Sloppiness",
			text: "I don't believe the design that you have done. It's even pass a rating of two out of 10. Please make sure you avoid stupidity, WTF, and do the work properly. Follow the instruction without sloppiness. Spend some time, read the documents, read the design aspects, make things better. Follow the instructions properly. Don't make me feel that you are stupid.",
		},
		{
			slug: "design-sense-accuracy",
			name: "Design Sense and Visual Verification",
			text: "Are you stupid with design? You cannot be that. Designs are very sensitive stuff. You have to be very much accurate and have the design sense. You cannot just do bunch of bad stuff and say it is done. You cannot do that. You have to make sure and verify after the design is done. In terms of looking into UI, you have the capability, and then understand what is missing and what I'm looking for.",
		},
		{
			slug: "spawn-agents-follow-commands",
			name: "Spawn Agents and Follow Directives",
			text: "It takes a genius to understand what is wrong. It takes a stupid to do the work. Well, they have no clue. You have done the stupidity, and I do not appreciate that. WTF? Why you are not following the instruction properly? Can you please follow the instruction properly and do the task effectively? Why you are not spawning agents? Why? Why you are not following the commands properly? Who pays you for this stupidity? Okay, follow the below instructions properly.",
		},
		{
			slug: "uiux-not-cheap",
			name: "UI/UX Is Not Cheap",
			text: "UI/UX is not a cheap thing, okay? Don't play cheap and dumb stuff with me. Try to have some UI/UX sense, and then try to do it. Can you please follow the below instructions properly?",
		},
		{
			slug: "clarity-hierarchy-communication",
			name: "Visual Hierarchy and Communication",
			text: "Design is communication, hierarchy, and clarity. Stop throwing unstyled elements onto the screen and build a cohesive, accessible interface for the task below.",
		},
		{
			slug: "contrast-and-typography-gstac",
			name: "Color Contrast and Typography Scale",
			text: "Verify color contrast, typography scale, and responsive padding according to GSTAC design guidelines. Complete the UI/UX task below with taste.",
		},
		{
			slug: "interaction-states-polish",
			name: "All Interaction States Polish",
			text: "User experience requires empathy and attention to interaction states: hover, active, focus, disabled, and loading. Polish every state in the task below.",
		},
		{
			slug: "a11y-keyboard-focus-rings",
			name: "Keyboard Accessibility and A11y",
			text: "Do not neglect keyboard accessibility (a11y), screen reader labels, and focus rings. Elevate the UI implementation below to enterprise standards.",
		},
		{
			slug: "visual-hierarchy-whitespace",
			name: "Whitespace and Visual Breathing Room",
			text: "Visual hierarchy is non-negotiable. Ensure primary actions dominate, secondary actions recede, and whitespace breathes naturally in the below task.",
		},
		{
			slug: "responsive-breakpoints-flow",
			name: "Responsive Breakpoints Flow",
			text: "Review the responsive layout across mobile, tablet, and widescreen breakpoints. Eliminate clipping, awkward wraps, and overflow in the task below.",
		},
		{
			slug: "snappy-transitions",
			name: "Snappy Intentional Motion",
			text: "Animations and transitions must feel snappy and intentional (150-250ms ease-out), never sluggish or disorienting. Refine the UI interactions below.",
		},
		{
			slug: "strict-spacing-grid",
			name: "Strict Spacing Grid Alignment",
			text: "Inconsistent spacing destroys visual credibility. Align all margins, gaps, and paddings to the 4px/8px design grid for the following UI component.",
		},
		{
			slug: "typography-contrast-ratios",
			name: "Typography and Letter Spacing",
			text: "Typography is 90% of user interface design. Check font weights, line heights, letter spacing, and contrast ratios in the below implementation.",
		},
		{
			slug: "eliminate-cls-layout-shifts",
			name: "Smooth Mount and No CLS",
			text: "Avoid abrupt layout shifts (CLS) and stuttering re-renders. Ensure components mount smoothly with appropriate skeletons and loaders in the task below.",
		},
		{
			slug: "iconography-visual-balance",
			name: "Icon Balance and Alignment",
			text: "Icons must be consistently sized, visually balanced, and accompanied by accessible text alternatives. Implement the UI elements below with care.",
		},
		{
			slug: "edge-case-content-overflow",
			name: "Data Edge Cases and Empty States",
			text: "Test the interface with extreme data lengths: empty states, single-character labels, and 200-character overflow strings. Harden the UI below.",
		},
		{
			slug: "modals-popovers-portals",
			name: "Modals and Popover Portals",
			text: "Refine modal dialogues, popovers, and tooltips with proper dismiss triggers, backdrop blurs, and portal mounting. Polish the design below.",
		},
		{
			slug: "portfolio-grade-craft",
			name: "Portfolio Grade Craftsmanship",
			text: "Treat every screen as a portfolio piece. Deliver a modern, polished, and delightful user experience for the instructions below.",
		},
	}

	result := make([]PromptTemplate, 0, len(entries))
	for _, e := range entries {
		result = append(result, PromptTemplate{
			ID:          "ui-" + e.slug,
			Name:        e.name,
			Category:    CategoryUIUX,
			Slug:        e.slug,
			Description: e.name,
			Content:     e.text,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	return result
}
