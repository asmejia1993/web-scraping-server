package handler

import (
	"net/http"
)

func (hf handlerFranchises) All() http.HandlerFunc {

	type errorServer struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		res, err := hf.fService.All(r.Context(), q)
		if err != nil {
			hf.logger.Errorf("error from all() service with: %v", err)
			msg := errorServer{
				Message: "internal server error",
			}
			hf.respond(w, msg, http.StatusInternalServerError)
		}
		hf.respond(w, res, http.StatusOK)
	}
}
