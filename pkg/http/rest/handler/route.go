package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	apiVersion      = "/v1"
	healthCheckRoot = "/health_check"
	webScraping     = "/web-scraping"
)

func Register(r *mux.Router, handler *handlerFranchises, ctx context.Context) {
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix(apiVersion).Subrouter()
	apiV1.Use(handler.MiddlewareLogger())

	apiV1.HandleFunc("/web-scraping/{id}", handler.Get()).Methods(http.MethodGet)
	apiV1.HandleFunc("/web-scraping", handler.Create()).Methods(http.MethodPost)
	apiV1.HandleFunc("/web-scraping", handler.All()).Methods(http.MethodGet)

}
