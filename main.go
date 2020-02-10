package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func report(server *core.Server, jsonFile, xmlFile string) {
	report := server.Report()

	if xmlFile != "" {
		xmlBytes, err := xml.MarshalIndent(report, "", "  ")
		if err != nil {
			term.Error("Unable to save XML report: %s\n", err)
		} else {
			err = ioutil.WriteFile(xmlFile, xmlBytes, 0640)
			if err != nil {
				term.Error("Unable to save XML report: %s\n", err)
			}
		}
	}

	if jsonFile != "" {
		jsonBytes, err := json.MarshalIndent(report, "", " ")
		if err != nil {
			term.Error("Unable to save JSON report: %s\n", err)
		} else {
			err = ioutil.WriteFile(jsonFile, jsonBytes, 0640)
			if err != nil {
				term.Error("Unable to save JSON report: %s\n", err)
			}
		}
	}

	report.Print()
}

func main() {

	if !term.Supported() {
		term.NoColors()
	}

	flag.Usage = opts.Usage

	options, err := opts.Parse()

	if err != nil {
		term.Error("Unable to start: %s\n", err)
		os.Exit(1)
	}

	if options.Unfancy {
		term.NoColors()
	}

	if options.Verbose {
		opts.Banner(options.Server, false)
	}

	if options.Server {

		server, err := core.Serve(options)

		if err != nil {
			term.Error("Unable to start server: %s\n", err)
			os.Exit(1)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			server.Terminate()
		}()

		<-server.Terminated
		server.KickAll()
		report(server, options.ReportJSONTarget, options.ReportXMLTarget)

	} else {

		found := ext.AircrackAvailable()
		if !found {
			term.Problem("Aircrack not found\n")
			os.Exit(1)
		}

		var tries uint = 0
		var started bool = false
		for {

			client, err := core.Connect(options)
			if err != nil {
				term.Warn("Could not connect to %s:%d, will try again in %d seconds (%d of %d)\n", options.Host, options.Port, 10, tries+1, options.ConnectionAttempts)
				tries++
				if tries < options.ConnectionAttempts {
					time.Sleep(10 * time.Second)
					continue
				} else {
					break
				}
			} else {
				started = true
				<-client.Terminated
				client.Shutdown()
				break
			}
		}

		if !started {
			fmt.Println()
			term.Error("Unable to connect to %s:%d, make sure %s is listening there and try again\n", options.Host, options.Port, opts.Project)
		}

	}
}
