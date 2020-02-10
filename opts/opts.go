package opts

import (
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/gosuri/uitable"
	"github.com/steps0x29a/alohomora/bigint"
	"github.com/steps0x29a/alohomora/term"
)

//The Options type wraps all command line options in a neat struct for easier handling.
type Options struct {
	Server             bool
	Port               uint
	Host               string
	Mode               string
	Verbose            bool
	Unfancy            bool
	Charset            string
	Jobsize            string
	Passlen            uint
	Offset             string
	Target             string
	Timeout            uint64
	ReportXMLTarget    string
	ReportJSONTarget   string
	QueueSize          uint64
	MaxJobs            string
	MaxTime            uint64
	EnableREST         bool
	RESTAddress        string
	RESTPort           uint
	ConnectionAttempts uint
}

const (
	maxWidth = 80

	nodeFlag    = "server"
	nodeDefault = false
	nodeHelp    = "If provided, alohomora will run in server mode. If not provided, it will run in client mode"

	modeFlag        = "mode"
	modeFlagShort   = "m"
	modeFlagDefault = "WPA2"

	portFlag        = "port"
	portFlagShort   = "p"
	portFlagDefault = 29100
	portFlagHelp    = "The port that alohomora should listen on (in server mode) or connect to (in client mode). Defaults to 29100."

	serverFlag        = "ip"
	serverFlagShort   = "i"
	serverFlagDefault = "0.0.0.0"
	serverFlagHelp    = "The IP address that alohomora should listen on (in server mode) or connect to (in client mode). Defaults to 0.0.0.0"

	unfancyFlag        = "unfancy"
	unfancyFlagShort   = "u"
	unfancyFlagDefault = false
	unfancyFlagHelp    = "If provided, alohomora's output will not be colored with ANSI escape codes. Not recommended"

	verboseFlag        = "verbose"
	verboseFlagShort   = "v"
	verboseFlagDefault = false
	verboseFlagHelp    = "If provided, alohomora will give (a lot of) additional output"

	jobsizeFlag      = "jobsize"
	jobsizeFlagShort = "j"
	jobsizeDefault   = "10000"
	jobsizeHelp      = "The amount of passwords that each connected client should bruteforce for each job. Defaults to 10000."

	charsetFlag        = "charset"
	charsetFlagShort   = "c"
	charsetFlagDefault = "0123456789"
	charsetFlagHelp    = "The charset from which the clients should generate passwords. Defaults to 0123456789. Enclose in \"\" for special characters and escape them properly!"

	lengthFlag        = "length"
	lengthFlagShort   = "l"
	lengthFlagDefault = 8
	lengthFlagHelp    = "The length (in characters) of passwords the clients should generate during cracking attempts. Defaults to 8. Must not be negative nor 0."

	offsetFlag        = "offset"
	offsetFlagShort   = "o"
	offsetFlagDefault = "0"
	offsetFlagHelp    = "The amount of passwords to skip before even trying to crack. Defaults to 0. Must not be negative."

	targetFlag        = "target"
	targetFlagShort   = "t"
	targetFlagDefault = ""
	targetFlagHelp    = "Full path to a valid .pcap file containing the handshake to crack"

	timeoutFlag        = "timeout"
	timeoutFlagShort   = "x"
	timeoutFlagDefault = 600
	timeoutFlagHelp    = "Amount of seconds before a crack job times out and is considered lost (job will be rescheduled). Defaults to 600 (10 minutes)."

	queueSizeFlag        = "queuesize"
	queueSizeFlagDefault = 50
	queueSizeFlagHelp    = "Amount of crackjobs to generate as a backlog. Defaults to 50. Must neither be negative nor 0."

	maxJobsFlag        = "maxjobs"
	maxJobsFlagDefault = "0"
	maxJobsFlagHelp    = "Maximum amount of jobs to dispatch to clients for current handshake. Defaults to 0 (disabled)."

	maxTimeFlag        = "maxtime"
	maxTimeFlagDefault = 0
	maxTimeFlagHelp    = "Amount of seconds before the server considers the current handshake as a failure. Defaults to 0 (disabled)."

	attemptsFlag        = "attempts"
	attemptsFlagShort   = "a"
	attemptsFlagDefault = 5
	attemptsFlagHelp    = "Number of connection attempts to a server (default is 5)."

	reportXMLFlag         = "oX"
	reportXMLFlagDefault  = ""
	reportJSONFlag        = "oJ"
	reportJSONFlagDefault = ""

	restFlag        = "rest"
	restFlagDefault = false

	restAddressFlag        = "restaddr"
	restAddressFlagDefault = "127.0.0.1"

	restPortFlag        = "restport"
	restPortFlagDefault = 29101
)

func intro() {
	Banner(true, true)
	fmt.Printf("Usage: \n\n ./alohomora --server --target <FILE> [options]\n ./alohomora --ip <IP> --port <PORT>\n\n")
}

func serverIntro() {
	fmt.Println(term.Bold("SERVER MODE USAGE"))
	fmt.Println()
	fmt.Println("  ./alohomora --server --target <FILE>")
	fmt.Println()
	t := uitable.New()
	t.MaxColWidth = 80
	t.Wrap = true // wrap columns
	t.AddRow("", "Runs alohomora server on 0.0.0.0:29100 targeting <FILE>. The character set used to generate passwords will be 0123456789. Each password will be 8 characters long and each client will be tasked with 10.000 passwords per job.")
	fmt.Println(t)
	fmt.Println()
	fmt.Println(term.Bold("SERVER OPTIONS"))
	fmt.Println()
}

func clientIntro() {
	fmt.Println(term.Bold("CLIENT MODE USAGE"))
	fmt.Println()
	fmt.Println("  ./alohomora --ip <IP> --port <PORT>")
	fmt.Println()
	t := uitable.New()
	t.MaxColWidth = 80
	t.Wrap = true // wrap columns
	t.AddRow("", "Runs alohomora client, connecting to <IP> on <PORT> and waiting for new jobs.")
	fmt.Println(t)
	fmt.Println()
	fmt.Println(term.Bold("CLIENT OPTIONS"))
	fmt.Println()
}

func Usage() {
	intro()
	serverIntro()

	t := uitable.New()
	t.MaxColWidth = 60
	t.Wrap = true // wrap columns

	t.AddRow(fmt.Sprintf("  -%s / --%s <IP>", serverFlagShort, serverFlag), "Set the server's listen address. Defaults to 0.0.0.0.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <PORT>", portFlagShort, portFlag), "Set the server's listen port. Defaults to 29100.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <MODE>", modeFlagShort, modeFlag), "Set the mode of operation. Default is 'WPA2'.\n\nCurrently supported modes:\n\n  - WPA2 [WPA2 handshake cracking] (default)\n\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <PATH>", targetFlagShort, targetFlag), "Set the cracking target (path to a file).\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <STRING>", charsetFlagShort, charsetFlag), "Set the charset used for password generation. Default charset is '0123456789'.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <NUM>", offsetFlagShort, offsetFlag), "Set the password generation offset. Default value is 0. This can come in handy if, for example, you know that the first 3 digits of a passphrase are '123' and you want to speed up the cracking process. Let's say that the passphrase is '123_____' (8 digits). Running alohomora with offset 12300000 will significantly reduce the amount of passphrases required to try before succeeding. This works with letters and symbols as well, of course.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <NUM>", jobsizeFlagShort, jobsizeFlag), "Set the amount of passwords that each client should attempt in each job. Default value is 10.000. Use this if you have high or low performance clients in order to improve cracking speeds.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <SECONDS>", timeoutFlagShort, timeoutFlag), "Set the timeout for jobs in seconds. Default timeout is 10 minutes (600 seconds) After a job has been passed to a client, the timeout counter will start. If the client fails to finish the job within that timeframe, the server will deem the job lost and reschedule it.\n")
	t.AddRow(fmt.Sprintf("  --%s <NUM>", queueSizeFlag), "Set the job generation backlog size. Default size is 50. You probably don't need to adjust it, but feel free to do so.\n")
	t.AddRow(fmt.Sprintf("  --%s <NUM>", maxJobsFlag), "Set the maximum amount of jobs to generate for the task at hand. Default value is 0 (= no limit). Depending on the passphrase to crack, A LOT of jobs can be generated over time. This option can limit this amount to a fixed number.\n")
	t.AddRow(fmt.Sprintf("  --%s <SECONDS>", maxTimeFlag), "Set the maximum amount of seconds the server will be running. Default is 0 (= no limit). Maybe you want to give cracking a particular passphrase a shot, but computing time is limited. Set the server lifetime to e.g. 3600 (= 1 hour) for a quick go.\n")
	t.AddRow(fmt.Sprintf("  --%s <FILE>", reportXMLFlag), "Generate an XML report upon server termination (write it to FILE).\n")
	t.AddRow(fmt.Sprintf("  --%s <FILE>", reportJSONFlag), "Generate a JSON report upon server termination (write it to FILE).\n")
	t.AddRow(fmt.Sprintf("  --%s", restFlag), "Enable alohomora's integrated REST server (highly experimental!). Default is disabled\n")
	t.AddRow(fmt.Sprintf("  --%s <IP>", restAddressFlag), "Make alohomora's integrated REST server listen on this address (highly experimental!). Default is 127.0.0.1\n")
	t.AddRow(fmt.Sprintf("  --%s <PORT>", restPortFlag), "Make alohomora's integrated REST server listen on this port (highly experimental!). Default is 29101\n")
	fmt.Println(t)

	clientIntro()
	t = uitable.New()
	t.MaxColWidth = 60
	t.Wrap = true // wrap columns
	t.AddRow(fmt.Sprintf("  -%s / --%s <IP>", serverFlagShort, serverFlag), "Set the address the client should connect to.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <PORT>", portFlagShort, portFlag), "Set the port to connect to. Defaults to 29100.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <PORT>", attemptsFlagShort, attemptsFlag), "Set how often the client will try to connect to the server before giving up. The client will wait for 10 seconds after each try. Default value is 5.\n")
	fmt.Println(t)

}

// Parse parses all command line parameters provided and encapsulates them in an
// instance of the Options type for easier handling. That instance will be returned
// if no error occurs.
func Parse() (*Options, error) {
	args := Options{}

	flag.BoolVar(&args.Server, nodeFlag, nodeDefault, nodeHelp)

	flag.StringVar(&args.Mode, modeFlag, modeFlagDefault, "undefined")
	flag.StringVar(&args.Mode, modeFlagShort, modeFlagDefault, "undefined")

	flag.UintVar(&args.Port, portFlag, portFlagDefault, portFlagHelp)
	flag.UintVar(&args.Port, portFlagShort, portFlagDefault, portFlagHelp)

	flag.StringVar(&args.Host, serverFlag, serverFlagDefault, serverFlagHelp)
	flag.StringVar(&args.Host, serverFlagShort, serverFlagDefault, serverFlagHelp)

	flag.BoolVar(&args.Unfancy, unfancyFlag, unfancyFlagDefault, unfancyFlagHelp)
	flag.BoolVar(&args.Unfancy, unfancyFlagShort, unfancyFlagDefault, unfancyFlagHelp)

	flag.BoolVar(&args.Verbose, verboseFlag, verboseFlagDefault, verboseFlagHelp)
	flag.BoolVar(&args.Verbose, verboseFlagShort, verboseFlagDefault, verboseFlagHelp)

	flag.StringVar(&args.Jobsize, jobsizeFlag, jobsizeDefault, jobsizeHelp)
	flag.StringVar(&args.Jobsize, jobsizeFlagShort, jobsizeDefault, jobsizeHelp)

	flag.StringVar(&args.Charset, charsetFlag, charsetFlagDefault, charsetFlagHelp)
	flag.StringVar(&args.Charset, charsetFlagShort, charsetFlagDefault, charsetFlagHelp)

	flag.UintVar(&args.Passlen, lengthFlag, lengthFlagDefault, lengthFlagHelp)
	flag.UintVar(&args.Passlen, lengthFlagShort, lengthFlagDefault, lengthFlagHelp)

	flag.StringVar(&args.Offset, offsetFlag, offsetFlagDefault, offsetFlagHelp)
	flag.StringVar(&args.Offset, offsetFlagShort, offsetFlagDefault, offsetFlagHelp)

	flag.StringVar(&args.Target, targetFlag, targetFlagDefault, targetFlagHelp)
	flag.StringVar(&args.Target, targetFlagShort, targetFlagDefault, targetFlagHelp)

	flag.Uint64Var(&args.Timeout, timeoutFlag, timeoutFlagDefault, timeoutFlagHelp)
	flag.Uint64Var(&args.Timeout, timeoutFlagShort, timeoutFlagDefault, timeoutFlagHelp)

	flag.Uint64Var(&args.QueueSize, queueSizeFlag, queueSizeFlagDefault, queueSizeFlagHelp)

	flag.StringVar(&args.MaxJobs, maxJobsFlag, maxJobsFlagDefault, maxJobsFlagHelp)

	flag.Uint64Var(&args.MaxTime, maxTimeFlag, maxTimeFlagDefault, maxTimeFlagHelp)

	flag.StringVar(&args.ReportXMLTarget, reportXMLFlag, reportXMLFlagDefault, "If provided, an XML report will be generated")
	flag.StringVar(&args.ReportJSONTarget, reportJSONFlag, reportJSONFlagDefault, "If provided, a JSON report will be generated")

	flag.BoolVar(&args.EnableREST, restFlag, restFlagDefault, "If set, a REST server is started on port 29100")

	flag.StringVar(&args.RESTAddress, restAddressFlag, restAddressFlagDefault, "Set REST listening address")
	flag.UintVar(&args.RESTPort, restPortFlag, restPortFlagDefault, "Set REST port")

	flag.UintVar(&args.ConnectionAttempts, attemptsFlag, attemptsFlagDefault, attemptsFlagHelp)
	flag.UintVar(&args.ConnectionAttempts, attemptsFlagShort, attemptsFlagDefault, attemptsFlagHelp)

	flag.Parse()

	return &args, args.validate()
}

func (opts Options) validate() error {
	if opts.Server {
		if len(opts.Host) == 0 {
			return errors.New("A valid listening address is required")
		}

		if opts.Port < 0 || opts.Port > 65535 {
			return errors.New("A valid port number is required")
		}

		if opts.EnableREST && (opts.RESTPort < 0 || opts.RESTPort > 65535 || opts.RESTPort == opts.Port) {
			return errors.New("REST port is invalid")
		}

		if bigint.LessThan(bigint.ToBigInt(opts.Offset), big.NewInt(0)) {
			return errors.New("Offset must be a positive number or 0")
		}

		if opts.Passlen == 0 {
			return errors.New("Minimum password length is 1")
		}

		if bigint.LTE(bigint.ToBigInt(opts.Jobsize), big.NewInt(0)) {
			return errors.New("Jobsize must be a positive number")
		}

		if len(opts.Charset) == 0 {
			return errors.New("Charset must contain at least one character")
		}

		if _, err := os.Stat(opts.Target); os.IsNotExist(err) {
			// path/to/whatever does not exist
			return errors.New("Target .pcap file must exist")
		}

	} else {

		// Everything else is a client
		if len(opts.Host) == 0 || opts.Host == "0.0.0.0" {
			return errors.New("A server's IP address or hostname is required")
		}

		//65535
		if opts.Port < 0 || opts.Port > 65535 {
			return errors.New("A valid port number is required")
		}
	}

	return nil
}
