package cmd

import (
	"path/filepath"
	"runtime"
)

func cmdPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed resolving cmd package dir")
	}

	return filepath.Dir(file)
}

func cmdPackagePath(parts ...string) string {
	return filepath.Join(append([]string{cmdPackageDir()}, parts...)...)
}

func resolveGoldenPath(name string) string {
	return cmdPackagePath(goldenDir, name)
}

func resolveSchemaDir() string {
	if filepath.IsAbs(schemaDir) {
		return schemaDir
	}

	return cmdPackagePath(schemaDir)
}
