package main

import (
	"github.com/pterm/pterm"
	"time"
)

func main() {
	s, _ := pterm.DefaultSpinner.Start("test")
	time.Sleep(100 * time.Millisecond)
	s.SuccessPrinter = &pterm.PrefixPrinter{
		MessageStyle: &pterm.Style{pterm.FgDarkGray},
		Prefix: pterm.Prefix{
			Text:  "~",
			Style: &pterm.Style{pterm.FgDarkGray},
		},
	}
	s.Success("up to date")
}
