package core

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/steps0x29a/alohomora/report"
	"github.com/steps0x29a/alohomora/rest"

	"github.com/steps0x29a/alohomora/handshakes"

	"github.com/steps0x29a/alohomora/jobs"
	"github.com/steps0x29a/alohomora/msg"

	uuid "github.com/satori/go.uuid"

	"github.com/steps0x29a/alohomora/bigint"
	"github.com/steps0x29a/alohomora/bytes"
	"github.com/steps0x29a/alohomora/opts"
	"github.com/steps0x29a/alohomora/term"
)

// A Server manages clients and jobs.
type Server struct {
	sync.Mutex

	Clients            map[*Client]bool
	StartedJobs        *big.Int
	TotalJobs          *big.Int
	FinishedJobs       *big.Int
	Queue              chan *jobs.CrackJob
	freeClients        chan *Client
	Terminated         chan bool
	register           chan *Client
	unregister         chan *Client
	Errors             chan error
	Pending            map[*Client]*jobs.CrackJob
	generationFinished bool
	maximumJobsReached bool
	verbose            bool
	timeout            uint64
	ReportFile         string
	queuesize          uint64
	maxjobs            *big.Int
	taskTimeout        uint64
	started            time.Time
	report             *report.Report
}

func newServer(opts *opts.Options) *Server {
	if opts.QueueSize <= 0 {
		term.Problem("Invalid queue size detected, defaulting to 1\n")
		opts.QueueSize = 1
	}

	server := &Server{
		Clients:            make(map[*Client]bool),
		StartedJobs:        big.NewInt(0),
		TotalJobs:          big.NewInt(0),
		FinishedJobs:       big.NewInt(0),
		Queue:              make(chan *jobs.CrackJob, opts.QueueSize),
		freeClients:        make(chan *Client),
		Terminated:         make(chan bool),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		Errors:             make(chan error),
		Pending:            make(map[*Client]*jobs.CrackJob),
		generationFinished: false,
		maximumJobsReached: false,
		verbose:            opts.Verbose,
		timeout:            opts.Timeout,
		maxjobs:            bigint.ToBigInt(opts.MaxJobs),
		taskTimeout:        opts.MaxTime,
		started:            time.Now(),
		report:             report.New(),
	}

	return server

}

func showOpts(opts *opts.Options) {
	if !opts.Verbose {
		return
	}
	term.Info("%s -> %s\n", term.LabelMagenta("ADR"), term.LabelGreen(fmt.Sprintf("%s:%d", opts.Host, opts.Port)))
	term.Info("%s -> %s\n", term.LabelMagenta("TGT"), term.LabelGreen(opts.Target))
	term.Info("%s -> %s\n", term.LabelMagenta("CST"), term.LabelGreen(opts.Charset))
	term.Info("%s -> %s\n", term.LabelMagenta("JOB"), term.LabelGreen(fmt.Sprintf("%s", opts.Jobsize)))
	term.Info("%s -> %s\n", term.LabelMagenta("LEN"), term.LabelGreen(fmt.Sprintf("%d", opts.Passlen)))
	term.Info("%s -> %s\n", term.LabelMagenta("OFF"), term.LabelGreen(fmt.Sprintf("%s", opts.Offset)))
	fmt.Println("")
}

func (server *Server) onClientConnected(client *Client) {
	server.Lock()
	defer server.Unlock()
	term.Info("Client connected: %s\n", term.BrightBlue(client.FullID()))
	clientCount := uint(len(server.Clients))
	if clientCount > server.report.MaxClientsConnected {
		server.report.MaxClientsConnected = clientCount
	}
}

func (server *Server) onClientLeft(client *Client) {
	term.Info("Client left: %s\n", term.Red(client.FullID()))
}

func (server *Server) checkErrors() {
	for {
		err := <-server.Errors
		term.Error("Server error: %s\n", term.BrightRed(err.Error()))
	}
}

func (server *Server) onProgress() {
	/*total := server.TotalJobs
	finished := term.Reverse(term.InsertAfterEvery(term.Reverse(server.FinishedJobs.String()), '.', 3))
	pending := len(server.Pending)
	clients := len(server.Clients)

	percent := bigint.Percent(total, server.FinishedJobs)
	numStrTotal := term.Reverse(term.InsertAfterEvery(term.Reverse(total.String()), '.', 3))
	*/
	/*if server.verbose {
		term.Info("Progress: %s/%s (%0.2f%%, %d clients connected, %d jobs pending)\n", finished, numStrTotal, percent, clients, pending)
	}*/

}

func (server *Server) loop() {
	for {
		select {
		case client := <-server.register:
			{
				server.Lock()
				server.Clients[client] = false
				server.Unlock()
				go server.receive(client)
				server.onClientConnected(client)
			}

		case client := <-server.unregister:
			{
				server.Lock()
				delete(server.Clients, client)
				server.Unlock()
				server.onClientLeft(client)
			}
		}
	}
}

func (server *Server) accept(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			term.Problem("Unable to accept new client: %s\n", err)
			continue
		}

		clientID, _ := uuid.NewV4()
		client := newClient(connection)
		client.ID = clientID
		client.connected = time.Now()
		server.register <- client
	}
}

func (server *Server) send(client *Client, message *msg.Message) {
	data, err := message.Encode()
	if err != nil {
		server.Errors <- err
		return
	}

	_, err = client.Socket.Write(data)
	// TODO: Handle incomplete writes
	if err != nil {
		server.Errors <- err
		return
	}

	_, err = client.Socket.Write(AlohomoraSuffix)
	if err != nil {
		server.Errors <- err
		return
	}
}

func (server *Server) onClientHello(client *Client, message *msg.Message) {
	// Ack the client
	server.Lock()
	defer server.Unlock()

	answer := msg.NewMessage(msg.Ack, nil)
	server.Clients[client] = true
	go server.send(client, answer)
}

func (server *Server) onClientIdle(client *Client, message *msg.Message) {
	if !server.Clients[client] {
		// Invalid
		term.Warn("Closing connection to %s\n", client)
		server.kick(client)
	} else {
		// Schedule client for work
		server.freeClients <- client
	}
}

func (server *Server) onClientResponse(client *Client, message *msg.Message) {
	server.Lock()

	result, err := jobs.DecodeResult(message.Payload)
	job := server.Pending[client]
	server.report.PasswordsTried = bigint.Add(server.report.PasswordsTried, big.NewInt(job.Gen.Amount))
	if client.tried == nil {
		client.tried = big.NewInt(0)
	}
	client.tried = bigint.Add(client.tried, big.NewInt(job.Gen.Amount))
	delete(server.Pending, client)
	server.FinishedJobs = bigint.Add(server.FinishedJobs, big.NewInt(1))
	server.report.FinishedRuns = bigint.Cp(server.FinishedJobs)
	server.Unlock()

	if err != nil {
		server.Errors <- err
		term.Error("Unable to decode result: %s\n", err)
	} else {
		client.finished++
		if result.Success {
			term.Success("Client %s cracked the password: %s\n", term.BrightBlue(client.ShortID()), term.LabelGreen(result.Payload))

			// As of now we can safely assume that the payload is a WPA2 payload
			wpaPayload, err := job.DecodeWPA2()
			if err != nil {
				server.Errors <- fmt.Errorf(fmt.Sprintf("Unable to decode job %s's payload as WPA2 payload: %s", job.ID.String()[:8], err.Error()))
			}
			server.onClientSuccess(client, fmt.Sprintf("%s %s", wpaPayload.ESSID, wpaPayload.BSSID), result.Payload)

			//server.Terminated <- true
			server.Terminate()
		} else {
			term.Info("Client %s %s to crack %s\n", term.BrightBlue(client.ShortID()), term.BrightRed("failed"), term.Cyan(result.JobID.String()[:8]))
		}
	}

	server.onClientIdle(client, message)
}

func (server *Server) onClientSuccess(client *Client, username, password string) {
	server.Lock()
	defer server.Unlock()
	server.report.Success = true
	server.report.AccessData.Username = username
	server.report.AccessData.Password = password
	server.report.SuccessClientAddress = client.Socket.RemoteAddr()
	server.report.SuccessClientID = client.ShortID()
}

// Terminate terminates the server, saves current time to report
func (server *Server) Terminate() {
	server.Lock()
	defer server.Unlock()
	server.report.EndTimestamp = time.Now()
	server.Terminated <- true

}

func (server *Server) onClientError(client *Client, message *msg.Message) {
	// Payload should be CrackJobResult

	result, err := jobs.DecodeResult(message.Payload)
	defer server.kick(client)

	if err != nil {
		server.Errors <- err
		term.Error("Client %s crashed with invalid error message\n", term.BrightBlue(client.ShortID()))
	} else {
		term.Error("Client %s crashed during %s\n", term.BrightBlue(client.ShortID()), term.Cyan(result.JobID.String()[:8]))
	}
}

// Report returns the server's report data at the current time
// (might be incomplete before all jobs have finished running)
func (server *Server) Report() *report.Report {
	return server.report
}

func (server *Server) handle(client *Client, message *msg.Message) {
	switch message.Type {
	case msg.Hello:
		{
			server.onClientHello(client, message)
			break
		}
	case msg.Idle:
		{
			server.onClientIdle(client, message)
			break
		}
	case msg.Finished:
		{
			server.onClientResponse(client, message)
			break
		}
	case msg.ClientError:
		{
			server.onClientError(client, message)
			break
		}
	}
}

func (server *Server) receive(client *Client) {
	var buffer = make([]byte, bufferSize)

	for {
		var message = make([]byte, 0)
		var size = 0

		for {
			read, err := client.Socket.Read(buffer)
			if read == 0 || err != nil && err != io.EOF {
				// Connection lost
				server.unregister <- client
				return
			}

			message = append(message, buffer[:read]...)
			size += read

			if bytes.EndsWith(message, AlohomoraSuffix) {

				decoded, err := msg.Decode(message[:size-len(AlohomoraSuffix)])
				if err != nil {
					term.Error("Unable to decode client message: %s\n", term.BrightRed(fmt.Sprintf("%s", err)))
				} else {
					go server.handle(client, decoded)
				}

				break
			}

		}

	}
}

func (server *Server) broadcast(message *msg.Message) {
	for client := range server.Clients {
		server.send(client, message)
	}
}

func (server *Server) kick(client *Client) {
	term.Info("Kicking %s...\n", term.BrightBlue(client.ShortID()))
	leaveMessage := msg.NewMessage(msg.Leave, make([]byte, 0))
	server.send(client, leaveMessage)
	defer client.Socket.Close()
}

// KickAll asks all clients to leave, closing their connections as well.
func (server *Server) KickAll() {
	for client := range server.Clients {
		server.kick(client)
	}
}

func (server *Server) dispatch() {
	for {
		select {
		case <-server.Terminated:
			{
				return
			}
		default:
			{
				// Server not yet terminated, dispatch jobs

				client, _ := <-server.freeClients

				job, _ := <-server.Queue

				// We need the payload should this client succeed. So save it before omitting it
				tmp := job.Payload

				if client.finished > 0 {
					// Assume client already knows the handshake
					job.Payload = nil
				}

				payload, err := job.Encode()
				if err != nil {
					//server.Terminated <- true
					server.Terminate()
					return
				}

				// Give the job its payload back
				job.Payload = tmp

				job.Started = time.Now()
				server.Lock()
				server.Pending[client] = job

				// Send job to client
				message := msg.NewMessage(msg.Task, payload)
				server.Unlock()

				if server.verbose {
					term.Info("Client %s %s with job %s\n", term.BrightBlue(client.ShortID()), term.BrightMagenta("tasked"), term.Cyan(job.ID.String()[:8]))
				}

				go server.send(client, message)
				client.assigned++
			}
		}

	}
}

func (server *Server) initCrackjobs(opts *opts.Options) {
	var filepath = opts.Target
	var handshake = handshakes.NewHandshake()
	err := handshake.Read(filepath)
	if err != nil {
		// This is bad
		term.Error("Unable to process target: %s\n", err)
		server.Terminate()
		return
	}

	charset := []rune(opts.Charset)
	length := int64(opts.Passlen)
	jobsize := bigint.ToBigInt(opts.Jobsize)
	offset := bigint.ToBigInt(opts.Offset)
	maxValue := bigint.Sub(bigint.Pow(big.NewInt(int64(len(charset))), big.NewInt(length)), offset)
	runs := bigint.Div(maxValue, jobsize)

	mod := bigint.Mod(maxValue, jobsize)
	if !bigint.Eq(mod, big.NewInt(0)) {
		runs = bigint.Add(runs, big.NewInt(1))
	}

	server.TotalJobs = bigint.Cp(runs)

	var jobindex *big.Int = big.NewInt(0)
	var remaining *big.Int = bigint.Cp(maxValue)

	if bigint.Lt(remaining, big.NewInt(0)) {
		term.Error("Invalid offset: %s\n", offset)
		//server.Terminated <- true
		server.Terminate()
		return
	}

	// While remaining > 0
	for bigint.Gt(remaining, big.NewInt(0)) {

		var runAmount *big.Int = bigint.Cp(jobsize)
		if bigint.Lt(remaining, jobsize) {
			runAmount.Set(remaining)
		}
		remaining = bigint.Sub(remaining, runAmount)

		var calcOffset = bigint.Add(offset, bigint.Mul(jobsize, jobindex))

		//var endOffset = bigint.Sub(bigint.Add(calcOffset, runAmount), big.NewInt(1))
		//first, _ := gen.GeneratePassword(charset, length, calcOffset)
		//last, _ := gen.GeneratePassword(charset, length, endOffset)

		job, err := jobs.NewWPA2Job(
			handshake.Data,
			charset,
			length,
			calcOffset,
			runAmount.Int64(),
			handshake.ESSID,
			handshake.BSSID,
		)

		if err != nil {
			term.Error("Unable to generate Crackjob: %s\n", err)
			server.Terminate()
			return
		}

		jobindex = bigint.Add(jobindex, big.NewInt(1))
		/*if server.verbose {
			term.Info("Generated Crackjob %s (%s - %s)\n", term.Cyan(job.String()), term.BrightBlue(first), term.BrightBlue(last))
		}*/

		server.Queue <- job
		if bigint.GtE(jobindex, server.maxjobs) && bigint.Gt(server.maxjobs, big.NewInt(0)) {
			term.Info("Maximum amount of jobs reached, stopping job generation\n")
			server.maximumJobsReached = true
			break
		}

		if server.taskTimeoutReached() {
			term.Info("Maximum time for task reached, stopping job generation\n")
			break
		}
	}

	server.generationFinished = true
}

func (server *Server) checkPending() {
	ticker := time.NewTicker(time.Second * 3)
	for {
		<-ticker.C
		for client, job := range server.Pending {
			if _, ok := server.Clients[client]; !ok {
				if server.verbose {
					term.Warn("Missing client %s - Rescheduling job %s\n", term.BrightBlue(client.ShortID()), term.Cyan(job.String()))
				}
				server.Queue <- job
				server.Lock()
				delete(server.Pending, client)
				server.Unlock()
			}

			// Test for job timeout
			now := time.Now()
			dur := now.Sub(job.Started)
			if dur.Seconds() > float64(server.timeout) {
				// Job timed out, kick client
				term.Warn("Client %s timed out on %s, kicking them\n", term.BrightBlue(client.ShortID()), term.Cyan(job.ID.String()[:8]))
				go server.kick(client)
			}
		}
	}
}

func (server *Server) updateProgress() {
	ticker := time.NewTicker(time.Millisecond * 2000)
	for {
		<-ticker.C
		if server.generationFinished && bigint.GtE(server.FinishedJobs, server.TotalJobs) {
			server.onProgress()
			term.Info("All jobs finished, terminating...\n")
			server.Terminate()
			//server.Terminated <- true
		} else if server.maximumJobsReached && bigint.GtE(server.FinishedJobs, server.maxjobs) {
			server.onProgress()
			term.Info("Maximum number of jobs reached, terminating...\n")
			//server.Terminated <- true
			server.Terminate()
		} else {
			server.onProgress()
		}

		// Check timeout
		if server.taskTimeoutReached() {
			// Task timed out, kill everything
			term.Info("Task timeout reached, terminating server...\n")
			//server.Terminated <- true
			server.Terminate()
		}
	}
}

func (server *Server) taskTimeoutReached() bool {
	return server.taskTimeout > 0 && time.Now().Sub(server.started).Seconds() > float64(server.taskTimeout)
}

// SlashRoot is the root endpoint of the REST API (Server implements the RESTHandler interface).
func (server *Server) SlashRoot(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hi there, this is alohomora!")
}

// ClientsHandleFunc wraps all connected clients in ClientInfo objects and marshals that
// information to JSON. Server implements RESTHandler interface.
func (server *Server) ClientsHandleFunc(res http.ResponseWriter, req *http.Request) {
	clients := make([]*ClientInfo, 0)
	for client, connected := range server.Clients {
		if connected {
			clients = append(clients, client.Info())
		}
	}

	data, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		term.Error("Unable to marshal clients to JSON: %s\n", err)
	} else {
		fmt.Fprint(res, string(data))
	}
}

// PendingJobsHandleFunc wraps all pending jobs in JobInfo objects and marshals that info to JSON.
// Server implements RESTHandler interface.
func (server *Server) PendingJobsHandleFunc(res http.ResponseWriter, req *http.Request) {
	mapping := make(map[string]*jobs.CrackJobInfo)
	for client, job := range server.Pending {
		mapping[client.ShortID()] = job.Info()
	}

	data, err := json.MarshalIndent(mapping, "", "  ")
	if err != nil {
		term.Error("Unable to marshal pending jobs to JSON: %s\n", err)
	} else {
		fmt.Fprint(res, string(data))
	}
}

func waitForTarget(target string, found chan bool) {
	for {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			time.Sleep(time.Second * 2)
			continue
		}

		// Found
		found <- true
		break
	}
}

// Serve builds a new Server instance and starts listening on the provided address/port.
func Serve(opts *opts.Options) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opts.Host, opts.Port))
	if err != nil {
		return nil, err
	}

	showOpts(opts)

	server := newServer(opts)
	server.report.StartTimestamp = time.Now()
	server.report.Charset = opts.Charset
	server.report.Offset = bigint.ToBigInt(opts.Offset)
	server.report.Length = opts.Passlen
	server.report.Jobsize = bigint.ToBigInt(opts.Jobsize)
	server.report.JobType = opts.Mode
	server.report.Target = opts.Target

	go server.accept(listener)
	go server.loop()

	found := make(chan bool)
	go waitForTarget(opts.Target, found)
	if opts.Verbose {
		term.Info("Waiting for target to become available...\n")
	}
	<-found
	if opts.Verbose {
		term.Info("Target available, let's go!\n")
	}

	go server.initCrackjobs(opts)
	go server.dispatch()
	go server.checkPending()
	go server.updateProgress()
	go server.checkErrors()

	if opts.EnableREST {
		if opts.Verbose {
			term.Info("Enabling REST server on %s:%d\n", opts.RESTAddress, opts.RESTPort)
		}
		api, err := rest.NewAPI(server, opts.RESTAddress, opts.RESTPort)
		if err != nil {
			term.Warn("Unable to start REST API: %s\n", err)
		} else {
			go api.Serve()
		}
	}

	term.Info("Server ready\n")

	return server, nil
}
