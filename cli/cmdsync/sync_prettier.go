package cmdsync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
)

var defaultPrettierrcBaseline = map[string]any{
	"printWidth":    100.0,
	"semi":          true,
	"singleQuote":   false,
	"trailingComma": "all",
}

// RunCommonPrettierRC does a JSON key-union. Existing keys are kept unless
// --force is passed; missing keys are inserted from the baseline.
func RunCommonPrettierRC(dry, force bool) error {
	path := constants.FilePrettierRC
	current, loadErr := loadExistingPrettierRC(path)
	if loadErr != nil {
		return loadErr
	}

	added, overwritten := mergePrettierBaseline(current, force)
	if len(added) == 0 && len(overwritten) == 0 {
		fmt.Printf("  ok  %s already has all curated keys\n", path)

		return nil
	}

	sort.Strings(added)
	sort.Strings(overwritten)

	if dry {
		printDryRunPrettierKeys(path, added, overwritten)

		return nil
	}

	return writePrettierRC(path, current, added, overwritten)
}

func loadExistingPrettierRC(path string) (map[string]any, error) {
	current := map[string]any{}
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return current, nil
	}

	if parseErr := parsePrettierRC(path, data, &current); parseErr != nil {
		return nil, parseErr
	}

	return current, nil
}

func mergePrettierBaseline(current map[string]any, force bool) (added, overwritten []string) {
	for k, v := range defaultPrettierrcBaseline {
		existing, has := current[k]
		if !has {
			current[k] = v
			added = append(added, k)
			continue
		}

		if force && !syncJSONEqual(existing, v) {
			current[k] = v
			overwritten = append(overwritten, k)
		}
	}

	return added, overwritten
}

func printDryRunPrettierKeys(path string, added, overwritten []string) {
	if len(added) > 0 {
		fmt.Printf("  +   %s would add: %s\n", path, strings.Join(added, ", "))
	}

	if len(overwritten) > 0 {
		fmt.Printf("  ~   %s would overwrite (--force): %s\n", path, strings.Join(overwritten, ", "))
	}
}

func writePrettierRC(path string, current map[string]any, added, overwritten []string) error {
	out, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return apperror.WrapWithDetails(
			err,
			"cmd.common.marshalPrettierRC",
			constants.ErrCodeMarshalPrettierRC,
			"failed to marshal .prettierrc JSON",
			"cmd.common",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
	}

	out = append(out, '\n')
	if writeErr := os.WriteFile(path, out, 0o644); writeErr != nil {
		return apperror.WrapWithDetails(
			writeErr,
			"cmd.common.writePrettierRC",
			constants.ErrCodeWritePrettierRC,
			"failed to write .prettierrc",
			"cmd.common",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
	}

	printAppliedPrettierKeys(path, added, overwritten)

	return nil
}

func printAppliedPrettierKeys(path string, added, overwritten []string) {
	if len(added) > 0 {
		fmt.Printf("  +   %s: added keys %s\n", path, strings.Join(added, ", "))
	}

	if len(overwritten) > 0 {
		fmt.Printf("  ~   %s: overwrote keys %s\n", path, strings.Join(overwritten, ", "))
	}
}

// syncJSONEqual compares two JSON-decoded values by re-marshaling.
func syncJSONEqual(a, b any) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	return bytes.Equal(ab, bb)
}

func parsePrettierRC(path string, data []byte, current *map[string]any) error {
	err := json.Unmarshal(data, current)
	if err != nil {
		return apperror.WrapWithDetails(
			err,
			"cmd.common.parsePrettierRC",
			constants.ErrCodeParsePrettierRC,
			"failed to parse existing .prettierrc JSON",
			"cmd.common",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
	}

	return nil
}
