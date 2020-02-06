package core

import (
	"fmt"
)

const (
	// Version determines the current version of alohomora
	Version string = "0.4"

	// Codename is a name for a version. Not used until 1.0 is released
	Codename string = ""

	// Author is me :)
	Author string = "Stefan 'steps0x29a' Matyba"

	// Website will point to the documentation website once there is one
	Website string = "https://github.com/steps0x29a/alohomora"
)

// Banner prints a fancy banner.
func Banner() {
	fmt.Println("")
	fmt.Println("    █████╗ ██╗      ██████╗ ██╗  ██╗ ██████╗ ███╗   ███╗ ██████╗ ██████╗  █████╗ ")
	fmt.Println("   ██╔══██╗██║     ██╔═══██╗██║  ██║██╔═══██╗████╗ ████║██╔═══██╗██╔══██╗██╔══██╗")
	fmt.Println("   ███████║██║     ██║   ██║███████║██║   ██║██╔████╔██║██║   ██║██████╔╝███████║")
	fmt.Println("   ██╔══██║██║     ██║   ██║██╔══██║██║   ██║██║╚██╔╝██║██║   ██║██╔══██╗██╔══██║")
	fmt.Println("   ██║  ██║███████╗╚██████╔╝██║  ██║╚██████╔╝██║ ╚═╝ ██║╚██████╔╝██║  ██║██║  ██║")
	fmt.Println("   ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝")
	fmt.Printf("                                    v%s                                          \n", Version)
	fmt.Println("\n                          Made with 🖤 by steps0x29a                          ")
	fmt.Println()
}
