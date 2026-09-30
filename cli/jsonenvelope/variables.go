package jsonenvelope

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ExpandVariables interpolates ${varName} and $variables.varName from vars map into target content.
func ExpandVariables(content []byte, vars map[string]any) []byte {
	if len(content) == 0 || len(vars) == 0 {
		return content
	}

	result := content
	for k, v := range vars {
		result = replaceSingleVariable(result, k, v)
	}
	return result
}

func replaceSingleVariable(content []byte, key string, val any) []byte {
	valStr := formatVariableValue(val)
	if valStr == "" {
		return content
	}

	patterns := buildVariablePatterns(key)
	escapedVal := escapeForJSON(valStr)

	res := content
	for _, p := range patterns {
		res = bytes.ReplaceAll(res, []byte(p), []byte(escapedVal))
	}
	return res
}

func formatVariableValue(val any) string {
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func buildVariablePatterns(key string) []string {
	return []string{
		"${" + key + "}",
		"${variables." + key + "}",
		"$variables." + key,
	}
}

func escapeForJSON(s string) string {
	b, err := json.Marshal(s)
	if err != nil || len(b) < 2 {
		return s
	}
	// Strip enclosing quotes to get inner escaped JSON string
	return string(b[1 : len(b)-1])
}

// MergeVariables combines top-level and work-directory variables into a single map.
func MergeVariables(topVars map[string]any, workDirCfg *WorkDirectoryConfig) map[string]any {
	out := make(map[string]any)
	for k, v := range topVars {
		out[k] = v
	}
	if workDirCfg != nil && workDirCfg.Variables != nil {
		for k, v := range workDirCfg.Variables {
			out[k] = v
		}
	}
	return out
}

// ResolveStringVariable performs direct expansion on a single Go string.
func ResolveStringVariable(input string, vars map[string]any) string {
	if !strings.Contains(input, "$") || len(vars) == 0 {
		return input
	}
	res := input
	for k, v := range vars {
		valStr := formatVariableValue(v)
		patterns := buildVariablePatterns(k)
		for _, p := range patterns {
			res = strings.ReplaceAll(res, p, valStr)
		}
	}
	return res
}
