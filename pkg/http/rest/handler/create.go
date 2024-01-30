package handler

import (
	"context"
	"net/http"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

func (hf handlerFranchises) Create() http.HandlerFunc {

	type inserted struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req model.FranchiseInfoReq

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

		//Queue new task
		hf.worker.Start(r.Context(), len(req.Company.Franchises))
		for _, v := range req.Company.Franchises {
			scrapReq := model.FranchiseScraper{
				Id:        response.ID,
				Franchise: v,
			}
			hf.worker.QueueTask(scrapReq)
		}

		hf.receiveResultFromWorker(r.Context())
		hf.respond(w, response, http.StatusCreated)
	}
}

func (hf handlerFranchises) receiveResultFromWorker(ctx context.Context) {
	for res := range hf.resultChan {
		select {
		case <-ctx.Done():
			return
		default:
			hf.logger.Infof("Received result from worker: %v", res)
		}
	}
}
