package funcintel

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestLanguageForPath(t *testing.T) {
	cases := map[string]string{
		"src/foo.go": constants.CommitInLanguageGo,
		"a/b/c.TS":   constants.CommitInLanguageTypeScript,
		"a.tsx":      constants.CommitInLanguageTypeScript,
		"a.mjs":      constants.CommitInLanguageJavaScript,
		"a.unknown":  "",
		"NoExt":      "",
	}
	for path, want := range cases {
		if got := LanguageForPath(path); got != want {
			t.Errorf("LanguageForPath(%q)=%q want %q", path, got, want)
		}
	}
}

func TestGoDetectorAddedNames(t *testing.T) {
	prev := "package x\nfunc Old() {}\n"
	new := "package x\nfunc Old() {}\nfunc New(a int) error { return nil }\nfunc (r *T) Method() {}\n"
	got := Get(constants.CommitInLanguageGo).Detect(prev, new)
	want := []string{"Method", "New"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestTsDetectorCoversArrowAndFunction(t *testing.T) {
	prev := ""
	new := "export function fooFn(x) {}\nexport const useDebounce = (v) => {}\n"
	got := Get(constants.CommitInLanguageTypeScript).Detect(prev, new)
	want := []string{"fooFn", "useDebounce"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestPythonDetector(t *testing.T) {
	got := Get(constants.CommitInLanguagePython).Detect("", "def foo(x):\n    return x\n")
	if len(got) != 1 || got[0] != "foo" {
		t.Fatalf("py detector: %v", got)
	}
}

func TestJavaDetector(t *testing.T) {
	got := Get(constants.CommitInLanguageJava).Detect("", "public static int compute(int x) {\n  return x;\n}\n")
	if len(got) != 1 || got[0] != "compute" {
		t.Fatalf("java detector: %v", got)
	}
}

func TestRenderOutputFormat(t *testing.T) {
	changes := []FileChange{
		{RelativePath: "src/z.go", Language: constants.CommitInLanguageGo, NewSource: "func B() {}\nfunc A() {}\n"},
		{RelativePath: "src/a.ts", Language: constants.CommitInLanguageTypeScript, NewSource: "export function useX() {}\n"},
		{RelativePath: "README.md", Language: "", NewlyAddedFile: true},
	}
	got := Render(changes)
	want := "- README.md\n- src/a.ts\n  - added: useX\n- src/z.go\n  - added: A, B"
	if got != want {
		t.Fatalf("render mismatch:\n--got--\n%s\n--want--\n%s", got, want)
	}
}

func TestEnabledLanguagesFiltersUnknown(t *testing.T) {
	got := EnabledLanguages([]string{"Go", "Klingon", "Python"})
	if strings.Join(got, ",") != "Go,Python" {
		t.Fatalf("enabled = %v", got)
	}
}

func TestDetectorIgnoresUnchangedDeclarations(t *testing.T) {
	src := "func Same() {}\n"
	got := Get(constants.CommitInLanguageGo).Detect(src, src)
	if len(got) != 0 {
		t.Fatalf("expected no added, got %v", got)
	}
}
