package cmddb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DbMountTarget describes a mounted SQLite database in the REST engine.
type DbMountTarget struct {
	Alias         string `json:"alias"`
	DatabasePath  string `json:"databasePath"`
	DatabaseType  string `json:"databaseType"`
	FileSizeBytes int64  `json:"fileSizeBytes"`
	TableCount    int    `json:"tableCount"`
	IsWritable    bool   `json:"isWritable"`
}

// ResolveMountTargets resolves a comma-separated list of database aliases or file paths.
func ResolveMountTargets(rawTargets string) ([]DbMountTarget, error) {
	trimmed := strings.TrimSpace(rawTargets)
	if trimmed == "" {
		trimmed = "master"
	}

	tokens := strings.Split(trimmed, ",")
	seenAliases := make(map[string]int, len(tokens))
	targets := make([]DbMountTarget, 0, len(tokens))

	for _, token := range tokens {
		cleanToken := strings.TrimSpace(token)
		if cleanToken == "" {
			continue
		}
		target, err := resolveSingleMountTarget(cleanToken, seenAliases)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid database targets specified")
	}
	return targets, nil
}

func resolveSingleMountTarget(token string, seenAliases map[string]int) (DbMountTarget, error) {
	resolvedPath, alias, dbType, err := resolveTargetLocation(token)
	if err != nil {
		return DbMountTarget{}, err
	}

	alias = disambiguateAlias(alias, seenAliases)
	size := getFileSize(resolvedPath)
	tableCount := countTablesInDb(resolvedPath)
	isWritable := checkFileWritable(resolvedPath)

	return DbMountTarget{
		Alias:         alias,
		DatabasePath:  resolvedPath,
		DatabaseType:  dbType,
		FileSizeBytes: size,
		TableCount:    tableCount,
		IsWritable:    isWritable,
	}, nil
}

func resolveTargetLocation(token string) (string, string, string, error) {
	if path, ok := resolveSystemDbAlias(token); ok {
		cleanAlias := strings.ToLower(token)
		if cleanAlias == "" || cleanAlias == "gitmap" {
			cleanAlias = "master"
		}
		return path, cleanAlias, inferDatabaseType(cleanAlias, path), nil
	}

	if path, ok := lookupAliasInRegistry(token); ok {
		cleanAlias := strings.ToLower(token)
		return path, cleanAlias, inferDatabaseType(cleanAlias, path), nil
	}

	return resolveFilesystemTarget(token)
}

func resolveFilesystemTarget(token string) (string, string, string, error) {
	cleanPath := filepath.Clean(token)
	stat, err := os.Stat(cleanPath)
	if err != nil {
		return "", "", "", fmt.Errorf("database target '%s' not found as registered alias or file path", token)
	}
	if stat.IsDir() {
		return "", "", "", fmt.Errorf("target '%s' is a directory, expected SQLite database file", token)
	}

	derivedAlias := deriveAliasFromPath(cleanPath)
	return cleanPath, derivedAlias, "custom", nil
}

func deriveAliasFromPath(path string) string {
	base := filepath.Base(path)
	lower := strings.ToLower(base)
	for _, ext := range []string{".sqlite3", ".sqlite", ".db"} {
		if strings.HasSuffix(lower, ext) {
			base = base[:len(base)-len(ext)]
			break
		}
	}
	clean := strings.TrimSpace(base)
	if clean == "" || clean == "." {
		return "database"
	}
	return strings.ToLower(clean)
}

func disambiguateAlias(alias string, seenAliases map[string]int) string {
	lower := strings.ToLower(alias)
	seenAliases[lower]++
	count := seenAliases[lower]
	if count > 1 {
		return fmt.Sprintf("%s-%d", alias, count)
	}
	return alias
}

func checkFileWritable(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}
