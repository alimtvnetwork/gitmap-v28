package cmdprompttemplate

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runListTemplates() error {
	list, err := LoadTemplates()
	if err != nil {
		return err
	}
	RenderTemplatesTable(list)

	return nil
}

func runAddTemplate(args []string) error {
	if len(args) < 2 {
		return apperror.NewValidationError("usage: gitmap prompts-template add <name> <content>")
	}
	name := args[0]
	content := strings.Join(args[1:], " ")
	tpl, err := AddTemplate(name, content)
	if err != nil {
		return err
	}
	fmt.Printf("  %s✔ Added prompt template '%s'%s\n", constants.ColorGreen, tpl.Name, constants.ColorReset)

	return nil
}

func runEditTemplate(args []string) error {
	if len(args) < 2 {
		return apperror.NewValidationError("usage: gitmap prompts-template edit <name|id> <content>")
	}
	tpl, err := EditTemplate(args[0], strings.Join(args[1:], " "))
	if err != nil {
		return err
	}
	fmt.Printf("  %s✔ Updated prompt template '%s'%s\n", constants.ColorGreen, tpl.Name, constants.ColorReset)

	return nil
}

func runDeleteTemplate(args []string) error {
	if len(args) < 1 {
		return apperror.NewValidationError("usage: gitmap prompts-template rm <name|id>")
	}
	if err := DeleteTemplate(args[0]); err != nil {
		return err
	}
	fmt.Printf("  %s✔ Removed prompt template '%s'%s\n", constants.ColorGreen, args[0], constants.ColorReset)

	return nil
}
