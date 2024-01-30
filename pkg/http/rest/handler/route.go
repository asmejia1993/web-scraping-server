package handler

import (
	"net/http"

	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/worker"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	apiVersion      = "/v1"
	healthCheckRoot = "/health_check"
	webScraping     = "/web-scraping"
)

func Register(r *mux.Router, lg *logrus.Logger, db *config.DBInfo, w worker.IWorker) {
	handler := newHandler(lg, db, w)
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix(apiVersion).Subrouter()
	apiV1.Use(handler.MiddlewareLogger())

	apiV1.HandleFunc("/web-scraping/{id}", handler.Get()).Methods(http.MethodGet)
	apiV1.HandleFunc("/web-scraping", handler.Create()).Methods(http.MethodPost)

}
