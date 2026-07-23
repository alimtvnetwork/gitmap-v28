// Package cmd — `gitmap downloader-config [path]` (shorthand: `dc`).
//
// Slice 1 of the downloader feature. Reads / validates / persists the
// downloader Seedable-Config. Two modes:
//
//  1. Path supplied  → load JSON from disk, validate, save to Setting DB.
//  2. No path        → interactive prompt, pre-populated from current
//     DB values (or downloaderconfig.Defaults() if none).
//
// All actual download / install logic ships in Slice 2 (aria2c installer
// + engine) and Slice 3 (download / download-unzip commands). This command
// exists so users can pre-tune the config before those slices land.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/downloaderconfig"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runDownloaderConfig is the dispatch entrypoint.
func runDownloaderConfig(args []string) {
	checkHelp(constants.CmdDownloaderConfig, args)
	fmt.Fprintf(os.Stderr, constants.MsgDownloaderConfigBanner+"\n", constants.Version)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Migrate: %v\n", err)
		os.Exit(1)
	}

	doc, source := loadDocOrPrompt(db, args)
	if err := db.SetDownloaderConfig(doc); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not save downloader config: %v\n", err)
		os.Exit(1)
	}

	if source != "" {
		fmt.Fprintf(os.Stderr, constants.MsgDownloaderConfigLoaded+"\n", source)
	}
	fmt.Fprintf(os.Stderr, constants.MsgDownloaderConfigSaved+"\n", constants.SettingDownloaderConfig)
	fmt.Fprintf(os.Stderr, constants.MsgDownloaderConfigDBVersion+"\n", doc.DatabaseVersion.LastKnownVersion)
}

// loadDocOrPrompt resolves the final Document either from a JSON file
// argument or via the interactive prompt. Returns the source description
// (for logging) which is empty for interactive mode.
func loadDocOrPrompt(db *store.DB, args []string) (downloaderconfig.Document, string) {
	if len(args) > 0 {
		path := args[0]
		doc, err := downloaderconfig.LoadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
			os.Exit(1)
		}

		return doc, path
	}

	current, ok := db.GetDownloaderConfig()
	if !ok {
		current = downloaderconfig.Defaults()
	}
	current.DatabaseVersion.LastKnownVersion = constants.Version

	return promptDownloaderConfig(current), ""
}

// promptDownloaderConfig walks the user through every field, showing the
// current value as the default. <Enter> on any prompt keeps the default.
func promptDownloaderConfig(current downloaderconfig.Document) downloaderconfig.Document {
	fmt.Fprintln(os.Stderr, constants.MsgDownloaderConfigPromptHeader)
	reader := bufio.NewReader(os.Stdin)
	dc := current.DownloaderConfig

	dc.PreferredDownloader = promptString(reader, "PreferredDownloader", dc.PreferredDownloader)
	dc.FallbackDownloader = promptString(reader, "FallbackDownloader", dc.FallbackDownloader)
	dc.ParallelDownloads = promptInt(reader, "ParallelDownloads", dc.ParallelDownloads)
	dc.SplitConnections = promptInt(reader, "SplitConnections", dc.SplitConnections)
	dc.DefaultSplitSize = promptString(reader, "DefaultSplitSize", dc.DefaultSplitSize)
	dc.LargeFileSplitSize = promptString(reader, "LargeFileSplitSize", dc.LargeFileSplitSize)
	dc.LargeFileThreshold = promptString(reader, "LargeFileThreshold", dc.LargeFileThreshold)
	dc.TinyFileThreshold = promptString(reader, "TinyFileThreshold", dc.TinyFileThreshold)
	dc.TinyFileSplitSize = promptString(reader, "TinyFileSplitSize", dc.TinyFileSplitSize)
	dc.TinyFileSplits = promptInt(reader, "TinyFileSplits", dc.TinyFileSplits)
	dc.AllowFallback = promptBool(reader, "AllowFallback", dc.AllowFallback)
	dc.OverwriteUserConfig = promptBool(reader, "OverwriteUserConfig", dc.OverwriteUserConfig)

	doc := downloaderconfig.Document{
		DownloaderConfig: dc,
		DatabaseVersion:  downloaderconfig.DatabaseVersion{LastKnownVersion: constants.Version},
	}

	if err := downloaderconfig.Validate(doc); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	return doc
}

// promptString reads a line and returns the trimmed input or the default.
func promptString(reader *bufio.Reader, label, def string) string {
	fmt.Fprintf(os.Stderr, "    %s [%s]: ", label, def)
	line, _ := reader.ReadString('\n')
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return def
	}

	return trimmed
}

// promptInt is promptString + strconv with re-prompt on parse error.
func promptInt(reader *bufio.Reader, label string, def int) int {
	for {
		raw := promptString(reader, label, strconv.Itoa(def))
		n, err := strconv.Atoi(raw)
		if err == nil {
			return n
		}
		fmt.Fprintf(os.Stderr, "      ⚠ %q is not a valid integer — try again\n", raw)
	}
}

// promptBool accepts y/yes/true/1 (case-insensitive) as true.
func promptBool(reader *bufio.Reader, label string, def bool) bool {
	defStr := "false"
	if def {
		defStr = "true"
	}
	raw := strings.ToLower(promptString(reader, label, defStr))
	switch raw {
	case "y", "yes", "true", "1":
		return true
	case "n", "no", "false", "0":
		return false
	}

	return def
}
