package opts

import (
	"fmt"

	"github.com/steps0x29a/alohomora/term"
)

const (
	// Project is the actual project name
	Project string = "alohomora"

	// Version determines the current version of alohomora
	Version string = "0.4"

	// Codename is a name for a version. Not used until 1.0 is released
	Codename string = ""

	// Author is me :)
	Author string = "Stefan 'steps0x29a' Matyba"

	// Website will point to the documentation website once there is one
	Website string = "https://github.com/steps0x29a/alohomora"
)

// LegacyBanner prints a fancy banner.
func LegacyBanner() {
	fmt.Println("")
	fmt.Println("    █████╗ ██╗      ██████╗ ██╗  ██╗ ██████╗ ███╗   ███╗ ██████╗ ██████╗  █████╗ ")
	fmt.Println("   ██╔══██╗██║     ██╔═══██╗██║  ██║██╔═══██╗████╗ ████║██╔═══██╗██╔══██╗██╔══██╗")
	fmt.Println("   ███████║██║     ██║   ██║███████║██║   ██║██╔████╔██║██║   ██║██████╔╝███████║")
	fmt.Println("   ██╔══██║██║     ██║   ██║██╔══██║██║   ██║██║╚██╔╝██║██║   ██║██╔══██╗██╔══██║")
	fmt.Println("   ██║  ██║███████╗╚██████╔╝██║  ██║╚██████╔╝██║ ╚═╝ ██║╚██████╔╝██║  ██║██║  ██║")
	fmt.Println("   ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝")
	fmt.Printf("                                    v%s                                          \n", Version)
	//fmt.Println("\n                          Made with 🖤 by steps0x29a                          ")
	fmt.Println()
}

// Banner prints a smaller, more likeable banner
func Banner(server, forUsage bool) {
	var mode string = term.BrightCyan("client")
	if server {
		mode = term.BrightMagenta("server")
	}
	if forUsage {
		fmt.Printf("%s v%s %s\n\n", Project, Version, term.Dim(fmt.Sprintf("[by %s]", Author)))
	} else {
		fmt.Printf("%s v%s %s - running in %s mode\n\n", Project, Version, term.Dim(fmt.Sprintf("[by %s]", Author)), mode)
	}
}
