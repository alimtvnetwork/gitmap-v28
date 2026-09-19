# Python Constants, Magic Number Elimination & Semantic Decomposition — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-python-constants`, `cg-python-magic-numbers`, `python-constants-enhancement`, `eliminate-magic-numbers`, `audit python constants`, `python-code-enhancement-constants`, `audit python magic numbers`

> **Prompt Version:** 2.3.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all magic numbers, repeated literal tuples, cryptic inline arithmetic, unlabelled multiplications, and oversized functions in Python code. Modifying source files directly, enforce typed constants (`Final`), semantic domain aliases, explicit geometric/typographical variables, strict function length caps (target <= 8 lines, hard cap <= 15 lines), and purposeful helper function extraction until 100% compliant.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all magic numbers, repeated literal tuples, naked multiplications, and functions exceeding 15 lines without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring: extract repeated tuples and numbers into typed `Final` constants, replace cryptic variables (`x2`, `y2`) with meaningful names (`right_x`, `bottom_y`), and assign explicit semantic variables to all calculations.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition: decompose compound rendering or transformation routines into focused single-responsibility helper functions.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.ai-memory/memory/01-index.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/14-constants-enums.md` for constants and enums standards.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase for Magic Numbers/Functions, Write .ai-memory/plans/pending/ Spec, Create .ai-memory/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Extract Constants, Decompose Functions, Rename Variables, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Python Constants, Meaningful Arithmetic & Function Length

Hardcoded numbers, repeated literal tuples, and cryptic calculations obscure intent, increase cognitive load, and introduce bugs during layout or styling changes. Functions combining multiple steps into long routines (> 15 lines) become difficult to test, verify, and understand.

---

### 1. Mandatory Python Constant Standards

1. **Typed `Final` Annotations:**
   All module constants MUST use `typing.Final` along with explicit type annotations:
   ```python
   from typing import Final, Tuple

   RGBAColor = Tuple[int, int, int, int]
   COLOR_RED: Final[RGBAColor] = (255, 0, 0, 255)
   COLOR_WHITE: Final[RGBAColor] = (255, 255, 255, 255)
   ```

2. **Semantic Domain Aliases:**
   Primitive color or dimension constants must be aliased to semantic purpose rather than used raw at callsites:
   ```python
   # ✅ REQUIRED: Semantic aliases convey domain intent
   BOX_BORDER_COLOR: Final[RGBAColor] = COLOR_RED
   LABEL_BACKGROUND_COLOR: Final[RGBAColor] = COLOR_WHITE
   ```

3. **Total Ban on Repeated Tuples & Magic Numbers:**
   Never write inline tuples like `(255, 0, 0, 255)` across multiple loop iterations or function bodies. A single constant MUST define the value.

---

### 2. Meaningful Arithmetic & Self-Documenting Variables

Every mathematical operation (multiplication, division, padding addition, boundary subtraction) MUST have a concrete, unambiguous explanation reflected in variable names and constant definitions:

1. **Explain Every Multiplier:**
   - ❌ `len(text) * 4`: Why `4`? What does `4` represent?
   - ✅ `DIGIT_GLYPH_WIDTH: Final[int] = 4` (Width allocated per digit: glyph raster + kerning spacing).
   - ✅ `glyph_block_width = len(text) * DIGIT_GLYPH_WIDTH`

2. **Explain Every Margin / Padding Addition:**
   - ❌ `+ 2`: Why `2`?
   - ✅ `LABEL_PADDING_X: Final[int] = 1` (1px margin on left and right)
   - ✅ `LABEL_TOTAL_PADDING_X: Final[int] = LABEL_PADDING_X * 2` (Combined horizontal padding)
   - ✅ `label_width = (len(text) * DIGIT_GLYPH_WIDTH) + LABEL_TOTAL_PADDING_X`

3. **Concrete Coordinate Naming (Total Ban on `x2`, `y2`, `t`, `val`):**
   - ❌ `x2 = box.x + box.width - 1`: Cryptic variable names obscure coordinate semantics.
   - ✅ `right_x = box.x + box.width - 1` (Inclusive 0-indexed rightmost pixel column).
   - ✅ `bottom_y = box.y + box.height - 1` (Inclusive 0-indexed bottommost pixel row).

---

### 3. Function Length Limit & Purposeful Extraction

- **Target Function Length:** Target <= 8 lines of body logic.
- **Hard Cap:** <= 15 lines of body logic.
- **Mandatory Purposeful Extraction:** If a function exceeds 15 lines, or performs more than one distinct responsibility (e.g. background bounding calculation, border stroke rendering, and glyph rasterization), extract sub-operations into dedicated, meaningfully named helper functions.
- **Never Compress Lines to Cheat the Cap:** Never cram multiple statements onto a single line, delete required line-gaps, or omit docstrings to fit under the line cap. Decompose with clean helper functions.

---

## 4. Code Review Reference: Before, Intermediate, and Advanced Enhanced Architectures

### 4.1 Before: Anti-Pattern with Magic Numbers & Cryptic Logic

```python
# ❌ ANTI-PATTERN: Repeated magic tuples, magic calculations (* 4 + 2, 7), cryptic coordinates (x2, y2)
def _draw_rect_before(rgba: bytearray, width: int, height: int, box: MarkedBox) -> None:
    x2 = box.x + box.width - 1
    y2 = box.y + box.height - 1
    for x in range(box.x, x2 + 1):
        _set_pixel(rgba, width, height, x, box.y, (255, 0, 0, 255))      # <-- Repeated magic tuple
        _set_pixel(rgba, width, height, x, y2, (255, 0, 0, 255))         # <-- Repeated magic tuple
    for y in range(box.y, y2 + 1):
        _set_pixel(rgba, width, height, box.x, y, (255, 0, 0, 255))      # <-- Repeated magic tuple
        _set_pixel(rgba, width, height, x2, y, (255, 0, 0, 255))         # <-- Repeated magic tuple

def _draw_label_before(rgba: bytearray, width: int, height: int, box: MarkedBox) -> None:
    text = str(box.number)
    # Magic numbers: * 4, + 2, 7, and (255, 255, 255, 255)
    _fill_rect(rgba, width, height, box.x, box.y, len(text) * 4 + 2, 7, (255, 255, 255, 255))
    for offset, digit in enumerate(text):
        _draw_digit(rgba, width, height, box.x + 1 + offset * 4, box.y + 1, digit)
```

---

### 4.2 Centralized Constants Definition (`constants.py`)

```python
"""Constants and typography definitions for image markup rendering."""

from typing import Final, Tuple

# Color definitions (R, G, B, A)
RGBAColor = Tuple[int, int, int, int]

COLOR_RED: Final[RGBAColor] = (255, 0, 0, 255)
COLOR_WHITE: Final[RGBAColor] = (255, 255, 255, 255)

# Semantic domain colors
BOX_BORDER_COLOR: Final[RGBAColor] = COLOR_RED
LABEL_BACKGROUND_COLOR: Final[RGBAColor] = COLOR_WHITE

# Typography and badge layout dimensions (in pixels)
DIGIT_GLYPH_WIDTH: Final[int] = 4     # Width allocated per digit (glyph + kerning)
LABEL_PADDING_X: Final[int] = 1        # 1px margin on left and right
LABEL_PADDING_Y: Final[int] = 1        # 1px margin on top
LABEL_TOTAL_PADDING_X: Final[int] = LABEL_PADDING_X * 2  # Total horizontal padding (+2px)
LABEL_BACKGROUND_HEIGHT: Final[int] = 7  # Fixed height for badge background
```

---

### 4.3 Stage 1 Refactor: Clean, Self-Documenting Constants

```python
def _draw_rect(rgba: bytearray, width: int, height: int, box: MarkedBox) -> None:
    """Draws a 1px border around the marked bounding box using BOX_BORDER_COLOR."""
    right_x = box.x + box.width - 1
    bottom_y = box.y + box.height - 1

    # Draw horizontal top and bottom borders
    for x in range(box.x, right_x + 1):
        _set_pixel(rgba, width, height, x, box.y, BOX_BORDER_COLOR)
        _set_pixel(rgba, width, height, x, bottom_y, BOX_BORDER_COLOR)

    # Draw vertical left and right borders
    for y in range(box.y, bottom_y + 1):
        _set_pixel(rgba, width, height, box.x, y, BOX_BORDER_COLOR)
        _set_pixel(rgba, width, height, right_x, y, BOX_BORDER_COLOR)

def _draw_label(rgba: bytearray, width: int, height: int, box: MarkedBox) -> None:
    """Draws the box number badge with a solid background and padded text."""
    text = str(box.number)

    # Compute dynamic badge dimensions based on number of digits
    label_width = (len(text) * DIGIT_GLYPH_WIDTH) + LABEL_TOTAL_PADDING_X

    # Draw label badge background
    _fill_rect(
        rgba,
        width,
        height,
        box.x,
        box.y,
        label_width,
        LABEL_BACKGROUND_HEIGHT,
        LABEL_BACKGROUND_COLOR,
    )

    # Draw each digit glyph with proper offset and padding
    digit_y = box.y + LABEL_PADDING_Y
    for offset, digit in enumerate(text):
        digit_x = box.x + LABEL_PADDING_X + (offset * DIGIT_GLYPH_WIDTH)
        _draw_digit(rgba, width, height, digit_x, digit_y, digit)
```

---

### 4.4 Stage 2 Advanced Refactor: Purposeful Semantic Decomposition (Strict <= 8–15 Lines)

To maximize readability for AI agents and human reviewers, each discrete responsibility is decomposed into focused helper functions with <= 8 lines of body logic:

```python
def _calculate_label_width(digit_count: int) -> int:
    """Calculates the total pixel width of the label badge including horizontal padding."""
    return (digit_count * DIGIT_GLYPH_WIDTH) + LABEL_TOTAL_PADDING_X

def _draw_horizontal_borders(
    rgba: bytearray, width: int, height: int, start_x: int, end_x: int, top_y: int, bottom_y: int
) -> None:
    """Renders the top and bottom horizontal borders of a bounding box."""
    for x in range(start_x, end_x + 1):
        _set_pixel(rgba, width, height, x, top_y, BOX_BORDER_COLOR)
        _set_pixel(rgba, width, height, x, bottom_y, BOX_BORDER_COLOR)

def _draw_vertical_borders(
    rgba: bytearray, width: int, height: int, start_y: int, end_y: int, left_x: int, right_x: int
) -> None:
    """Renders the left and right vertical borders of a bounding box."""
    for y in range(start_y, end_y + 1):
        _set_pixel(rgba, width, height, left_x, y, BOX_BORDER_COLOR)
        _set_pixel(rgba, width, height, right_x, y, BOX_BORDER_COLOR)

def _draw_rect_modular(rgba: bytearray, width: int, height: int, box: MarkedBox) -> None:
    """Draws a 1px border around the marked bounding box using modular boundary helpers."""
    right_x = box.x + box.width - 1
    bottom_y = box.y + box.height - 1

    _draw_horizontal_borders(rgba, width, height, box.x, right_x, box.y, bottom_y)
    _draw_vertical_borders(rgba, width, height, box.y, bottom_y, box.x, right_x)

def _render_label_digits(
    rgba: bytearray, width: int, height: int, base_x: int, base_y: int, text: str
) -> None:
    """Draws individual digit glyphs with horizontal offset kerning."""
    digit_y = base_y + LABEL_PADDING_Y

    for offset, digit in enumerate(text):
        digit_x = base_x + LABEL_PADDING_X + (offset * DIGIT_GLYPH_WIDTH)
        _draw_digit(rgba, width, height, digit_x, digit_y, digit)

def _draw_label_modular(rgba: bytearray, width: int, height: int, box: MarkedBox) -> None:
    """Draws the box number badge by coordinating background fill and digit rendering."""
    text = str(box.number)
    label_width = _calculate_label_width(len(text))

    _fill_rect(
        rgba, width, height, box.x, box.y, label_width, LABEL_BACKGROUND_HEIGHT, LABEL_BACKGROUND_COLOR
    )
    _render_label_digits(rgba, width, height, box.x, box.y, text)
```

---

## 5. Canonical Golang Equivalent Architecture

For cross-language consistency, here is the exact equivalent Go implementation adhering to repository Go coding guidelines (typed constants, parameter structs when needed, $\le 8$ lines per function, vertical line gaps before `if`, after `}`, and before `return`):

```go
package markup

// RGBAColor defines a 4-channel byte color tuple.
type RGBAColor struct {
	R byte
	G byte
	B byte
	A byte
}

var (
	// ColorRed defines standard solid red.
	ColorRed = RGBAColor{R: 255, G: 0, B: 0, A: 255}

	// ColorWhite defines standard solid white.
	ColorWhite = RGBAColor{R: 255, G: 255, B: 255, A: 255}

	// BoxBorderColor specifies the default color for marked box borders.
	BoxBorderColor = ColorRed

	// LabelBackgroundColor specifies the badge background fill color.
	LabelBackgroundColor = ColorWhite
)

const (
	// DigitGlyphWidth defines horizontal space allocated per digit (glyph + kerning).
	DigitGlyphWidth int = 4

	// LabelPaddingX defines horizontal margin padding on each side.
	LabelPaddingX int = 1

	// LabelPaddingY defines vertical margin padding on top.
	LabelPaddingY int = 1

	// LabelTotalPaddingX defines combined horizontal padding (left + right).
	LabelTotalPaddingX int = LabelPaddingX * 2

	// LabelBackgroundHeight defines the fixed height of the badge background.
	LabelBackgroundHeight int = 7
)

// CalculateLabelWidth computes total badge pixel width including horizontal padding.
func CalculateLabelWidth(digitCount int) int {
	return (digitCount * DigitGlyphWidth) + LabelTotalPaddingX
}
```

---

## 6. Verification Checklist

- [x] **Zero Magic Numbers:** All numerical dimensions, offsets, multipliers, and heights extracted into typed constants.
- [x] **Zero Repeated Literal Tuples:** Colors centralized into `Final` tuples and domain aliases.
- [x] **Concrete Variable Names:** Replaced single-letter/cryptic variables (`x2`, `y2`) with domain-accurate names (`right_x`, `bottom_y`).
- [x] **Meaningful Arithmetic:** Every multiplication and addition corresponds to an explicit, self-documenting concept.
- [x] **Strict Function Size (<= 15 Lines):** Compound routines decomposed into single-responsibility helper functions (target <= 8 lines).
