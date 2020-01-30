package core

import (
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/steps0x29a/alohomora/gen"
	"github.com/steps0x29a/alohomora/handshakes"

	"github.com/steps0x29a/alohomora/msg"

	uuid "github.com/satori/go.uuid"

	"github.com/steps0x29a/alohomora/opts"
	"github.com/steps0x29a/islazy/bigint"
	"github.com/steps0x29a/islazy/bytes"
	"github.com/steps0x29a/islazy/fio"
	"github.com/steps0x29a/islazy/term"
)

const (
	bufferSize uint32 = 4096
)

// A Server manages clients and jobs.
type Server struct {
	sync.Mutex

	Clients            map[*Client]bool
	StartedJobs        *big.Int
	TotalJobs          *big.Int
	FinishedJobs       *big.Int
	Queue              chan *CrackJob
	freeClients        chan *Client
	Terminated         chan bool
	register           chan *Client
	unregister         chan *Client
	Errors             chan error
	Pending            map[*Client]*CrackJob
	generationFinished bool
	maximumJobsReached bool
	verbose            bool
	timeout            uint64
	Report             string
	queuesize          uint64
	maxjobs            *big.Int
	taskTimeout        uint64
	started            time.Time
	ESSID              string
	BSSID              string
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
		Queue:              make(chan *CrackJob, opts.QueueSize),
		freeClients:        make(chan *Client),
		Terminated:         make(chan bool),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		Errors:             make(chan error),
		Pending:            make(map[*Client]*CrackJob),
		generationFinished: false,
		maximumJobsReached: false,
		verbose:            opts.Verbose,
		timeout:            opts.Timeout,
		maxjobs:            bigint.ToBigInt(opts.MaxJobs),
		taskTimeout:        opts.MaxTime,
		started:            time.Now(),
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

func (server *Server) onNewClient(client *Client) {
	term.Info("Client connected: %s\n", term.BrightBlue(client.FullID()))
}

func (server *Server) onClientLeft(client *Client) {
	term.Info("Client left: %s\n", term.Red(client.FullID()))
}

func (server *Server) onProgress() {
	total := server.TotalJobs
	finished := term.Reverse(term.InsertAfterEvery(term.Reverse(server.FinishedJobs.String()), '.', 3))
	pending := len(server.Pending)
	clients := len(server.Clients)
	queued := len(server.Queue)

	percent := bigint.Percent(total, server.FinishedJobs)
	numStrTotal := term.Reverse(term.InsertAfterEvery(term.Reverse(total.String()), '.', 3))
	if server.verbose {
		term.Info("Progress: %s/%s (%0.2f%%, %d clients connected, %d jobs pending, %d queued)\n", finished, numStrTotal, percent, clients, pending, queued)
	}

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
				server.onNewClient(client)
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
		client := Client{Socket: connection, ID: clientID}
		server.register <- &client
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
		term.Warn("Client not validated: %s\n", client)
		term.Warn("Closing connection to %s\n", client)
		server.kick(client)
	} else {
		// Schedule client for work
		server.freeClients <- client
	}
}

func (server *Server) onClientResponse(client *Client, message *msg.Message) {
	server.Lock()

	result, err := decodeResult(message.Payload)
	delete(server.Pending, client)
	server.FinishedJobs = bigint.Add(server.FinishedJobs, big.NewInt(1))
	server.Unlock()

	if err != nil {
		server.Errors <- err
		term.Error("Unable to decode result: %s\n", err)
	} else {
		if result.Success {
			term.Success("Client %s cracked the password: %s\n", term.BrightBlue(client.ShortID()), term.LabelGreen(result.Payload))

			// Write report
			server.onClientSuccess(client, result.Payload)

			server.Terminated <- true
		} else {
			term.Error("Client %s %s to crack %s\n", term.BrightBlue(client.ShortID()), term.BrightRed("failed"), term.Cyan(result.JobID.String()[:8]))
		}
	}

	server.onClientIdle(client, message)
}

func (server *Server) onClientSuccess(client *Client, password string) {
	fio.WriteTo(server.Report, password)
}

func (server *Server) onClientError(client *Client, message *msg.Message) {
	// Payload should be CrackJobResult

	result, err := decodeResult(message.Payload)
	defer server.kick(client)

	if err != nil {
		server.Errors <- err
		term.Error("Client %s crashed with invalid error message\n", term.BrightBlue(client.ShortID()))
	} else {
		term.Error("Client %s crashed during %s\n", term.BrightBlue(client.ShortID()), term.Cyan(result.JobID.String()[:8]))
	}
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
		job, _ := <-server.Queue

		payload, err := job.encode()
		if err != nil {
			server.Terminated <- true
			return
		}

		client, _ := <-server.freeClients

		job.Started = time.Now()
		server.Lock()
		server.Pending[client] = job
		message := msg.NewMessage(msg.Task, payload)
		server.Unlock()
		go server.send(client, message)
	}
}

func (server *Server) initCrackjobs(opts *opts.Options) {
	var filepath = opts.Target
	var handshake = handshakes.NewHandshake()
	err := handshake.Read(filepath)
	if err != nil {
		// This is bad
		term.Error("Unable to process target: %s\n", err)
		server.Terminated <- true
		return
	}

	charset := []rune(opts.Charset)
	length := int64(opts.Passlen)
	jobsize := bigint.ToBigInt(opts.Jobsize)
	offset := bigint.ToBigInt(opts.Offset)
	maxValue := bigint.Sub(bigint.Pow(big.NewInt(int64(len(charset))), big.NewInt(length)), offset)
	runs := bigint.Div(maxValue, jobsize)

	mod := bigint.Mod(maxValue, jobsize)
	if !bigint.SameAs(mod, big.NewInt(0)) {
		runs = bigint.Add(runs, big.NewInt(1))
	}

	server.TotalJobs = bigint.Copy(runs)

	var jobindex *big.Int = big.NewInt(0)
	var remaining *big.Int = bigint.Copy(maxValue)

	if bigint.LessThan(remaining, big.NewInt(0)) {
		term.Error("Invalid offset: %s\n", offset)
		server.Terminated <- true
		return
	}

	for bigint.GreaterThan(remaining, big.NewInt(0)) {

		var runAmount *big.Int = bigint.Copy(jobsize)
		if bigint.LessThan(remaining, jobsize) {
			runAmount.Set(remaining)
		}
		remaining = bigint.Sub(remaining, runAmount)

		var calcOffset = bigint.Add(offset, bigint.Mul(jobsize, jobindex))
		var endOffset = bigint.Sub(bigint.Add(calcOffset, runAmount), big.NewInt(1))

		first, err := gen.GeneratePassword(charset, length, calcOffset)
		if err != nil && server.verbose {
			term.Problem("Unable to preview first password for job: %s\n", err)
		}

		last, err := gen.GeneratePassword(charset, length, endOffset)
		if err != nil && server.verbose {
			term.Problem("Unable to preview last password for job: %s\n", err)
		}

		job, err := newWPA2Job(
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
			server.Terminated <- true
			return
		}

		jobindex = bigint.Add(jobindex, big.NewInt(1))
		if server.verbose {
			term.Success("Generated Crackjob %s (%s - %s)\n", term.Cyan(job.String()), term.BrightGreen(first), term.BrightGreen(last))
		}

		server.Queue <- job
		if bigint.GTE(jobindex, server.maxjobs) && bigint.GreaterThan(server.maxjobs, big.NewInt(0)) {
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
		if server.generationFinished && bigint.GTE(server.FinishedJobs, server.TotalJobs) {
			server.onProgress()
			term.Info("All jobs finished, terminating...\n")
			server.Terminated <- true
		} else if server.maximumJobsReached && bigint.GTE(server.FinishedJobs, server.maxjobs) {
			server.onProgress()
			term.Info("Maximum number of jobs reached, terminating...\n")
			server.Terminated <- true
		} else {
			server.onProgress()
		}

		// Check timeout
		if server.taskTimeoutReached() {
			// Task timed out, kill everything
			term.Info("Task timeout reached, terminating server...\n")
			server.Terminated <- true
		}
	}
}

func (server *Server) taskTimeoutReached() bool {
	return server.taskTimeout > 0 && time.Now().Sub(server.started).Seconds() > float64(server.taskTimeout)
}

// Serve builds a new Server instance and starts listening on the provided address/port.
func Serve(opts *opts.Options) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opts.Host, opts.Port))
	if err != nil {
		return nil, err
	}

	showOpts(opts)

	server := newServer(opts)

	go server.accept(listener)
	go server.loop()
	go server.initCrackjobs(opts)
	go server.dispatch()
	go server.checkPending()
	go server.updateProgress()

	if server.verbose {
		term.Info("Alohomora Server ready, waiting for clients...\n")
	}

	return server, nil
}
