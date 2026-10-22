package heatmap

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TrackFileTouch records a file access operation into the root heatmap database.
func TrackFileTouch(filePath string, touchType store.FileTouchType) error {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(filePath)))
	isValidPath := len(cleanPath) > 0 && cleanPath != "."
	if !isValidPath {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.cwd")
	}

	return recordFileTouchAtRoot(cwd, cleanPath, touchType)
}

func recordFileTouchAtRoot(cwd, cleanPath string, touchType store.FileTouchType) error {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		ActionType: string(touchType),
		FilePath:   cleanPath,
	})

	return logErr
}
