package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/formatter"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// writeAllOutputs writes terminal, CSV, JSON, text, folder structure, and clone scripts.
func writeAllOutputs(records []model.ScanRecord, outputDir, outFile string, quiet bool) {
	writeTerminalOutput(records, outputDir, quiet)
	writeCSVOutput(records, outputDir, outFile)
	writeJSONOutput(records, outputDir)
	writeTextOutput(records, outputDir)
	writeFolderStructure(records, outputDir)
	writeCloneScript(records, outputDir)
	writeDirectCloneScript(records, outputDir)
	writeDirectCloneSSHScript(records, outputDir)
	writeDesktopScript(records, outputDir)
}

// writeTerminalOutput renders records to stdout.
func writeTerminalOutput(records []model.ScanRecord, outputDir string, quiet bool) {
	err := formatter.Terminal(os.Stdout, records, outputDir, quiet)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrOutputFailed, err)
	}
}

// writeCSVOutput writes records to a CSV file.
func writeCSVOutput(records []model.ScanRecord, outputDir, outFile string) {
	path := resolveOutFile(outFile, outputDir, constants.DefaultCSVFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteCSV(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write CSV to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgCSVWritten, filepath.Base(path))
}

// writeJSONOutput writes records to a JSON file.
func writeJSONOutput(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultJSONFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteJSON(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write JSON to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgJSONWritten, filepath.Base(path))
}

// writeTextOutput writes records as plain text clone commands.
func writeTextOutput(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultTextFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteText(file, records)
	fmt.Printf(constants.MsgTextWritten, filepath.Base(path))
}

// writeFolderStructure writes a Markdown file showing the repo tree.
func writeFolderStructure(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultStructureFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteStructure(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write structure to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgStructureWritten, filepath.Base(path))
}

// writeCloneScript writes a PowerShell clone script.
func writeCloneScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultCloneScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteCloneScript(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write clone script to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgCloneScript, filepath.Base(path))
}

// writeDirectCloneScript writes a plain PS1 with one git clone per line.
func writeDirectCloneScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDirectCloneScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteDirectCloneScript(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write direct clone script to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgDirectClone, filepath.Base(path))
}

// writeDirectCloneSSHScript writes a plain SSH PS1 with one git clone per line.
func writeDirectCloneSSHScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDirectCloneSSHScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteDirectCloneSSHScript(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write SSH clone script to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgDirectCloneSSH, filepath.Base(path))
}

// writeDesktopScript writes a PowerShell script to register repos with GitHub Desktop.
func writeDesktopScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDesktopScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := formatter.WriteDesktopScript(file, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write desktop script to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgDesktopScript, filepath.Base(path))
}

// resolveOutFile determines the output file path.
func resolveOutFile(outFile, outputDir, defaultName string) string {
	if len(outFile) > 0 {
		return outFile
	}

	return filepath.Join(outputDir, defaultName)
}

// createOutputFile ensures the directory exists and creates the file.
func createOutputFile(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCreateDir, filepath.Dir(path), err)

		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCreateFile, path, err)

		return nil, err
	}

	return file, nil
}
