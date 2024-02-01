package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (hf handlerFranchises) Get() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id := vars["id"]

		type response struct {
			Id      string `json:"id"`
			Message string `json:"message"`
		}
		sizeKey := fmt.Sprintf("%s.franchises.size", id)
		processedKey := fmt.Sprintf("%s.franchises.processed", id)
		exist, err := hf.fService.Exist(r.Context(), sizeKey)
		if err == nil && exist > 0 {
			val, _ := hf.fService.GetKey(r.Context(), sizeKey)
			size, _ := strconv.Atoi(val)
			val, _ = hf.fService.GetKey(r.Context(), processedKey)
			processed, _ := strconv.Atoi(val)
			if processed < size {
				res := response{
					Id: id, Message: "scraping process still running in background",
				}
				hf.respond(w, res, http.StatusOK)
			}

		}

		res := hf.fService.Get(id, r.Context())

		hf.respond(w, res, http.StatusOK)
	}
}
