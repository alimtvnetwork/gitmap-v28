package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	varMu         sync.RWMutex
	varRegex      = regexp.MustCompile(`\$\{([a-zA-Z0-9_:-]+)\}|\$([a-zA-Z0-9_:-]+)`)
	cachedVarData *VariableStore
)

// VariableStore stores global and subcommand-scoped variables.
type VariableStore struct {
	Global map[string]string            `json:"global"`
	Scopes map[string]map[string]string `json:"scopes"`
}

func getVariablesFilePath() string {
	if localPath, isLocal := resolveLocalVariablesPath(); isLocal {
		return localPath
	}

	home, err := os.UserHomeDir()
	if err == nil {
		globalDir := filepath.Join(home, ".gitmap")
		_ = os.MkdirAll(globalDir, 0755)

		return filepath.Join(globalDir, "variables.json")
	}

	tempDir := filepath.Join(os.TempDir(), "gitmap")
	_ = os.MkdirAll(tempDir, 0755)

	return filepath.Join(tempDir, "variables.json")
}

func resolveLocalVariablesPath() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	localPath := filepath.Join(cwd, ".gitmap", "variables.json")
	if _, statErr := os.Stat(filepath.Dir(localPath)); statErr == nil {
		return localPath, true
	}
	return "", false
}

func loadVariableStore() *VariableStore {
	varMu.Lock()
	defer varMu.Unlock()

	if cachedVarData != nil {
		return cachedVarData
	}

	filePath := getVariablesFilePath()
	store := &VariableStore{
		Global: make(map[string]string),
		Scopes: make(map[string]map[string]string),
	}

	data, err := os.ReadFile(filePath)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, store)
	}

	cachedVarData = store

	return store
}

func saveVariableStore(store *VariableStore) error {
	filePath := getVariablesFilePath()
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// SetVariable sets a variable in either the global scope or a subcommand scope.
func SetVariable(scope, key, value string) error {
	store := loadVariableStore()
	varMu.Lock()
	defer varMu.Unlock()

	normScope := strings.ToLower(strings.TrimSpace(scope))
	normKey := strings.TrimPrefix(strings.TrimSpace(key), "$")

	if normScope == "" || normScope == "global" {
		store.Global[normKey] = value
	} else {
		ensureScopedMap(store, normScope)
		store.Scopes[normScope][normKey] = value
	}

	return saveVariableStore(store)
}

func ensureScopedMap(store *VariableStore, scope string) {
	if store.Scopes[scope] == nil {
		store.Scopes[scope] = make(map[string]string)
	}
}

// GetVariable retrieves a variable, searching the specific subcommand scope, then global scope, then OS env.
func GetVariable(scope, key string) (string, bool) {
	store := loadVariableStore()
	varMu.RLock()
	defer varMu.RUnlock()

	normScope := strings.ToLower(strings.TrimSpace(scope))
	normKey := strings.TrimPrefix(strings.TrimSpace(key), "$")

	if val, hasVal := getScopedVariable(store, normScope, normKey); hasVal {
		return val, true
	}

	if val, hasVal := store.Global[normKey]; hasVal {
		return val, true
	}

	if envVal, hasEnv := os.LookupEnv(normKey); hasEnv {
		return envVal, true
	}

	return "", false
}

func getScopedVariable(store *VariableStore, scope, key string) (string, bool) {
	if scope == "" || scope == "global" {
		return "", false
	}
	scopedMap, hasScope := store.Scopes[scope]
	if !hasScope {
		return "", false
	}
	val, hasVal := scopedMap[key]
	return val, hasVal
}

// DeleteVariable removes a variable from global or scoped storage.
func DeleteVariable(scope, key string) error {
	store := loadVariableStore()
	varMu.Lock()
	defer varMu.Unlock()

	normScope := strings.ToLower(strings.TrimSpace(scope))
	normKey := strings.TrimPrefix(strings.TrimSpace(key), "$")

	if normScope == "" || normScope == "global" {
		delete(store.Global, normKey)
	} else if scopedMap, hasScope := store.Scopes[normScope]; hasScope {
		delete(scopedMap, normKey)
	}

	return saveVariableStore(store)
}

// ListVariables returns copies of the global and scoped variables.
func ListVariables() (map[string]string, map[string]map[string]string) {
	store := loadVariableStore()
	varMu.RLock()
	defer varMu.RUnlock()

	glob := make(map[string]string)
	for k, v := range store.Global {
		glob[k] = v
	}

	scoped := make(map[string]map[string]string)
	for sc, m := range store.Scopes {
		inner := make(map[string]string)
		for k, v := range m {
			inner[k] = v
		}
		scoped[sc] = inner
	}

	return glob, scoped
}

// ExpandVariables replaces $VAR or ${VAR} in input with resolved variable values.
func ExpandVariables(input, scope string) string {
	if !strings.Contains(input, "$") {
		return input
	}

	return varRegex.ReplaceAllStringFunc(input, func(m string) string {
		varName := m
		if strings.HasPrefix(m, "${") && strings.HasSuffix(m, "}") {
			varName = m[2 : len(m)-1]
		} else if strings.HasPrefix(m, "$") {
			varName = m[1:]
		}

		if val, found := GetVariable(scope, varName); found {
			return val
		}

		return m
	})
}

// ExportVariablesToEnv exports all global and scoped variables into the OS environment.
func ExportVariablesToEnv() error {
	store := loadVariableStore()
	varMu.RLock()
	defer varMu.RUnlock()

	for k, v := range store.Global {
		_ = os.Setenv(k, v)
	}

	for scope, m := range store.Scopes {
		for k, v := range m {
			scopedKey := strings.ToUpper(scope) + "_" + k
			_ = os.Setenv(scopedKey, v)
		}
	}

	return nil
}
