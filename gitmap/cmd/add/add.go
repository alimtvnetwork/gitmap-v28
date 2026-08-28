package add

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/pterm/pterm"
)

// Run executes the add command logic.
func Run(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("AddRun", "E_ADD_MISSING_ARG")
	}
	target := args[0]
	switch target {
	case "common-attr":
		return writeCommonAttr()
	case "common-ignore":
		return writeCommonIgnore()
	default:
		return apperror.New("AddRun", "E_ADD_UNKNOWN_TARGET", map[string]any{"target": target})
	}
}

func writeCommonAttr() error {
	content := "* text=auto eol=lf\n*.png binary\n*.jpg binary\n*.jpeg binary\n*.gif binary\n*.ico binary\n"
	if err := os.WriteFile(".gitattributes", []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "writeCommonAttr: failed to write")
	}
	pterm.Success.Println("Created .gitattributes")
	return nil
}

func writeCommonIgnore() error {
	content := "node_modules/\n.DS_Store\nThumbs.db\n.idea/\n.vscode/\n*.log\n*.tmp\n*.bak\n.env\ndist/\nbuild/\ncoverage/\n"
	if err := os.WriteFile(".gitignore", []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "writeCommonIgnore: failed to write")
	}
	pterm.Success.Println("Created .gitignore")
	return nil
}
