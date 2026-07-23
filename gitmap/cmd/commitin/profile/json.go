package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// jsonShape mirrors spec §5.2 key order exactly. Field tags drive both
// strict decode (DisallowUnknownFields) and stable encode order.
type jsonShape struct {
	Name             string        `json:"Name"`
	SchemaVersion    int           `json:"SchemaVersion"`
	SourceRepoPath   string        `json:"SourceRepoPath"`
	IsDefault        bool          `json:"IsDefault"`
	ConflictMode     string        `json:"ConflictMode"`
	Author           *jsonAuthor   `json:"Author,omitempty"`
	Exclusions       []jsonKV      `json:"Exclusions"`
	MessageRules     []jsonKV      `json:"MessageRules"`
	MessagePrefix    []string      `json:"MessagePrefix"`
	MessageSuffix    []string      `json:"MessageSuffix"`
	TitlePrefix      string        `json:"TitlePrefix"`
	TitleSuffix      string        `json:"TitleSuffix"`
	OverrideMessages []string      `json:"OverrideMessages"`
	OverrideOnlyWeak bool          `json:"OverrideOnlyWeak"`
	WeakWords        []string      `json:"WeakWords"`
	FunctionIntel    jsonFuncIntel `json:"FunctionIntel"`
}

type jsonAuthor struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

type jsonKV struct {
	Kind  string `json:"Kind"`
	Value string `json:"Value"`
}

type jsonFuncIntel struct {
	IsEnabled bool     `json:"IsEnabled"`
	Languages []string `json:"Languages"`
}

// Decode parses a profile JSON byte slice into a Profile (strict).
func Decode(raw []byte) (*Profile, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var s jsonShape
	if err := dec.Decode(&s); err != nil {
		return nil, &LoadError{Reason: "invalid json", Cause: err}
	}
	if s.SchemaVersion != CurrentSchemaVersion {
		return nil, &LoadError{Reason: fmt.Sprintf("unsupported SchemaVersion %d", s.SchemaVersion)}
	}
	return shapeToProfile(&s), nil
}

// Encode serializes a Profile to canonical JSON (2-space indent, fixed
// key order, trailing newline) for byte-stable diffs.
func Encode(p *Profile) ([]byte, error) {
	s := profileToShape(p)
	out, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

func shapeToProfile(s *jsonShape) *Profile {
	p := &Profile{
		Name:             s.Name,
		SchemaVersion:    s.SchemaVersion,
		SourceRepoPath:   s.SourceRepoPath,
		IsDefault:        s.IsDefault,
		ConflictMode:     s.ConflictMode,
		MessagePrefix:    s.MessagePrefix,
		MessageSuffix:    s.MessageSuffix,
		TitlePrefix:      s.TitlePrefix,
		TitleSuffix:      s.TitleSuffix,
		OverrideMessages: s.OverrideMessages,
		OverrideOnlyWeak: s.OverrideOnlyWeak,
		WeakWords:        s.WeakWords,
		FunctionIntel:    FunctionIntel{IsEnabled: s.FunctionIntel.IsEnabled, Languages: s.FunctionIntel.Languages},
	}
	if s.Author != nil {
		p.Author = &Author{Name: s.Author.Name, Email: s.Author.Email}
	}
	for _, e := range s.Exclusions {
		p.Exclusions = append(p.Exclusions, Exclusion(e))
	}
	for _, r := range s.MessageRules {
		p.MessageRules = append(p.MessageRules, MessageRule(r))
	}
	return p
}

func profileToShape(p *Profile) *jsonShape {
	s := &jsonShape{
		Name:             p.Name,
		SchemaVersion:    p.SchemaVersion,
		SourceRepoPath:   p.SourceRepoPath,
		IsDefault:        p.IsDefault,
		ConflictMode:     p.ConflictMode,
		MessagePrefix:    p.MessagePrefix,
		MessageSuffix:    p.MessageSuffix,
		TitlePrefix:      p.TitlePrefix,
		TitleSuffix:      p.TitleSuffix,
		OverrideMessages: p.OverrideMessages,
		OverrideOnlyWeak: p.OverrideOnlyWeak,
		WeakWords:        p.WeakWords,
		FunctionIntel:    jsonFuncIntel{IsEnabled: p.FunctionIntel.IsEnabled, Languages: p.FunctionIntel.Languages},
	}
	if p.Author != nil {
		s.Author = &jsonAuthor{Name: p.Author.Name, Email: p.Author.Email}
	}
	if s.SchemaVersion == 0 {
		s.SchemaVersion = CurrentSchemaVersion
	}
	for _, e := range p.Exclusions {
		s.Exclusions = append(s.Exclusions, jsonKV(e))
	}
	for _, r := range p.MessageRules {
		s.MessageRules = append(s.MessageRules, jsonKV(r))
	}
	if s.Exclusions == nil {
		s.Exclusions = []jsonKV{}
	}
	if s.MessageRules == nil {
		s.MessageRules = []jsonKV{}
	}
	if s.MessagePrefix == nil {
		s.MessagePrefix = []string{}
	}
	if s.MessageSuffix == nil {
		s.MessageSuffix = []string{}
	}
	if s.OverrideMessages == nil {
		s.OverrideMessages = []string{}
	}
	if s.WeakWords == nil {
		s.WeakWords = []string{}
	}
	if s.FunctionIntel.Languages == nil {
		s.FunctionIntel.Languages = []string{}
	}
	return s
}
