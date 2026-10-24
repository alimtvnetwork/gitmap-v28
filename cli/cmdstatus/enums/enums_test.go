package enums

import "testing"

func TestStatusFilterModeType_Valid(t *testing.T) {
	modes := AllStatusFilterModes()
	if len(modes) != 4 {
		t.Fatalf("expected 4 filter modes, got %d", len(modes))
	}

	for _, m := range modes {
		if !m.IsValid() {
			t.Errorf("expected %s to be valid", m)
		}

		if m.String() != string(m) {
			t.Errorf("expected String() == %s, got %s", string(m), m.String())
		}
	}

	if StatusFilterModeType("unknown").IsValid() {
		t.Errorf("expected invalid filter mode to be false")
	}
}

func TestStatusFilterModeType_Booleans(t *testing.T) {
	if !StatusFilterAll.IsAll() || StatusFilterDirty.IsAll() {
		t.Errorf("unexpected IsAll behavior")
	}

	if !StatusFilterDirty.IsDirty() || StatusFilterAhead.IsDirty() {
		t.Errorf("unexpected IsDirty behavior")
	}

	if !StatusFilterAhead.IsAhead() || StatusFilterBehind.IsAhead() {
		t.Errorf("unexpected IsAhead behavior")
	}

	if !StatusFilterBehind.IsBehind() || StatusFilterAll.IsBehind() {
		t.Errorf("unexpected IsBehind behavior")
	}
}

func TestStatusFilterModeType_Parse(t *testing.T) {
	parsed, err := ParseStatusFilterMode("dirty")
	if err != nil || parsed != StatusFilterDirty {
		t.Errorf("expected dirty, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatusFilterModeType(" ALL ")
	if errAlias != nil || parsedAlias != StatusFilterAll {
		t.Errorf("expected all, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatusFilterMode("invalid-mode")
	if invalidErr == nil {
		t.Errorf("expected error for invalid mode")
	}

	if len(AllStatusFilterModeTypes()) != 4 {
		t.Errorf("expected AllStatusFilterModeTypes count to match")
	}
}

func TestStatusOutputModeType_Valid(t *testing.T) {
	modes := AllStatusOutputModes()
	if len(modes) != 3 {
		t.Fatalf("expected 3 output modes, got %d", len(modes))
	}

	for _, o := range modes {
		if !o.IsValid() {
			t.Errorf("expected %s to be valid", o)
		}

		if o.String() != string(o) {
			t.Errorf("expected String() == %s, got %s", string(o), o.String())
		}
	}

	if StatusOutputModeType("unknown").IsValid() {
		t.Errorf("expected invalid output mode to be false")
	}
}

func TestStatusOutputModeType_Booleans(t *testing.T) {
	if !StatusOutputTable.IsTable() || StatusOutputJSON.IsTable() {
		t.Errorf("unexpected IsTable behavior")
	}

	if !StatusOutputJSON.IsJSON() || StatusOutputCompact.IsJSON() {
		t.Errorf("unexpected IsJSON behavior")
	}

	if !StatusOutputCompact.IsCompact() || StatusOutputTable.IsCompact() {
		t.Errorf("unexpected IsCompact behavior")
	}
}

func TestStatusOutputModeType_Parse(t *testing.T) {
	parsed, err := ParseStatusOutputMode("json")
	if err != nil || parsed != StatusOutputJSON {
		t.Errorf("expected json, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatusOutputModeType("COMPACT")
	if errAlias != nil || parsedAlias != StatusOutputCompact {
		t.Errorf("expected compact, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatusOutputMode("xml")
	if invalidErr == nil {
		t.Errorf("expected error for invalid output mode")
	}

	if len(AllStatusOutputModeTypes()) != 3 {
		t.Errorf("expected AllStatusOutputModeTypes count to match")
	}
}

func TestStatusTableColumnType_Valid(t *testing.T) {
	cols := AllStatusTableColumns()
	if len(cols) != 6 {
		t.Fatalf("expected 6 columns, got %d", len(cols))
	}

	for _, c := range cols {
		if !c.IsValid() {
			t.Errorf("expected %s to be valid", c)
		}

		if c.String() != string(c) {
			t.Errorf("expected String() == %s, got %s", string(c), c.String())
		}
	}

	if StatusTableColumnType("unknown").IsValid() {
		t.Errorf("expected invalid column to be false")
	}
}

func TestStatusTableColumnType_Booleans(t *testing.T) {
	if !StatusTableColRepo.IsRepo() || StatusTableColBranch.IsRepo() {
		t.Errorf("unexpected IsRepo behavior")
	}

	if !StatusTableColBranch.IsBranch() || StatusTableColStatus.IsBranch() {
		t.Errorf("unexpected IsBranch behavior")
	}

	if !StatusTableColStatus.IsStatus() || StatusTableColSync.IsStatus() {
		t.Errorf("unexpected IsStatus behavior")
	}

	if !StatusTableColSync.IsSync() || StatusTableColStash.IsSync() {
		t.Errorf("unexpected IsSync behavior")
	}

	if !StatusTableColStash.IsStash() || StatusTableColFiles.IsStash() {
		t.Errorf("unexpected IsStash behavior")
	}

	if !StatusTableColFiles.IsFiles() || StatusTableColRepo.IsFiles() {
		t.Errorf("unexpected IsFiles behavior")
	}
}

func TestStatusTableColumnType_Parse(t *testing.T) {
	parsed, err := ParseStatusTableColumn("branch")
	if err != nil || parsed != StatusTableColBranch {
		t.Errorf("expected branch, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatusTableColumnType("FILES")
	if errAlias != nil || parsedAlias != StatusTableColFiles {
		t.Errorf("expected files, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatusTableColumn("nonexistent")
	if invalidErr == nil {
		t.Errorf("expected error for invalid table column")
	}

	if len(AllStatusTableColumnTypes()) != 6 {
		t.Errorf("expected AllStatusTableColumnTypes count to match")
	}
}

func TestStatusRepoStateType_Valid(t *testing.T) {
	states := AllStatusRepoStates()
	if len(states) != 6 {
		t.Fatalf("expected 6 repo states, got %d", len(states))
	}

	for _, s := range states {
		if !s.IsValid() {
			t.Errorf("expected %s to be valid", s)
		}

		if s.String() != string(s) {
			t.Errorf("expected String() == %s, got %s", string(s), s.String())
		}
	}

	if StatusRepoStateType("unknown").IsValid() {
		t.Errorf("expected invalid repo state to be false")
	}
}

func TestStatusRepoStateType_Booleans(t *testing.T) {
	if !StatusRepoStateClean.IsClean() || StatusRepoStateDirty.IsClean() {
		t.Errorf("unexpected IsClean behavior")
	}

	if !StatusRepoStateDirty.IsDirty() || StatusRepoStateAhead.IsDirty() {
		t.Errorf("unexpected IsDirty behavior")
	}

	if !StatusRepoStateAhead.IsAhead() || StatusRepoStateBehind.IsAhead() {
		t.Errorf("unexpected IsAhead behavior")
	}

	if !StatusRepoStateBehind.IsBehind() || StatusRepoStateStashed.IsBehind() {
		t.Errorf("unexpected IsBehind behavior")
	}

	if !StatusRepoStateStashed.IsStashed() || StatusRepoStateMissing.IsStashed() {
		t.Errorf("unexpected IsStashed behavior")
	}

	if !StatusRepoStateMissing.IsMissing() || StatusRepoStateClean.IsMissing() {
		t.Errorf("unexpected IsMissing behavior")
	}
}

func TestStatusRepoStateType_Parse(t *testing.T) {
	parsed, err := ParseStatusRepoState("missing")
	if err != nil || parsed != StatusRepoStateMissing {
		t.Errorf("expected missing, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatusRepoStateType("CLEAN")
	if errAlias != nil || parsedAlias != StatusRepoStateClean {
		t.Errorf("expected clean, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatusRepoState("broken")
	if invalidErr == nil {
		t.Errorf("expected error for invalid repo state")
	}

	if len(AllStatusRepoStateTypes()) != 6 {
		t.Errorf("expected AllStatusRepoStateTypes count to match")
	}
}
