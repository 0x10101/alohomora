package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// API wraps required information on the REST API in a convenient struct type
type API struct {
	router  *mux.Router
	address string
	port    uint
}

// NewAPI creates a new RestAPI object
func NewAPI(handler APIHandler, address string, port uint) (*API, error) {
	if handler == nil {
		return nil, errors.New("An API handler is required")
	}
	api := new(API)
	api.address = address
	api.port = port
	api.router = mux.NewRouter().StrictSlash(true)
	api.router.HandleFunc("/", handler.SlashRoot)
	api.router.HandleFunc("/clients", handler.ClientsHandler)
	api.router.HandleFunc("/clients/{id}", handler.ClientHandler)
	api.router.HandleFunc("/clients/{id}/history", handler.ClientHistoryHandler)
	api.router.HandleFunc("/jobs", handler.PendingJobsHandler)
	api.router.HandleFunc("/jobs/{id}", handler.JobHandler)
	api.router.HandleFunc("/clients/kick/{id}", handler.KickClientHandler).Methods("POST")
	api.router.HandleFunc("/server/terminate", handler.TerminateHandler).Methods("POST")
	api.router.HandleFunc("/report", handler.ReportHandler)
	api.router.HandleFunc("/history", handler.HistoryHandler)
	api.router.HandleFunc("/configure", handler.ConfigHandler).Methods("POST")

	return api, nil
}

// Serve starts the REST API
func (api *API) Serve() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", api.address, api.port), api.router))
}
