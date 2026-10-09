// Package cmddispatch aggregates every claimed command/alias name from the
// scattered dispatch tables (cli/cmd/rootcore.go, cli/cmd/root.go, …) so
// collisions — like the `agm`→gitignore hijack — fail fast instead of shipping.
//
// Aggregation is done by parsing the Go sources with go/ast (no changes to the
// dispatch flow itself). A full table-driven registry with a single Register()
// choke point is a documented follow-up; this package is the detection net.
package cmddispatch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

// DispatchName is one claimed command or alias and where it was claimed.
type DispatchName struct {
	Name  string // resolved name, e.g. "agm"
	Owner string // repo-relative file:line, e.g. "cli/cmd/rootcore.go:78"
	Table string // enclosing function, e.g. "coreBasicMaintenanceEntries"
}

// Collision is one name claimed by two or more distinct owners.
type Collision struct {
	Name   string
	Owners []string // sorted owner refs
}

// CollectNames aggregates every dispatch name from cli/cmd/*.go.
// It returns an error when the sources cannot be located (e.g. the binary was
// built elsewhere); callers must not treat that as a clean bill of health.
func CollectNames() ([]DispatchName, error) {
	root, err := repoRoot()
	if err != nil {
		return nil, err
	}
	consts, err := loadConstantStrings(root)
	if err != nil {
		return nil, err
	}
	files, err := filepath.Glob(filepath.Join(root, "cli", "cmd", "*.go"))
	if err != nil {
		return nil, err
	}
	var names []DispatchName
	for _, path := range files {
		if isTestFile(path) {
			continue
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		got, err := collectFromFile(path, filepath.ToSlash(rel), consts)
		if err != nil {
			return nil, fmt.Errorf("cmddispatch: parse %s: %w", rel, err)
		}
		names = append(names, got...)
	}
	sort.Slice(names, func(i, j int) bool {
		if names[i].Owner != names[j].Owner {
			return names[i].Owner < names[j].Owner
		}
		return names[i].Name < names[j].Name
	})
	return names, nil
}

// FindCollisions returns one Collision per name claimed by 2+ distinct owners.
func FindCollisions(names []DispatchName) []Collision {
	ownersByName := make(map[string][]string)
	order := make([]string, 0)
	for _, n := range names {
		if _, seen := ownersByName[n.Name]; !seen {
			order = append(order, n.Name)
		}
		ownersByName[n.Name] = appendUnique(ownersByName[n.Name], n.Owner)
	}
	var out []Collision
	for _, name := range order {
		owners := ownersByName[name]
		if len(owners) > 1 {
			out = append(out, Collision{Name: name, Owners: owners})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// RunCheck prints the registry report and returns an exit code:
// 0 = no collisions, 1 = collisions found, 2 = sources unavailable.
func RunCheck(w io.Writer) int {
	names, err := CollectNames()
	if err != nil {
		fmt.Fprintf(w, "dispatch registry: cannot collect names: %v\n", err)
		return 2
	}
	fmt.Fprintf(w, "dispatch registry: %d names collected\n", len(names))
	for _, n := range names {
		fmt.Fprintf(w, "  %-28s %s  [%s]\n", n.Name, n.Owner, n.Table)
	}
	collisions := FindCollisions(names)
	if len(collisions) == 0 {
		fmt.Fprintln(w, "dispatch registry: no collisions")
		return 0
	}
	fmt.Fprintf(w, "dispatch registry: %d COLLISION(S):\n", len(collisions))
	for _, c := range collisions {
		fmt.Fprintf(w, "  %q claimed by:\n", c.Name)
		for _, o := range c.Owners {
			fmt.Fprintf(w, "    - %s\n", o)
		}
	}
	return 1
}

// repoRoot locates the repository root from this package's compiled-in source
// path. It only works when running from a repo checkout (dev, CI, go run);
// installed binaries built elsewhere get a clear error, never a false pass.
func repoRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cmddispatch: runtime.Caller failed")
	}
	// file = <repo>/cli/cmddispatch/registry.go
	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	if st, err := os.Stat(filepath.Join(root, "cli", "cmd")); err != nil || !st.IsDir() {
		return "", fmt.Errorf("cmddispatch: repo root not found from %s", file)
	}
	return root, nil
}

func appendUnique(list []string, s string) []string {
	for _, v := range list {
		if v == s {
			return list
		}
	}
	return append(list, s)
}

func isTestFile(path string) bool {
	return len(path) > 8 && path[len(path)-8:] == "_test.go"
}
