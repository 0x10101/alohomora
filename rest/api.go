package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RestAPI struct {
	router  *mux.Router
	address string
	port    uint16
}

func NewRestAPI(handler RestHandler, address string, port uint16) *RestAPI {
	api := new(RestAPI)
	api.address = address
	api.port = port
	api.router = mux.NewRouter().StrictSlash(true)
	api.router.HandleFunc("/", handler.SlashRoot)
	api.router.HandleFunc("/clients", handler.ClientsHandleFunc)
	api.router.HandleFunc("/jobs", handler.PendingJobsHandleFunc)

	return api
}

func (api *RestAPI) Serve() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", api.address, api.port), api.router))
}
