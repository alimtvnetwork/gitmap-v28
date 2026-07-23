// Package cmd — seowritecreate.go scaffolds a sample seo-templates.json.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// sampleTemplateFile is the scaffold written by --create-template.
type sampleTemplateFile struct {
	Titles       []string `json:"titles"`
	Descriptions []string `json:"descriptions"`
	Placeholders []string `json:"placeholders"`
}

// createTemplateFile writes a sample seo-templates.json to the current directory.
func createTemplateFile() {
	sample := buildSampleTemplate()
	data, err := json.MarshalIndent(sample, "", constants.JSONIndent)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOCreateWrite, constants.SEOTemplateOutputFile, err)
		os.Exit(1)
	}

	if err := os.WriteFile(constants.SEOTemplateOutputFile, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOCreateWrite, constants.SEOTemplateOutputFile, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgSEOCreated, constants.SEOTemplateOutputFile)
}

// buildSampleTemplate returns a starter template with examples.
func buildSampleTemplate() sampleTemplateFile {
	return sampleTemplateFile{
		Titles: []string{
			"Top {service} in {area} — {url}",
			"{company}: Trusted {service} Provider in {area}",
			"Best {service} Near {area} | Visit {url}",
		},
		Descriptions: []string{
			"Looking for reliable {service} in {area}? {company} provides top-rated solutions. Visit {url} or call {phone}.",
			"{company} is the leading {service} provider serving {area}. Learn more at {url}. Contact us at {email}.",
		},
		Placeholders: []string{
			constants.PlaceholderService,
			constants.PlaceholderArea,
			constants.PlaceholderURL,
			constants.PlaceholderCompany,
			constants.PlaceholderPhone,
			constants.PlaceholderEmail,
			constants.PlaceholderAddress,
		},
	}
}
