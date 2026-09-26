package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func dispatchStateTemplatesSub(sub string, args []string) bool {
	switch sub {
	case "ls", "state-ls":
		runTemplatesStateList(args)
	case "add":
		runTemplatesStateAdd(args)
	case "edit":
		runTemplatesStateEdit(args)
	case "remove", "rm", "delete":
		runTemplatesStateRemove(args)
	default:
		return dispatchStateTemplatesIOAndUI(sub, args)
	}

	return true
}

func dispatchStateTemplatesIOAndUI(sub string, args []string) bool {
	switch sub {
	case "import":
		runTemplatesStateImport(args)
	case "export":
		runTemplatesStateExport(args)
	case "var", "vars":
		runTemplatesStateVar(args)
	case "ui", "web":
		runTemplatesStateUI(args)
	default:
		return false
	}

	return true
}

func openStateTemplatesDBOrExit() *store.TemplatesSplitDB {
	db, err := store.OpenTemplatesSplitDB()
	if err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmd.templates.open_db"), 1)

		return nil
	}

	return db
}

func runTemplatesStateList(args []string) {
	fs := flag.NewFlagSet("templates-ls", flag.ExitOnError)
	cat := fs.String("category", "", "Filter by category slug")
	isJSON := fs.Bool("json", false, "Emit JSON output")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	renderStateTemplatesList(db, strings.TrimSpace(*cat), *isJSON)
}

func renderStateTemplatesList(db *store.TemplatesSplitDB, category string, isJSON bool) {
	items, err := db.ListTemplateItems(category)
	if err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	cats, _ := db.ListCategories()
	if isJSON {
		writeJSONStdout(map[string]any{"categories": cats, "templates": items})

		return
	}

	printStateTemplatesTable(items)
}

func printStateTemplatesTable(items []store.StateTemplateItem) {
	if len(items) == 0 {
		fmt.Println("(no state templates stored in gitmap-templates.db)")

		return
	}

	fmt.Printf("%-18s  %-14s  %-26s  %s\n", "ID", "CATEGORY", "SLUG", "TITLE")
	for _, it := range items {
		fmt.Printf("%-18s  %-14s  %-26s  %s\n", it.ID, it.Category, it.Slug, it.Title)
	}
}

func runTemplatesStateAdd(args []string) {
	fs := flag.NewFlagSet("templates-add", flag.ExitOnError)
	id := fs.String("id", "", "Optional template ID")
	cat := fs.String("category", "seo", "Template category slug")
	slug := fs.String("slug", "", "Template slug")
	title := fs.String("title", "", "Template title / question")
	text := fs.String("text", "", "Template body text")
	addRaw := fs.String("additional", "", "Optional JSON object string")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	upsertFromAddFlags(*id, *cat, *slug, *title, *text, *addRaw)
}

func upsertFromAddFlags(id, cat, slug, title, text, addRaw string) {
	if strings.TrimSpace(title) == "" && strings.TrimSpace(text) == "" {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.add", "templates add requires --title or --text"), 1)

		return
	}

	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	item := store.StateTemplateItem{
		ID: id, Category: cat, Slug: slug, Title: title, Text: text,
		Additional: decodeAdditionalFlag(addRaw),
	}
	saveAndPrintTemplateItem(db, item, "Added")
}

func saveAndPrintTemplateItem(db *store.TemplatesSplitDB, item store.StateTemplateItem, verb string) {
	saved, err := db.UpsertTemplateItem(item)
	if err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	fmt.Printf("%s template %s (slug=%s, category=%s)\n", verb, saved.ID, saved.Slug, saved.Category)
}

func runTemplatesStateEdit(args []string) {
	fs := flag.NewFlagSet("templates-edit", flag.ExitOnError)
	cat := fs.String("category", "", "Updated category slug")
	title := fs.String("title", "", "Updated title")
	text := fs.String("text", "", "Updated body text")
	addRaw := fs.String("additional", "", "Updated additional JSON")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	if len(fs.Args()) == 0 {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.edit", "templates edit requires <id|slug>"), 1)

		return
	}

	applyTemplateEdit(fs.Args()[0], *cat, *title, *text, *addRaw)
}

func applyTemplateEdit(target, cat, title, text, addRaw string) {
	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	existing, err := db.GetTemplateItem(target)
	if err != nil || existing == nil {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.edit.not_found", "template not found: "+target), 1)

		return
	}

	mergeTemplateEdits(existing, cat, title, text, addRaw)
	saveAndPrintTemplateItem(db, *existing, "Updated")
}

func mergeTemplateEdits(item *store.StateTemplateItem, cat, title, text, addRaw string) {
	if strings.TrimSpace(cat) != "" {
		item.Category = strings.TrimSpace(cat)
	}
	if title != "" {
		item.Title = title
	}
	if text != "" {
		item.Text = text
	}
	if strings.TrimSpace(addRaw) != "" {
		item.Additional = decodeAdditionalFlag(addRaw)
	}
}

func runTemplatesStateRemove(args []string) {
	if len(args) == 0 {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.rm", "templates remove requires <id|slug>"), 1)

		return
	}

	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	if err := db.DeleteTemplateItem(args[0]); err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	fmt.Printf("Removed template %s\n", args[0])
}

func runTemplatesStateImport(args []string) {
	fs := flag.NewFlagSet("templates-import", flag.ExitOnError)
	isForce := fs.Bool("force", false, "Force re-import even if exportId matches history")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	if len(fs.Args()) == 0 {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.import", "templates import requires <file.json>"), 1)

		return
	}

	executeTemplateImport(fs.Args()[0], *isForce)
}

func executeTemplateImport(srcPath string, isForce bool) {
	count, isSkipped, exportID, err := store.ImportTemplatesFromFile(srcPath, isForce)
	if err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	if isSkipped {
		fmt.Printf("Skipped import: %s already imported (exportId=%s)\n", srcPath, exportID)

		return
	}

	fmt.Printf("Imported %d template(s) from %s (exportId=%s)\n", count, srcPath, exportID)
}

func runTemplatesStateExport(args []string) {
	fs := flag.NewFlagSet("templates-export", flag.ExitOnError)
	cat := fs.String("category", "", "Filter export by category slug")
	idOrSlug := fs.String("id", "", "Filter export by template ID or slug")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	if len(fs.Args()) == 0 {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.export", "templates export requires <dest.json>"), 1)

		return
	}

	executeTemplateExport(fs.Args()[0], *cat, *idOrSlug)
}

func executeTemplateExport(destPath, cat, idOrSlug string) {
	payload, err := store.ExportTemplatesToFile(destPath, cat, idOrSlug)
	if err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	fmt.Printf("Exported %d template(s) to %s (exportId=%s)\n", len(payload.Templates), destPath, payload.ExportID)
}

func runTemplatesStateVar(args []string) {
	if len(args) == 0 {
		runTemplatesVarList(nil)

		return
	}

	switch args[0] {
	case "set":
		runTemplatesVarSet(args[1:])
	case "ls", "list":
		runTemplatesVarList(args[1:])
	case "rm", "remove", "delete":
		runTemplatesVarRemove(args[1:])
	default:
		cliexit.HandleError(apperror.NewSimple("cmd.templates.var", "unknown var action: "+args[0]), 1)
	}
}

func runTemplatesVarSet(args []string) {
	fs := flag.NewFlagSet("templates-var-set", flag.ExitOnError)
	scope := fs.String("scope", "global", "Variable scope")
	desc := fs.String("desc", "", "Variable description")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	if len(fs.Args()) < 2 {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.var.set", "usage: gitmap templates var set <key> <value>"), 1)

		return
	}

	saveTemplateVar(fs.Args()[0], *scope, fs.Args()[1], *desc)
}

func saveTemplateVar(key, scope, val, desc string) {
	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	if err := db.UpsertVariable(key, scope, val, desc); err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	fmt.Printf("Set variable %s (%s)\n", key, scope)
}

func runTemplatesVarList(args []string) {
	fs := flag.NewFlagSet("templates-var-ls", flag.ExitOnError)
	scope := fs.String("scope", "", "Variable scope filter")
	isJSON := fs.Bool("json", false, "Emit JSON")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	vars, err := db.ListVariables(*scope)
	if err != nil {
		cliexit.HandleError(err, 1)

		return
	}

	printTemplateVars(vars, *isJSON)
}

func printTemplateVars(vars map[string]string, isJSON bool) {
	if isJSON {
		writeJSONStdout(vars)

		return
	}

	for k, v := range vars {
		fmt.Printf("%s=%s\n", k, v)
	}
}

func runTemplatesVarRemove(args []string) {
	fs := flag.NewFlagSet("templates-var-rm", flag.ExitOnError)
	scope := fs.String("scope", "global", "Variable scope")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	if len(fs.Args()) == 0 {
		cliexit.HandleError(apperror.NewSimple("cmd.templates.var.rm", "usage: gitmap templates var rm <key>"), 1)

		return
	}

	db := openStateTemplatesDBOrExit()
	if db == nil {
		return
	}
	defer db.Close()

	_ = db.DeleteVariable(fs.Args()[0], *scope)
	fmt.Printf("Removed variable %s\n", fs.Args()[0])
}

func decodeAdditionalFlag(raw string) map[string]any {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}

	var out map[string]any
	_ = json.Unmarshal([]byte(trimmed), &out)

	return out
}

func writeJSONStdout(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
