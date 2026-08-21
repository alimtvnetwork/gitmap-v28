package cmd

import (
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
)

type subCmdTestCase struct {
	name     string
	tokens   []string
	expected []cluster.ClusterSubCommand
	wantErr  bool
}

func clusterSubCommandTestCases() []subCmdTestCase {
	return []subCmdTestCase{
		{
			name:     "chaining ps and cmd",
			tokens:   []string{"ps", `"echo hello"`, ",", "cmd", `"dir"`},
			expected: []cluster.ClusterSubCommand{{Kind: db.CommandKindPsCommand, RawArg: `"echo hello"`}, {Kind: db.CommandKindCmdCommand, RawArg: `"dir"`}},
		},
		{
			name:     "chaining ps and install",
			tokens:   []string{"ps", `"echo hi"`, ",", "install", `"pkg1,pkg2"`},
			expected: []cluster.ClusterSubCommand{{Kind: db.CommandKindPsCommand, RawArg: `"echo hi"`}, {Kind: db.CommandKindInstall, RawArg: `"pkg1,pkg2"`}},
		},
		{name: "unknown subcommand", tokens: []string{"unknowncmd", `"foo"`}, wantErr: true},
	}
}

func assertSubCmdResult(t *testing.T, tt subCmdTestCase) {
	got, err := ParseSubCommands(tt.tokens)
	if (err != nil) != tt.wantErr {
		t.Errorf("ParseSubCommands() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
		t.Errorf("ParseSubCommands() = %v, expected %v", got, tt.expected)
	}
}

func TestClusterSubCommandParser(t *testing.T) {
	for _, tt := range clusterSubCommandTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			assertSubCmdResult(t, tt)
		})
	}
}
