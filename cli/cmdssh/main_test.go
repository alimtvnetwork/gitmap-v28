package cmdssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func setupHermeticTestDB(tempDir string) {
	testDBPath := filepath.Join(tempDir, "gitmap.db")
	testDB, openErr := store.OpenAt(testDBPath)
	if openErr != nil {
		return
	}
	defer testDB.Close()
	_ = store.EnsureSSHTables(testDB.SQL())
	_ = store.EnsureSSHHostsTable(testDB.SQL())
	_ = store.EnsureSSHHistoryTable(testDB.SQL())
	openSSHDBFunc = func() (*store.DB, error) { return store.OpenAt(testDBPath) }
	openSSHDB = func() (*store.DB, error) { return store.OpenAt(testDBPath) }
}

func runTestMain(m *testing.M) int {
	tempDir, err := os.MkdirTemp("", "cmdssh_hermetic_*")
	if err == nil {
		defer os.RemoveAll(tempDir)
		store.SetBinaryDataDirForTesting(tempDir)
		setupHermeticTestDB(tempDir)
	}
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(runTestMain(m))
}
