package opts

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/steps0x29a/alohomora/term"
)

func intro() {
	Banner(true, true)
	fmt.Printf("Usage: \n\n ./alohomora --server --target <FILE> [options]\n ./alohomora --ip <IP> --port <PORT>\n\n")
}

func serverIntro() {
	fmt.Println(term.Bold("SERVER MODE USAGE"))
	fmt.Println()
	fmt.Println("  ./alohomora --server")
	fmt.Println()
	t := uitable.New()
	t.MaxColWidth = 80
	t.Wrap = true // wrap columns
	t.AddRow("", "Runs alohomora in server mode on 0.0.0.0:29100 with the default target (which is 'handshake.pcap' in the current working directory). The character set used to generate passwords will be 0123456789. Each password will be 8 characters long and each client will be tasked with 10.000 passwords per job.")
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
	t.AddRow("", "Runs alohomora in client mode, connecting to <IP> on <PORT> and waiting for new jobs. By default, the client will try to connect to the server 5 times (in intervals of 10 seconds) before giving up.")
	fmt.Println(t)
	fmt.Println()
	fmt.Println(term.Bold("CLIENT OPTIONS"))
	fmt.Println()
}

// Usage prints a custom usage message. This function is passed to
// flag.Usage in order to automatically print usage information when needed.
func Usage() {
	intro()
	serverIntro()

	t := uitable.New()
	t.MaxColWidth = 60
	t.Wrap = true // wrap columns

	t.AddRow(fmt.Sprintf("  -%s / --%s <IP>", serverFlagShort, serverFlag), "Set the server's listen address. Defaults to 0.0.0.0.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <PORT>", portFlagShort, portFlag), "Set the server's listen port. Defaults to 29100.\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <MODE>", modeFlagShort, modeFlag), "Set the mode of operation. Default is 'WPA2'.\n\nCurrently supported modes:\n\n  - WPA2 [WPA2 handshake cracking] (default)\n\n")
	t.AddRow(fmt.Sprintf("  -%s / --%s <PATH>", targetFlagShort, targetFlag), "Set the cracking target (path to a file). By default, alohomora server will use 'handshake.pcap' in the current working directory. If that file is not present, alohomora will wait for it.\n")
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
	t.AddRow(fmt.Sprintf("  -%s / --%s", forceCharsetFlagShort, forceCharsetFlag), "Force charset, don't scramble and don't clean it up. This is useful if a server is restarted with a different offset (scrambling and cleaning mess up the charset's order, so a lot of passwords might be lost if they are applied and the server is restarted with a different offset)")
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
