package cmd

type CloneUIPart1 struct{}

func InitCloneUIPart1() error          { return nil }
func (x *CloneUIPart1) Process1() bool { return true }

type CloneUIPart2 struct{}

func InitCloneUIPart2() error          { return nil }
func (x *CloneUIPart2) Process2() bool { return true }
