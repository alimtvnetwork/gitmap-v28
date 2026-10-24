package cmdschedule

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	_ "modernc.org/sqlite"
)

func writeScheduleExportZIP(bundles []ScheduleExportBundle, filePath string) error {
	if filePath == "" {
		filePath = "schedules_export.zip"
	}

	_ = os.MkdirAll(filepath.Dir(filePath), constants.DirPermission)
	outFile, err := os.Create(filePath)
	if err != nil {
		return apperror.WrapSimple(err, "create zip file")
	}

	defer outFile.Close()
	zw := zip.NewWriter(outFile)
	defer zw.Close()

	return addBundlesToZIP(zw, bundles, filePath)
}

func addBundlesToZIP(zw *zip.Writer, bundles []ScheduleExportBundle, filePath string) error {
	for _, b := range bundles {
		raw, _ := json.MarshalIndent(b, "", constants.JSONIndent)
		f, err := zw.Create(b.Task.Slug + ".json")
		if err == nil {
			_, _ = f.Write(raw)
		}
	}

	printExportSuccessBanner(filePath, len(bundles), "ZIP")

	return nil
}

func printExportSuccessBanner(filePath string, count int, format string) {
	fmt.Printf("\n  \033[1;92m✔ Exported %d schedule(s)\033[0m to %s (\033[1m%s\033[0m)\n\n", count, format, filePath)
}
