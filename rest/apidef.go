package rest

import "net/http"

// The APIHandler interface defines - what is it called in golang? classes? - that are
// able to handle the REST API calls made to the server.
type APIHandler interface {
	// SlashRoot Handles calls to /
	SlashRoot(res http.ResponseWriter, req *http.Request)

	// ClientsHandleFunc handles calls to /clients
	ClientsHandleFunc(res http.ResponseWriter, req *http.Request)

	// PendingJobsHandleFunc handles calls to /jobs
	PendingJobsHandleFunc(res http.ResponseWriter, req *http.Request)
}
