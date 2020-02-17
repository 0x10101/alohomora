package core

/*
This file contains all functions required in order to implement the REST API provided by alohomora.
It was externalized because it started to mess up the main server.go file.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/steps0x29a/alohomora/handshakes"

	"github.com/steps0x29a/alohomora/bigint"

	"github.com/steps0x29a/alohomora/opts"

	"github.com/gorilla/mux"
	"github.com/steps0x29a/alohomora/jobs"
	"github.com/steps0x29a/alohomora/term"
)

// SlashRoot is the root endpoint of the REST API (Server implements the RESTHandler interface).
func (server *Server) SlashRoot(res http.ResponseWriter, req *http.Request) {
	respondWithJSON(res, http.StatusOK, "alohomora ready")
}

// ClientsHandler wraps all connected clients in ClientInfo objects and marshals that
// information to JSON. Server implements RESTHandler interface.
func (server *Server) ClientsHandler(res http.ResponseWriter, req *http.Request) {
	clients := make([]*ClientInfo, 0)
	for client, connected := range server.Clients {
		if connected {
			info := client.Info()
			// Get current job
			job := server.Pending[client]

			if job != nil {
				info.CurrentJob = job.ShortID()
			}
			clients = append(clients, info)
		}
	}

	respondWithJSON(res, http.StatusOK, clients)
}

// JobHandler handles REST calls to /job/{id} and returns a JSON representation of a matching job.
func (server *Server) JobHandler(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	if server.verbose {
		term.Info("REST client requested info on job %s\n", id)
	}

	job := server.findJob(id)
	if job != nil {
		respondWithJSON(res, http.StatusOK, job)
	} else {
		respondWithError(res, http.StatusNotFound, "Unknown job ID")
	}
}

func (server *Server) ClientHistoryHandler(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	if server.verbose {
		term.Info("REST client requested client history on client %s\n", id)
	}

	client := server.findClient(id)
	if client != nil {
		data := server.history[client]
		if data != nil {
			respondWithJSON(res, http.StatusOK, data)
		} else {
			respondWithError(res, http.StatusNotFound, "Unknown data")
		}
	} else {
		respondWithError(res, http.StatusNotFound, "Unknown client ID")
	}

}

func (server *Server) HistoryHandler(res http.ResponseWriter, req *http.Request) {

	if server.verbose {
		term.Info("REST client requested client history on all clients\n")
	}

	respondWithJSON(res, http.StatusOK, server.history)
}

// KickAllHandler is not directly called by REST calls, but from the KickClientHandler function.
// If the REST client wants to kick client 'ffffff', all clients are kicked.
func (server *Server) KickAllHandler(res http.ResponseWriter, req *http.Request) {
	server.KickAll()
	respondWithJSON(res, http.StatusOK, "OK")
}

func (server *Server) TargetHandler(res http.ResponseWriter, req *http.Request) {
	if server.verbose {
		term.Info("REST client requested target\n")
	}

	essid, bssid, err := handshakes.HandshakeInfo(server.opts.Target)
	if err != nil {
		term.Warn("Unable to parse %s: %s\n", server.opts.Target, err)
		respondWithError(res, http.StatusInternalServerError, "Unable to parse target file")
		return
	}

	obj := struct {
		Filename string `json:"filename"`
		BSSID    string `json:"bssid"`
		ESSID    string `json:"essid"`
	}{
		server.opts.Target,
		bssid,
		essid,
	}

	fmt.Println(obj)

	respondWithJSON(res, http.StatusOK, obj)
}

func (server *Server) ConfigHandler(res http.ResponseWriter, req *http.Request) {
	if server.verbose {
		term.Info("REST client requested server config\n")
	}

	respondWithJSON(res, http.StatusOK, server.opts)
}

func (server *Server) ConfigureHandler(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		respondWithError(res, http.StatusUnprocessableEntity, "Sorry")
		return
	}

	if err := req.Body.Close(); err != nil {
		respondWithError(res, http.StatusUnprocessableEntity, "Sorry")
		return
	}

	opts := opts.Options{}
	opts.Verbose = server.verbose
	opts.MaxTime = server.timeout
	opts.MaxJobs = server.maxjobs.String()
	opts.Timeout = server.taskTimeout
	if err := json.Unmarshal(body, &opts); err != nil {
		respondWithError(res, http.StatusUnprocessableEntity, "Invalid")
		return
	}

	// Apply new options to server
	server.verbose = opts.Verbose
	server.timeout = opts.MaxTime
	server.maxjobs = bigint.ToBigInt(opts.MaxJobs)
	server.taskTimeout = opts.Timeout

	fmt.Println(string(body))
	fmt.Println("**************")
	fmt.Println(opts)

	respondWithJSON(res, http.StatusOK, "OK")

}

// KickClientHandler finds and kicks a client.
func (server *Server) KickClientHandler(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	if id == "ffffff" {
		term.Warn("Client requested kicking all clients\n")
		server.KickAllHandler(res, req)
		return
	}

	if server.verbose {
		term.Info("REST client requested kicking client %s\n", id)
	}

	client := server.findClient(id)
	if client != nil {
		server.kick(client)
		respondWithJSON(res, http.StatusOK, "OK")
	} else {
		respondWithError(res, http.StatusNotFound, "Unknown client ID")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		term.Error("Unable to send JSON response to REST client: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "Sorry")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ClientHandler wraps a single client's information in a ClientInfo object and
// marshals that information to JSON. Server implements RESTHandler interface.
func (server *Server) ClientHandler(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	if server.verbose {
		term.Info("REST client requested info on client %s\n", id)
	}

	client := server.findClient(id)
	if client != nil {
		info := client.Info()
		job := server.Pending[client]
		if job != nil {
			info.CurrentJob = job.ShortID()
		}
		respondWithJSON(res, http.StatusOK, info)
	} else {
		respondWithError(res, http.StatusNotFound, "Unknown client ID")
	}
}

// TerminateHandler handles REST calls to /server/terminate and terminates the server immediately.
func (server *Server) TerminateHandler(res http.ResponseWriter, req *http.Request) {
	server.Terminate()
	report := server.Report()
	respondWithJSON(res, http.StatusOK, report)
}

// PendingJobsHandler wraps all pending jobs in JobInfo objects and marshals that info to JSON.
// Server implements RESTHandler interface.
func (server *Server) PendingJobsHandler(res http.ResponseWriter, req *http.Request) {
	mapping := make(map[string]*jobs.CrackJobInfo)
	for client, job := range server.Pending {
		mapping[client.ShortID()] = job.Info()
	}

	respondWithJSON(res, http.StatusOK, mapping)
}

// ReportHandler handles REST calls to /server/report and sends a JSONified report without an
// end timestamp to the calling client.
func (server *Server) ReportHandler(res http.ResponseWriter, req *http.Request) {
	report := server.Report()
	// Remove stopped time, has not yet been stopped.
	report.EndTimestamp = time.Time{}
	respondWithJSON(res, http.StatusOK, report)
}
