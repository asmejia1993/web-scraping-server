package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (hf handlerFranchises) Get() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id := vars["id"]

		res := hf.fService.Get(id, r.Context())

		hf.respond(w, res, http.StatusOK)
	}
}
