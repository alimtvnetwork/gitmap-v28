package cmdpull

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func checkCWDCacheOrProceed(raw bool) (bool, error) {
	if raw {
		return false, nil
	}
	served, err := tryServeCWDPullCache()
	if !served {
		return false, nil
	}
	handleCWDIgnoreChecks()

	return true, err
}

func tryServeCWDPullCache() (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, nil
	}

	ttl := resolvePullCacheTTL()
	if ttl <= 0 {
		return false, nil
	}

	pullDB, err := store.OpenPullSplitDB()
	if err != nil {
		return false, nil
	}
	defer pullDB.Close()

	history, err := queryRepoPullHistoryWithFallback(pullDB, cwd)
	if err != nil || len(history) == 0 {
		return false, nil
	}

	lastRun := history[0]
	parsedTime, err := parsePullTimestamp(lastRun.CreatedAt)
	if err != nil {
		return false, nil
	}

	elapsed := time.Since(parsedTime)
	if elapsed < 0 || elapsed > time.Duration(ttl)*time.Second {
		return false, nil
	}

	renderCachedPullStatus(cwd, lastRun, elapsed)
	return true, nil
}

func renderCachedPullStatus(cwd string, lastRun store.PullRepoRunRecord, elapsed time.Duration) {
	diag := gitutil.InspectDirtyState(cwd)
	repo := resolveSummaryRepoName(nil, cwd)
	branch := resolveSummaryBranch(nil, cwd)
	shortSHA := resolveSummaryShortSHA(nil, cwd)
	if shortSHA == "head" && lastRun.LastCommitSha != "" {
		shortSHA = ShortenSHA(lastRun.LastCommitSha)
	}
	elapsedSec := int(elapsed.Seconds())
	if elapsedSec < 0 {
		elapsedSec = 0
	}

	if diag.IsDirty {
		renderDirtyCachedStatus(repo, elapsedSec, shortSHA, branch, diag)
		return
	}

	fmt.Printf("  ⚡ %s: cached (checked %ds ago) - up-to-date (%s) [%s]\n",
		repo, elapsedSec, shortSHA, branch)
}

func renderDirtyCachedStatus(repo string, elapsedSec int, shortSHA, branch string, diag gitutil.DirtyDiagnosis) {
	reason := diag.SummaryReason
	if reason == "" {
		reason = "uncommitted changes"
	}
	fmt.Printf("  ⚡ %s: cached (checked %ds ago) - up-to-date (%s) [%s] (dirty: %s)\n",
		repo, elapsedSec, shortSHA, branch, reason)
	renderInlineDirtyFiles(diag)
}

func resolvePullCacheTTL() int {
	db, err := store.OpenDefault()
	if err != nil {
		return constants.DefaultPullCacheWindowSeconds
	}
	defer db.Close()

	val := db.GetSetting(constants.SettingPullCacheWindowSeconds)
	if strings.TrimSpace(val) == "" {
		return constants.DefaultPullCacheWindowSeconds
	}
	sec, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil || sec < 0 {
		return constants.DefaultPullCacheWindowSeconds
	}
	return sec
}

func queryRepoPullHistoryWithFallback(db *store.PullSplitDB, cwd string) ([]store.PullRepoRunRecord, error) {
	history, err := db.GetRecentRepoPullHistory(cwd, 1)
	if err == nil && len(history) > 0 {
		return history, nil
	}

	clean := filepath.Clean(cwd)
	if clean != cwd {
		h, errClean := db.GetRecentRepoPullHistory(clean, 1)
		if errClean == nil && len(h) > 0 {
			return h, nil
		}
	}

	alt := resolveDriveLetterVariant(cwd)
	if alt != "" {
		h, errAlt := db.GetRecentRepoPullHistory(alt, 1)
		if errAlt == nil && len(h) > 0 {
			return h, nil
		}
	}

	return history, err
}

func resolveDriveLetterVariant(cwd string) string {
	if len(cwd) <= 1 || cwd[1] != ':' {
		return ""
	}
	firstRune := rune(cwd[0])
	if unicode.IsUpper(firstRune) {
		return string(unicode.ToLower(firstRune)) + cwd[1:]
	}
	return string(unicode.ToUpper(firstRune)) + cwd[1:]
}

func parsePullTimestamp(createdAt string) (time.Time, error) {
	createdAt = strings.TrimSpace(createdAt)
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, createdAt, time.UTC)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("cannot parse pull timestamp")
}
