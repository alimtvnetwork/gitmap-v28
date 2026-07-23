// Package cmd — seowritetemplate.go handles template loading and placeholder substitution.
package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// templateFile represents the JSON structure for seed/custom templates.
type templateFile struct {
	Titles       []string `json:"titles"`
	Descriptions []string `json:"descriptions"`
}

// loadTemplateMessages generates commit messages from DB or custom file.
func loadTemplateMessages(flags seoWriteFlags) []commitMessage {
	titles, descriptions := loadTemplatePairs(flags)
	if len(titles) == 0 || len(descriptions) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrSEOTemplateEmpty)
		os.Exit(1)
	}

	return generateMessages(titles, descriptions, flags)
}

// loadTemplatePairs loads title/description slices from file or DB.
func loadTemplatePairs(flags seoWriteFlags) ([]string, []string) {
	if flags.templatePath != "" {
		return loadFromJSONFile(flags.templatePath)
	}

	return loadFromDatabase()
}

// loadFromJSONFile reads templates from a custom JSON file.
func loadFromJSONFile(path string) ([]string, []string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOTemplateRead, path, err)
		os.Exit(1)
	}

	var tf templateFile
	if err := json.Unmarshal(data, &tf); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOTemplateRead, path, err)
		os.Exit(1)
	}

	return tf.Titles, tf.Descriptions
}

// loadFromDatabase loads templates from the SQLite CommitTemplates table.
func loadFromDatabase() ([]string, []string) {
	db := openSEODatabase()
	defer db.Close()

	seedIfEmpty(db)

	return queryTemplates(db)
}

// openSEODatabase opens the gitmap database and migrates tables.
func openSEODatabase() *store.DB {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBOpen, "default", err)
		os.Exit(1)
	}

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBMigrate, err)
		os.Exit(1)
	}

	return db
}

// seedIfEmpty loads seed templates when the table is empty.
func seedIfEmpty(db *store.DB) {
	count, err := db.CountTemplates()
	if err != nil || count > 0 {
		return
	}

	seedFromFile(db, constants.SEOSeedFile)
}

// seedFromFile reads a JSON seed file and inserts templates.
func seedFromFile(db *store.DB, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOSeedRead, path, err)

		return
	}

	var tf templateFile
	if err := json.Unmarshal(data, &tf); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOSeedRead, path, err)

		return
	}

	insertTemplates(db, tf)
}

// insertTemplates inserts all titles and descriptions from a template file.
func insertTemplates(db *store.DB, tf templateFile) {
	total := 0

	for _, t := range tf.Titles {
		if err := db.InsertTemplate(constants.TemplateKindTitle, t); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not insert title template: %v\n", err)
		}
		total++
	}

	for _, d := range tf.Descriptions {
		if err := db.InsertTemplate(constants.TemplateKindDescription, d); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not insert description template: %v\n", err)
		}
		total++
	}

	fmt.Printf(constants.MsgSEOSeeded, total)
}

// queryTemplates fetches title and description templates from the DB.
func queryTemplates(db *store.DB) ([]string, []string) {
	titleRows, err := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if err != nil {
		return nil, nil
	}

	descRows, err := db.ListTemplatesByKind(constants.TemplateKindDescription)
	if err != nil {
		return nil, nil
	}

	titles := extractTemplateText(titleRows)
	descriptions := extractTemplateText(descRows)

	return titles, descriptions
}

// extractTemplateText converts template records to string slices.
func extractTemplateText(rows []store.CommitTemplate) []string {
	result := make([]string, len(rows))
	for i, r := range rows {
		result[i] = r.Template
	}

	return result
}

// generateMessages creates commit messages by pairing random templates.
func generateMessages(titles, descs []string, flags seoWriteFlags) []commitMessage {
	count := flags.maxCommits
	if count == 0 {
		count = len(titles) * len(descs)
	}

	replacer := buildReplacer(flags)
	messages := make([]commitMessage, count)

	for i := 0; i < count; i++ {
		t := titles[rand.Intn(len(titles))]
		d := descs[rand.Intn(len(descs))]
		messages[i] = commitMessage{
			title:       replacer.Replace(t),
			description: replacer.Replace(d),
		}
	}

	return messages
}

// buildReplacer creates a string replacer for all placeholders.
func buildReplacer(flags seoWriteFlags) *strings.Replacer {
	return strings.NewReplacer(
		constants.PlaceholderService, flags.service,
		constants.PlaceholderArea, flags.area,
		constants.PlaceholderURL, flags.url,
		constants.PlaceholderCompany, flags.company,
		constants.PlaceholderPhone, flags.phone,
		constants.PlaceholderEmail, flags.email,
		constants.PlaceholderAddress, flags.address,
	)
}
