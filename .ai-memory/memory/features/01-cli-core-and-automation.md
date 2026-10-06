# Feature Catalog 01: CLI Core and Native Automation

- **Domain:** Cobra Command Framework, Universal Help, and Go Native Automation
- **Status:** Authoritative Capability Catalog

## 1. Unified CLI Architecture
- Root dispatch table routes over 60 commands and 120 aliases with sub-millisecond dispatch.
- Universal help interceptor converts Markdown documentation to framed ANSI boxed terminal displays.
- Sub-millisecond Levenshtein typo suggestion engine provides instant fuzzy corrections.

## 2. Native Go Automation Engine (`gitmap automation` / `auto`)
- Replaces legacy Python scripts with compiled Go routines.
- Thread-safe lazy regex compilation registry (`sync.RWMutex`) compiles patterns on demand.
- Fast literal search bypass using Boyer-Moore string matching.
- Polyglot newline and whitespace normalizer enforcing Unix LF line endings across 20+ file formats.
- Sub-millisecond in-memory file cache (`<0.05ms`) with pre-warming and safe purging.
- Side-by-side benchmarking engine comparing native Go performance against script runners.
