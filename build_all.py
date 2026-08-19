import os

def write_file(path, content):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        f.write(content.strip() + "\n")

write_file("gitmap/db/migrations/004_node_path_alias.sql", """
ALTER TABLE ClusterNode ADD COLUMN DefaultPath TEXT NOT NULL DEFAULT '';

CREATE TABLE NodePathAlias (
    NodePathAliasId INTEGER PRIMARY KEY AUTOINCREMENT,
    NodeId TEXT NOT NULL,
    Alias TEXT NOT NULL,
    AbsolutePath TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL,
    UNIQUE(NodeId, Alias)
);
""")

write_file("gitmap/db/nodepath.go", """
package db

import "database/sql"

type NodePathAlias struct {
	NodePathAliasId int64
	NodeId          string
	Alias           string
	AbsolutePath    string
}

func UpsertPathAlias(db *sql.DB, nodeId, alias, absPath string) error { return nil }
func GetPathAlias(db *sql.DB, nodeId, alias string) (string, error) { return "", nil }
func SetDefaultPath(db *sql.DB, nodeId, path string) error { return nil }
func GetDefaultPath(db *sql.DB, nodeId string) (string, error) { return "", nil }
func ListPathAliases(db *sql.DB, nodeId string) ([]NodePathAlias, error) { return nil, nil }
""")

write_file("gitmap/cluster/pathalias.go", """
package cluster

type AliasEntry struct {
	Alias string
	Path  string
}

func ParseSetPathAliasArg(raw string) ([]AliasEntry, error) { return nil, nil }
""")

write_file("gitmap/cluster/pathresolver.go", """
package cluster

import "database/sql"

func ResolvePath(db *sql.DB, nodeId, pathOrAlias string) (string, error) { return "", nil }
""")

write_file("gitmap/cluster/exec_clone.go", """
package cluster

import "context"

func ExecCloneURL(ctx context.Context, node ClusterNode, url, workdir, subCmd string) (string, string, int, error) {
	return "", "", 0, nil
}
""")

write_file("gitmap/cluster/exec_update.go", """
package cluster

import "context"

func ExecUpdate(ctx context.Context, node ClusterNode, isAll bool, packages ...string) (string, string, int, error) {
	return "", "", 0, nil
}
""")

write_file("gitmap/cluster/lsrender_test.go", """
package cluster

import "testing"

func TestRenderNodeTable(t *testing.T) {}
""")

help_files = [
    "servers-ls.md", "clients-ls.md", "servers-clients-ls.md",
    "cluster-set-default-path.md", "cluster-set-path-alias.md",
    "cluster-cat.md", "cluster-write.md", "cluster-update.md", "cluster-update-all.md"
]
for h in help_files:
    write_file(f"gitmap/cmd/help/{hy", f"# {h}\nXelp content.")

write_file("gitmap/cmd/cluster_h_stubs.go", """
package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"

func runClusterLS(selector cluster.TargetSelectorType, args []string) {}
func runClusterCat(selector cluster.TargetSelectorType, args []string) {}
func runClusterWrite(selector cluster.TargetSelectorType, args []string) {}
func runClusterSetDefaultPath(selector cluster.TargetSelectorType, args []string) {}
func runClusterSetPathAlias(selector cluster.TargetSelectorType, args []string) {}
func runClusterUpdate(selector cluster.TargetSelectorType, isAll bool, args []string) {}
func runClusterClone(selector cluster.TargetSelectorType, subCmd string, args []string) {}
""")

