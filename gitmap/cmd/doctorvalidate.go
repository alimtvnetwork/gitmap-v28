package cmd

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// checkConfigFile verifies config.json exists and is valid JSON.
func checkConfigFile() int {
	configPath := resolveConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		printWarn(constants.DoctorConfigMissing)

		return 0
	}

	return validateConfigJSON(data, configPath)
}

// resolveConfigPath returns the config file path from RepoPath or default.
func resolveConfigPath() string {
	if len(constants.RepoPath) > 0 {
		return filepath.Join(constants.RepoPath, constants.GitMapSubdir, constants.DefaultConfigPath)
	}

	return constants.DefaultConfigPath
}

// validateConfigJSON checks if config data is valid JSON.
func validateConfigJSON(data []byte, path string) int {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		printIssue(constants.DoctorConfigInvalid, err.Error())

		return 1
	}

	printOK(constants.DoctorConfigOKFmt, path)

	return 0
}

// checkDatabase verifies the database can be opened and migrated.
func checkDatabase() int {
	db, err := store.OpenDefault()
	if err != nil {
		printIssue(constants.DoctorDBOpenFail, err.Error())

		return 1
	}
	defer db.Close()

	err = db.Migrate()
	if err != nil {
		printIssue(constants.DoctorDBMigrateFail, err.Error())

		return 1
	}

	printOK(constants.DoctorDBOK, store.DefaultDBPath())

	return 0
}

// checkLockFile reports if a stale lock file exists.
func checkLockFile() int {
	dir := store.BinaryDataDir()
	lockPath := filepath.Join(dir, constants.LockFileName)

	_, err := os.Stat(lockPath)
	if err != nil {
		printOK(constants.DoctorLockNone)

		return 0
	}

	printWarn(constants.DoctorLockExists)

	return 0
}

// checkNetwork reports basic connectivity status.
func checkNetwork() int {
	conn, err := net.DialTimeout(
		constants.NetworkProto,
		constants.NetworkCheckHost,
		time.Duration(constants.NetworkTimeoutSec)*time.Second,
	)
	if err != nil {
		printWarn(constants.DoctorNetworkOffline)

		return 0
	}

	conn.Close()
	printOK(constants.DoctorNetworkOK)

	return 0
}
