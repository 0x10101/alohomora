package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/steps0x29a/alohomora/core"
	"github.com/steps0x29a/alohomora/ext"
	"github.com/steps0x29a/alohomora/opts"
	"github.com/steps0x29a/alohomora/term"
)

/*func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))
}*/

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

		report := server.Report()

		if opts.ReportXMLTarget != "" {
			xmlBytes, err := xml.MarshalIndent(report, "", "  ")
			if err != nil {
				term.Error("Unable to save XML report: %s\n", err)
			} else {
				err = ioutil.WriteFile(opts.ReportXMLTarget, xmlBytes, 0640)
				if err != nil {
					term.Error("Unable to save XML report: %s\n", err)
				}
			}
		}

		if opts.ReportJSONTarget != "" {
			jsonBytes, err := json.MarshalIndent(report, "", " ")
			if err != nil {
				term.Error("Unable to save JSON report: %s\n", err)
			} else {
				err = ioutil.WriteFile(opts.ReportJSONTarget, jsonBytes, 0640)
				if err != nil {
					term.Error("Unable to save JSON report: %s\n", err)
				}
			}
		}

		report.Print()

	} else {

		found := ext.AircrackAvailable()
		if !found {
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
