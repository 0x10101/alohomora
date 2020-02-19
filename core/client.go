package core

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/steps0x29a/alohomora/bigint"
	"github.com/steps0x29a/alohomora/bytes"
	"github.com/steps0x29a/alohomora/ext"
	"github.com/steps0x29a/alohomora/fio"
	"github.com/steps0x29a/alohomora/gen"
	"github.com/steps0x29a/alohomora/handshakes"
	"github.com/steps0x29a/alohomora/jobs"
	"github.com/steps0x29a/alohomora/msg"
	"github.com/steps0x29a/alohomora/term"

	"github.com/steps0x29a/alohomora/opts"

	uuid "github.com/satori/go.uuid"
)

const (
	bufferSize uint32 = 4096
)

// A Client is basically a socket connection with some
// additional info.
type Client struct {
	Socket            net.Conn
	ID                uuid.UUID
	Terminated        chan bool
	Errors            chan error
	validated         chan bool
	connected         time.Time
	assigned          uint64
	finished          uint64
	tried             *big.Int
	handshakeFilepath string
}

// ClientInfo wraps information on a client for reporting.
type ClientInfo struct {
	ID         string   `json:"id"`
	Address    string   `json:"address"`
	Connected  string   `json:"connected"`
	Assigned   uint64   `json:"assigned"`
	Finished   uint64   `json:"finished"`
	Tried      *big.Int `json:"tried"`
	CurrentJob string   `json:"currentjob"`
}

// Info creates a ClientInfo instance from a Client pointer and returns it.
func (client *Client) Info() *ClientInfo {
	info := &ClientInfo{
		ID:         client.ShortID(),
		Address:    client.Socket.RemoteAddr().String(),
		Connected:  client.connected.String(),
		Assigned:   client.assigned,
		Finished:   client.finished,
		Tried:      client.tried,
		CurrentJob: "",
	}
	return info
}

// ShortID returns the first eight characters of a client's ID.
// Inspired by git's short commit hashes.
func (client *Client) ShortID() string {
	return client.ID.String()[:8]
}

// FullID returns a client's full ID consisting of its short ID (the
// first eight characters of its full ID) and its socket remote address.
func (client *Client) FullID() string {
	return fmt.Sprintf("%s | %s", client.ShortID(), client.Socket.RemoteAddr().String())
}

// String returns the same as client.FullID(), basically its short ID and socket
// remote address.
func (client *Client) String() string {
	return client.FullID()
}

func newClient(socket net.Conn) *Client {
	return &Client{
		Socket:     socket,
		ID:         uuid.Nil,
		Terminated: make(chan bool),
		Errors:     make(chan error),
		validated:  make(chan bool),
		assigned:   0,
		finished:   0,
		tried:      big.NewInt(0),
	}
}

func generatePasswords(params *jobs.GenerationParams) ([]string, error) {
	buffer := make([]string, params.Amount)
	term.Info("Generating %d passwords...", params.Amount)
	var i int64
	for i = 0; i < params.Amount; i++ {
		pw, err := gen.GeneratePassword(params.Charset, params.Length, bigint.Add(params.Offset, big.NewInt(i)))
		if err != nil {
			return nil, err
		}
		buffer[i] = pw
	}

	//term.Info("Generated %d passwords\n", params.Amount)
	fmt.Printf("%s\n", term.BrightGreen("OK"))
	return buffer, nil
}

func generatePasswordFile(params *jobs.GenerationParams) (string, error) {
	path, err := fio.TempFilePath()
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer fmt.Println()

	term.Info("Generating %d passwords...", params.Amount)

	var i int64
	for i = 0; i < params.Amount; i++ {
		pw, err := gen.GeneratePassword(params.Charset, params.Length, bigint.Add(params.Offset, big.NewInt(i)))
		if err != nil {
			return "", err
		}

		_, err = f.WriteString(fmt.Sprintf("%s\n", pw))
		if err != nil {
			return "", err
		}
	}

	//term.Info("Generated %d passwords\n", params.Amount)
	fmt.Printf(term.BrightGreen("OK"))
	return path, nil
}

func writeTmpBinFile(data []byte) (string, error) {
	path, err := fio.TempFilePath()
	if err != nil {
		return "", err
	}

	hs, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer hs.Close()

	_, err = hs.Write(data)
	if err != nil {
		return "", err
	}

	return path, nil

}

func (client *Client) work(job *jobs.CrackJob) ([]byte, bool, error) {

	// Decode crackjob' payload
	if job.Type == jobs.WPA2 {
		term.Info("Working on %s (%s)...\n", term.BrightBlue(job.ID.String()[:8]), term.Cyan("WPA2"))

		var handshake = handshakes.NewHandshake()

		if job.Payload != nil {
			term.Info("Payload data present, assuming new or unknown target\n")
			// Assume WPA2 payload, decode and save
			payload, err := job.DecodeWPA2()
			if err != nil {
				return nil, false, err
			}

			client.handshakeFilepath, err = writeTmpBinFile(payload.Data)
			if err != nil {
				return nil, false, err
			}

			err = handshake.Read(client.handshakeFilepath)
			if err != nil {
				return nil, false, err
			}
		} else {
			// Re-read the payload from file
			term.Info("No payload data, assuming same target as previous job\n")
			err := handshake.Read(client.handshakeFilepath)
			if err != nil {
				return nil, false, err
			}

		}

		// Generate passwords
		pwFilepath, err := generatePasswordFile(job.Gen)
		if err != nil {
			return nil, false, err
		}

		defer os.Remove(pwFilepath)
		//defer os.Remove(handshakeFilepath)

		term.Info("Running aircrack-ng...")
		output, err := ext.Aircrack(handshake.BSSID, handshake.ESSID, pwFilepath, client.handshakeFilepath)
		if err != nil {
			fmt.Printf("%s\n", term.BrightRed("ERROR"))
			return nil, false, err
		}
		fmt.Printf("%s\n", term.BrightGreen("OK"))

		password := ext.KeyFromOutput(output)
		found := password != ""
		if found {
			term.Success("%s\n", term.LabelGreen("Cracked the password!"))
		} else {
			term.Problem("%s\n", term.BrightRed("Too bad, password not cracked"))
		}
		fmt.Println()
		result := &jobs.CrackJobResult{Payload: password, JobID: job.ID, Success: found}

		enc, err := result.Encode()
		return enc, found, err
	} else if job.Type == jobs.MD5 {
		// Just tell server that we failed for now

		if job.Payload == nil {
			return nil, false, errors.New("Empty payloads not supported in MD5 mode")
		}

		decoded, err := job.DecodeMD5()
		if err != nil {
			return nil, false, err
		}

		passwords, err := generatePasswords(job.Gen)
		if err != nil {
			return nil, false, err
		}

		term.Info("Computing hashes...\n")
		for _, password := range passwords {

			salt := string(decoded.Salt)
			hash := string(decoded.Data)
			var candidates []string
			if len(salt) > 0 {
				candidates = []string{password, salt + password, password + salt, salt + "$" + password, password + "$" + salt}
			} else {
				candidates = []string{password}
			}

			for _, candidate := range candidates {

				h := md5.New()
				io.WriteString(h, candidate)
				digest := fmt.Sprintf("%x", h.Sum(nil))

				if digest == hash {
					term.Success("%s\n", term.LabelGreen("Cracked the hash!"))
					result := &jobs.CrackJobResult{Payload: candidate, JobID: job.ID, Success: true}
					enc, err := result.Encode()
					return enc, true, err
				}
			}

		}
		// Not cracked
		result := &jobs.CrackJobResult{Payload: "", JobID: job.ID, Success: false}
		enc, err := result.Encode()
		return enc, false, err
	}

	term.Warn("Only WPA2 jobs are implemented as of now\n")
	return nil, false, errors.New("Not a WPA2 payload")
}

func (client *Client) handle(message *msg.Message) {
	switch message.Type {
	case msg.Ack:
		{
			client.validated <- true
			break
		}
	case msg.Leave:
		{
			term.Info("Server asked me to leave, closing connection...\n")
			client.Terminated <- true
			break
		}

	case msg.Task:
		{

			// Decode payload
			job, err := jobs.DecodeJob(message.Payload)
			if err != nil {
				/// TODO: Send error message to server
				errMsg := msg.NewMessage(msg.ClientError, []byte(err.Error()))
				client.send(errMsg)
				client.Errors <- err
				return
			}

			term.Info("Received a new task: %s\n", term.BrightBlue(job.ID.String()[:8]))

			result, success, err := client.work(job)
			if err != nil {
				term.Error("Failed to attempt cracking: %s\n", err)
				result := jobs.CrackJobResult{JobID: job.ID, Payload: "", Success: false}
				encoded, err := result.Encode()
				if err != nil {
					term.Error("Unable to encode crackjobresult: %s\n", err)
					client.Terminated <- true
					break
				}
				term.Info("Sending fail message\n")
				client.send(msg.NewMessage(msg.ClientError, encoded))
				break
			}
			answer := msg.NewMessage(msg.Finished, result)
			client.send(answer)
			if success {
				client.Terminated <- true
			}
			break
		}
	}
}

func (client *Client) receive() {
	var buffer = make([]byte, bufferSize)

	for {
		var message = make([]byte, 0)
		var size = 0

		for {
			read, err := client.Socket.Read(buffer)
			if read == 0 || err != nil && err != io.EOF {
				// Connection lost
				client.Terminated <- true
				return
			}

			message = append(message, buffer[:read]...)
			size += read

			if bytes.EndsWith(message, AlohomoraSuffix) {
				decoded, err := msg.Decode(message[:size-len(AlohomoraSuffix)])
				if err != nil {
					term.Error("Unable to decode server message: %s\n", term.BrightRed(fmt.Sprintf("%s", err)))
				} else {
					client.handle(decoded)
				}

				break
			}

		}

	}
}

func (client *Client) snd(message *msg.Message) {
	data, err := message.Encode()
	if err != nil {
		term.Error("Unable to encode message: %s\n", err)
		client.Errors <- err
		return
	}

	writer := bufio.NewWriter(client.Socket)
	_, err = writer.Write(data)
	if err != nil {
		term.Error("Unable to send message to server: %s\n", err)
		client.Errors <- err
		return
	}

	_, err = writer.Write(AlohomoraSuffix)
	if err != nil {
		term.Error("Unable to send message to server: %s\n", err)
		client.Errors <- err
		return
	}
	fmt.Println(data)

}

func (client *Client) send(message *msg.Message) {
	data, err := message.Encode()
	if err != nil {
		term.Error("Unable to encode message: %s\n", err)
		client.Errors <- err
		return
	}

	_, err = client.Socket.Write(data)
	// TODO: Handle incomplete writes
	if err != nil {
		term.Error("Unable to send message: %s\n", err)
		client.Errors <- err
		return
	}

	_, err = client.Socket.Write(AlohomoraSuffix)
	if err != nil {
		term.Error("Unable to send suffix: %s\n", err)
		client.Errors <- err
		return
	}

}

// Shutdown cleans up the client's temp file(s)
func (client *Client) Shutdown() {
	// Delete payload file
	defer os.Remove(client.handshakeFilepath)
}

// Connect tries to establish a connection to a server.
// The server's IP and port are provided via an Options instance.
func Connect(opts *opts.Options) (*Client, error) {
	dialer := net.Dialer{
		Timeout: time.Second * 30,
	}

	term.Info("Connecting to %s:%d...\n", opts.Host, opts.Port)
	socket, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", opts.Host, opts.Port))
	if err != nil {
		return nil, err
	}

	client := newClient(socket)
	go client.receive()
	go client.send(msg.NewMessage(msg.Hello, nil))

	<-client.validated
	term.Success(term.BrightGreen("Connection established!\n"))

	// Tell server we are ready for action
	go client.send(msg.NewMessage(msg.Idle, nil))

	return client, nil
}
