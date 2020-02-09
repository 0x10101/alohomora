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

	t := term.NewTable()
	t.MaxColWidth = 6

	t.AddRow("Hi", "There", "this")
	t.AddRow("Is", "a test (a bit longer)")

	fmt.Printf(t.Format())

	if !term.Supported() {
		term.NoColors()
	}

	flag.Usage = opts.Usage

	opts, err := opts.Parse()

	if err != nil {
		term.Error("Unable to start: %s\n", err)
		os.Exit(1)
	}

	if opts.Unfancy {
		term.NoColors()
	}

	if opts.Verbose {
		core.Banner(opts.Server)
	}

	if opts.Server {

		server, err := core.Serve(opts)

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
		report(server, opts.ReportJSONTarget, opts.ReportXMLTarget)

	} else {

		found := ext.AircrackAvailable()
		if !found {
			term.Problem("Aircrack not found\n")
			os.Exit(1)
		}

		var tries uint = 0
		var started bool = false
		for {

			client, err := core.Connect(opts)
			if err != nil {
				term.Warn("Could not connect to %s:%d, will try again in %d seconds (%d of %d)\n", opts.Host, opts.Port, 10, tries+1, opts.ConnectionAttempts)
				tries++
				if tries < opts.ConnectionAttempts {
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
			term.Error("Unable to connect to %s:%d, make sure %s is listening there and try again\n", opts.Host, opts.Port, core.Project)
		}

	}
}
