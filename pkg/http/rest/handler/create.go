package handler

import (
	"net/http"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

func (hf handlerFranchises) Create() http.HandlerFunc {

	type inserted struct {
		ID  string `json:"id"`
		URL string `json:"link"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req model.FranchiseInfo

		err := hf.decode(r, &req)
		if err != nil {
			hf.respond(w, err, http.StatusInternalServerError)
			return
		}
		res, err := hf.fService.Create(r.Context(), req)
		response := inserted{
			ID: res,
		}
		hf.logger.Infof("response: %v", response)
		if err != nil {
			hf.respond(w, err, http.StatusInternalServerError)
			return
		}
		//call scraper fn
		hf.sc.InitScraping(req.Company.Franchises)
		hf.respond(w, response, http.StatusCreated)
	}
}
