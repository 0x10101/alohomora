package main

import (
	"os"

	"github.com/steps0x29a/alohomora/core"
	"github.com/steps0x29a/alohomora/opts"
	"github.com/steps0x29a/islazy/term"
)

func main() {

	if !term.Supported() {
		term.NoColors()
	}

	opts, err := opts.Parse()

	if err != nil {
		term.Error("Unable to start: %s\n", err)
		os.Exit(1)
	}

	if opts.Unfancy {
		term.NoColors()
	}

	if opts.Verbose {
		core.Banner()
	}

	if opts.Server {

		server, err := core.Serve(opts)
		if err != nil {
			term.Error("Unable to start server: %s\n", err)
			os.Exit(1)
		}

		<-server.Terminated
		server.KickAll()

	} else {

		found := core.AircrackAvailable()
		if found {
			term.Success("Aircrack found\n")
		} else {
			term.Problem("Aircrack not found\n")
			os.Exit(1)
		}

		client, err := core.Connect(opts)
		if err != nil {
			term.Error("Unable to start client: %s\n", err)
			os.Exit(1)
		}

		<-client.Terminated
	}
}
