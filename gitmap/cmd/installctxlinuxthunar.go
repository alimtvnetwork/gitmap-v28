package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// installThunar appends a marker-delimited <action> block to
// ~/.config/Thunar/uca.xml. Idempotent: a previous block is replaced
// rather than duplicated.
func installThunar(flat []flatCtxEntry, exe string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	path := filepath.Join(home, constants.CtxLinuxThunarRel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, path, err)

		return false
	}
	cur, _ := os.ReadFile(path)
	body := thunarMerged(string(cur), flat, exe)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, path, err)

		return false
	}

	return true
}

// thunarMerged splices a marker-delimited <action> block into Thunar's
// uca.xml. If a previous block exists it is replaced; otherwise we
// inject just before </actions>. A missing file is created with a
// minimal <actions/> root.
func thunarMerged(cur string, flat []flatCtxEntry, exe string) string {
	block := buildThunarBlock(flat, exe)
	if cur == "" {
		return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<actions>\n" + block + "</actions>\n"
	}
	begin, end := constants.CtxThunarMarkBegin, constants.CtxThunarMarkEnd
	if i := strings.Index(cur, begin); i >= 0 {
		if j := strings.Index(cur[i:], end); j >= 0 {
			return cur[:i] + block + cur[i+j+len(end):]
		}
	}
	if k := strings.LastIndex(cur, "</actions>"); k >= 0 {
		return cur[:k] + block + cur[k:]
	}

	return cur + "\n<actions>\n" + block + "</actions>\n"
}

// buildThunarBlock builds the marker-wrapped <action> entries.
func buildThunarBlock(flat []flatCtxEntry, exe string) string {
	var b strings.Builder
	b.WriteString(constants.CtxThunarMarkBegin + "\n")
	for _, e := range flat {
		fmt.Fprintf(&b, "<action><icon>utilities-terminal</icon><name>%s</name><unique-id>%s</unique-id><command>%s</command><patterns>*</patterns><directories/></action>\n",
			e.Label, e.Slug, dolphinExec(e, exe))
	}
	b.WriteString(constants.CtxThunarMarkEnd + "\n")

	return b.String()
}

// rmDirCtx / rmFileCtx — small helpers returning 1 on success.
func rmDirCtx(path string) int {
	if err := os.RemoveAll(path); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsRmFail, path, err)

		return 0
	}

	return 1
}

func rmFileCtx(path string) int {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsRmFail, path, err)

		return 0
	}

	return 1
}

// stripThunarBlock removes the marker block from uca.xml, leaving any
// user-managed entries intact. Returns 1 if the file was rewritten.
func stripThunarBlock(path string) int {
	cur, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	s := string(cur)
	begin, end := constants.CtxThunarMarkBegin, constants.CtxThunarMarkEnd
	i := strings.Index(s, begin)
	if i < 0 {
		return 0
	}
	j := strings.Index(s[i:], end)
	if j < 0 {
		return 0
	}
	out := s[:i] + s[i+j+len(end):]
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsRmFail, path, err)

		return 0
	}

	return 1
}
