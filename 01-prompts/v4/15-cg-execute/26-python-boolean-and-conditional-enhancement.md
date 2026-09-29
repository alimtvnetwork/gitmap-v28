[/goal](slashCommand:goal) Autonomously scan, audit, plan, and refactor inefficient Python conditionals, naive inline ternary branching in loops/comprehensions, and magic number thresholds into production-grade, typed, lookup-table-accelerated architectures. Enforce centralized constants, safe clamping, affirmative booleans, and provide high-performance Golang equivalents where required with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

#### High-Speed GitMap Acceleration Options (Run Everything Faster)

Always prefer native GitMap commands over slow generic shell pipelines:
1. **Ultra-Fast File & Directory Discovery (AUM Index & Walk):**
   - **Wildcard / Glob Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Search:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Search:** `gitmap find-files-any "<str>" [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Search:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Directory Tree & Scaffolding:** `gitmap folder-tree` (alias `gitmap ft`)
   - **Zero-Write File Stream:** `gitmap cat <filepath>`
   - **Instant Multi-Core Regex Search:** `gitmap search "<term>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
2. **Fast Repository Hygiene, Lowercase & Symlink Repair:**
   - **Auto-Lowercase Files (Safe 2-Step `git mv`):** `gitmap lowercase` (alias `gitmap lcf [--dry-run]`)
   - **Lowercase Root Readme:** `gitmap lowercase-readme`
   - **Sync Curated `.gitignore` / `.gitattributes` / `.prettierignore`:** `gitmap commons` (alias `gitmap co` or `gitmap sync all`)
   - **Repair Broken Symlinks:** `gitmap fix-link` (alias `gitmap fixlink`)
   - **Clean Update Temp & Inspect Storage:** `gitmap update-cleanup`, `gitmap storage` (alias `gitmap stor`)
3. **Fast Git State, Execution & Atomic Commits:**
   - **Repo Status & Remote Check:** `gitmap status` (`gitmap st`), `gitmap has-any-updates` (`gitmap hau`), `gitmap latest-branch` (`gitmap lb`)
   - **Fast Cross-Platform Shell Runner:** `gitmap pwsh "<command>"` (`gitmap ps`), `gitmap bash "<command>"` (`gitmap sh`), `gitmap async <cmd>` (`gitmap asyn`)
   - **Semantic Atomic Commit & Push:** `gitmap cpf "<summary>"` (Feature), `gitmap cpb "<summary>"` (Bug), `gitmap cpr "<summary>"` (Release), `gitmap pcp "<summary>"` (Pull-Commit-Push)
   - **Smart CI/CD Pipeline Waiting:** `gitmap pe`, `gitmap pipeline-ai status --json` (`gitmap pl-ai status -t <etaSeconds>`)

---

## 1. Problem Statement: Naive Python Conditionals & Magic Numbers

In image processing, stream filtering, and data transformation pipelines, engineers frequently write inline conditionals inside list comprehensions or loops:

```python
# ANTI-PATTERN: Magic numbers, per-pixel branching in hot loop, weak type hints
def _threshold_preview(gray: bytearray, white_threshold: int) -> bytearray:
    threshold = max(0, min(255, round(white_threshold)))
    return bytearray(255 if value >= threshold else value for value in gray)
```

### Why This Is An Anti-Pattern

1. **Magic Numbers:** Hardcoded `0`, `255`, and `round()` boundaries scatter domain constants across the codebase.
2. **CPU Branch Misprediction:** Evaluating `255 if value >= threshold else value` for every byte creates millions of conditional branch checks on large buffers.
3. **Incomplete Type Flexibility:** Restricting to `bytearray` blocks immutable `bytes` inputs.
4. **Poor Testability:** Clamping and thresholding logic are tightly coupled into a single un-composable expression.

---

## 2. The Enhanced Python Architecture

Refactor naive conditionals into centralized constants, safe numeric clamping, and high-performance pre-computed lookup tables (`bytes.translate` executes in C at native hardware speed):

### Constants Module (`constants.py`)

```python
"""Constants used across image processing modules."""

MIN_PIXEL_INTENSITY: int = 0
MAX_PIXEL_INTENSITY: int = 255
WHITE_PIXEL_VALUE: int = 255
```

### Enhanced Function (`threshold.py`)

```python
from typing import Union
import constants

def threshold_preview(
    gray: Union[bytearray, bytes],
    white_threshold: float,
) -> bytearray:
    """Highlight pixels at or above `white_threshold` by setting them to pure white.

    Pixels below the threshold retain their original values.

    Args:
        gray: Raw grayscale pixel byte buffer.
        white_threshold: The cutoff intensity (0-255). Values outside this
            range will be safely clamped.

    Returns:
        A new `bytearray` containing the thresholded preview image.
    """
    # 1. Normalize and safely clamp threshold into valid 8-bit range [0, 255]
    rounded_threshold = round(white_threshold)
    clamped_threshold = max(
        constants.MIN_PIXEL_INTENSITY,
        min(constants.MAX_PIXEL_INTENSITY, rounded_threshold),
    )

    # 2. Precompute 256-byte lookup table (LUT) eliminating per-pixel branching
    lookup_table = bytes(
        constants.WHITE_PIXEL_VALUE if value >= clamped_threshold else value
        for value in range(constants.MAX_PIXEL_INTENSITY + 1)
    )

    # 3. Apply hardware-accelerated translation
    result = bytearray(gray.translate(lookup_table))

    return result
```

---

## 3. Canonical Go (Golang) Translation

When porting or translating this pattern into Go, strictly follow the repository's Go architecture guidelines:
- Separate typed constants in `consts.go`.
- Pure clamping helper with guard clauses (<= 8 lines).
- Pre-computed 256-byte array lookup table (`[256]byte`) eliminating branches in hot loops.
- Sized pre-allocation (`make([]byte, len(gray))`) with zero-overhead indexing.
- Affirmative booleans with `is` prefix (`isAboveThreshold`).
- Mandatory vertical line gaps before `if`, after `}`, and before `return`.

### Constants (`consts.go`)

```go
package imageproc

const (
	// MinPixelIntensity represents the minimum 8-bit grayscale intensity.
	MinPixelIntensity int = 0

	// MaxPixelIntensity represents the maximum 8-bit grayscale intensity.
	MaxPixelIntensity int = 255

	// WhitePixelValue represents pure white intensity in 8-bit grayscale.
	WhitePixelValue byte = 255

	// LookupTableSize defines the number of discrete states in an 8-bit LUT.
	LookupTableSize int = 256
)
```

### Transformation Core (`threshold.go`)

```go
package imageproc

import "math"

// ClampPixelIntensity safely rounds and clamps an arbitrary threshold into [0, 255].
func ClampPixelIntensity(rawThreshold float64) byte {
	rounded := int(math.Round(rawThreshold))

	if rounded < MinPixelIntensity {
		return byte(MinPixelIntensity)
	}

	if rounded > MaxPixelIntensity {
		return byte(MaxPixelIntensity)
	}

	return byte(rounded)
}

// BuildThresholdLUT constructs a 256-byte lookup table mapping values >= threshold to pure white.
func BuildThresholdLUT(clampedThreshold byte) [LookupTableSize]byte {
	var lut [LookupTableSize]byte

	for value := 0; value < LookupTableSize; value++ {
		isAboveThreshold := byte(value) >= clampedThreshold

		if isAboveThreshold {
			lut[value] = WhitePixelValue
			continue
		}

		lut[value] = byte(value)
	}

	return lut
}

// ThresholdPreview highlights pixels at or above whiteThreshold by setting them to pure white.
// Pixels below the threshold retain their original values.
func ThresholdPreview(
	gray []byte,
	whiteThreshold float64,
) []byte {
	if len(gray) == 0 {
		return []byte{}
	}

	clampedThreshold := ClampPixelIntensity(whiteThreshold)
	lut := BuildThresholdLUT(clampedThreshold)

	result := make([]byte, len(gray))

	for i, pixel := range gray {
		result[i] = lut[pixel]
	}

	return result
}
```

### Table-Driven Tests (`threshold_test.go`)

```go
package imageproc_test

import (
	"bytes"
	"testing"

	"coding-guidelines/04-code/golang/pkg/imageproc"
)

func TestThresholdPreview_Scenarios(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		threshold float64
		expected  []byte
	}{
		{
			name:      "StandardThreshold_MiddleIntensity",
			input:     []byte{50, 128, 200},
			threshold: 128.0,
			expected:  []byte{50, 255, 255},
		},
		{
			name:      "NegativeThreshold_ClampsToZero",
			input:     []byte{0, 10, 255},
			threshold: -50.0,
			expected:  []byte{255, 255, 255},
		},
		{
			name:      "OverMaxThreshold_ClampsTo255",
			input:     []byte{100, 254, 255},
			threshold: 300.0,
			expected:  []byte{100, 254, 255},
		},
		{
			name:      "EmptyBuffer_ReturnsEmpty",
			input:     []byte{},
			threshold: 128.0,
			expected:  []byte{},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := imageproc.ThresholdPreview(tc.input, tc.threshold)

			if !bytes.Equal(actual, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
```

---

## 4. Verification Checklist

- [x] **Constants Extracted:** Magic numbers replaced with named domain constants.
- [x] **Safe Clamping:** Numeric inputs safely rounded and bounded before index generation.
- [x] **Zero Loop Branching:** Per-element `if/else` replaced with $O(1)$ lookup table indexing.
- [x] **Affirmative Booleans:** Booleans named with `is`/`has` prefixes and implicitly evaluated.
- [x] **Vertical Line Gaps:** Mandatory blank lines before `if`, after `}`, and before `return`.

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
