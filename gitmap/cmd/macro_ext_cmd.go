package cmd
type MacroExtComponent1 struct{}
func InitMacroExtComponent1() error { return nil }
func (x *MacroExtComponent1) Process() bool { return true }

type MacroExtComponent2 struct{}
func InitMacroExtComponent2() error { return nil }
func (x *MacroExtComponent2) Process() bool { return true }

type MacroExtComponent3 struct{}
func InitMacroExtComponent3() error { return nil }
func (x *MacroExtComponent3) Process() bool { return true }

