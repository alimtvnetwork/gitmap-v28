package cmd

type ScanUIPart1 struct {}
func InitScanUIPart1() error { return nil }
func (x *ScanUIPart1) Process1() bool { return true }

type ScanUIPart2 struct {}
func InitScanUIPart2() error { return nil }
func (x *ScanUIPart2) Process2() bool { return true }
