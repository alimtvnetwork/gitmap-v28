// Package cmd — agy_settings.go handles Antigravity settings export and import.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var agySettingsCmd = &cobra.Command{
	Use:     "settings",
	Aliases: []string{"setting", "config"},
	Short:   "Export and import Antigravity configuration settings",
}

var agySettingsExportCmd = &cobra.Command{
	Use:   "export [file.json]",
	Short: "Export Antigravity settings to a JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		outPath := "antigravity_settings.json"
		if len(args) > 0 {
			outPath = args[0]
		}

		return executeAgySettingsExport(outPath)
	},
}

var agySettingsImportCmd = &cobra.Command{
	Use:   "import [file.json]",
	Short: "Import Antigravity settings from a JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap agy settings import <file.json>")
		}

		return executeAgySettingsImport(args[0])
	},
}

func getAgyConfigPath() (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "home dir")
	}

	return filepath.Join(homeDir, ".gemini", "antigravity", "config.json"), nil
}

func executeAgySettingsExport(outPath string) error {
	cfgPath, pathErr := getAgyConfigPath()
	if pathErr != nil {
		return pathErr.Unwrap()
	}

	payload := map[string]interface{}{
		"version":    constants.Version,
		"exportedAt": time.Now().UTC().Format(time.RFC3339),
		"settings":   make(map[string]interface{}),
	}

	data, err := os.ReadFile(cfgPath)
	if err == nil {
		loadRawSettings(data, payload)
	}

	outBytes, mErr := json.MarshalIndent(payload, "", "  ")
	if mErr != nil {
		return mErr
	}

	if err := os.WriteFile(outPath, outBytes, constants.FilePermission); err != nil {
		return err
	}

	fmt.Printf("%s Exported Antigravity settings to %s\n", constants.ColorGreen+"✓"+constants.ColorReset, outPath)
	return nil
}

func executeAgySettingsImport(inPath string) error {
	data, readErr := os.ReadFile(inPath)
	if readErr != nil {
		return readErr
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	cfgPath, pathErr := getAgyConfigPath()
	if pathErr != nil {
		return pathErr.Unwrap()
	}

	if mkErr := os.MkdirAll(filepath.Dir(cfgPath), constants.DirPermission); mkErr != nil {
		return mkErr
	}

	if err := os.WriteFile(cfgPath, data, constants.FilePermission); err != nil {
		return err
	}

	fmt.Printf("%s Successfully imported Antigravity settings from %s\n", constants.ColorGreen+"✓"+constants.ColorReset, inPath)
	return nil
}

func loadRawSettings(data []byte, payload map[string]interface{}) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err == nil {
		payload["settings"] = raw
	}
}

func initAgySettings() {
	agySettingsCmd.AddCommand(agySettingsExportCmd)
	agySettingsCmd.AddCommand(agySettingsImportCmd)
}
