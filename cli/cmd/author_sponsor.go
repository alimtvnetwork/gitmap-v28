// Package cmd — author_sponsor.go displays author and sponsor credentials.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runAuthor(args []string) error {
	checkHelp("author", args)
	printAuthorCard()

	return nil
}

func printAuthorCard() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║                         CREATOR & ARCHITECT                          ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║                          MD ALIM UL KARIM                            ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  %sProfile & Leadership:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    Conceived, designed, and architected GitMap from the ground up.")
	fmt.Println("    Over 20+ years of professional software engineering leadership across")
	fmt.Println("    enterprise distributed systems, fintech, and AI-driven platforms.")
	fmt.Println()

	fmt.Printf("  %sInnovations:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    • Creator of the XProgramming Language (https://the-xproduct.com)")
	fmt.Println("    • Author of opinionated spec-first engineering guidelines & AI frameworks")
	fmt.Println("    • Architect of zero-nesting, shallow control flow developer automation engines")
	fmt.Println()

	fmt.Printf("  %sConnect & Web:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • Personal Website : %shttps://alimkarim.com/%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    • Google Search    : %shttps://www.google.com/search?q=MD+Alim+Ul+Karim%s\n\n", constants.ColorGreen, constants.ColorReset)
}

func runSponsor(args []string) error {
	checkHelp("sponsor", args)
	printSponsorCard()

	return nil
}

func printSponsorCard() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║                          PROUD SPONSOR                               ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║                         RISE UP ASIA LLC                             ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  %sAbout Rise Up Asia LLC:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    Elite software engineering company recognized for delivering world-class,")
	fmt.Println("    spec-driven software for California-based technology leaders (Silicon Valley")
	fmt.Println("    SaaS & fintech) and EU-based product innovators (Germany, Netherlands, Nordics).")
	fmt.Println()

	fmt.Printf("  %sMission & Sponsorship:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    Rise Up Asia LLC proudly backs GitMap as part of its commitment to")
	fmt.Println("    empowering global developers and autonomous AI engineers with")
	fmt.Println("    industrial-grade developer tooling.")
	fmt.Println()

	fmt.Printf("  %sWebsite:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • Official Website : %shttps://riseup-asia.com/%s\n\n", constants.ColorGreen, constants.ColorReset)
}

func runCredits(args []string) error {
	checkHelp("credits", args)
	printAuthorCard()
	printSponsorCard()

	return nil
}
