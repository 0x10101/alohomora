package rest

import "net/http"

type RestHandler interface {
	SlashRoot(res http.ResponseWriter, req *http.Request)
	ClientsHandleFunc(res http.ResponseWriter, req *http.Request)
	PendingJobsHandleFunc(res http.ResponseWriter, req *http.Request)
}
