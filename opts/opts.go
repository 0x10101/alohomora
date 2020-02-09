package opts

import (
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/cheynewallace/tabby"
	"github.com/steps0x29a/alohomora/bigint"
	"github.com/steps0x29a/alohomora/term"
)

//The Options type wraps all command line options in a neat struct for easier handling.
type Options struct {
	Server             bool
	Port               uint
	Host               string
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
	ConnectionAttempts uint
}

const (
	maxWidth = 80

	nodeFlag    = "server"
	nodeDefault = false
	nodeHelp    = "If provided, alohomora will run in server mode. If not provided, it will run in client mode"

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
	queueSizeFlagShort   = "q"
	queueSizeFlagDefault = 50
	queueSizeFlagHelp    = "Amount of crackjobs to generate as a backlog. Defaults to 50. Must neither be negative nor 0."

	maxJobsFlag        = "maxjobs"
	maxJobsFlagShort   = "m"
	maxJobsFlagDefault = "0"
	maxJobsFlagHelp    = "Maximum amount of jobs to dispatch to clients for current handshake. Defaults to 0 (disabled)."

	maxTimeFlag        = "maxtime"
	maxTimeFlagDefault = 0
	maxTimeFlagHelp    = "Amount of seconds before the server considers the current handshake as a failure. Defaults to 0 (disabled)."

	attemptsFlag        = "attempts"
	attemptsFlagShort   = "a"
	attemptsFlagDefault = 5
	attemptsFlagHelp    = "Number of connection attempts to a server (default is 5)."
)

/*func printOpt(short, long, help, def string) {
	head := fmt.Sprintf("  -%s / --%s   ", short, long)
	defP := fmt.Sprintf("(default: %s)", def)
	//ldef := len(defP)
	lhead := len(head)
	lhelp := len(help)
	fmt.Printf("%s", head)
	if lhead+lhelp < maxWidth {
		// Simply print
		fmt.Printf("%s\n", help)
	} else {
		// Loop

	}
}*/

func intro() {
	fmt.Printf("\nUsage: ./alohomora [--server] [options]\n\n")
}

func serverIntro() {
	fmt.Println(term.Bold("SERVER MODE USAGE"))
	fmt.Println()
	fmt.Println("  ./alohomora --server --target /path/to/pcap/file")
	fmt.Println()
	fmt.Println("  Runs alohomora server on 0.0.0.0:29100 targeting the pcap file")
	fmt.Println("  in /path/to/pcap/file. The character set used to generate passwords")
	fmt.Println("  will be 0123456789. Each password will be 8 characters long and each")
	fmt.Println("  client will be tasked with 10.000 passwords per job.")
	fmt.Println()
	fmt.Println(term.Bold("SERVER OPTIONS"))
	fmt.Println()
}

func Usage() {
	/*t := tabby.New()
	t.AddHeader("NAME", "TITLE", "DEPARTMENT")
	t.AddLine("John Smith", "Developer", "Engineering")
	t.Print()*/
	intro()
	serverIntro()
	t := tabby.New()

	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(serverFlagShort), term.Bold(serverFlag)), "Set the server's listen address")
	t.AddLine("", "Default value is 0.0.0.0.")

	//, fmt.Sprintf("(default: %s)", serverFlagDefault))
	/*t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(portFlagShort), term.Bold(portFlag)), "Set the server's listen port", fmt.Sprintf("(default: %d)", portFlagDefault))
	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(charsetFlagShort), term.Bold(charsetFlag)), "Set the charset used for password generation", fmt.Sprintf("(default: %s)", charsetFlagDefault))
	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(offsetFlagShort), term.Bold(offsetFlag)), "Set the password generation offset", fmt.Sprintf("(default: %s)", offsetFlagDefault))
	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(timeoutFlagShort), term.Bold(timeoutFlag)), "Set the timeout for jobs in seconds", fmt.Sprintf("(default: %d)", timeoutFlagDefault))
	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(queueSizeFlagShort), term.Bold(queueSizeFlag)), "Set the job generation backlog size", fmt.Sprintf("(default: %d)", queueSizeFlagDefault))
	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold(maxJobsFlagShort), term.Bold(maxJobsFlag)), "Set the maximum amount of jobs to generate for the task at hand", fmt.Sprintf("(default: %s)", maxJobsFlagDefault))
	t.AddLine(fmt.Sprintf("  -%s / --%s", term.Bold("k"), term.Bold(maxTimeFlag)), "Set the maximum amount of seconds the server will be running ", fmt.Sprintf("(default: %d)", maxTimeFlagDefault))
	*/
	t.Print()

}

func Usage2() {
	fmt.Println("")
	fmt.Println("Usage: ./alohomora [--server] [options]")
	fmt.Println()
	fmt.Println(term.Bold("SERVER MODE USAGE"))
	fmt.Println()
	fmt.Println("  ./alohomora --server --target /path/to/pcap/file")
	fmt.Println()
	fmt.Println("  Runs alohomora server on 0.0.0.0:29100 targeting the pcap file")
	fmt.Println("  in /path/to/pcap/file. The character set used to generate passwords")
	fmt.Println("  will be 0123456789. Each password will be 8 characters long and each")
	fmt.Println("  client will be tasked with 10.000 passwords per job.")
	fmt.Println()
	fmt.Println(term.Bold("SERVER OPTIONS"))
	fmt.Println()
	fmt.Printf("  -%s / --%s        Set the listen address   (default: %s)\n", term.Bold(serverFlagShort), term.Bold(serverFlag), serverFlagDefault)
	fmt.Printf("  -%s / --%s      Set the listen port      (default: %d)\n", term.Bold(portFlagShort), term.Bold(portFlag), portFlagDefault)
	fmt.Printf("  -%s / --%s   Set the character set    (default: %s)\n", term.Bold(charsetFlagShort), term.Bold(charsetFlag), charsetFlagDefault)
	fmt.Printf("  -%s / --%s    Set the password length  (default: %d)\n", term.Bold(lengthFlagShort), term.Bold(lengthFlag), lengthFlagDefault)
	fmt.Printf("  -%s / --%s    Set the generation offset  (default: %s)\n", term.Bold(offsetFlagShort), term.Bold(offsetFlag), offsetFlagDefault)
	fmt.Printf("  -%s / --%s    Set the job timeout  (default: %d)\n", term.Bold(timeoutFlagShort), term.Bold(timeoutFlag), timeoutFlagDefault)
	fmt.Printf("  -%s / --%s    Set the job backlog size  (default: %d)\n", term.Bold(queueSizeFlagShort), term.Bold(queueSizeFlag), queueSizeFlagDefault)
	fmt.Printf("  -%s / --%s    Set the maximum amount of jobs  (default: %s)\n", term.Bold(maxJobsFlagShort), term.Bold(maxJobsFlag), maxJobsFlagDefault)

}

// Parse parses all command line parameters provided and encapsulates them in an
// instance of the Options type for easier handling. That instance will be returned
// if no error occurs.
func Parse() (*Options, error) {
	args := Options{}

	flag.BoolVar(&args.Server, nodeFlag, nodeDefault, nodeHelp)

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
	flag.Uint64Var(&args.QueueSize, queueSizeFlagShort, queueSizeFlagDefault, queueSizeFlagHelp)

	flag.StringVar(&args.MaxJobs, maxJobsFlag, maxJobsFlagDefault, maxJobsFlagHelp)
	flag.StringVar(&args.MaxJobs, maxJobsFlagShort, maxJobsFlagDefault, maxJobsFlagHelp)

	flag.Uint64Var(&args.MaxTime, maxTimeFlag, maxTimeFlagDefault, maxTimeFlagHelp)

	flag.StringVar(&args.ReportXMLTarget, "oX", "", "If provided, an XML report will be generated")
	flag.StringVar(&args.ReportJSONTarget, "oJ", "", "If provided, a JSON report will be generated")

	flag.BoolVar(&args.EnableREST, "rest", false, "If set, a REST server is started on port 29100")

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
