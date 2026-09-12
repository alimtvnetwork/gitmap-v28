package cmd

import (
	"testing"
)

func TestInstallerTree_Render(t *testing.T) {
	t.Parallel()

	root := InstallerTreeNode{
		Title:       "Test Root",
		Description: "Root Description",
		Children: []InstallerTreeNode{
			{
				Title:       "Child 1",
				Description: "Desc 1",
				Children: []InstallerTreeNode{
					{
						Title:       "Subchild 1",
						Description: "Subdesc 1",
					},
				},
			},
			{
				Title:       "Child 2",
				Description: "Desc 2",
			},
		},
	}

	// Verify rendering does not panic
	printInstallerTree(root, "", true)
	printInstallSummaryHeader("test-slug")
}
