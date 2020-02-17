package rest

import "net/http"

// The APIHandler interface defines - what is it called in golang? classes? - that are
// able to handle the REST API calls made to the server.
type APIHandler interface {
	// SlashRoot Handles calls to /
	SlashRoot(res http.ResponseWriter, req *http.Request)

	// ClientsHandler handles calls to /clients
	ClientsHandler(res http.ResponseWriter, req *http.Request)

	// PendingJobsHandler handles calls to /jobs
	PendingJobsHandler(res http.ResponseWriter, req *http.Request)

	// JobHandler handles calls to /job/{id}
	JobHandler(rest http.ResponseWriter, req *http.Request)

	// ClientHandler handles calls to /client/{id}
	ClientHandler(res http.ResponseWriter, req *http.Request)

	// KickClientHandler handles calls to /client/kick/{id}
	KickClientHandler(res http.ResponseWriter, req *http.Request)

	// TerminateHandler handles calls to /server/terminate
	TerminateHandler(res http.ResponseWriter, req *http.Request)

	// ReportHandler handles calls to /server/report
	ReportHandler(res http.ResponseWriter, req *http.Request)

	// HistoryHandler handles calls to the history endpoint.
	HistoryHandler(res http.ResponseWriter, req *http.Request)

	// ClientHistoryHandler handles calls to the history-per-client endpoint
	ClientHistoryHandler(res http.ResponseWriter, req *http.Request)

	// ConfigureHandler handles calls to the configuration endpoint that enables the
	// REST client to configure limited parts of the server config.
	ConfigureHandler(res http.ResponseWriter, req *http.Request)

	// ConfigHandler handles calls to the config endpoint that pulls the
	// server's configuration from the API
	ConfigHandler(res http.ResponseWriter, req *http.Request)

	// TargetHandler handles calls to the target endpoint, getting information on the
	// server's target.
	TargetHandler(res http.ResponseWriter, req *http.Request)
}
