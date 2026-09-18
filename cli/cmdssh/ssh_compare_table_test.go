package cmdssh

import (
	"testing"
)

func TestCompareTableStructure(t *testing.T) {
	cols := compareTableColumns()
	if len(cols) != 6 {
		t.Fatalf("expected 6 columns in compare table, got %d", len(cols))
	}

	rows := compareTableRows()
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows for triad subsystems, got %d", len(rows))
	}

	subsystems := []string{"gitmap ssh", "gitmap cluster", "gitmap sc"}
	for i, sub := range subsystems {
		if rows[i].Cells[0] != sub {
			t.Errorf("row %d expected %s, got %s", i, sub, rows[i].Cells[0])
		}
	}
}

func TestRunSSHCompareCLI(t *testing.T) {
	err := runSSHCompareCLI(nil)
	if err != nil {
		t.Fatalf("runSSHCompareCLI returned unexpected error: %v", err)
	}
}
