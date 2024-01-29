package handler

import "net/http"

func (hf handlerFranchises) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hf.respond(w, "OK", http.StatusOK)
	}
}
