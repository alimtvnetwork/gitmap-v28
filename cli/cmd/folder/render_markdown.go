package folder

import (
	"fmt"
	"strings"
)

// RenderMarkdown converts a TreeNode graph into an indented, nested Markdown bullet list without backticks.
func RenderMarkdown(root *TreeNode, isDetailed bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- %s/\n", root.Name))
	renderMarkdownLevel(root, "  ", &sb, isDetailed)

	return sb.String()
}

func renderMarkdownLevel(node *TreeNode, indent string, sb *strings.Builder, isDetailed bool) {
	children := sortTreeChildren(node.Children)
	for _, child := range children {
		if child.IsDir {
			sb.WriteString(fmt.Sprintf("%s- %s/\n", indent, child.Name))
			renderMarkdownLevel(child, indent+"  ", sb, isDetailed)
		} else {
			metaInfo := formatMarkdownMeta(child.Meta, isDetailed)
			sb.WriteString(fmt.Sprintf("%s- %s%s\n", indent, child.Name, metaInfo))
		}
	}
}

func formatMarkdownMeta(meta *FileMeta, isDetailed bool) string {
	if !isDetailed || meta == nil {
		return ""
	}

	var parts []string
	if meta.Sequence > 0 {
		parts = append(parts, fmt.Sprintf("seq: %02d", meta.Sequence))
	}

	if !meta.IsBinary && meta.LinesOfCode > 0 {
		parts = append(parts, fmt.Sprintf("%d lines", meta.LinesOfCode))
	}

	if meta.SizeFormatted != "" {
		parts = append(parts, meta.SizeFormatted)
	}

	if len(parts) == 0 {
		return ""
	}

	return " (" + strings.Join(parts, ", ") + ")"
}
