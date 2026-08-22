package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

	"path/filepath"
	"runtime"
)

func cmdPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic(apperror.New("resolve cmd package dir", "ERR_TEST_PANIC", nil))
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
