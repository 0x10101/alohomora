package opts

import (
	"errors"
	"flag"
	"fmt"
	"math/big"

	"github.com/steps0x29a/alohomora/bigint"
	"github.com/steps0x29a/alohomora/term"
)

//The Options type wraps all command line options in a neat struct for easier handling.
type Options struct {
	Server             bool
	Port               uint   `json:"port"`
	Host               string `json:"ip"`
	Mode               string `json:"mode"`
	Verbose            bool   `json:"verbose"`
	Unfancy            bool   `json:"unfancy"`
	Charset            string `json:"charset"`
	Jobsize            string `json:"jobsize"`
	Passlen            uint   `json:"password_length"`
	Offset             string `json:"offset"`
	Target             string `json:"target"`
	Timeout            uint64 `json:"timeout"`
	ReportXMLTarget    string `json:"report_xml_target"`
	ReportJSONTarget   string `json:"report_json_target"`
	QueueSize          uint64 `json:"queue_size"`
	MaxJobs            string `json:"maxjobs"`
	MaxTime            uint64 `json:"jobtimeout"`
	EnableREST         bool
	RESTAddress        string
	RESTPort           uint
	ConnectionAttempts uint
	ForceCharset       bool `json:"force_charset"`
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
	targetFlagDefault = "handshake.pcap"
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

	forceCharsetFlag        = "force-charset"
	forceCharsetFlagShort   = "f"
	forceCharsetFlagDefault = false
)

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

	flag.BoolVar(&args.ForceCharset, forceCharsetFlag, forceCharsetFlagDefault, "")
	flag.BoolVar(&args.ForceCharset, forceCharsetFlagShort, forceCharsetFlagDefault, "")

	flag.Parse()

	// Use --unfancy option if present
	if args.Unfancy {
		term.NoColors()
	}

	// Cleanup the charset
	if !args.ForceCharset {
		args.Charset = string(CleanupCharset([]rune(args.Charset)))
	}

	return &args, args.validate()
}

func validPort(port uint) bool {
	return port > 0 && port <= 65535
}

func validateRESTOptions(opts Options) error {
	if !opts.EnableREST {
		return nil
	}

	if !validPort(opts.RESTPort) {
		return errors.New("A valid REST port is required")
	}

	if opts.RESTPort == opts.Port {
		return errors.New("REST port and listening port must not be the same")
	}

	return nil
}

func (opts Options) validate() error {
	if opts.Server {
		if len(opts.Host) == 0 {
			return errors.New("A valid listening address is required")
		}

		if !validPort(opts.Port) {
			return errors.New("A valid port number is required")
		}

		if err := validateRESTOptions(opts); err != nil {
			return err
		}

		if bigint.Lt(bigint.ToBigInt(opts.Offset), big.NewInt(0)) {
			return errors.New("Offset must be a positive number or 0")
		}

		if opts.Passlen == 0 {
			return errors.New("Minimum password length is 1")
		}

		if bigint.LtE(bigint.ToBigInt(opts.Jobsize), big.NewInt(0)) {
			return errors.New("Jobsize must be a positive number")
		}

		if len(opts.Charset) == 0 {
			return errors.New("Charset must contain at least one character")
		}

		if len(opts.Target) == 0 {
			return errors.New("A target is required")
		}

		if opts.Mode != "WPA2" && opts.Mode != "MD5" {
			return fmt.Errorf("Unknown job type: %s", opts.Mode)
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
