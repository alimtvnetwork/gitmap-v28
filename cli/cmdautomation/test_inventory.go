package cmdautomation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type testEstimate struct {
	duration float64
	tier     string
	isSlow   bool
}

// RunTestInventory discovers tests, calculates duration baselines, and generates inventory.
func RunTestInventory(opts TestInventoryOptions) TestInventoryResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	threshold := resolveThreshold(opts.SlowThreshold)
	tests := scanAllTests(root, threshold)
	summary := summarizeInventory(tests, threshold)
	outPath := resolveOutPath(opts.OutPath)
	res := TestInventoryResult{
		Version:    1,
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
		TotalTests: len(tests),
		Summary:    summary,
		Tests:      tests,
		OutPath:    outPath,
		Duration:   time.Since(start),
		IsPass:     len(tests) > 0,
	}
	_ = writeInventoryFile(outPath, res)
	return result.Ok(res)
}

func resolveThreshold(val float64) float64 {
	if val > 0.0 {
		return val
	}
	return 4.0
}

func resolveOutPath(custom string) string {
	if len(custom) > 0 {
		return custom
	}
	return filepath.Join(".ai-memory", "test-inventory.json")
}

func scanAllTests(root string, threshold float64) map[string]TestInventoryItem {
	tests := make(map[string]TestInventoryItem)
	files, _ := collectTextFiles(root, []string{".go"})
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			items := extractGoTests(f, root, threshold)
			for _, item := range items {
				tests[item.Id] = item
			}
		}
	}
	return tests
}

func extractGoTests(path, root string, threshold float64) []TestInventoryItem {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	relPath, _ := filepath.Rel(root, path)
	pkg := filepath.ToSlash(filepath.Dir(relPath))
	var items []TestInventoryItem
	for _, line := range strings.Split(string(data), "\n") {
		item := parseGoTestLine(line, filepath.ToSlash(relPath), pkg, string(data), threshold)
		if len(item.Id) > 0 {
			items = append(items, item)
		}
	}
	return items
}

func parseGoTestLine(line, relPath, pkg, content string, threshold float64) TestInventoryItem {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "func Test") {
		return TestInventoryItem{}
	}
	idx := strings.Index(trimmed, "(")
	if idx == -1 {
		return TestInventoryItem{}
	}
	name := strings.TrimSpace(trimmed[5:idx])
	est := estimateDuration(name, content, threshold)
	return TestInventoryItem{
		Id:          fmt.Sprintf("%s.%s", pkg, name),
		Package:     pkg,
		TestFile:    relPath,
		TestFunc:    name,
		DurationSec: est.duration,
		Tier:        est.tier,
		IsSlow:      est.isSlow,
		NeedsRun:    true,
	}
}

func estimateDuration(name, content string, threshold float64) testEstimate {
	dur := 0.01
	if strings.Contains(content, "exec.Command") {
		dur += 2.5
	}
	if strings.Contains(content, "time.Sleep") {
		dur += 1.5
	}
	if strings.Contains(strings.ToLower(name), "heavy") {
		dur += 3.0
	}
	isSlow := dur >= threshold
	tier := "fast"
	if isSlow {
		tier = "slow"
	}
	return testEstimate{duration: dur, tier: tier, isSlow: isSlow}
}

func summarizeInventory(tests map[string]TestInventoryItem, threshold float64) TestInventorySummary {
	slow := 0
	fast := 0
	slowSec := 0.0
	fastSec := 0.0
	pkgs := make(map[string]bool)
	for _, t := range tests {
		pkgs[t.Package] = true
		if t.IsSlow {
			slow++
			slowSec += t.DurationSec
		} else {
			fast++
			fastSec += t.DurationSec
		}
	}
	return TestInventorySummary{
		Total:            len(tests),
		Dirty:            len(tests),
		Packages:         len(pkgs),
		SlowTests:        slow,
		FastTests:        fast,
		SlowThresholdSec: threshold,
		EstimatedSlowSec: slowSec,
		EstimatedFastSec: fastSec,
	}
}

func writeInventoryFile(path string, res TestInventoryResult) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}
