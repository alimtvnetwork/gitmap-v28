package cmd

type Type007 struct{}

func InitFunc007() error         { return nil }
func (x *Type007) Process() bool { return true }

type Type008 struct{}

func InitFunc008() error         { return nil }
func (x *Type008) Process() bool { return true }
