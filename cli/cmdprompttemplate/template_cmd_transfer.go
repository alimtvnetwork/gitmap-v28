package cmdprompttemplate

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runExportTemplate(args []string) error {
	if len(args) < 2 {
		return apperror.NewValidationError("usage: gitmap prompts-template export <name|id> <dest.json>")
	}
	if err := ExportSingleTemplate(args[0], args[1]); err != nil {
		return err
	}
	fmt.Printf("  %s✔ Exported template '%s' to %s%s\n", constants.ColorGreen, args[0], args[1], constants.ColorReset)

	return nil
}

func runImportTemplate(args []string) error {
	if len(args) < 1 {
		return apperror.NewValidationError("usage: gitmap prompts-template import <src.json>")
	}
	tpl, err := ImportSingleTemplate(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("  %s✔ Imported prompt template '%s'%s\n", constants.ColorGreen, tpl.Name, constants.ColorReset)

	return nil
}

func runExportAllTemplates(args []string) error {
	dest := "prompts-templates.json"
	if len(args) > 0 {
		dest = args[0]
	}
	if err := ExportAllTemplates(dest); err != nil {
		return err
	}
	fmt.Printf("  %s✔ Exported all prompt templates to %s%s\n", constants.ColorGreen, dest, constants.ColorReset)

	return nil
}

func runImportAllTemplates(args []string) error {
	if len(args) < 1 {
		return apperror.NewValidationError("usage: gitmap prompts-template import-all <src.json>")
	}
	count, err := ImportAllTemplates(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("  %s✔ Imported %d prompt template(s) from %s%s\n", constants.ColorGreen, count, args[0], constants.ColorReset)

	return nil
}
