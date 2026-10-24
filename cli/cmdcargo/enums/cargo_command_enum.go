package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CargoCommandType represents standard upstream Cargo tool commands.
type CargoCommandType string

const (
	CargoCmdBuild   CargoCommandType = "build"
	CargoCmdCheck   CargoCommandType = "check"
	CargoCmdTest    CargoCommandType = "test"
	CargoCmdClippy  CargoCommandType = "clippy"
	CargoCmdRun     CargoCommandType = "run"
	CargoCmdClean   CargoCommandType = "clean"
	CargoCmdFmt     CargoCommandType = "fmt"
	CargoCmdDoc     CargoCommandType = "doc"
	CargoCmdUpdate  CargoCommandType = "update"
	CargoCmdInit    CargoCommandType = "init"
	CargoCmdNew     CargoCommandType = "new"
	CargoCmdBench   CargoCommandType = "bench"
	CargoCmdPublish CargoCommandType = "publish"
)

// IsValid reports whether the command matches a known Cargo command.
func (c CargoCommandType) IsValid() bool {
	switch c {
	case CargoCmdBuild, CargoCmdCheck, CargoCmdTest, CargoCmdClippy,
		CargoCmdRun, CargoCmdClean, CargoCmdFmt, CargoCmdDoc,
		CargoCmdUpdate, CargoCmdInit, CargoCmdNew, CargoCmdBench,
		CargoCmdPublish:
		return true
	default:
		return false
	}
}

// String returns the string representation of the command.
func (c CargoCommandType) String() string {
	return string(c)
}

// ParseCargoCommandType parses a raw string into a validated CargoCommandType.
func ParseCargoCommandType(s string) (CargoCommandType, *apperror.AppError) {
	val := CargoCommandType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseCargoCommandType", "E1000", fmt.Sprintf("invalid cargo command: %q", s)).WithContext("value", s)
}

// AllCargoCommandTypes returns all valid CargoCommandType variants.
func AllCargoCommandTypes() []CargoCommandType {
	return []CargoCommandType{
		CargoCmdBuild,
		CargoCmdCheck,
		CargoCmdTest,
		CargoCmdClippy,
		CargoCmdRun,
		CargoCmdClean,
		CargoCmdFmt,
		CargoCmdDoc,
		CargoCmdUpdate,
		CargoCmdInit,
		CargoCmdNew,
		CargoCmdBench,
		CargoCmdPublish,
	}
}

// IsBuild reports whether the command is build.
func (c CargoCommandType) IsBuild() bool {
	return c == CargoCmdBuild
}

// IsCheck reports whether the command is check.
func (c CargoCommandType) IsCheck() bool {
	return c == CargoCmdCheck
}

// IsTest reports whether the command is test.
func (c CargoCommandType) IsTest() bool {
	return c == CargoCmdTest
}

// IsClippy reports whether the command is clippy.
func (c CargoCommandType) IsClippy() bool {
	return c == CargoCmdClippy
}

// IsRun reports whether the command is run.
func (c CargoCommandType) IsRun() bool {
	return c == CargoCmdRun
}

// IsClean reports whether the command is clean.
func (c CargoCommandType) IsClean() bool {
	return c == CargoCmdClean
}

// IsFmt reports whether the command is fmt.
func (c CargoCommandType) IsFmt() bool {
	return c == CargoCmdFmt
}

// IsDoc reports whether the command is doc.
func (c CargoCommandType) IsDoc() bool {
	return c == CargoCmdDoc
}

// IsUpdate reports whether the command is update.
func (c CargoCommandType) IsUpdate() bool {
	return c == CargoCmdUpdate
}

// IsInit reports whether the command is init.
func (c CargoCommandType) IsInit() bool {
	return c == CargoCmdInit
}

// IsNew reports whether the command is new.
func (c CargoCommandType) IsNew() bool {
	return c == CargoCmdNew
}

// IsBench reports whether the command is bench.
func (c CargoCommandType) IsBench() bool {
	return c == CargoCmdBench
}

// IsPublish reports whether the command is publish.
func (c CargoCommandType) IsPublish() bool {
	return c == CargoCmdPublish
}

// HasCommand reports whether the command is non-empty and valid.
func (c CargoCommandType) HasCommand() bool {
	return c.IsValid()
}
