package cmdsequence

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
	"github.com/alimtvnetwork/gitmap-v28/cli/repodb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type seqFile struct {
	OriginalPath string
	Dir          string
	BaseName     string
	Seq          int
	HasSeq       bool
	Rest         string
	Time         int64
	NewSeq       int
}

// ParsePrettyFlag pulls --pretty / --no-pretty (and the --color /
// --no-color synonyms) out of args and returns the cleaned slice + the
// resolved render.PrettyMode. Accepted forms:
//
//	--pretty | --color                 → PrettyOn
//	--pretty=true|on|1|yes|y           → PrettyOn
//	--color=true|on|1|yes|y            → PrettyOn
//	--pretty=false|off|0|no|n          → PrettyOff
//	--color=false|off|0|no|n           → PrettyOff
//	--pretty=auto | --color=auto       → PrettyAuto (explicit reset)
//	--no-pretty | --no-color           → PrettyOff
//
// When the same flag is repeated, the **last** occurrence wins (matches
// stdlib flag.Parse semantics) — and "same flag" spans the synonym
// pair, so `--pretty --no-color` resolves to PrettyOff. When neither
// appears, the returned mode is PrettyAuto so callers can rely on
// Decide()'s default ladder.
//
// Unrecognized values fall through to PrettyAuto and the token is left
// in place so the downstream parser can produce a meaningful error.
func ParsePrettyFlag(args []string) ([]string, render.PrettyModeType) {
	mode := render.PrettyAuto
	out := make([]string, 0, len(args))
	for _, arg := range args {
		token, value, hasValue := splitPrettyToken(arg)
		mode = applyPrettyToken(token, value, hasValue, mode, &out, arg)
	}

	return out, mode
}

func parsePinMap(pinStr string, pinMap map[string]int) {
	pairs := strings.Split(pinStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			pinMap[strings.ToLower(strings.TrimSpace(parts[0]))] = val
		}
	}
}

func parseSeqFiles(entries []os.DirEntry, absDir string) []*seqFile {
	var parsedFiles []*seqFile
	re := lazyregex.NumberPrefixRegex.Compiled()
	for _, entry := range entries {
		if sf := buildSeqFile(entry, absDir, re); sf != nil {
			parsedFiles = append(parsedFiles, sf)
		}
	}

	return parsedFiles
}

func partitionFiles(files []*seqFile, pinMap map[string]int) ([]*seqFile, []*seqFile) {
	var pinned, unpinned []*seqFile
	for _, pf := range files {
		baseLower := strings.ToLower(strings.TrimSuffix(pf.Rest, filepath.Ext(pf.Rest)))
		if val, exists := pinMap[baseLower]; exists {
			pf.NewSeq = val
			pinned = append(pinned, pf)
		} else {
			unpinned = append(unpinned, pf)
		}
	}

	return pinned, unpinned
}

func splitPrettyToken(arg string) (token, value string, hasValue bool) {
	isMissingPrefix := !hasPrettyPrefix(arg)
	if isMissingPrefix {
		return arg, "", false
	}

	eq := strings.IndexByte(arg, '=')
	hasEqual := eq >= 0
	if hasEqual {
		return arg[:eq], arg[eq+1:], true
	}

	return arg, "", false
}
func applyPrettyToken(
	token,
	value string,
	hasValue bool,
	mode render.PrettyModeType,
	out *[]string,
	arg string,
) render.PrettyModeType {
	switch token {
	case flagPrettyPositive, flagColorPositive:
		return resolvePositivePretty(value, hasValue, mode, out, arg)
	case flagPrettyNegative, flagColorNegative:
		return render.PrettyOff
	default:
		*out = append(*out, arg)

		return mode
	}
}

const (
	flagPrettyPositive = "--pretty"
	flagPrettyNegative = "--no-pretty"
	flagColorPositive  = "--color"
	flagColorNegative  = "--no-color"
)

// prettyFlagPrefixes lists every token recognized by ParsePrettyFlag.
// Centralized so splitPrettyToken's prefix gate stays in sync with the
// switch in ParsePrettyFlag. `--color` / `--no-color` are accepted as
// synonyms for `--pretty` / `--no-pretty` because in this CLI the
// pretty-markdown pipeline is the only thing that emits ANSI color,
// and `--no-color` is the conventional spelling users reach for first
// (it also mirrors the widely-supported NO_COLOR env convention)
var prettyFlagPrefixes = []string{
	flagPrettyPositive, flagPrettyNegative,
	flagColorPositive, flagColorNegative,
}

// ParsePrettyFlag pulls --pretty / --no-pretty (and the --color /
// --no-color synonyms)
func hasPrettyPrefix(arg string) bool {
	for _, prefix := range prettyFlagPrefixes {
		hasMatch := strings.HasPrefix(arg, prefix)
		if hasMatch {
			return true
		}
	}

	return false
}
func resolvePositivePretty(
	value string,
	hasValue bool,
	current render.PrettyModeType,
	out *[]string,
	original string,
) render.PrettyModeType {
	isMissingValue := !hasValue
	if isMissingValue {
		return render.PrettyOn
	}

	switch strings.ToLower(value) {
	case "1", "t", "true", "on", "yes", "y":
		return render.PrettyOn
	case "0", "f", "false", "off", "no", "n":
		return render.PrettyOff
	case "auto", "":
		return render.PrettyAuto
	}

	*out = append(*out, original)

	return current
}

func runGitMv(src, dst string) bool {
	cmd := exec.Command("git", "mv", src, dst)
	cmd.Dir = filepath.Dir(src)
	err := cmd.Run()

	return err == nil
}

func executeRename(path, newPath, base, newBase string) {
	if runGitMv(path, newPath) {
		fmt.Printf("✅ git mv: %s -> %s\n", base, newBase)

		return
	}

	if err := os.Rename(path, newPath); err == nil {
		fmt.Printf("✅ os.rename: %s -> %s\n", base, newBase)
	} else {
		fmt.Printf("❌ rename failed: %v\n", err)
	}
}

func buildSeqFile(entry os.DirEntry, absDir string, re *regexp.Regexp) *seqFile {
	if entry.IsDir() {
		return nil
	}

	info, err := entry.Info()
	if err != nil {
		return nil
	}

	sf := &seqFile{
		OriginalPath: filepath.Join(absDir, entry.Name()),
		Dir:          absDir,
		BaseName:     entry.Name(),
		Time:         info.ModTime().UnixNano(),
	}

	applyRegexExtract(sf, entry.Name(), re)

	return sf
}

func applyRegexExtract(sf *seqFile, name string, re *regexp.Regexp) {
	match := re.FindStringSubmatch(name)
	if len(match) == 3 {
		sf.HasSeq = true
		sf.Seq, _ = strconv.Atoi(match[1])
		sf.Rest = match[2]
	} else {
		sf.Rest = name
	}
}

func sortKeepOldOrder(unpinned []*seqFile) {
	sort.Slice(unpinned, func(i, j int) bool {
		if unpinned[i].HasSeq && unpinned[j].HasSeq {
			return unpinned[i].Seq < unpinned[j].Seq
		}

		return unpinned[i].HasSeq
	})
}

func buildUsedSeqs(pinned []*seqFile) map[int]bool {
	used := make(map[int]bool)
	for _, pf := range pinned {
		used[pf.NewSeq] = true
	}

	return used
}

func findMaxSeq(files []*seqFile) int {
	maxSeq := 0
	for _, pf := range files {
		if pf.NewSeq > maxSeq {
			maxSeq = pf.NewSeq
		}
	}

	return maxSeq
}

func getRepoDB(ctx context.Context) (*store.DB, *sql.DB, error) {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return nil, nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		mainDB.Close()

		return nil, nil, err
	}

	repos, err := mainDB.FindByPath(cwd)
	if err != nil {
		mainDB.Close()

		return nil, nil, fmt.Errorf("find repo by path %s failed: %w", cwd, err)
	}
	if len(repos) == 0 {
		mainDB.Close()

		return nil, nil, fmt.Errorf("current directory is not a tracked gitmap repository. run 'gitmap scan' first")
	}

	repoDB, err := repodb.OpenRepoDB(ctx, constants.DefaultOutputDir, repos[0].AbsolutePath, repos[0].ID)
	if err != nil {
		mainDB.Close()

		return nil, nil, err
	}

	return mainDB, repoDB, nil
}
