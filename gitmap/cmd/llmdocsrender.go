package cmd

// JSON encoder for `gitmap llm-docs --format=json`.
//
// Migrated off json.MarshalIndent(local jsonDoc struct) onto
// gitmap/stablejson so every key (top-level + nested command groups +
// nested per-command rows) has compile-time-stable order rather than
// reflection-defined order.
//
// The llm-docs output is a single top-level object. The nested
// `commands` array is pre-rendered in compact mode and embedded as
// json.RawMessage. Optional top-level string sections and optional
// per-command `example` fields are conditionally appended so the legacy
// omitempty wire shape is preserved (absent rather than null/empty).
//
// Schema: spec/08-json-schemas/llm-docs.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// llm-docs top-level wire keys. Names + order are the contract.
const (
	llmDocsKeyCommands     = "commands"
	llmDocsKeyArchitecture = "architecture"
	llmDocsKeyFlags        = "flags"
	llmDocsKeyConventions  = "conventions"
	llmDocsKeyStructure    = "structure"
	llmDocsKeyDatabase     = "database"
	llmDocsKeyInstallation = "installation"
	llmDocsKeyPatterns     = "patterns"
)

// per-command-group wire keys.
const (
	llmDocsGroupKeyTitle    = "title"
	llmDocsGroupKeyCommands = "commands"
)

// per-command wire keys.
const (
	llmDocsCmdKeyName        = "name"
	llmDocsCmdKeyAlias       = "alias"
	llmDocsCmdKeyDescription = "description"
	llmDocsCmdKeyExample     = "example"
)

// encodeLLMDocsJSON writes the LLM reference document as stable JSON.
// sections controls which top-level keys are included (nil = all).
func encodeLLMDocsJSON(w io.Writer, sections map[string]bool) error {
	fields := make([]stablejson.Field, 0, 8)

	if wantSection(sections, "commands") {
		commandsRaw, err := renderLLMDocsCommandsRaw()
		if err != nil {
			return err
		}
		fields = append(fields, stablejson.Field{Key: llmDocsKeyCommands, Value: commandsRaw})
	}

	sectionMap := []struct {
		key   string
		write func(*strings.Builder)
		name  string
	}{
		{"architecture", writeLLMArchitecture, llmDocsKeyArchitecture},
		{"flags", writeLLMGlobalFlags, llmDocsKeyFlags},
		{"conventions", writeLLMCodingConventions, llmDocsKeyConventions},
		{"structure", writeLLMProjectStructure, llmDocsKeyStructure},
		{"database", writeLLMDatabase, llmDocsKeyDatabase},
		{"installation", writeLLMInstallation, llmDocsKeyInstallation},
		{"patterns", writeLLMPatterns, llmDocsKeyPatterns},
	}

	for _, sm := range sectionMap {
		if wantSection(sections, sm.key) {
			var sb strings.Builder
			sm.write(&sb)
			s := sb.String()
			if s != "" {
				fields = append(fields, stablejson.Field{Key: sm.name, Value: s})
			}
		}
	}

	return stablejson.WriteObject(w, fields)
}

// renderLLMDocsCommandsRaw pre-renders the command groups array in
// compact mode so it embeds cleanly as a top-level object value.
func renderLLMDocsCommandsRaw() (json.RawMessage, error) {
	var buf bytes.Buffer
	groups := buildCommandGroups()
	if err := stablejson.WriteArrayIndent(&buf, buildLLMDocsGroupItems(groups), ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// buildLLMDocsGroupItems is the single source of (field name, field
// order, value) for command groups.
func buildLLMDocsGroupItems(groups []llmCmdGroup) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(groups))
	for _, g := range groups {
		groupFields := make([]stablejson.Field, 0, 2)
		groupFields = append(groupFields, stablejson.Field{Key: llmDocsGroupKeyTitle, Value: g.title})

		cmdsRaw, _ := renderLLMDocsGroupCommandsRaw(g.commands)
		groupFields = append(groupFields, stablejson.Field{Key: llmDocsGroupKeyCommands, Value: cmdsRaw})

		items = append(items, groupFields)
	}

	return items
}

// renderLLMDocsGroupCommandsRaw pre-renders one group's commands array
// in compact mode.
func renderLLMDocsGroupCommandsRaw(commands []llmCmdEntry) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, buildLLMDocsCommandItems(commands), ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// buildLLMDocsCommandItems is the single source of (field name, field
// order, value) for individual commands within a group.
func buildLLMDocsCommandItems(commands []llmCmdEntry) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(commands))
	for _, c := range commands {
		fields := []stablejson.Field{
			{Key: llmDocsCmdKeyName, Value: c.name},
			{Key: llmDocsCmdKeyAlias, Value: c.alias},
			{Key: llmDocsCmdKeyDescription, Value: c.desc},
		}
		if c.example != "" {
			fields = append(fields, stablejson.Field{Key: llmDocsCmdKeyExample, Value: c.example})
		}
		items = append(items, fields)
	}

	return items
}
